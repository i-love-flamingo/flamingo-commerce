package cart

import "context"

type (
	// GuestCartService interface
	GuestCartService interface {
		CommonCartService

		//GetGuestCart - should return a new guest cart (including the id of the cart)
		GetNewCart(context.Context) (Cart, error)
	}

	// CustomerCartService  interface
	CustomerCartService interface {
		CommonCartService

		//MergeWithGuestCart
		MergeWithGuestCart(context.Context, Cart) error
	}

	CommonCartService interface {
		//GetGuestCart - should return the guest Cart with the given id
		GetCart(context.Context, string) (Cart, error)
		//AddToGuestCart - adds an item to a guest cart (cartid, marketplaceCode, qty)
		AddToCart(context.Context, string, AddRequest) error
	}

	AddRequest struct {
		MarketplaceCode        string
		Qty                    int
		VariantMarketplaceCode string
		//Identifier - some Adapters may need the Identifier instead of the MarketplaceCode, thats the reason why the AddRequest has it additionally to the MarketplaceCode attributes
		Identifier string
	}
)
