package application

import (
	"errors"

	authApplication "go.aoe.com/flamingo/core/auth/application"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
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
func (cs *CartReceiverService) Auth(c web.Context) cartDomain.Auth {
	ts, _ := cs.AuthManager.TokenSource(c)
	idToken, _ := cs.AuthManager.IDToken(c)

	return cartDomain.Auth{
		TokenSource: ts,
		IDToken:     idToken,
	}
}

// ShouldHaveCart - checks if there should be a cart. Indicated if a call to GetCart should return a real cart
func (cs *CartReceiverService) ShouldHaveCart(ctx web.Context) bool {
	if cs.UserService.IsLoggedIn(ctx) {
		return true
	}
	return cs.ShouldHaveGuestCart(ctx)
}

// ShouldHaveGuestCart - checks if there should be guest cart
func (cs *CartReceiverService) ShouldHaveGuestCart(ctx web.Context) bool {
	if _, ok := ctx.Session().Values[GuestCartSessionKey]; ok {
		return true
	}
	return false
}

// ViewDecoratedCart  return a Cart for view
func (cs *CartReceiverService) ViewDecoratedCart(ctx web.Context) (*cartDomain.DecoratedCart, error) {
	cart, e := cs.ViewCart(ctx)
	if e != nil {
		return nil, e
	}
	return cs.DecorateCart(ctx, cart)
}

// ViewCart  return a Cart for view
func (cs *CartReceiverService) ViewCart(ctx web.Context) (*cartDomain.Cart, error) {
	if cs.ShouldHaveCart(ctx) {
		cart, _, err := cs.GetCart(ctx)
		if err != nil {
			return cs.getEmptyCart(), err
		}
		return cart, nil
	}
	return cs.getEmptyCart(), nil
}

// GetCart Get the correct Cart (either Guest or User)
func (cs *CartReceiverService) GetCart(ctx web.Context) (*cartDomain.Cart, cartDomain.CartBehaviour, error) {
	if cs.UserService.IsLoggedIn(ctx) {
		cart, err := cs.CustomerCartService.GetCart(ctx, cs.Auth(ctx), "me")
		if err != nil {
			return nil, nil, err
		}
		behaviour, err := cs.CustomerCartService.GetCartOrderBehaviour(ctx, cs.Auth(ctx))
		if err != nil {
			return nil, nil, err
		}
		return cart, behaviour, nil
	}

	var guestCart *cartDomain.Cart
	var err error
	if cs.ShouldHaveGuestCart(ctx) {
		guestCart, err = cs.getSessionGuestCart(ctx)
		if err != nil {
			//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.Logger.WithField("category", "checkout.cartservice").Warnf("cart.application.cartservice: GetCart - No cart in session return empty")
			//delete(ctx.Session().Values, "cart.guestid")
			return nil, nil, TemporaryCartServiceError
		}
	} else {
		guestCart, err = cs.GuestCartService.GetNewCart(ctx)
		if err != nil {
			cs.Logger.WithField("category", "checkout.cartservice").Errorf("cart.application.cartservice: Cannot create a new guest cart. Error %s", err)
			return nil, nil, err
		}
		cs.Logger.WithField("category", "checkout.cartservice").Infof("cart.application.cartservice: Requested new Guestcart %v", guestCart)
		ctx.Session().Values[GuestCartSessionKey] = guestCart.ID
	}
	behaviour, err := cs.GuestCartService.GetCartOrderBehaviour(ctx)
	if err != nil {
		return guestCart, nil, err
	}
	return guestCart, behaviour, nil
}

//ViewGuestCart - ry to get the uest Cart - even if the user is logged in
func (cs *CartReceiverService) ViewGuestCart(ctx web.Context) (*cartDomain.Cart, error) {
	if cs.ShouldHaveGuestCart(ctx) {
		guestCart, err := cs.getSessionGuestCart(ctx)
		if err != nil {
			//TODO - decide on recoverable errors (where we should communicate "try again" / and not recoverable (where we should clean up guest cart in session and try to get a new one)
			cs.Logger.WithField("category", "checkout.cartservice").Warnf("cart.application.cartservice: GetCart - No cart in session return empty")
			return nil, TemporaryCartServiceError
		}
		return guestCart, nil
	} else {
		return cs.getEmptyCart(), nil
	}
}

// GetSessionGuestCart
func (cs *CartService) DeleteSessionGuestCart(ctx web.Context) error {
	delete(ctx.Session().Values, GuestCartSessionKey)
	//TODO - trigger backend also to be able to delete the cart there ( cs.GuestCartService.DeleteCart())
	return nil
}

// GetSessionGuestCart
func (cs *CartReceiverService) getSessionGuestCart(ctx web.Context) (*cartDomain.Cart, error) {
	if guestcartid, ok := ctx.Session().Values[GuestCartSessionKey]; ok {
		existingCart, err := cs.GuestCartService.GetCart(ctx, guestcartid.(string))
		if err != nil {
			cs.Logger.WithField("category", "checkout.cartservice").Errorf("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, err)
		}
		return existingCart, err
	}

	return nil, errors.New("No cart in session yet")
}

// DecorateCart Get the correct Cart
func (cs *CartReceiverService) DecorateCart(ctx web.Context, cart *cartDomain.Cart) (*cartDomain.DecoratedCart, error) {
	cs.Logger.WithField("category", "checkout.cartservice").Info("cart.application.cartservice: Get decorated cart ")
	if cart == nil {
		return nil, errors.New("no cart given")
	}
	return cs.CartDecoratorFactory.Create(ctx, *cart), nil
}

// GetDecoratedCart Get the correct Cart
func (cs *CartReceiverService) GetDecoratedCart(ctx web.Context) (*cartDomain.DecoratedCart, cartDomain.CartBehaviour, error) {
	cart, behaviour, err := cs.GetCart(ctx)
	if err != nil {
		return nil, nil, err
	}
	return cs.CartDecoratorFactory.Create(ctx, *cart), behaviour, nil
}

func (cs *CartReceiverService) getEmptyCart() *cartDomain.Cart {
	return &cartDomain.Cart{}
}
