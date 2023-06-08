package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	dp := &DataPack{}
	return dp
}

func (dp *DataPack) GetHeadLen() uint32 {
	// Id -> 4B   len -> 4B
	return 8
}

// Pack msg -> msgId + dataLen + data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 缓冲区提升拼接速度
	dataBuff := bytes.NewBuffer([]byte{})
	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(data)
	msg := &Message{}
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.GlobalConfig.MaxPacketSize > 0 && msg.DataLen > utils.GlobalConfig.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}
	return msg, nil
}
