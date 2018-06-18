package application

import cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"

type (

	// PaymentService
	PaymentService struct {
		DefaultPaymentMethod string `inject:"config:checkout.defaultPaymentMethod"`
	}
)

func (p PaymentService) GetDefaultCartPayment(cart *cartDomain.Cart) *cartDomain.CartPayment {
	payment := &cartDomain.CartPayment{}
	paymentInfo := cartDomain.PaymentInfo{
		Method: p.DefaultPaymentMethod,
	}
	payment.AddPayment(paymentInfo, cart.GetItemIds())
	return payment
}
