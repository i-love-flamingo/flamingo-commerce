package cart

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"

	price "flamingo.me/flamingo-commerce/v3/price/domain"
)

const (
	splitQualifierSeparator = "-"
)

type (
	// PaymentSelection value object - that represents the payment selection on the cart
	PaymentSelection interface {
		Gateway() string
		// ChargeSplits - the selected split per ChargeType and PaymentMethod
		CartSplit() PaymentSplit
		// ChargeSplits - the selected split per ChargeType and PaymentMethod
		ItemSplit() PaymentSplitByItem
		TotalValue() price.Price
		MethodByType(string) string
		IdempotencyKey() string
		GenerateNewIdempotencyKey() (PaymentSelection, error)
	}

	// SplitQualifier qualifies by Charge Type, Charge Reference and Payment Method
	SplitQualifier struct {
		ChargeType      string
		ChargeReference string
		Method          string
	}

	// PaymentSplit represents the Charges qualified by Charge Type, Charge Reference and Payment Method
	PaymentSplit map[SplitQualifier]price.Charge

	// PaymentSplitByItem - similar to value object that contains items of the different possible types, that have a price
	PaymentSplitByItem struct {
		CartItems     map[string]PaymentSplit
		ShippingItems map[string]PaymentSplit
		TotalItems    map[string]PaymentSplit
	}

	// PaymentSplitByItemBuilder - Builder to get valid PaymentSplitByItem instances
	PaymentSplitByItemBuilder struct {
		inBuilding *PaymentSplitByItem
	}

	// DefaultPaymentSelection value object - that implements the PaymentSelection interface
	DefaultPaymentSelection struct {
		// GatewayProp - the selected Gateway
		GatewayProp        string
		ChargedItemsProp   PaymentSplitByItem
		IdempotencyKeyUUID string
	}

	// PaymentSplitService enables the creation of a PaymentSplitByItem following different payment methods
	PaymentSplitService struct{}

	// builderAddFunc is a function used by the builder to add items
	// function which corresponds to builder addX function (addCartItem, addShipping, addTotal)
	builderAddFunc func(string, string, price.Charge) *PaymentSplitByItemBuilder

	// itemsWithAdd is a helper struct which holds items with their corresponding add function
	// from PaymentSplitByItemBuilder
	itemsWithAdd struct {
		// map of payable items corresponding to price.PricedItems
		ItemsToPay  map[string]price.Price
		AddFunction builderAddFunc
	}
)

var (
	_ PaymentSelection = new(DefaultPaymentSelection)

	// ErrSplitNoGiftCards indicates that there are no gift cards given to PaymentSplitWithGiftCards
	ErrSplitNoGiftCards = errors.New("no gift cards applied")

	// ErrSplitEmptyGiftCards indicates that there are gift cards given but with 0 applied balance
	ErrSplitEmptyGiftCards = errors.New("applied gift cards are empty")

	// ErrSplitGiftCardsExceedTotal indicates that gift card sum exceeds total of prices items
	ErrSplitGiftCardsExceedTotal = errors.New("gift card amount exceeds total priced items value")

	// ErrSplitGiftCardsNoChargeTypeMapping indicates that there is no mapping from the gift card charge type to an actual payment method
	ErrSplitGiftCardsNoChargeTypeMapping = fmt.Errorf("payment method for charge type %q not defined", price.ChargeTypeGiftCard)

	// ErrPaymentSelectionNotSet is used for nil PaymentSelection on cart
	ErrPaymentSelectionNotSet = errors.New("paymentSelection not set")
)

// NewDefaultPaymentSelection returns a PaymentSelection that can be used to update the cart
// is able to include gift card charges if applied to cart
func NewDefaultPaymentSelection(gateway string, chargeTypeToPaymentMethod map[string]string, cart Cart) (PaymentSelection, error) {
	pricedItems := cart.GetAllPaymentRequiredItems()
	giftCards := cart.AppliedGiftCards
	if _, ok := chargeTypeToPaymentMethod[price.ChargeTypeMain]; !ok {
		return nil, fmt.Errorf("payment method for charge type %q not defined", price.ChargeTypeMain)
	}
	result, err := newPaymentSelectionWithGiftCard(gateway, chargeTypeToPaymentMethod, pricedItems, giftCards)
	if err != nil {
		return result, err
	}
	// filter out zero charges from here on out
	result = RemoveZeroCharges(result, chargeTypeToPaymentMethod)
	// add an new Idempotency-Key to the payment selection
	return result.GenerateNewIdempotencyKey()
}

// RemoveZeroCharges removes charges which have an value of zero from selection as they are necessary
// for our internal calculations but not for external clients, we assume zero charges are ignored
// moreover charges are transformed to pay ables
func RemoveZeroCharges(selection PaymentSelection, chargeTypeToPaymentMethod map[string]string) PaymentSelection {
	// guard clause for nil selection
	if selection == nil {
		return nil
	}
	result := DefaultPaymentSelection{
		GatewayProp: selection.Gateway(),
	}
	builder := &PaymentSplitByItemBuilder{}
	// remove all zero charges from selection with helper function
	removeZeroChargesFromSplit(selection.ItemSplit().CartItems, chargeTypeToPaymentMethod, builder.AddCartItem)
	removeZeroChargesFromSplit(selection.ItemSplit().ShippingItems, chargeTypeToPaymentMethod, builder.AddShippingItem)
	removeZeroChargesFromSplit(selection.ItemSplit().TotalItems, chargeTypeToPaymentMethod, builder.AddTotalItem)

	result.ChargedItemsProp = builder.Build()

	resultWithIdempotencyKey, _ := result.GenerateNewIdempotencyKey()
	return resultWithIdempotencyKey
}

// removeZeroChargesFromSplit remove charges from single item splits
// helper which overwrites passed builder instance with adjusted charges
func removeZeroChargesFromSplit(
	paymentSplit map[string]PaymentSplit,
	chargeTypeToPaymentMethod map[string]string,
	add builderAddFunc,
) {
	for id, split := range paymentSplit {
		for qualifier, charge := range split.ChargesByType().GetAllCharges() {
			// charge should be transformed to payable
			charge = price.Charge{
				Price:     charge.Price.GetPayable(),
				Value:     charge.Value.GetPayable(),
				Type:      charge.Type,
				Reference: charge.Reference,
			}
			// skip charges with zero value
			if charge.Value.IsZero() {
				continue
			}
			// we assume that map of types and method matches
			method := chargeTypeToPaymentMethod[qualifier.Type]
			add(id, method, charge)
		}
	}
}

// newSimplePaymentSelection returns a PaymentSelection that can be used to update the cart.
// multiple charges by item are not used here: The complete grandtotal is selected to be paid in one charge with the given paymentgateway and paymentmethod
func newSimplePaymentSelection(gateway string, method string, pricedItems PricedItems) PaymentSelection {
	selection := DefaultPaymentSelection{
		GatewayProp: gateway,
	}
	builder := PaymentSplitByItemBuilder{}

	for k, itemPrice := range pricedItems.CartItems() {
		builder.AddCartItem(k, method, price.Charge{
			Price: itemPrice,
			Value: itemPrice,
			Type:  price.ChargeTypeMain,
		})

	}
	for k, itemPrice := range pricedItems.ShippingItems() {
		builder.AddShippingItem(k, method, price.Charge{
			Price: itemPrice,
			Value: itemPrice,
			Type:  price.ChargeTypeMain,
		})

	}
	for k, itemPrice := range pricedItems.TotalItems() {
		builder.AddTotalItem(k, method, price.Charge{
			Price: itemPrice,
			Value: itemPrice,
			Type:  price.ChargeTypeMain,
		})
	}
	selection.ChargedItemsProp = builder.Build()
	return selection
}

// newPaymentSelectionWithGiftCard returns Selection with given gift card charge type taken into account
func newPaymentSelectionWithGiftCard(gateway string, chargeTypeToPaymentMethod map[string]string, pricedItems PricedItems, appliedGiftCards []AppliedGiftCard) (PaymentSelection, error) {
	// create payment split by item with gift cards
	service := PaymentSplitService{}
	result, err := service.SplitWithGiftCards(chargeTypeToPaymentMethod, pricedItems, appliedGiftCards)
	// error handling
	if err != nil {
		switch err {
		case ErrSplitNoGiftCards:
			return newSimplePaymentSelection(gateway, chargeTypeToPaymentMethod[price.ChargeTypeMain], pricedItems), nil
		case ErrSplitEmptyGiftCards:
			return newSimplePaymentSelection(gateway, chargeTypeToPaymentMethod[price.ChargeTypeMain], pricedItems), nil
		case ErrSplitGiftCardsNoChargeTypeMapping:
			return newSimplePaymentSelection(gateway, chargeTypeToPaymentMethod[price.ChargeTypeMain], pricedItems), nil
		default:
			return nil, err
		}
	}
	// create selection
	selection := DefaultPaymentSelection{
		GatewayProp: gateway,
	}
	selection.ChargedItemsProp = *result
	return selection, nil
}

// NewPaymentSelection - with the passed PaymentSplitByItem
func NewPaymentSelection(gateway string, chargedItems PaymentSplitByItem) PaymentSelection {
	var selection PaymentSelection
	selection = DefaultPaymentSelection{
		GatewayProp:      gateway,
		ChargedItemsProp: chargedItems,
	}
	selection, _ = selection.GenerateNewIdempotencyKey()

	return selection
}

// Gateway returns the selected Gateway code
func (d DefaultPaymentSelection) Gateway() string {
	return d.GatewayProp
}

// CartSplit returns the selected split per ChargeType and PaymentMethod
func (d DefaultPaymentSelection) CartSplit() PaymentSplit {
	return d.ChargedItemsProp.Sum()
}

// ItemSplit returns the selected split per ChargeType and PaymentMethod
func (d DefaultPaymentSelection) ItemSplit() PaymentSplitByItem {
	return d.ChargedItemsProp
}

// TotalValue returns returns Valued price sum
func (d DefaultPaymentSelection) TotalValue() price.Price {
	return d.ChargedItemsProp.Sum().TotalValue()
}

// IdempotencyKey returns the Idempotency-Key for this payment selection
func (d DefaultPaymentSelection) IdempotencyKey() string {
	return d.IdempotencyKeyUUID
}

// GenerateNewIdempotencyKey updates the Idempotency-Key to a new value
func (d DefaultPaymentSelection) GenerateNewIdempotencyKey() (PaymentSelection, error) {
	key, err := uuid.NewRandom()
	if err != nil {
		return DefaultPaymentSelection{}, err
	}
	d.IdempotencyKeyUUID = key.String()
	return d, nil
}

// MarshalJSON adds the Idempotency-Key to the payment selection json
func (d DefaultPaymentSelection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GatewayProp      string             `json:"GatewayProp"`
		ChargedItemsProp PaymentSplitByItem `json:"ChargedItemsProp"`
		IdempotencyKey   string             `json:"IdempotencyKey"`
	}{
		GatewayProp:      d.GatewayProp,
		ChargedItemsProp: d.ChargedItemsProp,
		IdempotencyKey:   d.IdempotencyKey(),
	})
}

// Sum returns the resulting Split after sum all the included item split
func (c PaymentSplitByItem) Sum() PaymentSplit {
	sum := make(PaymentSplit)
	addToSum := func(splits PaymentSplit) {
		for qualifier, charge := range splits {
			_, ok := sum[qualifier]
			if ok {
				sum[qualifier], _ = sum[qualifier].Add(charge)
			} else {
				sum[qualifier] = charge
			}
		}
	}
	for _, itemSplit := range c.CartItems {
		addToSum(itemSplit)
	}
	for _, itemSplit := range c.ShippingItems {
		addToSum(itemSplit)
	}
	for _, itemSplit := range c.TotalItems {
		addToSum(itemSplit)
	}
	return sum
}

// TotalValue returns the sum of the valued Price in the included Charges in this Split
func (s PaymentSplit) TotalValue() price.Price {
	var prices []price.Price
	for _, v := range s {
		prices = append(prices, v.Value)
	}
	sum, _ := price.SumAll(prices...)
	return sum
}

// ChargesByType returns Charges (a list of Charges summed by Type)
func (s PaymentSplit) ChargesByType() price.Charges {
	charges := price.Charges{}
	for _, charge := range s {
		charges = charges.AddCharge(charge)
	}
	return charges
}

// MarshalJSON serialize to json
func (s PaymentSplit) MarshalJSON() ([]byte, error) {
	result := make(map[string]price.Charge)
	for qualifier, charge := range s {
		// explicit method and chargeType is necessary, otherwise keys could be overwritten
		if qualifier.Method == "" || qualifier.ChargeType == "" {
			return nil, errors.New("method or ChargeType is empty")
		}
		// SplitQualifier is parsed to a string method-chargeType-chargeReference
		result[qualifier.Method+splitQualifierSeparator+qualifier.ChargeType+splitQualifierSeparator+qualifier.ChargeReference] = charge
	}
	return json.Marshal(result)
}

// UnmarshalJSON deserialize from json
func (s *PaymentSplit) UnmarshalJSON(data []byte) error {
	var input map[string]price.Charge
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	result := PaymentSplit{}
	// parse string method-chargeType-chargeReference back to split qualifier
	for key, charge := range input {
		splitted := strings.Split(key, splitQualifierSeparator)
		// guard in case cannot be split
		if len(splitted) < 2 {
			return errors.New("SplitQualifier cannot be parsed for paymentsplit")
		}
		qualifier := SplitQualifier{
			Method:     splitted[0],
			ChargeType: splitted[1],
		}

		if len(splitted) == 3 {
			qualifier.ChargeReference = splitted[2]
		}

		result[qualifier] = charge
	}
	*s = result
	return nil
}

// AddCartItem adds a cart items charge to the PaymentSplitByItem
func (pb *PaymentSplitByItemBuilder) AddCartItem(id string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.CartItems[id] == nil {
		pb.inBuilding.CartItems[id] = make(PaymentSplit)
	}
	pb.inBuilding.CartItems[id][SplitQualifier{
		Method:          method,
		ChargeType:      charge.Type,
		ChargeReference: charge.Reference,
	}] = charge
	return pb
}

// AddShippingItem adds shipping charge
func (pb *PaymentSplitByItemBuilder) AddShippingItem(deliveryCode string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.ShippingItems[deliveryCode] == nil {
		pb.inBuilding.ShippingItems[deliveryCode] = make(PaymentSplit)
	}
	pb.inBuilding.ShippingItems[deliveryCode][SplitQualifier{
		Method:          method,
		ChargeType:      charge.Type,
		ChargeReference: charge.Reference,
	}] = charge
	return pb
}

// AddTotalItem adds total item charge
func (pb *PaymentSplitByItemBuilder) AddTotalItem(totalType string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.TotalItems[totalType] == nil {
		pb.inBuilding.TotalItems[totalType] = make(PaymentSplit)
	}
	pb.inBuilding.TotalItems[totalType][SplitQualifier{
		Method:          method,
		ChargeType:      charge.Type,
		ChargeReference: charge.Reference,
	}] = charge
	return pb
}

// Build returns the instance of PaymentSplitByItem
func (pb *PaymentSplitByItemBuilder) Build() PaymentSplitByItem {
	pb.init()
	return *pb.inBuilding
}

func (pb *PaymentSplitByItemBuilder) init() {
	if pb.inBuilding != nil {
		return
	}
	pb.inBuilding = &PaymentSplitByItem{
		CartItems:     make(map[string]PaymentSplit),
		ShippingItems: make(map[string]PaymentSplit),
		TotalItems:    make(map[string]PaymentSplit),
	}
}

// SplitWithGiftCards calculates a payment selection based on given method, priced items and applied gift cards
func (service PaymentSplitService) SplitWithGiftCards(chargeTypeToPaymentMethod map[string]string, items PricedItems, cards AppliedGiftCards) (*PaymentSplitByItem, error) {
	totalValue := items.Sum()
	// guard clause, if no gift cards no payment split with gift cards
	if len(cards) == 0 {
		return nil, ErrSplitNoGiftCards
	}
	// guard, gift card method is not defined
	if _, ok := chargeTypeToPaymentMethod[price.ChargeTypeGiftCard]; !ok {
		return nil, ErrSplitGiftCardsNoChargeTypeMapping
	}
	allGcAmounts := make([]price.Price, 0, len(cards))
	for _, gc := range cards {
		allGcAmounts = append(allGcAmounts, gc.Applied)
	}
	totalGCValue, err := price.SumAll(allGcAmounts...)
	if err != nil {
		return nil, err
	}
	// guard clause, all gift cards are empty
	if totalGCValue.IsZero() {
		return nil, ErrSplitEmptyGiftCards
	}
	// guard clause, can't split because gift card total exceeds payable amount of items
	if totalGCValue.GetPayable().IsGreaterThen(totalValue.GetPayable()) {
		return nil, ErrSplitGiftCardsExceedTotal
	}

	builder := &PaymentSplitByItemBuilder{}
	helpers := service.initItemsWithAdd(items, builder)
	// slices are passed by reference, avoid side effects on cart
	copiedCards := make(AppliedGiftCards, len(cards))
	copy(copiedCards, cards)
	// loop over helper containing the items to pay
	// and their corresponding helper function
	for _, helper := range helpers {
		itemKeys := service.sortItemsToPayKeys(helper.ItemsToPay)
		// distribute gift cards across items, this tries to spend the full gift card per item
		for i, card := range copiedCards {
			for _, k := range itemKeys {
				itemPrice := helper.ItemsToPay[k]

				// nothing to pay with gift card
				if itemPrice.IsZero() {
					continue
				}

				// burn gift card amount on item price
				remainingItem, appliedGiftCard, err := service.clearGiftCardWithItem(&card, &itemPrice)
				if err != nil {
					return nil, err
				}
				itemPrice = remainingItem

				// add calculated charges to builder for payment selection
				builder = helper.AddFunction(k, chargeTypeToPaymentMethod[price.ChargeTypeMain], price.Charge{
					Price: remainingItem,
					Value: remainingItem,
					Type:  price.ChargeTypeMain,
				})

				if !appliedGiftCard.IsZero() {
					builder = helper.AddFunction(k, chargeTypeToPaymentMethod[price.ChargeTypeGiftCard], price.Charge{
						Price:     appliedGiftCard,
						Value:     appliedGiftCard,
						Type:      price.ChargeTypeGiftCard,
						Reference: card.Code,
					})
				}

				copiedCards[i] = card
				helper.ItemsToPay[k] = itemPrice
			}
		}
	}

	result := builder.Build()
	return &result, nil
}

// clearGiftCardWithItem try to apply complete gift card on item
// otherwise rest will still be available to spend on applied
func (service PaymentSplitService) clearGiftCardWithItem(card *AppliedGiftCard, itemPrice *price.Price) (remaining,
	applied price.Price, err error) {
	// gift card is less than item price
	toApply := card.Applied
	if card.Applied.IsGreaterThen(*itemPrice) || card.Applied.Equal(*itemPrice) {
		// gift card is greater or equal item price
		toApply = *itemPrice
	}
	remaining, err = itemPrice.Sub(toApply)
	if err != nil {
		return remaining, applied, err
	}
	applied = toApply

	card.Applied, err = card.Applied.Sub(toApply)
	if err != nil {
		return remaining, applied, err
	}
	return remaining, applied, nil
}

// initItemsWithAdd init helper struct containing priced item entry with corresponding builder method
func (service PaymentSplitService) initItemsWithAdd(items PricedItems, builder *PaymentSplitByItemBuilder) []itemsWithAdd {
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

func (service PaymentSplitService) sortItemsToPayKeys(itemsToPay map[string]price.Price) []string {
	// sort item keys ascending, to stabilise later item access
	var itemKeys []string
	for k := range itemsToPay {
		itemKeys = append(itemKeys, k)
	}
	sort.Strings(itemKeys)
	return itemKeys
}

// MethodByType returns the payment method by charge type
func (d DefaultPaymentSelection) MethodByType(chargeType string) string {
	for qualifier := range d.CartSplit() {
		if qualifier.ChargeType == chargeType {
			return qualifier.Method
		}
	}

	return ""
}
