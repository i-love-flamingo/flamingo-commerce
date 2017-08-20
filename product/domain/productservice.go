package domain

import (
	"context"
	"fmt"
)

type (
	// ProductService interface
	ProductService interface {
		// Get a product
		Get(ctx context.Context, marketplaceCode string) (BasicProduct, error)
	}

	// ProductNotFound is an error
	ProductNotFound struct {
		MarketplaceCode string
	}
)

// Error implements the error interface
func (err ProductNotFound) Error() string {
	return fmt.Sprintf("Product with Marketplace Code %q Not Found", err.MarketplaceCode)
}
