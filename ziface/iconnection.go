package ziface

import (
	"net"
)

// IConnection 连接接口定义
type IConnection interface {
	// 启动连接，让当前连接工作
	Start()
	// 停止连接，结束当前连接状态
	Stop()
	// 从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接ID
	GetConnID() uint32
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 发送消息
	SendMsg(msgID uint32, data []byte) error
	// 发送消息，带缓冲区
	SendBuffMsg(msgID uint32, data []byte) error
	//设置连接属性
	SetProperty(key string, value interface{})
	// 获取连接属性
	GetProperty(key string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(key string)
}
