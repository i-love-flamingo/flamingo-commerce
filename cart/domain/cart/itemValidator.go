package cart

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"github.com/gorilla/sessions"
)

type (
	ItemValidator interface {
		Validate(ctx context.Context, session *sessions.Session, deliveryCode string, request AddRequest, product domain.BasicProduct) error
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
