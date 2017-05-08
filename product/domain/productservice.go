package domain

import (
	"flamingo/framework/web"
)

// ProductService interface
type ProductService interface {
	// Get a product
	Get(ctx web.Context, foreignId string) (*Product, error)
}
