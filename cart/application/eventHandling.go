package application

import (
	"go.aoe.com/flamingo/core/auth/domain"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
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
func (e *EventReceiver) Notify(event event.Event) {
	switch eventType := event.(type) {
	//Handle Logout
	case *domain.LogoutEvent:
		if e.CartCache != nil {
			e.CartCache.DeleteAll(eventType.Context)
		}
	//Handle LoginEvent and Merge Cart
	case *domain.LoginEvent:
		if eventType == nil {
			return
		}
		if eventType.Context == nil {
			return
		}
		if !e.CartReceiverService.ShouldHaveGuestCart(eventType.Context) {
			return
		}
		guestCart, err := e.CartReceiverService.ViewGuestCart(eventType.Context)
		if err != nil {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Error("LoginEvent - Guestcart cannot be received %v", err)
			return
		}
		if !e.CartReceiverService.UserService.IsLoggedIn(eventType.Context) {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Error("Received LoginEvent but user is not logged in!!!")
			return
		}
		for _, item := range guestCart.Cartitems {
			e.Logger.WithField(flamingo.LogKeyCategory, "cart").Debug("Merging item from guest to user cart %v", item)
			addRequest := e.CartService.BuildAddRequest(eventType.Context, item.MarketplaceCode, item.VariantMarketPlaceCode, item.Qty, item.OriginalDeliveryIntent.String())
			e.CartService.AddProduct(eventType.Context, addRequest)
		}

		if guestCart.HasAppliedCouponCode() {
			for _, code := range guestCart.AppliedCouponCodes {
				e.CartService.ApplyVoucher(eventType.Context, code.Code)
			}
		}

		if e.CartCache != nil {
			cacheId, err := BuildIdentifierFromCart(guestCart)
			if err != nil {
				e.CartCache.Delete(eventType.Context, *cacheId)
			}
		}
		e.CartService.DeleteSavedSessionGuestCartId(eventType.Context)

	}
}
