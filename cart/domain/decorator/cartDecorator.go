package decorator

import (
	"context"
	"sort"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	"flamingo.me/flamingo/v3/framework/flamingo"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// DecoratedCartFactory - Factory to be injected: If you need to create a new Decorator then get the factory injected and use the factory
	DecoratedCartFactory struct {
		productService domain.ProductService
		logger         flamingo.Logger
	}

	// DecoratedCart Decorates Access To a Cart
	DecoratedCart struct {
		Cart                cartDomain.Cart
		DecoratedDeliveries []DecoratedDelivery
		Ctx                 context.Context `json:"-"`
		Logger              flamingo.Logger `json:"-"`
	}

	// DecoratedDelivery Decorates a CartItem with its Product
	DecoratedDelivery struct {
		Delivery       cartDomain.Delivery
		DecoratedItems []DecoratedCartItem
		logger         flamingo.Logger
	}

	// DecoratedCartItem Decorates a CartItem with its Product
	DecoratedCartItem struct {
		Item    cartDomain.Item
		Product domain.BasicProduct
		logger  flamingo.Logger
	}

	// GroupedDecoratedCartItem - value object used for grouping (generated on the fly)
	GroupedDecoratedCartItem struct {
		DecoratedItems []DecoratedCartItem
		Group          string
	}

	ctxKey struct{}
)

// Inject dependencies
func (df *DecoratedCartFactory) Inject(
	productService domain.ProductService,
	logger flamingo.Logger,
) {
	df.productService = productService
	df.logger = logger
}

// Create Factory method to get Decorated Cart
func (df *DecoratedCartFactory) Create(ctx context.Context, cart cartDomain.Cart) *DecoratedCart {
	decoratedCart := DecoratedCart{Cart: cart, Logger: df.logger}

	for _, d := range cart.Deliveries {
		decoratedCart.DecoratedDeliveries = append(decoratedCart.DecoratedDeliveries, DecoratedDelivery{
			Delivery:       d,
			DecoratedItems: df.CreateDecorateCartItems(context.WithValue(ctx, ctxKey{}, cart), d.Cartitems),
			logger:         df.logger,
		})
	}

	decoratedCart.Ctx = ctx

	return &decoratedCart
}

// CreateDecorateCartItems Factory method to get Decorated Cart
func (df *DecoratedCartFactory) CreateDecorateCartItems(ctx context.Context, items []cartDomain.Item) []DecoratedCartItem {
	var decoratedItems []DecoratedCartItem

	for _, cartItem := range items {
		decoratedItem := df.decorateCartItem(ctx, cartItem)
		decoratedItems = append(decoratedItems, decoratedItem)
	}

	return decoratedItems
}

// decorateCartItem factory method
func (df *DecoratedCartFactory) decorateCartItem(ctx context.Context, cartItem cartDomain.Item) DecoratedCartItem {
	decoratedItem := DecoratedCartItem{Item: cartItem, logger: df.logger}

	product, err := df.productService.Get(ctx, cartItem.MarketplaceCode)
	if err != nil {
		df.logger.WithContext(ctx).Debug("cart.decorator - no product for item: ", err)

		if product == nil {
			// To avoid errors if consumers want to access the product data
			product = domain.SimpleProduct{
				BasicProductData: domain.BasicProductData{
					Title: cartItem.ProductName + "[outdated]",
				},
			}

			decoratedItem.Product = product
		}

		return decoratedItem
	}

	if product.Type() == domain.TypeConfigurable {
		if configurable, ok := product.(domain.ConfigurableProduct); ok {
			configurableWithVariant, err := configurable.GetConfigurableWithActiveVariant(cartItem.VariantMarketPlaceCode)
			if err != nil {
				product = domain.SimpleProduct{
					BasicProductData: domain.BasicProductData{
						Title: cartItem.ProductName + "[outdated]",
					},
				}
			} else {
				product = configurableWithVariant
			}
		}
	}

	if product.Type() == domain.TypeBundle {
		if bundle, ok := product.(domain.BundleProduct); ok {
			bundleWithActiveChoices, err := bundle.GetBundleProductWithActiveChoices(cartItem.BundleConfig)
			if err != nil {
				product = domain.SimpleProduct{
					BasicProductData: domain.BasicProductData{
						Title: cartItem.ProductName + "[outdated]",
					},
				}
			} else {
				product = bundleWithActiveChoices
			}
		}
	}

	decoratedItem.Product = product

	return decoratedItem
}

// IsConfigurable - checks if current CartItem is a Configurable Product
func (dci DecoratedCartItem) IsConfigurable() bool {
	if dci.Product == nil {
		return false
	}
	return dci.Product.Type() == domain.TypeConfigurableWithActiveVariant
}

// GetVariant getter
func (dci DecoratedCartItem) GetVariant() (*domain.Variant, error) {
	return dci.Product.(domain.ConfigurableProductWithActiveVariant).Variant(dci.Item.VariantMarketPlaceCode)
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

// GetVariantsVariationAttributes getter
func (dci DecoratedCartItem) GetVariantsVariationAttributes() domain.Attributes {
	attributes := domain.Attributes{}
	if dci.IsConfigurable() {
		variant, _ := dci.GetVariant()

		for _, attributeName := range dci.Product.(domain.ConfigurableProductWithActiveVariant).VariantVariationAttributes {
			attributes[attributeName] = variant.BaseData().Attributes[attributeName]
		}
	}
	return attributes
}

// GetVariantsVariationAttributeCodes getter
func (dci DecoratedCartItem) GetVariantsVariationAttributeCodes() []string {
	if dci.Product.Type() == domain.TypeConfigurableWithActiveVariant {
		return dci.Product.(domain.ConfigurableProductWithActiveVariant).VariantVariationAttributes
	}
	return nil
}

// GetChargesToPay getter
func (dci DecoratedCartItem) GetChargesToPay(wishedToPaySum *domain.WishedToPay) priceDomain.Charges {
	priceToPayForItem := dci.Item.RowPriceGrossWithDiscount
	return dci.Product.SaleableData().GetLoyaltyChargeSplit(&priceToPayForItem, wishedToPaySum, dci.Item.Qty)
}

// GetGroupedBy legacy function
// deprecated: only here to support the old structure of accesing DecoratedItems in the Decorated Cart
// Use instead:
//
//	or iterate over DecoratedCart.DecoratedDelivery
func (dc DecoratedCart) GetGroupedBy(group string, sortGroup bool, params ...string) []*GroupedDecoratedCartItem {

	if dc.Logger != nil {
		dc.Logger.Warn("DEPRECATED: DecoratedCart.GetGroupedBy()")
	}
	if len(dc.DecoratedDeliveries) != 1 {
		return nil
	}
	return dc.DecoratedDeliveries[0].GetGroupedBy(group, sortGroup, params...)
}

// GetAllDecoratedItems getter
func (dc DecoratedCart) GetAllDecoratedItems() []DecoratedCartItem {
	var allItems []DecoratedCartItem
	for _, dd := range dc.DecoratedDeliveries {
		allItems = append(allItems, dd.DecoratedItems...)
	}
	return allItems
}

// GetDecoratedDeliveryByCode getter
func (dc DecoratedCart) GetDecoratedDeliveryByCode(deliveryCode string) (*DecoratedDelivery, bool) {
	for _, dd := range dc.DecoratedDeliveries {
		if dd.Delivery.DeliveryInfo.Code == deliveryCode {
			return &dd, true
		}

	}
	return nil, false
}

// GetDecoratedDeliveryByCodeWithoutBool - used inside a template, therefor we need the method with a single return param
func (dc DecoratedCart) GetDecoratedDeliveryByCodeWithoutBool(deliveryCode string) *DecoratedDelivery {
	decoratedDelivery, _ := dc.GetDecoratedDeliveryByCode(deliveryCode)
	return decoratedDelivery
}

// GetGroupedBy getter
func (dc DecoratedDelivery) GetGroupedBy(group string, sortGroup bool, params ...string) []*GroupedDecoratedCartItem {
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

	// sort before return
	if sortGroup {
		direction := ""
		if len(params) > 0 {
			direction = params[0]
		}

		if direction == "DESC" {
			sort.Sort(sort.Reverse(sort.StringSlice(groupedItemsCollectionKeys)))
		} else {
			sort.Strings(groupedItemsCollectionKeys)
		}
	}

	var groupedItemsCollectionSorted []*GroupedDecoratedCartItem
	for _, keyName := range groupedItemsCollectionKeys {
		if groupedItemsEntry, ok := groupedItemsCollection[keyName]; ok {
			groupedItemsCollectionSorted = append(groupedItemsCollectionSorted, groupedItemsEntry)
		}
	}
	return groupedItemsCollectionSorted
}

// GetDecoratedCartItemByID getter
func (dc DecoratedDelivery) GetDecoratedCartItemByID(ID string) *DecoratedCartItem {
	for _, decoratedItem := range dc.DecoratedItems {
		if decoratedItem.Item.ID == ID {
			return &decoratedItem
		}
	}
	return nil
}

func CartFromDecoratedCartFactoryContext(ctx context.Context) *cartDomain.Cart {
	if cart, ok := ctx.Value(ctxKey{}).(cartDomain.Cart); ok {
		return &cart
	}

	return nil
}
