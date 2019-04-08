package cart

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// MaxQuantityRestrictor returns the maximum qty allowed for a given product and cart
	MaxQuantityRestrictor interface {
		// Restrict must return the maximum allowed qty or `math.MaxInt32` for infinity
		Restrict(ctx context.Context, product domain.BasicProduct, cart *Cart) int
	}
)
