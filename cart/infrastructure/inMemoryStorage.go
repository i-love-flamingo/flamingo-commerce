package infrastructure

import (
	"errors"
	"sync"

	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// InMemoryCartStorage - for now the default implementation of GuestCartStorage
	InMemoryCartStorage struct {
		guestCarts map[string]*domaincart.Cart
		locker     sync.Locker
	}
)

var _ CartStorage = &InMemoryCartStorage{}

func (s *InMemoryCartStorage) init() {
	if s.guestCarts == nil {
		s.guestCarts = make(map[string]*domaincart.Cart)
		s.locker = &sync.Mutex{}
	}
}

// HasCart checks if the cart storage has a cart with a given id
func (s *InMemoryCartStorage) HasCart(id string) bool {
	s.init()
	s.locker.Lock()
	defer s.locker.Unlock()
	if _, ok := s.guestCarts[id]; ok {
		return true
	}
	return false
}

// GetCart returns a cart with the given id from the cart storage
func (s *InMemoryCartStorage) GetCart(id string) (*domaincart.Cart, error) {
	s.init()
	s.locker.Lock()
	defer s.locker.Unlock()
	if cart, ok := s.guestCarts[id]; ok {
		return cart, nil
	}
	return nil, errors.New("no cart stored")
}

// StoreCart stores a cart in the storage
func (s *InMemoryCartStorage) StoreCart(cart *domaincart.Cart) error {
	s.init()
	s.locker.Lock()
	defer s.locker.Unlock()
	s.guestCarts[cart.ID] = cart
	return nil
}

// RemoveCart from storage
func (s *InMemoryCartStorage) RemoveCart(cart *domaincart.Cart) error {
	s.init()
	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.guestCarts, cart.ID)
	return nil
}
