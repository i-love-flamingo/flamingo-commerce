package cart

import (
	"context"
	"errors"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

// FakeItemValidatorCookie name to control behaviour
const FakeItemValidatorCookie = "X-FakeItemValidator"

type (
	// FakeItemValidator returns an AddToCartNotAllowed error if the Cookie FakeItemValidatorCookie is set
	FakeItemValidator struct{}
)

var _ validation.ItemValidator = FakeItemValidator{}

// Validate is only a fake implementation which is controlled by the Cookie FakeItemValidatorCookie.
// Always returns an error if the cookie is set
func (f FakeItemValidator) Validate(ctx context.Context, _ *web.Session, _ *decorator.DecoratedCart, _ string, _ cart.AddRequest, _ domain.BasicProduct) error {
	r := web.RequestFromContext(ctx)

	_, err := r.Request().Cookie(FakeItemValidatorCookie)
	if errors.Is(err, http.ErrNoCookie) {
		return nil
	}

	return &validation.AddToCartNotAllowed{Reason: "fake item validator reason"}
}
