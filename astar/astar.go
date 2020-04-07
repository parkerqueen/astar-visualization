package astar

const simpleCost = 10
const diagonalCost = 14

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

type fNode struct {
	node    Node
	gCost   uint
	hCost   uint
	fCost   uint
	closed  bool
	visited bool
	parent  *fNode
}

func (fnode *fNode) adjustFCost(start, end *fNode, visited bool) {
	current := fnode.node
	gCost := uint(0)
	hCost := uint(0)
	fCost := uint(0)

	if absDiff(current.R, start.node.R) == 1 &&
		absDiff(current.C, start.node.C) == 1 {
		gCost = start.gCost + diagonalCost
	} else {
		gCost = start.gCost + simpleCost
	}

	hCost = max(absDiff(current.R, end.node.R), absDiff(current.C, end.node.C))
	fCost = gCost + hCost

	if !visited || fCost < fnode.fCost {
		fnode.gCost = gCost
		fnode.hCost = hCost
		fnode.fCost = fCost
	}
}

//AStarSearch performs the A Star path finding on a grid
func (grid *Grid) AStarSearch(start, end Node) []Node {
	openList := newHeap()
	fnodeMap := map[Node]*fNode{}

	fnode := &fNode{node: end}
	fnodeMap[end] = fnode

	fnode = &fNode{node: start}
	openList.push(fnode)
	fnodeMap[start] = fnode

	pathlen := uint(0)
	for {
		current, ok := openList.pop()
		if !ok {
			break
		}
		current.closed = true

		if current.node == end {
			path := make([]Node, 0, pathlen)
			for ; current.parent != nil; current = current.parent {
				path = append(path, current.node)
			}
			return path
		}
		pathlen = pathlen + 1

		neighbours := grid.nodeNeighbours(current.node)
		for _, neighbour := range neighbours {
			fnode = fnodeMap[neighbour]

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
				openList.push(fnode)
			}
		}
	}

	return nil
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
