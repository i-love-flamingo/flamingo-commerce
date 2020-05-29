package domain_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/sourcing/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
)

type availableSourcesProviderMock struct {
	Sources []domain.Source
	Error   error
}

type stockProviderMock struct {
	Qty   int
	Error error
}

var _ domain.AvailableSourcesProvider = new(availableSourcesProviderMock)
var _ domain.StockProvider = new(stockProviderMock)
var _ domain.AvailableSourcesProvider = new(stockBySourceAndProductProviderMock)

func (a availableSourcesProviderMock) GetPossibleSources(ctx context.Context, product productDomain.BasicProduct, deliveryInfo *cart.DeliveryInfo) ([]domain.Source, error) {
	return a.Sources, a.Error
}

func (s stockProviderMock) GetStock(ctx context.Context, product productDomain.BasicProduct, source domain.Source) (int, error) {
	return s.Qty, s.Error
}

type stockBySourceAndProductProviderMock struct {
	// [source.LocationCode][product.Identifier] = Qty
	Qty   map[string]map[string]int
	Error error
}

func (s stockBySourceAndProductProviderMock) GetStock(ctx context.Context, product productDomain.BasicProduct, source domain.Source) (int, error) {
	return s.Qty[source.LocationCode][product.GetIdentifier()], s.Error
}

func (s stockBySourceAndProductProviderMock) GetPossibleSources(ctx context.Context, product productDomain.BasicProduct, deliveryInfo *cart.DeliveryInfo) ([]domain.Source, error) {
	panic("implement me")
}

func TestDefaultSourcingService_GetAvailableSources(t *testing.T) {
	t.Run("full qty with nil cart", func(t *testing.T) {
		stubbedSources := []domain.Source{{LocationCode: "loc1"}}
		stubbedStockQty := 10

		stockProviderMock := stockProviderMock{Qty: stubbedStockQty}
		sourcingService := newDefaultSourcingService(stockProviderMock, stubbedSources)

		sources, err := sourcingService.GetAvailableSources(context.Background(), productDomain.SimpleProduct{}, nil, nil)
		assert.NoError(t, err)

		expectedSources := domain.AvailableSources{stubbedSources[0]: stubbedStockQty}
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

		expectedSources := domain.AvailableSources{stubbedSources[0]: stubbedStockQty - stubbedQtyAlreadyInCart}
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
		assert.Equal(t, err, domain.ErrNoSourceAvailable)
	})
}

func TestDefaultSourcingService_AllocateItems(t *testing.T) {
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
		assert.Len(t, itemAllocation[domain.ItemID("item1")], 2)
		assert.Equal(t, 8, itemAllocation[domain.ItemID("item1")][source1])
		assert.Equal(t, 2, itemAllocation[domain.ItemID("item1")][source2])
		assert.Equal(t, 3, itemAllocation[domain.ItemID("item2")][source1])
		assert.Equal(t, 2, itemAllocation[domain.ItemID("item2")][source3])
		assert.Equal(t, 5, itemAllocation[domain.ItemID("item3")][source2])
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
