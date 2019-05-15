package validation

import (
	"context"
	"math"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// RestrictionService checks product restriction
	RestrictionService struct {
		qtyRestrictors []MaxQuantityRestrictor
	}

	// RestrictionResult contains the result of a restriction
	RestrictionResult struct {
		IsRestricted        bool
		MaxAllowed          int
		RemainingDifference int
		RestrictorName      string
	}

	// MaxQuantityRestrictor returns the maximum qty allowed for a given product and cart
	// it is possible to register many (MultiBind) MaxQuantityRestrictor implementations
	MaxQuantityRestrictor interface {
		// Name returns the code of the restrictor
		Name() string
		// Restrict must return a `RestrictionResult` which contains information regarding if a restriction is
		// applied and whats the max allowed quantity
		Restrict(ctx context.Context, product domain.BasicProduct, cart *cart.Cart) *RestrictionResult
	}
)

// Inject dependencies
func (rs *RestrictionService) Inject(
	qtyRestrictors []MaxQuantityRestrictor,
) *RestrictionService {
	rs.qtyRestrictors = qtyRestrictors

	return rs
}

// RestrictQty checks if there is an qty restriction present and returns an according result containing the max allowed
// quantity and the quantity difference to the current cart
func (rs *RestrictionService) RestrictQty(ctx context.Context, product domain.BasicProduct, currentCart *cart.Cart) *RestrictionResult {
	restrictionResult := &RestrictionResult{
		IsRestricted:        false,
		MaxAllowed:          math.MaxInt32,
		RemainingDifference: math.MaxInt32,
	}

	for _, r := range rs.qtyRestrictors {
		currentResult := r.Restrict(ctx, product, currentCart)

		if currentResult.IsRestricted {
			restrictionResult.IsRestricted = true

			if currentResult.MaxAllowed < restrictionResult.MaxAllowed {
				restrictionResult.MaxAllowed = currentResult.MaxAllowed
				restrictionResult.RestrictorName = currentResult.RestrictorName
			}

			if currentResult.RemainingDifference < restrictionResult.RemainingDifference {
				restrictionResult.RemainingDifference = currentResult.RemainingDifference
			}
		}
	}

	return restrictionResult
}
