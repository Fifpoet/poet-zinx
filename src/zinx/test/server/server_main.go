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

func (p *PingRouter) PreHandler(request ziface.IRequest) {
	// 回写到客户端显示
	fmt.Println("[INFO] PreHandler run! ")
}

func (p *PingRouter) Handler(request ziface.IRequest) {
	fmt.Println("[INFO] Handler run! ")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("[Error] Write back error")
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()

}
