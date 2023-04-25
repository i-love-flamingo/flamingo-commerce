package fake_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	commerceSourcingDomain "flamingo.me/flamingo-commerce/v3/sourcing/domain"
	"flamingo.me/flamingo-commerce/v3/sourcing/domain/fake"
)

func TestSourcingService_AllocateItems(t *testing.T) {
	t.Parallel()

	t.Run("success when product id is correct", func(t *testing.T) {
		t.Parallel()

		service := &fake.SourcingService{}
		service.Inject(&struct {
			FakeSourceData string `inject:"config:commerce.sourcing.fake.jsonPath,optional"`
		}{
			FakeSourceData: "testdata/fakeSourceData.json",
		})

		decoratedCart := &decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					Delivery: cartDomain.Delivery{
						DeliveryInfo: cartDomain.DeliveryInfo{
							Code: "inflight",
						},
					},
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Item: cartDomain.Item{
								ID: "item1",
							},
							Product: productDomain.SimpleProduct{
								Identifier: "0f0asdf-0asd0a9sd-askdlj123rw",
							},
						},
						{
							Item: cartDomain.Item{
								ID: "item2",
							},
							Product: productDomain.SimpleProduct{
								Identifier: "0f0asdf-0asd0a9sd-askdlj123rx",
							},
						},
					},
				},
			},
		}

		expectedItemAllocations := commerceSourcingDomain.ItemAllocations{
			"item1": {
				AllocatedQtys: map[commerceSourcingDomain.ProductID]commerceSourcingDomain.AllocatedQtys{
					"0f0asdf-0asd0a9sd-askdlj123rw": {
						commerceSourcingDomain.Source{
							LocationCode:         "0f0asdf-0asd0a9sd-askdlj123rw",
							ExternalLocationCode: "0f0asdf-0asd0a9sd-askdlj123rw",
						}: 10,
					},
				},
				Error: nil,
			},
			"item2": {
				AllocatedQtys: map[commerceSourcingDomain.ProductID]commerceSourcingDomain.AllocatedQtys{
					"0f0asdf-0asd0a9sd-askdlj123rx": {
						commerceSourcingDomain.Source{
							LocationCode:         "0f0asdf-0asd0a9sd-askdlj123rx",
							ExternalLocationCode: "0f0asdf-0asd0a9sd-askdlj123rx",
						}: 15,
					},
				},
				Error: nil,
			},
		}

		resultAllocations, err := service.AllocateItems(context.Background(), decoratedCart)
		assert.NoError(t, err)
		assert.NotNil(t, resultAllocations)
		assert.Equal(t, expectedItemAllocations, resultAllocations)
	})

	t.Run("empty result when there are no item ids", func(t *testing.T) {
		t.Parallel()

		service := &fake.SourcingService{}
		service.Inject(&struct {
			FakeSourceData string `inject:"config:commerce.sourcing.fake.jsonPath,optional"`
		}{
			FakeSourceData: "testdata/fakeSourceData.json",
		})

		decoratedCart := &decorator.DecoratedCart{
			DecoratedDeliveries: []decorator.DecoratedDelivery{
				{
					Delivery: cartDomain.Delivery{
						DeliveryInfo: cartDomain.DeliveryInfo{
							Code: "inflight",
						},
					},
					DecoratedItems: []decorator.DecoratedCartItem{
						{
							Product: productDomain.SimpleProduct{
								Identifier: "0f0asdf-0asd0a9sd-askdlj123rw",
							},
						},
						{
							Product: productDomain.SimpleProduct{
								Identifier: "0f0asdf-0asd0a9sd-askdlj123rx",
							},
						},
					},
				},
			},
		}

		expectedItemAllocations := commerceSourcingDomain.ItemAllocations{}

		resultAllocations, err := service.AllocateItems(context.Background(), decoratedCart)
		assert.NoError(t, err)
		assert.NotNil(t, resultAllocations)
		assert.Equal(t, expectedItemAllocations, resultAllocations)
	})
}

func TestSourcingService_GetAvailableSources(t *testing.T) {
	t.Parallel()

	t.Run("success when product id exists", func(t *testing.T) {
		t.Parallel()

		service := &fake.SourcingService{}
		service.Inject(&struct {
			FakeSourceData string `inject:"config:commerce.sourcing.fake.jsonPath,optional"`
		}{
			FakeSourceData: "testdata/fakeSourceData.json",
		})

		product := productDomain.SimpleProduct{
			Identifier: "0f0asdf-0asd0a9sd-askdlj123rw",
		}

		deliveryInfo := &cartDomain.DeliveryInfo{
			Code: "inflight",
		}

		expectedSources := commerceSourcingDomain.AvailableSourcesPerProduct{
			"0f0asdf-0asd0a9sd-askdlj123rw": {
				commerceSourcingDomain.Source{
					LocationCode:         "0f0asdf-0asd0a9sd-askdlj123rw",
					ExternalLocationCode: "0f0asdf-0asd0a9sd-askdlj123rw",
				}: 10,
			},
		}

		resultSources, err := service.GetAvailableSources(context.Background(), product, deliveryInfo, &decorator.DecoratedCart{})

		assert.NoError(t, err)
		assert.NotNil(t, resultSources)
		assert.Equal(t, expectedSources, resultSources)
	})

	t.Run("success when product id was not found but delivery code correct", func(t *testing.T) {
		t.Parallel()

		service := &fake.SourcingService{}
		service.Inject(&struct {
			FakeSourceData string `inject:"config:commerce.sourcing.fake.jsonPath,optional"`
		}{
			FakeSourceData: "testdata/fakeSourceData.json",
		})

		product := productDomain.SimpleProduct{
			Identifier: "some-fake-id",
		}

		deliveryInfo := &cartDomain.DeliveryInfo{
			Code: "inflight",
		}

		expectedSources := commerceSourcingDomain.AvailableSourcesPerProduct{
			"some-fake-id": {
				commerceSourcingDomain.Source{
					LocationCode:         "inflight",
					ExternalLocationCode: "inflight",
				}: 5,
			},
		}

		resultSources, err := service.GetAvailableSources(context.Background(), product, deliveryInfo, &decorator.DecoratedCart{})

		assert.NoError(t, err)
		assert.NotNil(t, resultSources)
		assert.Equal(t, expectedSources, resultSources)
	})

	t.Run("failure when product id is empty", func(t *testing.T) {
		t.Parallel()

		service := &fake.SourcingService{}
		service.Inject(&struct {
			FakeSourceData string `inject:"config:commerce.sourcing.fake.jsonPath,optional"`
		}{
			FakeSourceData: "testdata/fakeSourceData.json",
		})

		product := productDomain.SimpleProduct{
			Identifier: "",
		}

		deliveryInfo := &cartDomain.DeliveryInfo{
			Code: "inflight",
		}

		resultSources, err := service.GetAvailableSources(context.Background(), product, deliveryInfo, &decorator.DecoratedCart{})

		assert.Error(t, err)
		assert.Nil(t, resultSources)
		assert.Equal(t, commerceSourcingDomain.ErrEmptyProductIdentifier, err)
	})

	t.Run("failure when product type is bundle", func(t *testing.T) {
		t.Parallel()

		service := &fake.SourcingService{}
		service.Inject(&struct {
			FakeSourceData string `inject:"config:commerce.sourcing.fake.jsonPath,optional"`
		}{
			FakeSourceData: "testdata/fakeSourceData.json",
		})

		product := productDomain.BundleProduct{}

		deliveryInfo := &cartDomain.DeliveryInfo{
			Code: "inflight",
		}

		resultSources, err := service.GetAvailableSources(context.Background(), product, deliveryInfo, &decorator.DecoratedCart{})

		assert.Error(t, err)
		assert.Nil(t, resultSources)
		assert.Equal(t, commerceSourcingDomain.ErrUnsupportedProductType, err)
	})
}
