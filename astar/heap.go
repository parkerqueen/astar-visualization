package astar

import "github.com/emirpasic/gods/trees/binaryheap"

type heap struct {
	bheap *binaryheap.Heap
}

func newHeap() *heap {
	comparator := func(a, b interface{}) int {
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

	return &heap{bheap: binaryheap.NewWith(comparator)}
}

func (bheap *heap) push(fnodes ...*fNode) {
	bheap.bheap.Push(fnodes)
}

func (bheap *heap) pop() (*fNode, bool) {
	fnode, ok := bheap.bheap.Pop()
	return fnode.(*fNode), ok
}
