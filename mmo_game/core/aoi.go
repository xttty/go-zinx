package core

import "fmt"

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
