package dto

import (
	"flamingo.me/flamingo-commerce/v3/search/application"
	"flamingo.me/flamingo-commerce/v3/search/domain"
)

// CommerceSearchRequest - search request structure for graphQl
type CommerceSearchRequest struct {
	PageSize        int
	Page            int
	SortBy          string
	SortDirection   string
	KeyValueFilters []CommerceSearchKeyValueFilter
	Query           string
}

// CommerceSearchKeyValueFilter - key value filter for CommerceSearchRequest
type CommerceSearchKeyValueFilter struct {
	K string
	V []string
}

// SearchRequestToFilters maps CommerceSearchRequest to Filter
func SearchRequestToFilters(searchRequest *CommerceSearchRequest, defaultPageSize int) []domain.Filter {
	var filters []domain.Filter

	if searchRequest != nil {
		filters = application.BuildFilters(application.SearchRequest{
			AdditionalFilter: nil,
			PageSize:         searchRequest.PageSize,
			Page:             searchRequest.Page,
			SortBy:           searchRequest.SortBy,
			SortDirection:    searchRequest.SortDirection,
			Query:            searchRequest.Query,
			PaginationConfig: nil,
		}, defaultPageSize)

		for _, filter := range searchRequest.KeyValueFilters {
			filters = append(filters, domain.NewKeyValueFilter(filter.K, filter.V))
		}
	}

	return filters
}
