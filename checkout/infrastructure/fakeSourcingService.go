package infrastructure

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"

	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

// FakeSourcingService represents the fake source locator
type FakeSourcingService struct{}

var _ domain.SourcingService = new(FakeSourcingService)

// GetSourceID provides fake delivery locations
func (sl *FakeSourcingService) GetSourceID(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, item *decorator.DecoratedCartItem) (string, error) {
	return "mock_ispu_location1", nil
}
