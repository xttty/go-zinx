package ziface

// IMessage 消息封装接口
type IMessage interface {
	// 获取消息数据长度
	GetDataLen() uint32
	// 获取消息ID
	GetMsgID() uint32
	// 获取消息内容
	GetData() []byte
	// 设置消息ID
	SetMsgID(uint32)
	// 写入消息内容
	SetData([]byte)
	// 设置消息数据段长度
	SetDataLen(uint32)
}
