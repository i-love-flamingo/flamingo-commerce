package domain

import (
	"context"
	"fmt"

	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name ProductService --case snake

type (
	// ProductService interface
	ProductService interface {
		// Get a product
		Get(ctx context.Context, marketplaceCode string) (BasicProduct, error)
	}

	// SearchResult returns product hits
	SearchResult struct {
		searchDomain.Result
		Hits []BasicProduct
	}

	// SearchService is a typed search for products
	SearchService interface {
		//Search returns Products based on given Filters
		Search(ctx context.Context, filter ...searchDomain.Filter) (*SearchResult, error)
		// SearchBy returns Products prefiltered by the given attribute (also based on additional given Filters)
		// e.g. SearchBy(ctx,"brandCode","apple")
		SearchBy(ctx context.Context, attribute string, values []string, filter ...searchDomain.Filter) (*SearchResult, error)
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
