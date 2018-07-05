package domain

import (
	"context"
	"log"
	"sort"

	"flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	OrderDecoratorInterface interface {
		Create(context.Context, Order) DecoratedOrder
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

	//
	GroupedDecoratedOrderItem struct {
		DecoratedItems []*DecoratedOrderItem
		Group          string
	}
)

// Create creates a new decorated order
func (rd *OrderDecorator) Create(ctx context.Context, order *Order) DecoratedOrder {
	result := DecoratedOrder{Order: order}
	result.DecoratedItems = rd.createDecoratedItems(ctx, &order.OrderItems)

	return result
}

func (rd *OrderDecorator) createDecoratedItems(ctx context.Context, items *[]OrderItem) []*DecoratedOrderItem {
	result := make([]*DecoratedOrderItem, len(*items))
	for i, item := range *items {
		decoratedItem := rd.createDecoratedItem(ctx, &item)
		result[i] = decoratedItem
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
	log.Println(attributes)
	return attributes
}

// GetVariantsVariationAttribute getter
func (doi DecoratedOrderItem) GetVariantsVariationAttributeCodes() []string {
	return (*doi.Product).(domain.ConfigurableProductWithActiveVariant).VariantVariationAttributes
}

// GetGroupedBy
func (rd *DecoratedOrder) GetGroupedBy(group string, sortGroup bool) []GroupedDecoratedOrderItem {
	groupedItemsCollection := make(map[string]GroupedDecoratedOrderItem)
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
			groupedItemsCollection[groupKey] = GroupedDecoratedOrderItem{
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

	result := make([]GroupedDecoratedOrderItem, len(groupedItemKeys))
	for _, key := range groupedItemKeys {
		groupedItemsEntry, _ := groupedItemsCollection[key]
		result = append(result, groupedItemsEntry)
	}

	return result
}
