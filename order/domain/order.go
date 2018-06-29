package domain

import "time"

type (
	// Order
	Order struct {
		ID           string
		CreationTime time.Time
		UpdateTime   time.Time
		OrderItems   []OrderItem
		Status       string
		Total        float64
		CurrencyCode string
		Attributes   Attributes
	}
	// OrderItem
	OrderItem struct {
		Sku          string
		Name         string
		Price        float64
		PriceInclTax float64
	}

	// Attributes
	Attributes map[string]Attribute

	// Attribute
	Attribute interface{}
)
