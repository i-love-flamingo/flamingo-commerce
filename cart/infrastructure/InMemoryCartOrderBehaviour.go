package infrastructure

import (
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"context"
	"errors"
)

type (
	InMemoryCartOrderBehaviour struct {
		CartStorage CartStorage `inject:""`
	}

	//GuestCartStorage Interface - mya be implemnted by othe rpersitence types later as well
	CartStorage interface {
		GetCart(id string) (*cart.Cart, error)
		HasCart(id string) bool
		StoreCart(cart cart.Cart) error
	}

	// InMemoryCartStorage - for now the default implementation of GuestCartStorage
	InMemoryCartStorage struct {
		guestCarts map[string]cart.Cart
	}
)

var (
	_ cart.CartBehaviour = (*InMemoryCartOrderBehaviour)(nil)
)

func (cob *InMemoryCartOrderBehaviour) PlaceOrder(ctx context.Context, cart *cart.Cart, payment *cart.CartPayment) (string, error)                                                                                 {
	return "",  nil
}
func (cob *InMemoryCartOrderBehaviour) DeleteItem(ctx context.Context, cart *cart.Cart, itemId string) (*cart.Cart, error)                                                                                         {
	return nil,  nil
}
func (cob *InMemoryCartOrderBehaviour) UpdateItem(ctx context.Context, cart *cart.Cart, itemId string, itemUpdateCommand cart.ItemUpdateCommand) (*cart.Cart, error)                                               {
	return nil,  nil
}

/**
add item to cart and store in memory
 */
func (cob *InMemoryCartOrderBehaviour) AddToCart(ctx context.Context, cart *cart.Cart, addRequest cart.AddRequest) (*cart.Cart, error) {

	if !cob.CartStorage.HasCart(cart.ID) {

	}

	return nil,  nil
}

func (cob *InMemoryCartOrderBehaviour) UpdatePurchaser(ctx context.Context, cart *cart.Cart, purchaser *cart.Person, additionalData map[string]string) (*cart.Cart, error)                                         {
	return nil,  nil
}

func (cob *InMemoryCartOrderBehaviour) UpdateAdditionalData(ctx context.Context, cart *cart.Cart, additionalData map[string]string) (*cart.Cart, error)                                                            {
	return nil,  nil
}

func (cob *InMemoryCartOrderBehaviour) UpdateDeliveryInfosAndBilling(ctx context.Context, cart *cart.Cart, billingAddress *cart.Address, deliveryInfoUpdates []cart.DeliveryInfoUpdateCommand) (*cart.Cart, error) {
	return nil,  nil
}

/*
return the current cart from storage
 */
func (cob *InMemoryCartOrderBehaviour) GetCart(ctx context.Context, cartId string) (*cart.Cart, error) {
	return nil, nil
}




/** Implementation fo the storage **/
func (s *InMemoryCartStorage) init() {
	if s.guestCarts == nil {
		s.guestCarts = make(map[string]cart.Cart)
	}
}

func (s *InMemoryCartStorage) HasCart(id string) bool {
	s.init()
	if _, ok := s.guestCarts[id]; ok {
		return true
	}
	return false
}

func (s *InMemoryCartStorage) GetCart(id string) (*cart.Cart, error) {
	s.init()
	if cart, ok := s.guestCarts[id]; ok {
		return &cart, nil
	}
	return nil, errors.New("no cart stored")
}

func (s *InMemoryCartStorage) StoreCart(cart cart.Cart) error {
	s.init()
	s.guestCarts[cart.ID] = cart
	return nil
}
