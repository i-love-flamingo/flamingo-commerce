package domain


type Cartrepository interface {
	Add(Cart) (int, error)
	Update(Cart) error
	Delete(Cart) error
	Get(int) (*Cart, error)
}
