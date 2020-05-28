package validation

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"

	"fmt"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// ItemValidator checks a cart item
	ItemValidator interface {
		Validate(ctx context.Context, session *web.Session, cart *decorator.DecoratedCart, deliveryCode string, request cart.AddRequest, product domain.BasicProduct) error
	}

	// AddToCartNotAllowed error
	AddToCartNotAllowed struct {
		Reason              string
		RedirectHandlerName string
		RedirectParams      map[string]string
	}
)

// Error message
func (e *AddToCartNotAllowed) Error() string {
	return fmt.Sprintf("Product is not allowed: %v", e.Reason)
}

// MessageCode message code
func (e *AddToCartNotAllowed) MessageCode() string {
	return e.Reason
}
