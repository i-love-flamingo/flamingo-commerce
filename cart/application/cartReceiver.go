package application

import (
	"context"
	"errors"
	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/core/auth/domain"
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
	cartCache CartCache, // optional
) {
	cs.guestCartService = guestCartService
	cs.customerCartService = customerCartService
	cs.cartDecoratorFactory = cartDecoratorFactory
	cs.authManager = authManager
	cs.userService = userService
	cs.logger = logger
	cs.cartCache = cartCache
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

func (cs *CartReceiverService) getCartFromCache(ctx context.Context, session *web.Session, identifier CartCacheIdentifier) (*cartDomain.Cart, error) {
	if cs.cartCache == nil {
		cs.logger.Debug("no cache set")

		return nil, errors.New("no cache")
	}

	cs.logger.Debug("query cart cache %#v", identifier)

	return cs.cartCache.GetCart(ctx, session, identifier)
}

func (cs *CartReceiverService) storeCartInCache(ctx context.Context, session *web.Session, cart *cartDomain.Cart) error {
	if cs.cartCache == nil {
		return errors.New("no cache")
	}

	id, err := BuildIdentifierFromCart(cart)
	if err != nil {
		return err
	}

	return cs.cartCache.CacheCart(ctx, session, *id, cart)
}



// GetCart Get the correct Cart (either Guest or User)
func (cs *CartReceiverService) GetCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	if cs.userService.IsLoggedIn(ctx, session) {
		cacheID, err := cs.cartCache.BuildIdentifier(ctx, session)

		if err != nil {
			return nil, nil, err
		}

		cart, cacheErr := cs.getCartFromCache(ctx, session, cacheID)
		if cacheErr != nil {
			cart, err = cs.customerCartService.GetCart(ctx, cs.Auth(ctx, session), "me")
			if err != nil {
				return nil, nil, err
			}

			cs.storeCartInCache(ctx, session, cart)
		}

		behaviour, err := cs.customerCartService.GetModifyBehaviour(ctx, cs.Auth(ctx, session))
		if err != nil {
			return nil, nil, err
		}

		return cart, behaviour, nil
	}

	var guestCart *cartDomain.Cart
	var err error
	var cacheErr error

	if cs.ShouldHaveGuestCart(session) {
		cacheID, err := cs.cartCache.BuildIdentifier(ctx, session)

		if err != nil {
			return nil, nil, err
		}

		guestCart, cacheErr = cs.getCartFromCache(ctx, session, cacheID)

		if cacheErr != nil {
			guestCart, err = cs.getSessionGuestCart(ctx, session)

			if err != nil {
				//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
				cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Warn("cart.application.cartservice: GetCart - No cart in session return empty")

				//delete(ctx.Session().Values, "cart.guestid")
				return nil, nil, ErrTemporaryCartService
			}

			cs.storeCartInCache(ctx, session, guestCart)
			cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Info("guestcart not in cache - requested and passed to cache from service")
		}
	} else {
		guestCart, err = cs.guestCartService.GetNewCart(ctx)
		if err != nil {
			cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Error("cart.application.cartservice: Cannot create a new guest cart. Error %s", err)

			return nil, nil, err
		}

		cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Info("cart.application.cartservice: Requested new Guestcart %v", guestCart)
		session.Store(GuestCartSessionKey, guestCart.ID)
		cs.storeCartInCache(ctx, session, guestCart)
	}
	behaviour, err := cs.guestCartService.GetModifyBehaviour(ctx)

	if err != nil {
		return guestCart, nil, err
	}

	if guestCart == nil {
		cs.logger.WithField("category", "checkout.cartreceiver").Error("Something unexpected went wrong! No guestcart!")

		return nil, nil, errors.New("something unexpected went wrong - no guestcart")
	}

	return guestCart, behaviour, nil
}

//ViewGuestCart - ry to get the uest Cart - even if the user is logged in
func (cs *CartReceiverService) ViewGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if cs.ShouldHaveGuestCart(session) {
		guestCart, err := cs.getSessionGuestCart(ctx, session)
		if err != nil {
			//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Warn("cart.application.cartservice: GetCart - No cart in session return empty")

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

// GetSessionGuestCart
func (cs *CartReceiverService) getSessionGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if guestcartid, ok := session.Load(GuestCartSessionKey); ok {
		existingCart, err := cs.guestCartService.GetCart(ctx, guestcartid.(string))
		if err != nil {
			cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Error("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, err)
			// we seem to have an erratic session cart - remove it
			session.Delete(GuestCartSessionKey)
		}

		return existingCart, err
	}

	cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Error("No cart in session yet - getSessionGuestCart should be called only if HasSssionGuestCart returns true")

	return nil, errors.New("no cart in session yet")
}

// DecorateCart Get the correct Cart
func (cs *CartReceiverService) DecorateCart(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.DecoratedCart, error) {
	if cart == nil {
		return nil, errors.New("no cart given")
	}

	cs.logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Debug("cart.application.cartservice: Get decorated cart ")

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
