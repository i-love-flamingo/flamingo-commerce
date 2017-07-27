package domain

import "flamingo/core/product/domain"

type (
	// Factory to be injected: If you need to create a new Decorator then get the factory injected and use the factory
	DecoratorFactory struct {
		ProductService domain.ProductService `inject:""`
	}

	// Decorates Access To a Cart
	DecoratedCart struct {
		ProductService domain.ProductService
		Cart           Cart
	}

	// Decorates Access To a Cart
	DecoratedCartItem struct {
		Product  domain.Product
		Cartitem Cartitem
	}
)

// Native Factory
func CreateCartDecorator(Cart Cart, ProductService domain.ProductService) *DecoratedCart {
	var Decorator DecoratedCart
	Decorator.Cart = Cart
	Decorator.ProductService = ProductService
	return &Decorator
}

// Factory - with injected ProductService
func (df *DecoratorFactory) Create(Cart Cart) *DecoratedCart {
	return CreateCartDecorator(Cart, df.ProductService)
}

// GetLine gets an item - starting with 1
func (Decorator *DecoratedCart) GetLine(lineNr int) Cartitem {
	return Cart.Cartitems[lineNr-1]
}
