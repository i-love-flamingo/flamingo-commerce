package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	domaincart "flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/product/domain"
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
func (cob *InMemoryCartOrderBehaviour) PlaceOrder(ctx context.Context, cart *domaincart.Cart, payment *domaincart.CartPayment) (domaincart.PlacedOrderInfos, error) {
	return nil, nil
}

func (cob *InMemoryCartOrderBehaviour) DeleteItem(ctx context.Context, cart *domaincart.Cart, itemId string, deliveryCode string) (*domaincart.Cart, error) {
	if !cob.CartStorage.HasCart(cart.ID) {
		return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot delete - Guestcart with id %v not existend", cart.ID)
	}

	if delivery, ok := cart.GetDeliveryByCode(deliveryCode); ok {
		fmt.Printf("Inmemory Service Delete %v in %#v", itemId, delivery.Cartitems)
		for k, item := range delivery.Cartitems {
			if item.ID == itemId {
				if len(delivery.Cartitems) > k {
					delivery.Cartitems = append(delivery.Cartitems[:k], delivery.Cartitems[k+1:]...)
				} else {
					delivery.Cartitems = delivery.Cartitems[:k]
				}
			}
		}
	}

	cob.CartStorage.StoreCart(*cart)
	return cart, nil
}

func (cob *InMemoryCartOrderBehaviour) UpdateItem(ctx context.Context, cart *domaincart.Cart, itemId string, deliveryCode string, itemUpdateCommand domaincart.ItemUpdateCommand) (*domaincart.Cart, error) {
	if !cob.CartStorage.HasCart(cart.ID) {
		return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot add - Guestcart with id %v not existend", cart.ID)
	}

	if delivery, ok := cart.GetDeliveryByCode(deliveryCode); ok {
		fmt.Printf("Inmemory Service Update %v in %#v", itemId, delivery.Cartitems)

		for _, item := range delivery.Cartitems {
			if itemId == item.ID {
				item.Qty = *itemUpdateCommand.Qty

				calculateItemPrices(&item)
				for k, currentItem := range delivery.Cartitems {
					if currentItem.ID == itemId {
						delivery.Cartitems[k] = item
					}
				}
			}
		}

		cob.CartStorage.StoreCart(*cart)
	}

	return cart, nil
}

/**
add item to cart and store in memory
*/

func (cob *InMemoryCartOrderBehaviour) AddToCart(ctx context.Context, cart *domaincart.Cart, deliveryCode string, addRequest domaincart.AddRequest) (*domaincart.Cart, error) {

	if cart != nil && !cob.CartStorage.HasCart(cart.ID) {
		return nil, fmt.Errorf("cart.infrastructure.InMemoryCartOrderBehaviour: Cannot add - Guestcart with id %v not existend", cart.ID)
	}

	// has cart current delivery, check if there is an item present for this delivery
	if cart.HasDeliveryForCode(deliveryCode) {
		delivery, _ := cart.GetDeliveryByCode(deliveryCode)

		// does the item already exist?
		itemFound := false
		for _, item := range delivery.Cartitems {
			if item.MarketplaceCode == addRequest.MarketplaceCode {
				item.Qty = item.Qty + addRequest.Qty
				itemFound = true
			}
		}

		if !itemFound {
			// create and add new item
			product, _ := cob.ProductService.Get(ctx, addRequest.MarketplaceCode)
			cartItem := domaincart.Item{
				MarketplaceCode:        addRequest.MarketplaceCode,
				VariantMarketPlaceCode: addRequest.VariantMarketplaceCode,
				Qty:          addRequest.Qty,
				SinglePrice:  product.SaleableData().ActivePrice.GetFinalPrice(),
				ID:           strconv.Itoa(rand.Int()),
				CurrencyCode: product.SaleableData().ActivePrice.Currency,
			}

			calculateItemPrices(&cartItem)
			delivery.Cartitems = append(delivery.Cartitems, cartItem)
		}

	} else {
		// create delivery and add item
		delivery := new(domaincart.Delivery)
		delivery.DeliveryInfo.Code = deliveryCode

		// create and add new item
		product, _ := cob.ProductService.Get(ctx, addRequest.MarketplaceCode)
		cartItem := domaincart.Item{
			MarketplaceCode:        addRequest.MarketplaceCode,
			VariantMarketPlaceCode: addRequest.VariantMarketplaceCode,
			Qty:          addRequest.Qty,
			SinglePrice:  product.SaleableData().ActivePrice.GetFinalPrice(),
			ID:           strconv.Itoa(rand.Int()),
			CurrencyCode: product.SaleableData().ActivePrice.Currency,
		}

		calculateItemPrices(&cartItem)

		// append item
		delivery.Cartitems = append(delivery.Cartitems, cartItem)
		cart.Deliveries = append(cart.Deliveries, *delivery)
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
func (cob *InMemoryCartOrderBehaviour) UpdateBillingAddress(ctx context.Context, cart *domaincart.Cart, billingAddress *domaincart.Address) (*domaincart.Cart, error) {
	return nil, nil
}

// @todo implement when needed
func (cob *InMemoryCartOrderBehaviour) UpdateAdditionalData(ctx context.Context, cart *domaincart.Cart, additionalData map[string]string) (*domaincart.Cart, error) {
	return nil, nil
}

func (cob *InMemoryCartOrderBehaviour) UpdateDeliveryInfo(ctx context.Context, cart *domaincart.Cart, deliveryCode string, deliveryInfo domaincart.DeliveryInfo) (*domaincart.Cart, error) {
	return nil, nil
}

// @todo implement when needed
func (cob *InMemoryCartOrderBehaviour) UpdateDeliveryInfoAdditionalData(ctx context.Context, cart *domaincart.Cart, deliveryCode string, additionalData map[string]string) (*domaincart.Cart, error) {
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
