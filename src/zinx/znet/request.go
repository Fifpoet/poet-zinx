package znet

import "zinx/src/zinx/ziface"

type Request struct {
	conn ziface.IConnection
	data []byte
}

func (r *Request) GetConnection() interface{} {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
