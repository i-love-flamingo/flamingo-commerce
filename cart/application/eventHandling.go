package application

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/core/oauth/domain"
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
func (e *EventReceiver) Inject(
	logger flamingo.Logger,
	cartService *CartService,
	cartReceiverService *CartReceiverService,
	optionals *struct {
		CartCache CartCache `inject:",optional"`
	},
) {
	e.logger = logger.WithField(flamingo.LogKeyCategory, "cart").WithField(flamingo.LogKeySubCategory, "cart-events")
	e.cartService = cartService
	e.cartReceiverService = cartReceiverService
	if optionals != nil {
		e.cartCache = optionals.CartCache
	}
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
		web.RunWithDetachedContext(ctx, func(ctx context.Context) {
			if currentEvent == nil {
				return
			}
			if !e.cartReceiverService.ShouldHaveGuestCart(currentEvent.Session) {
				return
			}
			guestCart, err := e.cartReceiverService.ViewGuestCart(ctx, currentEvent.Session)
			if err != nil {
				e.logger.WithContext(ctx).Error("LoginEvent - GuestCart cannot be received %v", err)
				return
			}
			if !e.cartReceiverService.IsLoggedIn(ctx, currentEvent.Session) {
				e.logger.WithContext(ctx).Error("Received LoginEvent but user is not logged in!!!")
				return
			}
			customerCart, err := e.cartReceiverService.ViewCart(ctx, currentEvent.Session)
			if err != nil {
				e.logger.WithContext(ctx).Error("LoginEvent - customerCart cannot be received %v", err)
				return
			}
			err = e.cartService.DeleteSavedSessionGuestCartID(currentEvent.Session)
			if err != nil {
				e.logger.WithContext(ctx).Error("LoginEvent - DeleteSavedSessionGuestCartID Error", err)
			}
			for _, d := range guestCart.Deliveries {
				e.logger.WithContext(ctx).Info(fmt.Sprintf("Merging delivery with code %v of guestCart with ID %v into customerCart with ID %v", d.DeliveryInfo.Code, guestCart.ID, customerCart.ID))
				err := e.cartService.UpdateDeliveryInfo(ctx, currentEvent.Session, d.DeliveryInfo.Code, cartDomain.CreateDeliveryInfoUpdateCommand(d.DeliveryInfo))
				if err != nil {
					e.logger.WithContext(ctx).Error("LoginEvent - customerCart UpdateDeliveryInfo error", err)
					continue
				}
				for _, item := range d.Cartitems {
					e.logger.WithContext(ctx).Debugf("Merging item from guest to user cart %v", item)
					addRequest := e.cartService.BuildAddRequest(ctx, item.MarketplaceCode, item.VariantMarketPlaceCode, item.Qty)
					_, err := e.cartService.AddProduct(ctx, currentEvent.Session, d.DeliveryInfo.Code, addRequest)
					if err != nil {
						e.logger.WithContext(ctx).Error("LoginEvent - customerCart product has merge error", addRequest.MarketplaceCode, err)
					}
				}
			}
			if customerCart.BillingAddress == nil && guestCart.BillingAddress != nil {
				err := e.cartService.UpdateBillingAddress(ctx, currentEvent.Session, guestCart.BillingAddress)
				if err != nil {
					e.logger.WithContext(ctx).Error("LoginEvent - customerCart UpdateBillingAddress error", err)
				}
			}
			if customerCart.Purchaser == nil && guestCart.Purchaser != nil {
				err := e.cartService.UpdatePurchaser(ctx, currentEvent.Session, guestCart.Purchaser, &guestCart.AdditionalData)
				if err != nil {
					e.logger.WithContext(ctx).Error("LoginEvent - customerCart UpdatePurchaser error", err)
				}
			}
			if customerCart.PaymentSelection == nil && guestCart.PaymentSelection != nil {
				err := e.cartService.UpdatePaymentSelection(ctx, currentEvent.Session, guestCart.PaymentSelection)
				if err != nil {
					e.logger.WithContext(ctx).Error("LoginEvent - customerCart UpdatePaymentSelection error", err)
				}
			}
			if guestCart.HasAppliedCouponCode() {
				for _, code := range guestCart.AppliedCouponCodes {
					_, err := e.cartService.ApplyVoucher(ctx, currentEvent.Session, code.Code)
					if err != nil {
						e.logger.WithContext(ctx).Error("LoginEvent - customerCart ApplyVoucher has error", code.Code, err)
					}
				}
			}
			if guestCart.HasAppliedGiftCards() {
				for _, code := range guestCart.AppliedGiftCards {
					_, err := e.cartService.ApplyGiftCard(ctx, currentEvent.Session, code.Code)
					if err != nil {
						e.logger.WithContext(ctx).Error("LoginEvent - customerCart ApplyGiftCard has error", code.Code, err)
					}
				}
			}

			if e.cartCache != nil {
				session := web.SessionFromContext(ctx)
				cacheID, err := e.cartCache.BuildIdentifier(ctx, session)
				if err == nil {
					err = e.cartCache.Delete(ctx, currentEvent.Session, cacheID)
					if err != nil {
						e.logger.WithContext(ctx).Error("LoginEvent - Cache Delete Error", err)
					}
				}
			}
		})
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
