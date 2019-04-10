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
)

// Inject dependencies
func (rs *RestrictionService) Inject(
	qtyRestrictors []cart.MaxQuantityRestrictor,
) *RestrictionService {
	rs.qtyRestrictors = qtyRestrictors

	return rs
}

// RestrictQty checks if there is an qty restriction present and returns an according result containing the max allowed
// quantity and the quantity difference to the current cart
func (rs *RestrictionService) RestrictQty(ctx context.Context, product domain.BasicProduct, currentCart *cart.Cart) *cart.RestrictionResult {
	restrictionResult := &cart.RestrictionResult{
		IsRestricted:        false,
		RemainingDifference: math.MaxInt32,
		MaxAllowed:          math.MaxInt32,
	}

	for _, r := range rs.qtyRestrictors {
		currentResult := r.Restrict(ctx, product, currentCart)

		if currentResult.IsRestricted {
			restrictionResult.IsRestricted = true

			if currentResult.MaxAllowed < restrictionResult.MaxAllowed {
				restrictionResult.MaxAllowed = currentResult.MaxAllowed
			}

			if currentResult.RemainingDifference < restrictionResult.RemainingDifference {
				restrictionResult.RemainingDifference = currentResult.RemainingDifference
			}
		}
	}

	return restrictionResult
}
