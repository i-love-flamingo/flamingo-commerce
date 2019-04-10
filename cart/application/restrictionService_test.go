package application_test

import (
	"context"
	"math"
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type MockRestrictor struct {
	IsRestricted  bool
	MaxQty        int
	DifferenceQty int
}

func (r *MockRestrictor) Restrict(ctx context.Context, product domain.BasicProduct, currentCart *cart.Cart) *cart.RestrictionResult {
	return &cart.RestrictionResult{
		IsRestricted:        r.IsRestricted,
		MaxAllowed:          r.MaxQty,
		RemainingDifference: r.DifferenceQty,
	}
}

func TestRestrictionService_RestrictQty(t *testing.T) {
	type fields struct {
		qtyRestrictors []cart.MaxQuantityRestrictor
	}
	type args struct {
		ctx     context.Context
		product domain.BasicProduct
		cart    *cart.Cart
	}
	tests := []struct {
		name                      string
		fields                    fields
		args                      args
		expectedRestrictionResult *cart.RestrictionResult
	}{
		{
			name: "no restrictors",
			fields: fields{
				qtyRestrictors: nil,
			},
			args: args{
				ctx:     context.Background(),
				product: nil,
				cart:    nil,
			},
			expectedRestrictionResult: &cart.RestrictionResult{
				IsRestricted:        false,
				MaxAllowed:          math.MaxInt32,
				RemainingDifference: math.MaxInt32,
			},
		},
		{
			name: "no restricting restrictors",
			fields: fields{
				qtyRestrictors: []cart.MaxQuantityRestrictor{&MockRestrictor{IsRestricted: false}, &MockRestrictor{IsRestricted: false}},
			},
			args: args{
				ctx:     context.Background(),
				product: nil,
				cart:    nil,
			},
			expectedRestrictionResult: &cart.RestrictionResult{
				IsRestricted:        false,
				MaxAllowed:          math.MaxInt32,
				RemainingDifference: math.MaxInt32,
			},
		},
		{
			name: "restrict to 5",
			fields: fields{
				qtyRestrictors: []cart.MaxQuantityRestrictor{&MockRestrictor{IsRestricted: true, MaxQty: 5, DifferenceQty: 5}},
			},
			args: args{},
			expectedRestrictionResult: &cart.RestrictionResult{
				IsRestricted:        true,
				MaxAllowed:          5,
				RemainingDifference: 5,
			},
		},
		{
			name: "multiple restrictors to 17 / -7",
			fields: fields{
				qtyRestrictors: []cart.MaxQuantityRestrictor{
					&MockRestrictor{IsRestricted: true, MaxQty: 19, DifferenceQty: 19},
					&MockRestrictor{IsRestricted: true, MaxQty: 21, DifferenceQty: 5},
					&MockRestrictor{IsRestricted: false, MaxQty: -42, DifferenceQty: -42},
					&MockRestrictor{IsRestricted: true, MaxQty: 17, DifferenceQty: 6},
					&MockRestrictor{IsRestricted: true, MaxQty: 500, DifferenceQty: -7},
					&MockRestrictor{IsRestricted: true, MaxQty: math.MaxInt32, DifferenceQty: math.MaxInt32},
				},
			},
			args: args{},
			expectedRestrictionResult: &cart.RestrictionResult{
				IsRestricted:        true,
				MaxAllowed:          17,
				RemainingDifference: -7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &application.RestrictionService{}
			rs.Inject(tt.fields.qtyRestrictors)
			got := rs.RestrictQty(tt.args.ctx, tt.args.product, tt.args.cart)
			if !reflect.DeepEqual(got, tt.expectedRestrictionResult) {
				t.Errorf("RestrictionService.RestrictQty() got = %v, expected = %v", got, tt.expectedRestrictionResult)
			}
		})
	}
}
