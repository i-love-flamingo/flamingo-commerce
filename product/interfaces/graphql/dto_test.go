package graphql_test

import (
	"fmt"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/product/application"
	graphql "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
)

// testCustomFacet is a custom facet DTO for testing
type testCustomFacet struct {
	name     string
	label    string
	position int
}

func (f *testCustomFacet) Name() string          { return f.name }
func (f *testCustomFacet) Label() string         { return f.label }
func (f *testCustomFacet) Position() int         { return f.position }
func (f *testCustomFacet) HasSelectedItem() bool { return false }

// testFacetMapper maps "CustomFacet" type to testCustomFacet
type testFacetMapper struct{}

func (m *testFacetMapper) MapFacet(facet searchdomain.Facet) searchdto.CommerceSearchFacet {
	if facet.Type == "CustomFacet" {
		return &testCustomFacet{
			name:     facet.Name,
			label:    facet.Label,
			position: facet.Position,
		}
	}

	return nil
}

func TestMapFacet_BuiltInTypes(t *testing.T) {
	t.Parallel()

	logger := flamingo.NullLogger{}

	tests := []struct {
		name     string
		facet    searchdomain.Facet
		wantType string
		wantNil  bool
	}{
		{
			name:     "ListFacet is mapped",
			facet:    searchdomain.Facet{Type: searchdomain.ListFacet, Name: "color"},
			wantType: "*searchdto.CommerceSearchListFacet",
		},
		{
			name:     "TreeFacet is mapped",
			facet:    searchdomain.Facet{Type: searchdomain.TreeFacet, Name: "category"},
			wantType: "*searchdto.CommerceSearchTreeFacet",
		},
		{
			name:     "RangeFacet is mapped",
			facet:    searchdomain.Facet{Type: searchdomain.RangeFacet, Name: "price"},
			wantType: "*searchdto.CommerceSearchRangeFacet",
		},
		{
			name:    "unknown type returns nil",
			facet:   searchdomain.Facet{Type: "UnknownFacet", Name: "unknown"},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := graphql.MapFacetForTest(tt.facet, nil, logger)

			if tt.wantNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}

				return
			}

			if result == nil {
				t.Fatal("expected non-nil result")
			}

			gotType := fmt.Sprintf("%T", result)
			if gotType != tt.wantType {
				t.Errorf("expected type %q, got %q", tt.wantType, gotType)
			}

			if result.Name() != tt.facet.Name {
				t.Errorf("expected name %q, got %q", tt.facet.Name, result.Name())
			}
		})
	}
}

func TestMapFacet_CustomMapperTakesPrecedence(t *testing.T) {
	t.Parallel()

	logger := flamingo.NullLogger{}
	mappers := []searchdto.FacetMapper{&testFacetMapper{}}

	facet := searchdomain.Facet{
		Type:     "CustomFacet",
		Name:     "custom",
		Label:    "Custom Label",
		Position: 5,
	}

	result := graphql.MapFacetForTest(facet, mappers, logger)
	if result == nil {
		t.Fatal("expected non-nil result for custom facet type")
	}

	if result.Name() != "custom" {
		t.Errorf("expected name %q, got %q", "custom", result.Name())
	}

	if result.Label() != "Custom Label" {
		t.Errorf("expected label %q, got %q", "Custom Label", result.Label())
	}

	if result.Position() != 5 {
		t.Errorf("expected position %d, got %d", 5, result.Position())
	}
}

func TestMapFacet_CustomMapperFallsBackToBuiltIn(t *testing.T) {
	t.Parallel()

	logger := flamingo.NullLogger{}
	mappers := []searchdto.FacetMapper{&testFacetMapper{}}

	// A ListFacet should still work when custom mapper returns nil for it
	facet := searchdomain.Facet{
		Type: searchdomain.ListFacet,
		Name: "color",
	}

	result := graphql.MapFacetForTest(facet, mappers, logger)
	if result == nil {
		t.Fatal("expected non-nil result for list facet with custom mapper present")
	}

	if result.Name() != "color" {
		t.Errorf("expected name %q, got %q", "color", result.Name())
	}
}

func TestMapFacet_CustomMapperCanOverrideBuiltIn(t *testing.T) {
	t.Parallel()

	logger := flamingo.NullLogger{}

	// A mapper that intercepts ListFacet type
	overrideMapper := &overrideListFacetMapper{}
	mappers := []searchdto.FacetMapper{overrideMapper}

	facet := searchdomain.Facet{
		Type:     searchdomain.ListFacet,
		Name:     "color",
		Label:    "Color",
		Position: 1,
	}

	result := graphql.MapFacetForTest(facet, mappers, logger)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// The override mapper should have handled it
	if _, ok := result.(*testCustomFacet); !ok {
		t.Error("expected custom mapper to override built-in ListFacet mapping")
	}
}

// overrideListFacetMapper intercepts ListFacet and returns a custom DTO
type overrideListFacetMapper struct{}

func (m *overrideListFacetMapper) MapFacet(facet searchdomain.Facet) searchdto.CommerceSearchFacet {
	if facet.Type == searchdomain.ListFacet {
		return &testCustomFacet{
			name:     facet.Name,
			label:    facet.Label,
			position: facet.Position,
		}
	}

	return nil
}

func TestSearchResultDTO_Facets_WithCustomMapper(t *testing.T) {
	t.Parallel()

	dto := graphql.NewSearchResultDTOForTest(
		&application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"custom": searchdomain.Facet{
					Type:     "CustomFacet",
					Name:     "custom",
					Label:    "Custom",
					Position: 2,
				},
				"color": searchdomain.Facet{
					Type:     searchdomain.ListFacet,
					Name:     "color",
					Label:    "Color",
					Position: 1,
				},
			},
		},
		flamingo.NullLogger{},
		[]searchdto.FacetMapper{&testFacetMapper{}},
	)

	facets := dto.Facets()

	if len(facets) != 2 {
		t.Fatalf("expected 2 facets, got %d", len(facets))
	}

	// Should be sorted by position
	if facets[0].Name() != "color" {
		t.Errorf("expected first facet to be 'color', got %q", facets[0].Name())
	}

	if facets[1].Name() != "custom" {
		t.Errorf("expected second facet to be 'custom', got %q", facets[1].Name())
	}
}

func TestSearchResultDTO_Facets_UnknownFacetDroppedWithoutMapper(t *testing.T) {
	t.Parallel()

	dto := graphql.NewSearchResultDTOForTest(
		&application.SearchResult{
			Facets: searchdomain.FacetCollection{
				"custom": searchdomain.Facet{
					Type:     "CustomFacet",
					Name:     "custom",
					Label:    "Custom",
					Position: 1,
				},
			},
		},
		flamingo.NullLogger{},
		nil,
	)

	facets := dto.Facets()

	if len(facets) != 0 {
		t.Fatalf("expected 0 facets without mapper, got %d", len(facets))
	}
}
