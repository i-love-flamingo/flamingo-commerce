package cart

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// RestrictionResult contains the result of a restriction
	RestrictionResult struct {
		IsRestricted        bool
		MaxAllowed          int
		RemainingDifference int
	}

	// MaxQuantityRestrictor returns the maximum qty allowed for a given product and cart
	MaxQuantityRestrictor interface {
		// Restrict must return a `RestrictionResult` which contains information regarding if a restriction is
		// applied and whats the max allowed quantity
		Restrict(ctx context.Context, product domain.BasicProduct, cart *Cart) *RestrictionResult
	}
)
