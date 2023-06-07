package ziface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)          //非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //消息处理逻辑
	StartWorkPoll()                         //开启工作池
	SendMsgToTaskQueue(request IRequest)    //msg交给worker处理
}
