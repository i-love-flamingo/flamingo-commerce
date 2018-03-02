package controller

import (
	"fmt"
	"log"
	"strconv"

	"go.aoe.com/flamingo/core/cart/application"
	domaincart "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// CartApiController for cart api
	CartApiController struct {
		responder.JSONAware    `inject:""`
		ApplicationCartService *application.CartService `inject:""`
		DefaultDeliveryIntent  string                   `inject:"config:cart.defaultDeliveryIntent,optional"`
	}

	result struct {
		Message string
		Success bool
	}
)

// GetAction Get JSON Format of API
func (cc *CartApiController) GetAction(ctx web.Context) web.Response {
	cart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.cartapicontroller.get: %v", e.Error())
		return cc.JSONError(result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(cart)
}

// AddAction Add Item to cart
func (cc *CartApiController) AddAction(ctx web.Context) web.Response {
	addRequest := addRequestFromRequestContext(ctx, cc.DefaultDeliveryIntent)
	e := cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if e != nil {
		log.Printf("cart.cartapicontroller.add: %v", e.Error())
		return cc.JSONError(result{Message: e.Error(), Success: false}, 500)
	}
	return cc.JSON(result{
		Success: true,
		Message: fmt.Sprintf("added %v / %v Qty %v", addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty),
	})
}

func addRequestFromRequestContext(ctx web.Context, defaultDeliveryIntent string) domaincart.AddRequest {
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
	if qtyInt < 0 {
		qtyInt = 0
	}

	deliveryIntent, e := ctx.Param1("deliveryIntent")
	if e != nil {
		deliveryIntent = defaultDeliveryIntent
	}

	return domaincart.AddRequest{MarketplaceCode: marketplaceCode, Qty: qtyInt, VariantMarketplaceCode: variantMarketplaceCode, DeliveryIntent: deliveryIntent}
}
