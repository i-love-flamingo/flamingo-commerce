package domain

import (
	"context"
	"sort"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// OrderDecoratorInterface defines the interface of the order decorator
	OrderDecoratorInterface interface {
		Create(context.Context, *Order) *DecoratedOrder
	}

	// OrderDecorator struct defines the order decorator
	OrderDecorator struct {
		ProductService domain.ProductService `inject:""`
		Logger         flamingo.Logger       `inject:""`
	}

	// DecoratedOrder struct
	DecoratedOrder struct {
		Order          *Order
		DecoratedItems []*DecoratedOrderItem
	}

	// DecoratedOrderItem struct
	DecoratedOrderItem struct {
		Item    *OrderItem
		Product domain.BasicProduct
	}

	// GroupedDecoratedOrder struct
	GroupedDecoratedOrder struct {
		Order  *DecoratedOrder
		Groups []*GroupedDecoratedOrderItems
	}

	// GroupedDecoratedOrderItems struct
	GroupedDecoratedOrderItems struct {
		DecoratedItems []*DecoratedOrderItem
		Group          string
	}
)

// check interface implementation
var _ OrderDecoratorInterface = (*OrderDecorator)(nil)

// Create creates a new decorated order
func (rd *OrderDecorator) Create(ctx context.Context, order *Order) *DecoratedOrder {
	result := &DecoratedOrder{Order: order}
	result.DecoratedItems = rd.createDecoratedItems(ctx, order.OrderItems)

	return result
}

func (rd *OrderDecorator) createDecoratedItems(ctx context.Context, items []*OrderItem) []*DecoratedOrderItem {
	result := make([]*DecoratedOrderItem, len(items))
	for i, item := range items {
		result[i] = rd.createDecoratedItem(ctx, item)
	}

	return result
}

func (rd *OrderDecorator) createDecoratedItem(ctx context.Context, item *OrderItem) *DecoratedOrderItem {
	result := &DecoratedOrderItem{
		Item: item,
	}

	product, err := rd.ProductService.Get(ctx, item.MarketplaceCode)
	switch {
	case err != nil:
		rd.Logger.WithContext(ctx).Error("order.decorator - no product for item", err)
		// fallback to return something the frontend still could use
		product = rd.createFallbackProduct(item)
	case product.Type() == domain.TypeConfigurable && item.VariantMarketplaceCode != "":
		configurable, ok := product.(domain.ConfigurableProduct)
		if !ok {
			// not a usable configrable
			break
		}

		variant, err := configurable.GetConfigurableWithActiveVariant(item.VariantMarketplaceCode)
		if err == nil {
			product = variant
		} else {
			product = rd.createFallbackProduct(item)
		}
	}
	result.Product = product

	return result
}

func (rd *OrderDecorator) createFallbackProduct(item *OrderItem) *domain.SimpleProduct {
	return &domain.SimpleProduct{
		BasicProductData: domain.BasicProductData{
			Title: item.Name,
		},
		Saleable: domain.Saleable{
			IsSaleable: false,
		},
	}
}

// IsConfigurable - checks if current order item is a configurable product
func (doi DecoratedOrderItem) IsConfigurable() bool {
	return doi.Product.Type() == domain.TypeConfigurableWithActiveVariant
}

// GetVariant getter
func (doi DecoratedOrderItem) GetVariant() (*domain.Variant, error) {
	return doi.Product.(domain.ConfigurableProductWithActiveVariant).Variant(doi.Item.VariantMarketplaceCode)
}

// GetDisplayTitle getter
func (doi DecoratedOrderItem) GetDisplayTitle() string {
	if doi.IsConfigurable() {
		variant, e := doi.GetVariant()
		if e != nil {
			return "Error Getting Variant"
		}
		return variant.Title
	}
	return doi.Product.BaseData().Title
}

// GetDisplayMarketplaceCode getter
func (doi DecoratedOrderItem) GetDisplayMarketplaceCode() string {
	if doi.IsConfigurable() {
		variant, e := doi.GetVariant()
		if e != nil {
			return "Error Getting Variant"
		}
		return variant.MarketPlaceCode
	}
	return doi.Product.BaseData().MarketPlaceCode
}

// GetVariantsVariationAttributes gets the decorated order item variant attributes
func (doi DecoratedOrderItem) GetVariantsVariationAttributes() domain.Attributes {
	attributes := domain.Attributes{}
	if doi.IsConfigurable() {
		variant, _ := doi.GetVariant()

		for _, attributeName := range doi.Product.(domain.ConfigurableProductWithActiveVariant).VariantVariationAttributes {
			attributes[attributeName] = variant.BaseData().Attributes[attributeName]
		}
	}
	return attributes
}

// GetVariantsVariationAttributeCodes gets the decorated order item variant variation attributes
func (doi DecoratedOrderItem) GetVariantsVariationAttributeCodes() []string {
	if doi.Product.Type() == domain.TypeConfigurableWithActiveVariant {
		return doi.Product.(domain.ConfigurableProductWithActiveVariant).VariantVariationAttributes
	}
	return nil
}

// GetGroupedBy groups the decorated order into a *GroupedDecoratedOrder
func (rd *DecoratedOrder) GetGroupedBy(group string, sortGroup bool) *GroupedDecoratedOrder {
	result := &GroupedDecoratedOrder{
		Order: rd,
	}
	groupedItemsCollection := make(map[string]*GroupedDecoratedOrderItems)
	var groupedItemKeys []string

	var groupKey string
	for _, item := range rd.DecoratedItems {
		switch group {
		case "retailer_code":
			groupKey = item.Product.BaseData().RetailerCode
		default:
			groupKey = "default"
		}

		if _, ok := groupedItemsCollection[groupKey]; !ok {
			groupedItemsCollection[groupKey] = &GroupedDecoratedOrderItems{
				Group: groupKey,
			}

			groupedItemKeys = append(groupedItemKeys, groupKey)
		}

		groupedItemsEntry := groupedItemsCollection[groupKey]
		groupedItemsEntry.DecoratedItems = append(groupedItemsEntry.DecoratedItems, item)
	}

	if sortGroup {
		sort.Strings(groupedItemKeys)
	}

	groups := make([]*GroupedDecoratedOrderItems, len(groupedItemKeys))
	for i, key := range groupedItemKeys {
		groupedItemsEntry := groupedItemsCollection[key]
		groups[i] = groupedItemsEntry
	}
	result.Groups = groups

	return result
}

// GetSourceIds collects the source ids of the items of the group
func (i *GroupedDecoratedOrderItems) GetSourceIds() []string {
	// the group has at least one group in there
	sourceIds := make(map[string]bool, 1)
	result := make([]string, 1)
	for _, item := range i.DecoratedItems {
		sourceID := item.Item.SourceID
		if _, ok := sourceIds[sourceID]; ok {
			continue
		}

		sourceIds[sourceID] = true
		result = append(result, sourceID)
	}

	return result
}
