package infrastructure

import (
	"flamingo/core/cart/domain"
	"fmt"
)

type Cartrepository struct {
}

func (cr *Cartrepository) Add(domain.Cart) (int, error) {
	fmt.Println("Real impl called add")
	return 999, nil
}


func (cr *Cartrepository) Update(domain.Cart) error {
	fmt.Println("Real impl called Update")
	return nil
}

func (cr *Cartrepository) Delete(domain.Cart) error {
	fmt.Println("Real impl called delete")
	return nil
}

func (cr *Cartrepository) Get(Id int) (*domain.Cart, error) {
	fmt.Printf("Real impl called get for %s",Id)
	return &domain.Cart{
		Id,
		nil,
	}, nil
}


