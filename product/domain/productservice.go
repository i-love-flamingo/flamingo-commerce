package domain

import "context"

type (
	// ProductService interface
	ProductService interface {
		// Get a product
		Get(ctx context.Context, foreignID string) (*Product, error)
	}

	// ProductNotFound is an error
	ProductNotFound struct {
		ID string
	}
)

// Error implements the error interface
func (b ProductNotFound) Error() string {
	return "Product with ID " + b.ID + " Not Found"
}
