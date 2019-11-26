package dto

import (
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/application"
)

type CategorySearchResult struct {
	ProductSearchResult *application.SearchResult
	Category            domain.Category
}
