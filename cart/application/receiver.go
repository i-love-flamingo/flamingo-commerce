package application

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name Receiver --case snake --structname CartReceiver

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// Receiver that provides functionality regarding fetching the cart
	Receiver interface {
		ShouldHaveCart(ctx context.Context, session *web.Session) bool
		ShouldHaveGuestCart(session *web.Session) bool
		ViewCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error)
		ViewGuestCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, error)
		ModifyBehaviour(ctx context.Context) (cartDomain.ModifyBehaviour, error)
		GetCart(ctx context.Context, session *web.Session) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error)
		GetCartWithoutCache(ctx context.Context, session *web.Session) (*cartDomain.Cart, error)
	}
)
