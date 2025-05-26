package testutils

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

// BuildItemWithDiscounts helper for item building
var BuildItemWithDiscounts = &cart.Item{ID: "id-1",
	AppliedDiscounts: cart.AppliedDiscounts{
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
		{
			CampaignCode: "code-7",
			Label:        "title-7",
			Type:         "type-7",
			Applied:      domain.NewFromFloat(-10.0, "$"),
			SortOrder:    6,
			CustomAttributes: map[string]interface{}{
				"attr1": 3,
				"attr2": 1,
			},
		},
	},
}

// BuildItemWithAlternativeDiscounts helper for item building with different discounts
var BuildItemWithAlternativeDiscounts = &cart.Item{
	ID: "id-2",
	AppliedDiscounts: cart.AppliedDiscounts{
		{
			CampaignCode: "code-4",
			Label:        "title-4",
			Type:         "type-1",
			Applied:      domain.NewFromFloat(-10.0, "$"),
			SortOrder:    5,
		},
		{
			CampaignCode: "code-5",
			Label:        "title-5",
			Type:         "type-1",
			Applied:      domain.NewFromFloat(-15.0, "$"),
			SortOrder:    0,
		},
		{
			CampaignCode: "code-6",
			Label:        "title-6",
			Type:         "type-2",
			Applied:      domain.NewFromFloat(-5.0, "$"),
			SortOrder:    1,
		},
	},
} // todo: add discount total

// BuildItemWithDuplicateDiscounts helper for item building with duplicate discounts
var BuildItemWithDuplicateDiscounts = &cart.Item{
	ID: "id-1",
	AppliedDiscounts: cart.AppliedDiscounts{
		{
			CampaignCode: "code-1",
			Label:        "title-1",
			Type:         "type-1",
			Applied:      domain.NewFromFloat(-10.0, "$"),
			SortOrder:    0,
		},
		{
			CampaignCode: "code-1",
			Label:        "title-1",
			Type:         "type-1",
			Applied:      domain.NewFromFloat(-10.0, "$"),
			SortOrder:    0,
		},
	},
} // todo: add discount total

// BuildShippingItemWithDiscounts helper for shipping item building
var BuildShippingItemWithDiscounts = &cart.ShippingItem{
	Title:      "",
	PriceNet:   domain.NewFromFloat(20.0, "$"),
	TaxAmount:  domain.NewFromFloat(2.0, "$"),
	PriceGross: domain.NewFromFloat(22.0, "$"),
	AppliedDiscounts: cart.AppliedDiscounts{
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
			Applied:      domain.NewFromFloat(-5.0, "$"),
			SortOrder:    2,
		},
	},
}

// BuildShippingItemWithDuplicateDiscounts helper for shipping item building with duplicate discounts
var BuildShippingItemWithDuplicateDiscounts = &cart.ShippingItem{
	Title:      "",
	PriceNet:   domain.NewFromFloat(40.0, "$"),
	TaxAmount:  domain.NewFromFloat(2.0, "$"),
	PriceGross: domain.NewFromFloat(42.0, "$"),
	AppliedDiscounts: cart.AppliedDiscounts{
		{
			CampaignCode: "code-1",
			Label:        "title-1",
			Type:         "type-1",
			Applied:      domain.NewFromFloat(-15.0, "$"),
			SortOrder:    0,
		},
		{
			CampaignCode: "code-1",
			Label:        "title-1",
			Type:         "type-1",
			Applied:      domain.NewFromFloat(-15.0, "$"),
			SortOrder:    0,
		},
	},
}

// BuildDeliveryWithDiscounts helper for delivery building
// Adds an item with discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
var BuildDeliveryWithDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code"},
	Cartitems:    []cart.Item{*BuildItemWithDiscounts, *BuildItemWithDiscounts},
}

// BuildAlternativeDeliveryWithAlternativeDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
var BuildAlternativeDeliveryWithAlternativeDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code-2"},
	Cartitems:    []cart.Item{*BuildItemWithAlternativeDiscounts, *BuildItemWithAlternativeDiscounts},
}

// BuildDeliveryWithDifferentDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
var BuildDeliveryWithDifferentDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
	Cartitems:    []cart.Item{*BuildItemWithDiscounts, *BuildItemWithAlternativeDiscounts},
}

// BuildDeliveryWithDuplicateDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
var BuildDeliveryWithDuplicateDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
	Cartitems:    []cart.Item{*BuildItemWithDuplicateDiscounts},
}

// BuildDeliveryWithoutDiscounts helper for delivery building
var BuildDeliveryWithoutDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code"},
	Cartitems:    []cart.Item{{}, {}},
}

// BuildDeliveryWithoutDiscountsAndShippingDiscounts helper for delivery building
var BuildDeliveryWithoutDiscountsAndShippingDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code"},
	Cartitems:    []cart.Item{{}, {}},
	ShippingItem: *BuildShippingItemWithDiscounts,
}

// BuildDeliveryWithDifferentDiscountsAndShippingDiscounts helper for delivery building
// Adds an item with alternative discount twice
// Adds a shipping item with discounts
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
var BuildDeliveryWithDifferentDiscountsAndShippingDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
	Cartitems:    []cart.Item{*BuildItemWithDiscounts, *BuildItemWithAlternativeDiscounts},
	ShippingItem: *BuildShippingItemWithDiscounts,
}

// BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts helper for delivery building
// Adds an item with duplicate discounts
// Adds a shipping item with discounts
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
var BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts = &cart.Delivery{
	DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
	Cartitems:    []cart.Item{*BuildItemWithDuplicateDiscounts},
	ShippingItem: *BuildShippingItemWithDiscounts,
}
