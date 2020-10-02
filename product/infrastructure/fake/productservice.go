package fake

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

var (
	brands = []string{
		"Apple",
		"Bose",
		"Dior",
		"Hugo Boss",
	}
)

// ProductService is just mocking stuff
type ProductService struct {
	currencyCode  string
	testDataFiles map[string]string
	logger        flamingo.Logger
}

// Inject dependencies
func (ps *ProductService) Inject(logger flamingo.Logger,
	c *struct {
		CurrencyCode   string `inject:"config:commerce.product.fakeservice.currency,optional"`
		TestDataFolder string `inject:"config:commerce.product.fakeservice.jsonTestDataFolder,optional"`
	},
) *ProductService {
	ps.logger = logger
	if c != nil {
		ps.currencyCode = c.CurrencyCode
		if len(c.TestDataFolder) > 0 {
			ps.testDataFiles = registerTestData(c.TestDataFolder, ps.logger)
		}
	}

	return ps
}

// Get returns a product struct
func (ps *ProductService) Get(_ context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	switch marketplaceCode {
	case "fake_configurable":
		return ps.getFakeConfigurableWithVariants(marketplaceCode), nil

	case "fake_configurable_with_active_variant":
		return ps.getFakeConfigurableWithActiveVariant(marketplaceCode), nil

	case "fake_simple":
		return ps.FakeSimple(marketplaceCode, false, false, false, true, false), nil

	case "fake_simple_with_fixed_price":
		return ps.FakeSimple(marketplaceCode, false, false, false, true, true), nil

	case "fake_fixed_simple_without_discounts":
		return ps.FakeSimple(marketplaceCode, false, false, false, false, true), nil

	case "fake_simple_out_of_stock":
		return ps.FakeSimple(marketplaceCode, false, false, true, true, false), nil
	default:
		jsonProduct, err := ps.getProductFromJSON(marketplaceCode)
		if err != nil {
			if _, isProductNotFoundError := err.(domain.ProductNotFound); !isProductNotFoundError {
				return nil, err
			}
		} else {
			return jsonProduct, nil
		}
	}

	marketPlaceCodes := ps.GetMarketPlaceCodes()
	return nil, domain.ProductNotFound{
		MarketplaceCode: "Code " + marketplaceCode + " Not implemented in FAKE: Only following codes should be used" + strings.Join(marketPlaceCodes, ", "),
	}
}

// FakeSimple generates a simple fake product
func (ps *ProductService) FakeSimple(marketplaceCode string, isNew bool, isExclusive bool, isOutOfStock bool, isDiscounted bool, hasFixedPrice bool) domain.SimpleProduct {
	product := domain.SimpleProduct{}
	product.Title = "TypeSimple product"
	ps.addBasicData(&product.BasicProductData)

	product.Saleable = domain.Saleable{
		IsSaleable:   true,
		SaleableTo:   time.Now().Add(time.Hour * time.Duration(1)),
		SaleableFrom: time.Now().Add(time.Hour * time.Duration(-1)),
		LoyaltyPrices: []domain.LoyaltyPriceInfo{
			{
				Type:    "AwesomeLoyaltyProgram",
				Default: priceDomain.NewFromFloat(500, "BonusPoints"),
			},
		},
		LoyaltyEarnings: []domain.LoyaltyEarningInfo{
			{
				Type:    "AwesomeLoyaltyProgram",
				Default: priceDomain.NewFromFloat(23.23, "BonusPoints"),
			},
		},
	}

	discountedPrice := 0.0
	if isDiscounted {
		discountedPrice = 10.49 + float64(rand.Intn(10))
		if hasFixedPrice {
			discountedPrice = 10.49
		}
	}

	defaultPrice := 20.99 + float64(rand.Intn(10))
	if hasFixedPrice {
		defaultPrice = 20.99
	}

	product.ActivePrice = ps.getPrice(defaultPrice, discountedPrice)
	product.MarketPlaceCode = marketplaceCode

	product.CreatedAt = time.Date(2019, 6, 29, 00, 00, 00, 00, time.UTC)
	product.UpdatedAt = time.Date(2019, 7, 29, 12, 00, 00, 00, time.UTC)
	product.VisibleFrom = time.Date(2019, 7, 29, 12, 00, 00, 00, time.UTC)
	product.VisibleTo = time.Now().Add(time.Hour * time.Duration(10))

	product.Teaser = domain.TeaserData{
		ShortDescription: product.ShortDescription,
		ShortTitle:       product.Title,
		URLSlug:          product.BaseData().Attributes["urlSlug"].Value(),
		Media:            product.Media,
		MarketPlaceCode:  product.MarketPlaceCode,
		TeaserPrice: domain.PriceInfo{
			Default: priceDomain.NewFromFloat(9.99, "SD").GetPayable(),
		},
		TeaserLoyaltyPriceInfo: &domain.LoyaltyPriceInfo{
			Type:    "AwesomeLoyaltyProgram",
			Default: priceDomain.NewFromFloat(500, "BonusPoints"),
		},
		TeaserLoyaltyEarningInfo: &domain.LoyaltyEarningInfo{
			Type:    "AwesomeLoyaltyProgram",
			Default: priceDomain.NewFromFloat(23.23, "BonusPoints"),
		},
	}

	if isNew {
		product.BasicProductData.IsNew = true
	}

	if isExclusive {
		product.Attributes["exclusiveProduct"] = domain.Attribute{
			RawValue: "30002654_yes",
			Code:     "exclusiveProduct",
		}
	}

	product.StockLevel = domain.StockLevelInStock
	if isOutOfStock {
		product.StockLevel = domain.StockLevelOutOfStock
	}

	return product
}

// GetMarketPlaceCodes returns list of available marketplace codes which are supported by this fakeservice
func (ps *ProductService) GetMarketPlaceCodes() []string {
	marketPlaceCodes := []string{
		"fake_configurable",
		"fake_configurable_with_active_variant",
		"fake_simple",
		"fake_simple_with_fixed_price",
		"fake_simple_out_of_stock",
		"fake_fixed_simple_without_discounts",
	}

	return append(marketPlaceCodes, ps.jsonProductCodes()...)
}

func (ps *ProductService) getFakeConfigurable(marketplaceCode string) domain.ConfigurableProduct {
	product := domain.ConfigurableProduct{}
	product.Title = "TypeConfigurable product"
	ps.addBasicData(&product.BasicProductData)
	product.MarketPlaceCode = marketplaceCode
	product.Identifier = marketplaceCode + "_identifier"
	product.Teaser.TeaserPrice = ps.getPrice(30.99+float64(rand.Intn(10)), 20.49+float64(rand.Intn(10)))
	product.VariantVariationAttributes = []string{"color", "size"}
	product.VariantVariationAttributesSorting = map[string][]string{
		"size":  {"M", "L"},
		"color": {"Red", "White", "Black"},
	}

	return product
}

func (ps *ProductService) getFakeConfigurableWithVariants(marketplaceCode string) domain.ConfigurableProduct {
	product := ps.getFakeConfigurable(marketplaceCode)
	product.RetailerCode = "retailer"

	variants := []struct {
		marketplaceCode string
		title           string
		attributes      domain.Attributes
		stockLevel      string
	}{
		{"shirt-red-s", "Shirt Red S", domain.Attributes{
			"size":                  domain.Attribute{RawValue: "S", Code: "size", CodeLabel: "Size", Label: "S"},
			"manufacturerColor":     domain.Attribute{RawValue: "red", Code: "manufacturerColor", CodeLabel: "Manufacturer Color", Label: "Red"},
			"manufacturerColorCode": domain.Attribute{RawValue: "#ff0000", Code: "manufacturerColorCode", CodeLabel: "Manufacturer Color Code", Label: "BloodRed"}},
			"high",
		},
		{"shirt-white-s", "Shirt White S", domain.Attributes{
			"size":                  domain.Attribute{RawValue: "S", Code: "size", CodeLabel: "Size", Label: "S"},
			"manufacturerColor":     domain.Attribute{RawValue: "white", Code: "manufacturerColor", CodeLabel: "Manufacturer Color", Label: "White"},
			"manufacturerColorCode": domain.Attribute{RawValue: "#ffffff", Code: "manufacturerColorCode", CodeLabel: "Manufacturer Color Code", Label: "SnowWhite"}},
			"high",
		},
		{"shirt-white-m", "Shirt White M", domain.Attributes{
			"size":  domain.Attribute{RawValue: "M", Code: "size", CodeLabel: "Size", Label: "M"},
			"color": domain.Attribute{RawValue: "white", Code: "color", CodeLabel: "Color", Label: "White"}},
			"high",
		},
		{"shirt-black-m", "Shirt Black M", domain.Attributes{
			"size":                  domain.Attribute{RawValue: "M", Code: "size", CodeLabel: "Size", Label: "M"},
			"manufacturerColor":     domain.Attribute{RawValue: "blue", Code: "manufacturerColor", CodeLabel: "Manufacturer Color", Label: "Blue"},
			"manufacturerColorCode": domain.Attribute{RawValue: "#0000ff", Code: "manufacturerColorCode", CodeLabel: "Manufacturer Color Code", Label: "SkyBlue"}},
			"high",
		},
		{"shirt-black-l", "Shirt Black L", domain.Attributes{
			"size":  domain.Attribute{RawValue: "L", Code: "size", CodeLabel: "Size", Label: "L"},
			"color": domain.Attribute{RawValue: "black", Code: "color", CodeLabel: "Color", Label: "Black"}},
			"high",
		},
		{"shirt-red-l", "Shirt Red L", domain.Attributes{
			"size":  domain.Attribute{RawValue: "L", Code: "size", CodeLabel: "Size", Label: "L"},
			"color": domain.Attribute{RawValue: "red", Code: "color", CodeLabel: "Color", Label: "Red"}},
			"out",
		},
		{"shirt-red-m", "Shirt Red M", domain.Attributes{
			"size":  domain.Attribute{RawValue: "M", Code: "size", CodeLabel: "Size", Label: "M"},
			"color": domain.Attribute{RawValue: "red", Code: "color", CodeLabel: "Color", Label: "Red"}},
			"out",
		},
	}

	for _, variant := range variants {
		simpleVariant := ps.fakeVariant(variant.marketplaceCode)
		simpleVariant.Title = variant.title
		simpleVariant.Attributes = variant.attributes
		simpleVariant.BasicProductData.Attributes = variant.attributes
		simpleVariant.StockLevel = variant.stockLevel

		// Give new images for variants with custom colors
		if simpleVariant.Attributes.HasAttribute("manufacturerColorCode") {
			manufacturerColorCode := simpleVariant.Attributes["manufacturerColorCode"].Value()
			manufacturerColorCode = strings.TrimPrefix(manufacturerColorCode, "#")
			simpleVariant.Media[0] = domain.Media{Type: "image-external", Reference: "http://dummyimage.com/1024x768/000/" + manufacturerColorCode, Usage: "detail"}
		}

		product.Variants = append(product.Variants, simpleVariant)
	}

	return product
}

func (ps *ProductService) getFakeConfigurableWithActiveVariant(marketplaceCode string) domain.ConfigurableProductWithActiveVariant {
	configurable := ps.getFakeConfigurableWithVariants(marketplaceCode)
	product := domain.ConfigurableProductWithActiveVariant{
		Identifier:                        configurable.Identifier,
		BasicProductData:                  configurable.BasicProductData,
		Teaser:                            configurable.Teaser,
		VariantVariationAttributes:        configurable.VariantVariationAttributes,
		VariantVariationAttributesSorting: configurable.VariantVariationAttributesSorting,
		Variants:                          configurable.Variants,
		ActiveVariant:                     configurable.Variants[4], // shirt-black-l
	}

	product.Teaser.TeaserPrice = product.ActiveVariant.ActivePrice

	return product
}

func (ps *ProductService) fakeVariant(marketplaceCode string) domain.Variant {
	var simpleVariant domain.Variant
	simpleVariant.Attributes = make(map[string]domain.Attribute)

	ps.addBasicData(&simpleVariant.BasicProductData)

	simpleVariant.ActivePrice = ps.getPrice(30.99+float64(rand.Intn(10)), 20.49+float64(rand.Intn(10)))
	simpleVariant.MarketPlaceCode = marketplaceCode
	simpleVariant.IsSaleable = true

	return simpleVariant
}

func (ps *ProductService) addBasicData(product *domain.BasicProductData) {
	product.ShortDescription = "Short Description"
	product.Description = "Description"
	product.Keywords = []string{"keywords"}

	product.Media = append(product.Media, domain.Media{Type: "image-external", Reference: "http://dummyimage.com/1024x768/000/fff", Usage: "detail"})
	product.Media = append(product.Media, domain.Media{Type: "image-external", Reference: "http://dummyimage.com/200x200/000/fff", Usage: "list"})

	product.Attributes = domain.Attributes{
		"brandCode":        domain.Attribute{RawValue: brands[rand.Intn(len(brands))]},
		"brandName":        domain.Attribute{RawValue: brands[rand.Intn(len(brands))]},
		"collectionOption": domain.Attribute{RawValue: []interface{}{"departure", "arrival"}},
		"urlSlug":          domain.Attribute{RawValue: "product-slug"},
	}

	product.RetailerCode = "retailer"
	product.RetailerName = "Test Retailer"
	product.RetailerSku = "12345sku"

	categoryTeaser1 := domain.CategoryTeaser{
		Path: "Testproducts",
		Name: "Testproducts",
		Code: "testproducts",
	}
	categoryTeaser2 := domain.CategoryTeaser{
		Path: "Testproducts/Fake/Configurable",
		Name: "Configurable",
		Code: "configurable",
	}
	product.Categories = append(product.Categories, categoryTeaser1)
	product.Categories = append(product.Categories, categoryTeaser2)
	product.MainCategory = categoryTeaser1
}

func (ps *ProductService) getPrice(defaultP float64, discounted float64) domain.PriceInfo {
	defaultP, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", defaultP), 64)
	discounted, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", discounted), 64)

	var price domain.PriceInfo
	currency := "EUR"
	if ps.currencyCode != "" {
		currency = ps.currencyCode
	}

	price.Default = priceDomain.NewFromFloat(defaultP, currency).GetPayable()
	if discounted > 0 {
		price.Discounted = priceDomain.NewFromFloat(discounted, currency).GetPayable()
		price.DiscountText = "Super test campaign"
		price.IsDiscounted = true
	}
	price.ActiveBase = *big.NewFloat(1)
	price.ActiveBaseAmount = *big.NewFloat(10)
	price.ActiveBaseUnit = "ml"
	return price
}

func (ps *ProductService) getProductFromJSON(code string) (domain.BasicProduct, error) {
	file, ok := ps.testDataFiles[code]

	if !ok {
		return nil, &domain.ProductNotFound{MarketplaceCode: code}
	}

	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return unmarshalJSONProduct(jsonBytes)
}

// jsonProductCodes returns an ordered list of the json product codes
func (ps *ProductService) jsonProductCodes() []string {
	keys := make([]string, 0, len(ps.testDataFiles))
	for k := range ps.testDataFiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
