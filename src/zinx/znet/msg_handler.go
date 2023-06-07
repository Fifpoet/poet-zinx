package znet

import (
	"fmt"
	"strconv"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32                    //业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest    //Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: 5, // TODO 配置读取
		TaskQueue:      make([]chan ziface.IRequest, 10),
	}
}

// DoMsgHandler 非阻塞的方式运行router中的handler
func (m *MsgHandler) DoMsgHandler(req ziface.IRequest) {
	router, ok := m.Apis[req.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", req.GetMsgId(), " is not FOUND!")
		return
	}
	router.PreHandler(req)
	router.Handler(req)
	router.PostHandler(req)
}

func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api , msgId =" + strconv.Itoa(int(msgId)))
	}
	m.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

// StartOneWorker 用workerId对应taskQueue中的chan数组索引
func (m *MsgHandler) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkId = {", workerId, "}"+" is started")
	for {
		select {
		case req := <-taskQueue:
			m.DoMsgHandler(req)
		}
	}
}

func (m *MsgHandler) StartWorkPoll() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalConfig.MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandler) SendMsgToTaskQueue(req ziface.IRequest) {
	//TODO implement me
	panic("implement me")
}
