package infrastructure

import (
	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type CartItemMatcher interface {
	Matches(item domaincart.Item, addRequest domaincart.AddRequest) bool
}
type DefaultCartItemMatcher struct{}

var _ CartItemMatcher = (*DefaultCartItemMatcher)(nil)

const passengerIDKey = "passenger_id"

// Matches implements CartItemMatcher.
func (DefaultCartItemMatcher) Matches(item domaincart.Item, addRequest domaincart.AddRequest) bool {
	if item.MarketplaceCode != addRequest.MarketplaceCode {
		return false
	}

	if item.VariantMarketPlaceCode != addRequest.VariantMarketplaceCode {
		return false
	}

	if !item.BundleConfig.Equals(addRequest.BundleConfiguration) {
		return false
	}

	return true
}
