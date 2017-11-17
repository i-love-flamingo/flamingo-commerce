package application

import (
	"github.com/pkg/errors"
	domaincart "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

// CartService application struct
type (
	CartService struct {
		GuestCartService domaincart.GuestCartService `inject:""`
		//CustomerCartService  cart.CustomerCartService  `inject:""`
		CartDecoratorFactory domaincart.DecoratedCartFactory `inject:""`
		ProductService       productDomain.ProductService    `inject:""`
		Logger               flamingo.Logger                 `inject:""`
	}
)

// GetCart Get the correct Cart
func (cs *CartService) GetCart(ctx web.Context) (domaincart.Cart, error) {
	if cs.isLoggedIn(ctx) {
		return cs.getEmptyCart()
	} else {
		guestCart, e := cs.getSessionsGuestCart(ctx)
		if e != nil {
			cs.Logger.Warn("cart.application.cartservice: GetCart - No cart in session return empty")
			return cs.getEmptyCart()
		}
		return guestCart, nil
	}
}

// GetDecoratedCart Get the correct Cart
func (cs *CartService) GetDecoratedCart(ctx web.Context) (domaincart.DecoratedCart, error) {
	var empty domaincart.DecoratedCart
	cart, e := cs.GetCart(ctx)
	cs.Logger.Info("cart.application.cartservice: Get decorated cart ")
	if e != nil {
		return empty, e
	}
	return *cs.CartDecoratorFactory.Create(ctx, cart), nil
}

// AddProduct Add a product
func (cs *CartService) AddProduct(ctx web.Context, addRequest domaincart.AddRequest) error {
	addRequest, e := cs.checkProductAndEnrichAddRequest(ctx, addRequest)
	if e != nil {
		return e
	}
	if cs.isLoggedIn(ctx) {
		//TODO
		return nil
	} else {
		return cs.addProductToGuestCart(ctx, addRequest)
	}
}

// checkProductAndEnrichAddRequest existence and validate with productService
func (cs *CartService) checkProductAndEnrichAddRequest(ctx web.Context, addRequest domaincart.AddRequest) (domaincart.AddRequest, error) {
	product, e := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
	if e != nil {
		return addRequest, errors.New("cart.application.cartservice - AddProduct:Product not found")
	}
	if product.Type() == productDomain.TYPECONFIGURABLE {
		if addRequest.VariantMarketplaceCode == "" {
			return addRequest, errors.New("cart.application.cartservice - AddProduct:No Variant given for configurable product")
		}
		configurableProduct := product.(productDomain.ConfigurableProduct)
		_, e := configurableProduct.Variant(addRequest.VariantMarketplaceCode)
		if e != nil {
			return addRequest, errors.New("cart.application.cartservice - AddProduct:Product has not the given variant")
		}
		configurable, _ := product.(productDomain.ConfigurableProduct)
		addRequest.Identifier = configurable.Identifier
	}
	simple, _ := product.(productDomain.SimpleProduct)
	addRequest.Identifier = simple.Identifier
	return addRequest, nil
}

// addProductToGuestCart Handle Adding to Guest Cart
func (cs *CartService) addProductToGuestCart(ctx web.Context, addRequest domaincart.AddRequest) error {
	//check if we have a guest cart in the session
	if _, e := cs.getSessionsGuestCart(ctx); e != nil {
		// if not try to create a new one
		_, e := cs.createNewSessionGuestCart(ctx)
		if e != nil {
			//no mitigation - return error
			return e
		}
	}
	guestCartID := ctx.Session().Values["cart.guestid"].(string)
	// Add to guest cart
	e := cs.GuestCartService.AddToCart(ctx, guestCartID, addRequest)
	if e != nil {
		cs.Logger.Errorf("cart.application.cartservice: Failed Adding to cart %s Error %s", guestCartID, e)
		return e
	}
	cs.Logger.Infof("cart.application.cartservice: Added to cart %s", guestCartID)
	return nil
}

// isLoggedIn Checks if a user is logged in / authenticated
// @TODO
func (cs *CartService) isLoggedIn(ctx web.Context) bool {
	return false
}

// getSessionsGuestCart Checks if a valid guest cart exists for the session and tries to get it
// If no guest cart is registered or the existing one cannot be get it returns error that need to be handeled
func (cs *CartService) getSessionsGuestCart(ctx web.Context) (domaincart.Cart, error) {
	var cart domaincart.Cart
	if guestcartid, ok := ctx.Session().Values["cart.guestid"]; ok {
		existingCart, e := cs.GuestCartService.GetCart(ctx, guestcartid.(string))
		if e != nil {
			cs.Logger.Errorf("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, e)
		}
		return existingCart, e
	}
	return cart, errors.New("No cart in session yet")
}

// createNewSessionGuestCart Requests a new Guest Cart and stores the id in the session, if possible
func (cs *CartService) createNewSessionGuestCart(ctx web.Context) (domaincart.Cart, error) {
	newGuestCart, e := cs.GuestCartService.GetNewCart(ctx)
	if e != nil {
		cs.Logger.Errorf("cart.application.cartservice: Cannot create a new guest cart. Error %s", e)
		delete(ctx.Session().Values, "cart.guestid")
		return newGuestCart, e
	}
	cs.Logger.Infof("cart.application.cartservice: Requested new Guestcart %v", newGuestCart)
	ctx.Session().Values["cart.guestid"] = newGuestCart.ID
	return newGuestCart, nil
}

func (cs *CartService) getEmptyCart() (domaincart.Cart, error) {
	emptyCart := domaincart.Cart{
		Cartitems: nil,
	}
	return emptyCart, nil
}

/*
// OnLogin todo - this is test - would be a domain event not ccontroller
func (cs *CartService) OnLogin(event domain.LoginSucessEvent) {
	fmt.Printf("LoginSucess going to merge carts now %s", event)

}
*/
