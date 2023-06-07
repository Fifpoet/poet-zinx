package ziface

// IRouter 封装基础路由
type IRouter interface {
	PreHandler(request IRequest)
	Handler(request IRequest)
	PostHandler(request IRequest)
}
