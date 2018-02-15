package controller

import (
	"log"
	"strconv"

	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// CartViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart        cart.DecoratedCart
		CartValidationResult cart.CartValidationResult
	}

	// CartViewController for carts
	CartViewController struct {
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`

		ApplicationCartService *application.CartService `inject:""`
		Router                 *router.Router           `inject:""`
	}
)

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx web.Context) web.Response {
	decoratedCart, err := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if err != nil {
		log.Printf("cart.cartcontroller.viewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	return cc.Render(ctx, "checkout/cart", CartViewData{
		DecoratedCart:        decoratedCart,
		CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
	})

}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) AddAndViewAction(ctx web.Context) web.Response {
	addRequest := addRequestFromRequestContext(ctx)
	err := cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if err != nil {
		log.Printf("cart.cartcontroller.addandviewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	currentCart, err := cc.ApplicationCartService.GetCart(ctx)

	if err != nil {
		log.Printf("cart.cartcontroller.AddAndViewAction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	addedProduct, err := cc.ApplicationCartService.ProductService.Get(ctx, addRequest.MarketplaceCode)

	if err != nil {
		log.Printf("cart.cartcontroller.AddAndViewAction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	currentCart.EventPublisher.PublishAddToCartEvent(ctx, &addedProduct, addRequest.Qty)

	return cc.Redirect("cart.view", nil)
}

// UpdateQtyAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) UpdateQtyAndViewAction(ctx web.Context) web.Response {
	decoratedCart, err := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if err != nil {
		log.Printf("cart.cartcontroller.UpdateAndViewAction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

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

	err = decoratedCart.Cart.UpdateItemQty(ctx, cc.ApplicationCartService.Auth(ctx), id, qtyInt)
	if err != nil {
		log.Printf("cart.cartcontroller.UpdateAndViewAction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// DeleteAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) DeleteAndViewAction(ctx web.Context) web.Response {
	decoratedCart, err := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if err != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	id, err := ctx.Param1("id")
	if err != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", err)
		return cc.Redirect("cart.view", nil)
	}

	err = decoratedCart.Cart.DeleteItem(ctx, cc.ApplicationCartService.Auth(ctx), id)
	if err != nil {
		log.Printf("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}
