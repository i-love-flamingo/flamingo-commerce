package graphqlproductdto_test

import (
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func getProductDomainSimpleProduct() productDomain.SimpleProduct {
	return productDomain.SimpleProduct{
		Identifier: "simple_product",
		BasicProductData: productDomain.BasicProductData{
			Title: "product_title",
			Keywords: []string{
				"keywords",
			},
			MarketPlaceCode: "simple_product_code",
			Description:     "product_description",
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
					RawValue:  nil,
					UnitCode:  "attribute_a_unitCode",
				},
				"attribute_b_code": {
					Code:      "attribute_b_code",
					CodeLabel: "attribute_b_codeLabel",
					Label:     "attribute_b_label",
					RawValue:  nil,
					UnitCode:  "attribute_b_unitCode",
				},
			},
		},
		Saleable: productDomain.Saleable{},
		Teaser: productDomain.TeaserData{
			TeaserPrice: productDomain.PriceInfo{
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

func getSimpleProduct() graphqlProductDto.Product {
	product := getProductDomainSimpleProduct()
	return graphqlProductDto.NewGraphqlProductDto(product)
}

func TestSimpleProduct_Attributes(t *testing.T) {
	product := getSimpleProduct()

	assert.Equal(t, true, product.Attributes().HasAttribute("attribute_a_code"))
	assert.Equal(t, true, product.Attributes().HasAttribute("attribute_b_code"))
	assert.Equal(t, false, product.Attributes().HasAttribute("unknown"))
}

func TestSimpleProduct_Categories(t *testing.T) {
	product := getSimpleProduct()

	assert.Equal(t, productDomain.CategoryTeaser{
		Code:   "main_category",
		Path:   "main_category",
		Name:   "main_category",
		Parent: nil,
	}, product.Categories().Main)

	assert.DeepEqual(t, []productDomain.CategoryTeaser{
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

func TestSimpleProduct_Description(t *testing.T) {
	product := getSimpleProduct()
	assert.Equal(t, "product_description", product.Description())
}

func TestSimpleProduct_Loyalty(t *testing.T) {
	product := getSimpleProduct()

	assert.Equal(t, "AwesomeLoyaltyProgram", product.Loyalty().Earning.Type)
	assert.Equal(t, "AwesomeLoyaltyProgram", product.Loyalty().Price.Type)
}

func TestSimpleProduct_MarketPlaceCode(t *testing.T) {
	product := getSimpleProduct()

	assert.Equal(t, "simple_product_code", product.MarketPlaceCode())
}

func TestSimpleProduct_Media(t *testing.T) {
	product := getSimpleProduct()

	assert.DeepEqual(t, &productDomain.Media{
		Type:      "teaser",
		MimeType:  "teaser",
		Usage:     "teaser",
		Title:     "teaser",
		Reference: "teaser",
	}, product.Media().GetMedia("teaser"))
}

func TestSimpleProduct_Meta(t *testing.T) {
	product := getSimpleProduct()
	assert.DeepEqual(t, []string{"keywords"}, product.Meta().Keywords)
}

func TestSimpleProduct_Price(t *testing.T) {
	product := getSimpleProduct()
	assert.Equal(t, priceDomain.NewFromFloat(23.23, "EUR").FloatAmount(), product.Price().Default.GetPayable().FloatAmount())
}

func TestSimpleProduct_Product(t *testing.T) {
	product := getSimpleProduct()
	assert.DeepEqual(t, getProductDomainSimpleProduct().MarketPlaceCode, product.Product().BaseData().MarketPlaceCode)
}

func TestSimpleProduct_Title(t *testing.T) {
	product := getSimpleProduct()
	assert.Equal(t, "product_title", product.Title())
}

func TestSimpleProduct_Type(t *testing.T) {
	product := getSimpleProduct()
	assert.Equal(t, productDomain.TypeSimple, product.Type())
}