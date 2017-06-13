package infrastructure

import (
	"errors"
	"flamingo/core/cart/domain"
	"fmt"
	"math/rand"
)

// Fakecartrepository for testing/dev
type Fakecartrepository struct {
	Carts map[int]*domain.Cart
}

// FakecartrepositoryFactory factory
func FakecartrepositoryFactory() *Fakecartrepository {
	return &Fakecartrepository{
		Carts: make(map[int]*domain.Cart),
	}
}

func (cr *Fakecartrepository) init() {
	if cr.Carts == nil {
		cr.Carts = make(map[int]*domain.Cart)
	}
}

// Add to cart
func (cr *Fakecartrepository) Add(Cart domain.Cart) (int, error) {
	cr.init()
	fmt.Println("Fake cartrepo impl called add")
	if Cart.ID == 0 {
		Cart.ID = rand.Int()
	}
	cr.Carts[Cart.ID] = &Cart
	return Cart.ID, nil
}

// Update cart
func (cr *Fakecartrepository) Update(Cart domain.Cart) error {
	cr.init()
	fmt.Println("Fake cartrepo impl called Update")
	cr.Carts[Cart.ID] = &Cart
	return nil
}

// Delete cart
func (cr *Fakecartrepository) Delete(Cart domain.Cart) error {
	cr.init()
	fmt.Println("Fake cartrepo impl called delete")
	delete(cr.Carts, Cart.ID)
	return nil
}

// Get cart
func (cr *Fakecartrepository) Get(ID int) (*domain.Cart, error) {
	cr.init()
	fmt.Printf("Fake cartrepo impl called get for %s", ID)
	if val, ok := cr.Carts[ID]; ok {
		return val, nil
	}

	return nil, errors.New("No cart with that ID")
}
