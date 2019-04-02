package application_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type MockRestrictor struct {
	Qty uint
}

func (r *MockRestrictor) Restrict(ctx context.Context, product domain.BasicProduct, cart *cart.Cart) uint {
	return r.Qty
}

func validateRestriction(wantedRestriction uint, wantedError error) func(*testing.T, uint, error) {
	return func(t *testing.T, i uint, e error) {
		t.Helper()
		if i != wantedRestriction {
			t.Errorf("expected restriction %d, got %d", wantedRestriction, i)
		}
		if wantedError != e {
			t.Errorf("expected errror type %T, got %T", wantedError, e)
		}
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
		name      string
		fields    fields
		args      args
		validator func(*testing.T, uint, error)
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
			validator: validateRestriction(0, &application.ErrNoRestriction{}),
		},
		{
			name: "no restriction",
			fields: fields{
				qtyRestrictors: []cart.MaxQuantityRestrictor{&MockRestrictor{Qty: ^uint(0)}},
			},
			args: args{
				ctx:     context.Background(),
				product: nil,
				cart:    nil,
			},
			validator: validateRestriction(0, &application.ErrNoRestriction{}),
		},
		{
			name: "restrict to 5",
			fields: fields{
				qtyRestrictors: []cart.MaxQuantityRestrictor{&MockRestrictor{Qty: 5}},
			},
			args:      args{},
			validator: validateRestriction(5, nil),
		},
		{
			name: "multiple restrictors to 17",
			fields: fields{
				qtyRestrictors: []cart.MaxQuantityRestrictor{
					&MockRestrictor{Qty: 19},
					&MockRestrictor{Qty: 21},
					&MockRestrictor{Qty: 17},
					&MockRestrictor{Qty: 500},
					&MockRestrictor{Qty: ^uint(0)},
				},
			},
			args:      args{},
			validator: validateRestriction(17, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &application.RestrictionService{}
			rs.Inject(tt.fields.qtyRestrictors)
			got, err := rs.RestrictQty(tt.args.ctx, tt.args.product, tt.args.cart)
			tt.validator(t, got, err)
		})
	}
}
