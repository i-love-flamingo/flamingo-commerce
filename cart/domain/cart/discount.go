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
		CollectDiscounts() []*AppliedDiscount
	}

	// ByCode implements sort.Interface for []AppliedDiscount based on code
	ByCode []*AppliedDiscount
)

// CollectDiscounts sums up discounts of cart based on its deliveries
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (c *Cart) CollectDiscounts() []*AppliedDiscount {
	return nil
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (c *Cart) HasAppliedDiscounts() bool {
	return len(c.CollectDiscounts()) > 0
}

// CollectDiscounts sums up discounts of a delivery based on its single item discounts
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (d *Delivery) CollectDiscounts() []*AppliedDiscount {
	result := make([]*AppliedDiscount, 0)
	// guard if no items in delivery, no need to iterate
	if len(d.Cartitems) <= 0 {
		return result
	}
	return result
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the delivery
func (d *Delivery) HasAppliedDiscounts() bool {
	return len(d.CollectDiscounts()) > 0
}

// CollectDiscounts parses discounts of a single item
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (i *Item) CollectDiscounts() []*AppliedDiscount {
	// guard if no items in delivery, no need to iterate
	if len(i.AppliedDiscounts) <= 0 {
		return make([]*AppliedDiscount, 0)
	}
	// parse item discounts to applied discounts
	result := make([]*AppliedDiscount, len(i.AppliedDiscounts))
	for index, val := range i.AppliedDiscounts {
		result[index] = &AppliedDiscount{
			Code:    val.Code,
			Title:   val.Title,
			Applied: val.Amount,
			Type:    val.Type,
		}
	}
	return result
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (i *Item) HasAppliedDiscounts() bool {
	return len(i.CollectDiscounts()) > 0
}

// implementations for sort interface

func (a ByCode) Len() int {
	return len(a)
}

func (a ByCode) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCode) Less(i, j int) bool {
	return a[i].Code < a[j].Code
}
