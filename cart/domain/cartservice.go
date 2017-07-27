package domain

// CartService interface
type CartService interface {
	//Called when no Entry in a (Guest) Cart is expected. Should just return an empty cart and can save unrequired API Calls in case
	GetEmptyCart() (Cart, error)
	GetGuestCart(int) (Cart, error)
	GetNewGuestCart() (Cart, error)
	AddToGuestCart(int, string, int) error
}
