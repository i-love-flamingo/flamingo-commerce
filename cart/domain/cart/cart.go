package cart

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

type (
	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		CartOrderBehaviour CartOrderBehaviour `json:"-"`
		ID                 string
		Cartitems          []Item
		Totalitems         []Totalitem
		ShippingItem       ShippingItem
		GrandTotal         float64
		SubTotal           float64
		DiscountAmount     float64
		TaxAmount          float64
		//TODO - move to Item?
		CurrencyCode string
		//Intention is optional and expresses the intented use case for this cart - it is used when multiple carts are used to distinguish between them
		Intention string
	}

	// Item for Cart
	Item struct {
		ID              string
		MarketplaceCode string
		//VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of Ispu Location or Collection Point
		SourceId string

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
func (Cart Cart) GetByLineNr(lineNr int) (*Item, error) {
	var item Item
	if len(Cart.Cartitems) >= lineNr && lineNr > 0 {
		return &Cart.Cartitems[lineNr-1], nil
	} else {
		return &item, errors.New("Line in cart not existend")
	}
}

// GetByItemId gets an item by its id
func (Cart Cart) GetByItemId(itemId string) (*Item, error) {
	for _, currentItem := range Cart.Cartitems {
		log.Println(currentItem.ID)
		if currentItem.ID == itemId {
			return &currentItem, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("itemId %v in cart not existend", itemId))
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

//HasShippingItem
func (cart Cart) HasShippingItem() bool {
	if cart.ShippingItem.Title != "" {
		return true
	}
	return false
}

// SetShippingInformation
func (cart Cart) SetShippingInformation(ctx context.Context, auth Auth, shippingAddress *Address, billingAddress *Address, shippingCarrierCode string, shippingMethodCode string) error {
	if cart.CartOrderBehaviour == nil {
		return errors.New("This Cart has no Behaviour attached!")
	}
	return cart.CartOrderBehaviour.SetShippingInformation(ctx, auth, &cart, shippingAddress, billingAddress, shippingCarrierCode, shippingMethodCode)
}

// PlaceOrder
func (cart Cart) PlaceOrder(ctx context.Context, auth Auth, payment *Payment) (string, error) {
	if cart.CartOrderBehaviour == nil {
		return "", errors.New("This Cart has no Behaviour attached!")
	}
	return cart.CartOrderBehaviour.PlaceOrder(ctx, auth, &cart, payment)
}

// DeleteItem
func (cart Cart) DeleteItem(ctx context.Context, auth Auth, id string) error {
	if cart.CartOrderBehaviour == nil {
		return errors.New("This Cart has no Behaviour attached!")
	}
	return cart.CartOrderBehaviour.DeleteItem(ctx, auth, &cart, id)
}

// UpdateItemQty - delete item if qty =< 0
func (cart Cart) UpdateItemQty(ctx context.Context, auth Auth, id string, qty int) error {
	if cart.CartOrderBehaviour == nil {
		return errors.New("This Cart has no Behaviour attached!")
	}
	item, e := cart.GetByItemId(id)
	if e != nil {
		return e
	}
	if qty < 1 {
		return cart.DeleteItem(ctx, auth, id)
	}
	item.Qty = qty
	return cart.CartOrderBehaviour.UpdateItem(ctx, auth, &cart, id, *item)
}

// UpdateItem replaces value in Cart Item
func (cart Cart) UpdateItem(ctx context.Context, auth Auth, item Item) error {
	return cart.CartOrderBehaviour.UpdateItem(ctx, auth, &cart, item.ID, item)
}

// ItemCount - returns amount of Cartitems
func (Cart Cart) ItemCount() int {
	return len(Cart.Cartitems)
}
