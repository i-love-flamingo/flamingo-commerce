package dto

// CommerceSearchRequest - search request structure for graphQl
type CommerceSearchRequest struct {
	PageSize        int
	Page            int
	SortBy          string
	SortDirection   string
	KeyValueFilters []CommerceSearchKeyValueFilter
}

// CommerceSearchKeyValueFilter - key value filter for CommerceSearchRequest
type CommerceSearchKeyValueFilter struct {
	K string
	V []string
}
