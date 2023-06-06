package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/src/zinx/ziface"
)

type ServerConfig struct {
	TcpServer     ziface.IServer `json:"tcp_server,omitempty"`      //当前Zinx的全局Server对象
	Host          string         `json:"host,omitempty"`            //当前服务器主机IP
	TcpPort       int            `json:"tcp_port,omitempty"`        //当前服务器主机监听端口号
	Name          string         `json:"name,omitempty"`            //当前服务器名称
	Version       string         `json:"version,omitempty"`         //当前Zinx版本号
	MaxPacketSize uint32         `json:"max_packet_size,omitempty"` //都需数据包的最大值
	MaxConn       int            `json:"max_conn,omitempty"`        //当前服务器主机允许的最大链接个数
	Test          json.RawMessage
}

var GlobalConfig *ServerConfig

// ReloadConfig 读取用户配置文件
func (s *ServerConfig) ReloadConfig() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(GlobalConfig)
}

func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalConfig = &ServerConfig{
		Name:          "ZinxApp",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	// TODO 读取config
	//从配置文件中加载一些用户配置的参数
	//GlobalConfig.ReloadConfig()
}
