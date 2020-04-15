package astar

import (
	"github.com/emirpasic/gods/trees/binaryheap"
)

const simpleCost = 10
const diagonalCost = 14

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
	Rows   uint
	Cols   uint
	Walls  map[Node]bool
	Artist Painter
}

//---------------------------------------------
//Exported Interfaces
//---------------------------------------------

//Painter defines some GUI library that wishes to get updates
//on the changing nodes' status throughout the algorithm
type Painter interface {
	Paint(node Node, opened bool, closed bool)
}

//---------------------------------------------
//Local Structs
//---------------------------------------------

type fNode struct {
	node   Node
	gCost  uint
	hCost  uint
	fCost  uint
	opened bool
	closed bool
	parent *fNode
}

//---------------------------------------------
//Grid Methods
//---------------------------------------------

//AStarSearch performs the pathfinding from source to dest
func (grid *Grid) AStarSearch(source, dest Node) []Node {
	fnodeMap := map[Node]*fNode{}
	openList := binaryheap.NewWith(comparator)

	fnodeMap[dest] = &fNode{node: dest}
	fnodeMap[source] = &fNode{node: source, opened: true}

	pathlen := uint(0)
	openList.Push(fnodeMap[source])

	var current *fNode
	for {
		popped, ok := openList.Pop()
		if !ok {
			return nil
		}
		current = popped.(*fNode)

		current.closed = true
		if grid.Artist != nil {
			grid.Artist.Paint(current.node, current.opened, current.closed)
		}

		if current.node == dest {
			return makepath(current, pathlen)
		}

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

			wasOpened := fnode.opened
			fnode.adjust(current, fnodeMap[dest])

			if !wasOpened {
				openList.Push(fnode)
				if grid.Artist != nil {
					grid.Artist.Paint(fnode.node, true, false)
				}
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

//---------------------------------------------
//fNode Methods
//---------------------------------------------

func (fnode *fNode) adjust(current, dest *fNode) {
	gCost := uint(0)
	hCost := uint(0)
	fCost := uint(0)

	if absDiff(fnode.node.R, current.node.R) == 1 &&
		absDiff(fnode.node.C, current.node.C) == 1 {
		gCost = current.gCost + diagonalCost
	} else {
		gCost = current.gCost + simpleCost
	}
	hCost = max(absDiff(fnode.node.R, dest.node.R),
		absDiff(fnode.node.C, dest.node.C)) * 10
	fCost = gCost + hCost

	if !fnode.opened || gCost < fnode.gCost {
		fnode.gCost = gCost
		fnode.hCost = hCost
		fnode.fCost = fCost
		fnode.parent = current
	}
	fnode.opened = true
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
	for ; fnode != nil; fnode = fnode.parent {
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
