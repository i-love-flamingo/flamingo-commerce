package domain

// CartService interface
type CartService interface {
	//GetEmptyCart - Called when no Entry in a (Guest) Cart is expected. Should just return an empty cart and can save unrequired API Calls in case
	GetEmptyCart() (Cart, error)
	//GetGuestCart - should return the guest Cart with the given id
	GetGuestCart(int) (Cart, error)
	//GetGuestCart - should return a new guest cart (including the id of the cart)
	GetNewGuestCart() (Cart, error)
	//AddToGuestCart - adds an item to a guest cart (cartid, productIdendifier, qty)
	AddToGuestCart(int, string, int) error
}
