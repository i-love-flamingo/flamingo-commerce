package cart

import (
	"context"

	domainCart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/product/domain"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// FakeQtyRestrictor used to restrict
	FakeQtyRestrictor struct{}
)

var (
	_ validation.MaxQuantityRestrictor = &FakeQtyRestrictor{}
)

// Name fake implementation
func (f FakeQtyRestrictor) Name() string {
	return "Name"
}

// Restrict fake implementation
func (f FakeQtyRestrictor) Restrict(ctx context.Context, session *web.Session, product domain.BasicProduct, cart *domainCart.Cart, deliveryCode string) *validation.RestrictionResult {
	if product.BaseData().MarketPlaceCode == "fake_simple" {
		return &validation.RestrictionResult{
			IsRestricted:        true,
			MaxAllowed:          10,
			RemainingDifference: 10,
			RestrictorName:      f.Name(),
		}
	}
	return &validation.RestrictionResult{
		IsRestricted: false,
	}
}
