package ziface

// IRequest 请求抽象接口
type IRequest interface {
	// 获得连接方法
	GetConnection() IConnection
	// 获得请求数据方法
	GetData() []byte
	// 获取消息ID
	GetMsgID() uint32
}
