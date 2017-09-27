package application

import (
	"flamingo/core/cart/domain/cart"
	product "flamingo/core/product/domain"

	"flamingo/framework/web"

	"log"

	"github.com/pkg/errors"
)

// CartService application struct
type CartService struct {
	GuestCartService cart.GuestCartService `inject:""`
	//CustomerCartService  cart.CustomerCartService  `inject:""`
	CartDecoratorFactory cart.DecoratedCartFactory `inject:""`
	ProductService       product.ProductService    `inject:""`
}

// GetCart Get the correct Cart
func (cs *CartService) GetCart(ctx web.Context) (cart.Cart, error) {
	if cs.isLoggedIn(ctx) {
		return cs.getEmptyCart()
	} else {
		guestCart, e := cs.getSessionsGuestCart(ctx)
		if e != nil {
			return cs.getEmptyCart()
		}
		return guestCart, nil
	}
}

// GetDecoratedCart Get the correct Cart
func (cs *CartService) GetDecoratedCart(ctx web.Context) (cart.DecoratedCart, error) {
	var empty cart.DecoratedCart
	cart, e := cs.GetCart(ctx)
	if e != nil {
		return empty, e
	}
	return *cs.CartDecoratorFactory.Create(ctx, cart), nil
}

// AddProduct Add a product
func (cs *CartService) AddProduct(ctx web.Context, marketplaceCode string, amount int) error {
	_, e := cs.ProductService.Get(ctx, marketplaceCode)
	if e != nil {
		return errors.New("cart.application.cartservice - Product not found")
	}
	if cs.isLoggedIn(ctx) {
		//TODO
		return nil
	} else {
		return cs.addProductToGuestCart(ctx, marketplaceCode, amount)
	}
}

// addProductToGuestCart Handle Adding to Guest Cart
func (cs *CartService) addProductToGuestCart(ctx web.Context, productCode string, amount int) error {
	//check if we have a guest cart in the session
	if _, e := cs.getSessionsGuestCart(ctx); e != nil {
		// if not try to create a new one
		_, e := cs.createNewSessionGuestCart(ctx)
		if e != nil {
			//no mitigation - return error
			return e
		}
	}
	guestCartID := ctx.Session().Values["cart.guestid"].(int)
	// Add to guest cart
	e := cs.GuestCartService.AddToCart(guestCartID, productCode, amount)
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
		existingCart, e := cs.GuestCartService.GetCart(guestcartid.(int))
		if e != nil {
			log.Printf("cart.application.cartservice: Guestcart id in session cannot be retrieved. Id %s, Error: %s", guestcartid, e)
		}
		return existingCart, e
	}
	return cart, errors.New("No cart in session yet")
}

// createNewSessionGuestCart Requests a new Guest Cart and stores the id in the session, if possible
func (cs *CartService) createNewSessionGuestCart(ctx web.Context) (cart.Cart, error) {
	newGuestCart, e := cs.GuestCartService.GetNewCart()
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
