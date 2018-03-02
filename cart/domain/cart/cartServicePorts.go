package cart

import (
	"context"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

/**

CartServicePorts General Informations
 - GuestCartService: Used if no customer is logged in
 - CustomerCartService: Used if a customer is authenticated. You can access the users information in your Adapter implementation
 - When implementing the Ports in an own package, be sure to also set the correct "CartOrderBehaviour" on the cart!

*/

type (
	// GuestCartService interface
	GuestCartService interface {
		CommonCartService

		//GetGuestCart - should return a new guest cart (including the id of the cart)
		GetNewCart(ctx context.Context, auth Auth) (Cart, error)
	}

	// CustomerCartService  interface
	CustomerCartService interface {
		CommonCartService

		//MergeWithGuestCart
		//MergeWithGuestCart(ctx context.Context, auth Auth, cart Cart) error
	}

	CommonCartService interface {
		//GetGuestCart - should return the guest Cart with the given id
		GetCart(ctx context.Context, auth Auth, cartId string) (Cart, error)
		//AddToGuestCart - adds an item to a guest cart (cartid, marketplaceCode, qty)
		AddToCart(ctx context.Context, auth Auth, cartId string, addRequest AddRequest) error
	}

	// CartOrderBehaviour is a Port that can be implemented by other packages to implement  cart actions required for Ordering a Cart
	CartOrderBehaviour interface {
		PlaceOrder(ctx context.Context, auth Auth, cart *Cart, payment *Payment) (string, error)
		DeleteItem(ctx context.Context, auth Auth, cart *Cart, itemId string) error
		UpdateItem(ctx context.Context, auth Auth, cart *Cart, itemId string, item Item) error
		SetShippingInformation(ctx context.Context, auth Auth, cart *Cart, shippingAddress *Address, billingAddress *Address, shippingCarrierCode string, shippingMethodCode string) error
	}

	AddRequest struct {
		MarketplaceCode        string
		Qty                    int
		VariantMarketplaceCode string
		DeliveryIntent         string
	}

	// Auth defines cart authentication information
	Auth struct {
		TokenSource oauth2.TokenSource
		IDToken     *oidc.IDToken
	}
)
