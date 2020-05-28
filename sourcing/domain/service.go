package domain

import (
	"context"
	"flamingo.me/flamingo/v3/framework/flamingo"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/product/domain"

	"github.com/pkg/errors"
)

type (

	// SourcingService interface
	SourcingService interface {
		//AllocateItems returns Sources for the given item in the given cart
		// e.g. use this during place order to know
		// throws ErrInsufficientSourceQty if not enough stock is available for the amount of items in the cart
		// throws ErrNoSourceAvailable if no source is available at all for one of the items
		// throws ErrNeedMoreDetailsSourceCannotBeDetected  is informations on the cart (or delivery is missing)
		AllocateItems(ctx context.Context, decoratedCart *decorator.DecoratedCart) (ItemAllocations, error)

		// GetAvailableSources returns possible Sources for the product and the desired delivery.
		// Optional the existing cart can be passed so that existing items in the cart can be evaluated also (e.g. deduct stock)
		// e.g. use this before a product should be placed in the cart to know if and from where an item can be sourced
		// throws ErrNeedMoreDetailsSourceCannotBeDetected
		// throws ErrNoSourceAvailable if no source is available for the product and the given delivery
		GetAvailableSources(ctx context.Context, product domain.BasicProduct, deliveryInfo *cartDomain.DeliveryInfo, decoratedCart *decorator.DecoratedCart) (AvailableSources, error)
	}

	ItemID string
	//ItemAllocations represents the allocated Qtys per itemId
	ItemAllocations map[ItemID]AllocatedQtys

	//AllocatedQtys represents the allocated Qty per source
	AllocatedQtys map[Source]int

	Source struct {
		// LocationCode identifies the warehouse or stock location
		LocationCode string
		// ExternalLocationCode identifies the source location in an external system
		ExternalLocationCode string
	}

	// AvailableSources is the result value object containing the available Qty per Source
	AvailableSources map[Source]int

	//DefaultSourcingService - an example implementation
	DefaultSourcingService struct {
		availableSourcesProvider AvailableSourcesProvider
		stockProvider            StockProvider
		logger                   flamingo.Logger
	}

	AvailableSourcesProvider interface {
		GetPossibleSources(ctx context.Context, product domain.BasicProduct, deliveryInfo *cartDomain.DeliveryInfo) ([]Source, error)
	}

	StockProvider interface {
		GetStock(ctx context.Context, product domain.BasicProduct, source Source) (int, error)
	}
)

var (
	_ SourcingService = new(DefaultSourcingService)

	// ErrInsufficientSourceQty - use to indicate that the requested qty exceeds the available qty
	ErrInsufficientSourceQty = errors.New("Available Source Qty insufficient")

	// ErrNoSourceAvailable - use to indicate that no source for item is available at all
	ErrNoSourceAvailable = errors.New("No Available Source Qty")

	// ErrNeedMoreDetailsSourceCannotBeDetected - use to indicate that informations are missing to determine a source
	ErrNeedMoreDetailsSourceCannotBeDetected = errors.New("Source cannot be detected")
)

//Inject the dependencies
func (d *DefaultSourcingService) Inject(logger flamingo.Logger, dep *struct {
	AvailableSourcesProvider AvailableSourcesProvider `inject:",optional"`
	StockProvider            StockProvider            `inject:",optional"`
}) *DefaultSourcingService {
	if dep != nil {
		d.availableSourcesProvider = dep.AvailableSourcesProvider
		d.stockProvider = dep.StockProvider
	}
	d.logger = logger.WithField(flamingo.LogKeyModule, "sourcing").WithField(flamingo.LogKeyCategory, "DefaultSourcingService")
	return d
}

//GetAvailableSources
func (d *DefaultSourcingService) GetAvailableSources(ctx context.Context, product domain.BasicProduct, deliveryInfo *cartDomain.DeliveryInfo, decoratedCart *decorator.DecoratedCart) (AvailableSources, error) {
	availableSources := make(AvailableSources)
	if d.availableSourcesProvider == nil {
		d.logger.Error("no Source Provider bound")
		return nil, errors.New("no Source Provider bound")
	}
	if d.stockProvider == nil {
		d.logger.Error("no Stock Provider bound")
		return nil, errors.New("no Stock Provider bound")
	}

	sources, err := d.availableSourcesProvider.GetPossibleSources(ctx, product, deliveryInfo)
	if err != nil {
		return nil, err
	}

	var lastStockError error
	for _, source := range sources {
		qty, err := d.stockProvider.GetStock(ctx, product, source)
		if err != nil {
			d.logger.Error(err)
			lastStockError = err
			continue
		}
		if qty > 0 {
			availableSources[source] = qty
		}
	}

	// if a cart is given we need to deduct the possible allocated items in the cart
	if decoratedCart != nil {
		allocatedSources, err := d.AllocateItems(ctx, decoratedCart)
		if err != nil {
			return nil, err
		}
		itemIdsWithProduct := getItemIdsWithProduct(decoratedCart, product)
		for _, itemId := range itemIdsWithProduct {
			availableSources = availableSources.Reduce(allocatedSources[itemId])
		}
	}

	if len(availableSources) == 0 && lastStockError != nil {
		return availableSources, errors.Wrap(ErrNoSourceAvailable, lastStockError.Error())
	} else if len(availableSources) == 0 {
		return availableSources, ErrNoSourceAvailable
	}
	return availableSources, nil
}

func getItemIdsWithProduct(dc *decorator.DecoratedCart, product domain.BasicProduct) []ItemID {
	var result []ItemID
	for _, di := range dc.GetAllDecoratedItems() {
		if di.Product.GetIdentifier() == product.GetIdentifier() {
			result = append(result, ItemID(di.Item.ID))
		}
	}
	return result
}

func (d *DefaultSourcingService) AllocateItems(ctx context.Context, decoratedCart *decorator.DecoratedCart) (ItemAllocations, error) {
	if decoratedCart == nil {
		return nil, errors.New("Cart not given")
	}
	if d.availableSourcesProvider == nil {
		d.logger.Error("no Source Provider bound")
		return nil, errors.New("no Source Provider bound")
	}
	if d.stockProvider == nil {
		d.logger.Error("no Stock Provider bound")
		return nil, errors.New("no Stock Provider bound")
	}

	for _, delivery := range decoratedCart.DecoratedDeliveries {
		for _, _ = range delivery.DecoratedItems {

		}

	}
	return nil, nil
}

/*
// MainLocation returns first sourced location (or empty string)
func (s Sources) MainLocation() string {
	if len(s) < 1 {
		return ""
	}
	return s[0].LocationCode
}

// Next - returns the next source and the remaining, or error if nothing remains
func (s Sources) Next() (Source, Sources, error) {
	if s.QtySum() < 1 {
		return Source{}, Sources{}, ErrInsufficientSourceQty
	}
	for _, source := range s {
		if source.Qty > 0 {
			usedSource := Source{
				Qty:          1,
				LocationCode: source.LocationCode,
			}
			usedSources := Sources{usedSource}
			return usedSource, s.Reduce(usedSources), nil
		}
	}
	return Source{}, Sources{}, ErrInsufficientSourceQty
}

// QtySum returns the sum of all sourced items
func (s Sources) QtySum() int {
	qty := int(0)
	for _, source := range s {
		if source.Qty == math.MaxInt64 {
			return math.MaxInt64
		}
		qty = qty + source.Qty
	}
	return qty
}


*/

// Reduce returns new AvailableSources reduced by the given AvailableSources
func (s AvailableSources) Reduce(reduceby AllocatedQtys) AvailableSources {
	newAvailableSources := make(AvailableSources)
	for source, availableQty := range s {
		if allocated, ok := reduceby[source]; ok {
			newQty := availableQty - allocated
			if newQty > 0 {
				newAvailableSources[source] = newQty
			}
		} else {
			newAvailableSources[source] = availableQty
		}
	}
	return s
}
