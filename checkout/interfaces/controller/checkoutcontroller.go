package controller

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller/forms"

	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"

	"flamingo.me/flamingo-commerce/v3/payment/interfaces"

	"flamingo.me/flamingo-commerce/v3/price/domain"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CheckoutViewData represents the checkout view data
	CheckoutViewData struct {
		DecoratedCart        cart.DecoratedCart
		Form                 forms.CheckoutFormComposite
		CartValidationResult cart.ValidationResult
		ErrorInfos           ViewErrorInfos
		AvailablePayments    map[string][]paymentDomain.Method
		CustomerLoggedIn     bool
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
	}

	// ReviewStepViewData represents the success view data
	ReviewStepViewData struct {
		DecoratedCart cart.DecoratedCart
		ErrorInfos    ViewErrorInfos
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
		Gateway         string
		PaymentProvider string
		Method          string
		Amount          domain.Price
		Title           string
		CreditCardInfo  *cart.CreditCardInfo
	}

	// EmptyCartInfo struct defines the data info on empty carts
	EmptyCartInfo struct {
		CartExpired bool
	}

	// CheckoutController represents the checkout controller with its injectsions
	CheckoutController struct {
		responder *web.Responder
		router    *web.Router

		orderService         *application.OrderService
		decoratedCartFactory *cart.DecoratedCartFactory

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

		webCartPaymentGatewayProvider map[string]interfaces.WebCartPaymentGateway

		checkoutFormController *forms.CheckoutFormController
		baseURL                string
	}
)

func init() {
	gob.Register(PlaceOrderFlashData{})
}

// Inject dependencies
func (cc *CheckoutController) Inject(
	responder *web.Responder,
	router *web.Router,
	orderService *application.OrderService,
	decoratedCartFactory *cart.DecoratedCartFactory,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	userService *authApplication.UserService,
	logger flamingo.Logger,
	webCartPaymentGatewayProvider interfaces.WebCartPaymentGatewayProvider,
	checkoutFormController *forms.CheckoutFormController,
	config *struct {
		SkipStartAction                 bool   `inject:"config:checkout.skipStartAction,optional"`
		SkipReviewAction                bool   `inject:"config:checkout.skipReviewAction,optional"`
		ShowReviewStepAfterPaymentError bool   `inject:"config:checkout.showReviewStepAfterPaymentError,optional"`
		ShowEmptyCartPageIfNoItems      bool   `inject:"config:checkout.showEmptyCartPageIfNoItems,optional"`
		RedirectToCartOnInvalideCart    bool   `inject:"config:checkout.redirectToCartOnInvalideCart,optional"`
		PrivacyPolicyRequired           bool   `inject:"config:checkout.privacyPolicyRequired,optional"`
		BaseURL                         string `inject:"config:canonicalurl.baseurl,optional"`
	},
) {
	cc.responder = responder
	cc.router = router

	cc.checkoutFormController = checkoutFormController
	cc.orderService = orderService
	cc.decoratedCartFactory = decoratedCartFactory

	cc.skipStartAction = config.SkipStartAction
	cc.skipReviewAction = config.SkipReviewAction
	cc.showReviewStepAfterPaymentError = config.ShowReviewStepAfterPaymentError
	cc.showEmptyCartPageIfNoItems = config.ShowEmptyCartPageIfNoItems
	cc.redirectToCartOnInvalideCart = config.RedirectToCartOnInvalideCart
	cc.privacyPolicyRequired = config.PrivacyPolicyRequired

	cc.applicationCartService = applicationCartService
	cc.applicationCartReceiverService = applicationCartReceiverService

	cc.userService = userService

	cc.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "checkoutController")

	cc.webCartPaymentGatewayProvider = webCartPaymentGatewayProvider()
	cc.baseURL = config.BaseURL
}

/*
The checkoutController implements a default process for a checkout:
 * StartAction (supposed to show a switch to go to guest or customer)
 	* can be skipped with a configuration
 * SubmitCheckoutAction
 	* This step is supposed to show a big form (validation and default values are configurable as well)
	* payment can be selected in this step or in the next
	* In cases a customer is logged in the form is prepopulated

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
		cc.logger.Error("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	viewData := cc.getBasicViewData(ctx, r.Session(), *decoratedCart)
	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		if cc.showEmptyCartPageIfNoItems {
			return cc.responder.Render("checkout/emptycart", nil)
		}
		return cc.responder.Render("checkout/startcheckout", viewData).SetNoCache()
	}

	if cc.userService.IsLoggedIn(ctx, r.Session()) {
		return cc.responder.RouteRedirect("checkout", nil)
	}

	if cc.skipStartAction {
		return cc.responder.RouteRedirect("checkout", nil)
	}

	return cc.responder.Render("checkout/startcheckout", viewData).SetNoCache()
}

func (cc *CheckoutController) getPaymentGatewayAndSelectedMethod(ctx context.Context, paymentGatewayCode string) (interfaces.WebCartPaymentGateway, error) {

	gateway, ok := cc.webCartPaymentGatewayProvider[paymentGatewayCode]
	if !ok {
		return nil, errors.New("Payment gateway " + paymentGatewayCode + " not found")
	}

	return gateway, nil
}

// SubmitCheckoutAction handles the main checkout
func (cc *CheckoutController) SubmitCheckoutAction(ctx context.Context, r *web.Request) web.Result {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	return cc.showCheckoutFormAndHandleSubmit(ctx, r, "checkout/checkout")
}

// PlaceOrderAction functions as a return/notification URL for Payment Providers
func (cc *CheckoutController) PlaceOrderAction(ctx context.Context, r *web.Request) web.Result {
	session := web.SessionFromContext(ctx)
	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	guardRedirect := cc.getCommonGuardRedirects(ctx, r.Session(), decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.responder.Render("checkout/emptycart", nil).SetNoCache()
	}

	if !decoratedCart.Cart.PaymentSelection.IsSelected() {
		cc.logger.Error("cart.checkoutcontroller.submitaction: Error Gateway not in carts PaymentSelection")
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	gateway, err := cc.getPaymentGatewayAndSelectedMethod(ctx, decoratedCart.Cart.PaymentSelection.Gateway)
	if err != nil {
		cc.logger.Error("cart.checkoutcontroller.submitaction: Error %v", err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	cartPayment, err := gateway.GetFlowResult(ctx, &decoratedCart.Cart, session.ID())
	if err != nil {
		if cc.showReviewStepAfterPaymentError && !cc.skipReviewAction {
			return cc.showReviewFormWithErrors(ctx, *decoratedCart, err)
		}
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}
	err = gateway.ConfirmResult(ctx, &decoratedCart.Cart, cartPayment)
	if err != nil {
		if cc.showReviewStepAfterPaymentError && !cc.skipReviewAction {
			return cc.showReviewFormWithErrors(ctx, *decoratedCart, err)
		}
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}
	response, err := cc.placeOrder(ctx, r.Session(), *cartPayment, *decoratedCart)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
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

func (cc *CheckoutController) getPaymentReturnURL(PaymentProvider string) *url.URL {
	baseURL := cc.baseURL
	paymentURL, _ := cc.router.URL("checkout.placeorder", nil)

	rawURL := strings.TrimRight(baseURL, "/") + paymentURL.String()

	urlResult, _ := url.Parse(rawURL)

	return urlResult
}

func (cc *CheckoutController) getBasicViewData(ctx context.Context, session *web.Session, decoratedCart cart.DecoratedCart) CheckoutViewData {
	paymentGatewaysMethods := make(map[string][]paymentDomain.Method)
	for gatewayCode, gateway := range cc.webCartPaymentGatewayProvider {
		paymentGatewaysMethods[gatewayCode] = gateway.Methods()
	}
	return CheckoutViewData{
		DecoratedCart:        decoratedCart,
		CartValidationResult: cc.applicationCartService.ValidateCart(ctx, session, &decoratedCart),
		AvailablePayments:    paymentGatewaysMethods,
		CustomerLoggedIn:     cc.userService.IsLoggedIn(ctx, session),
	}
}

//showCheckoutFormAndHandleSubmit - Action that shows the form
func (cc *CheckoutController) showCheckoutFormAndHandleSubmit(ctx context.Context, r *web.Request, template string) web.Result {
	session := r.Session()

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, session)
	if e != nil {
		cc.logger.Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	if len(cc.webCartPaymentGatewayProvider) == 0 {
		cc.logger.Error("cart.checkoutcontroller.submitaction: Error No Payment set")
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	viewData := cc.getBasicViewData(ctx, session, *decoratedCart)
	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		if cc.showEmptyCartPageIfNoItems {
			return cc.responder.Render("checkout/emptycart", nil).SetNoCache()
		}
		return cc.responder.Render(template, viewData).SetNoCache()
	}

	if r.Request().Method != http.MethodPost {
		//Form not Submitted:
		form, err := cc.checkoutFormController.GetUnsubmittedForm(ctx, r)
		if err != nil {
			viewData.ErrorInfos = getViewErrorInfo(err)
			return cc.responder.Render(template, viewData).SetNoCache()
		}
		viewData.Form = *form
		return cc.responder.Render(template, viewData).SetNoCache()
	}

	//Form submitted:
	form, success, err := cc.checkoutFormController.HandleFormAction(ctx, r)
	if err != nil {
		viewData.ErrorInfos = getViewErrorInfo(err)
		return cc.responder.Render(template, viewData).SetNoCache()
	}
	viewData.Form = *form

	if success {
		cc.logger.Debug("submit checkout suceeded: redirect to checkout.review")
		if cc.skipReviewAction {
			return cc.processPaymentBeforePlaceOrder(ctx, r)
		}
		response := cc.responder.RouteRedirect("checkout.review", nil)
		response.SetNoCache()
		return response
	}

	//Default: show form with its validation result
	return cc.responder.Render(template, viewData).SetNoCache()
}

//showCheckoutFormWithErrors - error handling that is called from many places... It will show the checkoutform and the error
// template and form is optional - if it is not goven it is autodetected and prefilled from the infos in the cart
func (cc *CheckoutController) showCheckoutFormWithErrors(ctx context.Context, r *web.Request, decoratedCart cart.DecoratedCart, form *forms.CheckoutFormComposite, err error) web.Result {
	template := "checkout/checkout"

	cc.logger.Warn("showCheckoutFormWithErrors / Error: %s", err.Error())
	viewData := cc.getBasicViewData(ctx, r.Session(), decoratedCart)
	if form == nil {
		form, _ = cc.checkoutFormController.GetUnsubmittedForm(ctx, r)
	}
	if form == nil {
		viewData.Form = *form
	}
	viewData.ErrorInfos = getViewErrorInfo(err)
	return cc.responder.Render(template, viewData).SetNoCache()
}

//showReviewFormWithErrors
func (cc *CheckoutController) showReviewFormWithErrors(ctx context.Context, decoratedCart cart.DecoratedCart, err error) web.Result {
	cc.logger.Warn("Show Error (review step): %s", err.Error())
	viewData := ReviewStepViewData{
		DecoratedCart: decoratedCart,
		ErrorInfos:    getViewErrorInfo(err),
	}
	return cc.responder.Render("checkout/review", viewData).SetNoCache()
}

func getViewErrorInfo(err error) ViewErrorInfos {
	hasPaymentError := false

	if paymentErr, ok := err.(*paymentDomain.Error); ok {
		hasPaymentError = paymentErr.ErrorCode != paymentDomain.PaymentErrorCodeCancelled
	}

	errorInfos := ViewErrorInfos{
		HasError:        true,
		ErrorMessage:    err.Error(),
		HasPaymentError: hasPaymentError,
	}
	return errorInfos
}

func (cc *CheckoutController) processPaymentBeforePlaceOrder(ctx context.Context, r *web.Request) web.Result {
	session := web.SessionFromContext(ctx)
	//Guard Clause if Cart can not be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	gateway, err := cc.getPaymentGatewayAndSelectedMethod(ctx, decoratedCart.Cart.PaymentSelection.Gateway)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}

	//procces Payment:
	returnURL := cc.getPaymentReturnURL(decoratedCart.Cart.PaymentSelection.Gateway)

	//selected payment need to be set on cart before
	//Handover to selected gateway flow:
	flowResult, err := gateway.StartFlow(ctx, &decoratedCart.Cart, session.ID(), returnURL)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}

	return flowResult
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
		cc.logger.WithField("subcategory", "checkoutError").WithField("errorMsg", err.Error()).Error(fmt.Sprintf("place order failed: cart id: %v / total-amount: %v / cartPayment: %#v", decoratedCart.Cart.EntityID, decoratedCart.Cart.GrandTotal(), cartPayment))
		return nil, err
	}

	email := cc.getContactMail(decoratedCart.Cart)

	var placeOrderPaymentInfos []PlaceOrderPaymentInfo
	for _, transaction := range cartPayment.Transactions {
		placeOrderPaymentInfos = append(placeOrderPaymentInfos, PlaceOrderPaymentInfo{
			Gateway:         cartPayment.Gateway,
			Method:          transaction.Method,
			PaymentProvider: transaction.PaymentProvider,
			Title:           transaction.Title,
			Amount:          transaction.AmountPayed,
			CreditCardInfo:  transaction.CreditCardInfo,
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

// ReviewAction handles the cart review action
func (cc *CheckoutController) ReviewAction(ctx context.Context, r *web.Request) web.Result {
	if cc.skipReviewAction {
		return cc.responder.Render("checkout/carterror", nil)
	}

	//Guard Clause if cart can not be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCartWithoutCache(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

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

	//Everything valid then return
	if proceed == "1" && (!cc.privacyPolicyRequired || privacyPolicy == "1") && termsAndConditions == "1" && decoratedCart.Cart.PaymentSelection.IsSelected() {
		return cc.processPaymentBeforePlaceOrder(ctx, r)
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
