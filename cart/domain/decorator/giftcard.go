package decorator

import "flamingo.me/flamingo-commerce/v3/price/domain"

const (
	decoratedGiftCardError = "Unable to collect discounts, stopping and returning empty slice"
)

type (
	// DecoratedWithGiftCard interface for a decorated object to be able to handle giftcards
	// the difference to cart.WithGiftCard is, that these functions do NOT provide the client
	// with an error, errors are just logged
	DecoratedWithGiftCard interface {
		HasRemainingGiftCards() bool
		HasAppliedGiftCards() bool
		SumAppliedGiftCards() domain.Price
		SumGrandTotalWithGiftCards() domain.Price
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

// SumAppliedGiftCards sum up all applied amounts of giftcads
// price is returned as a payable
func (dc DecoratedCart) SumAppliedGiftCards() domain.Price {
	return dc.executeAndLog(dc.Cart.SumAppliedGiftCards)
}

// SumGrandTotalWithGiftCards calculate the grand total of the cart minus gift cards
func (dc DecoratedCart) SumGrandTotalWithGiftCards() domain.Price {
	return dc.executeAndLog(dc.Cart.SumGrandTotalWithGiftCards)
}

// executeAndLog executes given function and logs in case of an error
func (dc DecoratedCart) executeAndLog(toExecute func() (domain.Price, error)) domain.Price {
	result, err := toExecute()
	if err != nil {
		dc.Logger.Error(decoratedGiftCardError)
	}
	return result
}
