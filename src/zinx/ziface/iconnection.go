package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetConnID() uint32
	GetTCPConnection() *net.TCPConn
	SendMsg(msgId uint32, data []byte) error // 采用无缓冲管道发送消息，可能导致短暂阻塞
	SendBuffMsg(msgId uint32, data []byte) error
}
