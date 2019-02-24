package controller

import (
	"context"
	"fmt"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CartAPIController for cart api
	CartAPIController struct {
		responder           *web.Responder
		cartService         *application.CartService
		cartReceiverService *application.CartReceiverService
		logger              flamingo.Logger
	}

	result struct {
		Message     string
		MessageCode string
		Success     bool
		CartTeaser  *cart.Teaser
	}

	messageCodeAvailable interface {
		MessageCode() string
	}
)

// Inject dependencies
func (cc *CartAPIController) Inject(
	responder *web.Responder,
	ApplicationCartService *application.CartService,
	ApplicationCartReceiverService *application.CartReceiverService,
	Logger flamingo.Logger,
) {
	cc.responder = responder
	cc.cartService = ApplicationCartService
	cc.cartReceiverService = ApplicationCartReceiverService
	cc.logger = Logger
}

// GetAction Get JSON Format of API
func (cc *CartAPIController) GetAction(ctx context.Context, r *web.Request) web.Result {
	cart, e := cc.cartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "CartApiController").Error("cart.cartapicontroller.get: %v", e.Error())
		response := cc.responder.Data(result{Message: e.Error(), Success: false})
		response.Status(500)
		return response
	}
	return cc.responder.Data(cart)
}

// UpdateDeliveryInfo updates the delivery info
func (cc *CartAPIController) UpdateDeliveryInfo(ctx context.Context, r *web.Request) web.Result {
	// todo: not yet implemented?
	//cc.cartService.UpdateDeliveryInfo(ctx,)
	return cc.responder.Data(result{})
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
	if err != nil {
		cc.logger.WithField("category", "CartApiController").Error("cart.cartapicontroller.add: %v", err.Error())
		msgCode := ""
		if e, ok := err.(messageCodeAvailable); ok {
			msgCode = e.MessageCode()
		}
		response := cc.responder.Data(result{Message: err.Error(), MessageCode: msgCode, Success: false})
		response.Status(500)
		return response
	}
	cart, err := cc.cartReceiverService.ViewCart(ctx, r.Session())
	if err != nil {
		cc.logger.WithField("category", "CartApiController").Error("cart.cartapicontroller.add: %v", err.Error())
		return cc.responder.Data(result{Message: err.Error(), Success: false}).Status(500)
	}
	return cc.responder.Data(result{
		Success:    true,
		Message:    fmt.Sprintf("added %v / %v Qty %v", addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty),
		CartTeaser: cart.GetCartTeaser(),
	})
}

// ApplyVoucherAndGetAction applies the given voucher and returns the cart
func (cc *CartAPIController) ApplyVoucherAndGetAction(ctx context.Context, r *web.Request) web.Result {
	couponCode := r.Params["couponCode"]

	cart, err := cc.cartService.ApplyVoucher(ctx, r.Session(), couponCode)
	if err != nil {
		response := cc.responder.Data(result{Message: err.Error(), Success: false})
		response.Status(500)
		return response
	}
	return cc.responder.Data(cart)
}

// CleanAndGetAction cleans the cart and returns the cleaned cart
func (cc *CartAPIController) CleanAndGetAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.cartService.DeleteAllItems(ctx, r.Session())
	if err != nil {
		response := cc.responder.Data(result{Message: err.Error(), Success: false})
		response.Status(500)
		return response
	}

	return cc.responder.RouteRedirect("cart.api.get", nil)
}

// CleanDeliveryAndGetAction cleans the given delivery from the cart and returns the cleaned cart
func (cc *CartAPIController) CleanDeliveryAndGetAction(ctx context.Context, r *web.Request) web.Result {
	deliveryCode := r.Params["deliveryCode"]
	cart, err := cc.cartService.DeleteDelivery(ctx, r.Session(), deliveryCode)
	if err != nil {
		response := cc.responder.Data(result{Message: err.Error(), Success: false})
		response.Status(500)
		return response
	}

	return cc.responder.Data(cart)
}
