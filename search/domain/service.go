package domain

import (
	"context"
	"errors"
	"net/url"

	"go.aoe.com/flamingo/core/product/domain"
)

type (
	// Filter interface for search queries
	Filter interface {
		Values() url.Values
	}

	// KeyValueFilter allows simple k -> []values filtering
	KeyValueFilter struct {
		k string
		v []string
	}

	// SearchMeta data
	SearchMeta struct {
		Page       int
		NumPages   int
		NumResults int
	}

	// SearchService defines how to access search
	SearchService interface {
		GetProducts(
			ctx context.Context,
			searchMeta SearchMeta, // todo: refactor and make it a Filter
			filter ...Filter,
		) (
			meta SearchMeta,
			products []domain.BasicProduct,
			availableFilter []Filter,
			err error,
		)
	}
)

var (
	_ Filter = NewKeyValueFilter("a", []string{"b", "c"})

	// ErrNotFound error
	ErrNotFound = errors.New("search not found")
)

// NewKeyValueFilter factory
func NewKeyValueFilter(k string, v []string) *KeyValueFilter {
	return &KeyValueFilter{
		k: k,
		v: v,
	}
}

// Values of the current filter
func (f *KeyValueFilter) Values() url.Values {
	return url.Values{
		f.k: f.v,
	}
}
