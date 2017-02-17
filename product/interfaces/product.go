package interfaces

import "flamingo/core/product/models"

// ProductService interface
type ProductService interface {
	Get(string) models.Product
	GetByIDList([]string) []models.Product
}
