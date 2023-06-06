package ziface

type IDataPack interface {
	GetHeadLen() uint32                    //获取包头长度
	Pack(message IMessage) ([]byte, error) //包装
	Unpack([]byte) (IMessage, error)       //解包
}
