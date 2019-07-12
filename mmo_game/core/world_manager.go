package core

import (
	"errors"
	"sync"
)

// WorldManager 世界管理模块
type WorldManager struct {
	AoiMgr  *AOIManager       // aoi管理实例
	players map[int32]*Player // 玩家
	pLock   sync.RWMutex      // 玩家读写锁
}

// WorldMgrObj 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldManager

// Init 初始世界
func Init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(WorldMinX, WorldMaxX, WorldMinY, WorldMinY, WorldCntX, WorldCntY),
		players: make(map[int32]*Player),
	}
}

// AddPlayer 向世界中加入玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.players[player.Pid] = player
	wm.pLock.Unlock()

	wm.AoiMgr.AddPlayerToGridByPos(int(player.Pid), player.X, player.Y)
}

// RemovePlayerByPID 删除玩家
func (wm *WorldManager) RemovePlayerByPID(pID int32) {
	wm.pLock.Lock()
	player, exist := wm.players[pID]
	if exist {
		wm.AoiMgr.RemovePlayerFromGridByPos(int(pID), player.X, player.Y)
	}
	delete(wm.players, pID)
	wm.pLock.Unlock()
}

// GetPlayerByPID 根据PID获取玩家
func (wm *WorldManager) GetPlayerByPID(pID int32) (*Player, error) {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	player, exist := wm.players[pID]
	if !exist {
		return nil, errors.New("Player not found")
	}
	return player, nil
}

// GetAllPlayers 获取所有玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)
	for _, player := range wm.players {
		players = append(players, player)
	}
	return players
}
