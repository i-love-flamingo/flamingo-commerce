package application

import (
	"context"

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

//Notify should get called by flamingo Eventlogic
func (e *EventReceiver) NotifyWithContext(ctx context.Context, event event.Event) {
	switch event := event.(type) {
	//Handle Logout
	case *domain.LogoutEvent:
		if e.CartCache != nil {
			e.CartCache.DeleteAll(ctx, event.Session)
		}
	//Handle LoginEvent and Merge Cart
	case *domain.LoginEvent:
		if event == nil {
			return
		}
		if !e.CartReceiverService.ShouldHaveGuestCart(event.Session) {
			return
		}
		guestCart, err := e.CartReceiverService.ViewGuestCart(ctx, event.Session)
		if err != nil {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Error("LoginEvent - Guestcart cannot be received %v", err)
			return
		}
		if !e.CartReceiverService.UserService.IsLoggedIn(ctx, event.Session) {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Error("Received LoginEvent but user is not logged in!!!")
			return
		}
		for _, item := range guestCart.Cartitems {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Debug("Merging item from guest to user cart %v", item)
			addRequest := e.CartService.BuildAddRequest(ctx, item.MarketplaceCode, item.VariantMarketPlaceCode, item.Qty, item.OriginalDeliveryIntent.String())
			e.CartService.AddProduct(ctx, event.Session, addRequest)
		}

		if guestCart.HasAppliedCouponCode() {
			for _, code := range guestCart.AppliedCouponCodes {
				e.CartService.ApplyVoucher(ctx, event.Session, code.Code)
			}
		}

		if e.CartCache != nil {
			cacheId, err := BuildIdentifierFromCart(guestCart)
			if err != nil {
				e.CartCache.Delete(ctx, event.Session, *cacheId)
			}
		}
		e.CartService.DeleteSavedSessionGuestCartId(event.Session)
	}
}
