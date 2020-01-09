package placeorder

import (
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (

	//StartPlaceOrderCommand Command triggers new place order
	StartPlaceOrderCommand struct {
		Cart cartDomain.Cart
	}

	//RefreshPlaceOrder Command
	RefreshPlaceOrder struct {
	}

/*
	RefreshBlockingPlaceOrder struct {

	}

	CancelPlaceOrder struct {

	}
*/

)
