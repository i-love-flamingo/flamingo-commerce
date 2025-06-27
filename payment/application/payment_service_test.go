package application_test

import (
	"testing"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces/mocks"
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"github.com/stretchr/testify/assert"
)

func TestPaymentService_AvailablePaymentGateways(t *testing.T) {
	ps := application.PaymentService{}
	ps.Inject(func() map[string]interfaces.WebCartPaymentGateway {
		return map[string]interfaces.WebCartPaymentGateway{
			"gateway-code": &mocks.WebCartPaymentGateway{},
		}
	})

	assert.Equal(t, map[string]interfaces.WebCartPaymentGateway{
		"gateway-code": &mocks.WebCartPaymentGateway{},
	}, ps.AvailablePaymentGateways())
}

func TestPaymentService_PaymentGateway(t *testing.T) {
	ps := application.PaymentService{}
	ps.Inject(func() map[string]interfaces.WebCartPaymentGateway {
		return map[string]interfaces.WebCartPaymentGateway{
			"gateway-code": &mocks.WebCartPaymentGateway{},
		}
	})

	gateway, err := ps.PaymentGateway("non-existing")
	assert.Nil(t, gateway)
	assert.EqualError(t, err, "Payment gateway non-existing not found")

	gateway, err = ps.PaymentGateway("gateway-code")
	assert.Equal(t, &mocks.WebCartPaymentGateway{}, gateway)
	assert.Nil(t, err)
}

func TestPaymentService_PaymentGatewayByCart(t *testing.T) {
	ps := application.PaymentService{}
	ps.Inject(func() map[string]interfaces.WebCartPaymentGateway {
		return map[string]interfaces.WebCartPaymentGateway{
			"gateway-code": &mocks.WebCartPaymentGateway{},
		}
	})

	// cart without payment selection
	cart := cartDomain.Cart{}
	gateway, err := ps.PaymentGatewayByCart(cart)
	assert.Nil(t, gateway)
	assert.EqualError(t, err, "PaymentSelection not set")

	// cart with valid payment selection and working gateway
	cart = cartDomain.Cart{}
	paymentSelection, _ := cartDomain.NewDefaultPaymentSelection("gateway-code", map[string]string{domain.ChargeTypeMain: "main"}, cart)
	cart.PaymentSelection = paymentSelection
	gateway, err = ps.PaymentGatewayByCart(cart)
	assert.Equal(t, &mocks.WebCartPaymentGateway{}, gateway)
	assert.Nil(t, err)

}
