package cart

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// AppliedGiftCard value object represents a gift card (partial payment) on the cart
	AppliedGiftCard struct {
		Code             string
		Applied          domain.Price           // how much of the gift card has been subtracted from cart price
		Remaining        domain.Price           // how much of the gift card is still available
		CustomAttributes map[string]interface{} // additional custom attributes
	}

	// AppliedGiftCards convenience wrapper for array of applied gift cards
	AppliedGiftCards []AppliedGiftCard

	// WithGiftCard interface for a cart that is able to handle gift cards
	WithGiftCard interface {
		HasRemainingGiftCards() bool
		HasAppliedGiftCards() bool
		SumAppliedGiftCards() (domain.Price, error)
		SumGrandTotalWithGiftCards() (domain.Price, error)
	}
)

var (
	// interface assertion
	_ WithGiftCard = &Cart{}
)

// HasRemaining checks whether gift card has a remaining balance
func (card AppliedGiftCard) HasRemaining() bool {
	return !card.Remaining.IsZero()
}

// Total returns the total value of the gift card by adding what is applied and remaining
// In case the values cannot be added the function returns the remaining amount of the giftcard and an error
func (card AppliedGiftCard) Total() (domain.Price, error) {
	total, err := card.Applied.Add(card.Remaining)
	if err != nil {
		return card.Remaining, err
	}
	return total, nil
}

// HasAppliedGiftCards checks if a gift card is applied to the cart
func (c Cart) HasAppliedGiftCards() bool {
	return len(c.AppliedGiftCards) > 0
}

// SumAppliedGiftCards sum up all applied amounts of giftcads
// price is returned as a payable
func (c Cart) SumAppliedGiftCards() (domain.Price, error) {
	// guard for no gift cards applied
	if len(c.AppliedGiftCards) == 0 {
		return domain.Price{}.GetPayable(), nil
	}
	prices := make([]domain.Price, 0, len(c.AppliedGiftCards))
	// add prices to array
	for _, card := range c.AppliedGiftCards {
		prices = append(prices, card.Applied)
	}
	price, err := domain.SumAll(prices...)
	// in case of error regarding sum, pass on error
	if err != nil {
		return domain.Price{}.GetPayable(), err
	}
	return price.GetPayable(), nil
}

// SumGrandTotalWithGiftCards calculate the grand total of the cart minus gift cards
func (c Cart) SumGrandTotalWithGiftCards() (domain.Price, error) {
	giftCardTotal, err := c.SumAppliedGiftCards()
	if err != nil {
		return domain.Price{}.GetPayable(), err
	}
	// if there are no gift cards just return cart grand total
	total := c.GrandTotal()
	if giftCardTotal.IsZero() {
		return total.GetPayable(), err
	}
	// subtract gift card total from total for "remaining total"
	result, err := total.Sub(giftCardTotal)
	if err != nil {
		return domain.Price{}.GetPayable(), err
	}
	return result.GetPayable(), nil
}

// HasRemainingGiftCards check whether there are gift cards with remaining balance
func (c Cart) HasRemainingGiftCards() bool {
	for _, card := range c.AppliedGiftCards {
		if card.HasRemaining() {
			return true
		}
	}
	return false
}

// ByRemaining fetches gift cards that still have a remaining value from applied gift cards
func (cards *AppliedGiftCards) ByRemaining() AppliedGiftCards {
	result := AppliedGiftCards{}
	for _, card := range *cards {
		if card.HasRemaining() {
			result = append(result, card)
		}
	}
	return result
}

// GiftCardByCode returns a single gift card if the given code matches its code.
// First return parameter is the gift card if found and the second return parameter is a boolean depicting if a gift card was found
func (cards *AppliedGiftCards) GiftCardByCode(code string) (*AppliedGiftCard, bool) {
	for _, card := range *cards {
		if card.Code == code {
			return &card, true
		}
	}
	return nil, false
}
