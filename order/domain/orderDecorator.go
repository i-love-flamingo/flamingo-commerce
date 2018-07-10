package domain

import (
	"context"
	"sort"

	"flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	OrderDecoratorInterface interface {
		Create(context.Context, *Order) *DecoratedOrder
	}

	// OrderDecorator
	OrderDecorator struct {
		ProductService domain.ProductService `inject:""`
		Logger         flamingo.Logger       `inject:""`
	}

	// DecoratedOrder
	DecoratedOrder struct {
		Order          *Order
		DecoratedItems []*DecoratedOrderItem
	}

	// DecoratedOrderItem
	DecoratedOrderItem struct {
		Item    *OrderItem
		Product *domain.BasicProduct
	}

	// GroupedDecoratedOrder
	GroupedDecoratedOrder struct {
		Order  *DecoratedOrder
		Groups []*GroupedDecoratedOrderItems
	}

	// GroupedDecoratedOrderItem
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
		rd.Logger.Error("order.decorator - no product for item", err)
		// fallback to return something the frontend still could use
		product = rd.createFallbackProduct(item)
	case product.Type() == domain.TYPECONFIGURABLE:
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
	result.Product = &product

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
	return (*doi.Product).Type() == domain.TYPECONFIGURABLE_WITH_ACTIVE_VARIANT
}

// GetVariant getter
func (doi DecoratedOrderItem) GetVariant() (*domain.Variant, error) {
	return (*doi.Product).(domain.ConfigurableProductWithActiveVariant).Variant(doi.Item.VariantMarketplaceCode)
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
	return (*doi.Product).BaseData().Title
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
	return (*doi.Product).BaseData().MarketPlaceCode
}

// GetVariantsVariationAttribute getter
func (doi DecoratedOrderItem) GetVariantsVariationAttributes() domain.Attributes {
	attributes := domain.Attributes{}
	if doi.IsConfigurable() {
		variant, _ := doi.GetVariant()

		for _, attributeName := range (*doi.Product).(domain.ConfigurableProductWithActiveVariant).VariantVariationAttributes {
			attributes[attributeName] = variant.BaseData().Attributes[attributeName]
		}
	}
	return attributes
}

// GetVariantsVariationAttribute getter
func (doi DecoratedOrderItem) GetVariantsVariationAttributeCodes() []string {
	return (*doi.Product).(domain.ConfigurableProductWithActiveVariant).VariantVariationAttributes
}

// GetGroupedBy
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
			groupKey = (*item.Product).BaseData().RetailerCode
		default:
			groupKey = "default"
		}

		if _, ok := groupedItemsCollection[groupKey]; !ok {
			groupedItemsCollection[groupKey] = &GroupedDecoratedOrderItems{
				Group: groupKey,
			}

			groupedItemKeys = append(groupedItemKeys, groupKey)
		}

		groupedItemsEntry, _ := groupedItemsCollection[groupKey]
		groupedItemsEntry.DecoratedItems = append(groupedItemsEntry.DecoratedItems, item)
	}

	if sortGroup {
		sort.Strings(groupedItemKeys)
	}

	groups := make([]*GroupedDecoratedOrderItems, len(groupedItemKeys))
	for i, key := range groupedItemKeys {
		groupedItemsEntry, _ := groupedItemsCollection[key]
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
		sourceId := item.Item.SourceId
		if _, ok := sourceIds[sourceId]; ok {
			continue
		}

		sourceIds[sourceId] = true
		result = append(result, sourceId)
	}

	return result
}
