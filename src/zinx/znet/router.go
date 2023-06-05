package znet

import (
	"zinx/src/zinx/ziface"
)

// BaseRouter 空struct，可以不初始化pre和post方法
type BaseRouter struct {
}

func (b *BaseRouter) PreHandler(request ziface.IRequest) {
}

func (b *BaseRouter) Handler(request ziface.IRequest) {
}

func (b *BaseRouter) PostHandler(request ziface.IRequest) {
}
