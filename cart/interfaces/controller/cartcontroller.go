package controller

import (
	"log"

	"go.aoe.com/flamingo/core/breadcrumbs"
	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"

	"strconv"

	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// ViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart cart.DecoratedCart
	}

	// CartController for carts
	CartViewController struct {
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		ApplicationCartService  application.CartService `inject:""`
		Router                  *router.Router          `inject:""`
	}
)

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx web.Context) web.Response {

	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.cartcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	breadcrumbs.Add(ctx, breadcrumbs.Crumb{
		Title: "Shopping Bag",
		Url:   cc.Router.URL("cart.view", nil).String(),
	})

	return cc.Render(ctx, "checkout/cart", CartViewData{
		DecoratedCart: decoratedCart,
	})

}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) AddAndViewAction(ctx web.Context) web.Response {
	addRequest := AddRequestFromRequestContext(ctx)
	e := cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if e != nil {
		log.Printf("cart.cartcontroller.addandviewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}
	return cc.Redirect("cart.view", nil)
}

// UpdateAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) UpdateQtyAndViewAction(ctx web.Context) web.Response {
	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.cartcontroller.UpdateAndViewAction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}
	id, e := ctx.Param1("id")
	if e != nil {
		log.Printf("cart.cartcontroller.UpdateAndViewAction: Error %v", e)
		return cc.Redirect("cart.view", nil)
	}
	qty, e := ctx.Param1("qty")
	if e != nil {
		qty = "1"
	}
	var qtyInt int
	qtyInt, e = strconv.Atoi(qty)
	if e != nil {
		qtyInt = 1
	}
	e = decoratedCart.Cart.UpdateItemQty(ctx, id, qtyInt)
	if e != nil {
		log.Printf("cart.cartcontroller.UpdateAndViewAction: Error %v", e)
	}
	return cc.Redirect("cart.view", nil)

}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) DeleteAndViewAction(ctx web.Context) web.Response {
	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}
	id, e := ctx.Param1("id")
	if e != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", e)
		return cc.Redirect("cart.view", nil)
	}

	e = decoratedCart.Cart.DeleteItem(ctx, id)
	if e != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", e)
	}
	return cc.Redirect("cart.view", nil)

}
