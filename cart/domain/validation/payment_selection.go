package validation

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
)

type (
	// PaymentSelectionValidator decides if the PaymentSelection is valid
	PaymentSelectionValidator interface {
		Validate(ctx context.Context, cart *decorator.DecoratedCart, selection cart.PaymentSelection) error
	}
)
