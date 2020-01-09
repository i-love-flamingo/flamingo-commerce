package placeorder

import (
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (

	//StartPlaceOrder Command triggers new place order
	StartPlaceOrder struct {
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
