package domain

const (
	// TYPESIMPLE denotes simple products
	TYPESIMPLE = "simple"
)

type (

	// SimpleProduct - A product without Variants that can be teasered and being sold
	SimpleProduct struct {
		Identifier string
		BasicProductData
		Saleable
		Teaser TeaserData
	}
)

// Verify Interfaces
var _ BasicProduct = SimpleProduct{}

// Type interface implementation for SimpleProduct
func (p SimpleProduct) Type() string {
	return TYPESIMPLE
}

func (p SimpleProduct) IsSaleable() bool {
	return true
}

// BaseData interface implementation for SimpleProduct
func (p SimpleProduct) BaseData() BasicProductData {
	bp := p.BasicProductData
	return bp
}

// TeaserData interface implementation for SimpleProduct
func (p SimpleProduct) TeaserData() TeaserData {
	return p.Teaser
}

// SaleableData getter for SimpleProduct
func (p SimpleProduct) SaleableData() Saleable {
	return p.Saleable
}

// GetIdentifier interface implementation for SimpleProduct
func (p SimpleProduct) GetIdentifier() string {
	return p.Identifier
}

// HasMedia  for SimpleProduct
func (p SimpleProduct) HasMedia(group string, usage string) bool {
	media := findMediaInProduct(BasicProduct(p), group, usage)
	if media == nil {
		return false
	}
	return true
}

// GetMedia  for SimpleProduct
func (p SimpleProduct) GetMedia(group string, usage string) Media {
	return *findMediaInProduct(BasicProduct(p), group, usage)
}
