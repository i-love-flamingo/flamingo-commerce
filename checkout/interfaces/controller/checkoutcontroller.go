package controller

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/url"
	"strings"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	paymentDomain "flamingo.me/flamingo-commerce/v3/checkout/domain/payment"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller/formdto"
	customerApplication "flamingo.me/flamingo-commerce/v3/customer/application"
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	formApplicationService "go.aoe.com/flamingo/form/application"
	formDomain "go.aoe.com/flamingo/form/domain"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

type (
	// PaymentProviderProvider defines the map of providers for payment providers
	PaymentProviderProvider func() map[string]paymentDomain.Provider

	// CheckoutViewData represents the checkout view data
	CheckoutViewData struct {
		DecoratedCart        cart.DecoratedCart
		Form                 formDomain.Form
		CartValidationResult cart.ValidationResult
		ErrorInfos           ViewErrorInfos
		PaymentProviders     map[string]paymentDomain.Provider
	}

	// ViewErrorInfos defines the error info struct of the checkout controller views
	ViewErrorInfos struct {
		//HasError  indicates that an general error happend
		HasError bool
		//If there is a general error this field is filled and can be used in the template
		ErrorMessage string
		//if the Error happens during processing PAyment (can be used in template to behave special in case of payment errors)
		HasPaymentError bool
	}

	// SuccessViewData represents the success view data
	SuccessViewData struct {
		PaymentInfos        []PlaceOrderPaymentInfo
		PlacedOrderInfos    cart.PlacedOrderInfos
		Email               string
		PlacedDecoratedCart cart.DecoratedCart

		//PlacedDecoratedItems - DEPRECATED!!!
		// PlacedDecoratedItems []cart.DecoratedCartItem

		//CartTotals - Depricated
		CartTotals cart.Totals
	}

	// ReviewStepViewData represents the success view data
	ReviewStepViewData struct {
		DecoratedCart   cart.DecoratedCart
		SelectedPayment SelectedPayment
		ErrorInfos      ViewErrorInfos
	}

	// SelectedPayment represents the success view data
	SelectedPayment struct {
		Provider string
		Method   string
	}

	// PlaceOrderFlashData represents the data passed to the success page - they need to be "glob"able
	PlaceOrderFlashData struct {
		PlacedOrderInfos cart.PlacedOrderInfos
		Email            string
		PaymentInfos     []PlaceOrderPaymentInfo
		PlacedCart       cart.Cart
	}

	// PlaceOrderPaymentInfo struct defines the data of payments on placed orders
	PlaceOrderPaymentInfo struct {
		Provider       string
		Method         string
		Amount         float64
		Title          string
		CreditCardInfo *cart.CreditCardInfo
	}

	// EmptyCartInfo struct defines the data info on empty carts
	EmptyCartInfo struct {
		CartExpired bool
	}

	// CheckoutController represents the checkout controller with its injectsions
	CheckoutController struct {
		responder *web.Responder
		router    *web.Router

		checkoutFormService  *formdto.CheckoutFormService
		orderService         *application.OrderService
		paymentService       *application.PaymentService
		decoratedCartFactory *cart.DecoratedCartFactory
		eventRouter          flamingo.EventRouter

		skipStartAction                 bool
		skipReviewAction                bool
		showReviewStepAfterPaymentError bool
		showEmptyCartPageIfNoItems      bool
		redirectToCartOnInvalideCart    bool
		privacyPolicyRequired           bool

		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService

		userService *authApplication.UserService

		logger flamingo.Logger

		customerApplicationService *customerApplication.Service
		paymentProvider            PaymentProviderProvider

		baseURL string
	}
)

var rt = stats.Int64("flamingo-commerce/orderfailed", "my stat records 1 occurences per error", stats.UnitDimensionless)

func init() {
	gob.Register(PlaceOrderFlashData{})
	opencensus.View("flamingo-commerce/orderfailed/count", rt, view.Count())
}

// Inject dependencies
func (cc *CheckoutController) Inject(
	responder *web.Responder,
	router *web.Router,
	checkoutFormService *formdto.CheckoutFormService,
	orderService *application.OrderService,
	paymentService *application.PaymentService,
	decoratedCartFactory *cart.DecoratedCartFactory,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	userService *authApplication.UserService,
	logger flamingo.Logger,
	customerApplicationService *customerApplication.Service,
	paymentProvider PaymentProviderProvider,
	eventRouter flamingo.EventRouter,
	config *struct {
		SkipStartAction                 bool   `inject:"config:checkout.skipStartAction,optional"`
		SkipReviewAction                bool   `inject:"config:checkout.skipReviewAction,optional"`
		ShowReviewStepAfterPaymentError bool   `inject:"config:checkout.showReviewStepAfterPaymentError,optional"`
		ShowEmptyCartPageIfNoItems      bool   `inject:"config:checkout.showEmptyCartPageIfNoItems,optional"`
		RedirectToCartOnInvalideCart    bool   `inject:"config:checkout.redirectToCartOnInvalideCart,optional"`
		PrivacyPolicyRequired           bool   `inject:"config:checkout.privacyPolicyRequired,optional"`
		BaseURL                         string `inject:"config:canonicalurl.baseurl"`
	},
) {
	cc.responder = responder
	cc.router = router

	cc.checkoutFormService = checkoutFormService
	cc.orderService = orderService
	cc.paymentService = paymentService
	cc.decoratedCartFactory = decoratedCartFactory
	cc.eventRouter = eventRouter

	cc.skipStartAction = config.SkipStartAction
	cc.skipReviewAction = config.SkipReviewAction
	cc.showReviewStepAfterPaymentError = config.ShowReviewStepAfterPaymentError
	cc.showEmptyCartPageIfNoItems = config.ShowEmptyCartPageIfNoItems
	cc.redirectToCartOnInvalideCart = config.RedirectToCartOnInvalideCart
	cc.privacyPolicyRequired = config.PrivacyPolicyRequired

	cc.applicationCartService = applicationCartService
	cc.applicationCartReceiverService = applicationCartReceiverService

	cc.userService = userService

	cc.logger = logger

	cc.customerApplicationService = customerApplicationService
	cc.paymentProvider = paymentProvider

	cc.baseURL = config.BaseURL
}

/*
The checkoutController implements a default process for a checkout:
 * StartAction (supposed to show a switch to go to guest or customer)
 	* can be skipped with a configuration
 * SubmitUserCheckoutAction  OR  SubmitGuestCheckoutAction
 	* both actions are more or less the same (User checkout just populates the customer to the form and uses a different template)
 	* This step is supposed to show a big form (validation and default values are configurable as well)
	* payment can be selected in this step or in the next

 * ReviewAction
	* this step is supposed to show the current cart status just before checkout
		* optional the paymentmethod can also be selected here
	* This step can also be skipped - then directly the placeOrder is handled

*  Optional Payment Step (if the payment requires a redirect the payment page is shown and a redirect back to "ProcessPayment"
* SuccessStep

*/

// StartAction handles the checkout start action
func (cc *CheckoutController) StartAction(ctx context.Context, r *web.Request) web.Result {
	//Guard Clause if Cart cannout be fetched

	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		if cc.showEmptyCartPageIfNoItems {
			return cc.responder.Render("checkout/emptycart", nil)
		}
		return cc.responder.Render("checkout/startcheckout", CheckoutViewData{
			DecoratedCart: *decoratedCart,
		}).SetNoCache()
	}

	if cc.userService.IsLoggedIn(ctx, r.Session()) {
		return cc.responder.RouteRedirect("checkout.user", nil)
	}

	if cc.skipStartAction {
		return cc.responder.RouteRedirect("checkout.guest", nil)
	}

	return cc.responder.Render("checkout/startcheckout", CheckoutViewData{
		DecoratedCart: *decoratedCart,
	}).SetNoCache()
}

func (cc *CheckoutController) hasAvailablePaymentProvider() bool {
	return len(cc.getPaymentProviders()) > 0
}

func (cc *CheckoutController) getPayment(ctx context.Context, paymentProviderCode string, paymentMethodCode string) (paymentDomain.Provider, *paymentDomain.Method, error) {
	providers := cc.getPaymentProviders()

	provider := providers[paymentProviderCode]

	if provider == nil {
		return nil, nil, errors.New("Payment provider " + paymentProviderCode + " not found")
	}

	paymentMethods := provider.GetPaymentMethods()

	var paymentMethod *paymentDomain.Method
	for _, method := range paymentMethods {
		if method.Code == paymentMethodCode {
			paymentMethod = &method
			break
		}
	}

	if paymentMethod == nil {
		return nil, nil, errors.New("payment method not found")
	}
	return provider, paymentMethod, nil
}

// SubmitUserCheckoutAction handles the user order submit
func (cc *CheckoutController) SubmitUserCheckoutAction(ctx context.Context, r *web.Request) web.Result {
	//Guard
	if !cc.userService.IsLoggedIn(ctx, r.Session()) {
		r := cc.responder.RouteRedirect("checkout.start", nil)
		r.SetNoCache()
		return r
	}

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	customer, _ := cc.customerApplicationService.GetForAuthenticatedUser(ctx, r.Session())
	// set the customer on the form service even if nil + err is returned here
	cc.checkoutFormService.SetCustomer(customer)

	return cc.showCheckoutFormAndHandleSubmit(ctx, r, cc.checkoutFormService, "checkout/usercheckout")
}

// SubmitGuestCheckoutAction handles the guest order submit
func (cc *CheckoutController) SubmitGuestCheckoutAction(ctx context.Context, r *web.Request) web.Result {
	cc.checkoutFormService.SetCustomer(nil)
	if cc.userService.IsLoggedIn(ctx, r.Session()) {
		resp := cc.responder.RouteRedirect("checkout.user", nil)
		resp.SetNoCache()
		return resp
	}

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	return cc.showCheckoutFormAndHandleSubmit(ctx, r, cc.checkoutFormService, "checkout/guestcheckout")
}

// ProcessPaymentAction functions as a return/notification URL for Payment Providers
func (cc *CheckoutController) ProcessPaymentAction(ctx context.Context, r *web.Request) web.Result {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.responder.Render("checkout/emptycart", nil).SetNoCache()
	}

	providercode := r.Params["providercode"]
	methodcode := r.Params["methodcode"]

	provider, paymentMethod, err := cc.getPayment(ctx, providercode, methodcode)

	cartPayment, err := provider.ProcessPayment(ctx, r, &decoratedCart.Cart, paymentMethod, nil)
	if err != nil {
		if cc.showReviewStepAfterPaymentError && !cc.skipReviewAction {
			return cc.showReviewFormWithErrors(ctx, *decoratedCart, providercode, methodcode, err)
		}
		return cc.showCheckoutFormWithErrors(ctx, r, "", *decoratedCart, nil, err)
	}

	response, err := cc.placeOrder(ctx, r.Session(), *cartPayment, *decoratedCart)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, "", *decoratedCart, nil, err)
	}
	return response
}

// SuccessAction handles the order success action
func (cc *CheckoutController) SuccessAction(ctx context.Context, r *web.Request) web.Result {
	flashes := r.Session().Flashes("checkout.success.data")
	if len(flashes) > 0 {
		if placeOrderFlashData, ok := flashes[len(flashes)-1].(PlaceOrderFlashData); ok {
			decoratedCart := cc.decoratedCartFactory.Create(ctx, placeOrderFlashData.PlacedCart)
			viewData := SuccessViewData{
				CartTotals:          placeOrderFlashData.PlacedCart.CartTotals,
				Email:               placeOrderFlashData.Email,
				PaymentInfos:        placeOrderFlashData.PaymentInfos,
				PlacedDecoratedCart: *decoratedCart,
				PlacedOrderInfos:    placeOrderFlashData.PlacedOrderInfos,
			}

			return cc.responder.Render("checkout/success", viewData).SetNoCache()
		}
	}
	resp := cc.responder.RouteRedirect("checkout.expired", nil)
	resp.SetNoCache()
	return resp
}

// ExpiredAction handles the expired cart action
func (cc *CheckoutController) ExpiredAction(ctx context.Context, r *web.Request) web.Result {
	if cc.showEmptyCartPageIfNoItems {
		return cc.responder.Render("checkout/emptycart", EmptyCartInfo{
			CartExpired: true,
		}).SetNoCache()
	}
	return cc.responder.Render("checkout/expired", nil).SetNoCache()
}

func (cc *CheckoutController) getPaymentReturnURL(PaymentProvider string, PaymentMethod string) *url.URL {
	baseURL := cc.baseURL
	paymentURL, _ := cc.router.URL("checkout.processpayment", map[string]string{"providercode": PaymentProvider, "methodcode": PaymentMethod})

	rawURL := strings.TrimRight(baseURL, "/") + paymentURL.String()

	urlResult, _ := url.Parse(rawURL)

	return urlResult
}

//showCheckoutFormAndHandleSubmit - Action that shows the form (either customer or guest)
func (cc *CheckoutController) showCheckoutFormAndHandleSubmit(ctx context.Context, r *web.Request, formservice *formdto.CheckoutFormService, template string) web.Result {
	session := r.Session()

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, session)
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	if formservice == nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error CheckoutFormService not present!")
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	formservice.SetCart(&decoratedCart.Cart)

	if !cc.hasAvailablePaymentProvider() {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error No Payment set")
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	form, e := formApplicationService.ProcessFormRequest(ctx, r, formservice)
	// return on error (template need to handle error display)
	if e != nil {
		return cc.responder.Render(template, CheckoutViewData{
			DecoratedCart:        *decoratedCart,
			CartValidationResult: cc.applicationCartService.ValidateCart(ctx, session, decoratedCart),
			Form:                 form,
			PaymentProviders:     cc.getPaymentProviders(),
		}).SetNoCache()
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		if cc.showEmptyCartPageIfNoItems {
			return cc.responder.Render("checkout/emptycart", nil).SetNoCache()
		}

		return cc.responder.Render(template, CheckoutViewData{
			DecoratedCart:        *decoratedCart,
			CartValidationResult: cc.applicationCartService.ValidateCart(ctx, session, decoratedCart),
			Form:                 form,
			PaymentProviders:     cc.getPaymentProviders(),
		}).SetNoCache()
	}

	if form.IsValidAndSubmitted() {
		if checkoutFormData, ok := form.Data.(formdto.CheckoutFormData); ok {
			billingAddress, shippingAddress := formdto.MapAddresses(checkoutFormData)
			person := formdto.MapPerson(checkoutFormData)
			additionalData := formservice.GetAdditionalData(checkoutFormData)

			err := cc.orderService.CurrentCartSaveInfos(ctx, session, billingAddress, shippingAddress, person, additionalData)
			if err != nil {
				return cc.showCheckoutFormWithErrors(ctx, r, template, *decoratedCart, &form, err)
			}

			if cc.skipReviewAction {
				return cc.processPaymentOrPlaceOrderDirectly(ctx, r, template, &form)
			}

			resp := cc.responder.RouteRedirect("checkout.review", nil)
			resp.SetNoCache()
			return resp
		}

		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error cannot type convert to CheckoutFormData!")
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()

	}

	if form.IsSubmitted && form.HasGeneralErrors() {
		cc.logger.WithField("category", "checkout").Warn("CheckoutForm has general error: %#v", form.ValidationInfo.GeneralErrors)
	}

	cc.logger.Debug("paymentProviders %#v", cc.getPaymentProviders())
	//Default: Form not submitted yet or submitted with validation errors:
	return cc.responder.Render(template, CheckoutViewData{
		DecoratedCart:        *decoratedCart,
		CartValidationResult: cc.applicationCartService.ValidateCart(ctx, session, decoratedCart),
		Form:                 form,
		PaymentProviders:     cc.getPaymentProviders(),
	}).SetNoCache()
}

//showCheckoutFormWithErrors - error handling that is called form many places... It will show the checkoutform and the error
// template and form is optional - if it is not goven it is autodetected and prefilled from the infos in the cart
func (cc *CheckoutController) showCheckoutFormWithErrors(ctx context.Context, r *web.Request, template string, decoratedCart cart.DecoratedCart, form *formDomain.Form, err error) web.Result {
	if template == "" {
		template = "checkout/guestcheckout"
		if cc.userService.IsLoggedIn(ctx, r.Session()) {
			template = "checkout/usercheckout"
		}
	}
	cc.logger.Warn("Place Order Error: %s", err.Error())
	if form == nil {
		cc.checkoutFormService.SetCart(&decoratedCart.Cart)
		newForm, _ := formApplicationService.GetUnsubmittedForm(ctx, r, cc.checkoutFormService)
		form = &newForm
	}

	return cc.responder.Render(template, CheckoutViewData{
		DecoratedCart:        decoratedCart,
		CartValidationResult: cc.applicationCartService.ValidateCart(ctx, r.Session(), &decoratedCart),
		Form:                 *form,
		ErrorInfos:           getViewErrorInfo(err),
		PaymentProviders:     cc.getPaymentProviders(),
	}).SetNoCache()
}

//showReviewFormWithErrors
func (cc *CheckoutController) showReviewFormWithErrors(ctx context.Context, decoratedCart cart.DecoratedCart, selectedProvider string, selectedMethod string, err error) web.Result {
	cc.logger.Warn("Show Error (review step): %s", err.Error())
	viewData := ReviewStepViewData{
		DecoratedCart: decoratedCart,
		ErrorInfos:    getViewErrorInfo(err),
		SelectedPayment: SelectedPayment{
			Provider: selectedProvider,
			Method:   selectedMethod,
		},
	}
	return cc.responder.Render("checkout/review", viewData).SetNoCache()
}

func getViewErrorInfo(err error) ViewErrorInfos {
	hasPaymentError := false

	if paymentErr, ok := err.(*paymentDomain.Error); ok {
		hasPaymentError = paymentErr.ErrorCode != paymentDomain.PaymentCancelled
	}

	errorInfos := ViewErrorInfos{
		HasError:        true,
		ErrorMessage:    err.Error(),
		HasPaymentError: hasPaymentError,
	}
	return errorInfos
}

func (cc *CheckoutController) processPaymentOrPlaceOrderDirectly(ctx context.Context, r *web.Request, orderFormTemplate string, checkoutForm *formDomain.Form) web.Result {
	//Guard Clause if Cart can not be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	//procces Payment:
	paymentProvider, paymentMethod, err := cc.getPayment(ctx, decoratedCart.Cart.AdditionalData.SelectedPayment.Provider, decoratedCart.Cart.AdditionalData.SelectedPayment.Method)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, orderFormTemplate, *decoratedCart, checkoutForm, err)
	}
	//Payment Method requests an redirect - execute it
	if paymentMethod.IsExternalPayment {
		returnURL := cc.getPaymentReturnURL(paymentProvider.GetCode(), paymentMethod.Code)
		hostedPaymentPageResponse, err := paymentProvider.RedirectExternalPayment(ctx, r, &decoratedCart.Cart, paymentMethod, returnURL)
		if err != nil {
			return cc.showCheckoutFormWithErrors(ctx, r, orderFormTemplate, *decoratedCart, checkoutForm, err)
		}
		return hostedPaymentPageResponse
	}

	//Paymentmethod that need no external Redirect - can be processed directly
	cartPayment, err := paymentProvider.ProcessPayment(ctx, r, &decoratedCart.Cart, paymentMethod, nil)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, orderFormTemplate, *decoratedCart, checkoutForm, err)
	}

	response, err := cc.placeOrder(ctx, r.Session(), *cartPayment, *decoratedCart)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, orderFormTemplate, *decoratedCart, checkoutForm, err)
	}
	return response
}

func (cc *CheckoutController) getContactMail(cart cart.Cart) string {
	//Get Email from either the cart
	shippingEmail := cart.GetMainShippingEMail()
	if shippingEmail == "" {
		shippingEmail = cart.BillingAdress.Email
	}
	return shippingEmail
}

func (cc *CheckoutController) placeOrder(ctx context.Context, session *web.Session, cartPayment cart.Payment, decoratedCart cart.DecoratedCart) (web.Result, error) {
	placedOrderInfos, err := cc.orderService.CurrentCartPlaceOrder(ctx, session, cartPayment)
	if err != nil {
		name := decoratedCart.Cart.BillingAdress.Firstname + " " + decoratedCart.Cart.BillingAdress.Lastname
		subAmounts := ""
		for _, cartPayment := range cartPayment.PaymentInfos {
			retailer := cartPayment.Title
			if retailer == "" {
				retailer = cartPayment.Provider
			}

			if subAmounts != "" {
				subAmounts += ", "
			}

			subAmounts += retailer + ":" + fmt.Sprintf("%f", cartPayment.Amount)
		}

		// record 5ms per call
		stats.Record(ctx, rt.M(1))

		cc.logger.WithField("category", "checkout").WithField("subcategory", "checkoutError").WithField("errorMsg", err.Error()).Error(fmt.Sprintf("place order failed: cart id: %v / customer-name: %v / total-amount: %v / sub-amounts: %v", decoratedCart.Cart.EntityID, name, decoratedCart.Cart.CartTotals.GrandTotal, subAmounts))
		return nil, err
	}

	email := cc.getContactMail(decoratedCart.Cart)

	var placeOrderPaymentInfos []PlaceOrderPaymentInfo
	for _, cartPayment := range cartPayment.PaymentInfos {
		placeOrderPaymentInfos = append(placeOrderPaymentInfos, PlaceOrderPaymentInfo{
			Method:         cartPayment.Method,
			Provider:       cartPayment.Provider,
			Title:          cartPayment.Title,
			Amount:         cartPayment.Amount,
			CreditCardInfo: cartPayment.CreditCardInfo,
		})
	}

	r := web.RequestFromContext(ctx)
	r.Session().AddFlash(PlaceOrderFlashData{
		PlacedOrderInfos: placedOrderInfos,
		Email:            email,
		PlacedCart:       decoratedCart.Cart,
		PaymentInfos:     placeOrderPaymentInfos,
	}, "checkout.success.data")
	return cc.responder.RouteRedirect("checkout.success", nil), nil

}
func (cc *CheckoutController) getPaymentProviders() map[string]paymentDomain.Provider {
	result := make(map[string]paymentDomain.Provider)

	paymentProviders := cc.paymentProvider()

	if paymentProviders != nil {
		for name, paymentProvider := range cc.paymentProvider() {
			if paymentProvider.IsActive() {
				result[name] = paymentProvider
			}
		}
	}

	return result
}

// ReviewAction handles the cart review action
func (cc *CheckoutController) ReviewAction(ctx context.Context, r *web.Request) web.Result {
	if cc.skipReviewAction {
		return cc.responder.Render("checkout/carterror", nil)
	}

	// Invalidate cart cache
	cc.eventRouter.Dispatch(ctx, &cart.InvalidateCartEvent{Session: r.Session()})

	//Guard Clause if cart can not be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	selectedProvider := decoratedCart.Cart.AdditionalData.SelectedPayment.Provider
	selectedMethod := decoratedCart.Cart.AdditionalData.SelectedPayment.Method
	proceed, _ := r.Form1("proceed")
	termsAndConditions, _ := r.Form1("termsAndConditions")
	privacyPolicy, _ := r.Form1("privacyPolicy")

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.responder.Render("checkout/emptycart", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	viewData := ReviewStepViewData{
		DecoratedCart: *decoratedCart,
		SelectedPayment: SelectedPayment{
			Provider: selectedProvider,
			Method:   selectedMethod,
		},
	}

	errorMessage := ""
	// check for privacy policy if required
	if cc.privacyPolicyRequired && privacyPolicy != "1" && proceed == "1" {
		errorMessage = "privacy_policy_required"
		viewData.ErrorInfos = ViewErrorInfos{
			HasError:        true,
			ErrorMessage:    errorMessage,
			HasPaymentError: false,
		}
	}

	// check for terms and conditions if required
	if termsAndConditions != "1" && proceed == "1" {
		if errorMessage != "" {
			errorMessage = errorMessage + ","
		}
		errorMessage = errorMessage + "terms_and_conditions_required"
		viewData.ErrorInfos = ViewErrorInfos{
			HasError:        true,
			ErrorMessage:    errorMessage,
			HasPaymentError: false,
		}
	}

	if proceed == "1" && (!cc.privacyPolicyRequired || privacyPolicy == "1") && termsAndConditions == "1" && selectedProvider != "" && selectedMethod != "" {
		return cc.processPaymentOrPlaceOrderDirectly(ctx, r, "", nil)
	}

	return cc.responder.Render("checkout/review", viewData).SetNoCache()

}

//getCommonGuardRedirects - checks config and may return a redirect that should be executed before the common checkou actions
func (cc *CheckoutController) getCommonGuardRedirects(ctx context.Context, session *web.Session, decoratedCart *cart.DecoratedCart) web.Result {
	if cc.redirectToCartOnInvalideCart {
		result := cc.applicationCartService.ValidateCart(ctx, session, decoratedCart)
		if !result.IsValid() {
			cc.logger.WithField("category", "checkout").Info("StartAction > RedirectToCartOnInvalideCart")
			resp := cc.responder.RouteRedirect("cart.view", nil)
			resp.SetNoCache()
			return resp
		}
	}
	return nil
}
