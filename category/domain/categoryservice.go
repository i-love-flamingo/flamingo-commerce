package domain

import (
	"context"
	"errors"
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
	}
)
