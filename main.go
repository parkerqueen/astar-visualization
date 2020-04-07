package main

import (
	"fmt"

	"github.com/parkerqueen/a-star-go/astar"
)

func main() {
	walls := map[astar.Node]bool{
		{R: 5, C: 5}: true}

	grid := astar.Grid{Rows: 10, Cols: 10, Walls: walls}
	path := grid.AStarSearch(astar.Node{R: 0, C: 0}, astar.Node{R: 9, C: 9})
	fmt.Println(path)
}
