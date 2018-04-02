package controller

import (
	"log"
	"strconv"

	"go.aoe.com/flamingo/core/cart/application"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// CartViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart        cartDomain.DecoratedCart
		CartValidationResult cartDomain.CartValidationResult
	}

	// CartViewController for carts
	CartViewController struct {
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`

		ApplicationCartService         *application.CartService         `inject:""`
		ApplicationCartReceiverService *application.CartReceiverService `inject:""`
		Router                         *router.Router                   `inject:""`
	}
)

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx web.Context) web.Response {
	decoratedCart, err := cc.ApplicationCartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		log.Printf("cart.cartcontroller.viewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	return cc.Render(ctx, "checkout/cart", CartViewData{
		DecoratedCart:        *decoratedCart,
		CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
	})

}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) AddAndViewAction(ctx web.Context) web.Response {

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

	err := cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if notAllowedErr, ok := err.(*cartDomain.AddToCartNotAllowed); ok {
		if notAllowedErr.RedirectHandlerName != "" {
			return cc.Redirect(notAllowedErr.RedirectHandlerName, notAllowedErr.RedirectParams)
		}
	}
	if err != nil {
		log.Printf("cart.cartcontroller.addandviewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	return cc.Redirect("cart.view", nil)
}

// UpdateQtyAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) UpdateQtyAndViewAction(ctx web.Context) web.Response {

	id, err := ctx.Param1("id")
	if err != nil {
		log.Printf("cart.cartcontroller.UpdateAndViewAction: Error %v", err)
		return cc.Redirect("cart.view", nil)
	}

	qty, err := ctx.Param1("qty")
	if err != nil {
		qty = "1"
	}

	qtyInt, err := strconv.Atoi(qty)
	if err != nil {
		qtyInt = 1
	}

	err = cc.ApplicationCartService.UpdateItemQty(ctx, id, qtyInt)

	if err != nil {
		log.Printf("cart.cartcontroller.UpdateAndViewAction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// DeleteAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) DeleteAndViewAction(ctx web.Context) web.Response {

	id, err := ctx.Param1("id")
	if err != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", err)
		return cc.Redirect("cart.view", nil)
	}

	err = cc.ApplicationCartService.DeleteItem(ctx, id)
	if err != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}
