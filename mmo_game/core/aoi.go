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
				minX+j*aoi.gridWidth(),
				minX+(j+1)*aoi.gridWidth(),
				minY+i*aoi.gridHeight(),
				minY+(i+1)*aoi.gridHeight())
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
