package controller

import (
	"context"
	"encoding/gob"
	"strconv"

	"flamingo.me/flamingo-commerce/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
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
		responder.RenderAware
		responder.RedirectAware

		applicationCartService         *application.CartService
		applicationCartReceiverService *application.CartReceiverService
		router                         *router.Router
		logger                         flamingo.Logger

		showEmptyCartPageIfNoItems bool
	}

	// CartViewActionData for rendering results
	CartViewActionData struct {
		AddToCartProductsData []productDomain.BasicProductData
	}
)

func init() {
	gob.Register(CartViewActionData{})
}

// Inject dependencies
func (cc *CartViewController) Inject(
	renderAware responder.RenderAware,
	redirectAware responder.RedirectAware,

	applicationCartService *application.CartService,
	applicationCartReceiverService *application.CartReceiverService,
	router *router.Router,
	logger flamingo.Logger,
	config *struct {
		ShowEmptyCartPageIfNoItems bool `inject:"config:cart.showEmptyCartPageIfNoItems,optional"`
	},
) {
	cc.RenderAware = renderAware
	cc.RedirectAware = redirectAware
	cc.applicationCartService = applicationCartService
	cc.applicationCartReceiverService = applicationCartReceiverService
	cc.router = router
	cc.logger = logger

	if config != nil {
		cc.showEmptyCartPageIfNoItems = config.ShowEmptyCartPageIfNoItems
	}
}

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx context.Context, r *web.Request) web.Response {
	decoratedCart, err := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session().G())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.viewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, "checkout/emptycart", nil)
	}

	cartViewData := CartViewData{
		DecoratedCart:        *decoratedCart,
		CartValidationResult: cc.applicationCartService.ValidateCart(ctx, r.Session().G(), decoratedCart),
	}

	flashes := r.Session().Flashes("cart.view.data")
	if len(flashes) > 0 {
		if cartViewActionData, ok := flashes[0].(CartViewActionData); ok {
			cartViewData.AddToCartProductsData = cartViewActionData.AddToCartProductsData
		}
	}

	return cc.Render(ctx, "checkout/cart", cartViewData)
}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) AddAndViewAction(ctx context.Context, r *web.Request) web.Response {
	variantMarketplaceCode, _ := r.Param1("variantMarketplaceCode")

	qty, ok := r.Param1("qty")
	if !ok {
		qty = "1"
	}

	qtyInt, _ := strconv.Atoi(qty)
	deliveryCode, _ := r.Param1("deliveryCode")

	addRequest := cc.applicationCartService.BuildAddRequest(ctx, r.MustParam1("marketplaceCode"), variantMarketplaceCode, qtyInt)

	product, err := cc.applicationCartService.AddProduct(ctx, r.Session().G(), deliveryCode, addRequest)
	if notAllowedErr, ok := err.(*cartDomain.AddToCartNotAllowed); ok {
		if notAllowedErr.RedirectHandlerName != "" {
			return cc.Redirect(notAllowedErr.RedirectHandlerName, notAllowedErr.RedirectParams)
		}
	}
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.addandviewaction: Error %v", err)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	return cc.Redirect("cart.view", nil).With("cart.view.data", CartViewActionData{
		AddToCartProductsData: []productDomain.BasicProductData{product.BaseData()},
	})
}

// UpdateQtyAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) UpdateQtyAndViewAction(ctx context.Context, r *web.Request) web.Response {

	id, ok := r.Param1("id")
	if !ok {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.UpdateQtyAndViewAction: param id not found")
		return cc.Redirect("cart.view", nil)
	}

	qty, ok := r.Param1("qty")
	if !ok {
		qty = "1"
	}

	qtyInt, err := strconv.Atoi(qty)
	if err != nil {
		qtyInt = 1
	}
	deliveryCode, _ := r.Param1("deliveryCode")

	err = cc.applicationCartService.UpdateItemQty(ctx, r.Session().G(), id, deliveryCode, qtyInt)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.UpdateAndViewAction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// DeleteAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) DeleteAndViewAction(ctx context.Context, r *web.Request) web.Response {

	id, ok := r.Param1("id")
	if !ok {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: param id not found")
		return cc.Redirect("cart.view", nil)
	}
	deliveryCode, _ := r.Param1("deliveryCode")
	err := cc.applicationCartService.DeleteItem(ctx, r.Session().G(), id, deliveryCode)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// DeleteAllAndViewAction empties the cart and shows it
func (cc *CartViewController) DeleteAllAndViewAction(ctx context.Context, r *web.Request) web.Response {
	err := cc.applicationCartService.DeleteAllItems(ctx, r.Session().G())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// CleanAndViewAction empties the cart and shows it
func (cc *CartViewController) CleanAndViewAction(ctx context.Context, r *web.Request) web.Response {
	err := cc.applicationCartService.Clean(ctx, r.Session().G())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// CleanDeliveryAndViewAction empties a single delivery and shows it
func (cc *CartViewController) CleanDeliveryAndViewAction(ctx context.Context, r *web.Request) web.Response {
	deliveryCode := r.MustParam1("deliveryCode")
	_, err := cc.applicationCartService.DeleteDelivery(ctx, r.Session().G(), deliveryCode)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}
