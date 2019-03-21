package domain

import (
	"flamingo.me/flamingo-commerce/v3/search/domain"
)

type (
	// CategoryFacet search filter
	CategoryFacet struct {
		CategoryCode string
	}

	categoryKey string
)

var _ domain.Filter = CategoryFacet{}

const (
	// CategoryKey donates the default category facet key
	CategoryKey categoryKey = "category"
)

// NewCategoryFacet filter factory
func NewCategoryFacet(categoryCode string) CategoryFacet {
	return CategoryFacet{
		CategoryCode: categoryCode,
	}
}

// Value for category/domain.Filter
func (cf CategoryFacet) Value() (string, []string) {
	return string(CategoryKey), []string{cf.CategoryCode}
}
