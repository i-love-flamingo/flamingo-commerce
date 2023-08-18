package domain

import (
	"context"
	"errors"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New("category not found")
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name CategoryService --case snake

type (
	// CategoryService interface
	CategoryService interface {
		// Tree a category
		Tree(ctx context.Context, activeCategoryCode string) (Tree, error)

		// Get a category with more data
		Get(ctx context.Context, categoryCode string) (Category, error)
	}
)
