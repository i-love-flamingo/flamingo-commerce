package cart

import (
	"github.com/pkg/errors"
)

type (
	// Cart Value Object (immutable - because the cartservice is responsible to return a cart).
	Cart struct {
		ID             string
		Cartitems      []Item
		Totalitems     []Totalitem
		ShippingItem   ShippingItem
		GrandTotal     float64
		SubTotal       float64
		DiscountAmount float64
		TaxAmount      float64

		CurrencyCode string
		//Intention is optional and expresses the intented use case for this cart - it is used when multiple carts are used to distinguish between them
		Intention string
	}

	// Item for Cart
	Item struct {
		ID              int
		MarketplaceCode string
		//VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		Price float64
		Qty   int

		RowTotal       float64
		TaxAmount      float64
		DiscountAmount float64

		PriceInclTax    float64
		RowTotalInclTax float64
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price float64
	}

	// ShippingItem
	ShippingItem struct {
		Title string
		Price float64

		TaxAmount      float64
		DiscountAmount float64
	}
)

// GetByLineNr gets an item - starting with 1
func (Cart *Cart) GetByLineNr(lineNr int) (*Item, error) {
	var item Item
	if len(Cart.Cartitems) >= lineNr && lineNr > 0 {
		return &Cart.Cartitems[lineNr-1], nil
	} else {
		return &item, errors.New("Line in cart not existend")
	}
}

// HasItem checks if a cartitem for that sku exists and returns lineNr if found
func (Cart *Cart) HasItem(marketplaceCode string, variantMarketplaceCode string) (bool, int) {
	for lineNr, item := range Cart.Cartitems {
		if item.MarketplaceCode == marketplaceCode && item.VariantMarketPlaceCode == variantMarketplaceCode {
			return true, lineNr + 1
		}
	}
	return false, 0
}
