package graphqlproductdto

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
	"github.com/stretchr/testify/assert"
)

var (
	fakeService = fake.ProductService{}
)

func TestMapSimpleProduct(t *testing.T) {
	simpleProduct, _ := fakeService.Get(nil, "fake_simple")
	variationSelection := NewVariantsToVariationSelections(simpleProduct)

	t.Run("should have empty variationSelection", func(t *testing.T) {
		assert.Equal(t, 0, len(variationSelection))
	})
}

func TestMapConfigurable(t *testing.T) {
	t.Run("should map variationSelection for one variation attribute", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(nil, "fake_configurable")
		configurable := fakeConfigurable.(domain.ConfigurableProduct)
		configurable.VariantVariationAttributes = []string{"color"}
		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:                  "Red",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-red-l",
					},
					{
						Label:                  "White",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Label:                  "Black",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-black-l",
					},
				},
			},
		}, variationSelection)
	})

	t.Run("should map variationSelection for two variation attributes", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(nil, "fake_configurable")
		configurable := fakeConfigurable.(domain.ConfigurableProduct)
		configurable.VariantVariationAttributes = []string{"color", "size"}

		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:                  "Red",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-red-l",
					},
					{
						Label:                  "White",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Label:                  "Black",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-black-l",
					},
				},
			}, {
				Code:  "size",
				Label: "Size",
				Options: []VariationSelectionOption{
					{
						Label:                  "M",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Label:                  "L",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-black-l",
					},
				},
			},
		}, variationSelection)
	})

	t.Run("should return empty variation selection if no product has matching attributes", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(nil, "fake_configurable")
		configurable := fakeConfigurable.(domain.ConfigurableProduct)
		configurable.VariantVariationAttributes = []string{"foobar"}
		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection(nil), variationSelection)
	})
}

func TestMapConfigurableWithActiveVariant(t *testing.T) {
	t.Run("should map variationSelection for one variation attribute with active variant", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(nil, "fake_configurable_with_active_variant")
		configurable := fakeConfigurable.(domain.ConfigurableProductWithActiveVariant)
		configurable.VariantVariationAttributes = []string{"color"}
		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:                  "Red",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-red-l",
					},
					{
						Label:                  "White",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Label:                  "Black",
						State:                  VariationSelectionOptionStateActive,
						VariantMarketPlaceCode: "shirt-black-l",
					},
				},
			},
		}, variationSelection)
	})

	t.Run("should map variationSelection for two variation attributes with active variant", func(t *testing.T) {
		fakeConfigurable, _ := fakeService.Get(nil, "fake_configurable_with_active_variant")
		configurable := fakeConfigurable.(domain.ConfigurableProductWithActiveVariant)
		configurable.VariantVariationAttributes = []string{"color", "size"}

		variationSelection := NewVariantsToVariationSelections(configurable)
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Label:                  "Red",
						State:                  VariationSelectionOptionStateMatch,
						VariantMarketPlaceCode: "shirt-red-l",
					},
					{
						Label:                  "White",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Label:                  "Black",
						State:                  VariationSelectionOptionStateActive,
						VariantMarketPlaceCode: "shirt-black-l",
					},
				},
			}, {
				Code:  "size",
				Label: "Size",
				Options: []VariationSelectionOption{
					{
						Label:                  "M",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Label:                  "L",
						State:                  VariationSelectionOptionStateActive,
						VariantMarketPlaceCode: "shirt-black-l",
					},
				},
			},
		}, variationSelection)
	})
}
