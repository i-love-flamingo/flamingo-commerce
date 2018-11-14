package application

import (
	"context"

	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	//EventReceiver - handles events from other packages
	EventReceiver struct {
		Logger              flamingo.Logger      `inject:""`
		CartService         *CartService         `inject:""`
		CartReceiverService *CartReceiverService `inject:""`
		CartCache           CartCache            `inject:",optional"`
	}
)

//NotifyWithContext should get called by flamingo Eventlogic
func (e *EventReceiver) NotifyWithContext(ctx context.Context, event event.Event) {
	switch currentEvent := event.(type) {
	//Handle Logout
	case *domain.LogoutEvent:
		if e.CartCache != nil {
			e.CartCache.DeleteAll(ctx, currentEvent.Session)
		}
	//Handle LoginEvent and Merge Cart
	case *domain.LoginEvent:
		if currentEvent == nil {
			return
		}
		if !e.CartReceiverService.ShouldHaveGuestCart(currentEvent.Session) {
			return
		}
		guestCart, err := e.CartReceiverService.ViewGuestCart(ctx, currentEvent.Session)
		if err != nil {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Error("LoginEvent - Guestcart cannot be received %v", err)
			return
		}
		if !e.CartReceiverService.UserService.IsLoggedIn(ctx, currentEvent.Session) {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Error("Received LoginEvent but user is not logged in!!!")
			return
		}
		for _, d := range guestCart.Deliveries {
			for _, item := range d.Cartitems {
				e.Logger.WithField(flamingo.LogKeyCategory, "cart").Debug("Merging item from guest to user cart %v", item)
				addRequest := e.CartService.BuildAddRequest(ctx, item.MarketplaceCode, item.VariantMarketPlaceCode, item.Qty)
				e.CartService.AddProduct(ctx, currentEvent.Session, d.DeliveryInfo.Code, addRequest)
			}
		}

		if guestCart.HasAppliedCouponCode() {
			for _, code := range guestCart.AppliedCouponCodes {
				e.CartService.ApplyVoucher(ctx, currentEvent.Session, code.Code)
			}
		}

		if e.CartCache != nil {
			cacheID, err := BuildIdentifierFromCart(guestCart)
			if err == nil {
				e.CartCache.Delete(ctx, currentEvent.Session, *cacheID)
			}
		}
		e.CartService.DeleteSavedSessionGuestCartID(currentEvent.Session)
	// Handle Event to Invalidate the Cart Cache
	case *cartDomain.InvalidateCartEvent:
		if e.CartCache != nil {
			cartId, err := e.CartCache.BuildIdentifier(ctx, currentEvent.Session)
			if err == nil {
				e.CartCache.Invalidate(ctx, currentEvent.Session, cartId)
			}
		}
	}
}
