package domain

import "time"

type (
	// Order
	Order struct {
		ID           int
		CreationTime time.Time
		UpdateTime   string
		OrderItems   []OrderItem
		Status       string
		Total        float64
		CurrencyCode string
	}
	// OrderItem
	OrderItem struct {
		Sku          string
		Name         string
		Price        float64
		PriceInclTax float64
	}
)
