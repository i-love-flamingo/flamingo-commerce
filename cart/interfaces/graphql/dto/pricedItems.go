package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// PricedItems – provides custom graphql interface methods
	PricedItems struct {
		items cart.PricedItems
	}

	// PricedShippingItem – shipping item with price
	PricedShippingItem struct {
		Amount           domain.Price
		DeliveryInfoCode string
	}

	// PricedTotalItem – total item with price
	PricedTotalItem struct {
		Amount domain.Price
		Code   string
	}

	// PricedCartItem – cart item with price
	PricedCartItem struct {
		Amount domain.Price
		ItemID string
	}
)

// CartItems – return all cart items
func (pr PricedItems) CartItems() []PricedCartItem {
	var items []PricedCartItem

	for CartItemID, price := range pr.items.CartItems() {
		items = append(items, PricedCartItem{
			Amount: price,
			ItemID: CartItemID,
		})
	}

	return items
}

// TotalItems – return all total items
func (pr PricedItems) TotalItems() []PricedTotalItem {
	var items []PricedTotalItem

	for totalItemCode, price := range pr.items.TotalItems() {
		items = append(items, PricedTotalItem{
			Amount: price,
			Code:   totalItemCode,
		})
	}

	return items
}

// ShippingItems – return all shipping items
func (pr PricedItems) ShippingItems() []PricedShippingItem {
	var items []PricedShippingItem

	for deliveryInfoCode, price := range pr.items.ShippingItems() {
		items = append(items, PricedShippingItem{
			Amount:           price,
			DeliveryInfoCode: deliveryInfoCode,
		})
	}

	return items
}
