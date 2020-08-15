package cart

import (
	"context"
	"errors"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
)

// FakePaymentSelectionValidatorCookie name to control behaviour
const FakePaymentSelectionValidatorCookie = "X-FakePaymentSelectionValidator"

type (
	// FakePaymentSelectionValidator returns an error if the Cookie FakePaymentSelectionValidatorCookie is set
	FakePaymentSelectionValidator struct{}
)

// Validate is only a fake implementation which is controlled by the Cookie FakePaymentSelectionValidatorCookie.
// Always returns an error if the cookie is set
func (f FakePaymentSelectionValidator) Validate(ctx context.Context, _ *decorator.DecoratedCart, _ cart.PaymentSelection) error {
	r := web.RequestFromContext(ctx)
	_, err := r.Request().Cookie(FakePaymentSelectionValidatorCookie)
	if errors.Is(err, http.ErrNoCookie) {
		return nil
	}

	return errors.New("fake payment selection validator error")
}
