package categorydto

import (
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
)

// CategorySearchResult represents category search result
type CategorySearchResult struct {
	ProductSearchResult *graphql.SearchResultDTO
	Category            domain.Category
}
