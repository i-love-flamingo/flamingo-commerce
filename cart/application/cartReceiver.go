package application

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	//CartReceiverService provides methods to get the correct cart
	CartReceiverService struct {
		guestCartService     cartDomain.GuestCartService
		customerCartService  cartDomain.CustomerCartService
		cartDecoratorFactory *cartDomain.DecoratedCartFactory
		authManager          *authApplication.AuthManager
		userService          authApplication.UserServiceInterface
		eventRouter          flamingo.EventRouter
		logger               flamingo.Logger
		// CartCache is optional
		cartCache CartCache
	}
)

var (
	//ErrTemporaryCartService - should be returned if it is likely that the backend service will return a cart on a next try
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
	cartDecoratorFactory *cartDomain.DecoratedCartFactory,
	authManager *authApplication.AuthManager,
	userService authApplication.UserServiceInterface,
	logger flamingo.Logger,
	eventRouter flamingo.EventRouter,
	optionals *struct {
		CartCache CartCache `inject:",optional"`
	},
) {
	cs.guestCartService = guestCartService
	cs.customerCartService = customerCartService
	cs.cartDecoratorFactory = cartDecoratorFactory
	cs.authManager = authManager
	cs.userService = userService
	cs.logger = logger.WithField("module", "cart").WithField(flamingo.LogKeyCategory, "checkout.cartreceiver")
	cs.eventRouter = eventRouter
	if optionals != nil {
		cs.cartCache = optionals.CartCache
	}
}

// Auth tries to retrieve the authentication context for a active session
func (cs *CartReceiverService) Auth(c context.Context, session *web.Session) domain.Auth {
	ts, _ := cs.authManager.TokenSource(c, session)
	idToken, _ := cs.authManager.IDToken(c, session)

	return domain.Auth{
		TokenSource: ts,
		IDToken:     idToken,
	}
}

// IsLoggedIn returns the logged in state
func (cs *CartReceiverService) IsLoggedIn(ctx context.Context, session *web.Session) bool {
	return cs.userService.IsLoggedIn(ctx, session)
}

// ShouldHaveCart - checks if there should be a cart. Indicated if a call to GetCart should return a real cart
func (cs *CartReceiverService) ShouldHaveCart(ctx context.Context, session *web.Session) bool {
	if cs.userService.IsLoggedIn(ctx, session) {
		return true
	}

	return cs.ShouldHaveGuestCart(session)
}

// ShouldHaveGuestCart - checks if there should be guest cart
func (cs *CartReceiverService) ShouldHaveGuestCart(session *web.Session) bool {
	_, ok := session.Load(GuestCartSessionKey)
	return ok
}

// ViewDecoratedCart  return a Cart for view
func (cs *CartReceiverService) ViewDecoratedCart(ctx context.Context, session *web.Session) (*cartDomain.DecoratedCart, error) {
	cart, err := cs.ViewCart(ctx, session)
	if err != nil {
		return nil, err
	}

	return cs.DecorateCart(ctx, cart)
}

// ViewDecoratedCartWithoutCache  return a Cart for view
func (cs *CartReceiverService) ViewDecoratedCartWithoutCache(ctx context.Context, session *web.Session) (*cartDomain.DecoratedCart, error) {
	cart, err := cs.GetCartWithoutCache(ctx, session)
	if err != nil {
		return nil, err
	}

	return cs.DecorateCart(ctx, cart)
}

// ViewCart  return a Cart for view
func (cs *CartReceiverService) ViewCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if cs.ShouldHaveCart(ctx, session) {
		cart, _, err := cs.GetCart(ctx, session)
		if err != nil {
			return cs.getEmptyCart(), err
		}

		return cart, nil
	}

	return cs.getEmptyCart(), nil
}

func (cs *CartReceiverService) storeCartInCacheIfCacheIsEnabled(ctx context.Context, session *web.Session, cart *cartDomain.Cart) error {
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
func (cs *CartReceiverService) GetCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	if cs.userService.IsLoggedIn(ctx, session) {
		return cs.getCustomerCart(ctx, session)
	}
	if cs.ShouldHaveGuestCart(session) {
		return cs.getExistingGuestCart(ctx, session)
	}
	return cs.getNewGuestCart(ctx, session)
}

// getCustomerCart
func (cs *CartReceiverService) getCustomerCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {

	cart, found, err := cs.getCartFromCacheIfCacheIsEnabled(ctx, session)

	if err != nil {
		cs.logger.Error(err)
	}

	if err != nil || !found {
		cart, err = cs.customerCartService.GetCart(ctx, cs.Auth(ctx, session), "me")
		if err != nil {
			return nil, nil, err
		}
		cs.storeCartInCacheIfCacheIsEnabled(ctx, session, cart)
	}

	behaviour, err := cs.customerCartService.GetModifyBehaviour(ctx, cs.Auth(ctx, session))
	if err != nil {
		return nil, nil, err
	}

	return cart, behaviour, nil
}

func (cs *CartReceiverService) getCartFromCacheIfCacheIsEnabled(ctx context.Context, session *web.Session) (*cartDomain.Cart, bool, error) {
	if cs.cartCache == nil {
		return nil, false, nil
	}
	cacheID, err := cs.cartCache.BuildIdentifier(ctx, session)

	if err != nil {
		return nil, false, err
	}
	cs.logger.Debug("query cart cache %#v", cacheID)
	cart, cacheErr := cs.cartCache.GetCart(ctx, session, cacheID)
	if cacheErr != nil {
		return nil, false, cacheErr
	}
	return cart, true, nil
}

// getExistingGuestCart
func (cs *CartReceiverService) getExistingGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	cart, found, err := cs.getCartFromCacheIfCacheIsEnabled(ctx, session)

	if err != nil {
		cs.logger.Error(err)
	}
	if err != nil || !found {
		cart, err = cs.getSessionGuestCart(ctx, session)

		if err != nil {
			//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.logger.Warn("cart.application.cartservice: GetCart - No cart in session return empty")

			//delete(ctx.Session().Values, "cart.guestid")
			return nil, nil, ErrTemporaryCartService
		}

		cs.storeCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.logger.Debug("guestcart not in cache - requested and passed to cache")
	}
	behaviour, err := cs.guestCartService.GetModifyBehaviour(ctx)

	if err != nil {
		return nil, nil, err
	}

	if cart == nil {
		cs.logger.Error("Something unexpected went wrong! No guestcart!")

		return nil, nil, errors.New("something unexpected went wrong - no guestcart")
	}

	return cart, behaviour, nil
}

// getNewGuestCart
func (cs *CartReceiverService) getNewGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	guestCart, err := cs.guestCartService.GetNewCart(ctx)
	if err != nil {
		cs.logger.Error("cart.application.cartservice: Cannot create a new guest cart. Error %s", err)

		return nil, nil, err
	}

	cs.logger.Info("cart.application.cartservice: Requested new Guestcart %v", guestCart)
	session.Store(GuestCartSessionKey, guestCart.ID)
	cs.storeCartInCacheIfCacheIsEnabled(ctx, session, guestCart)
	behaviour, err := cs.guestCartService.GetModifyBehaviour(ctx)

	if err != nil {
		return guestCart, nil, err
	}

	if guestCart == nil {
		cs.logger.Error("Something unexpected went wrong! No guestcart!")

		return nil, nil, errors.New("something unexpected went wrong - no guestcart")
	}

	return guestCart, behaviour, nil
}

// GetCartWithoutCache - forces to get the cart without cache
func (cs *CartReceiverService) GetCartWithoutCache(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	// Invalidate cart cache
	if cs.eventRouter != nil {
		cs.eventRouter.Dispatch(ctx, &cart.InvalidateCartEvent{Session: session})
	}

	if cs.userService.IsLoggedIn(ctx, session) {
		return cs.customerCartService.GetCart(ctx, cs.Auth(ctx, session), "me")
	}

	if cs.ShouldHaveGuestCart(session) {
		return cs.getSessionGuestCart(ctx, session)
	}

	return cs.guestCartService.GetNewCart(ctx)

}

//ViewGuestCart - ry to get the uest Cart - even if the user is logged in
func (cs *CartReceiverService) ViewGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if cs.ShouldHaveGuestCart(session) {
		guestCart, err := cs.getSessionGuestCart(ctx, session)
		if err != nil {
			//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.logger.Warn("cart.application.cartservice: GetCart - No cart in session return empty")

			return nil, ErrTemporaryCartService
		}

		return guestCart, nil
	}

	return cs.getEmptyCart(), nil
}

// DeleteSavedSessionGuestCartID deletes a guest cart Key from the Session Values
func (cs *CartService) DeleteSavedSessionGuestCartID(session *web.Session) error {
	session.Delete(GuestCartSessionKey)

	//TODO - trigger backend also to be able to delete the cart there ( cs.GuestCartService.DeleteCart())
	return nil
}

// getSessionGuestCart
func (cs *CartReceiverService) getSessionGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if guestcartid, ok := session.Load(GuestCartSessionKey); ok {
		existingCart, err := cs.guestCartService.GetCart(ctx, guestcartid.(string))
		if err != nil {
			cs.logger.Error("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, err)
			// we seem to have an erratic session cart - remove it
			session.Delete(GuestCartSessionKey)
		}

		return existingCart, err
	}

	cs.logger.Error("No cart in session yet - getSessionGuestCart should be called only if HasSssionGuestCart returns true")

	return nil, errors.New("no cart in session yet")
}

// DecorateCart Get the correct Cart
func (cs *CartReceiverService) DecorateCart(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.DecoratedCart, error) {
	if cart == nil {
		return nil, errors.New("no cart given")
	}

	cs.logger.Debug("cart.application.cartservice: Get decorated cart ")

	return cs.cartDecoratorFactory.Create(ctx, *cart), nil
}

// GetDecoratedCart Get the correct Cart
func (cs *CartReceiverService) GetDecoratedCart(ctx context.Context, session *web.Session) (*cartDomain.DecoratedCart, cartDomain.ModifyBehaviour, error) {
	cart, behaviour, err := cs.GetCart(ctx, session)
	if err != nil {
		return nil, nil, err
	}

	return cs.cartDecoratorFactory.Create(ctx, *cart), behaviour, nil
}

func (cs *CartReceiverService) getEmptyCart() *cartDomain.Cart {
	return &cartDomain.Cart{}
}
