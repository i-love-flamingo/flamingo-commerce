package dto

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

// CartAppliedDiscounts DTO for cart.AppliedDiscounts
type CartAppliedDiscounts struct {
	discounts cart.AppliedDiscounts
}

// Items getter
func (d *CartAppliedDiscounts) Items() []cart.AppliedDiscount {
	return d.discounts.Items()
}

// ByCampaignCode getter and wrapper
func (d *CartAppliedDiscounts) ByCampaignCode(campaignCode string) *CartAppliedDiscounts {
	return &CartAppliedDiscounts{discounts: d.discounts.ByCampaignCode(campaignCode)}
}

// ByType getter and wrapper
func (d *CartAppliedDiscounts) ByType(filterType string) *CartAppliedDiscounts {
	return &CartAppliedDiscounts{discounts: d.discounts.ByType(filterType)}
}

// CartAppliedDiscountsResolver resolves discounts for items
type CartAppliedDiscountsResolver struct{}

// ForItem resolves for cart Items
func (*CartAppliedDiscountsResolver) ForItem(ctx context.Context, item *cart.Item) (*CartAppliedDiscounts, error) {
	return &CartAppliedDiscounts{discounts: item.AppliedDiscounts}, nil
}

// ForShippingItem resolves for shipping Items
func (*CartAppliedDiscountsResolver) ForShippingItem(ctx context.Context, item *cart.ShippingItem) (*CartAppliedDiscounts, error) {
	return &CartAppliedDiscounts{discounts: item.AppliedDiscounts}, nil
}
