package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetConnID() uint32
	GetTCPConnection() *net.TCPConn
	SendMsg(msgId uint32, data []byte) error
}
