package cart

import "context"

type (
	// GuestCartService interface
	GuestCartService interface {
		//GetGuestCart - should return the guest Cart with the given id
		GetCart(context.Context, string) (Cart, error)
		//GetGuestCart - should return a new guest cart (including the id of the cart)
		GetNewCart(context.Context) (Cart, error)
		//AddToGuestCart - adds an item to a guest cart (cartid, marketplaceCode, qty)
		AddToCart(context.Context, string, AddRequest) error
	}

	// CustomerCartService  interface
	CustomerCartService interface {
		GetCart(context.Context) ([]Cart, error)
		GetNewCart(context.Context) (Cart, error)
		AddToCart(context.Context, string, AddRequest) error
	}

	AddRequest struct {
		MarketplaceCode        string
		Qty                    int
		VariantMarketplaceCode string
	}
)
