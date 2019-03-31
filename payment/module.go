package payment

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

type (
	// Module registers our profiler
	Module struct {
		EnableOfflinePayment bool `inject:"config:commerce.payment.enableOfflinePaymentGateway,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.EnableOfflinePayment {
		injector.BindMap((*interfaces.WebCartPaymentGateway)(nil), interfaces.OfflineWebCartPaymentGatewayCode).To(interfaces.OfflineWebCartPaymentGateway{})
	}
}
