package cart

import (
	"fmt"

	"flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	ItemValidator interface {
		Validate(ctx web.Context, request AddRequest, product domain.BasicProduct) error
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
