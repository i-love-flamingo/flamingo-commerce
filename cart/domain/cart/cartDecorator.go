package cart

import (
	"context"
	"log"
	"sort"

	"go.aoe.com/flamingo/core/product/domain"
)

type (
	// DecoratedCartFactory - Factory to be injected: If you need to create a new Decorator then get the factory injected and use the factory
	DecoratedCartFactory struct {
		ProductService domain.ProductService `inject:""`
	}

	// DecoratedCart Decorates Access To a Cart
	DecoratedCart struct {
		Cart           Cart
		DecoratedItems []DecoratedCartItem
		Ctx            context.Context `json:"-"`
	}

	// DecoratedCartItem Decorates a CartItem with its Product
	GroupedDecoratedCartItem struct {
		DecoratedItems []DecoratedCartItem
		Group          string
	}

	// DecoratedCartItem Decorates a CartItem with its Product
	DecoratedCartItem struct {
		Item    Item
		Product domain.BasicProduct
	}
)

// CreateDecoratedCart Native Factory
func CreateDecoratedCart(ctx context.Context, Cart Cart, productService domain.ProductService) *DecoratedCart {
	decoratedCart := DecoratedCart{Cart: Cart}
	for _, cartitem := range Cart.Cartitems {
		decoratedItem := decorateCartItem(ctx, cartitem, productService)
		decoratedCart.DecoratedItems = append(decoratedCart.DecoratedItems, decoratedItem)
	}
	decoratedCart.Ctx = ctx
	return &decoratedCart
}

// Create Factory - with injected ProductService
func (df *DecoratedCartFactory) Create(ctx context.Context, Cart Cart) *DecoratedCart {
	return CreateDecoratedCart(ctx, Cart, df.ProductService)
}

//decorateCartItem factory method
func decorateCartItem(ctx context.Context, cartitem Item, productService domain.ProductService) DecoratedCartItem {
	decorateditem := DecoratedCartItem{Item: cartitem}
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

// GetVariant getter
func (dci DecoratedCartItem) GetVariant() (*domain.Variant, error) {
	return dci.Product.(domain.ConfigurableProduct).Variant(dci.Item.VariantMarketPlaceCode)
}

// GetDisplayTitle getter
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

// GetDisplayMarketplaceCode getter
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

// GetVariantsVariationAttribute getter
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

// GetGroupedBy getter
func (dc DecoratedCart) GetGroupedBy(group string) []*GroupedDecoratedCartItem {
	groupedItemsCollection := make(map[string]*GroupedDecoratedCartItem)
	var groupedItemsCollectionKeys []string

	var groupKey string
	for _, item := range dc.DecoratedItems {
		switch group {
		case "retailer_code":
			groupKey = item.Product.BaseData().RetailerCode
		default:
			groupKey = "default"
		}
		if _, ok := groupedItemsCollection[groupKey]; !ok {
			groupedItemsCollection[groupKey] = &GroupedDecoratedCartItem{
				Group: groupKey,
			}
			groupedItemsCollectionKeys = append(groupedItemsCollectionKeys, groupKey)
		}

		if groupedItemsEntry, ok := groupedItemsCollection[groupKey]; ok {
			groupedItemsEntry.DecoratedItems = append(groupedItemsEntry.DecoratedItems, item)
		}
	}

	//sort before return
	var groupedItemsCollectionSorted []*GroupedDecoratedCartItem
	sort.Strings(groupedItemsCollectionKeys)
	for _, keyName := range groupedItemsCollectionKeys {
		if groupedItemsEntry, ok := groupedItemsCollection[keyName]; ok {
			groupedItemsCollectionSorted = append(groupedItemsCollectionSorted, groupedItemsEntry)
		}
	}
	return groupedItemsCollectionSorted
}
