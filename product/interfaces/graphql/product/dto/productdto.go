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
		Identifier() string
		Media() ProductMedia
		Price() productDomain.PriceInfo
		AvailablePrices() []productDomain.PriceInfo
		Title() string
		Categories() ProductCategories
		Description() string
		ShortDescription() string
		Meta() ProductMeta
		Loyalty() ProductLoyalty
		Attributes() productDomain.Attributes
		Badges() ProductBadges
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

	// ProductBadges wrapper for badges of the product
	ProductBadges struct {
		All []productDomain.Badge
	}

	// VariationSelection represents possible combinations for attached variants
	VariationSelection struct {
		Code    string
		Label   string
		Options []VariationSelectionOption
	}

	// VariationSelectionOption one possible variation option
	VariationSelectionOption struct {
		Label   string
		State   VariationSelectionOptionState
		Variant VariationSelectionOptionVariant
	}

	// VariationSelectionOptionState state of the option
	VariationSelectionOptionState string

	// ActiveVariationSelection The variation for the currently active variant
	ActiveVariationSelection struct {
		Code  string
		Label string
		Value string
	}

	// VariationSelectionOptionVariant Information about the underlying variant
	VariationSelectionOptionVariant struct {
		variant productDomain.Variant
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
func NewGraphqlProductDto(product productDomain.BasicProduct, preSelectedVariantSku *string) Product {
	if product.Type() == productDomain.TypeConfigurable {
		configurableProduct := product.(productDomain.ConfigurableProduct)

		variantSku := ""

		if configurableProduct.Teaser.PreSelectedVariantSku != "" {
			variantSku = configurableProduct.Teaser.PreSelectedVariantSku
		}

		if preSelectedVariantSku != nil && *preSelectedVariantSku != "" {
			variantSku = *preSelectedVariantSku
		}

		if variantSku != "" {
			configurableProductWithActiveVariant, err := configurableProduct.GetConfigurableWithActiveVariant(variantSku)

			if err == nil {
				return ActiveVariantProduct{
					product: configurableProductWithActiveVariant,
				}
			}
		}

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

// NewVariationSelectionOptionVariant Creates a new option variant from the domain variant
func NewVariationSelectionOptionVariant(variant productDomain.Variant) VariationSelectionOptionVariant {
	return VariationSelectionOptionVariant{variant}
}

// MarketPlaceCode returns the marketPlaceCode of the variant
func (v *VariationSelectionOptionVariant) MarketPlaceCode() string {
	return v.variant.MarketPlaceCode
}

// BaseData of the variant
func (v *VariationSelectionOptionVariant) BaseData() productDomain.BasicProductData {
	return v.variant.BaseData()
}

// First badge of all badges, returns nil if there is no first badge
func (b *ProductBadges) First() *productDomain.Badge {
	badges := productDomain.Badges(b.All)

	return badges.First()
}
