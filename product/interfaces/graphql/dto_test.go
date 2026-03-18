package graphql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/product/application"
	productgraphql "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
)

// customFacet is a test custom facet type.
type customFacet struct {
	name     string
	label    string
	position int
}

func (f *customFacet) Name() string          { return f.name }
func (f *customFacet) Label() string         { return f.label }
func (f *customFacet) Position() int         { return f.position }
func (f *customFacet) HasSelectedItem() bool { return false }

// testFacetMapper is a test mapper that handles "CustomFacet" type.
type testFacetMapper struct{}

func (m *testFacetMapper) MapFacet(facet searchdomain.Facet, next searchdto.FacetMapperFunc) (searchdto.CommerceSearchFacet, bool) {
	if facet.Type == "CustomFacet" {
		return &customFacet{
			name:     facet.Name,
			label:    facet.Label,
			position: facet.Position,
		}, true
	}

	return next(facet)
}

// overrideListFacetMapper overrides the built-in ListFacet handling.
type overrideListFacetMapper struct{}

func (m *overrideListFacetMapper) MapFacet(facet searchdomain.Facet, next searchdto.FacetMapperFunc) (searchdto.CommerceSearchFacet, bool) {
	if facet.Type == searchdomain.ListFacet {
		return &customFacet{
			name:     facet.Name,
			label:    facet.Label,
			position: facet.Position,
		}, true
	}

	return next(facet)
}

// builtInMappers returns the default facet mappers for built-in types.
func builtInMappers() []searchdto.FacetMapper {
	return []searchdto.FacetMapper{
		&searchdto.ListFacetMapper{},
		&searchdto.TreeFacetMapper{},
		&searchdto.RangeFacetMapper{},
	}
}

func newSearchResultDTO(result *application.SearchResult, mappers []searchdto.FacetMapper) *productgraphql.SearchResultDTO {
	factory := &productgraphql.SearchResultDTOFactory{}
	factory.Inject(flamingo.NullLogger{}, mappers)

	return factory.NewSearchResultDTO(result)
}

func TestSearchResultDTO_Facets(t *testing.T) {
	t.Parallel()

	t.Run("built-in facet types", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			facetType searchdomain.FacetType
			wantCount int
		}{
			{
				name:      "ListFacet maps correctly",
				facetType: searchdomain.ListFacet,
				wantCount: 1,
			},
			{
				name:      "TreeFacet maps correctly",
				facetType: searchdomain.TreeFacet,
				wantCount: 1,
			},
			{
				name:      "RangeFacet maps correctly",
				facetType: searchdomain.RangeFacet,
				wantCount: 1,
			},
			{
				name:      "Unknown facet is skipped without custom mappers",
				facetType: "UnknownFacet",
				wantCount: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := &application.SearchResult{
					Facets: searchdomain.FacetCollection{
						"test": searchdomain.Facet{
							Type:  tt.facetType,
							Name:  "test",
							Label: "Test",
						},
					},
				}

				dto := newSearchResultDTO(result, builtInMappers())
				facets := dto.Facets()

				assert.Len(t, facets, tt.wantCount)
			})
		}
	})

	t.Run("mixed facets sorted by position", func(t *testing.T) {
		t.Parallel()

		mappers := append(builtInMappers(), &testFacetMapper{})

		result := &application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"list": searchdomain.Facet{
					Type:     searchdomain.ListFacet,
					Name:     "list",
					Label:    "List Facet",
					Position: 2,
				},
				"custom": searchdomain.Facet{
					Type:     "CustomFacet",
					Name:     "custom",
					Label:    "Custom Facet",
					Position: 1,
				},
			},
		}

		dto := newSearchResultDTO(result, mappers)
		facets := dto.Facets()

		require.Len(t, facets, 2)
		assert.Equal(t, "custom", facets[0].Name(), "expected first facet to be 'custom' (position 1)")
		assert.Equal(t, "list", facets[1].Name(), "expected second facet to be 'list' (position 2)")
	})

	t.Run("no mappers returns no facets", func(t *testing.T) {
		t.Parallel()

		result := &application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"list": searchdomain.Facet{
					Type:     searchdomain.ListFacet,
					Name:     "list",
					Label:    "List Facet",
					Position: 1,
				},
			},
		}

		dto := newSearchResultDTO(result, nil)
		facets := dto.Facets()

		assert.Empty(t, facets, "expected no facets without any mappers")
	})
}

func TestSearchResultDTO_FacetMapper(t *testing.T) {
	t.Parallel()

	mappers := append(builtInMappers(), &testFacetMapper{})

	t.Run("custom facet type is handled by mapper", func(t *testing.T) {
		t.Parallel()

		result := &application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"custom": searchdomain.Facet{
					Type:     "CustomFacet",
					Name:     "custom",
					Label:    "Custom",
					Position: 5,
				},
			},
		}

		dto := newSearchResultDTO(result, mappers)
		facets := dto.Facets()

		require.Len(t, facets, 1)
		assert.Equal(t, "custom", facets[0].Name())
		assert.Equal(t, "Custom", facets[0].Label())
		assert.Equal(t, 5, facets[0].Position())
	})

	t.Run("built-in type still works with custom mapper registered", func(t *testing.T) {
		t.Parallel()

		result := &application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"list": searchdomain.Facet{
					Type:  searchdomain.ListFacet,
					Name:  "list",
					Label: "List",
				},
			},
		}

		dto := newSearchResultDTO(result, mappers)
		facets := dto.Facets()

		require.Len(t, facets, 1)
		assert.Equal(t, "list", facets[0].Name())
	})

	t.Run("unknown type is skipped when no mapper handles it", func(t *testing.T) {
		t.Parallel()

		result := &application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"unknown": searchdomain.Facet{
					Type: "AnotherUnknownType",
					Name: "unknown",
				},
			},
		}

		dto := newSearchResultDTO(result, mappers)
		facets := dto.Facets()

		assert.Empty(t, facets, "expected no facets for unhandled type")
	})
}

func TestSearchResultDTO_FacetMapperPriority(t *testing.T) {
	t.Parallel()

	// Override mapper is registered after built-in mappers (by a downstream module),
	// so it becomes the outermost middleware and takes precedence.
	mappers := append(builtInMappers(), &overrideListFacetMapper{})
	result := &application.SearchResult{
		Facets: searchdomain.FacetCollection{
			"overridden": searchdomain.Facet{
				Type:     searchdomain.ListFacet,
				Name:     "overridden",
				Label:    "Overridden",
				Position: 3,
			},
		},
	}

	dto := newSearchResultDTO(result, mappers)
	facets := dto.Facets()

	require.Len(t, facets, 1)

	// The custom mapper should take precedence over the built-in mapper
	cf, ok := facets[0].(*customFacet)
	require.True(t, ok, "expected result to be *customFacet (from override mapper)")
	assert.Equal(t, "overridden", cf.Name())
}
