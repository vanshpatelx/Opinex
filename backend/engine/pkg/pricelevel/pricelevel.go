// dummyengine/pkg/pricelevel/pricelevel.go
package pricelevel

import (
	"math/big"
)

type Order struct {
	ID       *big.Int
	Price    int
	Quantity int
	UserID   *big.Int
}

type PriceLevel struct {
	Price    int
	Quantity int
	Index    int // Heap index
	Orders   []*Order
}
