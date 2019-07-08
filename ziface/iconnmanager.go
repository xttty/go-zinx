package ziface

// IConnManager 连接管理接口
type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 删除连接
	Get(connID uint32) (IConnection, error) // 获取连接
	Len() int                               // 连接数量
	ClearAll()                              // 清空并停止所有连接
}
