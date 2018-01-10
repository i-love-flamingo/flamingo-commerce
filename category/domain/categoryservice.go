package domain

import (
	"context"
	"errors"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New("category not found")
)

type (
	// CategoryService interface
	CategoryService interface {
		// Tree a category
		Tree(ctx context.Context, categoryCode string) (Category, error)

		// Get a category with more data
		Get(ctx context.Context, categoryCode string) (Category, error)
	}
)
