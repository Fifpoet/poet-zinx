package ziface

type IRequest interface {
	GetConnection() iConnection
	GetData() []byte
}
