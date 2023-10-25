package application_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/application/mocks"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/core/auth"
	authMock "flamingo.me/flamingo/v3/core/auth/mock"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/mock"
)

func TestEventReceiver_Notify(t *testing.T) {
	t.Parallel()

	t.Run("Invalidate cart cache on WebLogoutEvent", func(t *testing.T) {
		t.Parallel()

		receiver := &application.EventReceiver{}

		ctx := context.Background()
		session := web.EmptySession()

		cartCache := mocks.NewCartCache(t)
		cartCache.EXPECT().DeleteAll(ctx, session).Return(nil)

		receiver.Inject(flamingo.NullLogger{}, nil, nil, nil, &struct {
			CartCache application.CartCache `inject:",optional"`
		}{cartCache})

		receiver.Notify(ctx, &auth.WebLogoutEvent{
			Request: web.CreateRequest(nil, session),
			Broker:  "example",
		})
	})

	t.Run("Invalidate cart cache on InvalidateCartEvent", func(t *testing.T) {
		t.Parallel()

		receiver := &application.EventReceiver{}

		ctx := context.Background()
		session := web.EmptySession()

		cartCache := mocks.NewCartCache(t)
		cartCache.EXPECT().BuildIdentifier(ctx, session).Return(application.CartCacheIdentifier{
			GuestCartID:    "foo",
			IsCustomerCart: false,
			CustomerID:     "",
		}, nil)

		cartCache.EXPECT().Invalidate(ctx, session, application.CartCacheIdentifier{
			GuestCartID:    "foo",
			IsCustomerCart: false,
			CustomerID:     "",
		}).Return(nil)

		receiver.Inject(flamingo.NullLogger{}, nil, nil, nil, &struct {
			CartCache application.CartCache `inject:",optional"`
		}{cartCache})

		receiver.Notify(ctx, &cart.InvalidateCartEvent{Session: web.EmptySession()})
	})

	t.Run("CartMerger is called on WebLoginEvent", func(t *testing.T) {
		t.Parallel()

		session := web.EmptySession()
		request := web.CreateRequest(nil, session)

		guestCart := cart.Cart{ID: "guestCart", BelongsToAuthenticatedUser: false}
		customerCart := cart.Cart{ID: "customerCart", BelongsToAuthenticatedUser: true}

		cartReceiver := mocks.NewCartReceiver(t)
		cartReceiver.EXPECT().ShouldHaveGuestCart(session).Return(true)
		cartReceiver.EXPECT().ViewGuestCart(mock.Anything, session).Return(&guestCart, nil)
		cartReceiver.EXPECT().ViewCart(mock.Anything, session).Return(&customerCart, nil)
		cartMerger := mocks.NewCartMerger(t)
		cartMerger.EXPECT().Merge(mock.Anything, session, guestCart, customerCart)

		cartCache := mocks.NewCartCache(t)
		cartCache.EXPECT().BuildIdentifier(mock.Anything, session).Return(application.CartCacheIdentifier{}, nil)
		cartCache.EXPECT().Delete(mock.Anything, session, application.CartCacheIdentifier{}).Return(nil)
		eventRouter := &MockEventRouter{}

		receiver := &application.EventReceiver{}
		receiver.Inject(flamingo.NullLogger{}, cartReceiver, eventRouter, cartMerger, &struct {
			CartCache application.CartCache `inject:",optional"`
		}{cartCache})

		receiver.Notify(context.Background(), &auth.WebLoginEvent{
			Request:  request,
			Broker:   "example",
			Identity: &authMock.Identity{},
		})
	})
}

func TestCartMergeStrategyReplace_Merge(t *testing.T) {
	t.Parallel()

	session := web.EmptySession()

	c := &application.CartMergeStrategyReplace{}
	cartService := mocks.NewCartService(t)
	cartService.EXPECT().Clean(mock.Anything, session).Return(nil)
	cartService.EXPECT().UpdateDeliveryInfo(mock.Anything, session, "delivery1", mock.Anything).Return(nil)
	cartService.EXPECT().AddProduct(mock.Anything, session, "delivery1", cart.AddRequest{
		MarketplaceCode: "foo",
		Qty:             1,
	}).Return(nil, nil)
	cartService.EXPECT().AddProduct(mock.Anything, session, "delivery1", cart.AddRequest{
		MarketplaceCode: "bundle",
		Qty:             1,
		BundleConfiguration: map[domain.Identifier]domain.ChoiceConfiguration{"slot1": {
			MarketplaceCode:        "bar",
			VariantMarketplaceCode: "baz",
			Qty:                    2,
		}},
	}).Return(nil, nil)
	cartService.EXPECT().UpdateBillingAddress(mock.Anything, session, mock.Anything).Return(nil)
	cartService.EXPECT().UpdatePurchaser(mock.Anything, session, mock.Anything, mock.Anything).Return(nil)
	cartService.EXPECT().ApplyVoucher(mock.Anything, session, "SUMMER_SALE").Return(&cart.Cart{}, nil)
	cartService.EXPECT().ApplyGiftCard(mock.Anything, session, "GHDJAHJH-DADAD-2113").Return(&cart.Cart{}, nil)
	cartService.EXPECT().UpdatePaymentSelection(mock.Anything, session, mock.Anything).Return(nil)
	c.Inject(flamingo.NullLogger{}, cartService)
	c.Merge(context.Background(), session, cart.Cart{
		ID: "guest", BelongsToAuthenticatedUser: false,
		Deliveries: []cart.Delivery{{
			DeliveryInfo: cart.DeliveryInfo{Code: "delivery1"},
			Cartitems: []cart.Item{
				{MarketplaceCode: "foo", Qty: 1},
				{

					MarketplaceCode: "bundle",
					BundleConfig: map[domain.Identifier]domain.ChoiceConfiguration{"slot1": {
						MarketplaceCode:        "bar",
						VariantMarketplaceCode: "baz",
						Qty:                    2,
					}},
					Qty: 1,
				},
			},
		}},
		BillingAddress:     &cart.Address{},
		Purchaser:          &cart.Person{},
		AppliedCouponCodes: []cart.CouponCode{{Code: "SUMMER_SALE"}},
		AppliedGiftCards:   []cart.AppliedGiftCard{{Code: "GHDJAHJH-DADAD-2113"}},
		PaymentSelection:   cart.DefaultPaymentSelection{},
	}, cart.Cart{ID: "customer", BelongsToAuthenticatedUser: true})
}

func TestCartMergeStrategyMerge_Merge(t *testing.T) {
	t.Parallel()

	session := web.EmptySession()

	c := &application.CartMergeStrategyMerge{}
	cartService := mocks.NewCartService(t)
	cartService.EXPECT().UpdateDeliveryInfo(mock.Anything, session, "delivery1", mock.Anything).Return(nil)
	cartService.EXPECT().AddProduct(mock.Anything, session, "delivery1", cart.AddRequest{
		MarketplaceCode: "foo",
		Qty:             1,
	}).Return(nil, nil)
	cartService.EXPECT().AddProduct(mock.Anything, session, "delivery1", cart.AddRequest{
		MarketplaceCode: "bundle",
		Qty:             1,
		BundleConfiguration: map[domain.Identifier]domain.ChoiceConfiguration{"slot1": {
			MarketplaceCode:        "bar",
			VariantMarketplaceCode: "baz",
			Qty:                    2,
		}},
	}).Return(nil, nil)
	cartService.EXPECT().UpdateBillingAddress(mock.Anything, session, mock.Anything).Return(nil)
	cartService.EXPECT().UpdatePurchaser(mock.Anything, session, mock.Anything, mock.Anything).Return(nil)
	cartService.EXPECT().ApplyVoucher(mock.Anything, session, "SUMMER_SALE").Return(&cart.Cart{}, nil)
	cartService.EXPECT().ApplyGiftCard(mock.Anything, session, "GHDJAHJH-DADAD-2113").Return(&cart.Cart{}, nil)
	cartService.EXPECT().UpdatePaymentSelection(mock.Anything, session, mock.Anything).Return(nil)
	c.Inject(flamingo.NullLogger{}, cartService)
	c.Merge(context.Background(), session, cart.Cart{
		ID: "guest", BelongsToAuthenticatedUser: false,
		Deliveries: []cart.Delivery{{
			DeliveryInfo: cart.DeliveryInfo{Code: "delivery1"},
			Cartitems: []cart.Item{
				{MarketplaceCode: "foo", Qty: 1},
				{

					MarketplaceCode: "bundle",
					BundleConfig: map[domain.Identifier]domain.ChoiceConfiguration{"slot1": {
						MarketplaceCode:        "bar",
						VariantMarketplaceCode: "baz",
						Qty:                    2,
					}},
					Qty: 1,
				},
			},
		}},
		BillingAddress:     &cart.Address{},
		Purchaser:          &cart.Person{},
		AppliedCouponCodes: []cart.CouponCode{{Code: "SUMMER_SALE"}},
		AppliedGiftCards:   []cart.AppliedGiftCard{{Code: "GHDJAHJH-DADAD-2113"}},
		PaymentSelection:   cart.DefaultPaymentSelection{},
	}, cart.Cart{ID: "customer", BelongsToAuthenticatedUser: true})
}

func TestCartMergeStrategyNone_Merge(t *testing.T) {
	t.Parallel()

	c := &application.CartMergeStrategyNone{}
	c.Merge(context.Background(), nil, cart.Cart{}, cart.Cart{})
}
