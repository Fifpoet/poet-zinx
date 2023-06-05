package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetConnID() uint32
}

// HandFunc 统一处理业务链接
type HandFunc func(*net.TCPConn, []byte, int) error