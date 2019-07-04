package cart

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// AppliedDiscount value object - generic reference for a discount
	AppliedDiscount struct {
		Code    string       // unique code of discount
		Title   string       // readable name of discount
		Applied domain.Price // how much of the discount has been subtracted from cart price
		Type    string       // to distinguish between discounts
	}

	// WithDiscount interface for a cart that is able to handle discounts
	WithDiscount interface {
		HasAppliedDiscounts() bool
		AggregateDiscounts() []*AppliedDiscount
	}
)

// AggregateDiscounts sums up discounts of cart based on its deliveries
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (c *Cart) AggregateDiscounts() []*AppliedDiscount {
	return nil
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (c *Cart) HasAppliedDiscounts() bool {
	return len(c.AggregateDiscounts()) > 0
}

// AggregateDiscounts sums up discounts of a delivery based on its single item discounts
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (d *Delivery) AggregateDiscounts() []*AppliedDiscount {
	result := make([]*AppliedDiscount, 0)
	// guard if no items in delivery, no need to iterate
	if len(d.Cartitems) <= 0 {
		return result
	}
	return result
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the delivery
func (d *Delivery) HasAppliedDiscounts() bool {
	return len(d.AggregateDiscounts()) > 0
}

// AggregateDiscounts parses discounts of a single item
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (i *Item) AggregateDiscounts() []*AppliedDiscount {
	result := make([]*AppliedDiscount, 0)
	// guard if no items in delivery, no need to iterate
	if len(i.AppliedDiscounts) <= 0 {
		return result
	}
	// aggregate discounts by type + title
	collection := make(map[string]*AppliedDiscount)
	for _, discount := range i.AppliedDiscounts {
		key := discount.Type + discount.Title
		if collected, ok := collection[key]; ok {
			collected.Applied.Add(discount.Amount)
			break
		}
		collection[key] = &AppliedDiscount{
			Code:    discount.Code,
			Title:   discount.Title,
			Applied: discount.Amount,
			Type:    discount.Type,
		}
	}
	// restructure map to slice
	for _, val := range collection {
		result = append(result, val)
	}
	return result
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (i *Item) HasAppliedDiscounts() bool {
	return len(i.AggregateDiscounts()) > 0
}
