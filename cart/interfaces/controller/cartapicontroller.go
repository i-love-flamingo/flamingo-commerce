package controller

import (
	"fmt"
	"log"
	"strconv"

	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (

	// CartAPIController for cart api
	CartApiController struct {
		responder.JSONAware    `inject:""`
		ApplicationCartService application.CartService `inject:""`
	}
	Result struct {
		Message string
		Success bool
	}
)

// Get JSON Format of API
func (cc *CartApiController) GetAction(ctx web.Context) web.Response {
	cart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.cartapicontroller.get: %v", e.Error())
		return cc.JSONError(Result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(cart)
}

// Add Item to cart
func (cc *CartApiController) AddAction(ctx web.Context) web.Response {
	addRequest := AddRequestFromRequestContext(ctx)
	e := cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if e != nil {
		log.Printf("cart.cartapicontroller.add: %v", e.Error())
		return cc.JSONError(Result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(Result{
		Success: true,
		Message: fmt.Sprintf("added %v / %v Qty %v", addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty),
	})
}

func AddRequestFromRequestContext(ctx web.Context) cart.AddRequest {
	marketplaceCode := ctx.MustParam1("marketplaceCode")
	qty, e := ctx.Param1("qty")
	if e != nil {
		qty = "1"
	}
	variantMarketplaceCode, e := ctx.Param1("variantMarketplaceCode")
	if e != nil {
		variantMarketplaceCode = ""
	}
	qtyInt, _ := strconv.Atoi(qty)
	addRequest := cart.AddRequest{MarketplaceCode: marketplaceCode, Qty: qtyInt, VariantMarketplaceCode: variantMarketplaceCode}
	return addRequest
}
