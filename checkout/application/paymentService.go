package application

import "go.aoe.com/flamingo/core/cart/domain/cart"

type (

	// PaymentService
	PaymentService struct {
		DefaultPaymentMethod string `inject:"config:checkout.defaultPaymentMethod"`
	}
)

func (p PaymentService) GetPayment() *cart.Payment {
	return &cart.Payment{
		Method: p.DefaultPaymentMethod,
	}
}
