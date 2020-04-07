package astar

const simpleCost = 10
const diagonalCost = 14

type Node struct {
	R uint
	C uint
}

type Grid struct {
	Rows  uint
	Cols  uint
	Walls map[Node]bool
}

type edge struct {
	to   Node
	cost uint
}

type fNode struct {
	Node
	gCost  uint
	hCost  uint
	fCost  uint
	closed bool
}

func (grid *Grid) AStarSearch(start, end Node) []Node {
	return nil
}

func (grid *Grid) nodeEdges(node Node) []edge {
	nodeEdges := make([]edge, 0, 8)
	for r := safeDec(node.R); r < node.R+2; r++ {
		for c := safeDec(node.C); c < node.C+2; c++ {
			if (r == node.R && c == node.C) ||
				r >= grid.Rows || c >= grid.Cols {
				continue
			}

			nodeEdge := edge{to: Node{R: r, C: c}}
			if grid.Walls[nodeEdge.to] {
				continue
			}

			nodeEdge.cost = edgeCost(node, nodeEdge.to)
			nodeEdges = append(nodeEdges, nodeEdge)
		}
	}

	return nodeEdges
}

func edgeCost(start, end Node) uint {
	if absDiff(start.R, end.R) == 1 &&
		absDiff(start.C, end.C) == 1 {
		return diagonalCost
	}
	return simpleCost
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
