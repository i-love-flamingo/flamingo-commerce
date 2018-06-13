package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	domaincart "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/product/domain"
)

type (
	InMemoryCartOrderBehaviour struct {
		CartStorage    CartStorage           `inject:""`
		ProductService domain.ProductService `inject:""`
	}

	//GuestCartStorage Interface - mya be implemnted by othe rpersitence types later as well
	CartStorage interface {
		GetCart(id string) (*domaincart.Cart, error)
		HasCart(id string) bool
		StoreCart(cart domaincart.Cart) error
	}

	// InMemoryCartStorage - for now the default implementation of GuestCartStorage
	InMemoryCartStorage struct {
		guestCarts map[string]domaincart.Cart
	}
)

var (
	_ domaincart.CartBehaviour = (*InMemoryCartOrderBehaviour)(nil)
)

// @todo implement when needed
func (cob *InMemoryCartOrderBehaviour) PlaceOrder(ctx context.Context, cart *domaincart.Cart, payment *domaincart.CartPayment) (string, error) {
	return "", nil
}

func (cob *InMemoryCartOrderBehaviour) DeleteItem(ctx context.Context, cart *domaincart.Cart, itemId string) (*domaincart.Cart, error) {
	if !cob.CartStorage.HasCart(cart.ID) {
		return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot delete - Guestcart with id %v not existend", cart.ID)
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

	cob.CartStorage.StoreCart(*cart)
	return cart, nil
}

func (cob *InMemoryCartOrderBehaviour) UpdateItem(ctx context.Context, cart *domaincart.Cart, itemId string, itemUpdateCommand domaincart.ItemUpdateCommand) (*domaincart.Cart, error) {
	if !cob.CartStorage.HasCart(cart.ID) {
		return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot add - Guestcart with id %v not existend", cart.ID)
	}

	fmt.Printf("Inmemory Service Update %v in %#v", itemId, cart.Cartitems)

	item, err := cart.GetByItemId(itemId)
	if err != nil {
		return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot delete - Guestcart item with id %v not existend", itemId)
	}

	item.Qty = *itemUpdateCommand.Qty

	calculateItemPrices(item)
	for k, currentItem := range cart.Cartitems {
		if currentItem.ID == itemId {
			cart.Cartitems[k] = *item
		}
	}

	cob.CartStorage.StoreCart(*cart)

	return cart, nil
}

/**
add item to cart and store in memory
*/
func (cob *InMemoryCartOrderBehaviour) AddToCart(ctx context.Context, cart *domaincart.Cart, addRequest domaincart.AddRequest) (*domaincart.Cart, error) {

	if !cob.CartStorage.HasCart(cart.ID) {
		return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot add - Guestcart with id %v not existend", cart.ID)
	}

	// check if the current item is in the cart and add or increase qty
	// @todo for the future we have to check the deliveryIntent too, it will be possible to add the same item twice
	// with different delivery intents
	found, lineNr := cart.HasItem(addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode)
	if found {
		item, _ := cart.GetByLineNr(lineNr)
		item.Qty = item.Qty + addRequest.Qty
		calculateItemPrices(item)
	} else {
		product, _ := cob.ProductService.Get(ctx, addRequest.MarketplaceCode)
		cartItem := domaincart.Item{
			MarketplaceCode:        addRequest.MarketplaceCode,
			VariantMarketPlaceCode: addRequest.VariantMarketplaceCode,
			Qty:          addRequest.Qty,
			SinglePrice:  product.SaleableData().ActivePrice.GetFinalPrice(),
			ID:           strconv.Itoa(rand.Int()),
			CurrencyCode: product.SaleableData().ActivePrice.Currency,
			OriginalDeliveryIntent: &domaincart.DeliveryIntent{
				Method: addRequest.DeliveryIntent.Method,
				AutodetectDeliveryLocation: false,
				DeliveryLocationType:       addRequest.DeliveryIntent.DeliveryLocationType,
				DeliveryLocationCode:       addRequest.DeliveryIntent.DeliveryLocationCode,
			},
		}
		calculateItemPrices(&cartItem)
		cart.Cartitems = append(cart.Cartitems, cartItem)
	}

	cob.CartStorage.StoreCart(*cart)

	return cart, nil
}

func calculateItemPrices(item *domaincart.Item) {
	item.RowTotal, _ = new(big.Float).Mul(big.NewFloat(item.SinglePrice), new(big.Float).SetInt64(int64(item.Qty))).Float64()
}

// @todo implement when needed
func (cob *InMemoryCartOrderBehaviour) UpdatePurchaser(ctx context.Context, cart *domaincart.Cart, purchaser *domaincart.Person, additionalData map[string]string) (*domaincart.Cart, error) {
	return nil, nil
}

// @todo implement when needed
func (cob *InMemoryCartOrderBehaviour) UpdateAdditionalData(ctx context.Context, cart *domaincart.Cart, additionalData map[string]string) (*domaincart.Cart, error) {
	return nil, nil
}

// @todo implement when needed
func (cob *InMemoryCartOrderBehaviour) UpdateDeliveryInfosAndBilling(ctx context.Context, cart *domaincart.Cart, billingAddress *domaincart.Address, deliveryInfoUpdates []domaincart.DeliveryInfoUpdateCommand) (*domaincart.Cart, error) {
	return nil, nil
}

/*
return the current cart from storage
*/
func (cob *InMemoryCartOrderBehaviour) GetCart(ctx context.Context, cartId string) (*domaincart.Cart, error) {
	if cob.CartStorage.HasCart(cartId) {
		// if cart exists, there is no error ;)
		cart, _ := cob.CartStorage.GetCart(cartId)
		return cart, nil
	}
	return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot get - Guestcart with id %v not existend", cartId)
}

/** Implementation fo the storage **/
func (s *InMemoryCartStorage) init() {
	if s.guestCarts == nil {
		s.guestCarts = make(map[string]domaincart.Cart)
	}
}

func (s *InMemoryCartStorage) HasCart(id string) bool {
	s.init()
	if _, ok := s.guestCarts[id]; ok {
		return true
	}
	return false
}

func (s *InMemoryCartStorage) GetCart(id string) (*domaincart.Cart, error) {
	s.init()
	if cart, ok := s.guestCarts[id]; ok {
		return &cart, nil
	}
	return nil, errors.New("no cart stored")
}

func (s *InMemoryCartStorage) StoreCart(cart domaincart.Cart) error {
	s.init()
	s.guestCarts[cart.ID] = cart
	return nil
}

func (cob *InMemoryCartOrderBehaviour) ApplyVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, error) {
	if couponCode != "valid" || couponCode == "" {
		err := errors.New("Code invalid")
		return nil, err
	}

	coupon := domaincart.CouponCode{
		Code: couponCode,
	}
	cart.AppliedCouponCodes = append(cart.AppliedCouponCodes, coupon)
	err := cob.CartStorage.StoreCart(*cart)

	return cart, err
}
