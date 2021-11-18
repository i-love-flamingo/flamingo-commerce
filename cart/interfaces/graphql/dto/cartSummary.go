package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// CartSummary provides custom graphql interface methods
	CartSummary struct {
		cart *cart.Cart
	}
)

// Discounts collects up discounts of cart based on its deliveries
// All discounts with the same campaign code are aggregated and returned as one with a summed price
func (cs *CartSummary) Discounts() *CartAppliedDiscounts {
	result, err := cs.cart.MergeDiscounts()

	if err != nil {
		return nil
	}

	return &CartAppliedDiscounts{discounts: result}
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (cs *CartSummary) HasAppliedDiscounts() bool {
	result, _ := cs.cart.HasAppliedDiscounts()
	return result
}

// SumTotalDiscountWithGiftCardsAmount returns sum price of total discounts with applied gift cards
func (cs *CartSummary) SumTotalDiscountWithGiftCardsAmount() domain.Price {
	totalDiscount := cs.cart.TotalDiscountAmount
	appliedGiftCardsAmount := cs.cart.TotalGiftCardAmount

	price, _ := totalDiscount.Sub(appliedGiftCardsAmount)
	return price
}

// TotalDiscountAmount returns the sum of the applied values of the AppliedDiscounts
func (cs CartSummary) TotalDiscountAmount() *domain.Price {
	sum := cs.cart.TotalDiscountAmount

	return &sum
}

// TotalGiftCardAmount sums applied gift cards
func (cs CartSummary) TotalGiftCardAmount() *domain.Price {
	return &cs.cart.TotalGiftCardAmount
}

// GrandTotalWithGiftCards sums grand total with gift cards
func (cs CartSummary) GrandTotalWithGiftCards() *domain.Price {
	return &cs.cart.GrandTotalWithGiftCards
}

// SumTaxes sums taxes
func (cs CartSummary) SumTaxes() *Taxes {
	items := cs.cart.SumTaxes()
	taxes := make([]cart.Tax, 0, len(items))
	for _, tax := range items {
		taxes = append(taxes, tax)
	}

	if len(taxes) > 0 {
		return nil
	}

	return &Taxes{Items: taxes}
}

// SumPaymentSelectionCartSplitValueAmountByMethods sum
func (cs CartSummary) SumPaymentSelectionCartSplitValueAmountByMethods(methods []string) *domain.Price {
	if cs.cart.PaymentSelection == nil {
		return nil
	}

	prices := make([]domain.Price, 0, len(methods))

	for qualifier, charge := range cs.cart.PaymentSelection.CartSplit() {
		found := contains(methods, qualifier.Method)
		if !found {
			continue
		}

		prices = append(prices, charge.Value)
	}

	if len(prices) == 0 {
		return nil
	}

	price, err := domain.SumAll(prices...)
	if err != nil {
		return nil
	}

	return &price
}

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
