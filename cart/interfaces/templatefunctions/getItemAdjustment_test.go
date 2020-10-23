package templatefunctions

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
)

func TestGetQuantityAdjustmentDeletedItemsMessages_Func(t *testing.T) {
	type args struct {
		contextProvider func() context.Context
	}
	tests := []struct {
		name string
		args args
		want []QuantityAdjustment
	}{
		{
			name: "no adjustments returns empty slice",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: []QuantityAdjustment{},
		},
		{
			name: "adjustment that was not deleted returns empty slice",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{
						{
							WasDeleted: false,
						},
					})

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: []QuantityAdjustment{},
		},
		{
			name: "adjustment that was deleted returns respective adjustment",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{
						{
							OriginalItem: cart.Item{
								ID: "itemID",
							},
							DeliveryCode: "deliveryCode",
							WasDeleted:   true,
							RestrictionResult: &validation.RestrictionResult{
								RemainingDifference: -2,
								RestrictorName:      "restrictor",
							},
							NewQty: 0,
						},
					})

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: []QuantityAdjustment{
				{
					Item: cart.Item{
						ID: "itemID",
					},
					DeliveryCode: "deliveryCode",
					PrevQty:      2,
					CurrQty:      0,
					Reason:       "restrictor",
				},
			},
		},
		{
			name: "returns only deleted adjustment",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{
						{
							OriginalItem: cart.Item{
								ID: "itemID-A",
							},
							DeliveryCode: "deliveryCode-A",
							WasDeleted:   false,
							RestrictionResult: &validation.RestrictionResult{
								RemainingDifference: -2,
								RestrictorName:      "restrictor-A",
							},
							NewQty: 1,
						},
						{
							OriginalItem: cart.Item{
								ID: "itemID-B",
							},
							DeliveryCode: "deliveryCode-B",
							WasDeleted:   true,
							RestrictionResult: &validation.RestrictionResult{
								RemainingDifference: -4,
								RestrictorName:      "restrictor-B",
							},
							NewQty: 0,
						},
					})

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: []QuantityAdjustment{
				{
					Item: cart.Item{
						ID: "itemID-B",
					},
					DeliveryCode: "deliveryCode-B",
					PrevQty:      4,
					CurrQty:      0,
					Reason:       "restrictor-B",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdm := &GetQuantityAdjustmentDeletedItemsMessages{}
			getQuantityAdjustmentDeletedItemsMessages := gdm.Func(tt.args.contextProvider()).(func() []QuantityAdjustment)
			if got := getQuantityAdjustmentDeletedItemsMessages(); cmp.Diff(got, tt.want) != "" {
				t.Errorf("Func() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetQuantityAdjustmentUpdatedItemsMessage_Func(t *testing.T) {
	type args struct {
		contextProvider func() context.Context
		item            cart.Item
		deliveryCode    string
	}
	tests := []struct {
		name string
		args args
		want QuantityAdjustment
	}{
		{
			name: "no adjustments returns adjustment with prev qty = curr qty",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()

					return web.ContextWithSession(context.Background(), session)
				},
				item: cart.Item{
					ID:  "itemID",
					Qty: 3,
				},
				deliveryCode: "deliveryCode",
			},
			want: QuantityAdjustment{
				Item: cart.Item{
					ID:  "itemID",
					Qty: 3,
				},
				DeliveryCode: "deliveryCode",
				PrevQty:      3,
				CurrQty:      3,
			},
		},
		{
			name: "not found adjustment returns adjustment with prev qty = curr qty",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{})

					return web.ContextWithSession(context.Background(), session)
				},
				item: cart.Item{
					ID:  "itemID",
					Qty: 3,
				},
				deliveryCode: "deliveryCode",
			},
			want: QuantityAdjustment{
				Item: cart.Item{
					ID:  "itemID",
					Qty: 3,
				},
				DeliveryCode: "deliveryCode",
				PrevQty:      3,
				CurrQty:      3,
			},
		},
		{
			name: "found adjustment returns respective adjustment",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{
						{
							OriginalItem: cart.Item{
								ID:  "itemID",
								Qty: 1,
							},
							DeliveryCode: "deliveryCode",
							RestrictionResult: &validation.RestrictionResult{
								IsRestricted:        true,
								MaxAllowed:          1,
								RemainingDifference: -2,
								RestrictorName:      "restrictor",
							},
							NewQty: 1,
						},
					})
					return web.ContextWithSession(context.Background(), session)
				},
				item: cart.Item{
					ID:  "itemID",
					Qty: 3,
				},
				deliveryCode: "deliveryCode",
			},
			want: QuantityAdjustment{
				Item: cart.Item{
					ID:  "itemID",
					Qty: 1,
				},
				DeliveryCode: "deliveryCode",
				PrevQty:      3,
				CurrQty:      1,
				Reason:       "restrictor",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gum := &GetQuantityAdjustmentUpdatedItemsMessage{}
			getQuantityAdjustmentUpdatedItemsMessage := gum.Func(tt.args.contextProvider()).(func(cart.Item, string) QuantityAdjustment)
			if got := getQuantityAdjustmentUpdatedItemsMessage(tt.args.item, tt.args.deliveryCode); cmp.Diff(got, tt.want) != "" {
				t.Errorf("Func() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetQuantityAdjustmentCouponCodesRemoved_Func(t *testing.T) {
	type args struct {
		contextProvider func() context.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no adjustments returns false",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: false,
		},
		{
			name: "adjustment that has not removed coupon codes returns false",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{
						{
							HasRemovedCouponCodes: false,
						},
					})

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: false,
		},
		{
			name: "adjustment that has removed coupon codes returns true",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{
						{
							HasRemovedCouponCodes: true,
						},
					})

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: true,
		},
		{
			name: "at least one adjustment that has removed coupon codes returns true",
			args: args{
				contextProvider: func() context.Context {
					session := web.EmptySession()
					session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{
						{
							HasRemovedCouponCodes: false,
						},
						{
							HasRemovedCouponCodes: true,
						},
					})

					return web.ContextWithSession(context.Background(), session)
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcd := &GetQuantityAdjustmentCouponCodesRemoved{}
			getQuantityAdjustmentCouponCodesRemoved := gcd.Func(tt.args.contextProvider()).(func() bool)
			if got := getQuantityAdjustmentCouponCodesRemoved(); cmp.Diff(got, tt.want) != "" {
				t.Errorf("Func() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveQuantityAdjustmentMessages_Func(t *testing.T) {
	session := web.EmptySession()
	session.Store("cart.view.quantity.adjustments", application.QtyAdjustmentResults{{}})

	ctx := web.ContextWithSession(context.Background(), session)

	rm := &RemoveQuantityAdjustmentMessages{}
	removeQuantityAdjustmentMessages := rm.Func(ctx).(func() bool)

	_, found := session.Load("cart.view.quantity.adjustments")
	assert.Equal(t, true, found)

	removeQuantityAdjustmentMessages()

	_, found = session.Load("cart.view.quantity.adjustments")
	assert.Equal(t, false, found)
}
