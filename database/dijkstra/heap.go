//From: https://dev.to/douglasmakey/implementation-of-dijkstra-using-heap-in-go-6e3

package dijkstra

import "container/heap"

//*****************************
// Path

func (h minpath) Len() int           { return len(h) }
func (h minpath) Less(i, j int) bool { return h[i].value < h[j].value }
func (h minpath) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minpath) Push(x interface{}) {
	*h = append(*h, x.(Path))
}

func (h *minpath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//*****************************
// Heap

func NewHeap() *Heap {
	return &Heap{values: &minpath{}}
}

func (h *Heap) Push(x Path) {
	heap.Push(h.values, x)
}

func (h *Heap) Pop() Path {
	return heap.Pop(h.values).(Path)
}
