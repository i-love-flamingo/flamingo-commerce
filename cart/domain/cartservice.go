package domain

type CartService interface {
	Add(Cart) (int, error)
	Update(Cart) error
	Delete(Cart) error
	Get(int) (*Cart, error)
}
