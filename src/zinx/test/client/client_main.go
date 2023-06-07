package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/src/zinx/znet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for i := 1; i < 15; i++ {
		var msg []byte
		dp := znet.NewDataPack()
		msg, _ = dp.Pack(znet.NewMessage(uint32(i), []byte(fmt.Sprintf("Hello, %d", i))))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		// 客户端接受回写的数据 并打印
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			//receiveMsg 是有data数据的，需要再次读取data数据
			receiveMsg := msgHead.(*znet.Message)
			receiveMsg.Data = make([]byte, receiveMsg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, receiveMsg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Receive Msg: ID=", receiveMsg.Id, ", len=", receiveMsg.DataLen, ", data=", string(receiveMsg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
