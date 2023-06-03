package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for i := 1; i < 100; i++ {
		_, err := conn.Write([]byte(fmt.Sprintf("Hello, %d", i)))
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		// 客户端接受回写的数据 并打印
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		fmt.Printf(" server call back : %s, cnt = %d\n", buf, cnt)

		time.Sleep(2 * time.Second)
	}
}
