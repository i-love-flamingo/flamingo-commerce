package graphqlproductdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
)

func getProductDomainConfigurableProduct() productDomain.ConfigurableProduct {
	return productDomain.ConfigurableProduct{
		Identifier: "configurable_product",
		BasicProductData: productDomain.BasicProductData{
			Title: "product_title",
			Keywords: []string{
				"keywords",
			},
			MarketPlaceCode:  "configurable_product_code",
			Description:      "product_description",
			ShortDescription: "product_description_short",
			MainCategory: productDomain.CategoryTeaser{
				Code:   "main_category",
				Path:   "main_category",
				Name:   "main_category",
				Parent: nil,
			},
			Categories: []productDomain.CategoryTeaser{
				{
					Code:   "category_a",
					Path:   "category_a",
					Name:   "category_a",
					Parent: nil,
				},
				{
					Code:   "category_b",
					Path:   "category_b",
					Name:   "category_b",
					Parent: nil,
				},
			},
			Attributes: productDomain.Attributes{
				"attribute_a_code": {
					Code:      "attribute_a_code",
					CodeLabel: "attribute_a_codeLabel",
					Label:     "attribute_a_label",
					RawValue:  "attribute_a_variantValue",
					UnitCode:  "attribute_a_unitCode",
				},
				"attribute_b_code": {
					Code:      "attribute_b_code",
					CodeLabel: "attribute_b_codeLabel",
					Label:     "attribute_b_label",
					RawValue:  "attribute_b_variantValue",
					UnitCode:  "attribute_b_unitCode",
				},
			},
			Badges: []productDomain.Badge{
				{
					Code:  "hot",
					Label: "Hot Product",
				},
			},
		},
		VariantVariationAttributes: []string{
			"attribute_a_code",
		},
		VariantVariationAttributesSorting: map[string][]string{
			"attribute_a_code": {
				"attribute_a_variantValue",
				"attribute_b_variantValue",
			},
		},
		Variants: []productDomain.Variant{
			{
				BasicProductData: productDomain.BasicProductData{
					MarketPlaceCode: "variant_product_code",
					Attributes: productDomain.Attributes{
						"attribute_a_code": {
							Code:      "attribute_a_code",
							CodeLabel: "attribute_a_codeLabel",
							Label:     "attribute_a_variantLabel",
							RawValue:  "attribute_a_variantValue",
							UnitCode:  "attribute_a_unitCode",
						},
						"attribute_b_code": {
							Code:      "attribute_b_code",
							CodeLabel: "attribute_b_codeLabel",
							Label:     "attribute_b_variantLabel",
							RawValue:  "attribute_b_variantValue",
							UnitCode:  "attribute_b_unitCode",
						},
					},
					Badges: []productDomain.Badge{
						{
							Code:  "newvariant",
							Label: "New Variant Product",
						},
					},
				},
				Saleable: productDomain.Saleable{},
			},
		},
		Teaser: productDomain.TeaserData{
			TeaserLoyaltyPriceInfo: &productDomain.LoyaltyPriceInfo{
				Type:    "AwesomeLoyaltyProgram",
				Default: priceDomain.NewFromFloat(500, "BonusPoints"),
			},
			TeaserLoyaltyEarningInfo: &productDomain.LoyaltyEarningInfo{
				Type:    "AwesomeLoyaltyProgram",
				Default: priceDomain.NewFromFloat(23.23, "BonusPoints"),
			},

			Media: []productDomain.Media{
				{
					Type:      "teaser",
					MimeType:  "teaser",
					Usage:     "teaser",
					Title:     "teaser",
					Reference: "teaser",
				},
			},
		},
	}
}

func getConfigurableProduct() graphqlProductDto.Product {
	product := getProductDomainConfigurableProduct()
	return graphqlProductDto.NewGraphqlProductDto(product, nil, nil)
}

func TestConfigurableProduct_Attributes(t *testing.T) {
	product := getConfigurableProduct()

	assert.Equal(t, true, product.Attributes().HasAttribute("attribute_a_code"))
	assert.Equal(t, true, product.Attributes().HasAttribute("attribute_b_code"))
	assert.Equal(t, false, product.Attributes().HasAttribute("unknown"))
}

func TestConfigurableProduct_Categories(t *testing.T) {
	product := getConfigurableProduct()

	assert.Equal(t, productDomain.CategoryTeaser{
		Code:   "main_category",
		Path:   "main_category",
		Name:   "main_category",
		Parent: nil,
	}, product.Categories().Main)

	assert.Equal(t, []productDomain.CategoryTeaser{
		{
			Code:   "category_a",
			Path:   "category_a",
			Name:   "category_a",
			Parent: nil,
		},
		{
			Code:   "category_b",
			Path:   "category_b",
			Name:   "category_b",
			Parent: nil,
		},
	}, product.Categories().All)
}

func TestConfigurableProduct_Description(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(t, "product_description", product.Description())
}

func TestConfigurableProduct_ShortDescription(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(t, "product_description_short", product.ShortDescription())
}

func TestConfigurableProduct_Loyalty(t *testing.T) {
	product := getConfigurableProduct()

	assert.Equal(t, "AwesomeLoyaltyProgram", product.Loyalty().Earning.Type)
	assert.Equal(t, "AwesomeLoyaltyProgram", product.Loyalty().Price.Type)
}

func TestConfigurableProduct_MarketPlaceCode(t *testing.T) {
	product := getConfigurableProduct()

	assert.Equal(t, "configurable_product_code", product.MarketPlaceCode())
}

func TestConfigurableProduct_Identifier(t *testing.T) {
	product := getConfigurableProduct()

	assert.Equal(t, "configurable_product", product.Identifier())
}

func TestConfigurableProduct_Media(t *testing.T) {
	product := getConfigurableProduct()

	assert.Equal(t, &productDomain.Media{
		Type:      "teaser",
		MimeType:  "teaser",
		Usage:     "teaser",
		Title:     "teaser",
		Reference: "teaser",
	}, product.Media().GetMedia("teaser"))
}

func TestConfigurableProduct_Meta(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(t, []string{"keywords"}, product.Meta().Keywords)
}

func TestConfigurableProduct_Price(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(t, productDomain.PriceInfo{}, product.Price())
}

func TestConfigurableProduct_AvailablePrices(t *testing.T) {
	product := getConfigurableProduct()
	assert.Nil(t, product.AvailablePrices())
}

func TestConfigurableProduct_Product(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(t, getProductDomainConfigurableProduct().MarketPlaceCode, product.Product().BaseData().MarketPlaceCode)
}

func TestConfigurableProduct_Title(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(t, "product_title", product.Title())
}

func TestConfigurableProduct_Type(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(t, productDomain.TypeConfigurable, product.Type())
}

func TestConfigurableProduct_Badges(t *testing.T) {
	product := getConfigurableProduct()
	assert.Equal(
		t,
		[]productDomain.Badge{
			{
				Code:  "hot",
				Label: "Hot Product",
			},
		},
		product.Badges().All,
	)
}
