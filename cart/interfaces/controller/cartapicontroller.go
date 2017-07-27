package controller

import (
	"flamingo/core/cart/application"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"fmt"
	"strconv"
)

type (

	// CartAPIController for cart api
	CartApiAddController struct {
		responder.JSONAware    `inject:""`
		ApplicationCartService application.CartService `inject:""`
	}

	// CartAPIController for cart api
	CartApiGetController struct {
		responder.JSONAware    `inject:""`
		ApplicationCartService application.CartService `inject:""`
	}
)

// Get JSON Format of API
func (cc *CartApiGetController) Get(ctx web.Context) web.Response {
	cart, e := cc.ApplicationCartService.GetCart(ctx)
	if e != nil {
		fmt.Println(e.Error())
		return cc.JSON(struct{ Message string }{e.Error()})
	}
	return cc.JSON(CartViewData{
		Cart: cart,
	})
}

// Add Item to cart
func (cc *CartApiAddController) Get(ctx web.Context) web.Response {
	productCode := ctx.MustQuery1("productcode")
	qty, _ := ctx.Query1("qty")
	qtyInt, _ := strconv.Atoi(qty)

	e := cc.ApplicationCartService.AddProduct(ctx, productCode, qtyInt)
	if e != nil {
		return cc.JSON(struct{ Message string }{
			e.Error(),
		})
	}
	return cc.JSON(struct{ Message string }{
		"added " + productCode,
	})
}
