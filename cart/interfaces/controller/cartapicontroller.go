package controller

import (
	"context"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"

	formDomain "flamingo.me/form/domain"

	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
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
	}

	// CartAPIResult view data
	CartAPIResult struct {
		//Contains details if success is false
		Error                *resultError
		Success              bool
		CartTeaser           *cart.Teaser
		Data                 interface{}
		DataValidationInfo   *formDomain.ValidationInfo
		CartValidationResult *validation.Result
	}

	getCartResult struct {
		Cart                 *cart.Cart
		CartValidationResult *validation.Result
	}

	resultError struct {
		Message string
		Code    string
	}

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
	Logger flamingo.Logger,
) {
	cc.responder = responder
	cc.cartService = ApplicationCartService
	cc.cartReceiverService = ApplicationCartReceiverService
	cc.logger = Logger.WithField("category", "CartApiController")
	cc.billingAddressFormController = billingAddressFormController
	cc.deliveryFormController = deliveryFormController
}

// GetAction Get JSON Format of API
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

// AddAction Add Item to cart
func (cc *CartAPIController) AddAction(ctx context.Context, r *web.Request) web.Result {
	variantMarketplaceCode, _ := r.Params["variantMarketplaceCode"]

	qty, ok := r.Params["qty"]
	if !ok {
		qty = "1"
	}
	qtyInt, _ := strconv.Atoi(qty)
	deliveryCode, _ := r.Params["deliveryCode"]

	addRequest := cc.cartService.BuildAddRequest(ctx, r.Params["marketplaceCode"], variantMarketplaceCode, qtyInt)
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

// ApplyVoucherAndGetAction applies the given voucher and returns the cart
func (cc *CartAPIController) ApplyVoucherAndGetAction(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "voucher_error", cc.cartService.ApplyVoucher)
}

// RemoveVoucherAndGetAction removes the given voucher and returns the cart
func (cc *CartAPIController) RemoveVoucherAndGetAction(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "voucher_error", cc.cartService.RemoveVoucher)
}

// DeleteCartAction cleans the cart and returns the cleaned cart
func (cc *CartAPIController) DeleteCartAction(ctx context.Context, r *web.Request) web.Result {
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

// ApplyGiftCardAndGetAction applies the given giftcard and returns the cart
// the request needs a query string param "couponCode" which includes the corresponding giftcard code
func (cc *CartAPIController) ApplyGiftCardAndGetAction(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "giftcard_error", cc.cartService.ApplyGiftCard)
}

// ApplyCombinedVoucherGift applies a given code (which might be either a voucher or a giftcard code) to the
// cartService and returns the cart
func (cc *CartAPIController) ApplyCombinedVoucherGift(ctx context.Context, r *web.Request) web.Result {
	return cc.handlePromotionAction(ctx, r, "applyany_error", cc.cartService.ApplyAny)
}

// RemoveGiftCardAndGetAction removes the given giftcard and returns the cart
// the request needs a query string param "couponCode" which includes the corresponding giftcard code
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

// BillingAction handles the checkout start action
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

//newResult - factory to get new CartApiResult (with success true)
func newResult() CartAPIResult {
	return CartAPIResult{
		Success: true,
	}
}

//SetErrorByCode - sets the error on the CartApiResult data and success to false
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
