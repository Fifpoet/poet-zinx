package main

import "zinx/src/zinx/znet"

func main() {
	s := znet.NewServer("FIFPOET")
	s.Serve()

}
