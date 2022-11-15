package application

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo/v3/framework/flamingo"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// EventReceiver - handles events from other packages
	EventReceiver struct {
		logger              flamingo.Logger
		cartService         *CartService
		cartReceiverService *CartReceiverService
		cartCache           CartCache
		webIdentityService  *auth.WebIdentityService
		eventRouter         flamingo.EventRouter
	}

	// PreCartMergeEvent is dispatched after getting the (current) guest cart and the customer cart before merging
	PreCartMergeEvent struct {
		GuestCart    cartDomain.Cart
		CustomerCart cartDomain.Cart
	}

	// PostCartMergeEvent is dispatched after merging the guest cart and the customer cart
	PostCartMergeEvent struct {
		MergedCart cartDomain.Cart
	}
)

// Inject dependencies
func (e *EventReceiver) Inject(
	logger flamingo.Logger,
	cartService *CartService,
	cartReceiverService *CartReceiverService,
	webIdentityService *auth.WebIdentityService,
	eventRouter flamingo.EventRouter,
	optionals *struct {
		CartCache CartCache `inject:",optional"`
	},
) {
	e.logger = logger.WithField(flamingo.LogKeyCategory, "cart").WithField(flamingo.LogKeySubCategory, "cart-events")
	e.cartService = cartService
	e.cartReceiverService = cartReceiverService
	e.webIdentityService = webIdentityService
	e.eventRouter = eventRouter
	if optionals != nil {
		e.cartCache = optionals.CartCache
	}
}

// Notify should get called by flamingo Eventlogic
func (e *EventReceiver) Notify(ctx context.Context, event flamingo.Event) {
	switch currentEvent := event.(type) {
	// Handle Logout
	case *auth.WebLogoutEvent:
		if e.cartCache != nil {
			_ = e.cartCache.DeleteAll(ctx, currentEvent.Request.Session())
		}
	// Handle WebLoginEvent and Merge Cart
	case *auth.WebLoginEvent:
		web.RunWithDetachedContext(ctx, func(ctx context.Context) {
			if currentEvent == nil {
				return
			}
			session := currentEvent.Request.Session()
			if !e.cartReceiverService.ShouldHaveGuestCart(session) {
				return
			}
			guestCart, err := e.cartReceiverService.ViewGuestCart(ctx, session)
			if err != nil {
				e.logger.WithContext(ctx).Error("WebLoginEvent - GuestCart cannot be received %v", err)
				return
			}
			identity := e.webIdentityService.Identify(ctx, currentEvent.Request)
			if identity == nil {
				e.logger.WithContext(ctx).Error("Received WebLoginEvent but user is not logged in!")
				return
			}
			customerCart, err := e.cartReceiverService.ViewCart(ctx, session)
			if err != nil {
				e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart cannot be received %v", err)
				return
			}
			err = e.cartService.DeleteSavedSessionGuestCartID(session)
			if err != nil {
				e.logger.WithContext(ctx).Error("WebLoginEvent - DeleteSavedSessionGuestCartID Error", err)
			}

			clonedGuestCart, _ := guestCart.Clone()
			clonedCustomerCart, _ := customerCart.Clone()
			e.eventRouter.Dispatch(ctx, &PreCartMergeEvent{GuestCart: clonedGuestCart, CustomerCart: clonedCustomerCart})

			for _, d := range guestCart.Deliveries {
				e.logger.WithContext(ctx).Info(fmt.Sprintf("Merging delivery with code %v of guestCart with ID %v into customerCart with ID %v", d.DeliveryInfo.Code, guestCart.ID, customerCart.ID))
				err := e.cartService.UpdateDeliveryInfo(ctx, session, d.DeliveryInfo.Code, cartDomain.CreateDeliveryInfoUpdateCommand(d.DeliveryInfo))
				if err != nil {
					e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart UpdateDeliveryInfo error", err)
					continue
				}
				for _, item := range d.Cartitems {
					e.logger.WithContext(ctx).Debugf("Merging item from guest to user cart %v", item)
					addRequest := e.cartService.BuildAddRequest(ctx, item.MarketplaceCode, item.VariantMarketPlaceCode, item.Qty, item.AdditionalData)
					_, err := e.cartService.AddProduct(ctx, session, d.DeliveryInfo.Code, addRequest)
					if err != nil {
						e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart product has merge error", addRequest.MarketplaceCode, err)
					}
				}
			}
			if customerCart.BillingAddress == nil && guestCart.BillingAddress != nil {
				err := e.cartService.UpdateBillingAddress(ctx, session, guestCart.BillingAddress)
				if err != nil {
					e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart UpdateBillingAddress error", err)
				}
			}
			if customerCart.Purchaser == nil && guestCart.Purchaser != nil {
				err := e.cartService.UpdatePurchaser(ctx, session, guestCart.Purchaser, &guestCart.AdditionalData)
				if err != nil {
					e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart UpdatePurchaser error", err)
				}
			}
			if guestCart.HasAppliedCouponCode() {
				for _, code := range guestCart.AppliedCouponCodes {
					customerCart, err = e.cartService.ApplyVoucher(ctx, session, code.Code)
					if err != nil {
						e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart ApplyVoucher has error", code.Code, err)
					}
				}
			}
			if guestCart.HasAppliedGiftCards() {
				for _, code := range guestCart.AppliedGiftCards {
					customerCart, err = e.cartService.ApplyGiftCard(ctx, session, code.Code)
					if err != nil {
						e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart ApplyGiftCard has error", code.Code, err)
					}
				}
			}
			if customerCart.PaymentSelection == nil && guestCart.PaymentSelection != nil && customerCart.ItemCount() == 0 {
				err := e.cartService.UpdatePaymentSelection(ctx, session, guestCart.PaymentSelection)
				if err != nil {
					e.logger.WithContext(ctx).Error("WebLoginEvent - customerCart UpdatePaymentSelection error", err)
				}
			}

			if e.cartCache != nil {
				session := web.SessionFromContext(ctx)
				cacheID, err := e.cartCache.BuildIdentifier(ctx, session)
				if err == nil {
					err = e.cartCache.Delete(ctx, session, cacheID)
					if err != nil {
						e.logger.WithContext(ctx).Error("WebLoginEvent - Cache Delete Error", err)
					}
				}
			}

			mergedCart, _ := customerCart.Clone()
			e.eventRouter.Dispatch(ctx, &PostCartMergeEvent{MergedCart: mergedCart})
		})
	// Handle Event to Invalidate the Cart Cache
	case *cartDomain.InvalidateCartEvent:
		if e.cartCache != nil {
			cartID, err := e.cartCache.BuildIdentifier(ctx, currentEvent.Session)
			if err == nil {
				_ = e.cartCache.Invalidate(ctx, currentEvent.Session, cartID)
			}
		}
	}
}
