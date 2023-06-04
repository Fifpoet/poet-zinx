package znet

import (
	"errors"
	"fmt"
	"net"
	"time"
	"zinx/src/zinx/ziface"
)

type Server struct {
	// 服务器名称
	Name string
	// IP协议的版本 IP地址 端口
	IPVersion string
	IP        string
	Port      int
}

func CallBackFunc(conn *net.TCPConn, data []byte, cnt int) error {
	// 封装回显业务逻辑
	fmt.Println("[INFO] CallBackFunc running")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("[ERROR] Write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s Server) Start() {
	fmt.Printf("[Start] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		// 1 获取TCP连接对象
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[Error] resolve tcp addr err: ", err)
			return
		}
		// 2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("[Info] Start ZINX Server successfully!" + "server name: {" + s.Name + "}")
		var cid uint32
		cid = 0
		// 3 忙循环接受TCP连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[Error] Accept TCP Connection Error")
				continue
			}
			// TODO 超过TCP连接数则关闭新连接
			dealConn := NewConnection(conn, cid, CallBackFunc)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s Server) Stop() {
	fmt.Println("[Stop] Server stopped, name:{" + s.Name + "}")
	// TODO 释放连接资源
}

func (s Server) Serve() {
	s.Start()
	// 启动后处理

	// TODO 阻塞,否则主Go退出， listener的go将会退出. 忙循环的理由？
	for {
		time.Sleep(10 * time.Second)
	}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
	return s
}
