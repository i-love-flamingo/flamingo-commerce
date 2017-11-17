package cart

import "context"

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
