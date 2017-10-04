package application

import (
	"log"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/web"
)

// CartService application struct
type (
	CartService struct {
		GuestCartService cart.GuestCartService `inject:""`
		//CustomerCartService  cart.CustomerCartService  `inject:""`
		CartDecoratorFactory cart.DecoratedCartFactory    `inject:""`
		ProductService       productDomain.ProductService `inject:""`
	}
)

// GetCart Get the correct Cart
func (cs *CartService) GetCart(ctx web.Context) (cart.Cart, error) {
	if cs.isLoggedIn(ctx) {
		return cs.getEmptyCart()
	} else {
		guestCart, e := cs.getSessionsGuestCart(ctx)
		if e != nil {
			log.Printf("cart.application.cartservice: GetCart - No cart in session return empty")
			return cs.getEmptyCart()
		}
		return guestCart, nil
	}
}

// GetDecoratedCart Get the correct Cart
func (cs *CartService) GetDecoratedCart(ctx web.Context) (cart.DecoratedCart, error) {
	var empty cart.DecoratedCart
	cart, e := cs.GetCart(ctx)
	log.Printf("cart.application.cartservice: Get decorated cart ")
	if e != nil {
		return empty, e
	}
	return *cs.CartDecoratorFactory.Create(ctx, cart), nil
}

// AddProduct Add a product
func (cs *CartService) AddProduct(ctx web.Context, addRequest cart.AddRequest) error {
	e := cs.checkProduct(ctx, addRequest)
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

// checkProduct existence and validate with productService
func (cs *CartService) checkProduct(ctx web.Context, addRequest cart.AddRequest) error {
	product, e := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
	if product.Type() == productDomain.TYPECONFIGURABLE {
		if addRequest.VariantMarketplaceCode == "" {
			return errors.New("cart.application.cartservice - AddProduct:No Variant given for configurable product")
		}
		configurableProduct := product.(productDomain.ConfigurableProduct)
		_, e := configurableProduct.Variant(addRequest.VariantMarketplaceCode)
		if e != nil {
			return errors.New("cart.application.cartservice - AddProduct:Product has not the given variant")
		}
	}
	if e != nil {
		return errors.New("cart.application.cartservice - AddProduct:Product not found")
	}
	return nil
}

// addProductToGuestCart Handle Adding to Guest Cart
func (cs *CartService) addProductToGuestCart(ctx web.Context, addRequest cart.AddRequest) error {
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
		log.Printf("cart.application.cartservice: Failed Adding to cart %s Error %s", guestCartID, e)
		return e
	}
	log.Printf("cart.application.cartservice: Added to cart %s", guestCartID)
	return nil
}

// isLoggedIn Checks if a user is logged in / authenticated
// @TODO
func (cs *CartService) isLoggedIn(ctx web.Context) bool {
	return false
}

// getSessionsGuestCart Checks if a valid guest cart exists for the session and tries to get it
// If no guest cart is registered or the existing one cannot be get it returns error that need to be handeled
func (cs *CartService) getSessionsGuestCart(ctx web.Context) (cart.Cart, error) {
	var cart cart.Cart
	if guestcartid, ok := ctx.Session().Values["cart.guestid"]; ok {
		existingCart, e := cs.GuestCartService.GetCart(ctx, guestcartid.(string))
		if e != nil {
			log.Printf("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, e)
		}
		return existingCart, e
	}
	return cart, errors.New("No cart in session yet")
}

// createNewSessionGuestCart Requests a new Guest Cart and stores the id in the session, if possible
func (cs *CartService) createNewSessionGuestCart(ctx web.Context) (cart.Cart, error) {
	newGuestCart, e := cs.GuestCartService.GetNewCart(ctx)
	if e != nil {
		log.Printf("cart.application.cartservice: Cannot create a new guest cart. Error %s", e)
		delete(ctx.Session().Values, "cart.guestid")
		return newGuestCart, e
	}
	log.Printf("cart.application.cartservice: Requested new Guestcart %v", newGuestCart)
	ctx.Session().Values["cart.guestid"] = newGuestCart.ID
	return newGuestCart, nil
}

func (cs *CartService) getEmptyCart() (cart.Cart, error) {
	emptyCart := cart.Cart{
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
