package testutils

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

// BuildItemWithDiscounts helper for item building
func BuildItemWithDiscounts(t *testing.T) *cart.Item {
	t.Helper()
	item := cart.Item{ID: "id-1",
		AppliedDiscounts: []cart.AppliedDiscount{
			{
				CampaignCode: "code-1",
				Label:        "title-1",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-10.0, "$"),
				SortOrder:    3,
			},
			{
				CampaignCode: "code-2",
				Label:        "title-2",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-15.0, "$"),
				SortOrder:    2,
			},
			{
				CampaignCode: "code-3",
				Label:        "title-1",
				Type:         "type-2",
				Applied:      domain.NewFromFloat(-5.0, "$"),
				SortOrder:    4,
			},
		},
	}

	// todo: add discount total

	return &item
}

// BuildItemWithAlternativeDiscounts helper for item building with different discounts
func BuildItemWithAlternativeDiscounts(t *testing.T) *cart.Item {
	t.Helper()
	item := cart.Item{
		ID: "id-2",
		AppliedDiscounts: []cart.AppliedDiscount{
			cart.AppliedDiscount{
				CampaignCode: "code-4",
				Label:        "title-4",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-10.0, "$"),
				SortOrder:    5,
			},
			cart.AppliedDiscount{
				CampaignCode: "code-5",
				Label:        "title-5",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-15.0, "$"),
				SortOrder:    0,
			},
			cart.AppliedDiscount{
				CampaignCode: "code-6",
				Label:        "title-6",
				Type:         "type-2",
				Applied:      domain.NewFromFloat(-5.0, "$"),
				SortOrder:    1,
			},
		},
	} // todo: add discount total

	return &item
}

// BuildItemWithDuplicateDiscounts helper for item building with duplicate discounts
func BuildItemWithDuplicateDiscounts(t *testing.T) *cart.Item {
	t.Helper()

	item := cart.Item{
		ID: "id-1",
		AppliedDiscounts: []cart.AppliedDiscount{
			cart.AppliedDiscount{
				CampaignCode: "code-1",
				Label:        "title-1",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-10.0, "$"),
				SortOrder:    0,
			},
			cart.AppliedDiscount{
				CampaignCode: "code-1",
				Label:        "title-1",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-10.0, "$"),
				SortOrder:    0,
			},
		},
	} // todo: add discount total

	return &item
}

// BuildShippingItemWithDiscounts helper for shipping item building
func BuildShippingItemWithDiscounts(t *testing.T) *cart.ShippingItem {
	t.Helper()
	return &cart.ShippingItem{
		Title:      "",
		PriceNet:   domain.NewFromFloat(20.0, "$"),
		TaxAmount:  domain.NewFromFloat(2.0, "$"),
		PriceGross: domain.NewFromFloat(22.0, "$"),
		AppliedDiscounts: cart.AppliedDiscounts{
			cart.AppliedDiscount{
				CampaignCode: "code-1",
				Label:        "title-1",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-10.0, "$"),
				SortOrder:    3,
			},
			cart.AppliedDiscount{
				CampaignCode: "code-2",
				Label:        "title-2",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-5.0, "$"),
				SortOrder:    2,
			},
		},
	}
}

// BuildShippingItemWithAlternativeDiscounts helper for shipping item building with different discounts
func BuildShippingItemWithAlternativeDiscounts(t *testing.T) *cart.ShippingItem {
	t.Helper()
	return &cart.ShippingItem{
		Title:      "",
		PriceNet:   domain.NewFromFloat(30.0, "$"),
		TaxAmount:  domain.NewFromFloat(2.0, "$"),
		PriceGross: domain.NewFromFloat(32.0, "$"),
		AppliedDiscounts: cart.AppliedDiscounts{
			cart.AppliedDiscount{
				CampaignCode: "code-3",
				Label:        "title-1",
				Type:         "type-2",
				Applied:      domain.NewFromFloat(-5.0, "$"),
				SortOrder:    4,
			},
			cart.AppliedDiscount{
				CampaignCode: "code-4",
				Label:        "title-4",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-20.0, "$"),
				SortOrder:    5,
			},
		},
	}
}

// BuildShippingItemWithDuplicateDiscounts helper for shipping item building with duplicate discounts
func BuildShippingItemWithDuplicateDiscounts(t *testing.T) *cart.ShippingItem {
	t.Helper()
	return &cart.ShippingItem{
		Title:      "",
		PriceNet:   domain.NewFromFloat(40.0, "$"),
		TaxAmount:  domain.NewFromFloat(2.0, "$"),
		PriceGross: domain.NewFromFloat(42.0, "$"),
		AppliedDiscounts: cart.AppliedDiscounts{
			cart.AppliedDiscount{
				CampaignCode: "code-1",
				Label:        "title-1",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-15.0, "$"),
				SortOrder:    0,
			},
			cart.AppliedDiscount{
				CampaignCode: "code-1",
				Label:        "title-1",
				Type:         "type-1",
				Applied:      domain.NewFromFloat(-15.0, "$"),
				SortOrder:    0,
			},
		},
	}
}

// BuildDeliveryWithDiscounts helper for delivery building
// Adds an item with discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildDeliveryWithDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code"},
		Cartitems:    []cart.Item{*BuildItemWithDiscounts(t), *BuildItemWithDiscounts(t)},
	}
	return delivery
}

// BuildAlternativeDeliveryWithAlternativeDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildAlternativeDeliveryWithAlternativeDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code-2"},
		Cartitems:    []cart.Item{*BuildItemWithAlternativeDiscounts(t), *BuildItemWithAlternativeDiscounts(t)},
	}
	return delivery
}

// BuildDeliveryWithDifferentDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildDeliveryWithDifferentDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
		Cartitems:    []cart.Item{*BuildItemWithDiscounts(t), *BuildItemWithAlternativeDiscounts(t)},
	}
	return delivery
}

// BuildDeliveryWithDuplicateDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildDeliveryWithDuplicateDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
		Cartitems:    []cart.Item{*BuildItemWithDuplicateDiscounts(t)},
	}
	return delivery
}

// BuildDeliveryWithoutDiscounts helper for delivery building
func BuildDeliveryWithoutDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code"},
		Cartitems:    []cart.Item{{}, {}},
	}
	return delivery
}

// BuildDeliveryWithoutDiscountsAndShippingDiscounts helper for delivery building
func BuildDeliveryWithoutDiscountsAndShippingDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()

	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code"},
		Cartitems:    []cart.Item{{}, {}},
		ShippingItem: *BuildShippingItemWithDiscounts(t),
	}

	return delivery
}

// BuildDeliveryWithDifferentDiscountsAndShippingDiscounts helper for delivery building
// Adds an item with alternative discount twice
// Adds a shipping item with discounts
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildDeliveryWithDifferentDiscountsAndShippingDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
		Cartitems:    []cart.Item{*BuildItemWithDiscounts(t), *BuildItemWithAlternativeDiscounts(t)},
		ShippingItem: *BuildShippingItemWithDiscounts(t),
	}

	return delivery
}

// BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts helper for delivery building
// Adds an item with duplicate discounts
// Adds a shipping item with discounts
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	delivery := &cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
		Cartitems:    []cart.Item{*BuildItemWithDuplicateDiscounts(t)},
		ShippingItem: *BuildShippingItemWithDiscounts(t),
	}

	return delivery
}
