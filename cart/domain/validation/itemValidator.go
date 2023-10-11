package validation

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name ItemValidator --case snake

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
		AdditionalData      map[string]interface{}
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
