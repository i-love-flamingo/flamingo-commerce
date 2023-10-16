package fake_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchService_Search(t *testing.T) {
	s := fake.SearchService{}
	s.Inject(&fake.ProductService{}, &struct {
		LiveSearchJSON         string `inject:"config:commerce.product.fakeservice.jsonTestDataLiveSearch,optional"`
		CategoryFacetItemsJSON string `inject:"config:commerce.product.fakeservice.jsonTestDataCategoryFacetItems,optional"`
	}{})

	t.Run("Category Facet", func(t *testing.T) {

		t.Run("Selected category level 1", func(t *testing.T) {
			result, err := s.Search(context.Background(), searchDomain.NewKeyValueFilter("categoryCodes", []string{"clothing"}))
			require.Nil(t, err)
			assert.False(t, result.Facets["categoryCodes"].Items[0].Active, "Items[0].Active")
			assert.True(t, result.Facets["categoryCodes"].Items[1].Active, "Items[1].Active")
			assert.True(t, result.Facets["categoryCodes"].Items[1].Selected, "Items[1].Selected")
		})

		t.Run("Selected category level 2", func(t *testing.T) {
			result, err := s.Search(context.Background(), searchDomain.NewKeyValueFilter("categoryCodes", []string{"headphones"}))
			require.Nil(t, err)
			assert.True(t, result.Facets["categoryCodes"].Items[0].Active, "Items[0].Active")
			assert.True(t, result.Facets["categoryCodes"].Items[0].Items[1].Active, "Items[0].Items[1].Active")
			assert.True(t, result.Facets["categoryCodes"].Items[0].Items[1].Selected, "Items[0].Items[1].Selected")
			assert.False(t, result.Facets["categoryCodes"].Items[0].Items[1].Items[0].Active, "Items[0].Items[1].Items[0].Active")
		})

		t.Run("Selected category level 3", func(t *testing.T) {
			result, err := s.Search(context.Background(), searchDomain.NewKeyValueFilter("categoryCodes", []string{"headphone_accessories"}))
			require.Nil(t, err)
			assert.True(t, result.Facets["categoryCodes"].Items[0].Active, "Items[0].Active")
			assert.True(t, result.Facets["categoryCodes"].Items[0].Items[1].Active, "Items[0].Items[1].Active")
			assert.True(t, result.Facets["categoryCodes"].Items[0].Items[1].Items[0].Active, "Items[0].Items[1].Items[0].Active")
			assert.True(t, result.Facets["categoryCodes"].Items[0].Items[1].Items[0].Selected, "Items[0].Items[1].Items[0].Selected")
		})
	})
}

func TestSearchService_SearchBy(t *testing.T) {
	t.Parallel()

	type fields struct {
		liveSearchJSON string
	}

	type args struct {
		attribute string
		filters   []searchDomain.Filter
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.SearchResult
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "get livesearch results from json file",
			fields: fields{liveSearchJSON: filepath.Join("testdata", "livesearch.json")},
			args: args{
				attribute: "livesearch",
				filters: []searchDomain.Filter{
					searchDomain.NewQueryFilter("camera"),
				},
			},
			want: &domain.SearchResult{
				Result: searchDomain.Result{
					SearchMeta: searchDomain.SearchMeta{
						Query:          "",
						OriginalQuery:  "",
						Page:           1,
						NumPages:       1,
						NumResults:     0,
						SelectedFacets: []searchDomain.Facet(nil),
						SortOptions:    []searchDomain.SortOption(nil),
					},
					Hits: []searchDomain.Document{},
					Suggestion: []searchDomain.Suggestion{
						{
							Type:                 "type",
							Text:                 "text",
							Highlight:            "highlight",
							AdditionalAttributes: map[string]string{"additional": "value"},
						},
					},
					Facets: searchDomain.FacetCollection(nil),
					Promotions: []searchDomain.Promotion{
						{
							Title:   "Promotion title",
							Content: "Promotion content",
							URL:     "https://www.omnevo.net",
							Media: []searchDomain.Media{
								{
									Type:      "type",
									MimeType:  "mimetype",
									Usage:     "usage",
									Title:     "title",
									Reference: "reference"},
							},
							AdditionalAttributes: map[string]interface{}{"additional": "value"},
						},
					},
					Actions: []searchDomain.Action{
						{
							Type:                 "Redirect",
							Content:              "https://example.com",
							AdditionalAttributes: map[string]interface{}{"additional": "value"},
						},
					},
				},
				Hits: []domain.BasicProduct{},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := new(fake.SearchService).Inject(
				new(fake.ProductService),
				&struct {
					LiveSearchJSON         string `inject:"config:commerce.product.fakeservice.jsonTestDataLiveSearch,optional"`
					CategoryFacetItemsJSON string `inject:"config:commerce.product.fakeservice.jsonTestDataCategoryFacetItems,optional"`
				}{
					LiveSearchJSON: tt.fields.liveSearchJSON,
				},
			)
			got, err := s.SearchBy(context.Background(), tt.args.attribute, nil, tt.args.filters...)
			if !tt.wantErr(t, err, fmt.Sprintf("SearchBy(%v, %v, %v, %v)", context.Background(), tt.args.attribute, nil, tt.args.filters)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SearchBy(%v, %v, %v, %v)", context.Background(), tt.args.attribute, nil, tt.args.filters)
		})
	}
}
