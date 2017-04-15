package application

import (
	"flamingo/core/cart/domain"
	"fmt"
	"github.com/gorilla/sessions"
)


type Cartservice struct {
	Cartrepository domain.Cartrepository `inject:""`
	Session *sessions.Session `inject:""`
}


// Gets the main basket for the current session
func (This *Cartservice) GetSessionCart() *domain.Cart {
	fmt.Println("Load cart for Session")

	// get cartID from session - if there is no cart return empty new cart
	cartId := 1

	// TODO impl session
	Cart, e := This.Cartrepository.Get(cartId)
	if e != nil {
		//TODO evaluate errorcode and decide resilience and behavoiur
		// for now we create new cart
		newCart := domain.Cart{
			cartId,
			nil,
		}
		cartId, _ = This.Cartrepository.Add(newCart)
		return &newCart
	}
	return Cart
}



// Adds a item by Code and qty to the sessions Cart
func (This *Cartservice) AddItem(code string, qty int)  {
	fmt.Println("Application Cartservice. Add "+code)

	//TODO - inject productService/Repo
	//Todo - deöegate to Cart (since this is the aggregate root and too much logic in the applicationservice)- and use code as unique id for item (Look at API details first)
	cartItem := domain.Cartitem{
		code,
		qty,
		12.99,
	}
	Cart := This.GetSessionCart()
	Cart.Add(cartItem)
	This.Cartrepository.Update(*Cart)
}




//todo - this is test - would be a domain event not ccontroller
func (cs *Cartservice) OnLogin(event domain.LoginSucessEvent) {
	fmt.Printf("LoginSucess going to merge carts now %s",event)

}
