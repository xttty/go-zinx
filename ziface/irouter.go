package ziface

// IRouter 路由接口
type IRouter interface {
	// 在处理conn业务之前的钩子方法
	PreHandle(r IRequest)
	// 处理conn业务方法
	Handle(r IRequest)
	// 处理conn业务之后的钩子方法
	AfterHandle(r IRequest)
}