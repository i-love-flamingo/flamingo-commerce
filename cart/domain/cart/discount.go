package cart

import (
	"sort"

	"flamingo.me/flamingo-commerce/v3/price/domain"
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
	_ WithDiscount = new(ShippingItem)
)

// MergeDiscounts sums up discounts of cart based on its deliveries
// All discounts with the same campaign code are aggregated and returned as one with a summed price
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
	return hasAppliedDiscounts(c)
}

// MergeDiscounts sums up discounts of a delivery based on its single item discounts
// All discounts with the same campaign code are aggregated and returned as one with a summed price
func (d *Delivery) MergeDiscounts() (AppliedDiscounts, error) {
	// guard if no items in delivery, no need to iterate
	if len(d.Cartitems) <= 0 {
		return make([]AppliedDiscount, 0), nil
	}
	var err error
	collection := make(map[string]*AppliedDiscount)
	// collect different discounts on item level
	for _, item := range d.Cartitems {
		collection, err = mapDiscounts(collection, &item)
		if err != nil {
			return nil, err
		}
	}
	// collect different discounts for shipping cost
	collection, err = mapDiscounts(collection, &d.ShippingItem)
	if err != nil {
		return nil, err
	}
	// transform map to flat slice
	return mapToSlice(collection), nil
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the delivery
func (d *Delivery) HasAppliedDiscounts() (bool, error) {
	return hasAppliedDiscounts(d)
}

// MergeDiscounts parses discounts of a single item
// All discounts with the same campaign code are aggregated and returned as one with a summed price
func (i *Item) MergeDiscounts() (AppliedDiscounts, error) {
	return mergeDiscountsOnItemLevel(i.AppliedDiscounts)
}

// HasAppliedDiscounts check whether there are any discounts currently applied to the item
func (i *Item) HasAppliedDiscounts() (bool, error) {
	return hasAppliedDiscounts(i)
}

// TotalWithDiscountInclTax returns the final shipping price to pay
// Deprecated use public field PriceGrossWithDiscounts
func (s *ShippingItem) TotalWithDiscountInclTax() domain.Price {
	return s.PriceGrossWithDiscounts
}

// MergeDiscounts parses discounts of a shipping item
// All discounts with the same campaign code are aggregated and returned as one with a summed price
func (s *ShippingItem) MergeDiscounts() (AppliedDiscounts, error) {
	return mergeDiscountsOnItemLevel(s.AppliedDiscounts)
}

// HasAppliedDiscounts checks whether there are any discounts currently applied to the shipping item
func (s *ShippingItem) HasAppliedDiscounts() (bool, error) {
	return hasAppliedDiscounts(s)
}

// private helper functions

// mergeDiscountOnItemLevel merges the discounts based on the campaign code
func mergeDiscountsOnItemLevel(discounts AppliedDiscounts) (AppliedDiscounts, error) {
	sort.SliceStable(discounts, func(x, y int) bool {
		return discounts[x].SortOrder < discounts[y].SortOrder
	})

	return discounts, nil
}

// hasAppliedDiscounts returns whether the discountable has discounts applied
func hasAppliedDiscounts(discountable WithDiscount) (bool, error) {
	discounts, err := discountable.MergeDiscounts()
	if err != nil {
		return false, err
	}
	return len(discounts) > 0, nil
}

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

// Sum returns the sum of the applied values of the AppliedDiscounts
func (discounts AppliedDiscounts) Sum() (domain.Price, error) {
	result := domain.Price{}
	var err error
	for _, val := range discounts {
		result, err = result.Add(val.Applied)
		if err != nil {
			return domain.NewZero(""), err
		}
	}
	return result, nil
}

// Items getter for graphql integration
func (discounts AppliedDiscounts) Items() []AppliedDiscount {
	return []AppliedDiscount(discounts)
}
