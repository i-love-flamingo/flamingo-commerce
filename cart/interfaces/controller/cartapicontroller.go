package controller

import (
	"context"
	"fmt"
	"strconv"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	// CartApiController for cart api
	CartApiController struct {
		responder.JSONAware
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

func (cc *CartApiController) Inject(
	aware responder.JSONAware,
	ApplicationCartService *application.CartService,
	ApplicationCartReceiverService *application.CartReceiverService,
	Logger flamingo.Logger,
) {
	cc.JSONAware = aware
	cc.cartService = ApplicationCartService
	cc.cartReceiverService = ApplicationCartReceiverService
	cc.logger = Logger
}

// GetAction Get JSON Format of API
func (cc *CartApiController) GetAction(ctx context.Context, r *web.Request) web.Response {
	cart, e := cc.cartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if e != nil {
		cc.logger.WithField("category", "CartApiController").Error("cart.cartapicontroller.get: %v", e.Error())
		return cc.JSONError(result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(cart)
}

// AddAction Add Item to cart
func (cc *CartApiController) AddAction(ctx context.Context, r *web.Request) web.Response {
	variantMarketplaceCode, _ := r.Param1("variantMarketplaceCode")

	qty, ok := r.Param1("qty")
	if !ok {
		qty = "1"
	}
	qtyInt, _ := strconv.Atoi(qty)
	deliveryIntent, _ := r.Param1("deliveryIntent")

	addRequest := cc.cartService.BuildAddRequest(ctx, r.MustParam1("marketplaceCode"), variantMarketplaceCode, qtyInt, deliveryIntent)
	err, _ := cc.cartService.AddProduct(ctx, r.Session(), addRequest)
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

func (cc *CartApiController) ApplyVoucherAndGetAction(ctx context.Context, r *web.Request) web.Response {
	couponCode := r.MustParam1("couponCode")

	cart, err := cc.cartService.ApplyVoucher(ctx, r.Session(), couponCode)
	if err != nil {
		return cc.JSONError(result{Message: err.Error(), Success: false}, 500)
	}
	return cc.JSON(cart)
}
