package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
)

type (

	//PlaceOrderResult - represents the result
	PlaceOrderResult struct {
		Status               Status
		CartValidationResult *validation.Result
		OrderSuccessData     *OrderSuccessData
		Error                *Error
	}

	//Error value object
	Error struct {
		IsPaymentError bool
		ErrorKey       string
	}

	//Status value
	Status string

	// OrderSuccessData represents the infos available if the order was placed successfully
	OrderSuccessData struct {
		PaymentInfos        []application.PlaceOrderPaymentInfo
		PlacedOrderInfos    []placeorder.PlacedOrderInfo
		Email               string
		PlacedDecoratedCart decorator.DecoratedCart
	}
)

//allowed Values for Status Enum
const (
	INVALID        = "INVALID"
	ERROR          = "ERROR"
	ORDERSUCCESS   = "ORDERSUCCESS"
	PAYMENTPENDING = "PAYMENTPENDING"
)
