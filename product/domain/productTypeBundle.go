package domain

import (
	"encoding/json"
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
		Saleable
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
	ErrMarketplaceCodeDoNotExists                 = errors.New("selected marketplace code does not exist")
)

func (b BundleProduct) BaseData() BasicProductData {
	return b.BasicProductData
}

func (b BundleProduct) TeaserData() TeaserData {
	return b.Teaser
}

func (b BundleProduct) IsSaleable() bool {
	return true
}

func (b BundleProduct) SaleableData() Saleable {
	return b.Saleable
}

func (b BundleProduct) Type() string {
	return TypeBundle
}

func (b BundleProductWithActiveChoices) Type() string {
	return TypeBundleWithActiveChoices
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
			return BundleProductWithActiveChoices{}, ErrRequiredChoicesAreNotSelected
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

		if _, ok := bundleProductWithActiveChoices.ActiveChoices[Identifier(choice.Identifier)]; !ok && choice.Required {
			return BundleProductWithActiveChoices{}, ErrMarketplaceCodeDoNotExists
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
			Identifier: possibleChoice.Identifier,
			Product:    option.Product,
			Qty:        quantity,
			Label:      possibleChoice.Label,
			Required:   possibleChoice.Required,
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
			Identifier: possibleChoice.Identifier,
			Product:    option.Product,
			Qty:        quantity,
			Label:      possibleChoice.Label,
			Required:   possibleChoice.Required,
		}
	}

	return activeChoice, nil
}

func getQuantity(min, max, selected int) (int, error) {
	if selected >= min && selected <= max {
		return selected, nil
	}

	return 0, ErrSelectedQuantityOutOfRange
}

func (o *Option) UnmarshalJSON(optionData []byte) error {
	option := &struct {
		Product json.RawMessage
		MinQty  int
		MaxQty  int
	}{}

	err := json.Unmarshal(optionData, option)
	if err != nil {
		return errors.New("option product: " + err.Error())
	}

	product := &map[string]interface{}{}
	err = json.Unmarshal(option.Product, product)
	if err != nil {
		return errors.New("option product: " + err.Error())
	}

	productType, ok := (*product)["Type"]

	if !ok {
		return errors.New("option product: type is not specified")
	}

	o.MinQty = option.MinQty
	o.MaxQty = option.MaxQty

	switch productType {
	case TypeConfigurable:
		configurableProduct := &ConfigurableProduct{}
		err = json.Unmarshal(option.Product, configurableProduct)
		if err != nil {
			return errors.New("option product: " + err.Error())
		}
		o.Product = *configurableProduct
	default:
		simpleProduct := &SimpleProduct{}
		err = json.Unmarshal(option.Product, simpleProduct)
		if err != nil {
			return errors.New("option product: " + err.Error())
		}
		o.Product = *simpleProduct
	}
	return nil
}

func (b BundleProductWithActiveChoices) ExtractBundleConfig() BundleConfiguration {
	if len(b.ActiveChoices) == 0 {
		return nil
	}

	config := make(BundleConfiguration)

	for identifier, choice := range b.ActiveChoices {
		if choice.Product.Type() == TypeSimple {
			config[identifier] = ChoiceConfiguration{
				MarketplaceCode: choice.Product.BaseData().MarketPlaceCode,
				Qty:             choice.Qty,
			}
		}
		if choice.Product.Type() == TypeConfigurableWithActiveVariant {
			configurableWithActiveVariant, ok := choice.Product.(ConfigurableProductWithActiveVariant)
			if ok {
				config[identifier] = ChoiceConfiguration{
					MarketplaceCode:        configurableWithActiveVariant.ConfigurableBaseData().MarketPlaceCode,
					VariantMarketplaceCode: choice.Product.BaseData().MarketPlaceCode,
					Qty:                    choice.Qty,
				}
			}
		}
	}

	return config
}

// Equals compares the marketplace codes of all choices
func (bc BundleConfiguration) Equals(other BundleConfiguration) bool {
	if len(bc) != len(other) {
		return false
	}

	for choiceID, cartChoiceConfig := range bc {
		otherChoiceConfig, ok := other[choiceID]
		if !ok {
			return false
		}

		if cartChoiceConfig.MarketplaceCode != otherChoiceConfig.MarketplaceCode {
			return false
		}

		if cartChoiceConfig.VariantMarketplaceCode != otherChoiceConfig.VariantMarketplaceCode {
			return false
		}

		if cartChoiceConfig.Qty != otherChoiceConfig.Qty {
			return false
		}
	}

	return true
}
