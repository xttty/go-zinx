package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/util"
	"zinx/ziface"
)

// DataPack 封包拆包实例
type DataPack struct{}

// NewDataPack 构造拆封包实例
func NewDataPack() ziface.IDataPack {
	dp := &DataPack{}
	return dp
}

// GetHeadLen 返回包头部长度
func (dp *DataPack) GetHeadLen() uint32 {
	// id uint32(4 bytes)  +  DataLen uint32(4 bytes)
	return 8
}

// Pack 数据封包
func (dp *DataPack) Pack(m ziface.IMessage) ([]byte, error) {
	// 创建一个bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})
	// 写DataLen
	if err := binary.Write(dataBuf, binary.LittleEndian, m.GetDataLen()); err != nil {
		return nil, err
	}
	// 写msgID
	if err := binary.Write(dataBuf, binary.LittleEndian, m.GetMsgID()); err != nil {
		return nil, err
	}
	// 写data数据
	if err := binary.Write(dataBuf, binary.LittleEndian, m.GetData()); err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

// Unpack 数据解包，只解出head头信息
func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewBuffer(data)
	msg := &Message{}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	if util.GlobalObject.MaxPackageSize > 0 && msg.GetDataLen() > util.GlobalObject.MaxPackageSize {
		return nil, errors.New("Too large msg data received")
	}
	return msg, nil
}
