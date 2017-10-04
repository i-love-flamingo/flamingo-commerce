package cart

import (
	"context"

	"go.aoe.com/flamingo/core/product/domain"

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
	product, e := productService.Get(ctx, cartitem.MarketplaceCode)
	if e != nil {
		log.Println("cart.decorator - no product for item:", e)
		return decorateditem
	}
	decorateditem.Product = product
	return decorateditem
}

// IsConfigurable - checks if current CartItem is a Configurable Product
func (dci DecoratedCartItem) IsConfigurable() bool {
	return dci.Product.Type() == domain.TYPECONFIGURABLE
}

// GetVariant
func (dci DecoratedCartItem) GetVariant() (*domain.Variant, error) {
	return dci.Product.(domain.ConfigurableProduct).Variant(dci.Cartitem.VariantMarketPlaceCode)
}

// GetDisplayTitle
func (dci DecoratedCartItem) GetDisplayTitle() string {
	if dci.IsConfigurable() {
		variant, e := dci.GetVariant()
		if e != nil {
			return "Error Getting Variant"
		}
		return variant.Title
	}
	return dci.Product.BaseData().Title
}

// GetDisplayMarketplaceCode
func (dci DecoratedCartItem) GetDisplayMarketplaceCode() string {
	if dci.IsConfigurable() {
		variant, e := dci.GetVariant()
		if e != nil {
			return "Error Getting Variant"
		}
		return variant.MarketPlaceCode
	}
	return dci.Product.BaseData().MarketPlaceCode
}

// GetVariantsVariationAttribute
func (dci DecoratedCartItem) GetVariantsVariationAttributes() domain.Attributes {
	attributes := domain.Attributes{}
	if dci.IsConfigurable() {
		variant, _ := dci.GetVariant()

		for _, attributeName := range dci.Product.(domain.ConfigurableProduct).VariantVariationAttributes {
			attributes[attributeName] = variant.BaseData().Attributes[attributeName]
		}
	}
	log.Println(attributes)
	return attributes
}
