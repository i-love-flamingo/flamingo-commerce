package infrastructure

import (
	"context"
	"reflect"
	"testing"

	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/go-test/deep"
)

func TestInMemoryBehaviour_CleanCart(t *testing.T) {
	tests := []struct {
		name       string
		want       *domaincart.Cart
		wantDefers domaincart.DeferEvents
		wantErr    bool
	}{
		{
			name: "clean cart",
			want: &domaincart.Cart{
				ID:         "17",
				Deliveries: []domaincart.Delivery{},
			},
			wantDefers: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cob := &InMemoryBehaviour{}
			cob.Inject(
				&InMemoryCartStorage{},
				nil,
				flamingo.NullLogger{},
				func() *domaincart.ItemBuilder {
					return &domaincart.ItemBuilder{}
				},
				func() *domaincart.DeliveryBuilder {
					return &domaincart.DeliveryBuilder{}
				},
				func() *domaincart.Builder {
					return &domaincart.Builder{}
				},
				nil,
				nil,
			)
			cart := &domaincart.Cart{
				ID: "17",
				Deliveries: []domaincart.Delivery{
					{
						DeliveryInfo: domaincart.DeliveryInfo{
							Code: "dev-1",
						},
						Cartitems: nil,
					},
				},
			}

			if err := cob.cartStorage.StoreCart(cart); err != nil {
				t.Fatalf("cart could not be initialized")
			}

			got, gotDefers, err := cob.CleanCart(context.Background(), cart)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryCartOrderBehaviour.CleanCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("InMemoryCartOrderBehaviour.CleanCart() got!=want, diff: %#v", diff)
			}
			if diff := deep.Equal(gotDefers, tt.wantDefers); diff != nil {
				t.Errorf("InMemoryCartOrderBehaviour.CleanCart() gotDefers!=wantDefers, diff: %#v", diff)
			}
		})
	}
}

func TestInMemoryBehaviour_CleanDelivery(t *testing.T) {

	type args struct {
		cart         *domaincart.Cart
		deliveryCode string
	}
	tests := []struct {
		name       string
		args       args
		want       *domaincart.Cart
		wantDefers domaincart.DeferEvents
		wantErr    bool
	}{
		{
			name: "clean dev-1",
			args: args{
				cart: &domaincart.Cart{
					ID: "17",
					Deliveries: []domaincart.Delivery{
						{
							DeliveryInfo: domaincart.DeliveryInfo{
								Code: "dev-1",
							},
							Cartitems: nil,
						},
						{
							DeliveryInfo: domaincart.DeliveryInfo{
								Code: "dev-2",
							},
							Cartitems: nil,
						},
					},
				},
				deliveryCode: "dev-1",
			},
			want: &domaincart.Cart{
				ID: "17",
				Deliveries: []domaincart.Delivery{
					{
						DeliveryInfo: domaincart.DeliveryInfo{
							Code: "dev-2",
						},
						Cartitems: nil,
					},
				},
			},
			wantDefers: nil,
			wantErr:    false,
		},
		{
			name: "delivery not found",
			args: args{
				cart: &domaincart.Cart{
					ID: "17",
					Deliveries: []domaincart.Delivery{
						{
							DeliveryInfo: domaincart.DeliveryInfo{
								Code: "dev-1",
							},
							Cartitems: nil,
						},
						{
							DeliveryInfo: domaincart.DeliveryInfo{
								Code: "dev-2",
							},
							Cartitems: nil,
						},
					},
				},
				deliveryCode: "dev-3",
			},
			want:       nil,
			wantDefers: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cob := &InMemoryBehaviour{}
			cob.Inject(
				&InMemoryCartStorage{},
				nil,
				flamingo.NullLogger{},
				func() *domaincart.ItemBuilder {
					return &domaincart.ItemBuilder{}
				},
				func() *domaincart.DeliveryBuilder {
					return &domaincart.DeliveryBuilder{}
				},
				func() *domaincart.Builder {
					return &domaincart.Builder{}
				},
				nil,
				nil,
			)
			if err := cob.cartStorage.StoreCart(tt.args.cart); err != nil {
				t.Fatalf("cart could not be initialized")
			}

			got, gotDefers, err := cob.CleanDelivery(context.Background(), tt.args.cart, tt.args.deliveryCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryCartOrderBehaviour.CleanDelivery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("InMemoryCartOrderBehaviour.CleanDelivery() got!=want, diff: %#v", diff)
			}
			if diff := deep.Equal(gotDefers, tt.wantDefers); diff != nil {
				t.Errorf("InMemoryCartOrderBehaviour.CleanCart() gotDefers!=wantDefers, diff: %#v", diff)
			}
		})
	}
}

func TestInMemoryBehaviour_RemoveVoucher(t *testing.T) {
	type args struct {
		ctx                context.Context
		cart               *domaincart.Cart
		couponCodeToRemove string
	}
	tests := []struct {
		name string
		args args
		want *domaincart.Cart
	}{
		{
			name: "Remove voucher from cart with vouchers",
			args: args{
				ctx: nil,
				cart: &domaincart.Cart{
					AppliedCouponCodes: []domaincart.CouponCode{
						{Code: "OFF20"},
						{Code: "dummy-voucher-20"},
						{Code: "SALE"},
					},
				},
				couponCodeToRemove: "dummy-voucher-20",
			},
			want: &domaincart.Cart{
				AppliedCouponCodes: []domaincart.CouponCode{
					{Code: "OFF20"},
					{Code: "SALE"},
				},
			},
		},
		{
			name: "Remove voucher from cart without vouchers",
			args: args{
				ctx:                nil,
				cart:               &domaincart.Cart{},
				couponCodeToRemove: "dummy-voucher-20",
			},
			want: &domaincart.Cart{},
		},
		{
			name: "Remove voucher from cart that does not exist",
			args: args{
				ctx: nil,
				cart: &domaincart.Cart{
					AppliedCouponCodes: []domaincart.CouponCode{
						{Code: "OFF20"},
						{Code: "dummy-voucher-20"},
						{Code: "SALE"},
					},
				},
				couponCodeToRemove: "non-existing-voucher",
			},
			want: &domaincart.Cart{
				AppliedCouponCodes: []domaincart.CouponCode{
					{Code: "OFF20"},
					{Code: "dummy-voucher-20"},
					{Code: "SALE"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cob := &InMemoryBehaviour{}
			cob.Inject(
				&InMemoryCartStorage{},
				nil,
				flamingo.NullLogger{},
				func() *domaincart.ItemBuilder {
					return &domaincart.ItemBuilder{}
				},
				func() *domaincart.DeliveryBuilder {
					return &domaincart.DeliveryBuilder{}
				},
				func() *domaincart.Builder {
					return &domaincart.Builder{}
				},
				nil,
				nil,
			)

			if err := cob.cartStorage.StoreCart(tt.args.cart); err != nil {
				t.Fatalf("cart could not be initialized")
			}

			got, _, _ := cob.RemoveVoucher(tt.args.ctx, tt.args.cart, tt.args.couponCodeToRemove)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryBehaviour.RemoveVoucher() got = %v, want %v", got, tt.want)
			}

		})
	}
}
