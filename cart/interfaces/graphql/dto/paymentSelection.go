package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

// PaymentSelectionSplit is a GraphQL specific representation of `cart.PaymentSplit`
type PaymentSelectionSplit struct {
	Qualifier cart.SplitQualifier
	Charge    domain.Charge
}
