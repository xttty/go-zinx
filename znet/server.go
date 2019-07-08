package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/util"
	"zinx/ziface"
)

// Server IServer接口实现，定义一个服务器
type Server struct {
	Name              string                   // 服务器名字
	IPVersion         string                   // tcp4 或者其他
	IP                string                   // ip地址
	Port              int                      // 端口号
	MsgHandler        ziface.IMsgHandler       // 消息管理实例
	ConnMgr           ziface.IConnManager      // 连接管理
	ConnStartCallback func(ziface.IConnection) // 连接开始回调方法
	ConnStopCallback  func(ziface.IConnection) // 连接结束回调方法
}

// Start 实现IServer接口中的Start方法
func (s *Server) Start() {
	fmt.Printf("[START] %s listener at IP : %s, Port %d, is starting\n", s.Name, s.IP, s.Port)

	// 开启一个go routine去做服务端listener业务
	go func() {
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		// 监听成功
		fmt.Println("Start Zinx server ", s.Name, " success, now listenning...")

		//3 启动server网络连接业务
		//3.1 启动消息连接池
		s.MsgHandler.StartWorkerPool()

		// TODO 实现一个创建connID的方法
		var connID uint32
		connID = 0
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error ", err)
			}
			//3.2 设置服务器最大连接控制，如果超过最大连接，那么则关闭此新的连接
			if util.GlobalObject.MaxConn > 0 && util.GlobalObject.MaxConn <= s.ConnMgr.Len() {
				conn.Close()
				fmt.Println("count of connections excess!")
				continue
			}
			//3.3 处理该新连接请求的业务方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, connID, s.MsgHandler)
			connID++

			go dealConn.Start()
		}
	}()
}

// Stop 实现IServer接口中的Stop方法
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name ", s.Name)
	// 将其他需要清理的连接信息或者其他信息也要一并处理或者清理
	// 清空所有连接
	s.ConnMgr.ClearAll()
}

// Server 实现IServer接口中的Server方法
func (s *Server) Server() {
	defer s.Stop()

	s.Start()
	// TODO Server.Server() 服务启动时需要处理的其他事情

	for i := 0; i < 5; i++ {
		time.Sleep(10 * time.Second)
	}
}

// AddRouter 实现了IServer接口中的AddRouter方法
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

// GetConnMgr 获取连接管理实例
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置连接开始回调方法
func (s *Server) SetOnConnStart(callback func(ziface.IConnection)) {
	s.ConnStartCallback = callback
}

// SetOnConnStop 设置连接结束回调方法
func (s *Server) SetOnConnStop(callback func(ziface.IConnection)) {
	s.ConnStopCallback = callback
}

// CallOnConnStart 调用连接开始回调方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.ConnStartCallback != nil {
		fmt.Println("call connection start callback")
		s.ConnStartCallback(conn)
	}
}

// CallOnConnStop 调用连接结束回调方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.ConnStopCallback != nil {
		fmt.Println("call connection stop callback")
		s.ConnStopCallback(conn)
	}
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	// 初始化全局配置
	util.Init()

	s := &Server{
		Name:       util.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         util.GlobalObject.Host,
		Port:       util.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
