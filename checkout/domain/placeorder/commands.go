package placeorder

import (
	context "context"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	placeOrder struct {
		Cart cartDomain.Cart
		Ctx  context.Context
	}

	StartPlaceOrder struct {
		placeOrder
	}

	RefreshPlaceOrder struct {
		placeOrder
	}

	RefreshBlockingPlaceOrder struct {
		placeOrder
	}

	CancelPlaceOrder struct {
		placeOrder
	}
)
