package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/util"
	"zinx/ziface"
)

// Connection IConnection接口的实现
type Connection struct {
	TCPServer      ziface.IServer         // 该连接属于的tcp服务
	Conn           *net.TCPConn           // 当前连接原始的socket套接字
	ConnID         uint32                 // 当前连接ID
	isClosed       bool                   // 当前连接是否关闭
	ExitBuffChan   chan bool              // 告知该连接已经退出/停止的channel
	MsgHandler     ziface.IMsgHandler     // 消息管理
	MsgChan        chan []byte            // 读写分离后消息传递管道
	MsgBuffChan    chan []byte            //读写分离后带缓冲的消息管道
	ExitReaderChan chan bool              // 读取协程关闭通知
	ExitWriterChan chan bool              // 写协程关闭通知
	property       map[string]interface{} //连接属性
	propertyLock   sync.RWMutex           // 连接属性修改锁
}

// NewConnection 新建一个连接
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer:      server,
		Conn:           conn,
		ConnID:         connID,
		isClosed:       false,
		ExitBuffChan:   make(chan bool, 3),
		MsgHandler:     handler,
		MsgChan:        make(chan []byte),
		MsgBuffChan:    make(chan []byte, util.GlobalObject.MaxMsgChanLen),
		ExitReaderChan: make(chan bool, 1),
		ExitWriterChan: make(chan bool, 1),
		property:       make(map[string]interface{}),
	}
	// 连接添加进入连接池
	server.GetConnMgr().Add(c)
	return c
}

// startReader 处理conn读数据的goroutine
func (c *Connection) startReader() {
	fmt.Println("Reader Goroutine is Runing")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")

	for {
		select {
		case <-c.ExitReaderChan:
			return
		default:
			// 解包实例
			dp := NewDataPack()
			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
				fmt.Println("read msg head error ", err)
				c.ExitBuffChan <- true
				return
			}
			// 解包获取头信息放入msg中
			msg, err := dp.Unpack(headData)
			if err != nil {
				c.ExitBuffChan <- true
				return
			}
			// 如果有数据则存放到Data
			if msg.GetDataLen() > 0 {
				data := make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
					fmt.Println("read msg data error ", err)
					c.ExitBuffChan <- true
					return
				}
				msg.SetData(data)
			}
			r := &Request{
				conn: c,
				msg:  msg,
			}
			if util.GlobalObject.WorkerPoolSize > 0 {
				// 发送给消息队列
				c.MsgHandler.SendMsgToTaskQueue(r)
			} else {
				// 开启协程来处理请求
				go c.MsgHandler.DoMsgHandler(r)
			}
		}
	}
}

// startWriter 将数据写到客户端
func (c *Connection) startWriter() {
	defer fmt.Println(c.RemoteAddr().String(), " conn writer exit!")
	for {
		select {
		case data, ok := <-c.MsgChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("write data error ", err)
					c.ExitBuffChan <- true
					return
				}
			} else {
				fmt.Println("msg chan is close!")
			}
		case data, ok := <-c.MsgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("write data error", err)
					c.ExitBuffChan <- true
					return
				}
			} else {
				fmt.Println("msg buff chan is close!")
			}
		case <-c.ExitWriterChan:
			return
		}
	}
}

// Start 连接启动方法
func (c *Connection) Start() {
	defer c.Stop()

	// 开启处理该连接读取到客户端数据之后的请求业务
	// 读写分离
	go c.startReader()
	go c.startWriter()

	// connection start hook
	c.TCPServer.CallOnConnStart(c)

	for {
		select {
		case <-c.ExitBuffChan:
			// 得到退出消息，不再阻塞
			return
		}
	}
}

// Stop 连接停止方法
func (c *Connection) Stop() {
	// 如果当前连接已经关闭则直接return
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 如果用户注册了该连接的关闭回调，那么在此刻应该显示调用

	// 从连接池中remove
	c.TCPServer.GetConnMgr().Remove(c)
	// conn stop callback
	c.TCPServer.CallOnConnStop(c)
	// 关闭socket连接
	c.Conn.Close()
	// 清空连接属性
	for key := range c.property {
		c.RemoveProperty(key)
	}

	// 通知缓存队列中的读数据业务，该连接已经关闭
	fmt.Println("prepare closing channel")
	c.ExitReaderChan <- true
	c.ExitWriterChan <- true

	// 关闭管道
	close(c.ExitBuffChan)
	close(c.ExitReaderChan)
	close(c.ExitWriterChan)
	close(c.MsgChan)
	close(c.MsgBuffChan)
	fmt.Println("connection close!")
}

// GetTCPConnection 获得原始的socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获得连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获得客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 发送信息
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	dp := NewDataPack()
	sendMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack msg error ", err)
		return err
	}
	// 发送给写协程
	c.MsgChan <- sendMsg
	return nil
}

// SendBuffMsg 缓冲发送
func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	dp := NewDataPack()
	sendMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack msg error", err)
		return err
	}
	c.MsgBuffChan <- sendMsg
	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if _, exist := c.property[key]; !exist {
		c.property[key] = value
	}
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	value, exist := c.property[key]
	if !exist {
		return nil, errors.New("key is not exist")
	}
	return value, nil
}

// RemoveProperty 删除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
