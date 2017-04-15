package domain

//The cart aggregate
type (
	Cart struct {
		Id int
		Cartitems []Cartitem
	}

 	Cartitem struct {
		MarketplaceCode string
		Qty int
		Currentprice float32
	}

	Totalitem struct {
		Type string
		Price float32
	}
)

func (Cart *Cart) Add(Cartitem Cartitem) {
	Cart.Cartitems = append(Cart.Cartitems,Cartitem)
}
