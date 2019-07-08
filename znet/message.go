package znet

import (
	"zinx/ziface"
)

// Message 实现ziface.IMessage接口
type Message struct {
	// 消息ID
	ID uint32
	// 消息数据长度
	DataLen uint32
	// 消息内容
	Data []byte
}

// NewMsgPackage 生成一个消息包
func NewMsgPackage(id uint32, data []byte) ziface.IMessage {
	msg := &Message{
		ID:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
	return msg
}

// GetDataLen 获取消息数据长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// GetMsgID 获取消息ID
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

// GetData 获取消息数据内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgID 设置消息ID
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

// SetData 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// SetDataLen 设置消息长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
