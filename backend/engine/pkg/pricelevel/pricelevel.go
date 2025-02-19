// dummyengine/pkg/pricelevel/pricelevel.go
package pricelevel


import (
	"math/big"
)
type Order struct {
	ID       *big.Int
	Price    int
	Quantity int
}

type PriceLevel struct {
	Price    int
	Quantity int
	Index    int // Heap index
	Orders   []*Order
}
