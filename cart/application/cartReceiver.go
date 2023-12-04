package application

import (
	"context"
	"errors"
	"fmt"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo/v3/framework/flamingo"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/customer/application"
)

type (
	// CartReceiverService provides methods to get the correct cart
	CartReceiverService struct {
		cartDecoratorFactory *decorator.DecoratedCartFactory
		*BaseCartReceiver
	}

	// BaseCartReceiver get undecorated carts only
	BaseCartReceiver struct {
		guestCartService    cartDomain.GuestCartService
		customerCartService cartDomain.CustomerCartService
		webIdentityService  *auth.WebIdentityService
		eventRouter         flamingo.EventRouter
		logger              flamingo.Logger
		// CartCache is optional
		cartCache CartCache
	}
)

var (
	// ErrTemporaryCartService should be returned if it is likely that the backend service will return a cart on a next try
	ErrTemporaryCartService = errors.New("the cart could not be received currently - try again later")
)

const (
	// GuestCartSessionKey is a prefix
	GuestCartSessionKey = "cart.guestid"
)

// Inject the dependencies
func (cs *CartReceiverService) Inject(
	guestCartService cartDomain.GuestCartService,
	customerCartService cartDomain.CustomerCartService,
	cartDecoratorFactory *decorator.DecoratedCartFactory,
	webIdentityService *auth.WebIdentityService,
	logger flamingo.Logger,
	eventRouter flamingo.EventRouter,
	optionals *struct {
		CartCache CartCache `inject:",optional"`
	},
) {
	cs.cartDecoratorFactory = cartDecoratorFactory
	cs.BaseCartReceiver = &BaseCartReceiver{}
	cs.BaseCartReceiver.Inject(
		guestCartService,
		customerCartService,
		webIdentityService,
		logger,
		eventRouter,
		optionals)
}

// Inject the dependencies
func (cs *BaseCartReceiver) Inject(
	guestCartService cartDomain.GuestCartService,
	customerCartService cartDomain.CustomerCartService,
	webIdentityService *auth.WebIdentityService,
	logger flamingo.Logger,
	eventRouter flamingo.EventRouter,
	optionals *struct {
		CartCache CartCache `inject:",optional"`
	},
) {
	cs.guestCartService = guestCartService
	cs.customerCartService = customerCartService
	cs.webIdentityService = webIdentityService
	cs.logger = logger.WithField("module", "cart").WithField(flamingo.LogKeyCategory, "checkout.cartreceiver")
	cs.eventRouter = eventRouter

	if optionals != nil {
		cs.cartCache = optionals.CartCache
	}
}

// RestoreCart restores a previously used guest / customer cart
// deprecated: use CartService.RestoreCart(), ensure that your cart implements the CompleteBehaviour
func (cs *BaseCartReceiver) RestoreCart(ctx context.Context, session *web.Session, cartToRestore cartDomain.Cart) (*cartDomain.Cart, error) {
	identity := cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx))
	if identity != nil {
		restoredCart, err := cs.customerCartService.RestoreCart(ctx, identity, cartToRestore)
		if err != nil {
			return nil, err
		}

		_ = cs.storeCartInCacheIfCacheIsEnabled(ctx, session, restoredCart)
		return restoredCart, nil
	}

	restoredCart, err := cs.guestCartService.RestoreCart(ctx, cartToRestore)
	if err != nil {
		return nil, err
	}

	session.Store(GuestCartSessionKey, restoredCart.ID)
	_ = cs.storeCartInCacheIfCacheIsEnabled(ctx, session, restoredCart)
	return restoredCart, nil
}

// ShouldHaveCart - checks if there should be a cart. Indicated if a call to GetCart should return a real cart
func (cs *BaseCartReceiver) ShouldHaveCart(ctx context.Context, session *web.Session) bool {
	if cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx)) != nil {
		return true
	}

	return cs.ShouldHaveGuestCart(session)
}

// ShouldHaveGuestCart - checks if there should be guest cart
func (cs *BaseCartReceiver) ShouldHaveGuestCart(session *web.Session) bool {
	_, ok := session.Load(GuestCartSessionKey)
	return ok
}

// ViewDecoratedCart  return a Cart for view
func (cs *CartReceiverService) ViewDecoratedCart(ctx context.Context, session *web.Session) (*decorator.DecoratedCart, error) {
	cart, err := cs.ViewCart(ctx, session)
	if err != nil {
		return nil, err
	}

	return cs.DecorateCart(ctx, cart)
}

// ViewDecoratedCartWithoutCache  return a Cart for view
func (cs *CartReceiverService) ViewDecoratedCartWithoutCache(ctx context.Context, session *web.Session) (*decorator.DecoratedCart, error) {
	cart, err := cs.GetCartWithoutCache(ctx, session)
	if err != nil {
		return nil, err
	}

	return cs.DecorateCart(ctx, cart)
}

// ViewCart  return a Cart for view
func (cs *BaseCartReceiver) ViewCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if cs.ShouldHaveCart(ctx, session) {
		cart, _, err := cs.GetCart(ctx, session)
		if err != nil {
			return cs.getEmptyCart(), err
		}

		return cart, nil
	}

	return cs.getEmptyCart(), nil
}

func (cs *BaseCartReceiver) storeCartInCacheIfCacheIsEnabled(ctx context.Context, session *web.Session, cart *cartDomain.Cart) error {
	if cs.cartCache == nil {
		return errors.New("no cache")
	}

	id, err := cs.cartCache.BuildIdentifier(ctx, session)
	if err != nil {
		return err
	}

	return cs.cartCache.CacheCart(ctx, session, id, cart)
}

// GetCart Get the correct Cart (either Guest or User)
func (cs *BaseCartReceiver) GetCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	if cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx)) != nil {
		return cs.getCustomerCart(ctx, session)
	}
	if cs.ShouldHaveGuestCart(session) {
		return cs.getExistingGuestCart(ctx, session)
	}
	return cs.getNewGuestCart(ctx, session)
}

// ModifyBehaviour returns the correct behaviour to modify the cart for the current user (guest/customer)
func (cs *BaseCartReceiver) ModifyBehaviour(ctx context.Context) (cartDomain.ModifyBehaviour, error) {
	identity := cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx))
	if identity != nil {
		return cs.customerCartService.GetModifyBehaviour(ctx, identity)
	}
	return cs.guestCartService.GetModifyBehaviour(ctx)
}

// getCustomerCart
func (cs *BaseCartReceiver) getCustomerCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	cart, found, err := cs.getCartFromCacheIfCacheIsEnabled(ctx, session)

	switch err {
	case nil:
	case ErrCacheIsInvalid, ErrNoCacheEntry:
		cs.logger.WithContext(ctx).Info(err)
	default:
		cs.logger.WithContext(ctx).Error(err)
	}

	identitiy := cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx))
	if identitiy == nil {
		return nil, nil, application.ErrNoIdentity
	}

	if !found {
		cart, err = cs.customerCartService.GetCart(ctx, identitiy, "me")
		if err != nil {
			return nil, nil, err
		}
		_ = cs.storeCartInCacheIfCacheIsEnabled(ctx, session, cart)
	}

	behaviour, err := cs.customerCartService.GetModifyBehaviour(ctx, identitiy)
	if err != nil {
		return nil, nil, err
	}

	return cart, behaviour, nil
}

func (cs *BaseCartReceiver) getCartFromCacheIfCacheIsEnabled(ctx context.Context, session *web.Session) (*cartDomain.Cart, bool, error) {
	if cs.cartCache == nil {
		return nil, false, nil
	}
	cacheID, err := cs.cartCache.BuildIdentifier(ctx, session)

	if err != nil {
		return nil, false, err
	}
	cs.logger.WithContext(ctx).Debug("query cart cache %#v", cacheID)

	cart, cacheErr := cs.cartCache.GetCart(ctx, session, cacheID)
	if cacheErr == ErrNoCacheEntry {
		return nil, false, nil
	}

	if cacheErr != nil {
		return nil, false, cacheErr
	}

	return cart, true, nil
}

// getExistingGuestCart
func (cs *BaseCartReceiver) getExistingGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	cart, found, err := cs.getCartFromCacheIfCacheIsEnabled(ctx, session)

	if err != nil {
		if err == ErrCacheIsInvalid {
			cs.logger.WithContext(ctx).Info(err)
		} else if err == ErrNoCacheEntry {
			cs.logger.WithContext(ctx).Info(err)
		} else {
			cs.logger.WithContext(ctx).Error(err)
		}
	}

	if err != nil || !found {
		cart, err = cs.getSessionGuestCart(ctx, session)

		if err != nil {
			// TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.logger.WithContext(ctx).Warn("GetCart - No cart in session return empty")

			// delete(ctx.Session().Values, "cart.guestid")
			return nil, nil, ErrTemporaryCartService
		}

		_ = cs.storeCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.logger.WithContext(ctx).Debug("guestcart not in cache - requested and passed to cache")
	}

	behaviour, err := cs.guestCartService.GetModifyBehaviour(ctx)
	if err != nil {
		return nil, nil, err
	}

	if cart == nil {
		cs.logger.WithContext(ctx).Error("Something unexpected went wrong! No guestcart!")

		return nil, nil, errors.New("something unexpected went wrong - no guestcart")
	}

	return cart, behaviour, nil
}

// getNewGuestCart
func (cs *BaseCartReceiver) getNewGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	guestCart, err := cs.guestCartService.GetNewCart(ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).Error("Cannot create a new guest cart. Error: ", err)
		}

		return nil, nil, err
	}

	cs.logger.WithContext(ctx).Info("Requested new guest cart: ", guestCart)
	session.Store(GuestCartSessionKey, guestCart.ID)
	_ = cs.storeCartInCacheIfCacheIsEnabled(ctx, session, guestCart)
	behaviour, err := cs.guestCartService.GetModifyBehaviour(ctx)

	if err != nil {
		return guestCart, nil, err
	}

	if guestCart == nil {
		cs.logger.WithContext(ctx).Error("Something unexpected went wrong! No guest cart!")

		return nil, nil, errors.New("something unexpected went wrong - no guest cart")
	}

	return guestCart, behaviour, nil
}

// GetCartWithoutCache - forces to get the cart without cache
func (cs *BaseCartReceiver) GetCartWithoutCache(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	// Invalidate cart cache
	if cs.eventRouter != nil {
		cs.eventRouter.Dispatch(ctx, &cartDomain.InvalidateCartEvent{Session: session})
	}

	identitiy := cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx))
	if identitiy != nil {
		return cs.customerCartService.GetCart(ctx, identitiy, "me")
	}

	if cs.ShouldHaveGuestCart(session) {
		return cs.getSessionGuestCart(ctx, session)
	}

	return cs.guestCartService.GetNewCart(ctx)

}

// ViewGuestCart try to get the guest Cart - even if the user is logged in
func (cs *BaseCartReceiver) ViewGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if cs.ShouldHaveGuestCart(session) {
		guestCart, err := cs.getSessionGuestCart(ctx, session)
		if err != nil {
			// TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.logger.WithContext(ctx).Warn("GetCart - No cart in session return empty")

			return nil, ErrTemporaryCartService
		}

		return guestCart, nil
	}

	return cs.getEmptyCart(), nil
}

// DeleteSavedSessionGuestCartID deletes a guest cart Key from the Session Values
func (cs *CartService) DeleteSavedSessionGuestCartID(session *web.Session) error {
	session.Delete(GuestCartSessionKey)

	// TODO - trigger backend also to be able to delete the cart there ( cs.GuestCartService.DeleteCart())
	return nil
}

// getSessionGuestCart
func (cs *BaseCartReceiver) getSessionGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if guestcartid, ok := session.Load(GuestCartSessionKey); ok {
		existingCart, err := cs.guestCartService.GetCart(ctx, guestcartid.(string))
		if err != nil {
			cs.logger.WithContext(ctx).Warn(fmt.Sprintf("Guest cart with ID %q cannot be retrieved. Error: %s", guestcartid, err))
			// we seem to have an erratic session cart - remove it
			session.Delete(GuestCartSessionKey)
		}

		return existingCart, err
	}

	cs.logger.WithContext(ctx).Error("No cart in session yet - getSessionGuestCart should be called only if HasSessionGuestCart returns true")

	return nil, errors.New("no cart in session yet")
}

// DecorateCart Get the correct Cart
func (cs *CartReceiverService) DecorateCart(ctx context.Context, cart *cartDomain.Cart) (*decorator.DecoratedCart, error) {
	if cart == nil {
		return nil, errors.New("no cart given")
	}

	return cs.cartDecoratorFactory.Create(ctx, *cart), nil
}

// GetDecoratedCart Get the correct Cart
func (cs *CartReceiverService) GetDecoratedCart(ctx context.Context, session *web.Session) (*decorator.DecoratedCart, cartDomain.ModifyBehaviour, error) {
	cart, behaviour, err := cs.GetCart(ctx, session)
	if err != nil {
		return nil, nil, err
	}

	return cs.cartDecoratorFactory.Create(ctx, *cart), behaviour, nil
}

func (cs *BaseCartReceiver) getEmptyCart() *cartDomain.Cart {
	return &cartDomain.Cart{}
}
