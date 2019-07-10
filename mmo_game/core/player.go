package core

import (
	"fmt"
	"math/rand"
	"sync"
	"zinx/mmo_game/pb"
	"zinx/ziface"

	"github.com/golang/protobuf/proto"
)

// Player 玩家实例
type Player struct {
	Pid  int32              // 玩家ID
	Conn ziface.IConnection // 玩家连接
	X    float32            // 平面X轴坐标
	Y    float32            // 高度
	Z    float32            // 平面y轴坐标
	V    float32            // 旋转视角 0~360度
}

// PidGen 用来生成玩家ID的计数器
var PidGen int32 = 1

// IDLock 计数保护锁
var IDLock sync.Mutex

// NewPlayer 生成玩家
func NewPlayer(conn ziface.IConnection) *Player {
	IDLock.Lock()
	playerID := PidGen
	PidGen++
	IDLock.Unlock()
	player := &Player{
		Pid:  playerID,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(134 + rand.Intn(17)),
	}
	return player
}

// SendMsg 发送消息给客户端
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal message error", err)
		return
	}
	if p.Conn == nil {
		fmt.Println("player", p.Pid, "have no connection")
		return
	}
	err = p.Conn.SendBuffMsg(msgID, msg)
	if err != nil {
		fmt.Println("send msg error", err)
		return
	}
}

// SyncPid 同步给客户端玩家ID
func (p *Player) SyncPid() {
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(uint32(MsgSyncPid), data)
}

// BroadCastStartPosition 广播初始位置
func (p *Player) BroadCastStartPosition() {
	position := &pb.Position{
		X: p.X,
		Y: p.Y,
		Z: p.Z,
		V: p.V,
	}
	data := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  BroadCastPosType,
		Data: &pb.BroadCast_P{
			P: position,
		},
	}
	p.SendMsg(uint32(MsgBroadCast), data)
}
