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
}

// HandleFunc 定义一个统一处理链接业务的接口
// type HandleFunc func(*net.TCPConn, []byte, int) error