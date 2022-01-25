package fake_test

import (
	"context"
	"testing"

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
