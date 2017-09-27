package infrastructure

import (
	"flamingo/core/cart/domain/cart"
	"log"

	"math/rand"

	"fmt"

	"github.com/pkg/errors"
)

// In Session Cart Storage
type InMemoryCartService struct {
	GuestCarts map[int]cart.Cart
}

// TEst Assignement and if Interface is implemeted correct
//var _ cart.GuestCartService = InMemoryCartService{}

func (cs *InMemoryCartService) init() {
	if cs.GuestCarts == nil {
		cs.GuestCarts = make(map[int]cart.Cart)
	}
}

func (cs *InMemoryCartService) GetCart(guestcartid int) (cart.Cart, error) {
	var cart cart.Cart
	cs.init()
	if guestCart, ok := cs.GuestCarts[guestcartid]; ok {
		guestCart.CurrencyCode = "EUR"
		var total float32
		total = 0
		for _, item := range guestCart.Cartitems {
			total = total + item.RowTotal
		}
		guestCart.ShippingItem.Title = "Shipping"
		guestCart.ShippingItem.Price = 9.99
		guestCart.SubTotal = total
		guestCart.GrandTotal = total + guestCart.ShippingItem.Price
		return guestCart, nil
	}

	return cart, errors.New(fmt.Sprintf("cart.infrastructure.inmemorycartservice: Guest Cart with ID %v not exitend", guestcartid))
}

//Creates a new guest cart and returns it
func (cs *InMemoryCartService) GetNewCart() (cart.Cart, error) {
	cs.init()
	guestCart := cart.Cart{
		ID: rand.Int(),
	}
	cs.GuestCarts[guestCart.ID] = guestCart
	log.Println("New created:", cs.GuestCarts)
	return guestCart, nil
}

//TODO Get price from product package
func (cs InMemoryCartService) AddToCart(guestcartid int, marketplaceCode string, qty int) error {
	if _, ok := cs.GuestCarts[guestcartid]; !ok {
		return errors.New(fmt.Sprintf("cart.infrastructure.inmemorycartservice: Cannot add - Guestcart with id %v not existend", guestcartid))
	}
	guestcart := cs.GuestCarts[guestcartid]
	cartItem := cart.Cartitem{
		MarketplaceCode: marketplaceCode,
		Qty:             qty,
		Price:           12.99,
		RowTotal:        (12.99 * float32(qty)),
	}
	guestcart.Cartitems = append(guestcart.Cartitems, cartItem)
	cs.GuestCarts[guestcartid] = guestcart
	return nil
}

/*
// AddOrUpdateByCode if cartitem with code is already in the cart its updated. Otherwise added
func (Cart *Cart) AddOrUpdateByCode(code string, qty int, price float32) {
	for id, cartItem := range Cart.Cartitems {
		if cartItem.ProductIdendifier == code {
			cartItem.Qty = cartItem.Qty + qty
			Cart.Cartitems[id] = cartItem
			return
		}
	}
	newCartItem := Cartitem{
		code,
		qty,
		price,
	}
	Cart.Cartitems = append(Cart.Cartitems, newCartItem)
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
