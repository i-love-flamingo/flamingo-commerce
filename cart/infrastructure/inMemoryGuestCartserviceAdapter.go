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
	GuestCartServiceAdapter struct {
		ProductService          productDomain.ProductService `inject:""`
		GuestCartOrderBehaviour GuestCartOrderBehaviour      `inject:""`
	}
	GuestCartOrderBehaviour struct {
		GuestCartStorage GuestCartStorage `inject:""`
	}

	//GuestCartStorage Interface - mya be implemnted by othe rpersitence types later as well
	GuestCartStorage interface {
		GetCart(id string) (*domaincart.Cart, error)
		HasCart(id string) bool
		StoreCart(cart domaincart.Cart) error
	}

	//InmemoryGuestCartStorage - for now the default implementation of GuestCartStorage
	InmemoryGuestCartStorage struct {
		guestCarts map[string]domaincart.Cart
	}
)

// Test Assignement and if Interface is implemeted correct
var (
	_ domaincart.GuestCartService   = &GuestCartServiceAdapter{}
	_ domaincart.CartOrderBehaviour = &GuestCartOrderBehaviour{}
)

func (cs *GuestCartServiceAdapter) GetCart(ctx context.Context, auth domaincart.Auth, guestcartid string) (domaincart.Cart, error) {
	var cart domaincart.Cart
	if cs.GuestCartOrderBehaviour.GuestCartStorage == nil {
		return cart, fmt.Errorf("cart.infrastructure.GuestCartServiceAdapter: no GuestCartStorage given")
	}

	cart.CartOrderBehaviour = domaincart.CartOrderBehaviour(cs.GuestCartOrderBehaviour)

	if cs.GuestCartOrderBehaviour.GuestCartStorage.HasCart(guestcartid) {
		guestCart, e := cs.GuestCartOrderBehaviour.GuestCartStorage.GetCart(guestcartid)
		if e != nil {
			return cart, fmt.Errorf("cart.infrastructure.GuestCartServiceAdapter: Cart with ID %v could not be received from storage: %v", guestcartid, e)
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

	return cart, fmt.Errorf("cart.infrastructure.GuestCartServiceAdapter: Guest Cart with ID %v not exitend", guestcartid)
}

// GetNewCart Creates a new guest cart and returns it
func (cs *GuestCartServiceAdapter) GetNewCart(ctx context.Context, auth domaincart.Auth) (domaincart.Cart, error) {
	guestCart := domaincart.Cart{
		ID: strconv.Itoa(rand.Int()),
	}
	guestCart.CartOrderBehaviour = domaincart.CartOrderBehaviour(cs.GuestCartOrderBehaviour)
	cs.GuestCartOrderBehaviour.GuestCartStorage.StoreCart(guestCart)
	return guestCart, nil
}

// AddToCart adds products to a cart
func (cs GuestCartServiceAdapter) AddToCart(ctx context.Context, auth domaincart.Auth, guestcartid string, addRequest domaincart.AddRequest) error {
	if !cs.GuestCartOrderBehaviour.GuestCartStorage.HasCart(guestcartid) {
		return fmt.Errorf("cart.infrastructure.GuestCartServiceAdapter: Cannot add - Guestcart with id %v not existend", guestcartid)
	}

	guestcart, e := cs.GuestCartOrderBehaviour.GuestCartStorage.GetCart(guestcartid)
	if e != nil {
		return fmt.Errorf("cart.infrastructure.GuestCartServiceAdapter: Cart with ID %v could not be received from storage: %v", guestcartid, e)
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

	cs.GuestCartOrderBehaviour.GuestCartStorage.StoreCart(*guestcart)

	return nil
}

// SetShippingInformation adds a product
func (cs GuestCartOrderBehaviour) SetShippingInformation(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, shippingAddress *domaincart.Address, billingAddress *domaincart.Address, shippingCarrierCode string, shippingMethodCode string) error {
	return nil
}

// SetShippingInformation adds a product
func (cs GuestCartOrderBehaviour) PlaceOrder(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, payment *domaincart.Payment) (string, error) {
	rand.Seed(time.Now().Unix())
	return strconv.Itoa(rand.Int()), nil
}

func (cs GuestCartOrderBehaviour) DeleteItem(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, itemId string) error {
	if cs.GuestCartStorage == nil {
		return fmt.Errorf("cart.infrastructure.GuestCartOrderBehaviour: no GuestCartStorage given")
	}
	if !cs.GuestCartStorage.HasCart(cart.ID) {
		return fmt.Errorf("cart.infrastructure.GuestCartOrderBehaviour: Cannot delete - Guestcart with id %v not existend", cart.ID)
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
	cs.GuestCartStorage.StoreCart(*cart)
	return nil
}

func (cs GuestCartOrderBehaviour) UpdateItem(ctx context.Context, auth domaincart.Auth, cart *domaincart.Cart, itemId string, item domaincart.Item) error {
	if cs.GuestCartStorage == nil {
		return fmt.Errorf("cart.infrastructure.GuestCartOrderBehaviour: no GuestCartStorage given")
	}
	if !cs.GuestCartStorage.HasCart(cart.ID) {
		return fmt.Errorf("cart.infrastructure.GuestCartOrderBehaviour: Cannot update - Guestcart with id %v not existend", cart.ID)
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

	cs.GuestCartStorage.StoreCart(*cart)
	return nil
}

func calculateItemPrices(item *domaincart.Item) {
	item.RowTotal, _ = new(big.Float).Mul(big.NewFloat(item.Price), new(big.Float).SetInt64(int64(item.Qty))).Float64()
}

//********InmemoryGuestCartStorage************

func (s *InmemoryGuestCartStorage) init() {
	if s.guestCarts == nil {
		s.guestCarts = make(map[string]domaincart.Cart)
	}
}

func (s *InmemoryGuestCartStorage) HasCart(id string) bool {
	s.init()
	if _, ok := s.guestCarts[id]; ok {
		return true
	}
	return false
}

func (s *InmemoryGuestCartStorage) GetCart(id string) (*domaincart.Cart, error) {
	s.init()
	if cart, ok := s.guestCarts[id]; ok {
		return &cart, nil
	}
	return nil, errors.New("No cart stored")
}

func (s *InmemoryGuestCartStorage) StoreCart(cart domaincart.Cart) error {
	s.init()
	s.guestCarts[cart.ID] = cart
	return nil
}
