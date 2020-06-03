package payment

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces/controller"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module registers our payment module
	Module struct {
		EnableOfflinePayment bool `inject:"config:commerce.payment.enableOfflinePaymentGateway,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.EnableOfflinePayment {
		injector.BindMap((*interfaces.WebCartPaymentGateway)(nil), interfaces.OfflineWebCartPaymentGatewayCode).To(interfaces.OfflineWebCartPaymentGateway{})
	}

	web.BindRoutes(injector, new(routes))
}

type routes struct {
	paymentAPIController *controller.PaymentAPIController
}

func (r *routes) Inject(apiController *controller.PaymentAPIController) {
	r.paymentAPIController = apiController
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("payment.status", r.paymentAPIController.Status)
	registry.Route("/api/payment/status", "payment.status")
	registry.Route("/api/v1/payment/status", "payment.status")

}
