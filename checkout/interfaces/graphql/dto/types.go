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

	//State state iface
	State interface {
		Final() bool
	}

	//StateWait concrete state
	StateWait struct {
	}
	//StateSuccess concrete state
	StateSuccess struct {
	}
	//StateFatalError concrete state
	StateFatalError struct {
		//Error info
		Error string
	}
	//StateShowIframe concrete state
	StateShowIframe struct {
		//URL the url
		URL string
	}
	//StateShowHTML concrete state
	StateShowHTML struct {
		//HTML the HTML
		HTML string
	}
	//StateRedirect concrete state
	StateRedirect struct {
		//URL the url
		URL string
	}
	//StateCancelled concrete state
	StateCancelled struct {
		CancellationReason CancellationReason
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

//Final if final
func (s *StateWait) Final() bool {
	return false
}

//Final if final
func (s *StateSuccess) Final() bool {
	return true
}

//Final if final
func (s *StateFatalError) Final() bool {
	return true
}

//Final if final
func (s *StateShowIframe) Final() bool {
	return false
}

//Final if final
func (s *StateShowHTML) Final() bool {
	return false
}

//Final if final
func (s *StateRedirect) Final() bool {
	return false
}

//Final if final
func (s *StateCancelled) Final() bool {
	return true
}

//Reason returns reason
func (c *CancellationReasonPaymentError) Reason() string {
	return c.PaymentError.Error()
}

//Reason returns reason
func (c *CancellationReasonValidationError) Reason() string {
	return "cart-invalid"
}
