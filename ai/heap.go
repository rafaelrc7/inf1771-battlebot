package ai

import (
	"container/heap"
)

type nodeHeap []*node

func nodeHeapNew() *nodeHeap {
	h := &nodeHeap{}
	return h
}

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].priority < h[j].priority }

func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *nodeHeap) Push(x any) {
	n := len(*h)
	no := x.(*node)
	no.index = n
	*h = append(*h, no)
}

func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	no := old[n-1]
	no.index = -1
	old[n-1] = nil
	*h = old[0 : n-1]
	return no
}

func (h *nodeHeap) TryUpdate(no *node, p float64) bool {
	if no.index == -1 {
		return false
	}

	if p >= no.priority {
		return false
	}

	no.priority = p
	heap.Fix(h, no.index)

	return false
}
