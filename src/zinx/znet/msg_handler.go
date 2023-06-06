package znet

import (
	"fmt"
	"strconv"
	"zinx/src/zinx/ziface"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性
}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
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
