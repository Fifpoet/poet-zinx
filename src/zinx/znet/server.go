package znet

import (
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

		// 3 忙循环接受TCP连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[Error] Accept TCP Connection Error")
				continue
			}
			// TODO 超过TCP连接数则关闭新连接
			// 协程嵌套
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("[Error] Receive bytes Error")
						continue
					}
					// 接受到字节后 回显到客户端. 此时的buf中只读到了前cnt个字节
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("[Error] Write back buf err ", err)
						continue
					}
				}
			}()
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
