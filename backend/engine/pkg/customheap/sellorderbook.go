// dummyengine/pkg/customheap/sellorderbook.go

package customheap

type SellOrderBook struct{ 
	CommonHeap 
}

func (s *SellOrderBook) Less(i, j int) bool { 
	if i >= len(s.CommonHeap) || j >= len(s.CommonHeap) {
        return false // Prevent out-of-bounds error
    }
	return s.CommonHeap[i].Price < s.CommonHeap[j].Price 
}
