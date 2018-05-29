package domain

import (
	"context"
	"errors"
	"sort"
)

type (
	// Filter interface for search queries
	Filter interface {
		Value() (string, []string)
	}

	// KeyValueFilter allows simple k -> []values filtering
	KeyValueFilter struct {
		k string
		v []string
	}

	// SearchMeta data
	SearchMeta struct {
		Query          string
		OriginalQuery  string
		Page           int
		NumPages       int
		NumResults     int
		SelectedFacets []Facet
		SortOptions    []SortOption
	}

	// SortOption defines how sorting is possible, with both an asc and desc option
	SortOption struct {
		Label        string
		Asc          string
		Desc         string
		SelectedAsc  bool
		SelectedDesc bool
	}

	FacetType string

	FacetItem struct {
		Label    string
		Value    string
		Active   bool
		Selected bool

		// Tree Facet
		Items []*FacetItem

		// Range Facet
		Min, Max                 float64
		SelectedMin, SelectedMax float64
	}

	Facet struct {
		Type     string
		Name     string
		Label    string
		Items    []*FacetItem
		Position int
	}

	FacetCollection map[string]Facet
	facetSlice      []Facet

	// Result defines a search result for one type
	Result struct {
		SearchMeta SearchMeta
		Hits       []Document
		Suggestion []Suggestion
		Facets     FacetCollection
	}

	Suggestion struct {
		Text      string
		Highlight string
	}

	// Document holds a search result document
	Document interface{}

	// SearchService defines how to access search
	SearchService interface {
		//Types() []string
		Search(ctx context.Context, filter ...Filter) (results map[string]Result, err error)
		SearchFor(ctx context.Context, typ string, filter ...Filter) (result Result, err error)
	}

	RedirectError struct {
		To string
	}
)

func (re *RedirectError) Error() string {
	return "Error: enforced redirect to " + re.To
}

func (fs facetSlice) Len() int {
	return len(fs)
}

func (fs facetSlice) Less(i, j int) bool {
	return fs[i].Position < fs[j].Position
}

func (fs facetSlice) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fc FacetCollection) Order() []string {
	order := make(facetSlice, len(fc))
	i := 0
	for _, k := range fc {
		order[i] = k
		i++
	}

	sort.Stable(order)

	strings := make([]string, len(order))
	for i, v := range order {
		strings[i] = v.Name
	}

	return strings
}

const (
	ListFacet  FacetType = "ListFacet"
	TreeFacet            = "TreeFacet"
	RangeFacet           = "RangeFacet"
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

// Value of the current filter
func (f *KeyValueFilter) Value() (string, []string) {
	return f.k, f.v
}
