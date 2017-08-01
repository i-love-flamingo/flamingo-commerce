package domain

import (
	"context"
	"flamingo/core/product/domain"

	"log"
)

type (
	// Factory to be injected: If you need to create a new Decorator then get the factory injected and use the factory
	DecoratorFactory struct {
		ProductService domain.ProductService `inject:""`
	}

	// Decorates Access To a Cart
	DecoratedCart struct {
		Cart
		Cartitems      []DecoratedCartItem
		Ctx            context.Context       `json:"-"`
		ProductService domain.ProductService `json:"-"`
	}

	// Decorates Access To a Cart
	DecoratedCartItem struct {
		Cartitem
		Product *domain.Product
	}
)

// Native Factory
func CreateDecoratedCart(ctx context.Context, Cart Cart, productService domain.ProductService) *DecoratedCart {

	DecoratedCart := DecoratedCart{Cart: Cart}
	for _, cartitem := range Cart.Cartitems {
		decoratedItem := decorateCartItem(ctx, cartitem, productService)
		DecoratedCart.Cartitems = append(DecoratedCart.Cartitems, decoratedItem)
	}
	DecoratedCart.ProductService = productService
	DecoratedCart.Ctx = ctx
	return &DecoratedCart
}

// Factory - with injected ProductService
func (df *DecoratorFactory) Create(ctx context.Context, Cart Cart) *DecoratedCart {
	return CreateDecoratedCart(ctx, Cart, df.ProductService)
}

func decorateCartItem(ctx context.Context, cartitem Cartitem, productService domain.ProductService) DecoratedCartItem {
	decorateditem := DecoratedCartItem{Cartitem: cartitem}
	product, e := productService.Get(ctx, cartitem.ProductCode)
	if e != nil {
		log.Println("cart.decorator:", e)
		return decorateditem
	}
	decorateditem.Product = product
	return decorateditem
}
