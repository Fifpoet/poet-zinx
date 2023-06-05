package znet

import (
	"fmt"
	"net"
	"zinx/src/zinx/ziface"
)

// Connection 封装每一个连接，绑定对应的业务逻辑
type Connection struct {
	// TCP连接的原始套接字
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	Router   ziface.IRouter
	// 在New函数中初始化为1 chan用于阻塞
	ExitBufChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		isClosed:    false,
		Router:      router,
		Conn:        conn,
		ConnID:      connID,
		ExitBufChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Printf("[INFO] Reader Gorountine is running")
	defer fmt.Printf("[INFO] Reader Closed")
	defer c.Stop()

	//读取客户端字节流
	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			c.ExitBufChan <- true
			fmt.Printf("[Error] Read bytes error")
			continue
		}
		req := Request{
			conn: c,
			data: buf,
		}
		go func(request ziface.IRequest) {
			c.Router.PreHandler(request)
			c.Router.Handler(request)
			c.Router.PostHandler(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	for {
		select {
		case <-c.ExitBufChan:
			//在chan中得到了退出的消息（一个bool值）程序退出
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
