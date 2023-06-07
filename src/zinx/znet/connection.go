package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/src/zinx/ziface"
)

// Connection 封装每一个连接，绑定对应的业务逻辑
type Connection struct {
	// TCP连接的原始套接字
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	MsgHandler ziface.IMsgHandler
	// 在New函数中初始化为1 chan用于阻塞
	ExitBufChan chan bool
	MsgChan     chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		isClosed:    false,
		MsgHandler:  handler,
		Conn:        conn,
		ConnID:      connID,
		ExitBufChan: make(chan bool, 1),
		MsgChan:     make(chan []byte, 1),
	}
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
			fmt.Println("read msg head error ", err)
			c.ExitBufChan <- true // TODO 给到退出消息后，主线程不能立马select到，这里会多次执行
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
		go c.MsgHandler.DoMsgHandler(&req)
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

	//回写到客户端
	c.MsgChan <- msgBytes
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
	// 释放连接中的chan
	close(c.ExitBufChan)
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
