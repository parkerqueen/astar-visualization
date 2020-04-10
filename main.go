package main

import (
	"github.com/parkerqueen/a-star-go/astar"
	"github.com/parkerqueen/a-star-go/visualization"
)

func main() {
	walls := map[astar.Node]bool{
		{R: 5, C: 5}: true}

	grid := astar.Grid{Rows: 50, Cols: 50, Walls: walls}
	visualization.AStarVisualization(&grid, astar.Node{R: 0, C: 0}, astar.Node{R: 5, C: 9})
}
