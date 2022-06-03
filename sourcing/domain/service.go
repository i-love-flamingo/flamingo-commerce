package domain

import (
	"context"
	"math"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/product/domain"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"github.com/pkg/errors"
)

type (
	// SourcingService describes the main port used by the sourcing logic.
	SourcingService interface {
		// AllocateItems returns Sources for the given item in the given cart
		// e.g. use this during place order to know
		// throws ErrInsufficientSourceQty if not enough stock is available for the amount of items in the cart
		// throws ErrNoSourceAvailable if no source is available at all for one of the items
		// throws ErrNeedMoreDetailsSourceCannotBeDetected if information on the cart (or delivery is missing)
		AllocateItems(ctx context.Context, decoratedCart *decorator.DecoratedCart) (ItemAllocations, error)

		// GetAvailableSources returns possible Sources for the product and the desired delivery.
		// Optional the existing cart can be passed so that existing items in the cart can be evaluated also (e.g. deduct stock)
		// e.g. use this before a product should be placed in the cart to know if and from where an item can be sourced
		// throws ErrNeedMoreDetailsSourceCannotBeDetected
		// throws ErrNoSourceAvailable if no source is available for the product and the given delivery
		GetAvailableSources(ctx context.Context, product domain.BasicProduct, deliveryInfo *cartDomain.DeliveryInfo, decoratedCart *decorator.DecoratedCart) (AvailableSources, error)
	}

	// ItemID string alias
	ItemID string

	// ItemAllocations represents the allocated Qtys per itemID
	ItemAllocations map[ItemID]ItemAllocation

	// ItemAllocation info
	ItemAllocation struct {
		AllocatedQtys AllocatedQtys
		Error         error
	}

	// AllocatedQtys represents the allocated Qty per source
	AllocatedQtys map[Source]int

	// Source descriptor for a single location
	Source struct {
		// LocationCode identifies the warehouse or stock location
		LocationCode string
		// ExternalLocationCode identifies the source location in an external system
		ExternalLocationCode string
	}

	// AvailableSources is the result value object containing the available Qty per Source
	AvailableSources map[Source]int

	// DefaultSourcingService provides a default implementation of the SourcingService interface.
	// This default implementation is used unless a project overrides the interface binding.
	DefaultSourcingService struct {
		availableSourcesProvider AvailableSourcesProvider
		stockProvider            StockProvider
		logger                   flamingo.Logger
	}

	// AvailableSourcesProvider interface for DefaultSourcingService
	AvailableSourcesProvider interface {
		GetPossibleSources(ctx context.Context, product domain.BasicProduct, deliveryInfo *cartDomain.DeliveryInfo) ([]Source, error)
	}

	// StockProvider interface for DefaultSourcingService
	StockProvider interface {
		GetStock(ctx context.Context, product domain.BasicProduct, source Source, deliveryInfo *cartDomain.DeliveryInfo) (int, error)
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

// Inject the dependencies
func (d *DefaultSourcingService) Inject(
	logger flamingo.Logger,
	dep *struct {
		AvailableSourcesProvider AvailableSourcesProvider `inject:",optional"`
		StockProvider            StockProvider            `inject:",optional"`
	},
) *DefaultSourcingService {
	d.logger = logger.WithField(flamingo.LogKeyModule, "sourcing").WithField(flamingo.LogKeyCategory, "DefaultSourcingService")

	if dep != nil {
		d.availableSourcesProvider = dep.AvailableSourcesProvider
		d.stockProvider = dep.StockProvider
	}

	return d
}

// GetAvailableSources - see description in Interface
func (d *DefaultSourcingService) GetAvailableSources(ctx context.Context, product domain.BasicProduct, deliveryInfo *cartDomain.DeliveryInfo, decoratedCart *decorator.DecoratedCart) (AvailableSources, error) {
	if err := d.checkConfiguration(); err != nil {
		return nil, err
	}

	sources, err := d.availableSourcesProvider.GetPossibleSources(ctx, product, deliveryInfo)
	if err != nil {
		return nil, err
	}

	var lastStockError error
	availableSources := AvailableSources{}
	for _, source := range sources {
		qty, err := d.stockProvider.GetStock(ctx, product, source, deliveryInfo)
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

		for _, itemID := range itemIdsWithProduct {
			availableSources = availableSources.Reduce(allocatedSources[itemID].AllocatedQtys)
		}
	}

	if len(availableSources) == 0 {
		if lastStockError != nil {
			return availableSources, errors.Wrap(ErrNoSourceAvailable, lastStockError.Error())
		}
		return availableSources, ErrNoSourceAvailable
	}

	return availableSources, nil
}

// AllocateItems - see description in Interface
func (d *DefaultSourcingService) AllocateItems(ctx context.Context, decoratedCart *decorator.DecoratedCart) (ItemAllocations, error) {
	decoratedCart.DecoratedDeliveries[0].DecoratedItems[0].Product.GetIdentifier()
	if err := d.checkConfiguration(); err != nil {
		return nil, err
	}
	if decoratedCart == nil {
		return nil, errors.New("Cart not given")
	}

	// productSourcestock holds the available stock per product and source.
	// During allocation the initial retrieved available stock is reduced according to used allocation
	var productSourcestock = map[string]map[Source]int{}

	if len(decoratedCart.DecoratedDeliveries) == 0 {
		return nil, ErrNeedMoreDetailsSourceCannotBeDetected
	}

	resultItemAllocations := ItemAllocations{}

	// overallError that will be returned
	var overallError error
	for _, delivery := range decoratedCart.DecoratedDeliveries {
		for _, decoratedItem := range delivery.DecoratedItems {
			var itemAllocation ItemAllocation
			itemAllocation, productSourcestock = d.allocateItem(ctx, productSourcestock, decoratedItem, delivery.Delivery.DeliveryInfo)
			resultItemAllocations[ItemID(decoratedItem.Item.ID)] = itemAllocation
		}
	}

	return resultItemAllocations, overallError
}

func (d *DefaultSourcingService) checkConfiguration() error {
	if d.availableSourcesProvider == nil {
		d.logger.Error("no Source Provider bound")
		return errors.New("no Source Provider bound")
	}
	if d.stockProvider == nil {
		d.logger.Error("no Stock Provider bound")
		return errors.New("no Stock Provider bound")
	}

	return nil
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

// allocateItem returns the itemAllocation and the remaining stock for the given item.
// The passed productSourcestock is used - and the remaining productSourcestock is returned. In case a source is not yet given in productSourcestock it will be fetched
func (d *DefaultSourcingService) allocateItem(ctx context.Context, productSourcestock map[string]map[Source]int, decoratedItem decorator.DecoratedCartItem, deliveryInfo cartDomain.DeliveryInfo) (ItemAllocation, map[string]map[Source]int) {
	var resultItemAllocation = ItemAllocation{
		AllocatedQtys: make(AllocatedQtys),
	}
	// copy given known stock
	remainingSourcestock := productSourcestock

	productID := decoratedItem.Product.GetIdentifier()
	if productID == "" {
		return ItemAllocation{
			Error: errors.New("product id missing"),
		}, remainingSourcestock
	}
	sources, err := d.availableSourcesProvider.GetPossibleSources(ctx, decoratedItem.Product, &deliveryInfo)
	if err != nil {
		return ItemAllocation{
			Error: err,
		}, remainingSourcestock
	}
	if len(sources) == 0 {
		return ItemAllocation{
			Error: ErrNoSourceAvailable,
		}, remainingSourcestock
	}

	qtyToAllocate := decoratedItem.Item.Qty
	allocatedQty := 0

	if _, exists := productSourcestock[productID]; !exists {
		productSourcestock[productID] = make(map[Source]int)
	}

	for _, source := range sources {
		// if we have no stock given for source and productid we fetch it initially
		if _, exists := remainingSourcestock[productID][source]; !exists {
			sourceStock, err := d.stockProvider.GetStock(ctx, decoratedItem.Product, source, &deliveryInfo)
			if err != nil {
				d.logger.Error(err)
				continue
			}
			remainingSourcestock[productID][source] = sourceStock
		}

		if remainingSourcestock[productID][source] == 0 {
			continue
		}
		if allocatedQty < qtyToAllocate {
			// stock to write to result allocation is the lowest of either :
			// - the remaining qty that is to be allocated
			// OR
			// - the existing sourceStock that is then used completely
			stockToAllocate := min(qtyToAllocate-allocatedQty, productSourcestock[productID][source])

			resultItemAllocation.AllocatedQtys[source] = stockToAllocate

			// increment allocatedQty by allocated Stock
			allocatedQty = allocatedQty + stockToAllocate

			// decrement remaining productSourceStock accordingly as its not happening by itself
			remainingSourcestock[productID][source] = remainingSourcestock[productID][source] - stockToAllocate
		}
	}

	if allocatedQty < qtyToAllocate {
		resultItemAllocation.Error = ErrInsufficientSourceQty
	}
	return resultItemAllocation, remainingSourcestock
}

// QtySum returns the sum of all sourced items
func (s AvailableSources) QtySum() int {
	qty := 0
	for _, sqty := range s {
		if sqty == math.MaxInt64 {
			return sqty
		}
		qty = qty + sqty
	}
	return qty
}

// Reduce returns new AvailableSources reduced by the given AvailableSources
func (s AvailableSources) Reduce(reducedBy AllocatedQtys) AvailableSources {
	newAvailableSources := make(AvailableSources)
	for source, availableQty := range s {
		if allocated, ok := reducedBy[source]; ok {
			newQty := availableQty - allocated
			if newQty > 0 {
				newAvailableSources[source] = newQty
			}
		} else {
			newAvailableSources[source] = availableQty
		}
	}
	return newAvailableSources
}

// min returns minimum of 2 ints
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
