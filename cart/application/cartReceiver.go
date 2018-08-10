package application

import (
	"context"
	"errors"

	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
	authApplication "flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
)

type (
	//CartReceiverService provides methods to get the correct cart
	CartReceiverService struct {
		GuestCartService     cartDomain.GuestCartService      `inject:""`
		CustomerCartService  cartDomain.CustomerCartService   `inject:""`
		CartDecoratorFactory *cartDomain.DecoratedCartFactory `inject:""`
		AuthManager          *authApplication.AuthManager     `inject:""`
		UserService          *authApplication.UserService     `inject:""`
		Logger               flamingo.Logger                  `inject:""`
		CartCache            CartCache                        `inject:",optional"`
	}
)

var (
	//TemporaryCartServiceError - should be returned if it is likely that the backend service will return a cart on a next try
	TemporaryCartServiceError error = errors.New("The cart could not be received currently - try again later")
)

const (
	GuestCartSessionKey = "cart.guestid"
)

// Auth tries to retrieve the authentication context for a active session
func (cs *CartReceiverService) Auth(c context.Context, session *sessions.Session) cartDomain.Auth {
	ts, _ := cs.AuthManager.TokenSource(c, session)
	idToken, _ := cs.AuthManager.IDToken(c, session)

	return cartDomain.Auth{
		TokenSource: ts,
		IDToken:     idToken,
	}
}

// ShouldHaveCart - checks if there should be a cart. Indicated if a call to GetCart should return a real cart
func (cs *CartReceiverService) ShouldHaveCart(ctx context.Context, session *sessions.Session) bool {
	if cs.UserService.IsLoggedIn(ctx, session) {
		return true
	}
	return cs.ShouldHaveGuestCart(session)
}

// ShouldHaveGuestCart - checks if there should be guest cart
func (cs *CartReceiverService) ShouldHaveGuestCart(session *sessions.Session) bool {
	if _, ok := session.Values[GuestCartSessionKey]; ok {
		return true
	}
	return false
}

// ViewDecoratedCart  return a Cart for view
func (cs *CartReceiverService) ViewDecoratedCart(ctx context.Context, session *sessions.Session) (*cartDomain.DecoratedCart, error) {
	cart, e := cs.ViewCart(ctx, session)
	if e != nil {
		return nil, e
	}
	return cs.DecorateCart(ctx, cart)
}

// ViewCart  return a Cart for view
func (cs *CartReceiverService) ViewCart(ctx context.Context, session *sessions.Session) (*cartDomain.Cart, error) {
	if cs.ShouldHaveCart(ctx, session) {
		cart, _, err := cs.GetCart(ctx, session)
		if err != nil {
			return cs.getEmptyCart(), err
		}
		return cart, nil
	}
	return cs.getEmptyCart(), nil
}

func (cs *CartReceiverService) getCartFromCache(ctx context.Context, identifier CartCacheIdentifier) (*cartDomain.Cart, error) {
	if cs.CartCache == nil {
		cs.Logger.Debug("no cache set")
		return nil, errors.New("no cache")
	}
	cs.Logger.Debug("query cart cache %#v", identifier)
	return cs.CartCache.GetCart(web.ToContext(ctx), identifier)
}

func (cs *CartReceiverService) storeCartInCache(ctx context.Context, cart *cartDomain.Cart) error {
	if cs.CartCache == nil {
		return errors.New("no cache")
	}
	id, err := BuildIdentifierFromCart(cart)
	if err != nil {
		return err
	}
	return cs.CartCache.CacheCart(web.ToContext(ctx), *id, cart)
}

// GetCart Get the correct Cart (either Guest or User)
func (cs *CartReceiverService) GetCart(ctx context.Context, session *sessions.Session) (*cartDomain.Cart, cartDomain.CartBehaviour, error) {
	if cs.UserService.IsLoggedIn(ctx, session) {
		cacheId := CartCacheIdentifier{
			CustomerId:     cs.Auth(ctx, session).IDToken.Subject,
			IsCustomerCart: true,
		}
		var err error
		cart, cacheErr := cs.getCartFromCache(ctx, cacheId)
		if cacheErr != nil {
			cart, err = cs.CustomerCartService.GetCart(ctx, cs.Auth(ctx, session), "me")
			if err != nil {
				return nil, nil, err
			}
			cs.storeCartInCache(ctx, cart)
		}
		behaviour, err := cs.CustomerCartService.GetCartOrderBehaviour(ctx, cs.Auth(ctx, session))
		if err != nil {
			return nil, nil, err
		}
		return cart, behaviour, nil
	}

	var guestCart *cartDomain.Cart
	var err error
	var cacheErr error
	if cs.ShouldHaveGuestCart(session) {
		guestcartid, ok := session.Values[GuestCartSessionKey]
		if !ok {
			panic("Fatal - ShouldHaveGuestCart returned true but got no GuestCartSessionKey?")
		}
		guestcartidString, ok := guestcartid.(string)
		if !ok {
			panic("Fatal - ShouldHaveGuestCart returned true but got no GuestCartSessionKey string")
		}
		cacheId := CartCacheIdentifier{
			GuestCartId: guestcartidString,
		}
		guestCart, cacheErr = cs.getCartFromCache(ctx, cacheId)
		if cacheErr != nil {
			guestCart, err = cs.getSessionGuestCart(ctx, session)
			if err != nil {
				//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
				cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Warn("cart.application.cartservice: GetCart - No cart in session return empty")
				//delete(ctx.Session().Values, "cart.guestid")
				return nil, nil, TemporaryCartServiceError
			}
			cs.storeCartInCache(ctx, guestCart)
			cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Info("guestcart not in cache - requested and passed to cache from service")
		}
	} else {
		guestCart, err = cs.GuestCartService.GetNewCart(ctx)
		if err != nil {
			cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Error("cart.application.cartservice: Cannot create a new guest cart. Error %s", err)
			return nil, nil, err
		}
		cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Info("cart.application.cartservice: Requested new Guestcart %v", guestCart)
		session.Values[GuestCartSessionKey] = guestCart.ID
		cs.storeCartInCache(ctx, guestCart)
	}
	behaviour, err := cs.GuestCartService.GetCartOrderBehaviour(ctx)

	if err != nil {
		return guestCart, nil, err
	}
	if guestCart == nil {
		cs.Logger.WithField("category", "checkout.cartreceiver").Error("Something unexpected went wrong! No guestcart!")
		return nil, nil, errors.New("Something unexpected went wrong! No guestcart!")
	}
	return guestCart, behaviour, nil
}

//ViewGuestCart - ry to get the uest Cart - even if the user is logged in
func (cs *CartReceiverService) ViewGuestCart(ctx context.Context, session *sessions.Session) (*cartDomain.Cart, error) {
	if cs.ShouldHaveGuestCart(session) {
		guestCart, err := cs.getSessionGuestCart(ctx, session)
		if err != nil {
			//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Warn("cart.application.cartservice: GetCart - No cart in session return empty")
			return nil, TemporaryCartServiceError
		}
		return guestCart, nil
	} else {
		return cs.getEmptyCart(), nil
	}
}

// GetSessionGuestCart
func (cs *CartService) DeleteSavedSessionGuestCartId(session *sessions.Session) error {
	delete(session.Values, GuestCartSessionKey)
	//TODO - trigger backend also to be able to delete the cart there ( cs.GuestCartService.DeleteCart())
	return nil
}

// GetSessionGuestCart
func (cs *CartReceiverService) getSessionGuestCart(ctx context.Context, session *sessions.Session) (*cartDomain.Cart, error) {
	if guestcartid, ok := session.Values[GuestCartSessionKey]; ok {
		existingCart, err := cs.GuestCartService.GetCart(ctx, guestcartid.(string))
		if err != nil {
			cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Error("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, err)
		}
		return existingCart, err
	}
	cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Error("No cart in session yet - getSessionGuestCart should be called only if HasSssionGuestCart returns true")
	return nil, errors.New("No cart in session yet")
}

// DecorateCart Get the correct Cart
func (cs *CartReceiverService) DecorateCart(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.DecoratedCart, error) {
	if cart == nil {
		return nil, errors.New("no cart given")
	}
	cs.Logger.WithField(flamingo.LogKeyCategory, "checkout.cartreceiver").Debug("cart.application.cartservice: Get decorated cart ")
	return cs.CartDecoratorFactory.Create(ctx, *cart), nil
}

// GetDecoratedCart Get the correct Cart
func (cs *CartReceiverService) GetDecoratedCart(ctx context.Context, session *sessions.Session) (*cartDomain.DecoratedCart, cartDomain.CartBehaviour, error) {
	cart, behaviour, err := cs.GetCart(ctx, session)
	if err != nil {
		return nil, nil, err
	}
	return cs.CartDecoratorFactory.Create(ctx, *cart), behaviour, nil
}

func (cs *CartReceiverService) getEmptyCart() *cartDomain.Cart {
	return &cartDomain.Cart{}
}
