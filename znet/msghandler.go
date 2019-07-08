package znet

import (
	"fmt"
	"zinx/util"
	"zinx/ziface"
)

// MsgHandler 消息管理实例
type MsgHandler struct {
	// 路由
	Apis map[uint32]ziface.IRouter
	// 业务工作worker数量
	WorkerPoolSize uint32
	// 消息队列，worker从此队列中获得消息
	TaskQueue []chan ziface.IRequest
}

// NewMsgHandler 新建消息管理实例
func NewMsgHandler() *MsgHandler {
	mh := &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: util.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, util.GlobalObject.WorkerPoolSize),
	}
	return mh
}

// AddRouter 添加路由
func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("repeated api, msgID: ", msgID)
		return
	}
	mh.Apis[msgID] = router
}

// DoMsgHandler 执行路由
func (mh *MsgHandler) DoMsgHandler(rq ziface.IRequest) {
	router, ok := mh.Apis[rq.GetMsgID()]
	if !ok {
		fmt.Println("There is not router bind on this message")
		return
	}
	router.PreHandle(rq)
	router.Handle(rq)
	router.AfterHandle(rq)
}

// StartOneWorker 启动一个worker
func (mh *MsgHandler) StartOneWorker(workID uint32, taskQueue chan ziface.IRequest) {
	fmt.Printf("worker[%d] is started\n", workID)
	for {
		select {
		case rq := <-taskQueue:
			mh.DoMsgHandler(rq)
		}
	}
}

// StartWorkerPool 启动工作池
func (mh *MsgHandler) StartWorkerPool() {
	var i uint32
	for i < mh.WorkerPoolSize {
		mh.TaskQueue[i] = make(chan ziface.IRequest, util.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
		i++
	}
}

// SendMsgToTaskQueue 将消息传递给消息队列
func (mh *MsgHandler) SendMsgToTaskQueue(rq ziface.IRequest) {
	// 动态分配至其中一个任务队列
	workID := rq.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Printf("request[%d] assign to worker[%d]\n", rq.GetMsgID(), workID)
	mh.TaskQueue[workID] <- rq
}
