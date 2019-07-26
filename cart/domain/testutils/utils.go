package testutils

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// ByCode implements sort.Interface for []AppliedDiscount based on code
	ByCode cart.AppliedDiscounts
)

// implementations for sort interface

func (a ByCode) Len() int {
	return len(a)
}

func (a ByCode) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCode) Less(i, j int) bool {
	return a[i].CampaignCode < a[j].CampaignCode
}

// BuildItemWithDiscounts helper for item building
func BuildItemWithDiscounts(t *testing.T) *cart.Item {
	t.Helper()
	builder := cart.ItemBuilder{}
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-1",
		Label:        "title-1",
		Type:         "type-1",
		Applied:      domain.NewFromFloat(-10.0, "$"),
		SortOrder:    3,
	})
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-2",
		Label:        "title-2",
		Type:         "type-1",
		Applied:      domain.NewFromFloat(-15.0, "$"),
		SortOrder:    2,
	})
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-3",
		Label:        "title-1",
		Type:         "type-2",
		Applied:      domain.NewFromFloat(-5.0, "$"),
		SortOrder:    4,
	})
	builder.SetID("id-1")
	item, err := builder.Build()
	if err != nil {
		t.Fatalf("Could not build item %s", err.Error())
	}
	return item
}

// BuildItemWithAlternativeDiscounts helper for item building with different discounts
func BuildItemWithAlternativeDiscounts(t *testing.T) *cart.Item {
	t.Helper()
	builder := cart.ItemBuilder{}
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-4",
		Label:        "title-4",
		Type:         "type-1",
		Applied:      domain.NewFromFloat(-10.0, "$"),
		SortOrder:    5,
	})
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-5",
		Label:        "title-5",
		Type:         "type-1",
		Applied:      domain.NewFromFloat(-15.0, "$"),
		SortOrder:    0,
	})
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-6",
		Label:        "title-6",
		Type:         "type-2",
		Applied:      domain.NewFromFloat(-5.0, "$"),
		SortOrder:    1,
	})
	builder.SetID("id-2")
	item, err := builder.Build()
	if err != nil {
		t.Fatalf("Could not build item %s", err.Error())
	}
	return item
}

// BuildDeliveryWithDiscounts helper for delivery building
// Adds an item with discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildDeliveryWithDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	builder := cart.DeliveryBuilder{}
	builder.SetDeliveryCode("code")
	builder.AddItem(*BuildItemWithDiscounts(t))
	builder.AddItem(*BuildItemWithDiscounts(t))
	// add items with discounts
	delivery, err := builder.Build()
	if err != nil {
		t.Fatalf("Could not build delivery %s", err.Error())
	}
	return delivery
}

// BuildAlternativeDeliveryWithAlternativeDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildAlternativeDeliveryWithAlternativeDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	builder := cart.DeliveryBuilder{}
	builder.SetDeliveryCode("code-2")
	builder.AddItem(*BuildItemWithAlternativeDiscounts(t))
	builder.AddItem(*BuildItemWithAlternativeDiscounts(t))
	// add items with discounts
	delivery, err := builder.Build()
	if err != nil {
		t.Fatalf("Could not build delivery %s", err.Error())
	}
	return delivery
}

// BuildDeliveryWithDifferentDiscounts helper for delivery building
// Adds an item with alternative discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func BuildDeliveryWithDifferentDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	builder := cart.DeliveryBuilder{}
	builder.SetDeliveryCode("code-1")
	builder.AddItem(*BuildItemWithDiscounts(t))
	builder.AddItem(*BuildItemWithAlternativeDiscounts(t))
	// add items with discounts
	delivery, err := builder.Build()
	if err != nil {
		t.Fatalf("Could not build delivery %s", err.Error())
	}
	return delivery
}

// BuildDeliveryWithoutDiscounts helper for delivery building
func BuildDeliveryWithoutDiscounts(t *testing.T) *cart.Delivery {
	t.Helper()
	builder := cart.DeliveryBuilder{}
	builder.AddItem(cart.Item{})
	builder.AddItem(cart.Item{})
	builder.SetDeliveryCode("code")
	delivery, err := builder.Build()
	if err != nil {
		t.Fatalf("Could not build delivery %s", err.Error())
	}
	return delivery
}
