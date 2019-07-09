package main

import (
	"fmt"
	"zinx/mmo_game/core"
)

func main() {
	aoi := core.NewAOIManager(0, 100, 0, 100, 5, 5)
	grids := aoi.GetSurroundGridByGID(6, 2)
	for _, grid := range grids {
		fmt.Println(grid.String())
	}
}
