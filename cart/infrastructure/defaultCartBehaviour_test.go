package infrastructure

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
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
		t.Parallel()

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
		t.Parallel()

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

func TestDefaultCartBehaviour_DeleteItem(t *testing.T) {
	t.Parallel()

	t.Run("cart does not exist", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart := &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
					Cartitems: []domaincart.Item{
						{
							ID: "abc",
						},
					},
				},
			},
		}

		_, _, err := cob.DeleteItem(context.Background(), cart, "abc", "delivery")
		assert.Error(t, err)
	})

	t.Run("item in first place in delivery deleted", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
					Cartitems: []domaincart.Item{
						{
							ID: "id-1",
						},
						{
							ID: "id-2",
						},
						{
							ID: "id-3",
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.DeleteItem(context.Background(), cart, "id-1", "delivery")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "id-2", got.Deliveries[0].Cartitems[0].ID)
		assert.Equal(t, "id-3", got.Deliveries[0].Cartitems[1].ID)
	})

	t.Run("item in middle of delivery deleted", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
					Cartitems: []domaincart.Item{
						{
							ID: "id-1",
						},
						{
							ID: "id-2",
						},
						{
							ID: "id-3",
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.DeleteItem(context.Background(), cart, "id-2", "delivery")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "id-1", got.Deliveries[0].Cartitems[0].ID)
		assert.Equal(t, "id-3", got.Deliveries[0].Cartitems[1].ID)
	})

	t.Run("item in last place in delivery deleted", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
					Cartitems: []domaincart.Item{
						{
							ID: "id-1",
						},
						{
							ID: "id-2",
						},
						{
							ID: "id-3",
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.DeleteItem(context.Background(), cart, "id-3", "delivery")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "id-1", got.Deliveries[0].Cartitems[0].ID)
		assert.Equal(t, "id-2", got.Deliveries[0].Cartitems[1].ID)
	})

	t.Run("item in different delivery not deleted", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
					Cartitems: []domaincart.Item{
						{
							ID: "abc",
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.DeleteItem(context.Background(), cart, "abc", "delivery-2")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "abc", got.Deliveries[0].Cartitems[0].ID)
	})
}

func TestDefaultCartBehaviour_UpdateItems(t *testing.T) {
	t.Parallel()

	t.Run("cart does not exist", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart := &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
					Cartitems: []domaincart.Item{
						{
							ID: "abc",
						},
					},
				},
			},
		}

		_, _, err := cob.UpdateItems(context.Background(), cart, []domaincart.ItemUpdateCommand{})
		assert.Error(t, err)
	})

	t.Run("item not found in delivery", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
		})
		assert.NoError(t, err)

		got, _, err := cob.UpdateItems(context.Background(), cart, []domaincart.ItemUpdateCommand{
			{
				ItemID: "abc",
			},
		})
		assert.Empty(t, got)
		assert.EqualError(t, err, "cart.infrastructure.DefaultCartBehaviour: error on finding delivery of item: delivery not found for \"abc\"")
	})

	t.Run("item updated in delivery", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					Cartitems: []domaincart.Item{
						{
							ID:              "abc",
							Qty:             1,
							MarketplaceCode: "fake_fixed_simple_without_discounts",
							AdditionalData: map[string]string{
								"1": "a",
							},
							SinglePriceGross: priceDomain.NewFromFloat(20.99, "EUR"),
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		qty := 2
		got, _, err := cob.UpdateItems(context.Background(), cart, []domaincart.ItemUpdateCommand{
			{
				ItemID: "abc",
				Qty:    &qty,
				AdditionalData: map[string]string{
					"2": "b",
				},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"2": "b"}, got.Deliveries[0].Cartitems[0].AdditionalData)
		assert.Equal(t, 2, got.Deliveries[0].Cartitems[0].Qty)
		assert.Equal(t, 41.98, got.GrandTotal.FloatAmount())
	})

	t.Run("item updated with qty 0 is deleted", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					Cartitems: []domaincart.Item{
						{
							ID:              "abc",
							Qty:             1,
							MarketplaceCode: "fake_fixed_simple_without_discounts",
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		qty := 0
		got, _, err := cob.UpdateItems(context.Background(), cart, []domaincart.ItemUpdateCommand{
			{
				ItemID: "abc",
				Qty:    &qty,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, 0.0, got.GrandTotal.FloatAmount())
	})
}

func TestDefaultCartBehaviour_AddToCart(t *testing.T) {
	t.Parallel()

	t.Run("adding product to empty cart adds delivery", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
		})
		assert.NoError(t, err)

		got, _, err := cob.AddToCart(context.Background(), cart, "delivery", domaincart.AddRequest{
			MarketplaceCode: "fake_fixed_simple_without_discounts",
			Qty:             1,
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(got.Deliveries))
		assert.Equal(t, 1, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "fake_fixed_simple_without_discounts", got.Deliveries[0].Cartitems[0].MarketplaceCode)
		assert.Equal(t, 1, got.Deliveries[0].Cartitems[0].Qty)
		assert.Equal(t, 20.99, got.GrandTotal.FloatAmount())
	})

	t.Run("adding the same product increases qty", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			&fake.ProductService{},
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
					Cartitems: []domaincart.Item{
						{
							MarketplaceCode: "fake_fixed_simple_without_discounts",
							Qty:             1,
							AdditionalData: map[string]string{
								"1": "a",
							},
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.AddToCart(context.Background(), cart, "delivery", domaincart.AddRequest{
			MarketplaceCode: "fake_fixed_simple_without_discounts",
			Qty:             1,
			AdditionalData: map[string]string{
				"2": "b",
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "fake_fixed_simple_without_discounts", got.Deliveries[0].Cartitems[0].MarketplaceCode)
		assert.Equal(t, map[string]string{"1": "a", "2": "b"}, got.Deliveries[0].Cartitems[0].AdditionalData)
		assert.Equal(t, 2, got.Deliveries[0].Cartitems[0].Qty)
		assert.Equal(t, 41.98, got.GrandTotal.FloatAmount())
	})
}

func TestDefaultCartBehaviour_UpdatePurchaser(t *testing.T) {
	t.Parallel()

	t.Run("additional data custom attributes are merged", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery",
					},
				},
			},
			AdditionalData: domaincart.AdditionalData{
				CustomAttributes: map[string]string{
					"1": "a",
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.UpdatePurchaser(
			context.Background(),
			cart,
			&domaincart.Person{
				Address: &domaincart.Address{
					Firstname: "test",
				},
			},
			&domaincart.AdditionalData{
				CustomAttributes: map[string]string{
					"2": "b",
				},
			})
		assert.NoError(t, err)
		assert.Equal(t, "test", got.Purchaser.Address.Firstname)
		assert.Equal(t, map[string]string{"1": "a", "2": "b"}, got.AdditionalData.CustomAttributes)
	})
}

func TestDefaultCartBehaviour_UpdateBillingAddress(t *testing.T) {
	t.Parallel()

	t.Run("add billing address", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{ID: "1234"})
		assert.NoError(t, err)

		got, _, err := cob.UpdateBillingAddress(context.Background(), cart, domaincart.Address{
			Firstname: "first",
		})
		assert.NoError(t, err)
		assert.Equal(t, "first", got.BillingAddress.Firstname)
	})

	t.Run("update billing address", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			BillingAddress: &domaincart.Address{
				Firstname: "first-1",
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.UpdateBillingAddress(context.Background(), cart, domaincart.Address{
			Firstname: "first-2",
		})
		assert.NoError(t, err)
		assert.Equal(t, "first-2", got.BillingAddress.Firstname)
	})
}

func TestDefaultCartBehaviour_UpdateAdditionalData(t *testing.T) {
	t.Parallel()

	t.Run("add additional data", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{ID: "1234"})
		assert.NoError(t, err)

		got, _, err := cob.UpdateAdditionalData(context.Background(), cart, &domaincart.AdditionalData{
			ReservedOrderID:  "id-1",
			CustomAttributes: map[string]string{"1": "a"},
		})
		assert.NoError(t, err)
		assert.Equal(t, "id-1", got.AdditionalData.ReservedOrderID)
		assert.Equal(t, map[string]string{"1": "a"}, got.AdditionalData.CustomAttributes)
	})

	t.Run("update additional data", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			AdditionalData: domaincart.AdditionalData{
				ReservedOrderID:  "id-1",
				CustomAttributes: map[string]string{"1": "a"},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.UpdateAdditionalData(context.Background(), cart, &domaincart.AdditionalData{
			ReservedOrderID:  "id-2",
			CustomAttributes: map[string]string{"2": "b"},
		})
		assert.NoError(t, err)
		assert.Equal(t, "id-2", got.AdditionalData.ReservedOrderID)
		assert.Equal(t, map[string]string{"2": "b"}, got.AdditionalData.CustomAttributes)
	})
}

func TestDefaultCartBehaviour_UpdatePaymentSelection(t *testing.T) {
	t.Parallel()

	t.Run("add payment selection", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{ID: "1234"})
		assert.NoError(t, err)

		got, _, err := cob.UpdatePaymentSelection(context.Background(), cart, &domaincart.DefaultPaymentSelection{
			GatewayProp: "gateway-1",
		})
		assert.NoError(t, err)
		assert.Equal(t, "gateway-1", got.PaymentSelection.Gateway())
	})

	t.Run("update payment selection", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			PaymentSelection: &domaincart.DefaultPaymentSelection{
				GatewayProp: "gateway-1",
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.UpdatePaymentSelection(context.Background(), cart, &domaincart.DefaultPaymentSelection{
			GatewayProp: "gateway-2",
		})
		assert.NoError(t, err)
		assert.Equal(t, "gateway-2", got.PaymentSelection.Gateway())
	})
}

func TestDefaultCartBehaviour_UpdateDeliveryInfo(t *testing.T) {
	t.Parallel()

	t.Run("add new delivery info", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{ID: "1234"})
		assert.NoError(t, err)

		got, _, err := cob.UpdateDeliveryInfo(context.Background(), cart, "delivery", domaincart.DeliveryInfoUpdateCommand{
			DeliveryInfo: domaincart.DeliveryInfo{
				Code: "delivery",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(got.Deliveries))
		assert.Equal(t, "delivery", got.Deliveries[0].DeliveryInfo.Code)
	})

	t.Run("update delivery info", func(t *testing.T) {
		t.Parallel()

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			nil,
			nil,
		)

		cart, err := cob.StoreNewCart(context.Background(), &domaincart.Cart{
			ID: "1234",
			Deliveries: []domaincart.Delivery{
				{
					DeliveryInfo: domaincart.DeliveryInfo{
						Code: "delivery-1",
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.UpdateDeliveryInfo(context.Background(), cart, "delivery-1", domaincart.DeliveryInfoUpdateCommand{
			DeliveryInfo: domaincart.DeliveryInfo{
				Code: "delivery-2",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(got.Deliveries))
		assert.Equal(t, "delivery-2", got.Deliveries[0].DeliveryInfo.Code)
	})
}

func newInMemoryStorage() *InMemoryCartStorage {
	result := &InMemoryCartStorage{}
	result.Inject()

	return result
}

func TestDefaultCartBehaviour_createCartItemFromProduct(t *testing.T) {
	t.Parallel()

	t.Run("gross", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
