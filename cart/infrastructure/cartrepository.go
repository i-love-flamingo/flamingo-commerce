package infrastructure

import (
	"flamingo/core/cart/domain"
	"fmt"
	"math/rand"
	"errors"
)

type Fakecartrepository struct {
	Carts map[int]*domain.Cart
}

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


func (cr *Fakecartrepository) Add(Cart domain.Cart) (int, error) {
	cr.init()
	fmt.Println("Fake cartrepo impl called add")
	if Cart.Id == 0 {
		Cart.Id = rand.Int()
	}
	cr.Carts[Cart.Id]=&Cart
	return Cart.Id, nil
}


func (cr *Fakecartrepository) Update(Cart domain.Cart) error {
	cr.init()
	fmt.Println("Fake cartrepo impl called Update")
	cr.Carts[Cart.Id]=&Cart
	return nil
}

func (cr *Fakecartrepository) Delete(Cart domain.Cart) error {
	cr.init()
	fmt.Println("Fake cartrepo impl called delete")
	delete(cr.Carts,Cart.Id)
	return nil
}

func (cr *Fakecartrepository) Get(Id int) (*domain.Cart, error) {
	cr.init()
	fmt.Printf("Fake cartrepo impl called get for %s",Id)
	if val, ok := cr.Carts[Id]; ok {
		return val, nil
	}

	return nil, errors.New("No cart with that ID")
}


