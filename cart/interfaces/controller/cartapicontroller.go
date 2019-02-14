package controller

import (
	"context"
	"fmt"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CartApiController for cart api
	CartApiController struct {
		Responder *web.Responder
		responder.RedirectAware
		cartService         *application.CartService
		cartReceiverService *application.CartReceiverService
		logger              flamingo.Logger
	}

	result struct {
		Message     string
		MessageCode string
		Success     bool
	}

	messageCodeAvailable interface {
		MessageCode() string
	}
)

// Inject dependencies
func (cc *CartApiController) Inject(
	jsonAware Responder *web.Responder,
	redirectAware responder.RedirectAware,
	ApplicationCartService *application.CartService,
	ApplicationCartReceiverService *application.CartReceiverService,
	Logger flamingo.Logger,
) {
	cc.JSONAware = jsonAware
	cc.RedirectAware = redirectAware
	cc.cartService = ApplicationCartService
	cc.cartReceiverService = ApplicationCartReceiverService
	cc.logger = Logger
}

// GetAction Get JSON Format of API
func (cc *CartApiController) GetAction(ctx context.Context, r *web.Request) web.Result {
	cart, e := cc.cartReceiverService.ViewDecoratedCart(ctx, r.Session().G())
	if e != nil {
		cc.logger.WithField("category", "CartApiController").Error("cart.cartapicontroller.get: %v", e.Error())
		return cc.JSONError(result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(cart)
}

// AddAction Add Item to cart
func (cc *CartApiController) AddAction(ctx context.Context, r *web.Request) web.Result {
	variantMarketplaceCode, _ := r.Params["variantMarketplaceCode"]

	qty, ok := r.Params["qty"]
	if !ok {
		qty = "1"
	}
	qtyInt, _ := strconv.Atoi(qty)
	deliveryCode, _ := r.Params["deliveryCode"]

	addRequest := cc.cartService.BuildAddRequest(ctx, r.MustParam1("marketplaceCode"), variantMarketplaceCode, qtyInt)
	_, err := cc.cartService.AddProduct(ctx, r.Session().G(), deliveryCode, addRequest)
	if err != nil {
		cc.logger.WithField("category", "CartApiController").Error("cart.cartapicontroller.add: %v", err.Error())
		msgCode := ""
		if e, ok := err.(messageCodeAvailable); ok {
			msgCode = e.MessageCode()
		}
		return cc.JSONError(result{Message: err.Error(), MessageCode: msgCode, Success: false}, 500)
	}
	return cc.JSON(result{
		Success: true,
		Message: fmt.Sprintf("added %v / %v Qty %v", addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty),
	})
}

// ApplyVoucherAndGetAction applies the given voucher and returns the cart
func (cc *CartApiController) ApplyVoucherAndGetAction(ctx context.Context, r *web.Request) web.Result {
	couponCode := r.MustParam1("couponCode")

	cart, err := cc.cartService.ApplyVoucher(ctx, r.Session().G(), couponCode)
	if err != nil {
		return cc.JSONError(result{Message: err.Error(), Success: false}, 500)
	}
	return cc.JSON(cart)
}

// CleanAndGetAction cleans the cart and returns the cleaned cart
func (cc *CartApiController) CleanAndGetAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.cartService.DeleteAllItems(ctx, r.Session().G())
	if err != nil {
		return cc.JSONError(result{Message: err.Error(), Success: false}, 500)
	}

	return cc.Redirect("cart.api.get", nil)
}

// CleanDeliveryAndGetAction cleans the given delivery from the cart and returns the cleaned cart
func (cc *CartApiController) CleanDeliveryAndGetAction(ctx context.Context, r *web.Request) web.Result {
	deliveryCode := r.MustParam1("deliveryCode")
	cart, err := cc.cartService.DeleteDelivery(ctx, r.Session().G(), deliveryCode)
	if err != nil {
		return cc.JSONError(result{Message: err.Error(), Success: false}, 500)
	}

	return cc.JSON(cart)
}
