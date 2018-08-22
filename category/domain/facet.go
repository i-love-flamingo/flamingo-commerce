package domain

import (
	"flamingo.me/flamingo-commerce/search/domain"
)

type (
	// CategoryFacet search filter
	CategoryFacet struct {
		Category Category
	}

	categoryKey string
)

var _ domain.Filter = CategoryFacet{}

const (
	// CategoryKey donates the default category facet key
	CategoryKey categoryKey = "category"
)

// NewCategoryFacet filter factory
func NewCategoryFacet(category Category) CategoryFacet {
	return CategoryFacet{
		Category: category,
	}
}

// Value for category/domain.Filter
func (cf CategoryFacet) Value() (string, []string) {
	return string(CategoryKey), []string{cf.Category.Code()}
}
