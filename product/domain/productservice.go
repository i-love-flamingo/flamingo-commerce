package domain

import "context"

// ProductService interface
type ProductService interface {
	// Get a product
	Get(ctx context.Context, foreignID string) (*Product, error)
}
