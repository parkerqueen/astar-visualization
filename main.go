package main

import (
	"github.com/parkerqueen/a-star-go/astar"
)

func main() {
	walls := map[astar.Node]bool{
		{R: 5, C: 5}: true}

	grid := astar.Grid{Rows: 50, Cols: 50, Walls: walls}
	grid.AStarVisualization(astar.Node{R: 5, C: 6}, astar.Node{R: 8, C: 25})
}
