package controller

import (
	"context"
	"encoding/gob"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CartViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart         cartDomain.DecoratedCart
		CartValidationResult  cartDomain.ValidationResult
		AddToCartProductsData []productDomain.BasicProductData
	}

	// CartViewController for carts
	CartViewController struct {
		responder *web.Responder

		applicationCartService         *application.CartService
		applicationCartReceiverService *application.CartReceiverService
		router                         *web.Router
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
	responder *web.Responder,

	applicationCartService *application.CartService,
	applicationCartReceiverService *application.CartReceiverService,
	router *web.Router,
	logger flamingo.Logger,
	config *struct {
		ShowEmptyCartPageIfNoItems bool `inject:"config:cart.showEmptyCartPageIfNoItems,optional"`
	},
) {
	cc.responder = responder
	cc.applicationCartService = applicationCartService
	cc.applicationCartReceiverService = applicationCartReceiverService
	cc.router = router
	cc.logger = logger

	if config != nil {
		cc.showEmptyCartPageIfNoItems = config.ShowEmptyCartPageIfNoItems
	}
}

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx context.Context, r *web.Request) web.Result {
	decoratedCart, err := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.viewaction: Error %v", err)
		return cc.responder.Render("checkout/carterror", nil)
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.responder.Render("checkout/emptycart", nil)
	}

	cartViewData := CartViewData{
		DecoratedCart:        *decoratedCart,
		CartValidationResult: cc.applicationCartService.ValidateCart(ctx, r.Session(), decoratedCart),
	}

	flashes := r.Session().Flashes("cart.view.data")
	if len(flashes) > 0 {
		if cartViewActionData, ok := flashes[0].(CartViewActionData); ok {
			cartViewData.AddToCartProductsData = cartViewActionData.AddToCartProductsData
		}
	}

	return cc.responder.Render("checkout/cart", cartViewData)
}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) AddAndViewAction(ctx context.Context, r *web.Request) web.Result {
	variantMarketplaceCode, _ := r.Params["variantMarketplaceCode"]

	qty, ok := r.Params["qty"]
	if !ok {
		qty = "1"
	}

	qtyInt, _ := strconv.Atoi(qty)
	deliveryCode, _ := r.Params["deliveryCode"]

	addRequest := cc.applicationCartService.BuildAddRequest(ctx, r.Params["marketplaceCode"], variantMarketplaceCode, qtyInt)

	product, err := cc.applicationCartService.AddProduct(ctx, r.Session(), deliveryCode, addRequest)
	if notAllowedErr, ok := err.(*cartDomain.AddToCartNotAllowed); ok {
		if notAllowedErr.RedirectHandlerName != "" {
			return cc.responder.RouteRedirect(notAllowedErr.RedirectHandlerName, notAllowedErr.RedirectParams)
		}
	}
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.addandviewaction: Error %v", err)
		return cc.responder.Render("checkout/carterror", nil)
	}

	r.Session().AddFlash(CartViewActionData{
		AddToCartProductsData: []productDomain.BasicProductData{product.BaseData()},
	}, "cart.view.data")
	return cc.responder.RouteRedirect("cart.view", nil)
}

// UpdateQtyAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) UpdateQtyAndViewAction(ctx context.Context, r *web.Request) web.Result {

	id, ok := r.Params["id"]
	if !ok {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.UpdateQtyAndViewAction: param id not found")
		return cc.responder.RouteRedirect("cart.view", nil)
	}

	qty, ok := r.Params["qty"]
	if !ok {
		qty = "1"
	}

	qtyInt, err := strconv.Atoi(qty)
	if err != nil {
		qtyInt = 1
	}
	deliveryCode, _ := r.Params["deliveryCode"]

	err = cc.applicationCartService.UpdateItemQty(ctx, r.Session(), id, deliveryCode, qtyInt)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.UpdateAndViewAction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// DeleteAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) DeleteAndViewAction(ctx context.Context, r *web.Request) web.Result {

	id, ok := r.Params["id"]
	if !ok {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: param id not found")
		return cc.responder.RouteRedirect("cart.view", nil)
	}
	deliveryCode, _ := r.Params["deliveryCode"]
	err := cc.applicationCartService.DeleteItem(ctx, r.Session(), id, deliveryCode)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// DeleteAllAndViewAction empties the cart and shows it
func (cc *CartViewController) DeleteAllAndViewAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.applicationCartService.DeleteAllItems(ctx, r.Session())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// CleanAndViewAction empties the cart and shows it
func (cc *CartViewController) CleanAndViewAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.applicationCartService.Clean(ctx, r.Session())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// CleanDeliveryAndViewAction empties a single delivery and shows it
func (cc *CartViewController) CleanDeliveryAndViewAction(ctx context.Context, r *web.Request) web.Result {
	deliveryCode := r.Params["deliveryCode"]
	_, err := cc.applicationCartService.DeleteDelivery(ctx, r.Session(), deliveryCode)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}
