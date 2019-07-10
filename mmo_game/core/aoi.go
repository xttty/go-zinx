package core

import (
	"errors"
	"fmt"
)

// AOIManager AOI管理实例
type AOIManager struct {
	MinX  int           // 区域左边界坐标
	MaxX  int           // 区域右边界坐标
	MinY  int           // 区域上边界坐标
	MaxY  int           // 区域下边界坐标
	CntX  int           // x方向格子数
	CntY  int           // y方向格子数
	grids map[int]*Grid // 当前区域中有哪些格子，key是格子ID，value是格子实例指针
}

// NewAOIManager 新建AOI管理实例
func NewAOIManager(minX int, maxX int, minY int, maxY int, cntX int, cntY int) *AOIManager {
	aoi := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntX:  cntX,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}
	for i := 0; i < cntX; i++ {
		for j := 0; j < cntY; j++ {
			gid := j*cntX + i
			aoi.grids[gid] = NewGrid(gid,
				minX+i*aoi.gridWidth(),
				minX+(i+1)*aoi.gridWidth(),
				minY+j*aoi.gridHeight(),
				minY+(j+1)*aoi.gridHeight())
		}
	}
	return aoi
}

// gridWidth 计算格子宽度
func (aoi *AOIManager) gridWidth() int {
	return (aoi.MaxX - aoi.MinX) / aoi.CntX
}

// gridHeight 计算格子高度
func (aoi *AOIManager) gridHeight() int {
	return (aoi.MaxY - aoi.MinY) / aoi.CntY
}

func (aoi *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\nminX: %d, maxX: %d, minY: %d, maxY: %d, cntX: %d, cntY: %d\nGrids in AOI Manager:\n",
		aoi.MinX, aoi.MaxX, aoi.MinY, aoi.MaxY, aoi.CntX, aoi.CntY)
	for _, grid := range aoi.grids {
		s += grid.String()
	}
	return s
}

// GetSurroundGridByGID 获取周围格子
func (aoi *AOIManager) GetSurroundGridByGID(gID int, level int) []*Grid {
	var surroundGrids []*Grid
	grid, exist := aoi.grids[gID]
	if !exist {
		return surroundGrids
	}
	// 最左边界
	minX := grid.MinX - level*aoi.gridWidth()
	if minX < aoi.MinX {
		minX = aoi.MinX
	}
	// 最右边界
	maxX := grid.MaxX + level*aoi.gridWidth()
	if maxX > aoi.MaxX {
		maxX = aoi.MaxX
	}
	// 最上边界
	minY := grid.MinY - level*aoi.gridHeight()
	if minY < aoi.MinY {
		minY = aoi.MinY
	}
	// 最下边界
	maxY := grid.MaxY + level*aoi.gridHeight()
	if maxY > aoi.MaxY {
		maxY = aoi.MaxY
	}
	for x := minX; x < maxX; x += aoi.gridWidth() {
		for y := minY; y < maxY; y += aoi.gridHeight() {
			surroundGID := ((y-aoi.MinY)/aoi.gridHeight())*aoi.CntX + (x-aoi.MinX)/aoi.gridWidth()
			surroundGrid, exist := aoi.grids[surroundGID]
			if surroundGID != gID && exist {
				surroundGrids = append(surroundGrids, surroundGrid)
			}
		}
	}
	return surroundGrids
}

// GetGIDByPos 根据坐标得到格子ID
func (aoi *AOIManager) GetGIDByPos(x, y float32) (int, error) {
	if x < float32(aoi.MinX) || x > float32(aoi.MaxX) {
		return 0, errors.New("out of range")
	}
	if y < float32(aoi.MinY) || y > float32(aoi.MaxY) {
		return 0, errors.New("out of range")
	}
	gx := (int(x) - aoi.MinX) / aoi.gridWidth()
	gy := (int(y) - aoi.MinY) / aoi.gridHeight()
	gID := gy*aoi.CntX + gx
	return gID, nil
}

// GetSurroundPlayerIDs 根据格子ID获得周围所有player
func (aoi *AOIManager) GetSurroundPlayerIDs(gID, level int) ([]int, error) {
	baseGrid, exist := aoi.grids[gID]
	if !exist {
		return nil, errors.New("grid is not exist")
	}
	var playerIDs []int
	surroundGrids := aoi.GetSurroundGridByGID(gID, level)
	playerIDs = append(playerIDs, baseGrid.GetPlayerIDs()...)
	for _, grid := range surroundGrids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}
	return playerIDs, nil
}

// GetSurroundPlayerIDsByPos 根据坐标获得周围所有player的ID
func (aoi *AOIManager) GetSurroundPlayerIDsByPos(x, y float32, level int) ([]int, error) {
	baseGID, err := aoi.GetGIDByPos(x, y)
	if err != nil {
		return nil, err
	}
	playerIDs, err := aoi.GetSurroundPlayerIDs(baseGID, level)
	return playerIDs, err
}

// GetPlayerIDsByGID 获取特定格子中的所有playerID
func (aoi *AOIManager) GetPlayerIDsByGID(gID int) []int {
	grid, exist := aoi.grids[gID]
	if !exist {
		return nil
	}
	return grid.GetPlayerIDs()
}

// RemovePlayerFromGrid 将player从一个格子中移除
func (aoi *AOIManager) RemovePlayerFromGrid(pID, gID int) {
	grid, exist := aoi.grids[gID]
	if !exist {
		return
	}
	grid.Remove(pID)
}

// AddPlayerToGrid 将player添加到格子里
func (aoi *AOIManager) AddPlayerToGrid(pID, gID int) {
	grid, exist := aoi.grids[gID]
	if !exist {
		return
	}
	grid.Add(pID)
}

// AddPlayerToGridByPos 将player通过坐标添加到相应的格子中
func (aoi *AOIManager) AddPlayerToGridByPos(pID int, x, y float32) {
	gID, err := aoi.GetGIDByPos(x, y)
	if err != nil {
		return
	}
	aoi.AddPlayerToGrid(pID, gID)
}

// RemovePlayerFromGridByPos 将player通过坐标从相应格子中移除
func (aoi *AOIManager) RemovePlayerFromGridByPos(pID int, x, y float32) {
	gID, err := aoi.GetGIDByPos(x, y)
	if err != nil {
		return
	}
	aoi.RemovePlayerFromGrid(pID, gID)
}
