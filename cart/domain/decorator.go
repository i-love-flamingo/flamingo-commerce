package domain

import (
	"context"
	"flamingo/core/product/domain"

	"log"
)

type (
	// DecoratedCartFactory - Factory to be injected: If you need to create a new Decorator then get the factory injected and use the factory
	DecoratedCartFactory struct {
		ProductService domain.ProductService `inject:""`
	}

	// DecoratedCart Decorates Access To a Cart
	DecoratedCart struct {
		Cart
		Cartitems []DecoratedCartItem
		Ctx       context.Context `json:"-"`
	}

	// DecoratedCartItem Decorates a CartItem with its Product
	DecoratedCartItem struct {
		Cartitem
		Product domain.BasicProduct
	}
)

// CreateDecoratedCart Native Factory
func CreateDecoratedCart(ctx context.Context, Cart Cart, productService domain.ProductService) *DecoratedCart {

	DecoratedCart := DecoratedCart{Cart: Cart}
	for _, cartitem := range Cart.Cartitems {
		decoratedItem := decorateCartItem(ctx, cartitem, productService)
		DecoratedCart.Cartitems = append(DecoratedCart.Cartitems, decoratedItem)
	}
	DecoratedCart.Ctx = ctx
	return &DecoratedCart
}

// Create Factory - with injected ProductService
func (df *DecoratedCartFactory) Create(ctx context.Context, Cart Cart) *DecoratedCart {
	return CreateDecoratedCart(ctx, Cart, df.ProductService)
}

//decorateCartItem factory method
func decorateCartItem(ctx context.Context, cartitem Cartitem, productService domain.ProductService) DecoratedCartItem {
	decorateditem := DecoratedCartItem{Cartitem: cartitem}
	product, e := productService.Get(ctx, cartitem.ProductCode)
	if e != nil {
		log.Println("cart.decorator - no product for item:", e)
		return decorateditem
	}
	decorateditem.Product = product
	return decorateditem
}
