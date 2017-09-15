package domain

import (
	"context"
	"flamingo/core/product/domain"
)

type (
	// CategoryService interface
	CategoryService interface {
		// Get a category
		Get(ctx context.Context, categoryCode string) (Category, error)

		GetProducts(ctx context.Context, categoryCode string) ([]domain.BasicProduct, error)
	}
)
