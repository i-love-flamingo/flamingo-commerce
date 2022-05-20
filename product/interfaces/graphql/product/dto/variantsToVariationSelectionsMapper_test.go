package graphqlproductdto

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
	"github.com/stretchr/testify/assert"
)

var (
	fakeService = fake.ProductService{}
)

func getOptionVariantByMarketPlaceCode(variants []domain.Variant, marketPlaceCode string) VariationSelectionOptionVariant {
	for _, variant := range variants {
		if variant.MarketPlaceCode == marketPlaceCode {
			return NewVariationSelectionOptionVariant(variant)
		}
	}

	return NewVariationSelectionOptionVariant(domain.Variant{
		BasicProductData: domain.BasicProductData{},
		Saleable:         domain.Saleable{},
	})
}

func TestMapSimpleProduct(t *testing.T) {
	simpleProduct, _ := fakeService.Get(context.Background(), "fake_simple")
	variationSelection := NewVariantsToVariationSelections(simpleProduct)

	t.Run("should have empty variationSelection", func(t *testing.T) {
		assert.Equal(t, 0, len(variationSelection))
	})
}

func TestMapConfigurable(t *testing.T) {
	t.Run("should map variationSelection for one variation attribute", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(context.Background(), "fake_configurable")
		configurable := fakeConfigurable.(domain.ConfigurableProduct)
		configurable.VariantVariationAttributes = []string{"color"}
		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:   "Red",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-red-l"),
					},
					{
						Label:   "White",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-white-m"),
					},
					{
						Label:   "Black",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-black-l"),
					},
				},
			},
		}, variationSelection)
	})

	t.Run("should map variationSelection for two variation attributes", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(context.Background(), "fake_configurable")
		configurable := fakeConfigurable.(domain.ConfigurableProduct)

		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:   "Red",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-red-m"),
					},
					{
						Label:   "White",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-white-m"),
					},
					{
						Label:   "Black",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-black-l"),
					},
				},
			}, {
				Code:  "size",
				Label: "Size",
				Options: []VariationSelectionOption{
					{
						Label:   "M",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-red-m"),
					},
					{
						Label:   "L",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-red-l"),
					},
				},
			},
		}, variationSelection)
	})

	t.Run("should return empty variation selection if no product has matching attributes", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(context.Background(), "fake_configurable")
		configurable := fakeConfigurable.(domain.ConfigurableProduct)
		configurable.VariantVariationAttributes = []string{"foobar"}
		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection(nil), variationSelection)
	})
}

func TestMapConfigurableWithActiveVariant(t *testing.T) {
	t.Run("should map variationSelection for one variation attribute with active variant", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(context.Background(), "fake_configurable_with_active_variant")
		configurable := fakeConfigurable.(domain.ConfigurableProductWithActiveVariant)
		configurable.VariantVariationAttributes = []string{"color"}
		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:   "Red",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-red-l"),
					},
					{
						Label:   "White",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-white-m"),
					},
					{
						Label:   "Black",
						State:   VariationSelectionOptionStateActive,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-black-l"),
					},
				},
			},
		}, variationSelection)
	})

	t.Run("should map variationSelection for two variation attributes with active variant", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(context.Background(), "fake_configurable_with_active_variant")
		configurable := fakeConfigurable.(domain.ConfigurableProductWithActiveVariant)
		configurable.VariantVariationAttributes = []string{"color", "size"}
		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:   "Red",
						State:   VariationSelectionOptionStateMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-red-l"),
					},
					{
						Label:   "White",
						State:   VariationSelectionOptionStateNoMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-white-m"),
					},
					{
						Label:   "Black",
						State:   VariationSelectionOptionStateActive,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-black-l"),
					},
				},
			}, {
				Code:  "size",
				Label: "Size",
				Options: []VariationSelectionOption{
					{
						Label:   "M",
						State:   VariationSelectionOptionStateNoMatch,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-red-m"),
					},
					{
						Label:   "L",
						State:   VariationSelectionOptionStateActive,
						Variant: getOptionVariantByMarketPlaceCode(configurable.Variants, "shirt-black-l"),
					},
				},
			},
		}, variationSelection)
	})
}
