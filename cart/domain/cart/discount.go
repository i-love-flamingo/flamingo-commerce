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
		HasAppliedDiscounts() (bool, error)
		CollectDiscounts() ([]*AppliedDiscount, error)
	}

	// ByCode implements sort.Interface for []AppliedDiscount based on code
	ByCode []*AppliedDiscount
)

// CollectDiscounts sums up discounts of cart based on its deliveries
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (c *Cart) CollectDiscounts() ([]*AppliedDiscount, error) {
	// guard if no items in delivery, no need to iterate
	if len(c.Deliveries) <= 0 {
		return make([]*AppliedDiscount, 0), nil
	}
	// collect different discounts on item level
	var err error
	collection := make(map[string]*AppliedDiscount)
	for _, delivery := range c.Deliveries {
		collection, err = mapDiscounts(collection, &delivery)
		if err != nil {
			return nil, err
		}
	}
	// transform map to flat slice
	return mapToSlice(collection), nil
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (c *Cart) HasAppliedDiscounts() (bool, error) {
	discounts, err := c.CollectDiscounts()
	if err != nil {
		return false, err
	}
	return len(discounts) > 0, nil
}

// CollectDiscounts sums up discounts of a delivery based on its single item discounts
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (d *Delivery) CollectDiscounts() ([]*AppliedDiscount, error) {
	// guard if no items in delivery, no need to iterate
	if len(d.Cartitems) <= 0 {
		return make([]*AppliedDiscount, 0), nil
	}
	// collect different discounts on item level
	var err error
	collection := make(map[string]*AppliedDiscount)
	for _, item := range d.Cartitems {
		collection, err = mapDiscounts(collection, &item)
		if err != nil {
			return nil, err
		}
	}
	// transform map to flat slice
	return mapToSlice(collection), nil
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the delivery
func (d *Delivery) HasAppliedDiscounts() (bool, error) {
	discounts, err := d.CollectDiscounts()
	if err != nil {
		return false, err
	}
	return len(discounts) > 0, nil
}

// CollectDiscounts parses discounts of a single item
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (i *Item) CollectDiscounts() ([]*AppliedDiscount, error) {
	// guard if no items in delivery, no need to iterate
	if len(i.AppliedDiscounts) <= 0 {
		return make([]*AppliedDiscount, 0), nil
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
	return result, nil
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (i *Item) HasAppliedDiscounts() (bool, error) {
	discounts, err := i.CollectDiscounts()
	if err != nil {
		return false, err
	}
	return len(discounts) > 0, nil
}

// private helper functions

// mapToSlice transform map of discounts to flat slice
func mapToSlice(collection map[string]*AppliedDiscount) []*AppliedDiscount {
	result := make([]*AppliedDiscount, 0, len(collection))
	for _, val := range collection {
		result = append(result, val)
	}
	return result
}

// mapDiscounts map type + title of discount to corresponding discount
func mapDiscounts(result map[string]*AppliedDiscount, discountable WithDiscount) (map[string]*AppliedDiscount, error) {
	discounts, err := discountable.CollectDiscounts()
	// in case discounts cannot be collected, stop execution
	if err != nil {
		return nil, err
	}
	for _, discount := range discounts {
		key := discount.Type + discount.Title
		// discount has been collected before, increase amount
		if collected, ok := result[key]; ok {
			update, err := collected.Applied.Add(discount.Applied)
			if err != nil {
				return nil, err
			}
			collected.Applied = update
			continue
		}
		// discount is new, add to collection
		result[key] = &AppliedDiscount{
			Code:    discount.Code,
			Title:   discount.Title,
			Applied: discount.Applied,
			Type:    discount.Type,
		}
	}
	return result, nil
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
