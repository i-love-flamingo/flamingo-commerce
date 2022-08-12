package domain

import (
	"strconv"
)

type (
	// Filter interface for search queries
	Filter interface {
		// Value very generic method for filters - returning typical Parameter Name and its setted values
		Value() (string, []string)
	}

	// KeyValueFilter allows simple k -> []values filtering
	KeyValueFilter struct {
		k string
		v []string
	}

	// SortFilter - specifies the request to sort by some criteria(label) in a certain direction. Possible values for label and direction should be in SearchMeta.SortOption
	SortFilter struct {
		field     string
		direction string
	}

	// QueryFilter - represents a query string, normally given by a user in the search result
	QueryFilter struct {
		query string
	}

	// PaginationPage - if search supports pagination this filter tells which page to return
	PaginationPage struct {
		page int
	}

	// PaginationPageSize  - if search supports setting the amount (limit) per page
	PaginationPageSize struct {
		pageSize int
	}
)

var (
	_ Filter = NewKeyValueFilter("a", []string{"b", "c"})
)

const (
	// SortDirectionAscending general asc value
	SortDirectionAscending = "A"

	// SortDirectionDescending general desc value
	SortDirectionDescending = "D"

	// SortDirectionNone general not set value
	SortDirectionNone = ""
)

// NewKeyValueFilters - Factory method that you can use to get a list of KeyValueFilter based from url.Values
func NewKeyValueFilters(params map[string][]string) []Filter {
	var result []Filter
	for k, v := range params {
		if len(v) == 0 {
			continue
		}
		result = append(result, NewKeyValueFilter(k, v))
	}
	return result
}

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

// KeyValues of the current filter
func (f *KeyValueFilter) KeyValues() []string {
	return f.v
}

// Key of the current filter
func (f *KeyValueFilter) Key() string {
	return f.k
}

// NewSortFilter factory
func NewSortFilter(label string, direction string) *SortFilter {
	if direction != SortDirectionNone && direction != SortDirectionDescending && direction != SortDirectionAscending {
		direction = SortDirectionNone
	}
	return &SortFilter{
		field:     label,
		direction: direction,
	}
}

// Value of the current filter
func (f *SortFilter) Value() (string, []string) {
	return f.field, []string{f.direction}
}

// Field of the current filter
func (f *SortFilter) Field() string {
	return f.field
}

// Direction of the current filter
func (f *SortFilter) Direction() string {
	return f.direction
}

// Descending returns true if sort order is descending
func (f *SortFilter) Descending() bool {
	return f.direction == SortDirectionDescending
}

// NewQueryFilter factory
func NewQueryFilter(query string) *QueryFilter {
	return &QueryFilter{
		query: query,
	}
}

// Value of the current filter
func (f *QueryFilter) Value() (string, []string) {
	return "q", []string{f.query}
}

// Query of the current filter
func (f *QueryFilter) Query() string {
	return f.query
}

// NewPaginationPageFilter factory
func NewPaginationPageFilter(page int) *PaginationPage {
	return &PaginationPage{
		page: page,
	}
}

// Value of the current filter
func (f *PaginationPage) Value() (string, []string) {
	return "page", []string{strconv.Itoa(f.page)}
}

// GetPage of the current filter
func (f *PaginationPage) GetPage() int {
	return f.page
}

// NewPaginationPageSizeFilter factory
func NewPaginationPageSizeFilter(page int) *PaginationPageSize {
	return &PaginationPageSize{
		pageSize: page,
	}
}

// Value of the current filter
func (f *PaginationPageSize) Value() (string, []string) {
	return "limit", []string{strconv.Itoa(f.pageSize)}
}

// GetPageSize of the current filter
func (f *PaginationPageSize) GetPageSize() int {
	return f.pageSize
}
