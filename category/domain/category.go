package domain

import (
	"flamingo.me/flamingo-commerce/search/domain"
)

type (
	// Category domain model
	Category interface {
		Code() string
		Name() string
		Path() string
		Categories() []Category
		Active() bool
	}

	// CategoryFacet search filter
	CategoryFacet struct {
		Category Category
	}

	categoryKey string
)

const (
	// CategoryKey donates the default category facet key
	CategoryKey categoryKey = "category"
)

var (
	_ domain.Filter = CategoryFacet{}
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
