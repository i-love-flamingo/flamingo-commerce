package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/search/domain"
)

func TestBuildFilters(t *testing.T) {
	type args struct {
		request         SearchRequest
		defaultPageSize int
	}
	tests := []struct {
		name string
		args args
		want []domain.Filter
	}{
		{
			name: "default page size",
			args: args{
				request: SearchRequest{
					AdditionalFilter: []domain.Filter{domain.NewKeyValueFilter("key", []string{"value1", "value2"})},
					Page:             3,
					SortBy:           "price",
					SortDirection:    "desc",
					Query:            "query",
				},
				defaultPageSize: 15,
			},
			want: []domain.Filter{
				domain.NewQueryFilter("query"),
				domain.NewPaginationPageFilter(3),
				domain.NewPaginationPageSizeFilter(15),
				domain.NewSortFilter("price", "desc"),
				domain.NewKeyValueFilter("key", []string{"value1", "value2"}),
			},
		},
		{
			name: "given page size",
			args: args{
				request: SearchRequest{
					AdditionalFilter: []domain.Filter{domain.NewKeyValueFilter("key", []string{"value1", "value2"})},

					Page:          3,
					PageSize:      33,
					SortBy:        "price",
					SortDirection: "desc",
					Query:         "query",
				},
				defaultPageSize: 15,
			},
			want: []domain.Filter{
				domain.NewQueryFilter("query"),
				domain.NewPaginationPageFilter(3),
				domain.NewPaginationPageSizeFilter(33),
				domain.NewSortFilter("price", "desc"),
				domain.NewKeyValueFilter("key", []string{"value1", "value2"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildFilters(tt.args.request, tt.args.defaultPageSize)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSearchRequest_AddAdditionalFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request SearchRequest
		filter  domain.Filter
		want    []domain.Filter
	}{
		{
			name:    "add to empty filters",
			request: SearchRequest{},
			filter:  domain.NewKeyValueFilter("key1", []string{"value1"}),
			want: []domain.Filter{
				domain.NewKeyValueFilter("key1", []string{"value1"}),
			},
		},
		{
			name: "add to existing filters",
			request: SearchRequest{
				AdditionalFilter: []domain.Filter{
					domain.NewKeyValueFilter("existing", []string{"value"}),
				},
			},
			filter: domain.NewKeyValueFilter("key2", []string{"value2"}),
			want: []domain.Filter{
				domain.NewKeyValueFilter("existing", []string{"value"}),
				domain.NewKeyValueFilter("key2", []string{"value2"}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.request.AddAdditionalFilter(tt.filter)
			assert.Equal(t, tt.want, tt.request.AdditionalFilter)
		})
	}
}

func TestSearchRequest_SetAdditionalFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request SearchRequest
		filter  domain.Filter
		want    []domain.Filter
	}{
		{
			name:    "set on empty filters",
			request: SearchRequest{},
			filter:  domain.NewKeyValueFilter("key1", []string{"value1"}),
			want: []domain.Filter{
				domain.NewKeyValueFilter("key1", []string{"value1"}),
			},
		},
		{
			name: "replace existing filter",
			request: SearchRequest{
				AdditionalFilter: []domain.Filter{
					domain.NewKeyValueFilter("key1", []string{"old_value"}),
				},
			},
			filter: domain.NewKeyValueFilter("key1", []string{"new_value"}),
			want: []domain.Filter{
				domain.NewKeyValueFilter("key1", []string{"new_value"}),
			},
		},
		{
			name: "add new filter while keeping existing different ones",
			request: SearchRequest{
				AdditionalFilter: []domain.Filter{
					domain.NewKeyValueFilter("key1", []string{"value1"}),
				},
			},
			filter: domain.NewKeyValueFilter("key2", []string{"value2"}),
			want: []domain.Filter{
				domain.NewKeyValueFilter("key1", []string{"value1"}),
				domain.NewKeyValueFilter("key2", []string{"value2"}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.request.SetAdditionalFilter(tt.filter)
			assert.Equal(t, tt.want, tt.request.AdditionalFilter)
		})
	}
}
