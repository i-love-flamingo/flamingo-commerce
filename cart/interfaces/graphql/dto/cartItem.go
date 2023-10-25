package dto

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	BundleChoiceConfiguration struct {
		Identifier             string
		MarketplaceCode        string
		VariantMarketplaceCode string
		Qty                    int
	}

	CartItemResolver struct{}
)

func (c CartItemResolver) BundleConfiguration(ctx context.Context, item *cart.Item) ([]*BundleChoiceConfiguration, error) {
	if item.BundleConfig == nil {
		return nil, nil
	}

	choices := make([]*BundleChoiceConfiguration, 0)
	for identifier, configuration := range item.BundleConfig {
		choices = append(choices, &BundleChoiceConfiguration{
			Identifier:             string(identifier),
			MarketplaceCode:        configuration.MarketplaceCode,
			VariantMarketplaceCode: configuration.MarketplaceCode,
			Qty:                    configuration.Qty,
		})
	}

	return choices, nil
}
