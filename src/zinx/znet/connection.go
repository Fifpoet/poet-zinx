package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

// Connection 封装每一个连接，绑定对应的业务逻辑
type Connection struct {
	// conn可以感知隶属的server对象
	TcpServer ziface.IServer
	// TCP连接的原始套接字
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	MsgHandler ziface.IMsgHandler
	// 在New函数中初始化为1 chan用于阻塞
	ExitBufChan chan bool
	MsgChan     chan []byte
	// 有缓冲管道
	MsgBuffChan chan []byte
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:   server,
		isClosed:    false,
		MsgHandler:  handler,
		Conn:        conn,
		ConnID:      connID,
		ExitBufChan: make(chan bool, 1),
		MsgChan:     make(chan []byte, 1),
		MsgBuffChan: make(chan []byte),
	}
	c.TcpServer.GetConnManager().Add(c) // 获得server并add自己
	return c
}

func (c *Connection) StartWrite() {
	fmt.Println("[INFO] Write Goroutine is running")
	for {
		select {
		case data := <-c.MsgChan:
			// 开始读取
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("send data error!")
				return
			}
		case data, ok := <-c.MsgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				break
			}
		case <-c.ExitBufChan:
			return
		}
	}
}

// StartReader 封装读写分离
func (c *Connection) StartReader() {
	fmt.Println("[INFO] Reader Goroutine is running")
	defer fmt.Println("[INFO] Reader Closed")
	defer c.Stop()

	//由读取字节流 改成
	for {
		dp := NewDataPack()
		// 1. 读Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err) //TODO 读取时没有相应的router会阻塞
			c.ExitBufChan <- true                    // TODO 给到退出消息后，主线程不能立马select到，这里会多次执行
			continue
		}
		// 2. 先解包Head
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("[ERROR] unpack error ", err)
			c.ExitBufChan <- true
			continue
		}
		// 3. 按照DataPack的形式按长度读取data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("[ERROR] read msg data error ", err)
				c.ExitBufChan <- true
				continue
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}
		// 注意发送到taskQueue中不用起go
		if utils.GlobalConfig.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// SendMsg 提供封包方法 快捷的把发送的[]byte转换为msg.
// v0.7 读写分离 写到msgChan即可
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msgBytes, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("[ERROR] Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	c.MsgChan <- msgBytes
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msgBytes, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("[ERROR] Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}
	c.MsgBuffChan <- msgBytes
	return nil
}

func (c *Connection) Start() {
	go c.StartReader()
	go c.StartWrite()
	for {
		select {
		case <-c.ExitBufChan:
			//在chan中得到了退出的消息（true）程序退出
			fmt.Println("[INFO] Receive ExitBuf")
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	err := c.Conn.Close()
	if err != nil {
		return
	}
	c.ExitBufChan <- true
	// 把conn从管理器中删除
	c.TcpServer.GetConnManager().Remove(c)
	// 释放连接中的chan
	close(c.ExitBufChan)
	close(c.MsgBuffChan)
	close(c.MsgChan)
}

// GetConnID 获取封装的Connection对象ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// GetTCPConnection 通过Connection获取原始TCP连接对象
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetRemoteAddr 获取远程地址 ip:port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
