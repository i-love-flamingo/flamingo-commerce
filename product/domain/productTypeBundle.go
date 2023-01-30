package domain

import (
	"errors"
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
		MinQty  int
		MaxQty  int
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

	BundleConfiguration map[Identifier]ChoiceConfiguration

	ChoiceConfiguration struct {
		MarketplaceCode        string
		VariantMarketplaceCode string
		Qty                    int
	}
)

var (
	_                                BasicProduct = BundleProduct{}
	ErrRequiredChoicesAreNotSelected              = errors.New("required choices are not selected")
	ErrSelectedQuantityOutOfRange                 = errors.New("selected quantity is out of range")
)

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

func (b BundleProduct) GetBundleProductWithActiveChoices(bundleConfiguration BundleConfiguration) (BundleProductWithActiveChoices, error) {
	bundleProductWithActiveChoices := BundleProductWithActiveChoices{
		BundleProduct: b,
		ActiveChoices: make(map[Identifier]ActiveChoice, 0),
	}

	for _, choice := range b.Choices {
		choiceConfig, ok := bundleConfiguration[Identifier(choice.Identifier)]
		if !ok && choice.Required {
			return bundleProductWithActiveChoices, ErrRequiredChoicesAreNotSelected
		}

		for _, option := range choice.Options {
			if choiceConfig.MarketplaceCode != option.Product.BaseData().MarketPlaceCode {
				continue
			}

			activeChoice, err := mapChoiceToActiveProduct(option, choice, choiceConfig)
			if err != nil {
				return BundleProductWithActiveChoices{}, fmt.Errorf("bundle product: %w", err)
			}

			bundleProductWithActiveChoices.ActiveChoices[Identifier(choice.Identifier)] = activeChoice
		}
	}

	return bundleProductWithActiveChoices, nil
}

func mapChoiceToActiveProduct(option Option, possibleChoice Choice, selectedChoice ChoiceConfiguration) (ActiveChoice, error) {
	activeChoice := ActiveChoice{}

	quantity, err := getQuantity(option.MinQty, option.MaxQty, selectedChoice.Qty)

	if err != nil {
		return ActiveChoice{}, err
	}

	switch option.Product.Type() {
	case TypeConfigurable:
		activeChoice = ActiveChoice{
			Product:  option.Product,
			Qty:      quantity,
			Label:    possibleChoice.Label,
			Required: possibleChoice.Required,
		}

		if configurable, ok := option.Product.(ConfigurableProduct); ok {
			configurableWithActiveVariant, err := configurable.GetConfigurableWithActiveVariant(selectedChoice.VariantMarketplaceCode)
			if err != nil {
				return ActiveChoice{}, fmt.Errorf("error getting configurable with active variant: %w", err)
			}

			activeChoice.Product = configurableWithActiveVariant
		}
	case TypeSimple:
		activeChoice = ActiveChoice{
			Product:  option.Product,
			Qty:      quantity,
			Label:    possibleChoice.Label,
			Required: possibleChoice.Required,
		}
	}

	return activeChoice, nil
}

func getQuantity(min, max, selected int) (int, error) {
	if selected != 0 {
		if selected >= min && selected <= max {
			return selected, nil
		}

		return 0, ErrSelectedQuantityOutOfRange
	}

	return min, nil
}