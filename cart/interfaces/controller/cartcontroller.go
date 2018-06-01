package controller

import (
	"log"
	"strconv"

	"go.aoe.com/flamingo/core/cart/application"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
	"encoding/gob"
)

type (
	// CartViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart         cartDomain.DecoratedCart
		CartValidationResult  cartDomain.CartValidationResult
		AddToCartProductsData []productDomain.BasicProductData
	}

	// CartViewController for carts
	CartViewController struct {
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`

		ApplicationCartService         *application.CartService         `inject:""`
		ApplicationCartReceiverService *application.CartReceiverService `inject:""`
		Router                         *router.Router                   `inject:""`

		ShowEmptyCartPageIfNoItems bool `inject:"config:cart.showEmptyCartPageIfNoItems,optional"`
	}

	CartViewActionData struct {
		AddToCartProductsData []productDomain.BasicProductData
	}
)

func init() {
	gob.Register(CartViewActionData{})
}

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx web.Context) web.Response {
	decoratedCart, err := cc.ApplicationCartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		log.Printf("cart.cartcontroller.viewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if cc.ShowEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, "checkout/emptycart", nil)
	}

	cartViewData := CartViewData{
		DecoratedCart:        *decoratedCart,
		CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
	}

	flashes := ctx.Session().Flashes("cart.view.data")
	if len(flashes) > 0 {
		if cartViewActionData, ok := flashes[0].(CartViewActionData); ok {
			cartViewData.AddToCartProductsData = cartViewActionData.AddToCartProductsData
		}
	}

	return cc.Render(ctx, "checkout/cart", cartViewData)

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

	err, product := cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if notAllowedErr, ok := err.(*cartDomain.AddToCartNotAllowed); ok {
		if notAllowedErr.RedirectHandlerName != "" {
			return cc.Redirect(notAllowedErr.RedirectHandlerName, notAllowedErr.RedirectParams)
		}
	}
	if err != nil {
		log.Printf("cart.cartcontroller.addandviewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	return cc.Redirect("cart.view", nil).With("cart.view.data", CartViewActionData{
		AddToCartProductsData: []productDomain.BasicProductData{product.BaseData()},
	})
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
