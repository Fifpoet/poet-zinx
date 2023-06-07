package ziface

type IServer interface {
	Start()
	Stop()
	// Serve 开启业务服务方法
	Serve()
	// AddRouter 为一个服务注册业务路由
	AddRouter(msgId uint32, router IRouter)
}
