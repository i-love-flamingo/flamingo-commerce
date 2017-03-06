package interfaces

import (
	"flamingo/core/flamingo/web"
	"flamingo/core/product/models"
)

// ProductService interface
type ProductService interface {
	Get(web.Context, string) (models.Product, error)
	GetByIDList(web.Context, []string) []models.Product
}
