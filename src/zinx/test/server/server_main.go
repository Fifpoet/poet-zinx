package main

import (
	"fmt"
	"zinx/src/zinx/ziface"
	"zinx/src/zinx/znet"
)

type PingRouter struct {
	// 采用匿名组合的方式，可以不用重写三个handler
	znet.BaseRouter
}

type HelloRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handler(request ziface.IRequest) {
	fmt.Println("[INFO] PingHandler run! ")
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func (h *HelloRouter) Handler(request ziface.IRequest) {
	fmt.Println("[INFO] HelloHandler run! ")
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("hello, hello, hhh"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})
	s.AddRouter(3, &HelloRouter{})
	s.AddRouter(4, &HelloRouter{})
	s.AddRouter(5, &HelloRouter{})
	s.AddRouter(6, &HelloRouter{})
	s.AddRouter(7, &HelloRouter{})
	s.AddRouter(8, &HelloRouter{})
	s.AddRouter(9, &HelloRouter{})
	//s.AddRouter(10, &HelloRouter{})
	//s.AddRouter(11, &HelloRouter{})
	s.Serve()

}
