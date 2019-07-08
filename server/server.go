package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// NumberRouter 返回数字路由
type NumberRouter struct {
	znet.BaseRouter
}

// PreHandle PingRouter自己的实现
func (pr *PingRouter) PreHandle(rq ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
}

// Handle PingRouter自己的实现
func (pr *PingRouter) Handle(rq ziface.IRequest) {
	fmt.Println("Call Router Handle")
	data := rq.GetData()
	msgID := rq.GetMsgID()
	err := rq.GetConnection().SendMsg(msgID, data)
	if err != nil {
		fmt.Println("call back ping error:", err)
	}
}

// AfterHandle PingRouter自己的实现
func (pr *PingRouter) AfterHandle(rq ziface.IRequest) {
	fmt.Println("Call Router After Router")
}

// Handle NumberRouter Handle方法重写
func (nr *NumberRouter) Handle(rq ziface.IRequest) {
	data := []byte("1 2 3 4 5")
	msgID := rq.GetMsgID()
	err := rq.GetConnection().SendBuffMsg(msgID, data)
	if err != nil {
		fmt.Println("call back number router error: ", err)
	}
}

func connStartCallback(conn ziface.IConnection) {
	fmt.Println(conn.RemoteAddr(), "connected!")
}

func connStopCallback(conn ziface.IConnection) {
	fmt.Println(conn.RemoteAddr(), "has closed!")
}

func main() {
	s := znet.NewServer()
	pr := &PingRouter{}
	nr := &NumberRouter{}
	s.AddRouter(0, pr)
	s.AddRouter(1, nr)
	s.SetOnConnStart(connStartCallback)
	s.SetOnConnStop(connStopCallback)
	s.Server()
}
