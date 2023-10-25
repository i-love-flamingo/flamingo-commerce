package application

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name CartMerger --case snake

type (
	// EventReceiver handles events from other packages
	EventReceiver struct {
		logger              flamingo.Logger
		cartReceiverService Receiver
		cartCache           CartCache
		eventRouter         flamingo.EventRouter
		cartMerger          CartMerger
	}

	CartMerger interface {
		Merge(ctx context.Context, session *web.Session, guestCart cartDomain.Cart, customerCart cartDomain.Cart)
	}

	CartMergeStrategyMerge struct {
		cartService Service
		logger      flamingo.Logger
	}

	CartMergeStrategyReplace struct {
		cartService Service
		logger      flamingo.Logger
	}

	CartMergeStrategyNone struct{}

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

var _ CartMerger = &CartMergeStrategyReplace{}
var _ CartMerger = &CartMergeStrategyMerge{}
var _ CartMerger = &CartMergeStrategyNone{}

// Inject dependencies
func (e *EventReceiver) Inject(
	logger flamingo.Logger,
	cartReceiverService Receiver,
	eventRouter flamingo.EventRouter,
	cartMerger CartMerger,
	optionals *struct {
		CartCache CartCache `inject:",optional"`
	},
) {
	e.logger = logger
	e.cartReceiverService = cartReceiverService
	e.eventRouter = eventRouter
	e.cartMerger = cartMerger

	if optionals != nil {
		e.cartCache = optionals.CartCache
	}
}

func (e *EventReceiver) prepareLogger(ctx context.Context) flamingo.Logger {
	// do this on the fly to avoid wasting memory by holding a dedicated logger instance
	return e.logger.WithField(flamingo.LogKeyCategory, "cart").WithField(flamingo.LogKeySubCategory, "cart-events").WithContext(ctx)
}

//nolint:cyclop // grabbing both the customer cart and the guest cart is a complex process
func (e *EventReceiver) handleLoginEvent(ctx context.Context, loginEvent *auth.WebLoginEvent) {
	if loginEvent == nil {
		return
	}

	if loginEvent.Request == nil {
		return
	}

	if loginEvent.Identity == nil {
		return
	}

	session := loginEvent.Request.Session()
	if !e.cartReceiverService.ShouldHaveGuestCart(session) {
		return
	}

	guestCart, err := e.cartReceiverService.ViewGuestCart(ctx, session)
	if err != nil {
		e.prepareLogger(ctx).Error(fmt.Errorf("view guest cart failed: %w", err))
		return
	}

	customerCart, err := e.cartReceiverService.ViewCart(ctx, session)
	if err != nil {
		e.prepareLogger(ctx).Error(fmt.Errorf("view customer cart failed: %w", err))
		return
	}

	session.Delete(GuestCartSessionKey)

	clonedGuestCart, _ := guestCart.Clone()
	clonedCustomerCart, _ := customerCart.Clone()
	e.eventRouter.Dispatch(ctx, &PreCartMergeEvent{GuestCart: clonedGuestCart, CustomerCart: clonedCustomerCart})

	// merge the cart depending on the set strategy
	e.cartMerger.Merge(ctx, session, *guestCart, *customerCart)

	if e.cartCache != nil {
		cacheID, err := e.cartCache.BuildIdentifier(ctx, session)
		if err == nil {
			err = e.cartCache.Delete(ctx, session, cacheID)
			if err != nil {
				e.prepareLogger(ctx).Error(fmt.Errorf("can't delete cart cache entry %v: %w", cacheID, err))
			}
		}
	}

	customerCart, err = e.cartReceiverService.ViewCart(ctx, session)
	if err != nil {
		e.prepareLogger(ctx).Error(fmt.Errorf("view customer cart failed: %w", err))
		return
	}

	mergedCart, _ := customerCart.Clone()
	e.eventRouter.Dispatch(ctx, &PostCartMergeEvent{MergedCart: mergedCart})
}

// Notify should get called by flamingo Eventlogic
func (e *EventReceiver) Notify(ctx context.Context, event flamingo.Event) {
	switch currentEvent := event.(type) {
	// Clean cart cache on logout
	case *auth.WebLogoutEvent:
		if e.cartCache != nil {
			_ = e.cartCache.DeleteAll(ctx, currentEvent.Request.Session())
		}
	// Handle WebLoginEvent and merge the cart
	case *auth.WebLoginEvent:
		web.RunWithDetachedContext(ctx, func(ctx context.Context) {
			e.handleLoginEvent(ctx, currentEvent)
		})
	// Clean the cart cache when the cart should be invalidated
	case *cartDomain.InvalidateCartEvent:
		if e.cartCache != nil {
			cartID, err := e.cartCache.BuildIdentifier(ctx, currentEvent.Session)
			if err == nil {
				_ = e.cartCache.Invalidate(ctx, currentEvent.Session, cartID)
			}
		}
	}
}

func (c *CartMergeStrategyNone) Merge(_ context.Context, _ *web.Session, _ cartDomain.Cart, _ cartDomain.Cart) {
	// do nothing
}

// Inject dependencies
func (c *CartMergeStrategyReplace) Inject(
	logger flamingo.Logger,
	cartService Service,
) *CartMergeStrategyReplace {
	c.logger = logger
	c.cartService = cartService

	return c
}

//nolint:cyclop // setting all cart attributes is a complex task
func (c *CartMergeStrategyReplace) Merge(ctx context.Context, session *web.Session, guestCart cartDomain.Cart, _ cartDomain.Cart) {
	var err error

	c.logger.WithContext(ctx).Info("cleaning existing customer cart, to be able to replace the content with the guest one.")

	err = c.cartService.Clean(ctx, session)
	if err != nil {
		c.logger.WithContext(ctx).Error(fmt.Errorf("cleaning the customer cart didn't work: %w", err))
	}

	for _, delivery := range guestCart.Deliveries {
		err = c.cartService.UpdateDeliveryInfo(ctx, session, delivery.DeliveryInfo.Code, cartDomain.CreateDeliveryInfoUpdateCommand(delivery.DeliveryInfo))
		if err != nil {
			c.logger.WithContext(ctx).Error(fmt.Errorf("error during delivery info update: %w", err))
			continue
		}

		for _, item := range delivery.Cartitems {
			c.logger.WithContext(ctx).Debugf("adding guest cart item to customer cart: %v", item)
			addRequest := cartDomain.AddRequest{
				MarketplaceCode:        item.MarketplaceCode,
				Qty:                    item.Qty,
				VariantMarketplaceCode: item.VariantMarketPlaceCode,
				AdditionalData:         item.AdditionalData,
				BundleConfiguration:    item.BundleConfig,
			}

			_, err = c.cartService.AddProduct(ctx, session, delivery.DeliveryInfo.Code, addRequest)
			if err != nil {
				c.logger.WithContext(ctx).Error(fmt.Errorf("add to cart for guest item %v failed: %w", item, err))
			}
		}
	}

	if guestCart.BillingAddress != nil {
		err = c.cartService.UpdateBillingAddress(ctx, session, guestCart.BillingAddress)
		if err != nil {
			c.logger.WithContext(ctx).Error(fmt.Errorf("couldn't update billing address: %w", err))
		}
	}

	if guestCart.Purchaser != nil {
		err = c.cartService.UpdatePurchaser(ctx, session, guestCart.Purchaser, &guestCart.AdditionalData)
		if err != nil {
			c.logger.WithContext(ctx).Error(fmt.Errorf("couldn't update purchaser: %w", err))
		}
	}

	if guestCart.HasAppliedCouponCode() {
		for _, code := range guestCart.AppliedCouponCodes {
			_, err = c.cartService.ApplyVoucher(ctx, session, code.Code)
			if err != nil {
				c.logger.WithContext(ctx).Error(fmt.Errorf("couldn't apply voucher %q: %w", code.Code, err))
			}
		}
	}

	if guestCart.HasAppliedGiftCards() {
		for _, code := range guestCart.AppliedGiftCards {
			_, err = c.cartService.ApplyGiftCard(ctx, session, code.Code)
			if err != nil {
				c.logger.WithContext(ctx).Error(fmt.Errorf("couldn't apply gift card %q: %w", code.Code, err))
			}
		}
	}

	if guestCart.PaymentSelection != nil {
		err = c.cartService.UpdatePaymentSelection(ctx, session, guestCart.PaymentSelection)
		if err != nil {
			c.logger.WithContext(ctx).Error(fmt.Errorf("couldn't payment selection: %w", err))
		}
	}
}

// Inject dependencies
func (c *CartMergeStrategyMerge) Inject(
	logger flamingo.Logger,
	cartService Service,
) *CartMergeStrategyMerge {
	c.logger = logger
	c.cartService = cartService

	return c
}

//nolint:cyclop,gocognit // setting all cart attributes is a complex task
func (c *CartMergeStrategyMerge) Merge(ctx context.Context, session *web.Session, guestCart cartDomain.Cart, customerCart cartDomain.Cart) {
	var err error

	for _, delivery := range guestCart.Deliveries {
		c.logger.WithContext(ctx).Info(fmt.Sprintf("Merging delivery with code %v of guestCart with ID %v into customerCart with ID %v", delivery.DeliveryInfo.Code, guestCart.ID, customerCart.ID))

		err = c.cartService.UpdateDeliveryInfo(ctx, session, delivery.DeliveryInfo.Code, cartDomain.CreateDeliveryInfoUpdateCommand(delivery.DeliveryInfo))
		if err != nil {
			c.logger.WithContext(ctx).Error("WebLoginEvent customerCart UpdateDeliveryInfo error", err)
			continue
		}

		for _, item := range delivery.Cartitems {
			c.logger.WithContext(ctx).Debugf("Merging item from guest to user cart %v", item)
			addRequest := cartDomain.AddRequest{
				MarketplaceCode:        item.MarketplaceCode,
				Qty:                    item.Qty,
				VariantMarketplaceCode: item.VariantMarketPlaceCode,
				AdditionalData:         item.AdditionalData,
				BundleConfiguration:    item.BundleConfig,
			}

			_, err = c.cartService.AddProduct(ctx, session, delivery.DeliveryInfo.Code, addRequest)
			if err != nil {
				c.logger.WithContext(ctx).Error("WebLoginEvent - customerCart product has merge error", addRequest.MarketplaceCode, err)
			}
		}
	}

	if customerCart.BillingAddress == nil && guestCart.BillingAddress != nil {
		err = c.cartService.UpdateBillingAddress(ctx, session, guestCart.BillingAddress)
		if err != nil {
			c.logger.WithContext(ctx).Error("WebLoginEvent - customerCart UpdateBillingAddress error", err)
		}
	}

	if customerCart.Purchaser == nil && guestCart.Purchaser != nil {
		err = c.cartService.UpdatePurchaser(ctx, session, guestCart.Purchaser, &guestCart.AdditionalData)
		if err != nil {
			c.logger.WithContext(ctx).Error("WebLoginEvent - customerCart UpdatePurchaser error", err)
		}
	}

	if guestCart.HasAppliedCouponCode() {
		for _, code := range guestCart.AppliedCouponCodes {
			_, err = c.cartService.ApplyVoucher(ctx, session, code.Code)
			if err != nil {
				c.logger.WithContext(ctx).Error("WebLoginEvent - customerCart ApplyVoucher has error", code.Code, err)
			}
		}
	}

	if guestCart.HasAppliedGiftCards() {
		for _, code := range guestCart.AppliedGiftCards {
			_, err = c.cartService.ApplyGiftCard(ctx, session, code.Code)
			if err != nil {
				c.logger.WithContext(ctx).Error("WebLoginEvent - customerCart ApplyGiftCard has error", code.Code, err)
			}
		}
	}

	if customerCart.PaymentSelection == nil && guestCart.PaymentSelection != nil && customerCart.ItemCount() == 0 {
		err = c.cartService.UpdatePaymentSelection(ctx, session, guestCart.PaymentSelection)
		if err != nil {
			c.logger.WithContext(ctx).Error("WebLoginEvent - customerCart UpdatePaymentSelection error", err)
		}
	}
}
