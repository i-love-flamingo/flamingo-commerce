package cart

import "flamingo.me/flamingo-commerce/v3/price/domain"

type (
	// AppliedGiftCard value object represents a gift card (partial payment) on the cart
	AppliedGiftCard struct {
		Code      string
		Applied   domain.Price // how much of the gift card has been subtracted from cart price
		Remaining domain.Price // how much of the gift card is still available
	}

	// AppliedGiftCards
	AppliedGiftCards []AppliedGiftCard

	// WithGiftCard interface for a cart that is able to handle gift cards
	WithGiftCard interface {
		HasAppliedGiftCards() bool
		AppliedGiftCards() AppliedGiftCard
		SumAppliedGiftCards() (domain.Price, error)
		SumGrandTotalWithGiftCards() (domain.Price, error)
	}
)

var (
	// interface assertion
	_ WithGiftCard = &Cart{}
)

// HasAppliedGiftCards checks if a gift card is applied to the cart
func (c Cart) HasAppliedGiftCards() bool {
	return len(c.AppliedGiftCards()) > 0
}

func (c Cart) AppliedGiftCards() AppliedGiftCard {
	panic("implement me")
}

func (c Cart) SumAppliedGiftCards() (domain.Price, error) {
	panic("implement me")
}

func (c Cart) SumGrandTotalWithGiftCards() (domain.Price, error) {
	panic("implement me")
}
