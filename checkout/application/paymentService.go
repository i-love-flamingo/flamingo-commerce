package application

import cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

// PaymentService helper to get the default cart payment
type PaymentService struct {
	DefaultPaymentMethod string `inject:"config:checkout.defaultPaymentMethod"`
}

// GetDefaultCartPayment returns the payment to be used for a cart
func (p PaymentService) GetDefaultCartPayment(cart *cartDomain.Cart) *cartDomain.Payment {
	payment := &cartDomain.Payment{}
	paymentInfo := cartDomain.PaymentInfo{
		Method: p.DefaultPaymentMethod,
	}
	payment.AddPayment(paymentInfo, cart.GetItemCartReferences())
	return payment
}
