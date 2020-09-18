package graphqlproductdto_test

import (
	"flamingo.me/flamingo-commerce/v3/product/domain"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
	"gotest.tools/assert"
	"testing"
)

func TestProductMedia_GetMedia(t *testing.T) {
	product := graphqlProductDto.NewGraphqlProductDto(domain.ConfigurableProduct{
		Identifier:       "configurable",
		BasicProductData: domain.BasicProductData{},
		Teaser: domain.TeaserData{
			Media: []domain.Media{
				{
					Type:      "teaser",
					MimeType:  "teaser",
					Usage:     "teaser",
					Title:     "teaser",
					Reference: "teaser",
				}, {
					Type:      "detail",
					MimeType:  "detail",
					Usage:     "detail",
					Title:     "detail",
					Reference: "detail",
				},
				{
					Type:      "teaser2",
					MimeType:  "teaser2",
					Usage:     "teaser",
					Title:     "teaser2",
					Reference: "teaser2",
				},
				{
					Type:      "detail2",
					MimeType:  "detail2",
					Usage:     "detail",
					Title:     "detail2",
					Reference: "detail2",
				},
			},
		},
		VariantVariationAttributes:        nil,
		Variants:                          nil,
		VariantVariationAttributesSorting: nil,
	}, nil)

	assert.DeepEqual(t, &domain.Media{
		Type:      "teaser",
		MimeType:  "teaser",
		Usage:     "teaser",
		Title:     "teaser",
		Reference: "teaser",
	}, product.Media().GetMedia("teaser"))

	assert.DeepEqual(t, &domain.Media{
		Type:      "detail",
		MimeType:  "detail",
		Usage:     "detail",
		Title:     "detail",
		Reference: "detail",
	}, product.Media().GetMedia("detail"))

}

func TestNewGraphqlProductDto(t *testing.T) {
	simpleProduct := domain.SimpleProduct{
		Identifier:       "simple",
		BasicProductData: domain.BasicProductData{},
		Saleable:         domain.Saleable{},
		Teaser:           domain.TeaserData{},
	}

	graphqlSimpleProduct := graphqlProductDto.NewGraphqlProductDto(simpleProduct, nil)
	assert.Equal(t, simpleProduct.Type(), graphqlSimpleProduct.Type())

	configurableProduct := domain.SimpleProduct{
		Identifier:       "configurable",
		BasicProductData: domain.BasicProductData{},
		Saleable:         domain.Saleable{},
		Teaser:           domain.TeaserData{},
	}

	graphqlConfigurableProduct := graphqlProductDto.NewGraphqlProductDto(configurableProduct, nil)
	assert.Equal(t, configurableProduct.Type(), graphqlConfigurableProduct.Type())

	activeVariantProduct := domain.SimpleProduct{
		Identifier:       "activeVariant",
		BasicProductData: domain.BasicProductData{},
		Saleable:         domain.Saleable{},
		Teaser:           domain.TeaserData{},
	}

	graphqlActiveVariantProduct := graphqlProductDto.NewGraphqlProductDto(activeVariantProduct, nil)
	assert.Equal(t, activeVariantProduct.Type(), graphqlActiveVariantProduct.Type())
}
