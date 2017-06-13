package domain

type (
	// Cart aggregate
	Cart struct {
		ID        int
		Cartitems []Cartitem
	}

	// Cartitem for Cart
	Cartitem struct {
		MarketplaceCode string
		Qty             int
		Currentprice    float32
	}

	// Totalitem for totals
	Totalitem struct {
		Type  string
		Price float32
	}
)

// Add item to cart
func (Cart *Cart) Add(Cartitem Cartitem) {
	Cart.Cartitems = append(Cart.Cartitems, Cartitem)
}

// GetLine gets an item - starting with 1
func (Cart *Cart) GetLine(lineNr int) Cartitem {
	return Cart.Cartitems[lineNr-1]
}

// AddOrUpdateByCode if cartitem with code is already in the cart its updated. Otherwise added
func (Cart *Cart) AddOrUpdateByCode(code string, qty int, price float32) {
	for id, cartItem := range Cart.Cartitems {
		if cartItem.MarketplaceCode == code {
			cartItem.Qty = cartItem.Qty + qty
			Cart.Cartitems[id] = cartItem
			return
		}
	}
	newCartItem := Cartitem{
		code,
		qty,
		price,
	}
	Cart.Cartitems = append(Cart.Cartitems, newCartItem)
}
