package graphql_test

import (
	"testing"

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

func (m *testFacetMapper) MapFacet(facet searchdomain.Facet) (searchdto.CommerceSearchFacet, bool) {
	if facet.Type == "CustomFacet" {
		return &customFacet{
			name:     facet.Name,
			label:    facet.Label,
			position: facet.Position,
		}, true
	}

	return nil, false
}

// overrideListFacetMapper overrides the built-in ListFacet handling.
type overrideListFacetMapper struct{}

func (m *overrideListFacetMapper) MapFacet(facet searchdomain.Facet) (searchdto.CommerceSearchFacet, bool) {
	if facet.Type == searchdomain.ListFacet {
		return &customFacet{
			name:     facet.Name,
			label:    facet.Label,
			position: facet.Position,
		}, true
	}

	return nil, false
}

func newSearchResultDTO(result *application.SearchResult, mappers []searchdto.FacetMapper) *productgraphql.SearchResultDTO {
	dto := productgraphql.WrapSearchResult(result)
	dto.Inject(flamingo.NullLogger{}, &struct {
		FacetMappers []searchdto.FacetMapper `inject:",optional"`
	}{FacetMappers: mappers})

	return dto
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
				name:      "Unknown facet is skipped without mappers",
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

				dto := newSearchResultDTO(result, nil)
				facets := dto.Facets()

				if len(facets) != tt.wantCount {
					t.Errorf("expected %d facet(s), got %d", tt.wantCount, len(facets))
				}
			})
		}
	})

	t.Run("custom mapper", func(t *testing.T) {
		t.Parallel()

		mappers := []searchdto.FacetMapper{&testFacetMapper{}}

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

			if len(facets) != 1 {
				t.Fatalf("expected 1 facet, got %d", len(facets))
			}

			if facets[0].Name() != "custom" {
				t.Errorf("expected name 'custom', got %q", facets[0].Name())
			}

			if facets[0].Label() != "Custom" {
				t.Errorf("expected label 'Custom', got %q", facets[0].Label())
			}

			if facets[0].Position() != 5 {
				t.Errorf("expected position 5, got %d", facets[0].Position())
			}
		})

		t.Run("built-in type still works with mapper registered", func(t *testing.T) {
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

			if len(facets) != 1 {
				t.Fatalf("expected 1 facet, got %d", len(facets))
			}

			if facets[0].Name() != "list" {
				t.Errorf("expected name 'list', got %q", facets[0].Name())
			}
		})

		t.Run("unknown type is skipped when mapper does not handle it", func(t *testing.T) {
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

			if len(facets) != 0 {
				t.Errorf("expected 0 facets for unhandled type, got %d", len(facets))
			}
		})
	})

	t.Run("mapper priority", func(t *testing.T) {
		t.Parallel()

		mappers := []searchdto.FacetMapper{&overrideListFacetMapper{}}

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

		if len(facets) != 1 {
			t.Fatalf("expected 1 facet, got %d", len(facets))
		}

		// The custom mapper should take precedence over the built-in handler
		cf, ok := facets[0].(*customFacet)
		if !ok {
			t.Fatal("expected result to be *customFacet (from override mapper)")
		}

		if cf.Name() != "overridden" {
			t.Errorf("expected name 'overridden', got %q", cf.Name())
		}
	})

	t.Run("mixed facets sorted by position", func(t *testing.T) {
		t.Parallel()

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

		dto := newSearchResultDTO(result, []searchdto.FacetMapper{&testFacetMapper{}})
		facets := dto.Facets()

		if len(facets) != 2 {
			t.Fatalf("expected 2 facets, got %d", len(facets))
		}

		// Facets should be sorted by position
		if facets[0].Name() != "custom" {
			t.Errorf("expected first facet to be 'custom' (position 1), got %q", facets[0].Name())
		}

		if facets[1].Name() != "list" {
			t.Errorf("expected second facet to be 'list' (position 2), got %q", facets[1].Name())
		}
	})

	t.Run("without mappers only built-in types returned", func(t *testing.T) {
		t.Parallel()

		result := &application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"list": searchdomain.Facet{
					Type:     searchdomain.ListFacet,
					Name:     "list",
					Label:    "List Facet",
					Position: 1,
				},
				"custom": searchdomain.Facet{
					Type:     "CustomFacet",
					Name:     "custom",
					Label:    "Custom Facet",
					Position: 2,
				},
			},
		}

		dto := newSearchResultDTO(result, nil)
		facets := dto.Facets()

		// Without mappers, only built-in types are returned
		if len(facets) != 1 {
			t.Fatalf("expected 1 facet (only built-in), got %d", len(facets))
		}

		if facets[0].Name() != "list" {
			t.Errorf("expected facet 'list', got %q", facets[0].Name())
		}
	})
}
