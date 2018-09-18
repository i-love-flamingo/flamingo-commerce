package domain

import (
	"time"
)

type (
	// Order
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

	// OrderItem
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
		SourceId string
	}

	// Attributes
	Attributes map[string]Attribute

	// Attribute
	Attribute interface{}
)
