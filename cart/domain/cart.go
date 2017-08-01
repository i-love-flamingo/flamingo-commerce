package domain

import "github.com/pkg/errors"

type (
	// Cart Value Object (immutable - because the cartservice is responsible to return a cart).
	Cart struct {
		ID             int
		Cartitems      []Cartitem
		Totalitems     []Cartitem
		ShippingItem   ShippingItem
		GrandTotal     float32
		SubTotal       float32
		DiscountAmount float32
		TaxAmount      float32

		CurrencyCode string
	}

	// Cartitem for Cart
	Cartitem struct {
		ID          int
		ProductCode string
		ProductName string

		Price float32
		Qty   int

		RowTotal       float32
		TaxAmount      float32
		DiscountAmount float32

		PriceInclTax    float32
		RowTotalInclTax float32
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price float32
	}

	// ShippingItem
	ShippingItem struct {
		Title string
		Price float32

		TaxAmount      float32
		DiscountAmount float32
	}
)

// GetLine gets an item - starting with 1
func (Cart *Cart) GetLine(lineNr int) (Cartitem, error) {
	var item Cartitem
	if len(Cart.Cartitems) > lineNr {
		return Cart.Cartitems[lineNr-1], nil
	} else {
		return item, errors.New("Line in cart not existend")
	}
}
