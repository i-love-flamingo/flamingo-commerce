package domain

import (
	"time"
)

type (
	// Order struct
	Order struct {
		ID           string
		CreationTime time.Time
		UpdateTime   time.Time
		OrderItems   []*OrderItem
		Status       string
		Total        float64
		CurrencyCode string
		Attributes   Attributes
	}

	// OrderItem struct
	OrderItem struct {
		// DEPRECATED
		Sku string

		MarketplaceCode        string
		VariantMarketplaceCode string

		Qty float64

		CurrencyCode       string
		SinglePrice        float64
		SinglePriceInclTax float64
		RowTotal           float64
		TaxAmount          float64
		RowTotalInclTax    float64

		Name         string
		Price        float64
		PriceInclTax float64

		// Source Id where the item shoudl be picked
		SourceID string
	}

	// Attributes map
	Attributes map[string]Attribute

	// Attribute interface
	Attribute interface{}

	// PlacedOrderInfos represents a slice of PlacedOrderInfo
	PlacedOrderInfos []PlacedOrderInfo

	// PlacedOrderInfo defines the additional info struct for placed orders
	PlacedOrderInfo struct {
		OrderNumber  string
		DeliveryCode string
	}
)

// GetOrderNumberForDeliveryCode returns the order number for a delivery code
func (poi PlacedOrderInfos) GetOrderNumberForDeliveryCode(deliveryCode string) string {
	for _, v := range poi {
		if v.DeliveryCode == deliveryCode {
			return v.OrderNumber
		}
	}
	return ""
}
