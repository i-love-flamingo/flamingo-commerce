package domain

type (
	// Cart Value Object (immutable - because the cartservice is responsible to return a cart).
	Cart struct {
		ID        int
		Cartitems []Cartitem
	}

	// Cartitem for Cart
	Cartitem struct {
		ProductCode  string
		ProductName  string
		Qty          int
		Currentprice float32
	}

	// Totalitem for totals
	Totalitem struct {
		Type  string
		Price float32
	}
)

// GetLine gets an item - starting with 1
func (Cart *Cart) GetLine(lineNr int) Cartitem {
	return Cart.Cartitems[lineNr-1]
}
