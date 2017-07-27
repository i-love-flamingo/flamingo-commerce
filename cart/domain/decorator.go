package domain

import (
	"context"
	"flamingo/core/product/domain"
)

type (
	// Factory to be injected: If you need to create a new Decorator then get the factory injected and use the factory
	DecoratorFactory struct {
		ProductService domain.ProductService `inject:""`
	}

	// Decorates Access To a Cart
	DecoratedCart struct {
		ProductService domain.ProductService
		Cart           Cart
		Ctx            context.Context
	}

	// Decorates Access To a Cart
	DecoratedCartItem struct {
		Product  *domain.Product
		Cartitem Cartitem
	}
)

// Native Factory
func CreateCartDecorator(ctx context.Context, Cart Cart, ProductService domain.ProductService) *DecoratedCart {
	var Decorator DecoratedCart
	Decorator.Cart = Cart
	Decorator.ProductService = ProductService
	Decorator.Ctx = ctx
	return &Decorator
}

// Factory - with injected ProductService
func (df *DecoratorFactory) Create(ctx context.Context, Cart Cart) *DecoratedCart {
	return CreateCartDecorator(ctx, Cart, df.ProductService)
}

// GetLine gets an item - starting with 1
func (Decorator *DecoratedCart) GetLine(lineNr int) (DecoratedCartItem, error) {
	var decorateditem DecoratedCartItem
	item, e := Decorator.Cart.GetLine(lineNr)

	if e != nil {
		return decorateditem, e
	}
	product, e := Decorator.ProductService.Get(Decorator.Ctx, item.ProductCode)
	if e != nil {
		return decorateditem, e
	}
	return DecoratedCartItem{Cartitem: Decorator.Cart.Cartitems[lineNr-1], Product: product}, nil
}
