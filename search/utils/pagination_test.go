package utils

import (
	"log"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildWith(t *testing.T) {
	type args struct {
		currentResult    CurrentResultInfos
		paginationConfig PaginationConfig
		urlBase          *url.URL
	}
	tests := []struct {
		name string
		args args
		want PaginationInfo
	}{
		{
			name: "first page",
			args: args{
				currentResult: CurrentResultInfos{
					TotalHits:  10,
					PageSize:   10,
					ActivePage: 1,
					LastPage:   1,
				},
				paginationConfig: PaginationConfig{
					ShowAroundActivePageAmount: 1,
					ShowFirstPage:              true,
					ShowLastPage:               true,
				},
				urlBase: &url.URL{},
			},
			want: PaginationInfo{
				TotalHits: 10,
				//No next oage
				PageNavigation: []Page{
					{
						Page:     1,
						URL:      makeURL(&url.URL{}, 1, ""),
						IsActive: true,
					},
				},
			},
		},
		{
			name: "page 7",
			args: args{
				currentResult: CurrentResultInfos{
					TotalHits:  100,
					PageSize:   10,
					ActivePage: 7,
					LastPage:   10,
				},
				paginationConfig: PaginationConfig{
					ShowAroundActivePageAmount: 3,
					ShowFirstPage:              false,
					ShowLastPage:               false,
				},
				urlBase: &url.URL{},
			},
			want: PaginationInfo{
				TotalHits: 100,
				NextPage: &Page{
					Page: 8,
					URL:  makeURL(&url.URL{}, 8, ""),
				},
				PreviousPage: &Page{
					Page: 6,
					URL:  makeURL(&url.URL{}, 6, ""),
				},
				//No next oage
				PageNavigation: []Page{
					{
						IsSpacer: true,
					},
					{
						Page: 4,
						URL:  makeURL(&url.URL{}, 4, ""),
					},
					{
						Page: 5,
						URL:  makeURL(&url.URL{}, 5, ""),
					},
					{
						Page: 6,
						URL:  makeURL(&url.URL{}, 6, ""),
					},
					{
						Page:     7,
						URL:      makeURL(&url.URL{}, 7, ""),
						IsActive: true,
					},
					{
						Page: 8,
						URL:  makeURL(&url.URL{}, 8, ""),
					},
					{
						Page: 9,
						URL:  makeURL(&url.URL{}, 9, ""),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildWith(tt.args.currentResult, tt.args.paginationConfig, tt.args.urlBase)
			log.Printf("%#v", got)
			assert.Equal(t, tt.want, got)
		})
	}
}
