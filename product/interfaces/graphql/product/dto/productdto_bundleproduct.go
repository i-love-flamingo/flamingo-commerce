package graphqlproductdto

import productDomain "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	// BundleProduct A bundle Product with options
	BundleProduct struct {
		product productDomain.BundleProduct
		Choices []Choice
	}

	Choice struct {
		Identifier string
		Required   bool
		Label      string
		Options    []Option
		Active     Product
	}

	Option struct {
		Product Product
		Qty     int
	}
)

// Type the product type
func (sp BundleProduct) Type() string {
	return productDomain.TypeBundle
}

// Product the bundle product information
func (sp BundleProduct) Product() productDomain.BasicProduct {
	return sp.product
}

// MarketPlaceCode of the product
func (sp BundleProduct) MarketPlaceCode() string {
	return sp.product.BaseData().MarketPlaceCode
}

// Identifier of the product
func (sp BundleProduct) Identifier() string {
	return sp.product.GetIdentifier()
}

// Media of the product
func (sp BundleProduct) Media() ProductMedia {
	return ProductMedia{All: sp.product.TeaserData().Media}
}

// Price of the product
func (sp BundleProduct) Price() productDomain.PriceInfo {
	return sp.product.Saleable.ActivePrice
}

// AvailablePrices of the product
func (sp BundleProduct) AvailablePrices() []productDomain.PriceInfo {
	return sp.product.Saleable.AvailablePrices
}

// Title of the product
func (sp BundleProduct) Title() string {
	return sp.product.BaseData().Title
}

// Categories of the product
func (sp BundleProduct) Categories() ProductCategories {
	return ProductCategories{
		Main: sp.product.BaseData().MainCategory,
		All:  sp.product.BaseData().Categories,
	}
}

// Description of the product
func (sp BundleProduct) Description() string {
	return sp.product.BaseData().Description
}

// ShortDescription of the product
func (sp BundleProduct) ShortDescription() string {
	return sp.product.BaseData().ShortDescription
}

// Meta of the product
func (sp BundleProduct) Meta() ProductMeta {
	return ProductMeta{
		Keywords: sp.product.BaseData().Keywords,
	}
}

// Loyalty of the product
func (sp BundleProduct) Loyalty() ProductLoyalty {
	return ProductLoyalty{
		Price:   sp.product.TeaserData().TeaserLoyaltyPriceInfo,
		Earning: sp.product.TeaserData().TeaserLoyaltyEarningInfo,
	}
}

// Attributes of the product
func (sp BundleProduct) Attributes() productDomain.Attributes {
	return sp.product.BaseData().Attributes
}

// Badges of the product
func (sp BundleProduct) Badges() ProductBadges {
	return ProductBadges{
		All: sp.product.BaseData().Badges,
	}
}

func mapWithActiveChoices(domainChoices []productDomain.Choice, activeChoices map[productDomain.Identifier]productDomain.ActiveChoice) []Choice {
	choices := mapChoices(domainChoices)

	if activeChoices == nil {
		return choices
	}

	for i, choice := range choices {
		activeChoice, ok := activeChoices[productDomain.Identifier(choice.Identifier)]
		if ok {
			choices[i].Active = NewGraphqlProductDto(activeChoice.Product, nil, nil)
		}
	}

	return choices
}

func mapChoices(domainChoices []productDomain.Choice) []Choice {
	choices := make([]Choice, 0, len(domainChoices))
	for _, domainChoice := range domainChoices {
		choices = append(choices, mapChoice(domainChoice))
	}

	return choices
}

func mapChoice(domainChoice productDomain.Choice) Choice {
	return Choice{
		Identifier: domainChoice.Identifier,
		Required:   domainChoice.Required,
		Label:      domainChoice.Label,
		Options:    mapOptions(domainChoice.Options),
	}
}

func mapOptions(domainOptions []productDomain.Option) []Option {
	options := make([]Option, 0, len(domainOptions))

	for _, domainOption := range domainOptions {
		options = append(options, mapOption(domainOption))
	}

	return options
}

func mapOption(domainOption productDomain.Option) Option {
	return Option{
		Product: NewGraphqlProductDto(domainOption.Product, nil, nil),
		Qty:     domainOption.MinQty,
	}
}
