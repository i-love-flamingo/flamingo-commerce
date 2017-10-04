package controller

import (
	"flamingo/core/cart/application"
	"flamingo/core/cart/domain/cart"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"fmt"
	"strconv"
)

type (

	// CartAPIController for cart api
	CartApiController struct {
		responder.JSONAware    `inject:""`
		ApplicationCartService application.CartService `inject:""`
	}
)

// Get JSON Format of API
func (cc *CartApiController) GetAction(ctx web.Context) web.Response {
	cart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		fmt.Println(e.Error())
		return cc.JSON(struct{ Message string }{e.Error()})
	}
	return cc.JSON(cart)
}

// Add Item to cart
func (cc *CartApiController) AddAction(ctx web.Context) web.Response {
	marketplaceCode := ctx.MustQuery1("marketplaceCode")
	qty, e := ctx.Query1("qty")
	if e != nil {
		qty = "1"
	}
	variantMarketplaceCode, e := ctx.Query1("variantMarketplaceCode")
	if e != nil {
		variantMarketplaceCode = ""
	}
	qtyInt, _ := strconv.Atoi(qty)
	addRequest := cart.AddRequest{MarketplaceCode: marketplaceCode, Qty: qtyInt, VariantMarketplaceCode: variantMarketplaceCode}
	e = cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if e != nil {
		return cc.JSON(struct{ Message string }{
			e.Error(),
		})
	}
	return cc.JSON(struct{ Message string }{
		"added " + marketplaceCode + "/" + variantMarketplaceCode + " qty: " + qty,
	})
}
