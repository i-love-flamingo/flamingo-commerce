package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
)

type (
	// StartPlaceOrderResult result of start place order
	StartPlaceOrderResult struct {
		UUID string
	}

	// PlaceOrderContext infos
	PlaceOrderContext struct {
		Cart       *dto.DecoratedCart
		OrderInfos *PlacedOrderInfos
		State      State
		UUID       string
	}

	// PlacedOrderInfos infos
	PlacedOrderInfos struct {
		PaymentInfos        []application.PlaceOrderPaymentInfo
		PlacedOrderInfos    []placeorder.PlacedOrderInfo
		Email               string
		PlacedDecoratedCart *dto.DecoratedCart
	}
)
