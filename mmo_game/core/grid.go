package core

import (
	"fmt"
	"sync"
)

// Grid 地图格子结构
type Grid struct {
	GID       int          // 格子ID
	MinX      int          // 格子左边界坐标
	MaxX      int          // 格子右边界坐标
	MinY      int          // 格子上边界坐标
	MaxY      int          // 格子下边界坐标
	playerIDs map[int]bool // 当前格子内玩家或者物体成员ID
	pIDLock   sync.RWMutex // playerIDs的保护map锁
}

// NewGrid 新建一个格子实例
func NewGrid(gID int, minX int, maxX int, minY int, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add 向格子中添加玩家
func (g *Grid) Add(playID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playID] = true
}

// Remove 向当前格子中移除玩家
func (g *Grid) Remove(playID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playID)
}

// GetPlayerIDs 获取当前格子中的所有玩家ID
func (g *Grid) GetPlayerIDs() []int {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	var playIDs []int
	for pID := range g.playerIDs {
		playIDs = append(playIDs, pID)
	}
	return playIDs
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
