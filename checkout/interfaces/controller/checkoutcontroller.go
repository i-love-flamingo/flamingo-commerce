package controller

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller/forms"
	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

type (
	// CheckoutViewData represents the checkout view data
	CheckoutViewData struct {
		DecoratedCart        decorator.DecoratedCart
		Form                 forms.CheckoutFormComposite
		CartValidationResult validation.Result
		ErrorInfos           ViewErrorInfos
		AvailablePayments    map[string][]paymentDomain.Method
		CustomerLoggedIn     bool
	}

	// ViewErrorInfos defines the error info struct of the checkout controller views
	ViewErrorInfos struct {
		//HasError  indicates that an general error happened
		HasError bool
		//If there is a general error this field is filled and can be used in the template
		ErrorMessage string
		//if the Error happens during processing payment (can be used in template to behave special in case of payment errors)
		HasPaymentError bool
	}

	// SuccessViewData represents the success view data
	SuccessViewData struct {
		PaymentInfos        []application.PlaceOrderPaymentInfo
		PlacedOrderInfos    placeorder.PlacedOrderInfos
		Email               string
		PlacedDecoratedCart decorator.DecoratedCart
	}

	// ReviewStepViewData represents the success view data
	ReviewStepViewData struct {
		DecoratedCart decorator.DecoratedCart
		ErrorInfos    ViewErrorInfos
	}

	// PaymentStepViewData represents the payment flow view data
	PaymentStepViewData struct {
		FlowStatus paymentDomain.FlowStatus
		ErrorInfos ViewErrorInfos
	}

	// PlaceOrderFlashData represents the data passed to the success page - they need to be "glob"able
	PlaceOrderFlashData struct {
		PlacedOrderInfos placeorder.PlacedOrderInfos
		Email            string
		PaymentInfos     []application.PlaceOrderPaymentInfo
		PlacedCart       cart.Cart
	}

	// EmptyCartInfo struct defines the data info on empty carts
	EmptyCartInfo struct {
		CartExpired bool
	}

	// CheckoutController represents the checkout controller with its injections
	CheckoutController struct {
		responder *web.Responder
		router    *web.Router

		orderService         *application.OrderService
		decoratedCartFactory *decorator.DecoratedCartFactory

		skipStartAction                 bool
		skipReviewAction                bool
		showReviewStepAfterPaymentError bool
		showEmptyCartPageIfNoItems      bool
		redirectToCartOnInvalideCart    bool
		privacyPolicyRequired           bool

		devMode bool

		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService

		userService *authApplication.UserService

		logger flamingo.Logger

		checkoutFormController *forms.CheckoutFormController
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
	decoratedCartFactory *decorator.DecoratedCartFactory,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	userService *authApplication.UserService,
	logger flamingo.Logger,
	checkoutFormController *forms.CheckoutFormController,
	config *struct {
		SkipStartAction                 bool `inject:"config:checkout.skipStartAction,optional"`
		SkipReviewAction                bool `inject:"config:checkout.skipReviewAction,optional"`
		ShowReviewStepAfterPaymentError bool `inject:"config:checkout.showReviewStepAfterPaymentError,optional"`
		ShowEmptyCartPageIfNoItems      bool `inject:"config:checkout.showEmptyCartPageIfNoItems,optional"`
		RedirectToCartOnInvalideCart    bool `inject:"config:checkout.redirectToCartOnInvalideCart,optional"`
		PrivacyPolicyRequired           bool `inject:"config:checkout.privacyPolicyRequired,optional"`
		DevMode                         bool `inject:"config:debug.mode,optional"`
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

	cc.devMode = config.DevMode

	cc.applicationCartService = applicationCartService
	cc.applicationCartReceiverService = applicationCartReceiverService

	cc.userService = userService

	cc.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "checkoutController")
}

/*
The checkoutController implements a default process for a checkout:
 * StartAction (supposed to show a switch to go to guest or customer)
 	* can be skipped with a configuration
 * SubmitCheckoutAction
 	* This step is supposed to show a big form (validation and default values are configurable as well)
	* payment can be selected in this step or in the next
	* In cases a customer is logged in the form is pre populated
 * ReviewAction (can be skipped through configuration)
	* this step is supposed to show the current cart status just before checkout
		* optional the paymentmethod can also be selected here
 * PaymentAction
	* Payment gets initialized in SubmitCheckoutAction or ReviewAction
	* Handles different payment stages and reacts to payment status provided by gateway
 * PlaceOrderAction
	* Place the order if not already placed
	* Add Information about the order to the session flash messages
 * SuccessStep
	* Displays order success page containing the order infos from the previously set flash message
*/

// StartAction handles the checkout start action
func (cc *CheckoutController) StartAction(ctx context.Context, r *web.Request) web.Result {
	//Guard Clause if Cart cannot be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.viewaction: Error ", e)
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

// SubmitCheckoutAction handles the main checkout
func (cc *CheckoutController) SubmitCheckoutAction(ctx context.Context, r *web.Request) web.Result {
	//Guard Clause if Cart can not be fetched
	decoratedCart, e := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error ", e)
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

	decoratedCart, err := cc.orderService.LastPlacedOrCurrentCart(ctx)
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.placeorderaction: Error ", err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	guardRedirect := cc.getCommonGuardRedirects(ctx, session, decoratedCart)
	if guardRedirect != nil {
		return guardRedirect
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.responder.Render("checkout/emptycart", nil).SetNoCache()
	}
	return cc.placeOrderAction(ctx, r, session, decoratedCart)
}

func (cc *CheckoutController) placeOrderAction(ctx context.Context, r *web.Request, session *web.Session, decoratedCart *decorator.DecoratedCart) web.Result {

	err := cc.orderService.SetSources(ctx, session)
	if err != nil {
		cc.logger.Error("OnStepCurrentCartPlaceOrder SetSources Error ", err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}
	var placedOrderInfo *application.PlaceOrderInfo
	if cc.orderService.HasLastPlacedOrder(ctx) {
		placedOrderInfo, _ = cc.orderService.LastPlacedOrder(ctx)
		cc.orderService.ClearLastPlacedOrder(ctx)
	} else {
		placedOrderInfo, err = cc.orderService.CurrentCartPlaceOrderWithPaymentProcessing(ctx, session)

		if err != nil {
			cc.logger.WithContext(ctx).WithField("subcategory", "checkoutError").WithField("errorMsg", err.Error()).Error(fmt.Sprintf("place order failed: cart id: %v / total-amount: %v", decoratedCart.Cart.EntityID, decoratedCart.Cart.GrandTotal()))
			if paymentError, ok := err.(*paymentDomain.Error); ok {
				if cc.showReviewStepAfterPaymentError && !cc.skipReviewAction {
					return cc.showReviewFormWithErrors(ctx, *decoratedCart, paymentError)
				}
			}
			return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
		}
	}

	r.Session().AddFlash(PlaceOrderFlashData{
		PlacedOrderInfos: placedOrderInfo.PlacedOrders,
		Email:            placedOrderInfo.ContactEmail,
		PlacedCart:       decoratedCart.Cart,
		PaymentInfos:     placedOrderInfo.PaymentInfos,
	}, "checkout.success.data")
	return cc.responder.RouteRedirect("checkout.success", nil)
}

// SuccessAction handles the order success action
func (cc *CheckoutController) SuccessAction(ctx context.Context, r *web.Request) web.Result {
	flashes := r.Session().Flashes("checkout.success.data")
	if len(flashes) > 0 {

		// if in development mode, then restore the last order in flash session.
		if cc.devMode {
			r.Session().AddFlash(flashes[len(flashes)-1], "checkout.success.data")
		}

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
	return cc.responder.RouteRedirect("checkout.expired", nil).SetNoCache()
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

func (cc *CheckoutController) getPaymentReturnURL(r *web.Request) *url.URL {
	paymentURL, _ := cc.router.Absolute(r, "checkout.payment", nil)
	return paymentURL
}

func (cc *CheckoutController) getBasicViewData(ctx context.Context, session *web.Session, decoratedCart decorator.DecoratedCart) CheckoutViewData {
	paymentGatewaysMethods := make(map[string][]paymentDomain.Method)
	for gatewayCode, gateway := range cc.orderService.GetAvailablePaymentGateways(ctx) {
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
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error ", e)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	if len(cc.orderService.GetAvailablePaymentGateways(ctx)) == 0 {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error No Payment set")
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
		if cc.skipReviewAction {
			canProceed, err := cc.checkTermsAndPrivacyPolicy(r)
			if !canProceed || err != nil {
				viewData.ErrorInfos = getViewErrorInfo(err)

				return cc.responder.Render(template, viewData).SetNoCache()
			}

			cc.logger.WithContext(ctx).Debug("submit checkout succeeded: redirect to checkout.review")
			return cc.processPayment(ctx, r)
		}
		response := cc.responder.RouteRedirect("checkout.review", nil).SetNoCache()

		return response
	}

	//Default: show form with its validation result
	return cc.responder.Render(template, viewData).SetNoCache()
}

// showCheckoutFormWithErrors - error handling that is called from many places... It will show the checkoutform and the error
// template and form is optional - if it is not given it is autodetected and prefilled from the infos in the cart
func (cc *CheckoutController) showCheckoutFormWithErrors(ctx context.Context, r *web.Request, decoratedCart decorator.DecoratedCart, form *forms.CheckoutFormComposite, err error) web.Result {
	template := "checkout/checkout"

	cc.logger.WithContext(ctx).Warn("showCheckoutFormWithErrors / Error:", err)
	viewData := cc.getBasicViewData(ctx, r.Session(), decoratedCart)
	if form == nil {
		form, _ = cc.checkoutFormController.GetUnsubmittedForm(ctx, r)
	}
	if form != nil {
		viewData.Form = *form
	}
	viewData.ErrorInfos = getViewErrorInfo(err)
	return cc.responder.Render(template, viewData).SetNoCache()
}

// showReviewFormWithErrors
func (cc *CheckoutController) showReviewFormWithErrors(ctx context.Context, decoratedCart decorator.DecoratedCart, err error) web.Result {
	cc.logger.WithContext(ctx).Warn("Show Error (review step):", err)
	viewData := ReviewStepViewData{
		DecoratedCart: decoratedCart,
		ErrorInfos:    getViewErrorInfo(err),
	}
	return cc.responder.Render("checkout/review", viewData).SetNoCache()
}

func getViewErrorInfo(err error) ViewErrorInfos {
	if err == nil {
		return ViewErrorInfos{
			HasError:        true,
			HasPaymentError: false,
		}
	}

	hasPaymentError := false

	if paymentErr, ok := err.(*paymentDomain.Error); ok {
		hasPaymentError = true

		if paymentErr.ErrorCode == paymentDomain.PaymentErrorAbortedByCustomer {
			// in case of customer payment abort don't show error message in frontend
			return ViewErrorInfos{
				HasError:        false,
				HasPaymentError: false,
			}
		}
	}

	return ViewErrorInfos{
		HasError:        true,
		ErrorMessage:    err.Error(),
		HasPaymentError: hasPaymentError,
	}
}

func (cc *CheckoutController) processPayment(ctx context.Context, r *web.Request) web.Result {
	session := web.SessionFromContext(ctx)

	// reserve an unique order id for later order placing
	_, err := cc.applicationCartService.ReserveOrderIDAndSave(ctx, session)
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error ", err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	// guard clause if cart can not be fetched
	decoratedCart, err := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error ", err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	// get the payment gateway for the specified payment selection
	gateway, err := cc.orderService.GetPaymentGateway(ctx, decoratedCart.Cart.PaymentSelection.Gateway())
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}

	returnURL := cc.getPaymentReturnURL(r)

	// start the payment flow
	flowResult, err := gateway.StartFlow(ctx, &decoratedCart.Cart, application.PaymentFlowStandardCorrelationID, returnURL)
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}

	// payment flow requires an early place order
	if flowResult.EarlyPlaceOrder {
		payment, err := gateway.OrderPaymentFromFlow(ctx, &decoratedCart.Cart, application.PaymentFlowStandardCorrelationID)
		if err != nil {
			return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
		}

		err = cc.orderService.SetSources(ctx, session)
		if err != nil {
			return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
		}

		_, err = cc.orderService.CurrentCartPlaceOrder(ctx, session, *payment)
		if err != nil {
			return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
		}
	}

	return cc.responder.RouteRedirect("checkout.payment", nil)
}

// ReviewAction handles the cart review action
func (cc *CheckoutController) ReviewAction(ctx context.Context, r *web.Request) web.Result {
	if cc.skipReviewAction {
		return cc.responder.Render("checkout/carterror", nil)
	}

	// payment / order already in process redirect to payment status page
	if cc.orderService.HasLastPlacedOrder(ctx) {
		return cc.responder.RouteRedirect("checkout.payment", nil)
	}

	//Guard Clause if cart can not be fetched
	decoratedCart, err := cc.applicationCartReceiverService.ViewDecoratedCartWithoutCache(ctx, r.Session())
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error ", err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

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

	// check for terms and conditions and privacy policy
	canProceed, err := cc.checkTermsAndPrivacyPolicy(r)
	if err != nil {
		viewData.ErrorInfos = getViewErrorInfo(err)
	}

	//Everything valid then return
	if canProceed && err == nil && decoratedCart.Cart.IsPaymentSelected() {
		return cc.processPayment(ctx, r)
	}

	return cc.responder.Render("checkout/review", viewData).SetNoCache()

}

// PaymentAction asks the payment adapter about the current payment status and handle it
func (cc *CheckoutController) PaymentAction(ctx context.Context, r *web.Request) web.Result {
	session := web.SessionFromContext(ctx)

	decoratedCart, err := cc.orderService.LastPlacedOrCurrentCart(ctx)
	if err != nil {
		cc.logger.WithContext(ctx).Error(err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	if decoratedCart.Cart.PaymentSelection == nil {
		cc.logger.WithContext(ctx).Info("No PaymentSelection for cart with ID ", decoratedCart.Cart.ID)
		return cc.responder.RouteRedirect("checkout.expired", nil).SetNoCache()
	}

	gateway, err := cc.orderService.GetPaymentGateway(ctx, decoratedCart.Cart.PaymentSelection.Gateway())
	if err != nil {
		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}

	flowStatus, err := gateway.FlowStatus(ctx, &decoratedCart.Cart, application.PaymentFlowStandardCorrelationID)
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.paymentaction: Error ", err)

		// payment failed, reopen the cart to make it still usable.
		if cc.orderService.HasLastPlacedOrder(ctx) {
			infos, err := cc.orderService.LastPlacedOrder(ctx)
			if err != nil {
				cc.logger.WithContext(ctx).Error(err)
				return cc.responder.RouteRedirect("checkout", nil)
			}

			restoredCart, err := cc.orderService.CancelOrder(ctx, session, infos)
			if err != nil {
				cc.logger.WithContext(ctx).Error(err)
				return cc.responder.RouteRedirect("checkout", nil)
			}

			decoratedCart = cc.decoratedCartFactory.Create(ctx, *restoredCart)
			cc.orderService.ClearLastPlacedOrder(ctx)
		}

		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, err)
	}

	viewData := &PaymentStepViewData{
		FlowStatus: *flowStatus,
	}

	switch flowStatus.Status {
	case paymentDomain.PaymentFlowStatusUnapproved:
		// payment just started render payment page which handles actions
		return cc.responder.Render("checkout/payment", viewData).SetNoCache()
	case paymentDomain.PaymentFlowStatusApproved:
		// payment is done but not confirmed by customer, confirm it and place order
		orderPayment, err := gateway.OrderPaymentFromFlow(ctx, &decoratedCart.Cart, application.PaymentFlowStandardCorrelationID)
		if err != nil {
			viewData.ErrorInfos = getViewErrorInfo(err)
			cc.logger.WithContext(ctx).Error(err)
			return cc.responder.Render("checkout/payment", viewData).SetNoCache()
		}

		err = gateway.ConfirmResult(ctx, &decoratedCart.Cart, orderPayment)
		if err != nil {
			viewData.ErrorInfos = getViewErrorInfo(err)
			cc.logger.WithContext(ctx).Error(err)
			return cc.responder.Render("checkout/payment", viewData).SetNoCache()
		}

		return cc.responder.RouteRedirect("checkout.placeorder", nil)
	case paymentDomain.PaymentFlowStatusCompleted:
		// payment is done and confirmed, place order
		return cc.responder.RouteRedirect("checkout.placeorder", nil)
	case paymentDomain.PaymentFlowStatusAborted:
		// mark payment selection as new payment to allow the user to retry
		newPaymentSelection := decoratedCart.Cart.PaymentSelection.GenerateNewIdempotencyKey()
		err := cc.applicationCartService.UpdatePaymentSelection(ctx, session, newPaymentSelection)
		if err != nil {
			cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.paymentaction: Error during UpdatePaymentSelection:", err)
		}

		// payment was aborted by user, redirect to checkout so a new payment can be started
		if cc.orderService.HasLastPlacedOrder(ctx) {
			infos, err := cc.orderService.LastPlacedOrder(ctx)
			if err != nil {
				viewData.ErrorInfos = getViewErrorInfo(err)
				cc.logger.WithContext(ctx).Error(err)
				return cc.responder.Render("checkout/payment", viewData).SetNoCache()
			}

			// ignore restored cart here since it gets fetched newly in checkout
			_, err = cc.orderService.CancelOrder(ctx, session, infos)
			if err != nil {
				viewData.ErrorInfos = getViewErrorInfo(err)
				cc.logger.WithContext(ctx).Error(err)
				return cc.responder.Render("checkout/payment", viewData).SetNoCache()
			}

			cc.orderService.ClearLastPlacedOrder(ctx)
		}

		if cc.showReviewStepAfterPaymentError && !cc.skipReviewAction {
			return cc.responder.RouteRedirect("checkout.review", nil)
		}

		return cc.responder.RouteRedirect("checkout", nil)
	case paymentDomain.PaymentFlowStatusFailed, paymentDomain.PaymentFlowStatusCancelled:
		// mark payment selection as new payment to allow the user to retry
		newPaymentSelection := decoratedCart.Cart.PaymentSelection.GenerateNewIdempotencyKey()
		err := cc.applicationCartService.UpdatePaymentSelection(ctx, session, newPaymentSelection)
		if err != nil {
			cc.logger.WithContext(ctx).Error("cart.checkoutcontroller.paymentaction: Error during UpdatePaymentSelection:", err)
		}

		// payment failed or is cancelled by payment provider, redirect back to checkout
		if cc.orderService.HasLastPlacedOrder(ctx) {
			infos, err := cc.orderService.LastPlacedOrder(ctx)
			if err != nil {
				viewData.ErrorInfos = getViewErrorInfo(err)
				cc.logger.WithContext(ctx).Error(err)
				return cc.responder.Render("checkout/payment", viewData).SetNoCache()
			}

			restoredCart, err := cc.orderService.CancelOrder(ctx, session, infos)
			if err != nil {
				viewData.ErrorInfos = getViewErrorInfo(err)
				cc.logger.WithContext(ctx).Error(err)
				return cc.responder.Render("checkout/payment", viewData).SetNoCache()
			}

			decoratedCart = cc.decoratedCartFactory.Create(ctx, *restoredCart)
			cc.orderService.ClearLastPlacedOrder(ctx)
		}

		if flowStatus.Error == nil {
			flowStatus.Error = &paymentDomain.Error{}
		}

		if cc.showReviewStepAfterPaymentError && !cc.skipReviewAction {
			return cc.showReviewFormWithErrors(ctx, *decoratedCart, flowStatus.Error)
		}

		return cc.showCheckoutFormWithErrors(ctx, r, *decoratedCart, nil, flowStatus.Error)
	case paymentDomain.PaymentFlowWaitingForCustomer:
		// payment pending, waiting for customer
		return cc.responder.Render("checkout/payment", viewData).SetNoCache()
	default:
		// show payment page which can react to unknown payment status
		return cc.responder.Render("checkout/payment", viewData).SetNoCache()
	}
}

// getCommonGuardRedirects checks config and may return a redirect that should be executed before the common checkou actions
func (cc *CheckoutController) getCommonGuardRedirects(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart) web.Result {
	if cc.redirectToCartOnInvalideCart {
		result := cc.applicationCartService.ValidateCart(ctx, session, decoratedCart)
		if !result.IsValid() {
			cc.logger.WithContext(ctx).Info("StartAction > RedirectToCartOnInvalideCart")
			resp := cc.responder.RouteRedirect("cart.view", nil)
			resp.SetNoCache()
			return resp
		}
	}
	return nil
}

// checkTermsAndPrivacyPolicy checks if TermsAndConditions and PrivacyPolicy is set as required
// the returned error indicates that the check failed
func (cc *CheckoutController) checkTermsAndPrivacyPolicy(r *web.Request) (bool, error) {
	proceed, _ := r.Form1("proceed")
	termsAndConditions, _ := r.Form1("termsAndConditions")
	privacyPolicy, _ := r.Form1("privacyPolicy")

	// prepare a minimal slice for error messages
	errorMessages := make([]string, 0, 2)

	// check for privacy policy if required
	if cc.privacyPolicyRequired && privacyPolicy != "1" && proceed == "1" {
		errorMessages = append(errorMessages, "privacy_policy_required")
	}

	// check for terms and conditions if required
	if termsAndConditions != "1" {
		errorMessages = append(errorMessages, "terms_and_conditions_required")
	}

	canProceed := (proceed == "1" && (!cc.privacyPolicyRequired || privacyPolicy == "1") && termsAndConditions == "1")

	if 0 == len(errorMessages) {
		return canProceed, nil
	}

	return canProceed, errors.New(strings.Join(errorMessages, ", "))
}
