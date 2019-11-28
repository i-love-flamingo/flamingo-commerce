package dto

import (
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/application"
)

// CategorySearchResult represents category search result
type CategorySearchResult struct {
	ProductSearchResult *application.SearchResult
	Category            domain.Category
}
