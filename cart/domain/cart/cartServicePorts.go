package cart

import (
	"context"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type (
	// GuestCartService interface
	GuestCartService interface {
		GetCart(ctx context.Context, cartId string) (*Cart, error)

		//GetGuestCart - should return a new guest cart (including the id of the cart)
		GetNewCart(ctx context.Context) (*Cart, error)

		GetCartOrderBehaviour(context.Context) (CartBehaviour, error)
	}

	// CustomerCartService  interface
	CustomerCartService interface {
		GetCartOrderBehaviour(context.Context, Auth) (CartBehaviour, error)
		GetCart(ctx context.Context, auth Auth, cartId string) (*Cart, error)
	}

	// CartBehaviour is a Port that can be implemented by other packages to implement  cart actions required for Ordering a Cart
	CartBehaviour interface {
		PlaceOrder(ctx context.Context, cart *Cart, payment *PaymentInfo) (string, error)
		DeleteItem(ctx context.Context, cart *Cart, itemId string) error
		UpdateItem(ctx context.Context, cart *Cart, itemId string, item Item) error
		AddToCart(ctx context.Context, cart *Cart, addRequest AddRequest) error
		SetShippingInformation(ctx context.Context, cart *Cart, shippingAddress *Address, billingAddress *Address, shippingCarrierCode string, shippingMethodCode string) error
	}

	AddRequest struct {
		MarketplaceCode        string
		Qty                    int
		VariantMarketplaceCode string
		DeliveryIntent         DeliveryIntent
	}

	// Auth defines cart authentication information
	Auth struct {
		TokenSource oauth2.TokenSource
		IDToken     *oidc.IDToken
	}
)
