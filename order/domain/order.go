package domain

type (
	// Order
	Order struct {
		ID           int
		CreationTime string
		UpdateTime   string
		OrderItems   []OrderItem
		Status       string
		Total        float64
	}
	// OrderItem
	OrderItem struct {
		Sku          string
		Name         string
		Price        float64
		PriceInclTax float64
	}
)
