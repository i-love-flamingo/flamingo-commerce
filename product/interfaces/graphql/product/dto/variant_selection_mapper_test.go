package graphqlproductdto

import (
	"testing"

	"gotest.tools/v3/assert"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

func TestVariantSelectionMappingOfConfigurableProducts(t *testing.T) {
	t.Parallel()

	t.Run("just work", func(t *testing.T) {
		redS := domain.Variant{
			BasicProductData: domain.BasicProductData{
				MarketPlaceCode: "red-s",
				Attributes:      domain.Attributes{"color": {Label: "Red", CodeLabel: "Color", RawValue: "red"}, "size": {Label: "S", CodeLabel: "Clothing Size", RawValue: "s"}},
			},
		}
		redM := domain.Variant{
			BasicProductData: domain.BasicProductData{
				MarketPlaceCode: "red-m",
				Attributes:      domain.Attributes{"color": {Label: "Red", CodeLabel: "Colour", RawValue: "red"}, "size": {Label: "M", CodeLabel: "Clothing Size", RawValue: "m"}},
			},
		}
		redL := domain.Variant{
			BasicProductData: domain.BasicProductData{
				MarketPlaceCode: "red-l",
				Attributes:      domain.Attributes{"color": {Label: "Red", CodeLabel: "Colour", RawValue: "red"}, "size": {Label: "L", CodeLabel: "Clothing Size", RawValue: "l"}},
			},
		}
		blueS := domain.Variant{
			BasicProductData: domain.BasicProductData{
				MarketPlaceCode: "blue-s",
				Attributes:      domain.Attributes{"color": {Label: "Blue", CodeLabel: "Colour", RawValue: "blue"}, "size": {Label: "S", CodeLabel: "Clothing Size", RawValue: "s"}},
			},
		}
		blueM := domain.Variant{
			BasicProductData: domain.BasicProductData{
				MarketPlaceCode: "blue-m",
				Attributes:      domain.Attributes{"color": {Label: "Blue", CodeLabel: "Colour", RawValue: "blue"}, "size": {Label: "M", CodeLabel: "Clothing Size", RawValue: "m"}},
			},
		}

		configurable := domain.ConfigurableProduct{
			VariantVariationAttributes:        []string{"color", "size"},
			VariantVariationAttributesSorting: map[string][]string{"color": {"red", "blue"}, "size": {"s", "m", "l"}},
			Variants:                          []domain.Variant{redS, redM, redL, blueS, blueM},
		}

		got := MapVariantSelections(configurable)

		redSMarchingSelection := []VariantSelectionMatchAttributes{
			{
				Key:   "color",
				Value: "Red",
			},
			{
				Key:   "size",
				Value: "S",
			},
		}
		redMMarchingSelection := []VariantSelectionMatchAttributes{
			{
				Key:   "color",
				Value: "Red",
			},
			{
				Key:   "size",
				Value: "M",
			},
		}
		redLMatchingSelection := []VariantSelectionMatchAttributes{
			{
				Key:   "color",
				Value: "Red",
			},
			{
				Key:   "size",
				Value: "L",
			},
		}
		blueSMatchingSelection := []VariantSelectionMatchAttributes{
			{
				Key:   "color",
				Value: "Blue",
			},
			{
				Key:   "size",
				Value: "S",
			},
		}
		blueMMatchingSelection := []VariantSelectionMatchAttributes{
			{
				Key:   "color",
				Value: "Blue",
			},
			{
				Key:   "size",
				Value: "M",
			},
		}

		want := VariantSelection{
			Attributes: []VariantSelectionAttribute{
				{
					Label: "Color",
					Code:  "color",
					Options: []VariantSelectionAttributeOption{
						{
							Label:    "Red",
							UnitCode: "",
							OtherAttributesRestrictions: []OtherAttributesRestriction{
								{
									Code:             "size",
									AvailableOptions: []string{"S", "M", "L"},
								},
							},
						},
						{
							Label:    "Blue",
							UnitCode: "",
							OtherAttributesRestrictions: []OtherAttributesRestriction{
								{
									Code:             "size",
									AvailableOptions: []string{"S", "M"},
								},
							},
						},
					},
				},
				{
					Label: "Clothing Size",
					Code:  "size",
					Options: []VariantSelectionAttributeOption{
						{
							Label:    "S",
							UnitCode: "",
							OtherAttributesRestrictions: []OtherAttributesRestriction{
								{
									Code:             "color",
									AvailableOptions: []string{"Red", "Blue"},
								},
							},
						},
						{
							Label:    "M",
							UnitCode: "",
							OtherAttributesRestrictions: []OtherAttributesRestriction{
								{
									Code:             "color",
									AvailableOptions: []string{"Red", "Blue"},
								},
							},
						},
						{
							Label:    "L",
							UnitCode: "",
							OtherAttributesRestrictions: []OtherAttributesRestriction{
								{
									Code:             "color",
									AvailableOptions: []string{"Red"},
								},
							},
						},
					},
				},
			},
			Variants: []VariantSelectionMatch{
				{
					Variant:    VariantSelectionMatchVariant{MarketplaceCode: redS.MarketPlaceCode},
					Attributes: redSMarchingSelection,
				},
				{
					Variant:    VariantSelectionMatchVariant{MarketplaceCode: redM.MarketPlaceCode},
					Attributes: redMMarchingSelection,
				},
				{
					Variant:    VariantSelectionMatchVariant{MarketplaceCode: redL.MarketPlaceCode},
					Attributes: redLMatchingSelection,
				},
				{
					Variant:    VariantSelectionMatchVariant{MarketplaceCode: blueS.MarketPlaceCode},
					Attributes: blueSMatchingSelection,
				},
				{
					Variant:    VariantSelectionMatchVariant{MarketplaceCode: blueM.MarketPlaceCode},
					Attributes: blueMMatchingSelection,
				},
			},
		}

		assert.DeepEqual(t, want.Attributes, got.Attributes)

		assert.DeepEqual(t, got.Variants[0].Attributes, redSMarchingSelection)
		assert.Equal(t, got.Variants[0].Variant.MarketplaceCode, redS.MarketPlaceCode)

		assert.DeepEqual(t, got.Variants[1].Attributes, redMMarchingSelection)
		assert.Equal(t, got.Variants[1].Variant.MarketplaceCode, redM.MarketPlaceCode)

		assert.DeepEqual(t, got.Variants[2].Attributes, redLMatchingSelection)
		assert.Equal(t, got.Variants[2].Variant.MarketplaceCode, redL.MarketPlaceCode)

		assert.DeepEqual(t, got.Variants[3].Attributes, blueSMatchingSelection)
		assert.Equal(t, got.Variants[3].Variant.MarketplaceCode, blueS.MarketPlaceCode)

		assert.DeepEqual(t, got.Variants[4].Attributes, blueMMatchingSelection)
		assert.Equal(t, got.Variants[4].Variant.MarketplaceCode, blueM.MarketPlaceCode)
	})
}
