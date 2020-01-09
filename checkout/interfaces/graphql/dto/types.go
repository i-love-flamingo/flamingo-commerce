package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//PlaceOrderContext infos
	PlaceOrderContext struct {
		Cart       *decorator.DecoratedCart
		OrderInfos *PlacedOrderInfos
		State      process.State
	}

	//PlacedOrderInfos infos
	PlacedOrderInfos struct {
		PaymentInfos     []application.PlaceOrderPaymentInfo
		PlacedOrderInfos []placeorder.PlacedOrderInfo
		Email            string
		PlacedDecoratedCart *dto.DecoratedCart
	}

	//CancellationReason iface
	CancellationReason interface {
		Reason() string
	}

	//CancellationReasonPaymentError error
	CancellationReasonPaymentError struct {
		PaymentError error
	}

	//CancellationReasonValidationError error
	CancellationReasonValidationError struct {
		ValidationResult validation.Result
	}
)

//Reason returns reason
func (c *CancellationReasonPaymentError) Reason() string {
	return c.PaymentError.Error()
}

//Reason returns reason
func (c *CancellationReasonValidationError) Reason() string {
	return "cart-invalid"
}
