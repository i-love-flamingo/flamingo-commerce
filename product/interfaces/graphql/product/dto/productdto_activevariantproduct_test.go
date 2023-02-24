package graphqlproductdto_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
)

func getProductDomainConfigurableWithActiveVariantProduct() productDomain.ConfigurableProductWithActiveVariant {
	return productDomain.ConfigurableProductWithActiveVariant{
		Identifier: "configurable_with_active_variant_product",
		BasicProductData: productDomain.BasicProductData{
			MarketPlaceCode: "configurable_with_active_variant_product",
			Badges: []productDomain.Badge{
				{
					Code:  "hot",
					Label: "Hot Product",
				},
			},
		},

		ActiveVariant: productDomain.Variant{
			BasicProductData: productDomain.BasicProductData{
				Title: "product_title",
				Keywords: []string{
					"keywords",
				},
				MarketPlaceCode:  "active_variant_product_code_a",
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
				Media: []productDomain.Media{
					{
						Type:      "teaser",
						MimeType:  "teaser",
						Usage:     "teaser",
						Title:     "teaser",
						Reference: "teaser",
					},
				},
				Badges: []productDomain.Badge{
					{
						Code:  "hotvariant",
						Label: "Hot Variant Product",
					},
				},
			},
			Saleable: productDomain.Saleable{
				ActivePrice: productDomain.PriceInfo{
					Default:           priceDomain.NewFromFloat(23.23, "EUR"),
					Discounted:        priceDomain.Price{},
					DiscountText:      "",
					ActiveBase:        big.Float{},
					ActiveBaseAmount:  big.Float{},
					ActiveBaseUnit:    "",
					IsDiscounted:      false,
					CampaignRules:     nil,
					DenyMoreDiscounts: false,
					Context:           productDomain.PriceContext{},
					TaxClass:          "",
				},
				AvailablePrices: []productDomain.PriceInfo{{
					Default: priceDomain.NewFromFloat(10.00, "EUR"),
					Context: productDomain.PriceContext{CustomerGroup: "gold-members"},
				}},
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
					MarketPlaceCode: "active_variant_product_code_a",
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
				},
				Saleable: productDomain.Saleable{},
			},
			{
				BasicProductData: productDomain.BasicProductData{
					MarketPlaceCode: "active_variant_product_code_b",
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
			PreSelectedVariantSku: "active_variant_product_code_a",
			Badges: []productDomain.Badge{
				{
					Code:  "new",
					Label: "New Product",
				},
			},
		},
	}
}

func getActiveVariantProduct() graphqlProductDto.ActiveVariantProduct {
	product := getProductDomainConfigurableWithActiveVariantProduct()
	return graphqlProductDto.NewGraphqlProductDto(product, nil, nil).(graphqlProductDto.ActiveVariantProduct)
}

func TestActiveVariantProduct_Attributes(t *testing.T) {
	product := getActiveVariantProduct()

	assert.Equal(t, true, product.Attributes().HasAttribute("attribute_a_code"))
	assert.Equal(t, true, product.Attributes().HasAttribute("attribute_b_code"))
	assert.Equal(t, false, product.Attributes().HasAttribute("unknown"))
}

func TestActiveVariantProduct_Categories(t *testing.T) {
	product := getActiveVariantProduct()

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

func TestActiveVariantProduct_Description(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, "product_description", product.Description())
}

func TestActiveVariantProduct_ShortDescription(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, "product_description_short", product.ShortDescription())
}

func TestActiveVariantProduct_Loyalty(t *testing.T) {
	product := getActiveVariantProduct()

	assert.Equal(t, "AwesomeLoyaltyProgram", product.Loyalty().Earning.Type)
	assert.Equal(t, "AwesomeLoyaltyProgram", product.Loyalty().Price.Type)
}

func TestActiveVariantProduct_MarketPlaceCode(t *testing.T) {
	product := getActiveVariantProduct()

	assert.Equal(t, "configurable_with_active_variant_product", product.MarketPlaceCode())
}

func TestActiveVariantProduct_VariantMarketPlaceCode(t *testing.T) {
	product := getActiveVariantProduct()

	assert.Equal(t, "active_variant_product_code_a", product.VariantMarketPlaceCode())
}

func TestActiveVariantProduct_Identifier(t *testing.T) {
	product := getActiveVariantProduct()

	assert.Equal(t, "configurable_with_active_variant_product", product.Identifier())
}

func TestActiveVariantProduct_Media(t *testing.T) {
	product := getActiveVariantProduct()

	assert.Equal(t, &productDomain.Media{
		Type:      "teaser",
		MimeType:  "teaser",
		Usage:     "teaser",
		Title:     "teaser",
		Reference: "teaser",
	}, product.Media().GetMedia("teaser"))
}

func TestActiveVariantProduct_Meta(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, []string{"keywords"}, product.Meta().Keywords)
}

func TestActiveVariantProduct_Price(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, priceDomain.NewFromFloat(23.23, "EUR").FloatAmount(), product.Price().Default.GetPayable().FloatAmount())
}

func TestActiveVariantProduct_AvailablePrices(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, []productDomain.PriceInfo{{
		Default: priceDomain.NewFromFloat(10.00, "EUR"),
		Context: productDomain.PriceContext{CustomerGroup: "gold-members"},
	}}, product.AvailablePrices())
}

func TestActiveVariantProduct_Product(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, "active_variant_product_code_a", product.Product().BaseData().MarketPlaceCode)
}

func TestActiveVariantProduct_Title(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, "product_title", product.Title())
}

func TestActiveVariantProduct_Type(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(t, productDomain.TypeConfigurableWithActiveVariant, product.Type())
}

func TestActiveVariantProduct_VariationSelections(t *testing.T) {
	configurableProduct := getProductDomainConfigurableWithActiveVariantProduct()
	product := getActiveVariantProduct()
	assert.Equal(t, []graphqlProductDto.VariationSelection{
		{
			Code:  "attribute_a_code",
			Label: "attribute_a_codeLabel",
			Options: []graphqlProductDto.VariationSelectionOption{
				{
					Label:    "attribute_a_variantLabel",
					UnitCode: "attribute_a_unitCode",
					State:    graphqlProductDto.VariationSelectionOptionStateActive,
					Variant:  graphqlProductDto.NewVariationSelectionOptionVariant(configurableProduct.Variants[0]),
				},
			},
		},
	}, product.VariationSelections())
}

func TestActiveVariantProduct_ActiveVariationSelections(t *testing.T) {
	product := getActiveVariantProduct()

	assert.Equal(t, []graphqlProductDto.ActiveVariationSelection{{
		Code:     "attribute_a_code",
		Label:    "attribute_a_codeLabel",
		Value:    "attribute_a_variantLabel",
		UnitCode: "attribute_a_unitCode",
	}}, product.ActiveVariationSelections())
}

func TestActiveVariantProduct_Badges(t *testing.T) {
	product := getActiveVariantProduct()
	assert.Equal(
		t,
		[]productDomain.Badge{
			{
				Code:  "hotvariant",
				Label: "Hot Variant Product",
			},
		},
		product.Badges().All,
	)
}
