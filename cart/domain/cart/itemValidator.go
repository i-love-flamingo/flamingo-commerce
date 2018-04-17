package cart

import (
	"fmt"

	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/web"
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
