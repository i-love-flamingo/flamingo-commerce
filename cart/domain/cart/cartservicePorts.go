package cart

type (
	// GuestCartService interface
	GuestCartService interface {
		//GetGuestCart - should return the guest Cart with the given id
		GetCart(int) (Cart, error)
		//GetGuestCart - should return a new guest cart (including the id of the cart)
		GetNewCart() (Cart, error)
		//AddToGuestCart - adds an item to a guest cart (cartid, marketplaceCode, qty)
		AddToCart(int, string, int) error
	}
	// CustomerCartService  interface
	CustomerCartService interface {
		GetAllCarts() ([]Cart, error)
		GetCart(int) (Cart, error)
		GetNewCart() (Cart, error)
		AddToCart(int, string, int) error
	}
)
