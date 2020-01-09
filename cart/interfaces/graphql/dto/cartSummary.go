package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// CartSummary – provides custom graphql interface methods
	CartSummary struct {
		cart *cart.Cart
	}
)

// Discounts collects up discounts of cart based on its deliveries
// All discounts with the same campaign code are aggregated and returned as one with a summed price
func (cs *CartSummary) Discounts() *cart.AppliedDiscounts {
	result, err := cs.cart.MergeDiscounts()

	if err != nil {
		return nil
	}

	return &result
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (cs *CartSummary) HasAppliedDiscounts() bool {
	result, _ := cs.cart.HasAppliedDiscounts()
	return result
}

// SumTotalDiscountWithGiftCardsAmount – returns sum price of total discounts with applied gift cards
func (cs *CartSummary) SumTotalDiscountWithGiftCardsAmount() domain.Price {
	totalDiscount := cs.cart.SumTotalDiscountAmount()
	appliedGiftCardsAmount, _ := cs.cart.SumAppliedGiftCards()

	price, _ := totalDiscount.Sub(appliedGiftCardsAmount)
	return price
}

// SumAppliedDiscounts – returns the sum of the applied values of the AppliedDiscounts
func (cs CartSummary) SumAppliedDiscounts() *domain.Price {
	result, err := cs.cart.MergeDiscounts()
	if err != nil {
		return nil
	}

	sum, err := result.Sum()
	if err != nil {
		return nil
	}

	return &sum
}

// SumAppliedGiftCards – sums applied gift cards
func (cs CartSummary) SumAppliedGiftCards() *domain.Price {
	sum, err := cs.cart.SumAppliedGiftCards()
	if err != nil {
		return nil
	}
	return &sum
}

// SumGrandTotalWithGiftCards – sums grand total with gift cards
func (cs CartSummary) SumGrandTotalWithGiftCards() *domain.Price {
	sum, err := cs.cart.SumGrandTotalWithGiftCards()
	if err != nil {
		return nil
	}
	return &sum
}

// SumTaxes – sums taxes
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
