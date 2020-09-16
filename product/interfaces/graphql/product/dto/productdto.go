package graphqlProductDto

import (
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

const (
	VariationSelectionOptionStateActive  VariationSelectionOptionState = "ACTIVE"   // Option is currently active
	VariationSelectionOptionStateMatch   VariationSelectionOptionState = "MATCH"    // Product exists for this option
	VariationSelectionOptionStateNoMatch VariationSelectionOptionState = "NO_MATCH" // Product does not exists for this option
)

// Product contains normalized Product information regardless of being a variant or simple Product
type (

	// Product interface - needs to be implemented by SimpleProducts and ConfigurableProducts
	Product interface {
		Product() productDomain.BasicProduct
		Type() string
		MarketPlaceCode() string
		Media() ProductMedia
		Price() productDomain.PriceInfo
		Title() string
		Categories() ProductCategories
		Description() string
		Meta() ProductMeta
		Loyalty() ProductLoyalty
		Attributes() productDomain.Attributes
	}

	// All loyalty related information
	ProductLoyalty struct {
		Price   *productDomain.LoyaltyPriceInfo
		Earning *productDomain.LoyaltyEarningInfo
	}

	ProductMedia struct {
		All []productDomain.Media
	}

	// ProductCategories
	ProductCategories struct {
		Main productDomain.CategoryTeaser
		All  []productDomain.CategoryTeaser
	}

	// Normalized ProductMeta data
	ProductMeta struct {
		Keywords []string
	}

	// A selection for a Product variation
	VariationSelection struct {
		Code    string
		Label   string
		Options []VariationSelectionOption
	}

	// One possible variation for the Product
	VariationSelectionOption struct {
		Code                   string
		Label                  string
		State                  VariationSelectionOptionState
		VariantMarketPlaceCode string
	}

	// Possible state of option depending on active variant
	VariationSelectionOptionState string

	// The variation for the currently active variant
	ActiveVariationSelection struct {
		AttributeLabel string
		OptionLabel    string
	}
)

var (
	_ Product = SimpleProduct{}
	_ Product = ConfigurableProduct{}
	_ Product = ActiveVariantProduct{}
)

// GetMedia returns the FIRST found Product media by usage
func (pm ProductMedia) GetMedia(usage string) *productDomain.Media {
	for _, media := range pm.All {
		if media.Usage == usage {
			return &media
		}
	}
	return nil
}

func NewGraphqlProductDto(product productDomain.BasicProduct) Product {
	if product.Type() == productDomain.TypeConfigurable {
		configurableProduct := product.(productDomain.ConfigurableProduct)
		return ConfigurableProduct{
			product: configurableProduct,
		}
	}

	if product.Type() == productDomain.TypeConfigurableWithActiveVariant {
		configurableProduct := product.(productDomain.ConfigurableProductWithActiveVariant)
		return ActiveVariantProduct{
			product: configurableProduct,
		}
	}

	simpleProduct := product.(productDomain.SimpleProduct)
	return SimpleProduct{
		product: simpleProduct,
	}
}
