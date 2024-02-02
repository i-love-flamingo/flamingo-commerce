package infrastructure

import (
	"context"
	"errors"
	"sync"

	"go.opencensus.io/trace"

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

// Inject dependencies and prepare storage
// Important: InMemoryStorage MUST be bound AsEagerSingleton, Inject MUST be called in tests to behave as expected
func (s *InMemoryCartStorage) Inject() *InMemoryCartStorage {
	s.locker = &sync.Mutex{}
	s.guestCarts = make(map[string]*domaincart.Cart)

	return s
}

// HasCart checks if the cart storage has a cart with a given id
func (s *InMemoryCartStorage) HasCart(ctx context.Context, id string) bool {
	_, span := trace.StartSpan(ctx, "cart/InMemoryCartStorage/HasCart")
	defer span.End()

	s.locker.Lock()
	defer s.locker.Unlock()

	if _, ok := s.guestCarts[id]; ok {
		return true
	}
	return false
}

// GetCart returns a cart with the given id from the cart storage
func (s *InMemoryCartStorage) GetCart(ctx context.Context, id string) (*domaincart.Cart, error) {
	_, span := trace.StartSpan(ctx, "cart/InMemoryCartStorage/GetCart")
	defer span.End()

	s.locker.Lock()
	defer s.locker.Unlock()

	if cart, ok := s.guestCarts[id]; ok {
		return cart, nil
	}
	return nil, errors.New("no cart stored")
}

// StoreCart stores a cart in the storage
func (s *InMemoryCartStorage) StoreCart(ctx context.Context, cart *domaincart.Cart) error {
	_, span := trace.StartSpan(ctx, "cart/InMemoryCartStorage/StoreCart")
	defer span.End()

	s.locker.Lock()
	defer s.locker.Unlock()

	s.guestCarts[cart.ID] = cart
	return nil
}

// RemoveCart from storage
func (s *InMemoryCartStorage) RemoveCart(ctx context.Context, cart *domaincart.Cart) error {
	_, span := trace.StartSpan(ctx, "cart/InMemoryCartStorage/RemoveCart")
	defer span.End()

	s.locker.Lock()
	defer s.locker.Unlock()

	delete(s.guestCarts, cart.ID)
	return nil
}
