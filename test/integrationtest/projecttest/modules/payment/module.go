package payment

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

type (
	// Module registers our fake payment profile
	Module struct{}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMap((*interfaces.WebCartPaymentGateway)(nil), FakePaymentGateway).To(new(FakeGateway)).In(dingo.Singleton)
}
