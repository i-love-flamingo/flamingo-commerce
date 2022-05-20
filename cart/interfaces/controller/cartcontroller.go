package controller

import (
	"context"
	"encoding/gob"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CartViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart         decorator.DecoratedCart
		CartValidationResult  validation.Result
		AddToCartProductsData []productDomain.BasicProductData
		CartRestrictionError  application.RestrictionError
	}

	// CartViewController for carts
	CartViewController struct {
		responder *web.Responder

		applicationCartService         *application.CartService
		applicationCartReceiverService *application.CartReceiverService
		router                         *web.Router
		logger                         flamingo.Logger

		showEmptyCartPageIfNoItems bool
		adjustItemsToRestrictedQty bool
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
		ShowEmptyCartPageIfNoItems bool `inject:"config:commerce.cart.showEmptyCartPageIfNoItems,optional"`
		AdjustItemsToRestrictedQty bool `inject:"config:commerce.cart.adjustItemsToRestrictedQty,optional"`
	},
) {
	cc.responder = responder
	cc.applicationCartService = applicationCartService
	cc.applicationCartReceiverService = applicationCartReceiverService
	cc.router = router
	cc.logger = logger.WithField(flamingo.LogKeyCategory, "cartcontroller").WithField(flamingo.LogKeyModule, "cart")

	if config != nil {
		cc.showEmptyCartPageIfNoItems = config.ShowEmptyCartPageIfNoItems
		cc.adjustItemsToRestrictedQty = config.AdjustItemsToRestrictedQty
	}
}

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx context.Context, r *web.Request) web.Result {
	if cc.adjustItemsToRestrictedQty {
		err := cc.adjustItemsToRestrictedQtyAndAddToSession(ctx, r)
		if err != nil {
			cc.logger.WithContext(ctx).Warn("cart.cartcontroller.viewaction: Error %v", err)
		}
	}

	decoratedCart, err := cc.applicationCartReceiverService.ViewDecoratedCart(ctx, r.Session())
	if err != nil {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.viewaction: Error %v", err)
		return cc.responder.Render("checkout/carterror", nil).SetNoCache()
	}

	if cc.showEmptyCartPageIfNoItems && decoratedCart.Cart.ItemCount() == 0 {
		return cc.responder.Render("checkout/emptycart", nil).SetNoCache()
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
	restrictionFlashes := r.Session().Flashes("cart.view.error.restriction")
	if len(restrictionFlashes) > 0 {
		if restrictionError, ok := restrictionFlashes[0].(application.RestrictionError); ok {
			cartViewData.CartRestrictionError = restrictionError
		}
	}

	return cc.responder.Render("checkout/cart", cartViewData).SetNoCache()
}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) AddAndViewAction(ctx context.Context, r *web.Request) web.Result {
	variantMarketplaceCode := r.Params["variantMarketplaceCode"]

	qty, ok := r.Params["qty"]
	if !ok {
		qty = "1"
	}

	qtyInt, _ := strconv.Atoi(qty)
	deliveryCode := r.Params["deliveryCode"]

	addRequest := cc.applicationCartService.BuildAddRequest(ctx, r.Params["marketplaceCode"], variantMarketplaceCode, qtyInt, nil)

	product, err := cc.applicationCartService.AddProduct(ctx, r.Session(), deliveryCode, addRequest)
	if notAllowedErr, ok := err.(*validation.AddToCartNotAllowed); ok {
		if notAllowedErr.RedirectHandlerName != "" {
			return cc.responder.RouteRedirect(notAllowedErr.RedirectHandlerName, notAllowedErr.RedirectParams)
		}
	}
	if err != nil {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.addandviewaction: Error %v", err)
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
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.UpdateQtyAndViewAction: param id not found")
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
	deliveryCode := r.Params["deliveryCode"]

	err = cc.applicationCartService.UpdateItemQty(ctx, r.Session(), id, deliveryCode, qtyInt)
	if err != nil {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.UpdateAndViewAction: Error %v", err)

		if e, ok := err.(*application.RestrictionError); ok {
			r.Session().AddFlash(e, "cart.view.error.restriction")
		}
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// DeleteAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) DeleteAndViewAction(ctx context.Context, r *web.Request) web.Result {

	id, ok := r.Params["id"]
	if !ok {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.deleteaction: param id not found")
		return cc.responder.RouteRedirect("cart.view", nil)
	}
	deliveryCode := r.Params["deliveryCode"]
	err := cc.applicationCartService.DeleteItem(ctx, r.Session(), id, deliveryCode)
	if err != nil {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// DeleteAllAndViewAction empties the cart and shows it
func (cc *CartViewController) DeleteAllAndViewAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.applicationCartService.DeleteAllItems(ctx, r.Session())
	if err != nil {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// CleanAndViewAction empties the cart and shows it
func (cc *CartViewController) CleanAndViewAction(ctx context.Context, r *web.Request) web.Result {
	err := cc.applicationCartService.Clean(ctx, r.Session())
	if err != nil {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// CleanDeliveryAndViewAction empties a single delivery and shows it
func (cc *CartViewController) CleanDeliveryAndViewAction(ctx context.Context, r *web.Request) web.Result {
	deliveryCode := r.Params["deliveryCode"]
	_, err := cc.applicationCartService.DeleteDelivery(ctx, r.Session(), deliveryCode)
	if err != nil {
		cc.logger.WithContext(ctx).Warn("cart.cartcontroller.deleteaction: Error %v", err)
	}

	return cc.responder.RouteRedirect("cart.view", nil)
}

// adjustItemsToRestrictedQtyAndAddToSession checks the items of the cart against their qty restrictions and adds adjustments to the session
func (cc *CartViewController) adjustItemsToRestrictedQtyAndAddToSession(ctx context.Context, r *web.Request) error {
	adjustments, err := cc.applicationCartService.AdjustItemsToRestrictedQty(ctx, r.Session())
	if err != nil {
		return err
	}

	cc.addAdjustmentsToSession(adjustments, r)

	return nil
}

func (cc *CartViewController) addAdjustmentsToSession(adjustments application.QtyAdjustmentResults, r *web.Request) {
	var storedAdjustments application.QtyAdjustmentResults
	var ok bool

	if sessionStoredAdjustments, found := r.Session().Load("cart.view.quantity.adjustments"); found {
		if storedAdjustments, ok = sessionStoredAdjustments.(application.QtyAdjustmentResults); !ok {
			storedAdjustments = application.QtyAdjustmentResults{}
		}
	} else {
		storedAdjustments = application.QtyAdjustmentResults{}
	}

	for _, adjustment := range adjustments {
		if i, contains := cc.containsAdjustment(storedAdjustments, adjustment); contains {
			storedAdjustments[i] = adjustment
		} else {
			storedAdjustments = append(storedAdjustments, adjustment)
		}
	}

	r.Session().Store("cart.view.quantity.adjustments", storedAdjustments)
}

func (cc *CartViewController) containsAdjustment(adjustments application.QtyAdjustmentResults, adjustment application.QtyAdjustmentResult) (index int, contains bool) {
	for i, a := range adjustments {
		if a.OriginalItem.ID == adjustment.OriginalItem.ID && a.DeliveryCode == adjustment.DeliveryCode {
			return i, true
		}
	}

	return -1, false
}
