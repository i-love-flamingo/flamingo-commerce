package graphqlproductdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
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
	}, nil, nil)

	assert.Equal(t, &domain.Media{
		Type:      "teaser",
		MimeType:  "teaser",
		Usage:     "teaser",
		Title:     "teaser",
		Reference: "teaser",
	}, product.Media().GetMedia("teaser"))

	assert.Equal(t, &domain.Media{
		Type:      "detail",
		MimeType:  "detail",
		Usage:     "detail",
		Title:     "detail",
		Reference: "detail",
	}, product.Media().GetMedia("detail"))

}

func TestNewGraphqlProductDto(t *testing.T) {
	// simple input
	simpleProduct := domain.SimpleProduct{
		Identifier:       "simple",
		BasicProductData: domain.BasicProductData{},
		Saleable:         domain.Saleable{},
		Teaser:           domain.TeaserData{},
	}
	graphqlSimpleProduct := graphqlProductDto.NewGraphqlProductDto(simpleProduct, nil, nil)
	assert.Equal(t, "simple", graphqlSimpleProduct.Type())

	// Configurable input
	configurableProduct := domain.ConfigurableProduct{
		Identifier:       "configurable",
		BasicProductData: domain.BasicProductData{},
		Teaser:           domain.TeaserData{},
	}
	graphqlConfigurableProduct := graphqlProductDto.NewGraphqlProductDto(configurableProduct, nil, nil)
	assert.Equal(t, "configurable", graphqlConfigurableProduct.Type())

	// Configurable input with active variant preselected
	configurableWithPreselectedVariantProduct := domain.ConfigurableProduct{
		Identifier: "configurable",
		BasicProductData: domain.BasicProductData{
			MarketPlaceCode: "configurable_code",
		},
		Teaser: domain.TeaserData{
			PreSelectedVariantSku: "active_variant_code",
		},
		Variants: []domain.Variant{
			{
				BasicProductData: domain.BasicProductData{
					MarketPlaceCode: "active_variant_code",
				},
				Saleable: domain.Saleable{},
			},
		},
	}
	graphqlConfigurableWithPreselectedVariantProduct := graphqlProductDto.NewGraphqlProductDto(configurableWithPreselectedVariantProduct, nil, nil).(graphqlProductDto.ActiveVariantProduct)
	assert.Equal(t, "configurable_with_activevariant", graphqlConfigurableWithPreselectedVariantProduct.Type())
	assert.Equal(t, "configurable_code", graphqlConfigurableWithPreselectedVariantProduct.MarketPlaceCode())
	assert.Equal(t, "active_variant_code", graphqlConfigurableWithPreselectedVariantProduct.VariantMarketPlaceCode())

	// Configurable input with active variant override
	configurableWithManualPreselectedVariantProduct := domain.ConfigurableProduct{
		Identifier: "configurable",
		BasicProductData: domain.BasicProductData{
			MarketPlaceCode: "configurable_code",
		},
		Teaser: domain.TeaserData{
			PreSelectedVariantSku: "variant_code",
		},
		Variants: []domain.Variant{
			{
				BasicProductData: domain.BasicProductData{
					MarketPlaceCode: "variant_code",
				},
				Saleable: domain.Saleable{},
			},
			{
				BasicProductData: domain.BasicProductData{
					MarketPlaceCode: "second_active_variant_code",
				},
				Saleable: domain.Saleable{},
			},
		},
	}
	customVariantCode := "second_active_variant_code"
	graphqlConfigurableWithManualPreselectedVariantProduct := graphqlProductDto.NewGraphqlProductDto(configurableWithManualPreselectedVariantProduct, &customVariantCode, nil).(graphqlProductDto.ActiveVariantProduct)
	assert.Equal(t, "configurable_with_activevariant", graphqlConfigurableWithManualPreselectedVariantProduct.Type())
	assert.Equal(t, "configurable_code", graphqlConfigurableWithManualPreselectedVariantProduct.MarketPlaceCode())
	assert.Equal(t, customVariantCode, graphqlConfigurableWithManualPreselectedVariantProduct.VariantMarketPlaceCode())

	// Active variant input
	activeVariantProduct := domain.ConfigurableProductWithActiveVariant{
		Identifier:       "activeVariant",
		BasicProductData: domain.BasicProductData{},
		Teaser:           domain.TeaserData{},
	}
	graphqlActiveVariantProduct := graphqlProductDto.NewGraphqlProductDto(activeVariantProduct, nil, nil)
	assert.Equal(t, "configurable_with_activevariant", graphqlActiveVariantProduct.Type())
}

func TestProductBadges_GetFirst(t *testing.T) {
	t.Parallel()

	var b graphqlProductDto.ProductBadges
	assert.Nil(t, b.First(), "get first on nil badges does not fail")

	b = graphqlProductDto.ProductBadges{}
	assert.Nil(t, b.First(), "get nil first badge on empty badges")

	b = graphqlProductDto.ProductBadges{
		domain.Badges{
			{
				Code:  "first",
				Label: "First",
			},
			{
				Code:  "second",
				Label: "Secons",
			},
		},
	}
	assert.Equal(t, &domain.Badge{Code: "first", Label: "First"}, b.First(), "got the first badge")
}
