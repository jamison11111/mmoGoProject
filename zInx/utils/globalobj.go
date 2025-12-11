package utils

import (
	"encoding/json"
	"lwc/zInx/ziface"
	"os"
)

/*
存储一切与zInx框架有关的全局参数,供其他模块使用，一些参数也可通过用户根据zinx.json文件来配置
*/
type GlobalObj struct {
	TcpServer        ziface.IServer
	Host             string
	TcpPort          int
	Name             string
	Version          string //zInx版本号
	MaxPacketSize    uint32 //单次读取的数据包的最大值
	MaxConn          int
	WorkerPoolSize   uint32 //业务工作池的worker数量
	MaxWorkerTaskLen uint32 //每个队列的最大容量
	ConfFilePath     string //配置文件的路径
	MaxMsgChanLen    int    //一条连接中,有缓冲的通信通道的最大容量,也即能允许的写协程来不及写出去但业务协程已经处理且打包好的消息的最大数量
}

// 定义一个全局对象
var GlobalObject *GlobalObj

func (G *GlobalObj) Reload() {
	data, err := os.ReadFile(G.ConfFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	//这些都是默认值,会被配置文件的值给覆盖
	GlobalObject = &GlobalObj{
		Name:             "default server name",
		Version:          "V0.4",
		TcpPort:          7776,
		Host:             "0.0.0.0",
		MaxConn:          10000,
		MaxPacketSize:    4096,
		ConfFilePath:     "D:/goProject/mmo_game/conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    55,
	}
	//配置文件的参数会覆盖上面的一些默认配置
	GlobalObject.Reload()
}
