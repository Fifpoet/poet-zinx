package znet

import (
	"fmt"
	"testing"
)

func TestDataPack(t *testing.T) {
	pack := DataPack{}
	unPack, err := pack.Unpack([]byte{11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(unPack)
}
