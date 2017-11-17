package cart

import (
	"context"

	"github.com/pkg/errors"
)

type (
	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		CartOrderBehaviour CartOrderBehaviour
		ID                 string
		Cartitems          []Item
		Totalitems         []Totalitem
		ShippingItem       ShippingItem
		GrandTotal         float64
		SubTotal           float64
		DiscountAmount     float64
		TaxAmount          float64

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

	//CartBehaviour is a Port that can be implemented by other packages to implement  cart actions required for Ordering a Cart
	CartOrderBehaviour interface {
		PlaceOrder(ctx context.Context, cart *Cart, payment *Payment) (string, error)
		SetShippingInformation(ctx context.Context, cart *Cart, shippingAddress *Address, billingAddress *Address, shippingCarrierCode string, shippingMethodCode string) error
	}
)

// GetByLineNr gets an item - starting with 1
func (Cart Cart) GetByLineNr(lineNr int) (*Item, error) {
	var item Item
	if len(Cart.Cartitems) >= lineNr && lineNr > 0 {
		return &Cart.Cartitems[lineNr-1], nil
	} else {
		return &item, errors.New("Line in cart not existend")
	}
}

// HasItem checks if a cartitem for that sku exists and returns lineNr if found
func (cart Cart) HasItem(marketplaceCode string, variantMarketplaceCode string) (bool, int) {
	for lineNr, item := range cart.Cartitems {
		if item.MarketplaceCode == marketplaceCode && item.VariantMarketPlaceCode == variantMarketplaceCode {
			return true, lineNr + 1
		}
	}
	return false, 0
}

// SetShippingInformation
func (cart Cart) SetShippingInformation(ctx context.Context, shippingAddress *Address, billingAddress *Address, shippingCarrierCode string, shippingMethodCode string) error {
	if cart.CartOrderBehaviour == nil {
		return errors.New("This Cart has no Behaviour attached!")
	}
	return cart.CartOrderBehaviour.SetShippingInformation(ctx, &cart, shippingAddress, billingAddress, shippingCarrierCode, shippingMethodCode)
}

// PlaceOrder
func (cart Cart) PlaceOrder(ctx context.Context, payment *Payment) (string, error) {
	if cart.CartOrderBehaviour == nil {
		return "", errors.New("This Cart has no Behaviour attached!")
	}
	return cart.CartOrderBehaviour.PlaceOrder(ctx, &cart, payment)
}

// ItemCount - returns amount of Cartitems
func (Cart Cart) ItemCount() int {
	return len(Cart.Cartitems)
}
