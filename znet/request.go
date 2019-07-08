package znet

import (
	"zinx/ziface"
)

// Request 请求实现结构
type Request struct {
	// 请求连接
	conn ziface.IConnection
	// 请求数据
	// data []byte
	msg ziface.IMessage
}

// GetConnection 获取连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
