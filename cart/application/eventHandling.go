package application

import (
	"context"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	//EventReceiver - handles events from other packages
	EventReceiver struct {
		logger              flamingo.Logger
		cartService         *CartService
		cartReceiverService *CartReceiverService
		cartCache           CartCache
	}
)

// Inject dependencies
func (e EventReceiver) Inject(
	logger flamingo.Logger,
	cartService *CartService,
	cartReceiverService *CartReceiverService,
	cartCache CartCache,
) {
	e.logger = logger
	e.cartService = cartService
	e.cartReceiverService = cartReceiverService
	e.cartCache = cartCache
}

//Notify should get called by flamingo Eventlogic
func (e *EventReceiver) Notify(ctx context.Context, event flamingo.Event) {
	switch currentEvent := event.(type) {
	//Handle Logout
	case *domain.LogoutEvent:
		if e.cartCache != nil {
			e.cartCache.DeleteAll(ctx, currentEvent.Session)
		}
	//Handle LoginEvent and Merge Cart
	case *domain.LoginEvent:
		if currentEvent == nil {
			return
		}
		if !e.cartReceiverService.ShouldHaveGuestCart(currentEvent.Session) {
			return
		}
		guestCart, err := e.cartReceiverService.ViewGuestCart(ctx, currentEvent.Session)
		if err != nil {
			e.logger.WithField(flamingo.LogKeyCategory, "cart").Error("LoginEvent - Guestcart cannot be received %v", err)
			return
		}
		if !e.cartReceiverService.IsLoggedIn(ctx, currentEvent.Session) {
			e.logger.WithField(flamingo.LogKeyCategory, "cart").Error("Received LoginEvent but user is not logged in!!!")
			return
		}
		for _, d := range guestCart.Deliveries {
			for _, item := range d.Cartitems {
				e.logger.WithField(flamingo.LogKeyCategory, "cart").Debug("Merging item from guest to user cart %v", item)
				addRequest := e.cartService.BuildAddRequest(ctx, item.MarketplaceCode, item.VariantMarketPlaceCode, item.Qty)
				e.cartService.AddProduct(ctx, currentEvent.Session, d.DeliveryInfo.Code, addRequest)
			}
		}

		if guestCart.HasAppliedCouponCode() {
			for _, code := range guestCart.AppliedCouponCodes {
				e.cartService.ApplyVoucher(ctx, currentEvent.Session, code.Code)
			}
		}

		if e.cartCache != nil {
			cacheID, err := BuildIdentifierFromCart(guestCart)
			if err == nil {
				e.cartCache.Delete(ctx, currentEvent.Session, *cacheID)
			}
		}
		e.cartService.DeleteSavedSessionGuestCartID(currentEvent.Session)
	// Handle Event to Invalidate the Cart Cache
	case *cartDomain.InvalidateCartEvent:
		if e.cartCache != nil {
			cartID, err := e.cartCache.BuildIdentifier(ctx, currentEvent.Session)
			if err == nil {
				e.cartCache.Invalidate(ctx, currentEvent.Session, cartID)
			}
		}
	}
}
