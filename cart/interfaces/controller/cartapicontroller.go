package controller

import (
	"context"
	"net/url"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"

	formDomain "flamingo.me/form/domain"

	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// CartAPIController for cart api
	CartAPIController struct {
		responder                    *web.Responder
		cartService                  *application.CartService
		cartReceiverService          *application.CartReceiverService
		logger                       flamingo.Logger
		billingAddressFormController *forms.BillingAddressFormController
		deliveryFormController       *forms.DeliveryFormController
		simplePaymentFormController  *forms.SimplePaymentFormController
	}

	// CartAPIResult view data
	CartAPIResult struct {
		// Contains details if success is false
		Error                *resultError
		Success              bool
		CartTeaser           *cart.Teaser
		Data                 interface{}
		DataValidationInfo   *formDomain.ValidationInfo `swaggertype:"object"`
		CartValidationResult *validation.Result
	}

	getCartResult struct {
		Cart                 *cart.Cart
		CartValidationResult *validation.Result
	}

	resultError struct {
		Message string
		Code    string
	} // @name cartResultError

	messageCodeAvailable interface {
		MessageCode() string
	}

	// PromotionFunction type takes ctx, cart, couponCode and applies the promotion
	promotionFunc func(ctx context.Context, session *web.Session, couponCode string) (*cart.Cart, error)
)

// Inject dependencies
func (cc *CartAPIController) Inject(
	responder *web.Responder,
	ApplicationCartService *application.CartService,
	ApplicationCartReceiverService *application.CartReceiverService,
	billingAddressFormController *forms.BillingAddressFormController,
	deliveryFormController *forms.DeliveryFormController,
	simplePaymentFormController *forms.SimplePaymentFormController,
	Logger flamingo.Logger,
) {
	cc.responder = responder
	cc.cartService = ApplicationCartService
	cc.cartReceiverService = ApplicationCartReceiverService
	cc.logger = Logger.WithField("category", "CartApiController")
	cc.billingAddressFormController = billingAddressFormController
	cc.deliveryFormController = deliveryFormController
	cc.simplePaymentFormController = simplePaymentFormController
}

// GetAction Get JSON Format of API
// @Summary Get the current cart
// @Tags Cart
// @Produce json
// @Success 200 {object} getCartResult
// @Failure 500 {object} CartAPIResult
// @Router /api/v1/cart [get]
func (cc *CartAPIController) GetAction(ctx context.Context, r *web.Request) web.Result {
	decoratedCart, e := cc.cartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		result := newResult()
		result.SetError(e, "get_error")
		cc.logger.WithContext(ctx).Error("cart.cartapicontroller.get: %v", e.Error())
		return cc.responder.Data(result).Status(500)
	}
	validationResult := cc.cartService.ValidateCart(ctx, web.SessionFromContext(ctx), decoratedCart)
	return cc.responder.Data(getCartResult{
		CartValidationResult: &validationResult,
		Cart:                 &decoratedCart.Cart,
	})
}

// DeleteCartAction removes all cart content and returns a blank cart
// @Summary Remove all stored cart information e.g. items, deliveries, billing address and returns the empty cart.
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Router /api/v1/cart [delete]
func (cc *CartAPIController) DeleteCartAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.cartService.Clean(ctx, r.Session())

	result := newResult()
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.cartapicontroller.delete: %v", err.Error())

		result.SetError(err, "delete_cart_error")
		response := cc.responder.Data(result)
		response.Status(500)
		return response
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// AddAction Add Item to cart
// @Summary Add Item to cart
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param deliveryCode path string true "the identifier for the delivery in the cart"
// @Param marketplaceCode query string true "the product identifier that should be added"
// @Param variantMarketplaceCode query string false "optional the product identifier of the variant (for configurable products) that should be added"
// @Param qty query integer false "optional the qty that should be added"
// @Router /api/v1/cart/delivery/{deliveryCode}/item [post]
func (cc *CartAPIController) AddAction(ctx context.Context, r *web.Request) web.Result {
	variantMarketplaceCode := r.Params["variantMarketplaceCode"]

	qty, ok := r.Params["qty"]
	if !ok {
		qty = "1"
	}
	qtyInt, _ := strconv.Atoi(qty)
	deliveryCode := r.Params["deliveryCode"]

	addRequest := cc.cartService.BuildAddRequest(ctx, r.Params["marketplaceCode"], variantMarketplaceCode, qtyInt, nil)
	_, err := cc.cartService.AddProduct(ctx, r.Session(), deliveryCode, addRequest)

	result := newResult()
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.cartapicontroller.add: %v", err.Error())

		result.SetError(err, "add_product_error")
		response := cc.responder.Data(result)
		response.Status(500)
		return response
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// DeleteItemAction deletes an item from the cart
// @Summary Delete item from cart
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param deliveryCode path string true "the identifier for the delivery in the cart"
// @Param itemID query string true "the item that should be deleted"
// @Router /api/v1/cart/delivery/{deliveryCode}/item [delete]
func (cc *CartAPIController) DeleteItemAction(ctx context.Context, r *web.Request) web.Result {
	itemID, _ := r.Query1("itemID")
	deliveryCode := r.Params["deliveryCode"]

	err := cc.cartService.DeleteItem(ctx, r.Session(), itemID, deliveryCode)

	result := newResult()
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.cartapicontroller.delete: %v", err.Error())

		result.SetError(err, "delete_item_error")
		response := cc.responder.Data(result)
		response.Status(500)
		return response
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// UpdateItemAction updates the item qty in the current cart
// @Summary Update item in the cart
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param deliveryCode path string true "the identifier for the delivery in the cart"
// @Param itemID query string true "the item that should be updated"
// @Param qty query integer true "the new qty"
// @Router /api/v1/cart/delivery/{deliveryCode}/item [put]
func (cc *CartAPIController) UpdateItemAction(ctx context.Context, r *web.Request) web.Result {
	itemID, _ := r.Query1("itemID")
	deliveryCode := r.Params["deliveryCode"]
	qty, ok := r.Params["qty"]
	if !ok {
		qty = "1"
	}
	qtyInt, _ := strconv.Atoi(qty)

	err := cc.cartService.UpdateItemQty(ctx, r.Session(), itemID, deliveryCode, qtyInt)

	result := newResult()
	if err != nil {
		cc.logger.WithContext(ctx).Error("cart.cartapicontroller.updateItem: %v", err.Error())

		result.SetError(err, "update_item_error")
		response := cc.responder.Data(result)
		response.Status(500)
		return response
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// ApplyVoucherAndGetAction applies the given voucher and returns the cart
// @Summary Apply Voucher Code
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param couponCode query string true "the couponCode that should be applied"
// @Router /api/v1/cart/voucher [post]
func (cc *CartAPIController) ApplyVoucherAndGetAction(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "voucher_error", cc.cartService.ApplyVoucher)
}

// RemoveVoucherAndGetAction removes the given voucher and returns the cart
// @Summary Remove Voucher Code
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param couponCode query string true "the couponCode that should be applied"
// @Router /api/v1/cart/voucher [delete]
func (cc *CartAPIController) RemoveVoucherAndGetAction(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "voucher_error", cc.cartService.RemoveVoucher)
}

// DeleteAllItemsAction removes all cart items and returns the cart
// @Summary Remove all cart items from all deliveries and return the cart, keeps the delivery info untouched.
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Router /api/v1/cart/deliveries/items [delete]
func (cc *CartAPIController) DeleteAllItemsAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.cartService.DeleteAllItems(ctx, r.Session())
	result := newResult()
	if err != nil {
		result.SetError(err, "delete_items_error")
		response := cc.responder.Data(result)
		response.Status(500)
		return response
	}
	return cc.responder.Data(result)
}

// ApplyGiftCardAndGetAction applies the given gift card and returns the cart
// the request needs a query string param "couponCode" which includes the corresponding gift card code
// @Summary Apply Gift Card
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param couponCode query string true "the gift card code"
// @Router /api/v1/cart/gift-card [post]
func (cc *CartAPIController) ApplyGiftCardAndGetAction(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "giftcard_error", cc.cartService.ApplyGiftCard)
}

// ApplyCombinedVoucherGift applies a given code (which might be either a voucher or a Gift Card code) to the
// cartService and returns the cart
// @Summary Apply Gift Card or Voucher (auto detected)
// @Description Use this if you have one user input and that input can be used to either enter a voucher or a gift card
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param couponCode query string true "the couponCode that should be applied as gift card or voucher"
// @Router /api/v1/cart/voucher-gift-card [post]
func (cc *CartAPIController) ApplyCombinedVoucherGift(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "applyany_error", cc.cartService.ApplyAny)
}

// RemoveGiftCardAndGetAction removes the given gift card and returns the cart
// the request needs a query string param "couponCode" which includes the corresponding gift card code
// @Summary Remove Gift Card
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param couponCode query string true "the couponCode that should be deleted as gift card"
// @Router /api/v1/cart/gift-card [delete]
func (cc *CartAPIController) RemoveGiftCardAndGetAction(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "giftcard_error", cc.cartService.RemoveGiftCard)
}

// handles promotion action
func (cc *CartAPIController) handlePromotionAction(ctx context.Context, r *web.Request, errorCode string, fn promotionFunc) web.Result {
	couponCode := r.Params["couponCode"]
	result := newResult()
	_, err := fn(ctx, r.Session(), couponCode)
	if err != nil {
		cc.enrichResultWithCartInfos(ctx, &result)
		result.SetError(err, errorCode)
		response := cc.responder.Data(result)
		response.Status(500)

		return response
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// DeleteDelivery cleans the given delivery from the cart and returns the cleaned cart
// @Summary Cleans the given delivery from the cart
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param deliveryCode path string true "the identifier for the delivery in the cart"
// @Router /api/v1/cart/delivery/{deliveryCode} [delete]
func (cc *CartAPIController) DeleteDelivery(ctx context.Context, r *web.Request) web.Result {
	result := newResult()
	deliveryCode := r.Params["deliveryCode"]
	_, err := cc.cartService.DeleteDelivery(ctx, r.Session(), deliveryCode)
	if err != nil {
		result.SetError(err, "delete_delivery_error")
		return cc.responder.Data(result).Status(500)
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// BillingAction adds billing infos
// @Summary Adds billing infos to the current cart
// @Tags Cart
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param vat formData string false "vat"
// @Param firstname formData string true "firstname"
// @Param lastname formData string true "lastname"
// @Param middlename formData string false "middlename"
// @Param title formData string false "title"
// @Param salutation formData string false "salutation"
// @Param street formData string false "street"
// @Param streetNr formData string false "streetNr"
// @Param addressLine1 formData string false "addressLine1"
// @Param addressLine2 formData string false "addressLine2"
// @Param company formData string false "company"
// @Param postCode formData string false "postCode"
// @Param city formData string false "city"
// @Param state formData string false "state"
// @Param regionCode formData string false "regionCode"
// @Param country formData string false "country"
// @Param countryCode formData string false "countryCode"
// @Param phoneAreaCode formData string false "phoneAreaCode"
// @Param phoneCountryCode formData string false "phoneCountryCode"
// @Param phoneNumber formData string false "phoneNumber"
// @Param email formData string true "email"
// @Router /api/v1/cart/billing [put]
func (cc *CartAPIController) BillingAction(ctx context.Context, r *web.Request) web.Result {
	result := newResult()
	form, success, err := cc.billingAddressFormController.HandleFormAction(ctx, r)
	result.Success = success
	if err != nil {
		result.SetError(err, "form_error")
		return cc.responder.Data(result)
	}

	if form != nil {
		result.Data = form.Data
		result.DataValidationInfo = &form.ValidationInfo
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// UpdateDeliveryInfoAction updates the delivery info
// @Summary Adds delivery infos, such as shipping address to the delivery for the cart
// @Tags Cart
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param deliveryCode path string true "the identifier for the delivery in the cart"
// @Param deliveryAddress.vat formData string false "vat"
// @Param deliveryAddress.firstname formData string true "firstname"
// @Param deliveryAddress.lastname formData string true "lastname"
// @Param deliveryAddress.middlename formData string false "middlename"
// @Param deliveryAddress.title formData string false "title"
// @Param deliveryAddress.salutation formData string false "salutation"
// @Param deliveryAddress.street formData string false "street"
// @Param deliveryAddress.streetNr formData string false "streetNr"
// @Param deliveryAddress.addressLine1 formData string false "addressLine1"
// @Param deliveryAddress.addressLine2 formData string false "addressLine2"
// @Param deliveryAddress.company formData string false "company"
// @Param deliveryAddress.postCode formData string false "postCode"
// @Param deliveryAddress.city formData string false "city"
// @Param deliveryAddress.state formData string false "state"
// @Param deliveryAddress.regionCode formData string false "regionCode"
// @Param deliveryAddress.country formData string false "country"
// @Param deliveryAddress.countryCode formData string false "countryCode"
// @Param deliveryAddress.phoneAreaCode formData string false "phoneAreaCode"
// @Param deliveryAddress.phoneCountryCode formData string false "phoneCountryCode"
// @Param deliveryAddress.phoneNumber formData string false "phoneNumber"
// @Param deliveryAddress.email formData string true "email"
// @Param useBillingAddress formData bool false "useBillingAddress"
// @Param shippingMethod formData string false "shippingMethod"
// @Param shippingCarrier formData string false "shippingCarrier"
// @Param locationCode formData string false "locationCode"
// @Param desiredTime formData string false "desired date/time in RFC3339" format(date-time)
// @Router /api/v1/cart/delivery/{deliveryCode} [put]
func (cc *CartAPIController) UpdateDeliveryInfoAction(ctx context.Context, r *web.Request) web.Result {
	result := newResult()
	form, success, err := cc.deliveryFormController.HandleFormAction(ctx, r)
	result.Success = success
	if err != nil {
		result.SetError(err, "form_error")
		return cc.responder.Data(result)
	}
	if form != nil {
		result.Data = form.Data
		result.DataValidationInfo = &form.ValidationInfo
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

// UpdatePaymentSelectionAction to set / update the cart payment selection
// @Summary Update/set the PaymentSelection for the current cart
// @Tags Cart
// @Produce json
// @Success 200 {object} CartAPIResult
// @Failure 500 {object} CartAPIResult
// @Param gateway query string true "name of the payment gateway - e.g. 'offline'"
// @Param method query string true "name of the payment method - e.g. 'offlinepayment_cashondelivery'"
// @Router /api/v1/cart/payment-selection [put]
func (cc *CartAPIController) UpdatePaymentSelectionAction(ctx context.Context, r *web.Request) web.Result {
	result := newResult()
	gateway, _ := r.Query1("gateway")
	method, _ := r.Query1("method")

	urlValues := make(url.Values)
	urlValues["gateway"] = []string{gateway}
	urlValues["method"] = []string{method}
	newRequest := web.CreateRequest(web.RequestFromContext(ctx).Request(), web.SessionFromContext(ctx))
	newRequest.Request().Form = urlValues

	form, success, err := cc.simplePaymentFormController.HandleFormAction(ctx, newRequest)
	result.Success = success
	if err != nil {
		result.SetError(err, "form_error")
		response := cc.responder.Data(result)
		response.Status(500)
		return response
	}
	if form != nil {
		result.Data = form.Data
		result.DataValidationInfo = &form.ValidationInfo
	}
	cc.enrichResultWithCartInfos(ctx, &result)
	return cc.responder.Data(result)
}

func (cc *CartAPIController) enrichResultWithCartInfos(ctx context.Context, result *CartAPIResult) {
	session := web.SessionFromContext(ctx)
	decoratedCart, err := cc.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		result.SetError(err, "view_cart_error")

	}
	validationResult := cc.cartService.ValidateCart(ctx, session, decoratedCart)
	result.CartTeaser = decoratedCart.Cart.GetCartTeaser()
	result.CartValidationResult = &validationResult
}

// newResult factory to get new CartApiResult (with success true)
func newResult() CartAPIResult {
	return CartAPIResult{
		Success: true,
	}
}

// SetErrorByCode sets the error on the CartApiResult data and success to false
func (r *CartAPIResult) SetErrorByCode(message string, code string) *CartAPIResult {
	r.Success = false
	r.Error = &resultError{
		Message: message,
		Code:    code,
	}
	return r
}

// SetError updates the cart error field
func (r *CartAPIResult) SetError(err error, fallbackCode string) *CartAPIResult {
	if e, ok := err.(messageCodeAvailable); ok {
		fallbackCode = e.MessageCode()
	}
	return r.SetErrorByCode(err.Error(), fallbackCode)
}
