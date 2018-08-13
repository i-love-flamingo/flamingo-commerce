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
	// categoryService interface
	CategoryService interface {
		// Tree of a category
		Tree(ctx context.Context, categoryCode string) (Category, error)

		// Children of a category
		Children(ctx context.Context, categoryCode string) (Category, error)

		// Get a category with more data
		Get(ctx context.Context, categoryCode string) (Category, error)
	}
)
