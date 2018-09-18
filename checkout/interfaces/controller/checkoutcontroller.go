package controller

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/url"
	"strings"

	cartApplication "flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/checkout/application"
	paymentDomain "flamingo.me/flamingo-commerce/checkout/domain/payment"
	"flamingo.me/flamingo-commerce/checkout/interfaces/controller/formDto"
	customerApplication "flamingo.me/flamingo-commerce/customer/application"
	authApplication "flamingo.me/flamingo/core/auth/application"
	formApplicationService "flamingo.me/flamingo/core/form/application"
	formDomain "flamingo.me/flamingo/core/form/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/opencensus"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
	"github.com/gorilla/sessions"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

type (
	PaymentProviderProvider func() map[string]paymentDomain.PaymentProvider

	// CheckoutViewData represents the checkout view data
	CheckoutViewData struct {
		DecoratedCart        cart.DecoratedCart
		Form                 formDomain.Form
		CartValidationResult cart.CartValidationResult
		ErrorInfos           ViewErrorInfos
		PaymentProviders     map[string]paymentDomain.PaymentProvider
	}

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
		CartTotals cart.CartTotals
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

	PlaceOrderPaymentInfo struct {
		Provider       string
		Method         string
		Amount         float64
		Title          string
		CreditCardInfo *cart.CreditCardInfo
	}

	EmptyCartInfo struct {
		CartExpired bool
	}

	// CheckoutController represents the checkout controller with its injectsions
	CheckoutController struct {
		responder.RenderAware
		responder.RedirectAware
		router *router.Router

		checkoutFormService  *formDto.CheckoutFormService
		orderService         *application.OrderService
		paymentService       *application.PaymentService
		decoratedCartFactory *cart.DecoratedCartFactory

		skipStartAction                 bool
		skipReviewAction                bool
		showReviewStepAfterPaymentError bool
		showEmptyCartPageIfNoItems      bool
		redirectToCartOnInvalideCart    bool

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

func (cc *CheckoutController) Inject(
	renderAware responder.RenderAware,
	redirectAware responder.RedirectAware,
	router *router.Router,
	checkoutFormService *formDto.CheckoutFormService,
	orderService *application.OrderService,
	paymentService *application.PaymentService,
	decoratedCartFactory *cart.DecoratedCartFactory,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	userService *authApplication.UserService,
	logger flamingo.Logger,
	customerApplicationService *customerApplication.Service,
	paymentProvider PaymentProviderProvider,
	config *struct {
		SkipStartAction                 bool   `inject:"config:checkout.skipStartAction,optional"`
		SkipReviewAction                bool   `inject:"config:checkout.skipReviewAction,optional"`
		ShowReviewStepAfterPaymentError bool   `inject:"config:checkout.showReviewStepAfterPaymentError,optional"`
		ShowEmptyCartPageIfNoItems      bool   `inject:"config:checkout.showEmptyCartPageIfNoItems,optional"`
		RedirectToCartOnInvalideCart    bool   `inject:"config:checkout.redirectToCartOnInvalideCart,optional"`
		BaseURL                         string `inject:"config:canonicalurl.baseurl"`
	},
) {
	cc.RenderAware = renderAware
	cc.RedirectAware = redirectAware
	cc.router = router

	cc.checkoutFormService = checkoutFormService
	cc.orderService = orderService
	cc.paymentService = paymentService
	cc.decoratedCartFactory = decoratedCartFactory

	cc.skipStartAction = config.SkipStartAction
	cc.skipReviewAction = config.SkipReviewAction
	cc.showReviewStepAfterPaymentError = config.ShowReviewStepAfterPaymentError
	cc.showEmptyCartPageIfNoItems = config.ShowEmptyCartPageIfNoItems
	cc.redirectToCartOnInvalideCart = config.RedirectToCartOnInvalideCart

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
func (cc *CheckoutController) StartAction(ctx context.Context, r *web.Request) web.Response {
	//Guard Clause if Cart cannout be fetched

	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil).Hook(web.NoCache)
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		if cc.showEmptyCartPageIfNoItems {
			return cc.Render(ctx, "checkout/emptycart", nil)
		}
		return cc.Render(ctx, "checkout/startcheckout", CheckoutViewData{
			DecoratedCart: *decoratedCart,
		}).Hook(web.NoCache)
	}

	if cc.userService.IsLoggedIn(ctx, r.Session()) {
		return cc.Redirect("checkout.user", nil)
	}

	if cc.skipStartAction {
		return cc.Redirect("checkout.guest", nil)
	}

	return cc.Render(ctx, "checkout/startcheckout", CheckoutViewData{
		DecoratedCart: *decoratedCart,
	}).Hook(web.NoCache)
}

func (cc *CheckoutController) hasAvailablePaymentProvider() bool {
	return len(cc.getPaymentProviders()) > 0
}

func (cc *CheckoutController) getPayment(ctx context.Context, paymentProviderCode string, paymentMethodCode string) (paymentDomain.PaymentProvider, *paymentDomain.PaymentMethod, error) {
	providers := cc.getPaymentProviders()

	provider := providers[paymentProviderCode]

	if provider == nil {
		return nil, nil, errors.New("Payment provider " + paymentProviderCode + " not found")
	}

	paymentMethods := provider.GetPaymentMethods()

	var paymentMethod *paymentDomain.PaymentMethod
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
func (cc *CheckoutController) SubmitUserCheckoutAction(ctx context.Context, r *web.Request) web.Response {
	//Guard
	if !cc.userService.IsLoggedIn(ctx, r.Session()) {
		return cc.Redirect("checkout.start", nil).Hook(web.NoCache)
	}

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil).Hook(web.NoCache)
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	customer, _ := cc.customerApplicationService.GetForAuthenticatedUser(ctx, r.Session())
	// set the customer on the form service even if nil + err is returned here
	cc.checkoutFormService.Customer = customer

	return cc.showCheckoutFormAndHandleSubmit(ctx, r, cc.checkoutFormService, "checkout/usercheckout").Hook(web.NoCache)
}

// SubmitGuestCheckoutAction handles the guest order submit
func (cc *CheckoutController) SubmitGuestCheckoutAction(ctx context.Context, r *web.Request) web.Response {
	cc.checkoutFormService.Customer = nil
	if cc.userService.IsLoggedIn(ctx, r.Session()) {
		return cc.Redirect("checkout.user", nil).Hook(web.NoCache)
	}

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil).Hook(web.NoCache)
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	return cc.showCheckoutFormAndHandleSubmit(ctx, r, cc.checkoutFormService, "checkout/guestcheckout").Hook(web.NoCache)
}

// ProcessPaymentAction functions as a return/notification URL for Payment Providers
func (cc *CheckoutController) ProcessPaymentAction(ctx context.Context, r *web.Request) web.Response {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil).Hook(web.NoCache)
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, "checkout/emptycart", nil).Hook(web.NoCache)
	}

	providercode := r.MustParam1("providercode")
	methodcode := r.MustParam1("methodcode")

	provider, paymentMethod, err := cc.getPayment(ctx, providercode, methodcode)

	cartPayment, err := provider.ProcessPayment(ctx, r, &decoratedCart.Cart, paymentMethod, nil)
	if err != nil {
		if cc.showReviewStepAfterPaymentError && !cc.skipReviewAction {
			return cc.showReviewFormWithErrors(ctx, *decoratedCart, providercode, methodcode, err).Hook(web.NoCache)
		}
		return cc.showCheckoutFormWithErrors(ctx, r, "", *decoratedCart, nil, err).Hook(web.NoCache)
	}

	response, err := cc.placeOrder(ctx, r.Session(), *cartPayment, *decoratedCart)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, "", *decoratedCart, nil, err).Hook(web.NoCache)
	}
	return response
}

// SuccessAction handles the order success action
func (cc *CheckoutController) SuccessAction(ctx context.Context, r *web.Request) web.Response {
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

			return cc.Render(ctx, "checkout/success", viewData).Hook(web.NoCache)
		}
	}
	return cc.Redirect("checkout.expired", nil).Hook(web.NoCache)
}

func (cc *CheckoutController) ExpiredAction(ctx context.Context, r *web.Request) web.Response {
	if cc.showEmptyCartPageIfNoItems {
		return cc.Render(ctx, "checkout/emptycart", EmptyCartInfo{
			CartExpired: true,
		}).Hook(web.NoCache)
	}
	return cc.Render(ctx, "checkout/expired", nil).Hook(web.NoCache)
}

func (cc *CheckoutController) getPaymentReturnUrl(PaymentProvider string, PaymentMethod string) *url.URL {
	baseUrl := cc.baseURL
	paymentUrl := cc.router.URL("checkout.processpayment", router.P{"providercode": PaymentProvider, "methodcode": PaymentMethod})

	rawUrl := strings.TrimRight(baseUrl, "/") + paymentUrl.String()

	urlResult, _ := url.Parse(rawUrl)

	return urlResult
}

//showCheckoutFormAndHandleSubmit - Action that shows the form (either customer or guest)
func (cc *CheckoutController) showCheckoutFormAndHandleSubmit(ctx context.Context, r *web.Request, formservice *formDto.CheckoutFormService, template string) web.Response {
	session := r.Session()

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, session)
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if formservice == nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error CheckoutFormService not present!")
		return cc.Render(ctx, "checkout/carterror", nil)
	}
	formservice.Cart = &decoratedCart.Cart

	if !cc.hasAvailablePaymentProvider() {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error No Payment set")
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	form, e := formApplicationService.ProcessFormRequest(ctx, r, formservice)
	// return on error (template need to handle error display)
	if e != nil {
		return cc.Render(ctx, template, CheckoutViewData{
			DecoratedCart:        *decoratedCart,
			CartValidationResult: cc.applicationCartService.ValidateCart(ctx, session, decoratedCart),
			Form:                 form,
			PaymentProviders:     cc.getPaymentProviders(),
		})
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		if cc.showEmptyCartPageIfNoItems {
			return cc.Render(ctx, "checkout/emptycart", nil)
		}

		return cc.Render(ctx, template, CheckoutViewData{
			DecoratedCart:        *decoratedCart,
			CartValidationResult: cc.applicationCartService.ValidateCart(ctx, session, decoratedCart),
			Form:                 form,
			PaymentProviders:     cc.getPaymentProviders(),
		})
	}

	if form.IsValidAndSubmitted() {

		if checkoutFormData, ok := form.Data.(formDto.CheckoutFormData); ok {

			billingAddress, shippingAddress := formDto.MapAddresses(checkoutFormData)
			person := formDto.MapPerson(checkoutFormData)

			err := cc.orderService.CurrentCartSaveInfos(ctx, session, billingAddress, shippingAddress, person)
			if err != nil {
				return cc.showCheckoutFormWithErrors(ctx, r, template, *decoratedCart, &form, err)
			}

			if cc.skipReviewAction {

				return cc.processPaymentOrPlaceOrderDirectly(ctx, r, checkoutFormData.SelectedPaymentProvider, checkoutFormData.SelectedPaymentProviderMethod, template, &form)
			}
			return cc.Redirect("checkout.review", nil)
		} else {
			cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error cannot type convert to CheckoutFormData!")
			return cc.Render(ctx, "checkout/carterror", nil)
		}
	} else {
		if form.IsSubmitted && form.HasGeneralErrors() {
			cc.logger.WithField("category", "checkout").Warn("CheckoutForm has general error: %#v", form.ValidationInfo.GeneralErrors)
		}
	}

	cc.logger.Debug("paymentProviders %#v", cc.getPaymentProviders())
	//Default: Form not submitted yet or submitted with validation errors:
	return cc.Render(ctx, template, CheckoutViewData{
		DecoratedCart:        *decoratedCart,
		CartValidationResult: cc.applicationCartService.ValidateCart(ctx, session, decoratedCart),
		Form:                 form,
		PaymentProviders:     cc.getPaymentProviders(),
	})
}

//showCheckoutFormWithErrors - error handling that is called form many places... It will show the checkoutform and the error
// template and form is optional - if it is not goven it is autodetected and prefilled from the infos in the cart
func (cc *CheckoutController) showCheckoutFormWithErrors(ctx context.Context, r *web.Request, template string, decoratedCart cart.DecoratedCart, form *formDomain.Form, err error) web.Response {
	if template == "" {
		template = "checkout/guestcheckout"
		if cc.userService.IsLoggedIn(ctx, r.Session()) {
			template = "checkout/usercheckout"
		}
	}
	cc.logger.Warn("Place Order Error: %s", err.Error())
	if form == nil {
		cc.checkoutFormService.Cart = &decoratedCart.Cart
		newForm, _ := formApplicationService.GetUnsubmittedForm(ctx, r, cc.checkoutFormService)
		form = &newForm
	}

	return cc.Render(ctx, template, CheckoutViewData{
		DecoratedCart:        decoratedCart,
		CartValidationResult: cc.applicationCartService.ValidateCart(ctx, r.Session(), &decoratedCart),
		Form:                 *form,
		ErrorInfos:           getViewErrorInfo(err),
		PaymentProviders:     cc.getPaymentProviders(),
	})
}

//showReviewFormWithErrors
func (cc *CheckoutController) showReviewFormWithErrors(ctx context.Context, decoratedCart cart.DecoratedCart, selectedProvider string, selectedMethod string, err error) web.Response {
	cc.logger.Warn("Show Error (review step): %s", err.Error())
	viewData := ReviewStepViewData{
		DecoratedCart: decoratedCart,
		ErrorInfos:    getViewErrorInfo(err),
		SelectedPayment: SelectedPayment{
			Provider: selectedProvider,
			Method:   selectedMethod,
		},
	}
	return cc.Render(ctx, "checkout/review", viewData)
}

func getViewErrorInfo(err error) ViewErrorInfos {
	hasPaymentError := false

	if paymentErr, ok := err.(*paymentDomain.PaymentError); ok {
		hasPaymentError = paymentErr.ErrorCode != paymentDomain.PaymentCancelled
	}

	errorInfos := ViewErrorInfos{
		HasError:        true,
		ErrorMessage:    err.Error(),
		HasPaymentError: hasPaymentError,
	}
	return errorInfos
}

func (cc *CheckoutController) processPaymentOrPlaceOrderDirectly(ctx context.Context, r *web.Request, selectedPaymentProvider string, selectedPaymentProviderMethod string, orderFormTemplate string, checkoutForm *formDomain.Form) web.Response {
	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	//procces Payment:
	paymentProvider, paymentMethod, err := cc.getPayment(ctx, selectedPaymentProvider, selectedPaymentProviderMethod)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, orderFormTemplate, *decoratedCart, checkoutForm, err)
	}
	//Payment Method requests an redirect - execute it
	if paymentMethod.IsExternalPayment {
		returnUrl := cc.getPaymentReturnUrl(paymentProvider.GetCode(), paymentMethod.Code)
		hostedPaymentPageResponse, err := paymentProvider.RedirectExternalPayment(ctx, r, &decoratedCart.Cart, paymentMethod, returnUrl)
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

func (cc *CheckoutController) placeOrder(ctx context.Context, session *sessions.Session, cartPayment cart.CartPayment, decoratedCart cart.DecoratedCart) (web.Response, error) {
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

	return cc.Redirect("checkout.success", nil).With("checkout.success.data", PlaceOrderFlashData{
		PlacedOrderInfos: placedOrderInfos,
		Email:            email,
		PlacedCart:       decoratedCart.Cart,
		PaymentInfos:     placeOrderPaymentInfos,
	}), nil

}
func (cc *CheckoutController) getPaymentProviders() map[string]paymentDomain.PaymentProvider {
	result := make(map[string]paymentDomain.PaymentProvider)

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

// ReviewAction
func (cc *CheckoutController) ReviewAction(ctx context.Context, r *web.Request) web.Response {
	if cc.skipReviewAction {
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	selectedProvider, _ := r.Form1("selectedPaymentProvider")
	selectedMethod, _ := r.Form1("selectedPaymentProviderMethod")
	proceed, _ := r.Form1("proceed")
	termsAndConditions, _ := r.Form1("termsAndConditions")

	cc.logger.Debug("ReviewAction: selectedProvider: %v / selectedMethod: %v / proceed: %v / termsAndConditions: %v", selectedProvider, selectedMethod, proceed, termsAndConditions)

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil).Hook(web.NoCache)
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, "checkout/emptycart", nil).Hook(web.NoCache)
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	if proceed == "1" && termsAndConditions == "1" && selectedProvider != "" && selectedMethod != "" {
		return cc.processPaymentOrPlaceOrderDirectly(ctx, r, selectedProvider, selectedMethod, "", nil).Hook(web.NoCache)
	}

	viewData := ReviewStepViewData{
		DecoratedCart: *decoratedCart,
		SelectedPayment: SelectedPayment{
			Provider: selectedProvider,
			Method:   selectedMethod,
		},
	}
	return cc.Render(ctx, "checkout/review", viewData).Hook(web.NoCache)

}

//getCommonGuardRedirects - checks config and may return a redirect that should be executed before the common checkou actions
func (cc *CheckoutController) getCommonGuardRedirects(ctx context.Context, session *sessions.Session, decoratedCart *cart.DecoratedCart) web.Response {
	if cc.redirectToCartOnInvalideCart {
		result := cc.applicationCartService.ValidateCart(ctx, session, decoratedCart)
		if !result.IsValid() {
			cc.logger.WithField("category", "checkout").Info("StartAction > RedirectToCartOnInvalideCart")
			return cc.Redirect("cart.view", nil).Hook(web.NoCache)
		}
	}
	return nil
}
