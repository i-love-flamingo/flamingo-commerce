package fake

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	commerceSourcingDomain "flamingo.me/flamingo-commerce/v3/sourcing/domain"
)

type (
	SourcingService struct {
		fakeSource *fakeSource
	}

	fakeSource struct {
		DeliveryCodes map[string]int
		Products      map[string]int
	}
)

var _ commerceSourcingDomain.SourcingService = &SourcingService{}

func (s *SourcingService) Inject(
	cfg *struct {
		FakeSourceData string `inject:"config:commerce.sourcing.fake.jsonPath,optional"`
	},
) {
	if cfg.FakeSourceData == "" {
		panic("fake sourcing service enabled but jsonPath was not set")
	}

	fileBytes, err := os.ReadFile(cfg.FakeSourceData)
	if err != nil {
		panic(err)
	}

	fakeSource := &fakeSource{}
	err = json.Unmarshal(fileBytes, fakeSource)
	if err != nil {
		panic(err)
	}

	s.fakeSource = fakeSource
}

func (s *SourcingService) AllocateItems(ctx context.Context, decoratedCart *decorator.DecoratedCart) (commerceSourcingDomain.ItemAllocations, error) {
	var (
		itemAllocations = commerceSourcingDomain.ItemAllocations{}
		overallError    error
	)
	for _, decoratedDelivery := range decoratedCart.DecoratedDeliveries {
		for _, decoratedItem := range decoratedDelivery.DecoratedItems {
			if decoratedItem.Item.ID == "" {
				continue
			}

			availableSources, err := s.GetAvailableSources(ctx, decoratedItem.Product, &decoratedDelivery.Delivery.DeliveryInfo, decoratedCart)
			availableSourcesForProduct := availableSources[commerceSourcingDomain.ProductID(decoratedItem.Product.GetIdentifier())]

			itemAllocation := commerceSourcingDomain.ItemAllocation{
				AllocatedQtys: map[commerceSourcingDomain.ProductID]commerceSourcingDomain.AllocatedQtys{
					commerceSourcingDomain.ProductID(decoratedItem.Product.GetIdentifier()): commerceSourcingDomain.AllocatedQtys(availableSourcesForProduct),
				},
				Error: err,
			}

			itemAllocations[commerceSourcingDomain.ItemID(decoratedItem.Item.ID)] = itemAllocation
		}
	}

	return itemAllocations, overallError
}

func (s *SourcingService) GetAvailableSources(_ context.Context, product productDomain.BasicProduct, deliveryInfo *cartDomain.DeliveryInfo, _ *decorator.DecoratedCart) (commerceSourcingDomain.AvailableSourcesPerProduct, error) {
	if product.Type() == productDomain.TypeBundle || product.Type() == productDomain.TypeConfigurable {
		return nil, commerceSourcingDomain.ErrUnsupportedProductType
	}

	productId := product.GetIdentifier()

	if productId == "" {
		return nil, commerceSourcingDomain.ErrEmptyProductIdentifier
	}

	allocationQty, found := s.fakeSource.Products[productId]
	if found {
		return commerceSourcingDomain.AvailableSourcesPerProduct{
			commerceSourcingDomain.ProductID(productId): s.makeAvailableSources(allocationQty, productId),
		}, nil
	}

	if deliveryInfo.Code != "" {
		allocationQty, found := s.fakeSource.DeliveryCodes[deliveryInfo.Code]
		if found {
			return commerceSourcingDomain.AvailableSourcesPerProduct{
				commerceSourcingDomain.ProductID(productId): s.makeAvailableSources(allocationQty, deliveryInfo.Code),
			}, nil
		}
	}

	return nil, commerceSourcingDomain.ErrNoSourceAvailable
}

func (s *SourcingService) makeAvailableSources(qty int, id string) commerceSourcingDomain.AvailableSources {
	return map[commerceSourcingDomain.Source]int{
		commerceSourcingDomain.Source{
			LocationCode:         id,
			ExternalLocationCode: id,
		}: qty,
	}
}
