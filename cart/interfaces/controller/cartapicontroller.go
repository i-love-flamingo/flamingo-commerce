package controller

import (
	"fmt"
	"log"
	"strconv"

	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// CartApiController for cart api
	CartApiController struct {
		responder.JSONAware            `inject:""`
		ApplicationCartService         *application.CartService         `inject:""`
		ApplicationCartReceiverService *application.CartReceiverService `inject:""`
	}

	result struct {
		Message string
		Success bool
	}
)

// GetAction Get JSON Format of API
func (cc *CartApiController) GetAction(ctx web.Context) web.Response {
	cart, e := cc.ApplicationCartReceiverService.ViewDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.cartapicontroller.get: %v", e.Error())
		return cc.JSONError(result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(cart)
}

// AddAction Add Item to cart
func (cc *CartApiController) AddAction(ctx web.Context) web.Response {
	variantMarketplaceCode, e := ctx.Param1("variantMarketplaceCode")
	if e != nil {
		variantMarketplaceCode = ""
	}
	qty, e := ctx.Param1("qty")
	if e != nil {
		qty = "1"
	}
	qtyInt, _ := strconv.Atoi(qty)
	deliveryIntent, e := ctx.Param1("deliveryIntent")

	addRequest := cc.ApplicationCartService.BuildAddRequest(ctx, ctx.MustParam1("marketplaceCode"), variantMarketplaceCode, qtyInt, deliveryIntent)
	e = cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if e != nil {
		log.Printf("cart.cartapicontroller.add: %v", e.Error())
		return cc.JSONError(result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(result{
		Success: true,
		Message: fmt.Sprintf("added %v / %v Qty %v", addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty),
	})
}
