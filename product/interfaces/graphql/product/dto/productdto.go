package graphqlproductdto

import (
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

const (
	// VariationSelectionOptionStateActive signals, that option is currently active because the active product has this attribute
	VariationSelectionOptionStateActive VariationSelectionOptionState = "ACTIVE"
	// VariationSelectionOptionStateMatch signals, that product exists for this option but is not the active variant
	VariationSelectionOptionStateMatch VariationSelectionOptionState = "MATCH"
	// VariationSelectionOptionStateNoMatch signals, that product does not exists for this option, fallback is used
	VariationSelectionOptionStateNoMatch VariationSelectionOptionState = "NO_MATCH"
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

	// ProductLoyalty contains all loyalty related information
	ProductLoyalty struct {
		Price   *productDomain.LoyaltyPriceInfo
		Earning *productDomain.LoyaltyEarningInfo
	}

	// ProductMedia returns media for the product
	ProductMedia struct {
		All []productDomain.Media
	}

	// ProductCategories wrapper for categories
	ProductCategories struct {
		Main productDomain.CategoryTeaser
		All  []productDomain.CategoryTeaser
	}

	// ProductMeta contains meta information about the product
	ProductMeta struct {
		Keywords []string
	}

	// VariationSelection represents possible combinations for attached variants
	VariationSelection struct {
		Code    string
		Label   string
		Options []VariationSelectionOption
	}

	// VariationSelectionOption one possible variation option
	VariationSelectionOption struct {
		Label                  string
		State                  VariationSelectionOptionState
		VariantMarketPlaceCode string
	}

	// VariationSelectionOptionState state of the option
	VariationSelectionOptionState string

	// ActiveVariationSelection The variation for the currently active variant
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

// NewGraphqlProductDto returns a new Product dto
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