package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
)

// DecoratedCart – provides custom graphql interface methods
type DecoratedCart struct {
	decoratedCart *decorator.DecoratedCart
}

// Cart – provides the cart
func (dc DecoratedCart) Cart() cart.Cart {
	return dc.decoratedCart.Cart
}

// DecoratedDeliveries – returns decorated deliveries
func (dc DecoratedCart) DecoratedDeliveries() []decorator.DecoratedDelivery {
	return dc.decoratedCart.DecoratedDeliveries
}

// GetDecoratedDeliveryByCode – returns decorated delivery filtered by code
func (dc *DecoratedCart) GetDecoratedDeliveryByCode(deliveryCode string) *decorator.DecoratedDelivery {
	decoratedDelivery, _ := dc.decoratedCart.GetDecoratedDeliveryByCode(deliveryCode)
	return decoratedDelivery
}

// GetAllPaymentRequiredItems – returns all payment required items
func (dc *DecoratedCart) GetAllPaymentRequiredItems() PricedItems {
	dcCart := dc.Cart()
	return PricedItems{items: dcCart.GetAllPaymentRequiredItems()}
}

// CartSummary – returns cart summary
func (dc *DecoratedCart) CartSummary() CartSummary {
	dcCart := dc.Cart()
	return CartSummary{cart: &dcCart}
}

// NewDecoratedCart – factory method
func NewDecoratedCart(dc *decorator.DecoratedCart) *DecoratedCart {
	return &DecoratedCart{decoratedCart: dc}
}
