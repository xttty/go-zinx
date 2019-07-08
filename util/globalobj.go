package util

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

// GlobalObj 全局配置对象
type GlobalObj struct {
	TCPServer        ziface.IServer // 当前server对象
	Host             string         // 当前服务器ip
	TCPPort          int            // 当前服务器端口号
	Name             string         // 服务名
	Version          string         // 服务版本号
	MaxPackageSize   uint32         // 数据包最大值
	MaxConn          int            // 当前服务器最大连接数
	WorkerPoolSize   uint32         // 业务工作worker数量
	MaxWorkerTaskLen uint32         // 业务工作worker对应的任务队列最大任务数
	MaxMsgChanLen    uint32         // 读写分离后写入缓冲区大小
}

// GlobalObject 服务全局对象
var GlobalObject *GlobalObj

// Reload 重载配置
func (g *GlobalObj) Reload() {
	config, err := ioutil.ReadFile("../conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(config, GlobalObject)
	if err != nil {
		panic(err)
	}
}

// Init 服务全局对象初始化
func Init() {
	GlobalObject = &GlobalObj{
		TCPServer:        nil,
		Host:             "0.0.0.0",
		TCPPort:          7777,
		Name:             "ZinxServerApp",
		Version:          "0.0.1",
		MaxPackageSize:   12000,
		MaxConn:          4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}
	GlobalObject.Reload()
}
