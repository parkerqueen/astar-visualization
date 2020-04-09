//TODOS
//Add comments
//Think of a way to carve out the GUI code in a separate file

package astar

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"github.com/emirpasic/gods/trees/binaryheap"
)

const simpleCost = 10
const diagonalCost = 14

var simpleCellCol = color.RGBA{R: 255, G: 255, B: 255, A: 1}
var sourceCellCol = color.RGBA{R: 135, G: 101, B: 53, A: 1}
var destCellCol = color.RGBA{R: 53, G: 135, B: 105, A: 1}
var openCellCol = color.RGBA{R: 20, G: 30, B: 100, A: 1}
var closedCellCol = color.RGBA{R: 251, G: 25, B: 10, A: 1}
var blockedCellCol = color.RGBA{R: 74, G: 71, B: 71, A: 1}

var guiApp fyne.App
var guiWin fyne.Window

//---------------------------------------------
//Exported Structs
//---------------------------------------------

//A Node struct denoting a location on the grid
type Node struct {
	R uint
	C uint
}

//A Grid with a specific number of rows & columns
type Grid struct {
	Rows  uint
	Cols  uint
	Walls map[Node]bool
}

//---------------------------------------------
//Local Structs
//---------------------------------------------

type fNode struct {
	node    Node
	gCost   uint
	hCost   uint
	fCost   uint
	closed  bool
	visited bool
	parent  *fNode
}

//---------------------------------------------
//Grid Methods
//---------------------------------------------

//AStarVisualization performs the A Star path finding animation
func (grid *Grid) AStarVisualization(start, end Node) {
	guiApp = app.New()

	guiWin = guiApp.NewWindow("A* Path Visualization")
	guiWin.Resize(fyne.Size{Width: 1200, Height: 800})
	guiWin.SetFixedSize(true)
	guiWin.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeySpace {
			grid.aStarSearch(start, end)
		}
	})

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(int(grid.Cols)))
	for r := uint(0); r < grid.Rows; r++ {
		for c := uint(0); c < grid.Rows; c++ {
			cellCol := simpleCellCol
			node := Node{R: r, C: c}
			if start == node {
				cellCol = sourceCellCol
			} else if end == node {
				cellCol = destCellCol
			} else if grid.Walls[node] {
				cellCol = blockedCellCol
			}
			container.AddObject(canvas.NewRectangle(cellCol))
		}
	}

	guiWin.SetContent(container)
	guiWin.Show()
	guiApp.Run()
}

func (grid *Grid) aStarSearch(start, end Node) []Node {
	fnodeMap := map[Node]*fNode{}
	openList := binaryheap.NewWith(comparator)

	fnodeMap[end] = &fNode{node: end}
	fnodeMap[start] = &fNode{node: start, visited: true}

	pathlen := uint(0)
	openList.Push(fnodeMap[start])

	var current *fNode
	for {
		popped, ok := openList.Pop()
		if !ok {
			return nil
		}
		current = popped.(*fNode)

		current.closed = true
		if current.node == end {
			return makepath(current, pathlen)
		}
		grid.updateCellCol(current)

		pathlen = pathlen + 1

		neighbours := grid.nodeNeighbours(current.node)
		for _, neighbour := range neighbours {
			fnode := fnodeMap[neighbour]

			if fnode == nil {
				fnode = &fNode{node: neighbour}
				fnodeMap[neighbour] = fnode
			}

			if fnode.closed {
				continue
			}

			fnode.adjustFCost(current, fnodeMap[end], fnode.visited)
			if !fnode.visited {
				fnode.visited = true
				openList.Push(fnode)
				grid.updateCellCol(fnode)
			}
		}
	}
}

func (grid *Grid) nodeNeighbours(node Node) []Node {
	neighbours := make([]Node, 0, 8)
	for r := safeDec(node.R); r < node.R+2; r++ {
		for c := safeDec(node.C); c < node.C+2; c++ {
			if (r == node.R && c == node.C) ||
				r >= grid.Rows || c >= grid.Cols {
				continue
			}

			neighbour := Node{R: r, C: c}
			if grid.Walls[neighbour] {
				continue
			}
			neighbours = append(neighbours, neighbour)
		}
	}

	return neighbours
}

func (grid *Grid) updateCellCol(fnode *fNode) {
	guiContainer := guiWin.Content().(*fyne.Container)
	cellIndex := int(fnode.node.R*grid.Cols + fnode.node.C)

	gridCell := guiContainer.Objects[cellIndex].(*canvas.Rectangle)

	color := openCellCol
	if fnode.closed {
		color = closedCellCol
	}

	gridCell.FillColor = color
	gridCell.Refresh()
}

//---------------------------------------------
//fNode Methods
//---------------------------------------------

func (fnode *fNode) adjustFCost(start, end *fNode, visited bool) {
	gCost := uint(0)
	hCost := uint(0)
	fCost := uint(0)
	current := fnode.node

	if absDiff(current.R, start.node.R) == 1 &&
		absDiff(current.C, start.node.C) == 1 {
		gCost = start.gCost + diagonalCost
	} else {
		gCost = start.gCost + simpleCost
	}
	hCost = max(absDiff(current.R, end.node.R), absDiff(current.C, end.node.C)) * 10
	fCost = gCost + hCost

	if !visited || gCost < fnode.gCost {
		fnode.parent = start
		fnode.gCost = gCost
		fnode.hCost = hCost
		fnode.fCost = fCost
	}
}

//---------------------------------------------
//Helper Functions for A* Pathfinding
//---------------------------------------------

func comparator(a, b interface{}) int {
	x := a.(*fNode)
	y := b.(*fNode)

	if x.fCost < y.fCost {
		return -1
	} else if x.fCost == y.fCost {
		return 0
	} else {
		return 1
	}
}

func makepath(fnode *fNode, pathlen uint) []Node {
	path := make([]Node, 0, pathlen)
	for ; fnode.parent != nil; fnode = fnode.parent {
		path = append(path, fnode.node)
	}
	return path
}

//---------------------------------------------
//Utility Functions
//---------------------------------------------

//A function that finds the absolute difference between two unisgned integers.
func absDiff(x, y uint) uint {
	if x > y {
		return x - y
	}
	return y - x
}

//A function to decrement an unsigned integer but not below 0.
func safeDec(x uint) uint {
	if x == 0 {
		return 0
	}
	return x - 1
}

//A simple function to return the maximum of two integers (unsigned)
func max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}
