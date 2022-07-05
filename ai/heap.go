package ai

import (
	"container/heap"
	"github.com/rafaelrc7/inf1771-battlebot/gamemap"
)

type CoordCost struct {
	Coord    gamemap.Coord
	Priority float64
}

type CoordHeap []*CoordCost

func CoordHeapNew() *CoordHeap {
	h := &CoordHeap{}
	heap.Init(h)
	return h
}

func (h CoordHeap) Len() int           { return len(h) }
func (h CoordHeap) Less(i, j int) bool { return h[i].Priority < h[j].Priority }
func (h CoordHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *CoordHeap) Push(x any) {
	*h = append(*h, x.(*CoordCost))
}
func (h *CoordHeap) Pop() any {
	old := *h
	n := len(old)
	c := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return c
}

func (h *CoordHeap) PushCoord(c gamemap.Coord, Priority float64) {
	heap.Push(h, &CoordCost{Coord: c, Priority: Priority})
}

func (h *CoordHeap) TryUpdate(c gamemap.Coord, Priority float64) bool {
	for i, v := range *h {
		if c == v.Coord {
			if Priority < v.Priority {
				v.Priority = Priority
				heap.Fix(h, i)
				return true
			} else {
				return false
			}
		}
	}

	return false
}
