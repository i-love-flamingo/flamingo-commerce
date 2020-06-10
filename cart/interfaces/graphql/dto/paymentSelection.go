package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
)

// PaymentSelection struct for extending with CartSplit
type PaymentSelection struct {
	paymentSelection *cart.PaymentSelection
}

// CartSplit to extend PaymentSelection
func (ps *PaymentSelection) CartSplit() []cart.PaymentSplit {
	return []cart.PaymentSplit{}
}

// PaymentSelectionSplitQualifier struct
type PaymentSelectionSplitQualifier struct {
	Type      string
	Method    string
	Reference string
}

// DefaultPaymentSelection struct for extending with CartSplit
type DefaultPaymentSelection struct {
	defaultPaymentSelection *cart.DefaultPaymentSelection
}

// DefaultPaymentSelection method
func (dp *DefaultPaymentSelection) DefaultPaymentSelection() *cart.DefaultPaymentSelection {
	return dp.defaultPaymentSelection
}

// GenerateNewIdempotencyKey returns a new key
func (dp *DefaultPaymentSelection) GenerateNewIdempotencyKey() (cart.PaymentSelection, error) {
	return dp.defaultPaymentSelection, nil
}

// IdempotencyKey returns a new key
func (dp *DefaultPaymentSelection) IdempotencyKey() string {
	return dp.defaultPaymentSelection.IdempotencyKeyUUID
}

// MethodByType returns a new key
func (dp *DefaultPaymentSelection) MethodByType(chargeType string) string {
	return dp.defaultPaymentSelection.MethodByType(chargeType)
}

// ItemSplit returns a new key
func (dp *DefaultPaymentSelection) ItemSplit() cart.PaymentSplitByItem {
	return dp.defaultPaymentSelection.ItemSplit()
}

// Gateway returns the gateway from defaultPaymentSelection
func (dp *DefaultPaymentSelection) Gateway() string {
	return dp.defaultPaymentSelection.Gateway()
}

// TotalValue returns the total value from defaultPaymentSelection
func (dp *DefaultPaymentSelection) TotalValue() price.Price {
	return dp.defaultPaymentSelection.TotalValue()
}

// CartSplit to extend DefaultPaymentSelection
func (dp *DefaultPaymentSelection) CartSplit() cart.PaymentSplit {
	return cart.PaymentSplit{}
}
