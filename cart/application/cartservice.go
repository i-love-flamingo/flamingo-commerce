package application

import (
	"flamingo/core/cart/domain"
	"fmt"

	"github.com/gorilla/sessions"
)

// Cartservice application
type Cartservice struct {
	DomainCartService domain.CartService `inject:""`
	Session           *sessions.Session  `inject:""`
}

// GetSessionCart gets the main basket for the current session
func (cs *Cartservice) GetSessionCart() *domain.Cart {
	fmt.Println("Load cart for Session")

	// get cartID from session - if there is no cart return empty new cart
	cartID := 1

	// TODO impl session
	Cart, e := cs.DomainCartService.Get(cartID)
	if e != nil {
		//TODO evaluate errorcode and decide resilience and behavoiur
		// for now we create new cart
		newCart := domain.Cart{
			cartID,
			nil,
		}
		cartID, _ = cs.DomainCartService.Add(newCart)
		return &newCart
	}
	return Cart
}

// AddItem adds a item by Code and qty to the sessions Cart
func (cs *Cartservice) AddItem(code string, qty int) {
	fmt.Println("Application Cartservice. Add " + code)

	//TODO - inject productService/Repo
	//Todo - deöegate to Cart (since this is the aggregate root and too much logic in the applicationservice)- and use code as unique id for item (Look at API details first)
	cartItem := domain.Cartitem{
		code,
		qty,
		12.99,
	}
	Cart := cs.GetSessionCart()
	Cart.Add(cartItem)
	cs.DomainCartService.Update(*Cart)
}

// OnLogin todo - this is test - would be a domain event not ccontroller
func (cs *Cartservice) OnLogin(event domain.LoginSucessEvent) {
	fmt.Printf("LoginSucess going to merge carts now %s", event)

}
