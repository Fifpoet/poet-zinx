package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/src/zinx/ziface"
)

type ServerConfig struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称

	/*
		Zinx
	*/
	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量

	/*
		config file path
	*/
	ConfFilePath string
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
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          7777,
		Name:             "FIF",
		Version:          "v0.7",
		MaxPacketSize:    4096,
		MaxConn:          500,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "conf/zinx.json",
	}
	// TODO 读取config
	//从配置文件中加载一些用户配置的参数
	//GlobalConfig.ReloadConfig()
}
