package infrastructure

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
)

// In Session Cart Storage
type InMemoryCartServiceAdapter struct {
	GuestCarts     map[string]cart.Cart
	ProductService productDomain.ProductService `inject:""`
}

// Test Assignement and if Interface is implemeted correct
var _ cart.GuestCartService = cart.GuestCartService(&InMemoryCartServiceAdapter{})

func (cs *InMemoryCartServiceAdapter) init() {
	if cs.GuestCarts == nil {
		cs.GuestCarts = make(map[string]cart.Cart)
	}

}

func (cs *InMemoryCartServiceAdapter) GetCart(ctx context.Context, guestcartid string) (cart.Cart, error) {
	var cart cart.Cart
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

	return cart, errors.New(fmt.Sprintf("cart.infrastructure.inmemorycartservice: Guest Cart with ID %v not exitend", guestcartid))
}

//Creates a new guest cart and returns it
func (cs *InMemoryCartServiceAdapter) GetNewCart(ctx context.Context) (cart.Cart, error) {
	cs.init()
	guestCart := cart.Cart{
		ID: strconv.Itoa(rand.Int()),
	}
	cs.GuestCarts[guestCart.ID] = guestCart
	log.Println("New created:", cs.GuestCarts)
	return guestCart, nil
}

//AddToCart
func (cs InMemoryCartServiceAdapter) AddToCart(ctx context.Context, guestcartid string, addRequest cart.AddRequest) error {
	if _, ok := cs.GuestCarts[guestcartid]; !ok {
		return errors.New(fmt.Sprintf("cart.infrastructure.inmemorycartservice: Cannot add - Guestcart with id %v not existend", guestcartid))
	}
	guestcart := cs.GuestCarts[guestcartid]
	found, lineNr := guestcart.HasItem(addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode)
	if found {
		item, _ := guestcart.GetByLineNr(lineNr)
		item.Qty = item.Qty + addRequest.Qty
	} else {
		product, _ := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
		rowTotal, _ := new(big.Float).Mul(big.NewFloat(product.SaleableData().ActivePrice.GetFinalPrice()), new(big.Float).SetInt64(int64(addRequest.Qty))).Float64()
		cartItem := cart.Item{
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

/*
// AddOrUpdateByCode if cartitem with code is already in the cart its updated. Otherwise added
func (Cart *Cart) AddOrUpdateByCode(code string, qty int, price float32) {
	for id, cartItem := range Cart.DecoratedItems {
		if cartItem.ProductIdendifier == code {
			cartItem.Qty = cartItem.Qty + qty
			Cart.DecoratedItems[id] = cartItem
			return
		}
	}
	newCartItem := Item{
		code,
		qty,
		price,
	}
	Cart.DecoratedItems = append(Cart.DecoratedItems, newCartItem)
}
*/

/*


// FakecartrepositoryFactory factory
func FakecartrepositoryFactory() *Fakecartrepository {
	return &Fakecartrepository{
		GuestCarts: make(map[int]*domain.Cart),
	}
}


// Add to cart
func (cr *Fakecartrepository) Add(Cart domain.Cart) (int, error) {
	cr.init()
	fmt.Println("Fake cartrepo impl called add")
	if Cart.ID == 0 {
		Cart.ID = rand.Int()
	}
	cr.GuestCarts[Cart.ID] = &Cart
	return Cart.ID, nil
}

// Update cart
func (cr *Fakecartrepository) Update(Cart domain.Cart) error {
	cr.init()
	fmt.Println("Fake cartrepo impl called Update")
	cr.GuestCarts[Cart.ID] = &Cart
	return nil
}

// Delete cart
func (cr *Fakecartrepository) Delete(Cart domain.Cart) error {
	cr.init()
	fmt.Println("Fake cartrepo impl called delete")
	delete(cr.GuestCarts, Cart.ID)
	return nil
}

// Get cart
func (cr *Fakecartrepository) Get(ID int) (*domain.Cart, error) {
	cr.init()
	fmt.Printf("Fake cartrepo impl called get for %s", ID)
	if val, ok := cr.GuestCarts[ID]; ok {
		return val, nil
	}

	return nil, errors.New("No cart with that ID")
}
*/
