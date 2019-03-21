package application

import (
	"reflect"
	"testing"

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
			for i, want := range tt.want {
				if len(got) <= i {
					t.Fatalf("too few entries in filter: want %d, got %d", len(tt.want), len(got))
				}
				if !reflect.DeepEqual(got[i], want) {
					t.Errorf("BuildFilters() = %#v, want %#v", got[i], want)
				}
			}
		})
	}
}
