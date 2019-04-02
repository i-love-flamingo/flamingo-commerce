package application

import (
	"context"
	"math"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// RestrictionService checks product restriction
	RestrictionService struct {
		qtyRestrictors []cart.MaxQuantityRestrictor
	}

	// ErrNoRestriction is used to indicate that there is no restriction
	ErrNoRestriction struct{}
)

// Inject dependencies
func (rs *RestrictionService) Inject(
	qtyRestrictors []cart.MaxQuantityRestrictor,
) *RestrictionService {
	rs.qtyRestrictors = qtyRestrictors

	return rs
}

func (ErrNoRestriction) Error() string {
	return "qty is not restricted"
}

// RestrictQty checks if there is an allowed max qty, ErrNoRestriction is returned if there is no qta restriction at all for the  given product
func (rs *RestrictionService) RestrictQty(ctx context.Context, product domain.BasicProduct, cart *cart.Cart) (int, error) {
	var maximumAllowed = math.MaxInt32
	for _, r := range rs.qtyRestrictors {
		currentMax := r.Restrict(ctx, product, cart)
		if currentMax < maximumAllowed {
			maximumAllowed = currentMax
		}
	}

	if maximumAllowed == math.MaxInt32 {
		return 0, &ErrNoRestriction{}
	}

	return maximumAllowed, nil

}
