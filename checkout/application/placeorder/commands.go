package placeorder

import (
	"net/url"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (

	//StartPlaceOrderCommand Command triggers new place order
	StartPlaceOrderCommand struct {
		Cart      cartDomain.Cart
		ReturnURL *url.URL
	}

	//RefreshPlaceOrderCommand Command
	RefreshPlaceOrderCommand struct {
	}

/*
	RefreshBlockingPlaceOrder struct {

	}

	CancelPlaceOrder struct {

	}
*/

)
