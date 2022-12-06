package domain

const (
	// TypeBundle denotes bundle products
	TypeBundle = "bundle"
)

type (
	Option struct {
		Product BasicProduct
		Qty     int
	}

	Choice struct {
		Identifier string
		Required   bool
		Label      string
		Options    []Option
	}

	BundleProduct struct {
		Identifier string
		Choices    []Choice
		BasicProductData
		Teaser TeaserData
	}
)

var _ BasicProduct = BundleProduct{}

func (b BundleProduct) BaseData() BasicProductData {
	return b.BasicProductData
}

func (b BundleProduct) TeaserData() TeaserData {
	return b.Teaser
}

func (b BundleProduct) IsSaleable() bool {
	return false
}

func (b BundleProduct) SaleableData() Saleable {
	return Saleable{}
}

func (b BundleProduct) Type() string {
	return TypeBundle
}

func (b BundleProduct) GetIdentifier() string {
	return b.Identifier
}

func (b BundleProduct) HasMedia(group string, usage string) bool {
	media := findMediaInProduct(BasicProduct(b), group, usage)
	return media != nil
}

func (b BundleProduct) GetMedia(group string, usage string) Media {
	return *findMediaInProduct(BasicProduct(b), group, usage)
}
