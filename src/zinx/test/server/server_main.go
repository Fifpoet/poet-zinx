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

func (p *PingRouter) Handler(request ziface.IRequest) {
	fmt.Println("[INFO] Handler run! ")
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()

}
