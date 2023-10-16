package infrastructure

import (
	"context"
	"errors"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure/mocks"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
)

func TestDefaultCartBehaviour_CleanCart(t *testing.T) {
	t.Parallel()

	t.Run("clean cart", func(t *testing.T) {
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

		err := cob.cartStorage.StoreCart(context.Background(), cart)
		assert.NoError(t, err)

		got, _, err := cob.CleanCart(context.Background(), cart)
		assert.NoError(t, err)
		assert.Empty(t, got.Deliveries)
	})
}

func TestDefaultCartBehaviour_CleanDelivery(t *testing.T) {
	t.Parallel()

	t.Run("clean dev-1", func(t *testing.T) {
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

		cart := &domaincart.Cart{
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
		}

		err := cob.cartStorage.StoreCart(context.Background(), cart)
		assert.NoError(t, err)

		got, _, err := cob.CleanDelivery(context.Background(), cart, "dev-1")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(got.Deliveries))
		assert.Equal(t, "dev-2", got.Deliveries[0].DeliveryInfo.Code)
	})

	t.Run("delivery not found", func(t *testing.T) {
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

		cart := &domaincart.Cart{
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
		}

		err := cob.cartStorage.StoreCart(context.Background(), cart)
		assert.NoError(t, err)

		got, _, err := cob.CleanDelivery(context.Background(), cart, "dev-3")
		assert.EqualError(t, err, "DefaultCartBehaviour: delivery dev-3 not found")
		assert.Nil(t, got)
	})
}

var errInvalidVoucher = errors.New("invalid voucher")
var errInvalidGiftCard = errors.New("invalid gift card")

func TestDefaultCartBehaviour_ApplyVoucher(t *testing.T) {
	t.Parallel()

	t.Run("apply voucher successful", func(t *testing.T) {
		t.Parallel()

		voucherHandler := mocks.NewVoucherHandler(t)
		voucherHandler.EXPECT().ApplyVoucher(mock.Anything, mock.Anything, "voucher").
			Return(&domaincart.Cart{ID: "voucher"}, nil)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			voucherHandler,
			nil,
			nil,
		)

		got, _, err := cob.ApplyVoucher(context.Background(), &domaincart.Cart{ID: "test"}, "voucher")
		assert.NoError(t, err)
		assert.Equal(t, "voucher", got.ID)
		voucherHandler.AssertCalled(t, "ApplyVoucher", mock.Anything, mock.Anything, "voucher")
	})

	t.Run("apply voucher error", func(t *testing.T) {
		t.Parallel()

		voucherHandler := mocks.NewVoucherHandler(t)
		voucherHandler.EXPECT().ApplyVoucher(mock.Anything, mock.Anything, "voucher").
			Return(nil, errInvalidVoucher)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			voucherHandler,
			nil,
			nil,
		)

		got, _, err := cob.ApplyVoucher(context.Background(), &domaincart.Cart{ID: "test"}, "voucher")
		assert.EqualError(t, err, "invalid voucher")
		assert.Nil(t, got)
		voucherHandler.AssertCalled(t, "ApplyVoucher", mock.Anything, mock.Anything, "voucher")
	})
}

func TestDefaultCartBehaviour_RemoveVoucher(t *testing.T) {
	t.Parallel()

	t.Run("remove voucher successful", func(t *testing.T) {
		t.Parallel()

		voucherHandler := mocks.NewVoucherHandler(t)
		voucherHandler.EXPECT().RemoveVoucher(mock.Anything, mock.Anything, "voucher").
			Return(&domaincart.Cart{ID: "voucher"}, nil)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			voucherHandler,
			nil,
			nil,
		)

		got, _, err := cob.RemoveVoucher(context.Background(), &domaincart.Cart{ID: "test"}, "voucher")
		assert.NoError(t, err)
		assert.Equal(t, "voucher", got.ID)
		voucherHandler.AssertCalled(t, "RemoveVoucher", mock.Anything, mock.Anything, "voucher")
	})

	t.Run("remove voucher error", func(t *testing.T) {
		t.Parallel()

		voucherHandler := mocks.NewVoucherHandler(t)
		voucherHandler.EXPECT().RemoveVoucher(mock.Anything, mock.Anything, "voucher").
			Return(nil, errInvalidVoucher)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			voucherHandler,
			nil,
			nil,
		)

		got, _, err := cob.RemoveVoucher(context.Background(), &domaincart.Cart{ID: "test"}, "voucher")
		assert.EqualError(t, err, "invalid voucher")
		assert.Nil(t, got)
		voucherHandler.AssertCalled(t, "RemoveVoucher", mock.Anything, mock.Anything, "voucher")
	})
}

func TestDefaultCartBehaviour_ApplyGiftCard(t *testing.T) {
	t.Parallel()

	t.Run("apply gift card successful", func(t *testing.T) {
		t.Parallel()

		giftCardHandler := mocks.NewGiftCardHandler(t)
		giftCardHandler.EXPECT().ApplyGiftCard(mock.Anything, mock.Anything, "giftCard").
			Return(&domaincart.Cart{ID: "giftCard"}, nil)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			giftCardHandler,
			nil,
		)

		got, _, err := cob.ApplyGiftCard(context.Background(), &domaincart.Cart{ID: "test"}, "giftCard")
		assert.NoError(t, err)
		assert.Equal(t, "giftCard", got.ID)
		giftCardHandler.AssertCalled(t, "ApplyGiftCard", mock.Anything, mock.Anything, "giftCard")
	})

	t.Run("apply gift card error", func(t *testing.T) {
		t.Parallel()

		giftCardHandler := mocks.NewGiftCardHandler(t)
		giftCardHandler.EXPECT().ApplyGiftCard(mock.Anything, mock.Anything, "giftCard").
			Return(nil, errInvalidGiftCard)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			giftCardHandler,
			nil,
		)

		got, _, err := cob.ApplyGiftCard(context.Background(), &domaincart.Cart{ID: "test"}, "giftCard")
		assert.EqualError(t, err, "invalid gift card")
		assert.Nil(t, got)
		giftCardHandler.AssertCalled(t, "ApplyGiftCard", mock.Anything, mock.Anything, "giftCard")
	})
}

func TestDefaultCartBehaviour_RemoveGiftCard(t *testing.T) {
	t.Parallel()

	t.Run("remove gift card successful", func(t *testing.T) {
		t.Parallel()

		giftCardHandler := mocks.NewGiftCardHandler(t)
		giftCardHandler.EXPECT().RemoveGiftCard(mock.Anything, mock.Anything, "giftCard").
			Return(&domaincart.Cart{ID: "giftCard"}, nil)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			giftCardHandler,
			nil,
		)

		got, _, err := cob.RemoveGiftCard(context.Background(), &domaincart.Cart{ID: "test"}, "giftCard")
		assert.NoError(t, err)
		assert.Equal(t, "giftCard", got.ID)
		giftCardHandler.AssertCalled(t, "RemoveGiftCard", mock.Anything, mock.Anything, "giftCard")
	})

	t.Run("remove gift card error", func(t *testing.T) {
		t.Parallel()

		giftCardHandler := mocks.NewGiftCardHandler(t)
		giftCardHandler.EXPECT().RemoveGiftCard(mock.Anything, mock.Anything, "giftCard").
			Return(nil, errInvalidGiftCard)

		cob := &DefaultCartBehaviour{}
		cob.Inject(
			newInMemoryStorage(),
			nil,
			flamingo.NullLogger{},
			nil,
			giftCardHandler,
			nil,
		)

		got, _, err := cob.RemoveGiftCard(context.Background(), &domaincart.Cart{ID: "test"}, "giftCard")
		assert.EqualError(t, err, "invalid gift card")
		assert.Nil(t, got)
		giftCardHandler.AssertCalled(t, "RemoveGiftCard", mock.Anything, mock.Anything, "giftCard")
	})
}

func TestDefaultCartBehaviour_Complete(t *testing.T) {
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

func TestDefaultCartBehaviour_Restore(t *testing.T) {
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
		assert.EqualError(t, err, "DefaultCartBehaviour: error finding delivery of item: delivery not found for \"abc\"")
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
							SinglePriceGross:          priceDomain.NewFromFloat(20.99, "EUR"),
							RowPriceGross:             priceDomain.NewFromFloat(20.99, "EUR"),
							RowPriceGrossWithDiscount: priceDomain.NewFromFloat(20.98, "EUR"),
							TotalDiscountAmount:       priceDomain.NewFromFloat(0.01, "EUR"),
							ItemRelatedDiscountAmount: priceDomain.NewFromFloat(0.01, "EUR"),
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
		assert.Equal(t, 20.99, got.Deliveries[0].Cartitems[0].SinglePriceGross.FloatAmount())
		assert.Equal(t, 41.98, got.Deliveries[0].Cartitems[0].RowPriceGross.FloatAmount())
		assert.Equal(t, 41.97, got.Deliveries[0].Cartitems[0].RowPriceGrossWithDiscount.FloatAmount())
		assert.Equal(t, 0.01, got.Deliveries[0].Cartitems[0].TotalDiscountAmount.FloatAmount())
		assert.Equal(t, 0.01, got.Deliveries[0].Cartitems[0].ItemRelatedDiscountAmount.FloatAmount())
		assert.Equal(t, 41.97, got.GrandTotal.FloatAmount())
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

	t.Run("update without qty still updates other fields", func(t *testing.T) {
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
			ID: "cartID",
			Deliveries: []domaincart.Delivery{
				{
					Cartitems: []domaincart.Item{
						{
							ID:              "firstItem",
							Qty:             1,
							MarketplaceCode: "fake_fixed_simple_without_discounts",
						},
						{
							ID:              "secondItem",
							MarketplaceCode: "fake_bundle",
							Qty:             1,
							BundleConfig: map[domain.Identifier]domain.ChoiceConfiguration{
								"identifier1": {
									MarketplaceCode: "simple_option1",
									Qty:             1,
								},
								"identifier2": {
									MarketplaceCode:        "configurable_option2",
									VariantMarketplaceCode: "shirt-black-s",
									Qty:                    1,
								},
							},
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		source := "baz"
		got, _, err := cob.UpdateItems(context.Background(), cart, []domaincart.ItemUpdateCommand{
			{
				ItemID:         "firstItem",
				AdditionalData: map[string]string{"foo": "bar"},
				SourceID:       &source,
			},
			{
				ItemID: "secondItem",
				BundleConfiguration: map[domain.Identifier]domain.ChoiceConfiguration{
					"identifier1": {
						MarketplaceCode: "simple_option2",
						Qty:             1,
					},
					"identifier2": {
						MarketplaceCode:        "configurable_option1",
						VariantMarketplaceCode: "shirt-black-l",
						Qty:                    1,
					},
				},
			},
		})

		assert.NoError(t, err)

		assert.Equal(t, "firstItem", got.Deliveries[0].Cartitems[0].ID)
		assert.Equal(t, "fake_fixed_simple_without_discounts", got.Deliveries[0].Cartitems[0].MarketplaceCode)
		assert.Equal(t, map[string]string{"foo": "bar"}, got.Deliveries[0].Cartitems[0].AdditionalData)
		assert.Equal(t, "baz", got.Deliveries[0].Cartitems[0].SourceID)

		assert.Equal(t, "secondItem", got.Deliveries[0].Cartitems[1].ID)
		assert.Equal(t, "fake_bundle", got.Deliveries[0].Cartitems[1].MarketplaceCode)
		assert.Equal(t, "simple_option2", got.Deliveries[0].Cartitems[1].BundleConfig["identifier1"].MarketplaceCode)
		assert.Equal(t, 1, got.Deliveries[0].Cartitems[1].BundleConfig["identifier1"].Qty)
		assert.Equal(t, "configurable_option1", got.Deliveries[0].Cartitems[1].BundleConfig["identifier2"].MarketplaceCode)
		assert.Equal(t, "shirt-black-l", got.Deliveries[0].Cartitems[1].BundleConfig["identifier2"].VariantMarketplaceCode)
		assert.Equal(t, 1, got.Deliveries[0].Cartitems[1].BundleConfig["identifier2"].Qty)
	})

	t.Run("update bundle configuration for a cart item", func(t *testing.T) {
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
							ID:              "1234",
							MarketplaceCode: "fake_bundle",
							Qty:             1,
							BundleConfig: map[domain.Identifier]domain.ChoiceConfiguration{
								"identifier1": {
									MarketplaceCode: "simple_option1",
									Qty:             1,
								},
								"identifier2": {
									MarketplaceCode:        "configurable_option2",
									VariantMarketplaceCode: "shirt-black-s",
									Qty:                    1,
								},
							},
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		qty := 1
		got, _, err := cob.UpdateItems(context.Background(), cart, []domaincart.ItemUpdateCommand{{
			ItemID: "1234",
			Qty:    &qty,
			BundleConfiguration: map[domain.Identifier]domain.ChoiceConfiguration{
				"identifier1": {
					MarketplaceCode: "simple_option2",
					Qty:             2,
				},
				"identifier2": {
					MarketplaceCode:        "configurable_option1",
					VariantMarketplaceCode: "shirt-black-l",
					Qty:                    2,
				},
			},
		}})

		assert.NoError(t, err)
		assert.Equal(t, "fake_bundle", got.Deliveries[0].Cartitems[0].MarketplaceCode)
		assert.Equal(t, "simple_option2", got.Deliveries[0].Cartitems[0].BundleConfig["identifier1"].MarketplaceCode)
		assert.Equal(t, 2, got.Deliveries[0].Cartitems[0].BundleConfig["identifier1"].Qty)
		assert.Equal(t, "configurable_option1", got.Deliveries[0].Cartitems[0].BundleConfig["identifier2"].MarketplaceCode)
		assert.Equal(t, "shirt-black-l", got.Deliveries[0].Cartitems[0].BundleConfig["identifier2"].VariantMarketplaceCode)
		assert.Equal(t, 2, got.Deliveries[0].Cartitems[0].BundleConfig["identifier2"].Qty)
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

	t.Run("adding the same configurable product with different active variant", func(t *testing.T) {
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
							MarketplaceCode:        "fake_configurable_with_active_variant",
							VariantMarketPlaceCode: "shirt-red-s",
							Qty:                    1,
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.AddToCart(context.Background(), cart, "delivery", domaincart.AddRequest{
			MarketplaceCode:        "fake_configurable_with_active_variant",
			VariantMarketplaceCode: "shirt-black-l",
			Qty:                    1,
		})

		assert.NoError(t, err)
		assert.Equal(t, 2, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "fake_configurable_with_active_variant", got.Deliveries[0].Cartitems[0].MarketplaceCode)
		assert.Equal(t, "shirt-red-s", got.Deliveries[0].Cartitems[0].VariantMarketPlaceCode)
		assert.Equal(t, "fake_configurable_with_active_variant", got.Deliveries[0].Cartitems[1].MarketplaceCode)
		assert.Equal(t, "shirt-black-l", got.Deliveries[0].Cartitems[1].VariantMarketPlaceCode)
	})

	t.Run("adding the same bundle product with different choices", func(t *testing.T) {
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
							MarketplaceCode: "fake_bundle",
							Qty:             1,
							BundleConfig: map[domain.Identifier]domain.ChoiceConfiguration{
								"identifier1": {
									MarketplaceCode: "simple_option1",
									Qty:             1,
								},
								"identifier2": {
									MarketplaceCode:        "configurable_option1",
									VariantMarketplaceCode: "shirt-black-l",
									Qty:                    1,
								},
							},
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		got, _, err := cob.AddToCart(context.Background(), cart, "delivery", domaincart.AddRequest{
			MarketplaceCode: "fake_bundle",
			Qty:             1,
			BundleConfiguration: map[domain.Identifier]domain.ChoiceConfiguration{
				"identifier1": {
					MarketplaceCode: "simple_option2",
					Qty:             1,
				},
				"identifier2": {
					MarketplaceCode:        "configurable_option1",
					VariantMarketplaceCode: "shirt-black-l",
					Qty:                    1,
				},
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, 2, len(got.Deliveries[0].Cartitems))
		assert.Equal(t, "fake_bundle", got.Deliveries[0].Cartitems[0].MarketplaceCode)
		assert.Equal(t, "simple_option1", got.Deliveries[0].Cartitems[0].BundleConfig["identifier1"].MarketplaceCode)
		assert.Equal(t, "fake_bundle", got.Deliveries[0].Cartitems[1].MarketplaceCode)
		assert.Equal(t, "simple_option2", got.Deliveries[0].Cartitems[1].BundleConfig["identifier1"].MarketplaceCode)
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
