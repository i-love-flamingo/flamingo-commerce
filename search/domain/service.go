package domain

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"sort"
)

const (
	// SuggestionTypeProduct represents product suggestions
	SuggestionTypeProduct = "product"
	// SuggestionTypeCategory represents category suggestions
	SuggestionTypeCategory = "category"
)

type (
	// SearchService defines how to access search
	SearchService interface {
		// Types() []string
		Search(ctx context.Context, filter ...Filter) (results map[string]Result, err error)
		SearchFor(ctx context.Context, typ string, filter ...Filter) (result *Result, err error)
	}

	// Result defines a search result for one type
	Result struct {
		SearchMeta SearchMeta
		Hits       []Document
		Suggestion []Suggestion
		Facets     FacetCollection
		Promotions []Promotion
		Actions    []Action
	}

	// Action might be considered on the frontend to be taken depending on search results
	Action struct {
		Type                 string
		Content              string
		AdditionalAttributes map[string]interface{}
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

	// SortOption defines how sorting is possible, and which of them are activated with both an asc and desc option
	SortOption struct {
		// Label that you normally want to show in the frontend (e.g. "Price")
		Label string
		// Field that you need to use in SearchRequest>SortFilter
		Field string
		// SelectedAsc true if sorting by this field is active
		SelectedAsc bool
		// SelectedDesc true if sorting by this field is active
		SelectedDesc bool
		// Asc - represents the field that is used to trigger ascending search.
		// Deprecated: use "Field" and "SelectedAsc" instead to set which field should be sortable
		Asc string
		// Desc - represents the field that is used to trigger descending search.
		// Deprecated: use "Field" and "SelectedDesc" instead to set which field should be sortable
		Desc string
	}

	// FacetType for type facets
	FacetType string

	// FacetItem contains information about a facet item
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

		Count int64
	}

	// Facet provided by the search backend
	Facet struct {
		Type     FacetType
		Name     string
		Label    string
		Items    []*FacetItem
		Position int
	}

	// FacetCollection for all available facets
	FacetCollection map[string]Facet
	facetSlice      []Facet

	// Suggestion hint
	Suggestion struct {
		Type                 string
		Text                 string
		Highlight            string
		AdditionalAttributes map[string]string
	}

	// Promotion result during search
	Promotion struct {
		Title                string
		Content              string
		URL                  string
		Media                []Media
		AdditionalAttributes map[string]interface{}
	}

	// Media contains promotion media data
	Media struct {
		Type      string
		MimeType  string
		Usage     string
		Title     string
		Reference string
	}

	// Document holds a search result document
	Document interface{}

	// RedirectError suggests to redirect
	RedirectError struct {
		To string
	}

	// RequestQueryHook can be used to enforce redirect errors
	RequestQueryHook interface {
		Hook(ctx context.Context, path string, query *url.Values) error
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

// Order a facet collection
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

// Facet types
const (
	ListFacet  FacetType = "ListFacet"
	TreeFacet  FacetType = "TreeFacet"
	RangeFacet FacetType = "RangeFacet"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New("search not found")
)

// ValidatePageSize checks if the pageSize is logical for current result
func (sm *SearchMeta) ValidatePageSize(pageSize int) error {
	if pageSize == 0 {
		return errors.New("cannot validate - no expected pageSize given")
	}
	expectedNumPages := math.Ceil(float64(sm.NumResults) / float64(pageSize))
	if expectedNumPages != float64(sm.NumPages) {
		return fmt.Errorf("pagesize not valid expected %f / given in result: %d", expectedNumPages, sm.NumPages)
	}
	return nil
}
