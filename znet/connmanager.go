package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// ConnManager 连接管理实例
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 连接池
	connLock    sync.RWMutex                  // 连接读写锁
}

// NewConnManager 新建连接管理实例
func NewConnManager() *ConnManager {
	cm := &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
	return cm
}

// Add 实现连接管理Add方法
func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	id := conn.GetConnID()
	if _, exist := cm.connections[id]; !exist {
		cm.connections[id] = conn
	}
	fmt.Println("add connection, connID =", conn.GetConnID(), "current conn count:", cm.Len())
}

// Remove 实现连接管理Remove方法
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection remove connID =", conn.GetConnID(), "current conn count:", cm.Len())
}

// Len 连接管理中连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// Get 获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	conn, exist := cm.connections[connID]
	var err error
	if !exist {
		err = errors.New("connection is not exist")
	}
	return conn, err
}

// ClearAll 清空并停止所有连接
func (cm *ConnManager) ClearAll() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
	fmt.Println("clear all connections, current connection count:", cm.Len())
}
