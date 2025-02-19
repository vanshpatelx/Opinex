// dummyengine/pkg/customheap/buyorderbook.go

package customheap

type BuyOrderBook struct {
	CommonHeap
}

func (b *BuyOrderBook) Less(i, j int) bool {
    if i >= len(b.CommonHeap) || j >= len(b.CommonHeap) {
        return false // Prevent out-of-bounds error
    }
    return b.CommonHeap[i].Price > b.CommonHeap[j].Price
}

