// dummyengine/pkg/customheap/heap.go

package customheap

import (
	"dummyengine/pkg/pricelevel"
)

type CommonHeap []*pricelevel.PriceLevel

func (h CommonHeap) Len() int {
	return len(h)
}

func (h CommonHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *CommonHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*pricelevel.PriceLevel)
	item.Index = n
	*h = append(*h, item)
}

func (h *CommonHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}
