package decorator

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

const (
	decoratedDiscountError = "Unable to collect discounts, stopping and returning empty slice"
)

type (
	// DecoratedWithDiscount interface for a decorated object to be able to handle discounts
	// the difference to cart.WithDiscount is, that these functions do NOT provide the client
	// with an error, errors are just logged
	DecoratedWithDiscount interface {
		HasAppliedDiscounts() bool
		MergeDiscounts() cart.AppliedDiscounts
	}
)

var (
	// interface assertions
	_ DecoratedWithDiscount = new(DecoratedCartItem)
	_ DecoratedWithDiscount = new(DecoratedDelivery)
	_ DecoratedWithDiscount = new(DecoratedCart)
)

// HasAppliedDiscounts checks whether decorated item has discounts
func (dci *DecoratedCartItem) HasAppliedDiscounts() bool {
	return hasAppliedDiscounts(dci)
}

// MergeDiscounts sum up discounts applied to item
func (dci DecoratedCartItem) MergeDiscounts() cart.AppliedDiscounts {
	return collectDiscounts(&dci.Item, dci.logger)
}

// HasAppliedDiscounts checks whether decorated delivery has discounts
func (dc DecoratedDelivery) HasAppliedDiscounts() bool {
	return hasAppliedDiscounts(dc)
}

// MergeDiscounts sum up discounts applied to delivery
func (dc DecoratedDelivery) MergeDiscounts() cart.AppliedDiscounts {
	return collectDiscounts(&dc.Delivery, dc.logger)
}

// HasAppliedDiscounts checks whether decorated cart has discounts
func (dc DecoratedCart) HasAppliedDiscounts() bool {
	return hasAppliedDiscounts(dc)
}

// MergeDiscounts sum up discounts applied to cart
func (dc DecoratedCart) MergeDiscounts() cart.AppliedDiscounts {
	return collectDiscounts(&dc.Cart, dc.Logger)
}

// private helpers as all implementations of decorated entities are based on underlying
// interface cart.WithDiscount and can therefore be abstracted

func hasAppliedDiscounts(discountable DecoratedWithDiscount) bool {
	return len(discountable.MergeDiscounts()) > 0
}

func collectDiscounts(discountable cart.WithDiscount, logger flamingo.Logger) cart.AppliedDiscounts {
	discounts, err := discountable.MergeDiscounts()
	if err != nil {
		logger.Error(decoratedDiscountError)
		return make(cart.AppliedDiscounts, 0)
	}
	return discounts
}
