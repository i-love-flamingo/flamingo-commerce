package infrastructure

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

func TestInMemoryBehaviour_CleanCart(t *testing.T) {
	t.Parallel()
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
			cob := &DefaultCartBehaviour{}
			cob.Inject(
				newInMemoryStorage(),
				nil,
				flamingo.NullLogger{},
				nil,
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

			if err := cob.cartStorage.StoreCart(context.Background(), cart); err != nil {
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
	t.Parallel()
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
			cob := &DefaultCartBehaviour{}
			cob.Inject(
				newInMemoryStorage(),
				nil,
				flamingo.NullLogger{},
				nil,
				nil,
				nil,
			)
			if err := cob.cartStorage.StoreCart(context.Background(), tt.args.cart); err != nil {
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

func TestInMemoryBehaviour_ApplyVoucher(t *testing.T) {
	t.Parallel()
	type args struct {
		cart        *domaincart.Cart
		voucherCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []domaincart.CouponCode
		wantErr bool
	}{
		{
			name: "apply valid voucher - success",
			args: args{
				cart:        &domaincart.Cart{},
				voucherCode: "valid_voucher",
			},
			want: []domaincart.CouponCode{
				{
					Code: "valid_voucher",
				},
			},
			wantErr: false,
		},
		{
			name: "apply invalid giftcard - failure",
			args: args{
				cart:        &domaincart.Cart{},
				voucherCode: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cob := &DefaultCartBehaviour{}
			cob.Inject(
				newInMemoryStorage(),
				nil,
				flamingo.NullLogger{},
				&DefaultVoucherHandler{},
				&DefaultGiftCardHandler{},
				nil,
			)
			got, _, err := cob.ApplyVoucher(context.Background(), tt.args.cart, tt.args.voucherCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultCartBehaviour.ApplyVoucher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got.AppliedCouponCodes)
			}
		})
	}
}

func TestInMemoryBehaviour_RemoveVoucher(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx                context.Context
		cart               *domaincart.Cart
		couponCodeToRemove string
	}
	tests := []struct {
		name string
		args args
		want []domaincart.CouponCode
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
			want: []domaincart.CouponCode{
				{Code: "OFF20"},
				{Code: "SALE"},
			},
		},
		{
			name: "Remove voucher from cart without vouchers",
			args: args{
				ctx:                nil,
				cart:               &domaincart.Cart{},
				couponCodeToRemove: "dummy-voucher-20",
			},
			want: nil,
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
			want: []domaincart.CouponCode{
				{Code: "OFF20"},
				{Code: "dummy-voucher-20"},
				{Code: "SALE"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cob := &DefaultCartBehaviour{}
			cob.Inject(
				newInMemoryStorage(),
				nil,
				flamingo.NullLogger{},
				&DefaultVoucherHandler{},
				&DefaultGiftCardHandler{},
				nil,
			)

			if err := cob.cartStorage.StoreCart(context.Background(), tt.args.cart); err != nil {
				t.Fatalf("cart could not be initialized")
			}

			got, _, err := cob.RemoveVoucher(tt.args.ctx, tt.args.cart, tt.args.couponCodeToRemove)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got.AppliedCouponCodes)
		})
	}
}

func TestInMemoryBehaviour_ApplyGiftCard(t *testing.T) {
	t.Parallel()
	type args struct {
		cart         *domaincart.Cart
		giftCardCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []domaincart.AppliedGiftCard
		wantErr bool
	}{
		{
			name: "apply valid giftcard - success",
			args: args{
				cart:         &domaincart.Cart{DefaultCurrency: "$"},
				giftCardCode: "valid_giftcard",
			},
			want: []domaincart.AppliedGiftCard{
				{
					Code:      "valid_giftcard",
					Applied:   priceDomain.NewFromInt(10, 100, "$"),
					Remaining: priceDomain.NewFromInt(0, 100, "$"),
				},
			},
			wantErr: false,
		},
		{
			name: "apply invalid giftcard - failure",
			args: args{
				cart:         &domaincart.Cart{},
				giftCardCode: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cob := &DefaultCartBehaviour{}
			cob.Inject(
				newInMemoryStorage(),
				nil,
				flamingo.NullLogger{},
				&DefaultVoucherHandler{},
				&DefaultGiftCardHandler{},
				&struct {
					DefaultTaxRate  float64 `inject:"config:commerce.cart.defaultCartAdapter.defaultTaxRate,optional"`
					ProductPricing  string  `inject:"config:commerce.cart.defaultCartAdapter.productPrices"`
					DefaultCurrency string  `inject:"config:commerce.cart.defaultCartAdapter.defaultCurrency"`
				}{DefaultCurrency: "$"},
			)
			got, _, err := cob.ApplyGiftCard(context.Background(), tt.args.cart, tt.args.giftCardCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultCartBehaviour.ApplyGiftCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got.AppliedGiftCards)
			}
		})
	}
}

func TestInMemoryBehaviour_RemoveGiftCard(t *testing.T) {
	t.Parallel()
	type args struct {
		cart         *domaincart.Cart
		giftCardCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []domaincart.AppliedGiftCard
		wantErr bool
	}{
		{
			name: "remove giftcard successfully",
			args: args{
				cart: &domaincart.Cart{
					AppliedGiftCards: []domaincart.AppliedGiftCard{
						{
							Code:      "to-remove",
							Applied:   priceDomain.NewFromInt(10, 100, "$"),
							Remaining: priceDomain.NewFromInt(0, 100, "$"),
						},
						{
							Code:      "valid",
							Applied:   priceDomain.NewFromInt(10, 100, "$"),
							Remaining: priceDomain.NewFromInt(0, 100, "$"),
						},
					},
				},
				giftCardCode: "to-remove",
			},
			want: []domaincart.AppliedGiftCard{
				{
					Code:      "valid",
					Applied:   priceDomain.NewFromInt(10, 100, "$"),
					Remaining: priceDomain.NewFromInt(0, 100, "$"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cob := &DefaultCartBehaviour{}
			cob.Inject(
				newInMemoryStorage(),
				nil,
				flamingo.NullLogger{},
				&DefaultVoucherHandler{},
				&DefaultGiftCardHandler{},
				nil,
			)
			got, _, err := cob.RemoveGiftCard(context.Background(), tt.args.cart, tt.args.giftCardCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultCartBehaviour.ApplyGiftCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got.AppliedGiftCards)
			}
		})
	}
}

func TestInMemoryBehaviour_Complete(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)
		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{ID: "test-id"})
		assert.NoError(t, err)

		got, _, err := cob.Complete(context.Background(), cart)
		assert.NoError(t, err)
		assert.Equal(t, cart, got)

		_, err = cob.GetCart(context.Background(), "test-id")
		assert.Error(t, err, "Cart should not be stored any more")
	})
}

func TestInMemoryBehaviour_Restore(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)
		cart := &domaincart.Cart{ID: "1234"}

		got, _, err := cob.Restore(context.Background(), cart)
		assert.NoError(t, err)

		_, err = cob.GetCart(context.Background(), got.ID)
		assert.Nil(t, err)
	})
}

func newInMemoryStorage() *InMemoryCartStorage {
	result := &InMemoryCartStorage{}
	result.Inject()

	return result
}

func TestDefaultCartBehaviour_createCartItemFromProduct(t *testing.T) {
	t.Run("gross", func(t *testing.T) {
		cob := DefaultCartBehaviour{}
		cob.Inject(nil, nil, flamingo.NullLogger{}, nil, nil, &struct {
			DefaultTaxRate  float64 `inject:"config:commerce.cart.defaultCartAdapter.defaultTaxRate,optional"`
			ProductPricing  string  `inject:"config:commerce.cart.defaultCartAdapter.productPrices"`
			DefaultCurrency string  `inject:"config:commerce.cart.defaultCartAdapter.defaultCurrency"`
		}{ProductPricing: "gross", DefaultTaxRate: 10.0, DefaultCurrency: "€"})

		item, err := cob.createCartItemFromProduct(2, "ma", "", map[string]string{}, nil, domain.SimpleProduct{
			Saleable: domain.Saleable{
				IsSaleable: true,
				ActivePrice: domain.PriceInfo{
					Default: priceDomain.NewFromFloat(50.00, "USD"),
				},
			},
		})

		require.NoError(t, err)
		assert.True(t, item.SinglePriceGross.Equal(priceDomain.NewFromFloat(50.00, "USD")))
		assert.Equal(t, 45.45, item.SinglePriceNet.FloatAmount())
		assert.Equal(t, 100.00, item.RowPriceGross.FloatAmount())
		assert.Equal(t, 45.45*2, item.RowPriceNet.FloatAmount())
		assert.Equal(t, 100.00-45.45*2, item.TotalTaxAmount().FloatAmount())
	})

	t.Run("net", func(t *testing.T) {
		cob := DefaultCartBehaviour{}
		cob.Inject(nil, nil, flamingo.NullLogger{}, nil, nil, &struct {
			DefaultTaxRate  float64 `inject:"config:commerce.cart.defaultCartAdapter.defaultTaxRate,optional"`
			ProductPricing  string  `inject:"config:commerce.cart.defaultCartAdapter.productPrices"`
			DefaultCurrency string  `inject:"config:commerce.cart.defaultCartAdapter.defaultCurrency"`
		}{ProductPricing: "net", DefaultTaxRate: 10.0, DefaultCurrency: "€"})

		item, err := cob.createCartItemFromProduct(2, "ma", "", map[string]string{}, nil, domain.SimpleProduct{
			Saleable: domain.Saleable{
				IsSaleable: true,
				ActivePrice: domain.PriceInfo{
					Default: priceDomain.NewFromFloat(50.00, "USD"),
				},
			},
		})

		require.NoError(t, err)
		assert.True(t, item.SinglePriceNet.Equal(priceDomain.NewFromFloat(50.00, "USD")))
		assert.Equal(t, 55.00, item.SinglePriceGross.FloatAmount())
		assert.Equal(t, 55.00*2, item.RowPriceGross.FloatAmount())
		assert.Equal(t, 50.00*2, item.RowPriceNet.FloatAmount())
		assert.Equal(t, 10.0, item.TotalTaxAmount().FloatAmount())
	})

}
