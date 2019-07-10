package ziface

import (
	"net"
)

// IConnection 连接接口定义
type IConnection interface {
	Start()                                      // 启动连接，让当前连接工作
	Stop()                                       // 停止连接，结束当前连接状态
	GetTCPConnection() *net.TCPConn              // 从当前连接获取原始的socket TCPConn
	GetConnID() uint32                           // 获取当前连接ID
	RemoteAddr() net.Addr                        // 获取远程客户端地址信息
	SendMsg(msgID uint32, data []byte) error     // 发送消息
	SendBuffMsg(msgID uint32, data []byte) error // 发送消息，带缓冲区
	SetProperty(key string, value interface{})   // 设置连接属性
	GetProperty(key string) (interface{}, error) // 获取连接属性
	RemoveProperty(key string)                   // 移除连接属性
}
