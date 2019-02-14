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
		CartValidationResult  cartDomain.CartValidationResult
		AddToCartProductsData []productDomain.BasicProductData
	}

	// CartViewController for carts
	CartViewController struct {
		Responder *web.Responder
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
	renderAware Responder *web.Responder,
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
func (cc *CartViewController) ViewAction(ctx context.Context, r *web.Request) web.Result {
	decoratedCart, err := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session().G())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.viewaction: Error %v", err)
		return cc.Responder.Render( "checkout/carterror", nil)
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.Responder.Render( "checkout/emptycart", nil)
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

	return cc.Responder.Render( "checkout/cart", cartViewData)
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

	addRequest := cc.applicationCartService.BuildAddRequest(ctx, r.MustParam1("marketplaceCode"), variantMarketplaceCode, qtyInt)

	product, err := cc.applicationCartService.AddProduct(ctx, r.Session().G(), deliveryCode, addRequest)
	if notAllowedErr, ok := err.(*cartDomain.AddToCartNotAllowed); ok {
		if notAllowedErr.RedirectHandlerName != "" {
			return cc.Redirect(notAllowedErr.RedirectHandlerName, notAllowedErr.RedirectParams)
		}
	}
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.addandviewaction: Error %v", err)
		return cc.Responder.Render( "checkout/carterror", nil)
	}

	return cc.Redirect("cart.view", nil).With("cart.view.data", CartViewActionData{
		AddToCartProductsData: []productDomain.BasicProductData{product.BaseData()},
	})
}

// UpdateQtyAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) UpdateQtyAndViewAction(ctx context.Context, r *web.Request) web.Result {

	id, ok := r.Params["id"]
	if !ok {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.UpdateQtyAndViewAction: param id not found")
		return cc.Redirect("cart.view", nil)
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

	err = cc.applicationCartService.UpdateItemQty(ctx, r.Session().G(), id, deliveryCode, qtyInt)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.UpdateAndViewAction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// DeleteAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) DeleteAndViewAction(ctx context.Context, r *web.Request) web.Result {

	id, ok := r.Params["id"]
	if !ok {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: param id not found")
		return cc.Redirect("cart.view", nil)
	}
	deliveryCode, _ := r.Params["deliveryCode"]
	err := cc.applicationCartService.DeleteItem(ctx, r.Session().G(), id, deliveryCode)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// DeleteAllAndViewAction empties the cart and shows it
func (cc *CartViewController) DeleteAllAndViewAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.applicationCartService.DeleteAllItems(ctx, r.Session().G())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// CleanAndViewAction empties the cart and shows it
func (cc *CartViewController) CleanAndViewAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.applicationCartService.Clean(ctx, r.Session().G())
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}

// CleanDeliveryAndViewAction empties a single delivery and shows it
func (cc *CartViewController) CleanDeliveryAndViewAction(ctx context.Context, r *web.Request) web.Result {
	deliveryCode := r.MustParam1("deliveryCode")
	_, err := cc.applicationCartService.DeleteDelivery(ctx, r.Session().G(), deliveryCode)
	if err != nil {
		cc.logger.WithField(flamingo.LogKeyCategory, "cartcontroller").Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.Redirect("cart.view", nil)
}
