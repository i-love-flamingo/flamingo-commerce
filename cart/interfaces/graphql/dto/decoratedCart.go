package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
)

type (
	// DecoratedCart – provides custom graphql interface methods
	DecoratedCart struct {
		decoratedCart *decorator.DecoratedCart
	}

	// DecoratedDelivery Decorates a CartItem with its Product
	DecoratedDelivery struct {
		Delivery       cart.Delivery
		DecoratedItems []DecoratedCartItem
	}

	// DecoratedCartItem Decorates a CartItem with its Product
	DecoratedCartItem struct {
		Item    cart.Item
		Product graphqlProductDto.Product
	}
)

// Cart – provides the cart
func (dc DecoratedCart) Cart() cart.Cart {
	return dc.decoratedCart.Cart
}

// DecoratedCart – provides the cart
func (dc DecoratedCart) DecoratedCart() *decorator.DecoratedCart {
	return dc.decoratedCart
}

// DecoratedDeliveries – returns decorated deliveries
func (dc DecoratedCart) DecoratedDeliveries() []DecoratedDelivery {
	return mapDecoratedDeliveries(dc.decoratedCart.DecoratedDeliveries)
}

// GetDecoratedDeliveryByCode – returns decorated delivery filtered by code
func (dc *DecoratedCart) GetDecoratedDeliveryByCode(deliveryCode string) *DecoratedDelivery {
	decoratedDelivery, _ := dc.decoratedCart.GetDecoratedDeliveryByCode(deliveryCode)
	return &DecoratedDelivery{
		Delivery:       decoratedDelivery.Delivery,
		DecoratedItems: mapDecoratedItems(decoratedDelivery.DecoratedItems),
	}
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

// mapDecoratedDeliveries
func mapDecoratedDeliveries(decoratedDeliveries []decorator.DecoratedDelivery) []DecoratedDelivery {
	if len(decoratedDeliveries) == 0 {
		return nil
	}

	deliveries := make([]DecoratedDelivery, 0, len(decoratedDeliveries))

	for _, dd := range decoratedDeliveries {
		deliveries = append(deliveries, DecoratedDelivery{
			Delivery:       dd.Delivery,
			DecoratedItems: mapDecoratedItems(dd.DecoratedItems),
		})
	}

	return deliveries
}

func mapDecoratedItems(decoratedItems []decorator.DecoratedCartItem) []DecoratedCartItem {
	items := make([]DecoratedCartItem, 0, len(decoratedItems))

	for _, di := range decoratedItems {
		items = append(items, DecoratedCartItem{
			Item:    di.Item,
			Product: graphqlProductDto.NewGraphqlProductDto(di.Product, &di.Item.VariantMarketPlaceCode, nil),
		})
	}

	return items
}
