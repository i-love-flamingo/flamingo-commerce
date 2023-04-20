package domain_test

import (
	"context"
	"errors"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/sourcing/domain"

	"github.com/stretchr/testify/assert"
)

type (
	availableSourcesProviderMock struct {
		Sources []domain.Source
		Error   error
	}
	stockProviderMock struct {
		Qty   int
		Error error
	}
	stockBySourceAndProductProviderMock struct {
		// [source.LocationCode][product.Identifier] = Qty
		Qty   map[string]map[string]int
		Error error
	}
)

var (
	_ domain.AvailableSourcesProvider = new(availableSourcesProviderMock)
	_ domain.StockProvider            = new(stockProviderMock)
	_ domain.AvailableSourcesProvider = new(stockBySourceAndProductProviderMock)
)

func (a availableSourcesProviderMock) GetPossibleSources(_ context.Context, _ productDomain.BasicProduct, _ *cart.DeliveryInfo) ([]domain.Source, error) {
	return a.Sources, a.Error
}

func (s stockProviderMock) GetStock(_ context.Context, _ productDomain.BasicProduct, _ domain.Source, _ *cart.DeliveryInfo) (int, error) {
	return s.Qty, s.Error
}

func (s stockBySourceAndProductProviderMock) GetStock(_ context.Context, product productDomain.BasicProduct, source domain.Source, _ *cart.DeliveryInfo) (int, error) {
	return s.Qty[source.LocationCode][product.GetIdentifier()], s.Error
}

func (s stockBySourceAndProductProviderMock) GetPossibleSources(_ context.Context, _ productDomain.BasicProduct, _ *cart.DeliveryInfo) ([]domain.Source, error) {
	panic("implement me")
}

func TestDefaultSourcingService_GetAvailableSources(t *testing.T) {
	t.Run("error handling on unbound providers", func(t *testing.T) {
		sourcingService := domain.DefaultSourcingService{}
		sourcingService.Inject(flamingo.NullLogger{}, nil)
		_, err := sourcingService.GetAvailableSources(context.Background(), nil, nil, nil)
		assert.EqualError(t, err, "no Source Provider bound", "received error if available sources provider and stock provider are not configured")

		sourcingService = newDefaultSourcingService(nil, nil)
		_, err = sourcingService.GetAvailableSources(context.Background(), productDomain.SimpleProduct{}, nil, nil)
		assert.EqualError(t, err, "no Stock Provider bound", "received error if stock provider is not set")
	})

	t.Run("error handing on error fetching available sources", func(t *testing.T) {
		sourcingService := domain.DefaultSourcingService{}
		sourcingService.Inject(flamingo.NullLogger{}, &struct {
			AvailableSourcesProvider domain.AvailableSourcesProvider `inject:",optional"`
			StockProvider            domain.StockProvider            `inject:",optional"`
		}{
			AvailableSourcesProvider: availableSourcesProviderMock{
				Sources: nil,
				Error:   errors.New("mocked available sources provider error"),
			},
			StockProvider: stockProviderMock{},
		})

		_, err := sourcingService.GetAvailableSources(context.Background(), productDomain.SimpleProduct{Identifier: "example"}, nil, nil)
		assert.Contains(t, err.Error(), "mocked available sources provider error", "result contains the error message of the available sources provider")
	})

	t.Run("full qty with nil cart", func(t *testing.T) {
		stubbedSources := []domain.Source{{LocationCode: "loc1"}}
		stubbedStockQty := 10

		stockProviderMock := stockProviderMock{Qty: stubbedStockQty}
		sourcingService := newDefaultSourcingService(stockProviderMock, stubbedSources)

		sources, err := sourcingService.GetAvailableSources(context.Background(), productDomain.SimpleProduct{Identifier: "simple_test"}, nil, nil)
		assert.NoError(t, err)

		expectedSources := domain.AvailableSourcesPerProduct{domain.ProductID("simple_test"): domain.AvailableSources{
			stubbedSources[0]: stubbedStockQty,
		}}

		assert.Equal(t, expectedSources, sources)
	})

	t.Run("qty reduced with existing cart", func(t *testing.T) {
		stubbedSources := []domain.Source{{LocationCode: "loc1"}}
		stubbedStockQty := 10
		stubbedQtyAlreadyInCart := 2
		stubbedProduct := productDomain.SimpleProduct{Identifier: "productid"}

		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: stubbedProduct,
							Item:    cart.Item{Qty: stubbedQtyAlreadyInCart, ID: "item1"},
						},
					},
				},
			},
		}

		stockProviderMock := stockProviderMock{Qty: stubbedStockQty}
		sourcingService := newDefaultSourcingService(stockProviderMock, stubbedSources)

		sources, err := sourcingService.GetAvailableSources(context.Background(), stubbedProduct, nil, &testCart)
		assert.NoError(t, err)

		expectedSources := domain.AvailableSourcesPerProduct{domain.ProductID("productid"): domain.AvailableSources{
			stubbedSources[0]: stubbedStockQty - stubbedQtyAlreadyInCart,
		}}
		assert.Equal(t, expectedSources, sources)
	})

	t.Run("all available qty is already in cart", func(t *testing.T) {
		stubbedSources := []domain.Source{{LocationCode: "loc1"}, {LocationCode: "loc2"}}
		stubbedStockQty := 5
		stubbedQtyAlreadyInCart := 10
		stubbedProduct := productDomain.SimpleProduct{
			Identifier: "marketPlaceCode1",
		}
		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: stubbedProduct,
							Item: cart.Item{
								ID:  "itemID1",
								Qty: stubbedQtyAlreadyInCart,
							},
						},
					},
				},
			},
		}
		stockProviderMock := stockProviderMock{Qty: stubbedStockQty}
		sourcingService := newDefaultSourcingService(stockProviderMock, stubbedSources)

		availableSources, err := sourcingService.GetAvailableSources(context.Background(), stubbedProduct, nil, &testCart)

		t.Log(availableSources)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNoSourceAvailable)
	})

	t.Run("get sources for bundle product, cart is nil, sources shouldn't be allocated", func(t *testing.T) {
		t.Parallel()

		simpleInBundle1 := productDomain.SimpleProduct{Identifier: "product1"}
		simpleInBundle2 := productDomain.SimpleProduct{Identifier: "product2"}

		bundleProduct := productDomain.BundleProductWithActiveChoices{
			BundleProduct: productDomain.BundleProduct{
				Identifier: "bundle_product",
			},
			ActiveChoices: map[productDomain.Identifier]productDomain.ActiveChoice{
				"identifier1": {
					Qty:     1,
					Product: simpleInBundle1,
				},
				"identifier2": {
					Qty:     2,
					Product: simpleInBundle2,
				},
			},
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		stubbedSources := []domain.Source{source1, source2}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"product1": 6,
					"product2": 4,
				},
				"Source2": {
					"product1": 4,
					"product2": 1,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		availableSources, err := sourcingService.GetAvailableSources(context.Background(), bundleProduct, nil, nil)
		assert.NoError(t, err)

		assert.Equal(t, 10, availableSources[domain.ProductID(simpleInBundle1.GetIdentifier())].QtySum())
		assert.Equal(t, 5, availableSources[domain.ProductID(simpleInBundle2.GetIdentifier())].QtySum())
	})

	t.Run("get sources for simple product, cart is nil, sources shouldn't be allocated", func(t *testing.T) {
		t.Parallel()

		simpleProduct := productDomain.SimpleProduct{
			Identifier: "product2",
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		stubbedSources := []domain.Source{source1, source2}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"product1": 6,
					"product2": 4,
				},
				"Source2": {
					"product1": 4,
					"product2": 1,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		availableSources, err := sourcingService.GetAvailableSources(context.Background(), simpleProduct, nil, nil)
		assert.NoError(t, err)

		assert.Equal(t, 5, availableSources[domain.ProductID(simpleProduct.GetIdentifier())].QtySum())
	})

	t.Run("get sources for simple product, cart is not nil, sources should be allocated", func(t *testing.T) {
		t.Parallel()

		simpleProduct := productDomain.SimpleProduct{
			Identifier: "product2",
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		stubbedSources := []domain.Source{source1, source2}

		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: simpleProduct,
							Item:    cart.Item{Qty: 2, ID: "item2"},
						},
					},
				},
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: simpleProduct,
							Item:    cart.Item{Qty: 1, ID: "item2"},
						},
					},
				},
			},
		}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"product1": 5,
					"product2": 4,
				},
				"Source2": {
					"product1": 5,
					"product2": 1,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		availableSources, err := sourcingService.GetAvailableSources(context.Background(), simpleProduct, nil, &testCart)
		assert.NoError(t, err)

		assert.Equal(t, 3, availableSources[domain.ProductID(simpleProduct.GetIdentifier())].QtySum())
	})

	t.Run("get sources for bundle product, cart is not nil, sources should be allocated", func(t *testing.T) {
		t.Parallel()

		simpleInBundle1 := productDomain.SimpleProduct{Identifier: "gucciSlippers"}
		simple2 := productDomain.SimpleProduct{Identifier: "gucciTShirt"}

		bundleProduct := productDomain.BundleProductWithActiveChoices{
			BundleProduct: productDomain.BundleProduct{
				Identifier: "bundle_product",
			},
			ActiveChoices: map[productDomain.Identifier]productDomain.ActiveChoice{
				"identifier1": {
					Qty:     1,
					Product: simpleInBundle1,
				},
				"identifier2": {
					Qty:     1,
					Product: simple2,
				},
			},
		}

		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: bundleProduct,
							Item:    cart.Item{Qty: 2, ID: "item1"},
						},
						{
							Product: simple2,
							Item:    cart.Item{Qty: 1, ID: "item2"},
						},
					},
				},
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: simple2,
							Item:    cart.Item{Qty: 1, ID: "item3"},
						},
					},
				},
			},
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		stubbedSources := []domain.Source{source1, source2}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"gucciSlippers": 5,
					"gucciTShirt":   4,
				},
				"Source2": {
					"gucciSlippers": 5,
					"gucciTShirt":   1,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		availableSources, err := sourcingService.GetAvailableSources(context.Background(), bundleProduct, nil, &testCart)
		assert.NoError(t, err)

		assert.Equal(t, 8, availableSources[domain.ProductID(simpleInBundle1.GetIdentifier())].QtySum())
		assert.Equal(t, 1, availableSources[domain.ProductID(simple2.GetIdentifier())].QtySum())
	})
}

func TestDefaultSourcingService_AllocateItems(t *testing.T) {
	t.Parallel()
	/**
	Given:
	Cart:
		Delivery1:
			product1 - 10
			product2 - 5
		Delivery2:
			product1 - 5

	existing Stock Source :
		Source1:
			product1: 8
			product2: 3
		Source2:
			product1: 10
		Source3:
			product2: 10

	=> Expected Result:

		Cart:
			Delivery1:
				item1:product1 - 10
						sourced: Source1 -> 8 & Source2 -> 2
				item2:product2 - 5
						sourced: Source1 -> 3 & Source3 -> 2
			Delivery2:
				item3: product1 - 5
						sourced: Source2 -> 5

	*/
	t.Run("allocate easy", func(t *testing.T) {
		t.Parallel()

		stubbedProduct1 := productDomain.SimpleProduct{Identifier: "product1"}
		stubbedProduct2 := productDomain.SimpleProduct{Identifier: "product2"}

		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: stubbedProduct1,
							Item:    cart.Item{Qty: 10, ID: "item1"},
						},
						{
							Product: stubbedProduct2,
							Item:    cart.Item{Qty: 5, ID: "item2"},
						},
					},
				},
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: stubbedProduct1,
							Item:    cart.Item{Qty: 5, ID: "item3"},
						},
					},
				},
			},
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		source3 := domain.Source{LocationCode: "Source3"}
		stubbedSources := []domain.Source{source1, source2, source3}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"product1": 8,
					"product2": 3,
				},
				"Source2": {
					"product1": 10,
				},
				"Source3": {
					"product2": 10,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		itemAllocation, err := sourcingService.AllocateItems(context.Background(), &testCart)
		assert.NoError(t, err)
		assert.NoError(t, itemAllocation[domain.ItemID("item1")].Error)
		assert.Len(t, itemAllocation[domain.ItemID("item1")].AllocatedQtys[domain.ProductID(stubbedProduct1.GetIdentifier())], 2)

		assert.Equal(t, 8, itemAllocation[domain.ItemID("item1")].AllocatedQtys[domain.ProductID(stubbedProduct1.GetIdentifier())][source1])
		assert.Equal(t, 2, itemAllocation[domain.ItemID("item1")].AllocatedQtys[domain.ProductID(stubbedProduct1.GetIdentifier())][source2])
		assert.Equal(t, 3, itemAllocation[domain.ItemID("item2")].AllocatedQtys[domain.ProductID(stubbedProduct2.GetIdentifier())][source1])
		assert.Equal(t, 2, itemAllocation[domain.ItemID("item2")].AllocatedQtys[domain.ProductID(stubbedProduct2.GetIdentifier())][source3])
		assert.Equal(t, 5, itemAllocation[domain.ItemID("item3")].AllocatedQtys[domain.ProductID(stubbedProduct1.GetIdentifier())][source2])
	})

	t.Run("allocate cart with bundle item", func(t *testing.T) {
		t.Parallel()

		bundleProduct := productDomain.BundleProductWithActiveChoices{
			BundleProduct: productDomain.BundleProduct{
				Identifier: "bundle_product",
			},
			ActiveChoices: map[productDomain.Identifier]productDomain.ActiveChoice{
				"identifier1": {
					Qty: 1,
					Product: productDomain.SimpleProduct{
						Identifier: "product4",
					},
				},
				"identifier2": {
					Qty: 2,
					Product: productDomain.SimpleProduct{
						Identifier: "product3",
					},
				},
			},
		}

		simple2 := productDomain.SimpleProduct{Identifier: "product2"}
		simple5 := productDomain.SimpleProduct{Identifier: "product5"}

		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: bundleProduct,
							Item:    cart.Item{Qty: 1, ID: "item1"},
						},
						{
							Product: simple2,
							Item:    cart.Item{Qty: 1, ID: "item2"},
						},
					},
				},
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: simple5,
							Item:    cart.Item{Qty: 4, ID: "item3"},
						},
					},
				},
			},
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		source3 := domain.Source{LocationCode: "Source3"}
		stubbedSources := []domain.Source{source1, source2, source3}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"product2": 3,
					"product4": 27,
					"product5": 3,
				},
				"Source2": {
					"product1": 10,
					"product4": 27,
				},
				"Source3": {
					"product2": 10,
					"product3": 4,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		itemAllocation, err := sourcingService.AllocateItems(context.Background(), &testCart)
		assert.NoError(t, err)
		assert.NoError(t, itemAllocation[domain.ItemID("item1")].Error)
		assert.NoError(t, itemAllocation[domain.ItemID("item2")].Error)

		assert.ErrorIs(t, domain.ErrInsufficientSourceQty, itemAllocation[domain.ItemID("item3")].Error)

		assert.Equal(t, 1,
			itemAllocation[domain.ItemID("item1")].AllocatedQtys["product4"][source1])
		assert.Equal(t, 2,
			itemAllocation[domain.ItemID("item1")].AllocatedQtys["product3"][source3])
		assert.Equal(t, 1,
			itemAllocation[domain.ItemID("item2")].AllocatedQtys[domain.ProductID(simple2.GetIdentifier())][source1])
	})

	t.Run("if too many products are allocated to a bundle, they won't be available for the next item", func(t *testing.T) {
		t.Parallel()

		simpleInBundle1 := productDomain.SimpleProduct{Identifier: "product1"}
		simple2 := productDomain.SimpleProduct{Identifier: "product2"}

		bundleProduct := productDomain.BundleProductWithActiveChoices{
			BundleProduct: productDomain.BundleProduct{
				Identifier: "bundle_product",
			},
			ActiveChoices: map[productDomain.Identifier]productDomain.ActiveChoice{
				"identifier1": {
					Qty:     1,
					Product: simpleInBundle1,
				},
				"identifier2": {
					Qty:     2,
					Product: simple2,
				},
			},
		}

		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: bundleProduct,
							Item:    cart.Item{Qty: 2, ID: "item1"},
						},
						{
							Product: simple2,
							Item:    cart.Item{Qty: 1, ID: "item2"},
						},
					},
				},
			},
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		stubbedSources := []domain.Source{source1, source2}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"product1": 28,
					"product2": 4,
				},
				"Source2": {
					"product1": 28,
					"product2": 0,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		itemAllocation, err := sourcingService.AllocateItems(context.Background(), &testCart)
		assert.NoError(t, err)
		assert.NoError(t, itemAllocation[domain.ItemID("item1")].Error)

		assert.ErrorIs(t, domain.ErrInsufficientSourceQty, itemAllocation[domain.ItemID("item2")].Error)

		assert.Equal(t, 2,
			itemAllocation[domain.ItemID("item1")].AllocatedQtys["product1"][source1])
		assert.Equal(t, 4,
			itemAllocation[domain.ItemID("item1")].AllocatedQtys["product2"][source1])
	})

	t.Run("if an item from a bundle is purchased separately and insufficient quantity remains, it won't be available for the bundle", func(t *testing.T) {
		t.Parallel()

		simpleInBundle1 := productDomain.SimpleProduct{Identifier: "product1"}
		simple2 := productDomain.SimpleProduct{Identifier: "product2", Saleable: productDomain.Saleable{IsSaleable: true}}

		bundleProduct := productDomain.BundleProductWithActiveChoices{
			BundleProduct: productDomain.BundleProduct{
				Identifier: "bundle_product",
			},
			ActiveChoices: map[productDomain.Identifier]productDomain.ActiveChoice{
				"identifier1": {
					Qty:     1,
					Product: simpleInBundle1,
				},
				"identifier2": {
					Qty:     3,
					Product: simple2,
				},
			},
		}

		testCart := decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: simple2,
							Item:    cart.Item{Qty: 3, ID: "item1"},
						},
						{
							Product: bundleProduct,
							Item:    cart.Item{Qty: 1, ID: "item2"},
						},
					},
				},
			},
		}

		source1 := domain.Source{LocationCode: "Source1"}
		source2 := domain.Source{LocationCode: "Source2"}
		stubbedSources := []domain.Source{source1, source2}

		stockBySourceAndProductProviderMock := stockBySourceAndProductProviderMock{
			Qty: map[string]map[string]int{
				"Source1": {
					"product1": 28,
					"product2": 3,
				},
				"Source2": {
					"product1": 28,
					"product2": 0,
				},
			},
		}

		sourcingService := newDefaultSourcingService(stockBySourceAndProductProviderMock, stubbedSources)

		itemAllocation, err := sourcingService.AllocateItems(context.Background(), &testCart)
		assert.NoError(t, err)
		assert.NoError(t, itemAllocation[domain.ItemID("item1")].Error)

		assert.ErrorIs(t, domain.ErrInsufficientSourceQty, itemAllocation[domain.ItemID("item2")].Error)

		assert.Equal(t, 3,
			itemAllocation[domain.ItemID("item1")].AllocatedQtys[domain.ProductID(simple2.GetIdentifier())][source1])
	})
}

func newDefaultSourcingService(stockProvider domain.StockProvider, expectedSources []domain.Source) domain.DefaultSourcingService {
	sourcingService := domain.DefaultSourcingService{}
	availableSourcesProviderMock := availableSourcesProviderMock{Sources: expectedSources}

	sourcingService.Inject(flamingo.NullLogger{}, &struct {
		AvailableSourcesProvider domain.AvailableSourcesProvider `inject:",optional"`
		StockProvider            domain.StockProvider            `inject:",optional"`
	}{
		StockProvider:            stockProvider,
		AvailableSourcesProvider: availableSourcesProviderMock,
	})

	return sourcingService
}
