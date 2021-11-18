package decorator

import "flamingo.me/flamingo-commerce/v3/price/domain"

type (
	// DecoratedWithGiftCard interface for a decorated object to be able to handle giftcards
	// the difference to cart.WithGiftCard is, that these functions do NOT provide the client
	// with an error, errors are just logged
	DecoratedWithGiftCard interface {
		HasRemainingGiftCards() bool
		HasAppliedGiftCards() bool
		TotalGiftCardAmount() domain.Price
		GrandTotalWithGiftCards() domain.Price
	}
)

var (
	// interface assertion
	_ DecoratedWithGiftCard = &DecoratedCart{}
)

// HasRemainingGiftCards check whether there are gift cards with remaining balance
func (dc DecoratedCart) HasRemainingGiftCards() bool {
	return dc.Cart.HasRemainingGiftCards()
}

// HasAppliedGiftCards checks if a gift card is applied to the cart
func (dc DecoratedCart) HasAppliedGiftCards() bool {
	return dc.Cart.HasAppliedGiftCards()
}

// TotalGiftCardAmount sum up all applied amounts of giftcads
// price is returned as a payable
func (dc DecoratedCart) TotalGiftCardAmount() domain.Price {
	return dc.Cart.TotalGiftCardAmount
}

// GrandTotalWithGiftCards calculate the grand total of the cart minus gift cards
func (dc DecoratedCart) GrandTotalWithGiftCards() domain.Price {
	return dc.Cart.GrandTotalWithGiftCards
}
