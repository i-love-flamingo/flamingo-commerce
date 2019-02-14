package cart

import (
	"context"
	"flamingo.me/flamingo/v3/framework/web"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	ItemValidator interface {
		Validate(ctx context.Context, session *web.Session, deliveryCode string, request AddRequest, product domain.BasicProduct) error
	}

	AddToCartNotAllowed struct {
		Reason              string
		RedirectHandlerName string
		RedirectParams      map[string]string
	}
)

func (e *AddToCartNotAllowed) Error() string {
	return fmt.Sprintf("Product is not allowed: %v", e.Reason)
}

func (e *AddToCartNotAllowed) MessageCode() string {
	return e.Reason
}
