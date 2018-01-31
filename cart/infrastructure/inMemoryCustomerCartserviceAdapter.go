package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	domaincart "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
)

// In Session Cart Storage
type (
	CustomerCartServiceAdapter struct {
		ProductService             productDomain.ProductService `inject:""`
		CustomerCartOrderBehaviour *CustomerCartOrderBehaviour  `inject:""`
		CartProvider               domaincart.CartProvider      `inject:""`
	}
	CustomerCartOrderBehaviour struct {
		CustomerCartStorage CustomerCartStorage `inject:""`
	}

	//CustomerCartStorage Interface - mya be implemnted by othe rpersitence types later as well
	CustomerCartStorage interface {
		GetCart(id string) (*domaincart.Cart, error)
		HasCart(id string) bool
		StoreCart(cart domaincart.Cart) error
	}

	//InmemoryCustomerCartStorage - for now the default implementation of CustomerCartStorage
	InmemoryCustomerCartStorage struct {
		guestCarts map[string]domaincart.Cart
	}
)

// Test Assignement and if Interface is implemeted correct
var (
	_ domaincart.CustomerCartService = &CustomerCartServiceAdapter{}
	_ domaincart.CartOrderBehaviour  = &CustomerCartOrderBehaviour{}
)

func (cs *CustomerCartServiceAdapter) GetCart(ctx context.Context, auth domaincart.Auth, guestcartid string) (domaincart.Cart, error) {
	cart := cs.CartProvider()
	if cs.CustomerCartOrderBehaviour.CustomerCartStorage == nil {
		return *cart, fmt.Errorf("cart.infrastructure.CustomerCartServiceAdapter: no CustomerCartStorage given")
	}

	cart.CartOrderBehaviour = domaincart.CartOrderBehaviour(cs.CustomerCartOrderBehaviour)

	if cs.CustomerCartOrderBehaviour.CustomerCartStorage.HasCart(guestcartid) {
		guestCart, e := cs.CustomerCartOrderBehaviour.CustomerCartStorage.GetCart(guestcartid)
		if e != nil {
			return *cart, fmt.Errorf("cart.infrastructure.CustomerCartServiceAdapter: Cart with ID %v could not be received from storage: %v", guestcartid, e)
		}

		var total big.Float
		total.SetFloat64(0)
		for _, item := range guestCart.Cartitems {
			total.Add(&total, big.NewFloat(item.RowTotal))
		}
		guestCart.ShippingItem.Title = "Shipping"
		guestCart.ShippingItem.Price = 9.99
		guestCart.SubTotal, _ = total.Float64()
		guestCart.GrandTotal, _ = new(big.Float).Add(&total, big.NewFloat(guestCart.ShippingItem.Price)).Float64()
		return *guestCart, nil
	}

	return *cart, fmt.Errorf("cart.infrastructure.CustomerCartServiceAdapter: Customer Cart with ID %v not exitend", guestcartid)
}

// AddToCart adds products to a cart
func (cs CustomerCartServiceAdapter) AddToCart(ctx context.Context, auth domaincart.Auth, guestcartid string, addRequest domaincart.AddRequest) error {
	if !cs.CustomerCartOrderBehaviour.CustomerCartStorage.HasCart(guestcartid) {
		return fmt.Errorf("cart.infrastructure.CustomerCartServiceAdapter: Cannot add - Customercart with id %v not existend", guestcartid)
	}

	guestcart, e := cs.CustomerCartOrderBehaviour.CustomerCartStorage.GetCart(guestcartid)
	if e != nil {
		return fmt.Errorf("cart.infrastructure.CustomerCartServiceAdapter: Cart with ID %v could not be received from storage: %v", guestcartid, e)
	}
	found, lineNr := guestcart.HasItem(addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode)
	if found {
		item, _ := guestcart.GetByLineNr(lineNr)
		item.Qty = item.Qty + addRequest.Qty
		calculateItemPrices(item)
	} else {
		product, _ := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
		cartItem := domaincart.Item{
			MarketplaceCode:        addRequest.MarketplaceCode,
			VariantMarketPlaceCode: addRequest.VariantMarketplaceCode,
			Qty:   addRequest.Qty,
			Price: product.SaleableData().ActivePrice.GetFinalPrice(),
			ID:    strconv.Itoa((rand.Int())),
		}
		guestcart.CurrencyCode = product.SaleableData().ActivePrice.Currency
		calculateItemPrices(&cartItem)
		guestcart.Cartitems = append(guestcart.Cartitems, cartItem)
	}

	cs.CustomerCartOrderBehaviour.CustomerCartStorage.StoreCart(*guestcart)

	return nil
}

// SetShippingInformation adds a product
func (cs CustomerCartOrderBehaviour) SetShippingInformation(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, shippingAddress *domaincart.Address, billingAddress *domaincart.Address, shippingCarrierCode string, shippingMethodCode string) error {
	return nil
}

// SetShippingInformation adds a product
func (cs CustomerCartOrderBehaviour) PlaceOrder(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, payment *domaincart.Payment) (string, error) {
	rand.Seed(time.Now().Unix())
	return strconv.Itoa(rand.Int()), nil
}

func (cs CustomerCartOrderBehaviour) DeleteItem(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, itemId string) error {
	if cs.CustomerCartStorage == nil {
		return fmt.Errorf("cart.infrastructure.CustomerCartOrderBehaviour: no CustomerCartStorage given")
	}
	if !cs.CustomerCartStorage.HasCart(cart.ID) {
		return fmt.Errorf("cart.infrastructure.CustomerCartOrderBehaviour: Cannot delete - Customercart with id %v not existend", cart.ID)
	}

	fmt.Printf("Inmemory Service Delete %v in %#v", itemId, cart.Cartitems)
	for k, item := range cart.Cartitems {
		if item.ID == itemId {
			if len(cart.Cartitems) > k {
				cart.Cartitems = append(cart.Cartitems[:k], cart.Cartitems[k+1:]...)
			} else {
				cart.Cartitems = cart.Cartitems[:k]
			}
		}
	}
	cs.CustomerCartStorage.StoreCart(*cart)
	return nil
}

func (cs CustomerCartOrderBehaviour) UpdateItem(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, itemId string, item domaincart.Item) error {
	if cs.CustomerCartStorage == nil {
		return fmt.Errorf("cart.infrastructure.CustomerCartOrderBehaviour: no CustomerCartStorage given")
	}
	if !cs.CustomerCartStorage.HasCart(cart.ID) {
		return fmt.Errorf("cart.infrastructure.CustomerCartOrderBehaviour: Cannot update - Customercart with id %v not existend", cart.ID)
	}

	if item.Qty < 1 {
		item.Qty = 1
	}
	calculateItemPrices(&item)
	for k, currentItem := range cart.Cartitems {
		if currentItem.ID == itemId {
			cart.Cartitems[k] = item
		}
	}

	cs.CustomerCartStorage.StoreCart(*cart)
	return nil
}

//********InmemoryCustomerCartStorage************

func (s *InmemoryCustomerCartStorage) init() {
	if s.guestCarts == nil {
		s.guestCarts = make(map[string]domaincart.Cart)
	}
}

func (s *InmemoryCustomerCartStorage) HasCart(id string) bool {
	s.init()
	if _, ok := s.guestCarts[id]; ok {
		return true
	}
	return false
}

func (s *InmemoryCustomerCartStorage) GetCart(id string) (*domaincart.Cart, error) {
	s.init()
	if cart, ok := s.guestCarts[id]; ok {
		return &cart, nil
	}
	return nil, errors.New("No cart stored")
}

func (s *InmemoryCustomerCartStorage) StoreCart(cart domaincart.Cart) error {
	s.init()
	s.guestCarts[cart.ID] = cart
	return nil
}
