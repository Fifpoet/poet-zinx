package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/src/zinx/ziface"
)

// Server 封装服务器信息，起循环监听特定地址
type Server struct {
	// 服务器名称
	Name string
	// IP协议的版本 IP地址 端口
	IPVersion  string
	IP         string
	Port       int
	MsgHandler ziface.IMsgHandler
}

// NewServer 更新硬编码为从json文件读取
func NewServer() ziface.IServer { // TODO 使用接口还是Server
	//utils.GlobalConfig.ReloadConfig()
	//conf := utils.GlobalConfig
	//s := &Server{
	//	Name:      conf.Name,
	//	IPVersion: conf.Version,
	//	IP:        conf.Host,
	//	Port:      conf.TcpPort,
	//	Router:    nil,
	//}
	s := &Server{
		Name:       "FIF",
		IPVersion:  "tcp4",
		IP:         "0.0.0.0",
		Port:       7777,
		MsgHandler: NewMsgHandle(),
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	go func() {
		// 0 启动WorkerQueue
		s.MsgHandler.StartWorkPoll()
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
		fmt.Println("[INFO] Start ZINX Server successfully!" + "server name: {" + s.Name + "}")
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
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++ //TODO 线程安全？？
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] Server stopped, name:{" + s.Name + "}")
	// TODO 释放连接资源
}

func (s *Server) Serve() {
	s.Start()
	// 启动后处理

	// TODO 阻塞,否则主Go退出， listener的go将会退出. 忙循环的理由？
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}
