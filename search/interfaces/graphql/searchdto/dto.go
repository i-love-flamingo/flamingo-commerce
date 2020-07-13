package searchdto

// CommerceSearchRequest - search request structure for GraphQL
type CommerceSearchRequest struct {
	PageSize        int
	Page            int
	SortBy          string
	KeyValueFilters []CommerceSearchKeyValueFilter
	Query           string
}

// CommerceSearchKeyValueFilter - key value filter for CommerceSearchRequest
type CommerceSearchKeyValueFilter struct {
	K string
	V []string
}

// CommerceSearchSortOption – search option structure for GraphQL
type CommerceSearchSortOption struct {
	Label    string
	Field    string
	Selected bool
}
