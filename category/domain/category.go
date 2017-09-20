package domain

import (
	"flamingo/core/search/domain"
	"net/url"
)

type (
	// Category domain model
	Category interface {
		Code() string
		Name() string
		Categories() []Category
	}

	// CategoryFacet search filter
	CategoryFacet struct {
		category string
	}

	categoryKey string
)

const (
	CategoryKey categoryKey = "category"
)

var (
	_ domain.Filter = new(CategoryFacet)
)

// NewCategoryFacet filter factory
func NewCategoryFacet(category string) *CategoryFacet {
	return &CategoryFacet{
		category: category,
	}
}

// Values for category/domain.Filter
func (cf *CategoryFacet) Values() url.Values {
	return url.Values{
		string(CategoryKey): {cf.category},
	}
}
