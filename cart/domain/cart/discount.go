package cart

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"sort"
)

type (
	// AppliedDiscount value object - generic reference for a discount
	AppliedDiscount struct {
		CampaignCode  string       // unique code of the underlying campaign or rule e.g. "summer-campaign-2018"
		CouponCode    string       // code of discount e.g. provided by user "summer2018"
		Label         string       // readable name of discount "Super Summer Sale 2018"
		Applied       domain.Price // how much of the discount has been subtracted from cart price, IMPORTANT: always negative
		Type          string       // to distinguish between discounts
		IsItemRelated bool         // flag indicating if the discount is applied due to item in cart
		SortOrder     int          // indicates in which order discount have been applied, low value has been applied before high value
	}

	// WithDiscount interface for a cart that is able to handle discounts
	WithDiscount interface {
		HasAppliedDiscounts() (bool, error)
		MergeDiscounts() (AppliedDiscounts, error)
	}

	// AppliedDiscounts represents multiple discounts that are subtracted from total price of cart
	AppliedDiscounts []AppliedDiscount
)

var (
	// interface assertions
	_ WithDiscount = new(Item)
	_ WithDiscount = new(Delivery)
	_ WithDiscount = new(Cart)
)

// MergeDiscounts sums up discounts of cart based on its deliveries
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (c *Cart) MergeDiscounts() (AppliedDiscounts, error) {
	// guard if no items in delivery, no need to iterate
	if len(c.Deliveries) <= 0 {
		return make([]AppliedDiscount, 0), nil
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
	discounts, err := c.MergeDiscounts()
	if err != nil {
		return false, err
	}
	return len(discounts) > 0, nil
}

// MergeDiscounts sums up discounts of a delivery based on its single item discounts
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (d *Delivery) MergeDiscounts() (AppliedDiscounts, error) {
	// guard if no items in delivery, no need to iterate
	if len(d.Cartitems) <= 0 {
		return make([]AppliedDiscount, 0), nil
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
	discounts, err := d.MergeDiscounts()
	if err != nil {
		return false, err
	}
	return len(discounts) > 0, nil
}

// MergeDiscounts parses discounts of a single item
// All discounts with the same type and title are aggregated and returned as one with a summed price
func (i *Item) MergeDiscounts() (AppliedDiscounts, error) {
	sort.SliceStable(i.AppliedDiscounts, func(x, y int) bool {
		return i.AppliedDiscounts[x].SortOrder < i.AppliedDiscounts[y].SortOrder
	})

	return i.AppliedDiscounts, nil
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the cart
func (i *Item) HasAppliedDiscounts() (bool, error) {
	discounts, err := i.MergeDiscounts()
	if err != nil {
		return false, err
	}
	return len(discounts) > 0, nil
}

// private helper functions

// mapToSlice transform map of discounts to flat slice
func mapToSlice(collection map[string]*AppliedDiscount) []AppliedDiscount {
	result := make([]AppliedDiscount, 0, len(collection))
	for _, val := range collection {
		result = append(result, *val)
	}

	sort.SliceStable(result, func(x, y int) bool {
		return result[x].SortOrder < result[y].SortOrder
	})

	return result
}

// mapDiscounts map type + title of discount to corresponding discount
func mapDiscounts(result map[string]*AppliedDiscount, discountable WithDiscount) (map[string]*AppliedDiscount, error) {
	discounts, err := discountable.MergeDiscounts()
	// in case discounts cannot be collected, stop execution
	if err != nil {
		return nil, err
	}
	for _, discount := range discounts {
		// discount has been collected before, increase amount
		if collected, ok := result[discount.CampaignCode]; ok {
			update, err := collected.Applied.Add(discount.Applied)
			if err != nil {
				return nil, err
			}
			collected.Applied = update
			continue
		}
		// discount is new, add to collection
		result[discount.CampaignCode] = &AppliedDiscount{
			CampaignCode:  discount.CampaignCode,
			CouponCode:    discount.CouponCode,
			Label:         discount.Label,
			Applied:       discount.Applied,
			Type:          discount.Type,
			IsItemRelated: discount.IsItemRelated,
			SortOrder:     discount.SortOrder,
		}
	}
	return result, nil
}

// ByCampaignCode filter AppliedDiscounts based on provided campaign code
func (discounts AppliedDiscounts) ByCampaignCode(campaignCode string) AppliedDiscounts {
	f := func(discount AppliedDiscount) bool {
		return discount.CampaignCode == campaignCode
	}
	return discounts.filter(f)

}

// ByType filter AppliedDiscounts based on type
func (discounts AppliedDiscounts) ByType(filterType string) AppliedDiscounts {
	f := func(discount AppliedDiscount) bool {
		return discount.Type == filterType
	}
	return discounts.filter(f)
}

func (discounts AppliedDiscounts) filter(filterFunc func(AppliedDiscount) bool) AppliedDiscounts {
	result := make(AppliedDiscounts, 0)
	for _, discount := range discounts {
		if filterFunc(discount) {
			result = append(result, discount)
		}
	}
	return result
}
