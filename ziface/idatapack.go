package ziface

// IDataPack 数据打包接口
// 封包和拆包：直接面向TCP连接中的数据流，为传递数据添加头部信息，用于解决TCP粘包问题
type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
