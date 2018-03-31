package infrastructure

import (
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/checkout/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// FakeDeliveryLocationsService represents the fake source locator
	FakeSourcingService struct {
	}
)

var (
	_ domain.SourcingService = new(FakeSourcingService)
)

// GetDeliveryLocations provides fake delivery locations
func (sl *FakeSourcingService) GetSourceId(ctx web.Context, decoratedCart *cartDomain.DecoratedCart, item *cartDomain.DecoratedCartItem) (string, error) {
	return "mock_ispu_location1", nil
}
