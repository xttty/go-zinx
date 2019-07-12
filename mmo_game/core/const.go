package core

const (
	_ int32 = iota
	// MsgSyncPid 同步pid消息
	MsgSyncPid
	// MsgBroadCast 广播消息
	MsgBroadCast
)

const (
	_ int32 = iota
	// BroadCastChatType 广播聊天类型
	BroadCastChatType
	// BroadCastPosType 广播位置类型
	BroadCastPosType
	// BroadCastActionType 广播动作类型
	BroadCastActionType
	// BroadCastPosUpdateType 广播位置更新类型
	BroadCastPosUpdateType
)

const (
	// WorldMinX 世界地图X轴起点
	WorldMinX int = 0
	// WorldMaxX 世界地图X轴终点点
	WorldMaxX int = 500
	// WorldMinY 世界地图Y轴起点
	WorldMinY int = 0
	// WorldMaxY 世界地图Y轴终点
	WorldMaxY int = 1000
	// WorldCntX 世界地图X轴格子数
	WorldCntX int = 200
	// WorldCntY 世界地图Y轴格子数
	WorldCntY int = 500
)
