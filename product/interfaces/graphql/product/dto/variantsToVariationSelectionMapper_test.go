package graphqlProductDto

import (
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	fakeService = fake.ProductService{}
)

func TestMapSimpleProduct(t *testing.T) {
	simpleProduct, _ := fakeService.Get(nil, "fake_simple")
	mapper := New(simpleProduct)
	variationSelection := mapper.Map()

	t.Run("should have empty variationSelection", func(t *testing.T) {
		assert.Equal(t, 0, len(variationSelection))
	})
}

func TestMapConfigurable(t *testing.T) {
	fakeConfigurable, _ := fakeService.Get(nil, "fake_configurable")
	configurable := fakeConfigurable.(domain.ConfigurableProduct)


	t.Run("should map variationSelection for one variation attribute", func(t *testing.T) {
		configurable.VariantVariationAttributes = []string{"color"}
		mapper := New(configurable)
		variationSelection := mapper.Map()
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Code:                   "white",
						Label:                  "White",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Code:                   "black",
						Label:                  "Black",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-black-l",
					},
					{
						Code:                   "red",
						Label:                  "Red",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-red-l",
					},
				},
			},
		}, variationSelection)
	})

	t.Run("should map variationSelection for two variation attributes", func(t *testing.T) {
		configurable.VariantVariationAttributes = []string{"color", "size"}

		mapper := New(configurable)
		variationSelection := mapper.Map()
		assert.Equal(t, []VariationSelection{
			{
				Code:  "color",
				Label: "Color",
				Options: []VariationSelectionOption{
					{
						Code:                   "white",
						Label:                  "White",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Code:                   "black",
						Label:                  "Black",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-black-l",
					},
					{
						Code:                   "red",
						Label:                  "Red",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-red-l",
					},
				},
			}, {
				Code:  "size",
				Label: "Size",
				Options: []VariationSelectionOption{
					{
						Code:                   "m",
						Label:                  "M",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-white-m",
					},
					{
						Code:                   "l",
						Label:                  "L",
						State:                  VariationSelectionOptionStateNoMatch,
						VariantMarketPlaceCode: "shirt-black-l",
					},
				},
			},
		}, variationSelection)
	})
}

//func TestMapConfigurableWithActiveVariant(t *testing.T) {
//	mapper := New(configurableProduct)
//	variationSelection := mapper.Map()
//
//	t.Run("should have empty variationSelection", func(t *testing.T) {
//		assert.Equal(t, 0, len(variationSelection))
//	})
//}
