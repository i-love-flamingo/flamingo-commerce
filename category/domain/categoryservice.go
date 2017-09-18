package domain

import (
	"context"
	"errors"
	"flamingo/core/product/domain"
)

var (
	// NotFound error
	NotFound = errors.New("category not found")
)

type (
	// CategoryService interface
	CategoryService interface {
		// Get a category
		Get(ctx context.Context, categoryCode string) (Category, error)

		// GetProducts for a given category
		GetProducts(ctx context.Context, categoryCode string) ([]domain.BasicProduct, error)
	}
)
