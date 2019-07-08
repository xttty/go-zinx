package znet

import (
	"zinx/ziface"
)

// BaseRouter 基础路由
type BaseRouter struct{}

// PreHandle base empty function
func (router *BaseRouter) PreHandle(rq ziface.IRequest) {}

// Handle base empty function
func (router *BaseRouter) Handle(rq ziface.IRequest) {}

// AfterHandle base empty function
func (router *BaseRouter) AfterHandle(rq ziface.IRequest) {}
