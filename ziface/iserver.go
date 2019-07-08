package ziface

// IServer 定义服务器接口
type IServer interface {
	Start()                                 // Start 启动服务器
	Stop()                                  // Stop 停止服务器
	Server()                                // Server 开启业务服务
	AddRouter(msgID uint32, router IRouter) // 添加路由功能，给当前服务注册一个路由业务
	GetConnMgr() IConnManager               // 获取连接管理实例
	SetOnConnStart(func(IConnection))       // 设置连接开始hook
	SetOnConnStop(func(IConnection))        // 设置连接断开hook
	CallOnConnStart(conn IConnection)       // 调用连接开始hook
	CallOnConnStop(conn IConnection)        // 调用连接断开hook
}
