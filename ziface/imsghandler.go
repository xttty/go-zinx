package ziface

// IMsgHandler 消息管理接口
type IMsgHandler interface {
	// 马上以非阻塞的方式处理消息
	DoMsgHandler(r IRequest)
	// 为消息绑定路由
	AddRouter(msgID uint32, router IRouter)
	// 启动工作池
	StartWorkerPool()
	// 将消息传递给任务队列
	SendMsgToTaskQueue(rq IRequest)
}
