package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
)

type (
	PlaceOrderContext struct {
		Cart       *decorator.DecoratedCart
		OrderInfos *PlacedOrderInfos
		State      State
	}

	PlacedOrderInfos struct {
		PaymentInfos     []application.PlaceOrderPaymentInfo
		PlacedOrderInfos []placeorder.PlacedOrderInfo
		Email            string
		PlacedDecoratedCart *dto.DecoratedCart
	}

	State interface {
		Final() bool
	}

	StateWait struct {
	}

	StateSuccess struct {
	}

	StateFatalError struct {
		Error string
	}

	StateShowIframe struct {
		Url string
	}

	StateShowHtml struct {
		Html string
	}

	StateRedirect struct {
		Url string
	}

	StateCancelled struct {
		CancellationReason CancellationReason
	}

	CancellationReason interface {
		Reason() string
	}

	CancellationReasonPaymentError struct {
		PaymentError error
	}

	CancellationReasonValidationError struct {
		ValidationResult validation.Result
	}
)

func (s *StateWait) Final() bool {
	return false
}

func (s *StateSuccess) Final() bool {
	return true
}

func (s *StateFatalError) Final() bool {
	return true
}

func (s *StateShowIframe) Final() bool {
	return false
}

func (s *StateShowHtml) Final() bool {
	return false
}

func (s *StateRedirect) Final() bool {
	return false
}

func (s *StateCancelled) Final() bool {
	return true
}

func (c *CancellationReasonPaymentError) Reason() string {
	return c.PaymentError.Error()
}

func (c *CancellationReasonValidationError) Reason() string {
	return "cart-invalid"
}
