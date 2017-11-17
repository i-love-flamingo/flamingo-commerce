package infrastructure

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strconv"

	"time"

	domaincart "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
)

// In Session Cart Storage
type (
	InMemoryCartServiceAdapter struct {
		GuestCarts     map[string]domaincart.Cart
		ProductService productDomain.ProductService `inject:""`
	}
	CartOrderBehaviour struct{}
)

// Test Assignement and if Interface is implemeted correct
var _ domaincart.GuestCartService = &InMemoryCartServiceAdapter{}

func (cs *InMemoryCartServiceAdapter) init() {
	if cs.GuestCarts == nil {
		cs.GuestCarts = make(map[string]domaincart.Cart)
	}

}

func (cs *InMemoryCartServiceAdapter) GetCart(ctx context.Context, guestcartid string) (domaincart.Cart, error) {
	var cart domaincart.Cart
	cart.CartOrderBehaviour = new(CartOrderBehaviour)
	cs.init()
	if guestCart, ok := cs.GuestCarts[guestcartid]; ok {
		guestCart.CurrencyCode = "EUR"
		var total big.Float
		total.SetFloat64(0)
		for _, item := range guestCart.Cartitems {
			total.Add(&total, big.NewFloat(item.RowTotal))
		}
		guestCart.ShippingItem.Title = "Shipping"
		guestCart.ShippingItem.Price = 9.99
		guestCart.SubTotal, _ = total.Float64()
		guestCart.GrandTotal, _ = new(big.Float).Add(&total, big.NewFloat(guestCart.ShippingItem.Price)).Float64()
		return guestCart, nil
	}

	return cart, fmt.Errorf("cart.infrastructure.inmemorycartservice: Guest Cart with ID %v not exitend", guestcartid)
}

//Creates a new guest cart and returns it
func (cs *InMemoryCartServiceAdapter) GetNewCart(ctx context.Context) (domaincart.Cart, error) {
	cs.init()
	guestCart := domaincart.Cart{
		ID: strconv.Itoa(rand.Int()),
	}
	guestCart.CartOrderBehaviour = new(CartOrderBehaviour)
	cs.GuestCarts[guestCart.ID] = guestCart
	log.Println("New created:", cs.GuestCarts)
	return guestCart, nil
}

//AddToCart
func (cs InMemoryCartServiceAdapter) AddToCart(ctx context.Context, guestcartid string, addRequest domaincart.AddRequest) error {
	if _, ok := cs.GuestCarts[guestcartid]; !ok {
		return fmt.Errorf("cart.infrastructure.inmemorycartservice: Cannot add - Guestcart with id %v not existend", guestcartid)
	}
	guestcart := cs.GuestCarts[guestcartid]
	found, lineNr := guestcart.HasItem(addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode)
	if found {
		item, _ := guestcart.GetByLineNr(lineNr)
		item.Qty = item.Qty + addRequest.Qty
	} else {
		product, _ := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
		rowTotal, _ := new(big.Float).Mul(big.NewFloat(product.SaleableData().ActivePrice.GetFinalPrice()), new(big.Float).SetInt64(int64(addRequest.Qty))).Float64()
		cartItem := domaincart.Item{
			MarketplaceCode:        addRequest.MarketplaceCode,
			VariantMarketPlaceCode: addRequest.VariantMarketplaceCode,
			Qty:      addRequest.Qty,
			Price:    product.SaleableData().ActivePrice.GetFinalPrice(),
			RowTotal: rowTotal,
		}
		guestcart.Cartitems = append(guestcart.Cartitems, cartItem)
		cs.GuestCarts[guestcartid] = guestcart
	}

	return nil
}

// SetShippingInformation adds a product
func (cs *CartOrderBehaviour) SetShippingInformation(ctx context.Context, cart *domaincart.Cart, shippingAddress *domaincart.Address, billingAddress *domaincart.Address, shippingCarrierCode string, shippingMethodCode string) error {
	return nil
}

// SetShippingInformation adds a product
func (cs *CartOrderBehaviour) PlaceOrder(ctx context.Context, cart *domaincart.Cart, payment *domaincart.Payment) (string, error) {
	rand.Seed(time.Now().Unix())
	return string(rand.Int()), nil
}

func (cs *CartOrderBehaviour) DeleteItem(ctx context.Context, cart *domaincart.Cart, itemId string) error {
	for k, item := range cart.Cartitems {
		if item.ID == itemId {
			if len(cart.Cartitems) > k {
				cart.Cartitems = append(cart.Cartitems[:k], cart.Cartitems[k+1:]...)
			} else {
				cart.Cartitems = cart.Cartitems[:k]
			}
		}
	}
	return nil
}
