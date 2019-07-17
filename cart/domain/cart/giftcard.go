package cart

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// AppliedGiftCard value object represents a gift card (partial payment) on the cart
	AppliedGiftCard struct {
		Code      string
		Applied   domain.Price // how much of the gift card has been subtracted from cart price
		Remaining domain.Price // how much of the gift card is still available
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

	// PaymentSplitWithGiftCards handles payment split on item level including gift cards
	PaymentSplitWithGiftCards struct {
		Method      string
		PricedItems PricedItems
		GiftCards   AppliedGiftCards
	}

	// ErrSplitNoGiftCards indicates that there are no gift cards given to PaymentSplitWithGiftCards
	ErrSplitNoGiftCards struct {
		message string
	}

	// ErrSplitEmptyGiftCards indicates that there are gift cards given but with 0 applied balance
	ErrSplitEmptyGiftCards struct {
		message string
	}

	// ErrSplitGiftCardsExceedTotal indicates that gift card sum exceeds total of prices items
	ErrSplitGiftCardsExceedTotal struct {
		message string
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
		return domain.Price{}.GetPayable(), nil
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

// Split calculates a payment selection based on given method, priced items and applied gift cards
func (split PaymentSplitWithGiftCards) Split() (*PaymentSplitByItem, error) {
	totalValue := split.PricedItems.Sum()
	// guard clause, if no gift cards no payment split with gift cards
	if len(split.GiftCards) == 0 {
		return nil, &ErrSplitNoGiftCards{message: "No gift cards applied"}
	}
	allGcAmounts := make([]price.Price, 0, len(split.GiftCards))
	for _, gc := range split.GiftCards {
		allGcAmounts = append(allGcAmounts, gc.Applied)
	}
	totalGCValue, err := price.SumAll(allGcAmounts...)
	if err != nil {
		return nil, err
	}
	// guard clause, all gift cards are empty
	if totalGCValue.IsZero() {
		return nil, &ErrSplitEmptyGiftCards{message: "Applied gift cards are empty"}
	}
	// guard clause, can't split because gift card total exceeds payable amount of items
	if totalGCValue.IsGreaterThen(totalValue) {
		return nil, &ErrSplitGiftCardsExceedTotal{"gift card amount exceeds total priced items value"}
	}
	// distribute gift card amounts relatively across all items
	giftCartAmountRatio := totalGCValue.FloatAmount() / totalValue.FloatAmount()
	builder := &PaymentSplitByItemBuilder{}
	helpers := split.initItemsWithAdd(split.PricedItems, builder)
	// loop over helper containing the items to pay
	// and their corresponding helper function
	for _, helper := range helpers {
		builder, totalGCValue, err = split.splitWithGiftCards(builder, helper, giftCartAmountRatio, totalGCValue)
		if err != nil {
			return nil, err
		}
	}
	result := builder.Build()
	return &result, nil
}

// Split calculates a payment selection based on given method, priced items and applied gift cards
func (split PaymentSplitWithGiftCards) initItemsWithAdd(items PricedItems, builder *PaymentSplitByItemBuilder) []itemsWithAdd {
	return []itemsWithAdd{
		// cart items
		{
			ItemsToPay:  items.CartItems(),
			AddFunction: builder.AddCartItem,
		},
		// shipping
		{
			ItemsToPay:  items.ShippingItems(),
			AddFunction: builder.AddShippingItem,
		},
		// total
		{
			ItemsToPay:  items.TotalItems(),
			AddFunction: builder.AddTotalItem,
		},
	}
}

// splitWithGiftCards distribute gift card charges across item prices
func (split PaymentSplitWithGiftCards) splitWithGiftCards(builder *PaymentSplitByItemBuilder, helper itemsWithAdd,
	ratio float64, totalGCValue price.Price) (*PaymentSplitByItemBuilder, price.Price, error) {
	var remainingItemValue, appliedGcAmount price.Price
	var err error
	// loop over helper containing the items to pay
	for k, itemPrice := range helper.ItemsToPay {
		remainingItemValue, totalGCValue, appliedGcAmount, err = split.calcRelativeGiftCardAmount(itemPrice, totalGCValue, ratio)
		if err != nil {
			return nil, totalGCValue, err
		}
		// only add values if there are not zero
		if !remainingItemValue.IsZero() {
			builder = helper.AddFunction(k, split.Method, price.Charge{
				Price: remainingItemValue,
				Value: remainingItemValue,
				Type:  price.ChargeTypeMain,
			})
		}
		if !appliedGcAmount.IsZero() {
			builder = helper.AddFunction(k, split.Method, price.Charge{
				Price: appliedGcAmount,
				Value: appliedGcAmount,
				Type:  price.ChargeTypeGiftCard,
			})
		}
	}
	return builder, totalGCValue, nil
}

// calcRelativeGiftCardAmount calc amount of applied gift card relative to item price
func (split PaymentSplitWithGiftCards) calcRelativeGiftCardAmount(value price.Price, remainingGcAmount price.Price,
	ratio float64) (remainingItemValue price.Price,
	newRemainingGcAmount price.Price, appliedGcAmount price.Price, err error) {
	//relativeItemGcAmount the gift card amount that relates to the given item Value
	relativeItemGcAmount := price.NewFromFloat(ratio*value.FloatAmount(), value.Currency()).GetPayable()
	// if the relative amount is greater than the complete the remaining, just remove the remaining
	if relativeItemGcAmount.IsGreaterThen(remainingGcAmount) {
		relativeItemGcAmount = remainingGcAmount
	}
	// if the relative amount is greater than the item price, just use the item price
	if relativeItemGcAmount.IsGreaterThen(value) {
		relativeItemGcAmount = value
	}
	appliedGcAmount = relativeItemGcAmount
	newRemainingGcAmount, err = remainingGcAmount.Sub(appliedGcAmount)
	if err != nil {
		return
	}
	remainingItemValue, err = value.Sub(appliedGcAmount)
	if err != nil {
		return
	}
	return
}

func (e *ErrSplitNoGiftCards) Error() string {
	return e.message
}

func (e *ErrSplitEmptyGiftCards) Error() string {
	return e.message
}

func (e *ErrSplitGiftCardsExceedTotal) Error() string {
	return e.message
}
