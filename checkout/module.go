package checkout

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/contextstore"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/payment"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	flamingographql "flamingo.me/graphql"
	"github.com/go-playground/form"

	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/locker"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql"
)

type (
	// Module registers our profiler
	Module struct {
		UseFakeSourcingService bool   `inject:"config:checkout.useFakeSourcingService,optional"`
		PlaceOrderLockType     string `inject:"config:checkout.placeorder.lockType,optional"`
		PlaceOrderContextStore string `inject:"config:checkout.placeorder.contextStore,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {

	injector.Bind((*form.Decoder)(nil)).ToProvider(form.NewDecoder).AsEagerSingleton()
	if m.UseFakeSourcingService {
		injector.Override((*domain.SourcingService)(nil), "").To(infrastructure.FakeSourcingService{})
	}

	if m.PlaceOrderLockType == "clusterlock" {
		injector.Bind(new(placeorder.TryLock)).ToProvider(locker.NewRedis).In(dingo.Singleton)
	} else {
		injector.Bind(new(placeorder.TryLock)).To(&locker.Simple{}).In(dingo.Singleton)
	}

	if m.PlaceOrderContextStore == "redis" {
		injector.Bind(new(process.ContextStore)).To(new(contextstore.Redis)).In(dingo.Singleton)
	} else {
		injector.Bind(new(process.ContextStore)).To(new(contextstore.Memory)).In(dingo.Singleton)
	}

	injector.Bind(new(process.PaymentValidatorFunc)).ToInstance(placeorder.PaymentValidator)

	injector.Bind(new(process.State)).AnnotatedWith("startState").To(states.New{})
	injector.Bind(new(process.State)).AnnotatedWith("failedState").To(states.Failed{})
	injector.BindMap(new(process.State), new(states.New).Name()).To(states.New{})
	injector.BindMap(new(process.State), new(states.ValidateCart).Name()).To(states.ValidateCart{})
	injector.BindMap(new(process.State), new(states.CreatePayment).Name()).To(states.CreatePayment{})
	injector.BindMap(new(process.State), new(states.CompleteCart).Name()).To(states.CompleteCart{})
	injector.BindMap(new(process.State), new(states.CompletePayment).Name()).To(states.CompletePayment{})
	injector.BindMap(new(process.State), new(states.PlaceOrder).Name()).To(states.PlaceOrder{})
	injector.BindMap(new(process.State), new(states.ValidatePayment).Name()).To(states.ValidatePayment{})
	injector.BindMap(new(process.State), new(states.WaitForCustomer).Name()).To(states.WaitForCustomer{})
	injector.BindMap(new(process.State), new(states.Success).Name()).To(states.Success{})
	injector.BindMap(new(process.State), new(states.Failed).Name()).To(states.Failed{})
	injector.BindMap(new(process.State), new(states.ShowIframe).Name()).To(states.ShowIframe{})
	injector.BindMap(new(process.State), new(states.ShowHTML).Name()).To(states.ShowHTML{})
	injector.BindMap(new(process.State), new(states.Redirect).Name()).To(states.Redirect{})
	injector.BindMap(new(process.State), new(states.PostRedirect).Name()).To(states.PostRedirect{})

	// bind internal states to graphQL states
	injector.BindMap(new(dto.State), new(states.New).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.ValidateCart).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.CreatePayment).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.CompleteCart).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.CompletePayment).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.PlaceOrder).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.ValidatePayment).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.WaitForCustomer).Name()).To(dto.WaitForCustomer{})
	injector.BindMap(new(dto.State), new(states.Success).Name()).To(dto.Success{})
	injector.BindMap(new(dto.State), new(states.Failed).Name()).To(dto.Failed{})
	injector.BindMap(new(dto.State), new(states.ShowHTML).Name()).To(dto.ShowHTML{})
	injector.BindMap(new(dto.State), new(states.ShowIframe).Name()).To(dto.ShowIframe{})
	injector.BindMap(new(dto.State), new(states.Redirect).Name()).To(dto.Redirect{})
	injector.BindMap(new(dto.State), new(states.PostRedirect).Name()).To(dto.PostRedirect{})

	web.BindRoutes(injector, new(routes))

	injector.BindMulti(new(flamingographql.Service)).To(graphql.Service{})
}

// DefaultConfig for checkout module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"checkout": config.Map{
			"useDeliveryForms":    true,
			"usePersonalDataForm": false,
			"skipReviewAction":    false,
		},
	}
}

type routes struct {
	controller *controller.CheckoutController
}

// Inject required controller
func (r *routes) Inject(controller *controller.CheckoutController) {
	r.controller = controller
}

// Routes  configuration for checkout controllers
func (r *routes) Routes(registry *web.RouterRegistry) {
	// routes
	registry.HandleAny("checkout.start", r.controller.StartAction)
	registry.Route("/checkout/start", "checkout.start")

	registry.HandleAny("checkout.review", r.controller.ReviewAction)
	registry.Route("/checkout/review", `checkout.review`)

	registry.HandleAny("checkout", r.controller.SubmitCheckoutAction)
	registry.Route("/checkout", "checkout")

	registry.HandleAny("checkout.payment", r.controller.PaymentAction)
	registry.Route("/checkout/payment", "checkout.payment")

	registry.HandleAny("checkout.success", r.controller.SuccessAction)
	registry.Route("/checkout/success", "checkout.success")

	registry.HandleAny("checkout.expired", r.controller.ExpiredAction)
	registry.Route("/checkout/expired", "checkout.expired")

	registry.HandleAny("checkout.placeorder", r.controller.PlaceOrderAction)
	registry.Route("/checkout/placeorder", "checkout.placeorder")
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(cart.Module),
		new(payment.Module),
		new(flamingo.SessionModule),
	}
}
