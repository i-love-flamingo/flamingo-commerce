package domain

import (
	"fmt"
)

const (
	// TypeBundle denotes bundle products
	TypeBundle                  = "bundle"
	TypeBundleWithActiveChoices = "bundle_with_active_choices"
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

	ActiveChoice struct {
		Identifier string
		Required   bool
		Label      string
		Product    BasicProduct
		Qty        int
	}

	// BundleProductWithActiveChoices - A product with preselected choices
	BundleProductWithActiveChoices struct {
		BundleProduct
		ActiveChoices map[Identifier]ActiveChoice
	}

	Identifier string

	BundleConfiguration struct {
		Identifier             string
		MarketplaceCode        string
		VariantMarketplaceCode string
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

func (b BundleProduct) GetBundleProductWithActiveChoices(bundleConfiguration []*BundleConfiguration) (BundleProductWithActiveChoices, error) {
	bundleProductWithActiveChoices := BundleProductWithActiveChoices{
		BundleProduct: b,
		ActiveChoices: make(map[Identifier]ActiveChoice, 0),
	}

	for _, selectedChoice := range bundleConfiguration {
		for _, possibleChoice := range b.Choices {
			if selectedChoice.Identifier != possibleChoice.Identifier {
				continue
			}

			for _, option := range possibleChoice.Options {
				if selectedChoice.MarketplaceCode != option.Product.BaseData().MarketPlaceCode {
					continue
				}

				activeChoice, err := mapChoiceToActiveProduct(option, possibleChoice, selectedChoice)
				if err != nil {
					return BundleProductWithActiveChoices{}, fmt.Errorf("bundle product: %w", err)
				}

				bundleProductWithActiveChoices.ActiveChoices[Identifier(selectedChoice.Identifier)] = activeChoice
			}
		}
	}

	return bundleProductWithActiveChoices, nil
}

func mapChoiceToActiveProduct(option Option, possibleChoice Choice, selectedChoice *BundleConfiguration) (ActiveChoice, error) {
	activeChoice := ActiveChoice{}

	if option.Product.Type() == TypeSimple {
		activeChoice = ActiveChoice{
			Product:  option.Product,
			Qty:      option.Qty,
			Label:    possibleChoice.Label,
			Required: possibleChoice.Required,
		}
	}
	if option.Product.Type() == TypeConfigurable {
		activeChoice = ActiveChoice{
			Product:  option.Product,
			Qty:      option.Qty,
			Label:    possibleChoice.Label,
			Required: possibleChoice.Required,
		}

		if selectedChoice.VariantMarketplaceCode != "" {
			if configurable, ok := option.Product.(ConfigurableProduct); ok {
				configurableWithActiveVariant, err := configurable.GetConfigurableWithActiveVariant(selectedChoice.VariantMarketplaceCode)
				if err != nil {
					return ActiveChoice{}, fmt.Errorf("error getting configurable with active variant: %w", err)
				}

				activeChoice.Product = configurableWithActiveVariant
			}
		}
	}

	return activeChoice, nil
}
