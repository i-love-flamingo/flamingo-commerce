package application

import (
	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/auth/application"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

// CartService application struct
type (
	CartService struct {
		GuestCartService     cartDomain.GuestCartService      `inject:""`
		CustomerCartService  cartDomain.CustomerCartService   `inject:""`
		CartDecoratorFactory *cartDomain.DecoratedCartFactory `inject:""`
		ProductService       productDomain.ProductService     `inject:""`
		Logger               flamingo.Logger                  `inject:""`
		CartValidator        cartDomain.CartValidator         `inject:",optional"`
		AuthManager          *application.AuthManager         `inject:""`
		UserService          *application.UserService         `inject:""`

		ItemValidator cartDomain.ItemValidator `inject:",optional"`

		//DefaultDeliveryMethodForValidation - used for calling the CartValidator (this is something that might get obsolete if the Cart and the CartItems have theire Deliverymethod "saved")
		DefaultDeliveryMethodForValidation string `inject:"config:cart.validation.defaultDeliveryMethod,optional"`
	}
)

// Auth tries to retrieve the authentication context for a active session
func (cs *CartService) Auth(c web.Context) cartDomain.Auth {
	ts, _ := cs.AuthManager.TokenSource(c)
	idToken, _ := cs.AuthManager.IDToken(c)

	return cartDomain.Auth{
		TokenSource: ts,
		IDToken:     idToken,
	}
}

// GetCart Get the correct Cart
func (cs *CartService) GetCart(ctx web.Context) (cartDomain.Cart, error) {
	if cs.isLoggedIn(ctx) {
		return cs.CustomerCartService.GetCart(ctx, cs.Auth(ctx), "me")
	}

	guestCart, err := cs.GetSessionGuestCart(ctx)
	if err != nil {
		cs.Logger.Warn("cart.application.cartservice: GetCart - No cart in session return empty")

		return cs.getEmptyCart()
	}

	return guestCart, nil
}

// ValidateCart validates a carts content
func (cs CartService) ValidateCart(ctx web.Context, decoratedCart cartDomain.DecoratedCart) cartDomain.CartValidationResult {
	if cs.CartValidator != nil {
		// TODO pass delivery Method from CART - once cart supports this!
		result := cs.CartValidator.Validate(ctx, decoratedCart, cs.DefaultDeliveryMethodForValidation)

		return result
	}

	return cartDomain.CartValidationResult{}
}

// GetDecoratedCart Get the correct Cart
func (cs *CartService) GetDecoratedCart(ctx web.Context) (cartDomain.DecoratedCart, error) {
	var empty cartDomain.DecoratedCart
	cs.Logger.Info("cart.application.cartservice: Get decorated cart ")

	cart, err := cs.GetCart(ctx)
	if err != nil {
		return empty, err
	}

	return *cs.CartDecoratorFactory.Create(ctx, cart), nil
}

// AddProduct Add a product
func (cs *CartService) AddProduct(ctx web.Context, addRequest cartDomain.AddRequest) error {
	addRequest, err := cs.checkProductAndEnrichAddRequest(ctx, addRequest)
	if err != nil {
		return err
	}

	if cs.isLoggedIn(ctx) {
		cart, err := cs.CustomerCartService.GetCart(ctx, cs.Auth(ctx), "me")
		if err != nil {
			return err
		}

		err = cs.CustomerCartService.AddToCart(ctx, cs.Auth(ctx), cart.ID, addRequest)

		if err == nil {
			// publish cart event if addtocart was successful
			cs.publishAddtoCartEvent(ctx, cart, addRequest)
		}

		return err
	}

	return cs.addProductToGuestCart(ctx, addRequest)
}

// checkProductAndEnrichAddRequest existence and validate with productService
func (cs *CartService) checkProductAndEnrichAddRequest(ctx web.Context, addRequest cartDomain.AddRequest) (cartDomain.AddRequest, error) {
	product, err := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
	if err != nil {
		return addRequest, errors.New("cart.application.cartservice - AddProduct:Product not found")
	}

	if product.Type() == productDomain.TYPECONFIGURABLE {
		if addRequest.VariantMarketplaceCode == "" {
			return addRequest, errors.New("cart.application.cartservice - AddProduct:No Variant given for configurable product")
		}

		configurableProduct := product.(productDomain.ConfigurableProduct)
		_, err := configurableProduct.Variant(addRequest.VariantMarketplaceCode)
		if err != nil {
			return addRequest, errors.New("cart.application.cartservice - AddProduct:Product has not the given variant")
		}
	}

	//Now Validate the Item with the optional registered ItemValidator
	if cs.ItemValidator != nil {
		return addRequest, cs.ItemValidator.Validate(ctx, addRequest, product)
	}

	return addRequest, nil
}

// addProductToGuestCart Handle Adding to Guest Cart
func (cs *CartService) addProductToGuestCart(ctx web.Context, addRequest cartDomain.AddRequest) error {
	//check if we have a guest cart in the session
	guestCart, err := cs.GetSessionGuestCart(ctx)

	if err != nil {
		// if not try to create a new one
		guestCart, err = cs.createNewSessionGuestCart(ctx)
		if err != nil {
			//no mitigation - return error
			return err
		}
	}

	guestCartID := ctx.Session().Values["cart.guestid"].(string)
	// Add to guest cart

	err = cs.GuestCartService.AddToCart(ctx, cs.Auth(ctx), guestCartID, addRequest)

	if err != nil {
		cs.Logger.Errorf("cart.application.cartservice: Failed Adding to cart %s Error %s", guestCartID, err)
		return err
	}

	if err != nil {
		return err
	}

	cs.publishAddtoCartEvent(ctx, guestCart, addRequest)

	cs.Logger.Infof("cart.application.cartservice: Added to cart %s", guestCartID)

	return nil
}

// isLoggedIn Checks if a user is logged in / authenticated
func (cs *CartService) isLoggedIn(ctx web.Context) bool {
	return cs.UserService.IsLoggedIn(ctx)
}

// HasSessionAGuestCart
func (cs *CartService) HasSessionAGuestCart(ctx web.Context) bool {
	if _, ok := ctx.Session().Values["cart.guestid"]; ok {
		return true
	}
	
	return false
}

// GetSessionGuestCart
func (cs *CartService) GetSessionGuestCart(ctx web.Context) (cartDomain.Cart, error) {
	var cart cartDomain.Cart

	if guestcartid, ok := ctx.Session().Values["cart.guestid"]; ok {
		existingCart, err := cs.GuestCartService.GetCart(ctx, cs.Auth(ctx), guestcartid.(string))
		if err != nil {
			cs.Logger.Errorf("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, err)
		}

		return existingCart, err
	}

	return cart, errors.New("No cart in session yet")
}

// GetSessionGuestCart
func (cs *CartService) DeleteSessionGuestCart(ctx web.Context) error {
	delete(ctx.Session().Values, "cart.guestid")
	//TODO - trigger backend also to be able to delete the cart there ( cs.GuestCartService.DeleteCart())

	return nil
}

// createNewSessionGuestCart Requests a new Guest Cart and stores the id in the session, if possible
func (cs *CartService) createNewSessionGuestCart(ctx web.Context) (cartDomain.Cart, error) {
	newGuestCart, err := cs.GuestCartService.GetNewCart(ctx, cs.Auth(ctx))

	if err != nil {
		cs.Logger.Errorf("cart.application.cartservice: Cannot create a new guest cart. Error %s", err)
		delete(ctx.Session().Values, "cart.guestid")

		return newGuestCart, err
	}

	cs.Logger.Infof("cart.application.cartservice: Requested new Guestcart %v", newGuestCart)
	ctx.Session().Values["cart.guestid"] = newGuestCart.ID

	return newGuestCart, nil
}

func (cs *CartService) getEmptyCart() (cartDomain.Cart, error) {
	return cartDomain.Cart{}, nil
}

func (cs *CartService) publishAddtoCartEvent(ctx web.Context, currentCart cartDomain.Cart, addRequest cartDomain.AddRequest) {
	currentCart.EventPublisher.PublishAddToCartEvent(ctx, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty)
}
