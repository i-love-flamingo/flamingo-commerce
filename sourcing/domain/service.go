package domain

import (
	"context"
	"fmt"
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
		CustomAllocations map[ProductID]AllocatedQtys
		AllocatedQtys     AllocatedQtys
		Error             error
	}

	ProductID string

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

	// ErrNeedMoreDetailsSourceCannotBeDetected - use to indicate that information are missing to determine a source
	ErrNeedMoreDetailsSourceCannotBeDetected = errors.New("Source cannot be detected")

	// ErrUnsupportedProduct return when product type is not supported by service
	ErrUnsupportedProduct = errors.New("unsupported product")

	// ErrProductIDMissing return when product id is missing
	ErrProductIDMissing = errors.New("product id missing")
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
func (d *DefaultSourcingService) GetAvailableSources(
	ctx context.Context,
	product domain.BasicProduct,
	deliveryInfo *cartDomain.DeliveryInfo,
	decoratedCart *decorator.DecoratedCart,
) (AvailableSources, error) {
	if err := d.checkConfiguration(); err != nil {
		return nil, err
	}

	sources, err := d.availableSourcesProvider.GetPossibleSources(ctx, product, deliveryInfo)
	if err != nil {
		return nil, err
	}

	availableSources, lastStockError := d.fetchAvailableSources(ctx, product, deliveryInfo, sources)
	if err != nil {
		return nil, err
	}

	// if a cart is given we need to deduct the possible allocated items in the cart
	if decoratedCart != nil {
		availableSources, err = d.deductAllocatedItems(ctx, availableSources, decoratedCart, product)
		if err != nil {
			return nil, err
		}
	}

	if len(availableSources) == 0 {
		if lastStockError != nil {
			return availableSources, errors.Wrap(ErrNoSourceAvailable, lastStockError.Error())
		}

		return availableSources, fmt.Errorf("%w %s", ErrNoSourceAvailable, formatSources(sources))
	}

	return availableSources, nil
}

func (d *DefaultSourcingService) fetchAvailableSources(
	ctx context.Context,
	product domain.BasicProduct,
	deliveryInfo *cartDomain.DeliveryInfo,
	sources []Source) (AvailableSources, error) {
	availableSources := make(AvailableSources)
	var lastStockError error

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

	return availableSources, lastStockError
}

func (d *DefaultSourcingService) deductAllocatedItems(ctx context.Context, availableSources AvailableSources, decoratedCart *decorator.DecoratedCart, product domain.BasicProduct) (AvailableSources, error) {
	allocatedSources, err := d.AllocateItems(ctx, decoratedCart)
	if err != nil {
		return nil, err
	}

	itemIdsWithProduct := getItemIdsWithProduct(decoratedCart, product)

	for _, itemID := range itemIdsWithProduct {
		availableSources = availableSources.Reduce(allocatedSources[itemID].AllocatedQtys)
	}

	return availableSources, nil
}

func (d *DefaultSourcingService) AllocateItems(ctx context.Context, decoratedCart *decorator.DecoratedCart) (ItemAllocations, error) {
	if err := d.checkConfiguration(); err != nil {
		return nil, err
	}
	if decoratedCart == nil {
		return nil, errors.New("Cart not given")
	}

	productSourcestock := make(map[string]map[Source]int)

	if len(decoratedCart.DecoratedDeliveries) == 0 {
		return nil, ErrNeedMoreDetailsSourceCannotBeDetected
	}

	resultItemAllocations := make(ItemAllocations)

	for _, delivery := range decoratedCart.DecoratedDeliveries {
		deliveryInfo := delivery.Delivery.DeliveryInfo // create a new variable to avoid memory aliasing

		for _, decoratedItem := range delivery.DecoratedItems {
			item := decoratedItem // create a new variable to avoid memory aliasing

			itemAllocation, updatedSourcestock, err := d.allocateItemForProduct(ctx, productSourcestock, &item, deliveryInfo)
			if err != nil {
				return nil, err
			}

			resultItemAllocations[ItemID(item.Item.ID)] = itemAllocation
			productSourcestock = updatedSourcestock
		}
	}

	return resultItemAllocations, nil
}

func (d *DefaultSourcingService) allocateItemForProduct(
	ctx context.Context,
	productSourcestock map[string]map[Source]int,
	decoratedItem *decorator.DecoratedCartItem,
	deliveryInfo cartDomain.DeliveryInfo,
) (ItemAllocation, map[string]map[Source]int, error) {
	if decoratedItem.Product.Type() == domain.TypeBundle || decoratedItem.Product.Type() == domain.TypeConfigurableWithActiveVariant {
		return ItemAllocation{}, productSourcestock, ErrUnsupportedProduct
	}

	if bundleProduct, ok := decoratedItem.Product.(domain.BundleProductWithActiveChoices); ok {
		itemAllocation, productStock := d.allocateBundleWithActiveChoices(ctx, decoratedItem.Item.Qty, productSourcestock, bundleProduct, deliveryInfo)
		return itemAllocation, productStock, nil
	}

	allocatedQtys, updatedSourcestock, err := d.allocateItem(ctx, productSourcestock, decoratedItem.Product, decoratedItem.Item.Qty, deliveryInfo)

	itemAllocation := ItemAllocation{
		AllocatedQtys: allocatedQtys,
		Error:         err,
	}

	return itemAllocation, updatedSourcestock, nil
}

func (d *DefaultSourcingService) allocateBundleWithActiveChoices(
	ctx context.Context,
	itemQty int,
	productSourcestock map[string]map[Source]int,
	bundleProduct domain.BundleProductWithActiveChoices,
	deliveryInfo cartDomain.DeliveryInfo,
) (ItemAllocation, map[string]map[Source]int) {
	var resultItemAllocation ItemAllocation
	// copy given known stock
	remainingSourcestock := productSourcestock

	for choiceID, activeChoice := range bundleProduct.ActiveChoices {
		qty := activeChoice.Qty * itemQty

		allocatedQtys, productSourceStockForBundle, err := d.allocateItem(ctx, productSourcestock, activeChoice.Product, qty, deliveryInfo)

		if resultItemAllocation.CustomAllocations == nil {
			resultItemAllocation.CustomAllocations = make(map[ProductID]AllocatedQtys)
		}

		if err != nil {
			resultItemAllocation.Error = err
		}

		remainingSourcestock = productSourceStockForBundle
		resultItemAllocation.CustomAllocations[ProductID(choiceID)] = allocatedQtys
	}

	return resultItemAllocation, remainingSourcestock
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

func (d *DefaultSourcingService) allocateItem(
	ctx context.Context,
	productSourcestock map[string]map[Source]int,
	product domain.BasicProduct,
	qtyToAllocate int,
	deliveryInfo cartDomain.DeliveryInfo,
) (AllocatedQtys, map[string]map[Source]int, error) {
	sources, err := d.availableSourcesProvider.GetPossibleSources(ctx, product, &deliveryInfo)
	if err != nil {
		return nil, productSourcestock, fmt.Errorf("error getting possible sources: %w", err)
	}

	if len(sources) == 0 {
		return nil, productSourcestock, ErrNoSourceAvailable
	}

	allocatedQtys := make(AllocatedQtys)

	allocatedQty, remainingSourcestock, err := d.allocateFromSources(ctx, productSourcestock, product, qtyToAllocate, sources, &deliveryInfo, allocatedQtys)
	if err != nil {
		return nil, productSourcestock, err
	}

	if allocatedQty < qtyToAllocate {
		return allocatedQtys, remainingSourcestock, ErrInsufficientSourceQty
	}

	return allocatedQtys, remainingSourcestock, nil
}

func (d *DefaultSourcingService) allocateFromSources(
	ctx context.Context,
	productSourcestock map[string]map[Source]int,
	product domain.BasicProduct,
	qtyToAllocate int,
	sources []Source,
	deliveryInfo *cartDomain.DeliveryInfo,
	allocatedQtys AllocatedQtys,
) (int, map[string]map[Source]int, error) {
	productID := product.GetIdentifier()
	if productID == "" {
		return 0, productSourcestock, ErrProductIDMissing
	}

	allocatedQty := 0
	remainingSourcestock := productSourcestock

	if _, exists := productSourcestock[productID]; !exists {
		productSourcestock[productID] = make(map[Source]int)
	}

	for _, source := range sources {
		sourceStock, err := d.getSourceStock(ctx, remainingSourcestock, product, productID, source, deliveryInfo)
		if err != nil {
			d.logger.Error(err)

			continue
		}

		if sourceStock == 0 {
			continue
		}

		stockToAllocate := min(qtyToAllocate-allocatedQty, sourceStock)
		remainingSourcestock[productID][source] -= stockToAllocate
		allocatedQty += stockToAllocate
		allocatedQtys[source] = stockToAllocate // Added this line to update allocatedQtys map
	}

	return allocatedQty, remainingSourcestock, nil
}

func (d *DefaultSourcingService) getSourceStock(
	ctx context.Context,
	remainingSourcestock map[string]map[Source]int,
	product domain.BasicProduct,
	productID string,
	source Source,
	deliveryInfo *cartDomain.DeliveryInfo,
) (int, error) {
	if _, exists := remainingSourcestock[productID][source]; !exists {
		sourceStock, err := d.stockProvider.GetStock(ctx, product, source, deliveryInfo)
		if err != nil {
			return 0, fmt.Errorf("error getting stock product: %w", err)
		}

		remainingSourcestock[productID][source] = sourceStock
	}

	return remainingSourcestock[productID][source], nil
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

func formatSources(sources []Source) string {
	checkedSources := "Checked sources:"

	for _, source := range sources {
		checkedSources += fmt.Sprintf(" SourceCode: %q ExternalSourceCode: %q", source.LocationCode, source.ExternalLocationCode)
	}

	return checkedSources
}
