package domain

import (
	"context"
	"fmt"
	"math"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/product/domain"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// SourcingService helps in retrieving item sources
	// Deprecated: Sourcing moved to separate module
	SourcingService interface {
		//GetSourceID  returns one source location code where the product should be sourced
		//@todo will be Deprecated in future in favor of SourcingServiceDetail interface
		GetSourceID(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, item *decorator.DecoratedCartItem) (string, error)
	}
	// SourcingServiceDetail additional interface to return
	// Deprecated: Sourcing moved to separate module
	SourcingServiceDetail interface {
		//GetSourcesForItem returns Sources for the given item in the cart
		GetSourcesForItem(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, item *decorator.DecoratedCartItem) (Sources, error)
		//GetAvailableSources returns Sources for the product - containing the maximum possible qty per source
		GetAvailableSources(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, product domain.BasicProduct) (Sources, error)
	}

	// Sources is the result value object containing all sources (for the request item or product)
	// Deprecated: Sourcing moved to separate module
	Sources []Source

	// Source represents the Sourcing info
	// Deprecated: Sourcing moved to separate module
	Source struct {
		// LocationCode identifies the warehouse or stock location
		LocationCode string
		// Qty for the sources items
		Qty int
		// ExternalLocationCode identifies the source location in an external system
		ExternalLocationCode string
	}

	// SourcingEngine computes item sources
	// Deprecated: Sourcing moved to separate module
	SourcingEngine struct {
		SourcingService SourcingService          `inject:",optional"`
		Logger          flamingo.Logger          `inject:""`
		Cartservice     *application.CartService `inject:""`
	}
)

var (
	// ErrInsufficientSourceQty - use to indicate that the requested qty exceeds the available qty
	ErrInsufficientSourceQty = errors.New("Available Source Qty insufficient")
	// ErrNoSourceAvailable - use to indicate that no source for item is available at all
	ErrNoSourceAvailable = errors.New("No Available Source Qty")
)

const (
	// ExternalSourceIDKey specifies the key for the ItemUpdateCommand.AdditionalData map where the external source id should be stored
	ExternalSourceIDKey = "external_source_id"
)

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

// Reduce returns new Source
func (s Sources) Reduce(reduceby Sources) Sources {
	for k, source := range s {
		for _, reducebySource := range reduceby {
			if source.LocationCode == reducebySource.LocationCode {
				s[k].Qty = s[k].Qty - reducebySource.Qty
			}
		}
	}
	return s
}

// SetSourcesForCartItems gets Sources and modifies the Cart Items
// todo move to application layer ?
func (se *SourcingEngine) SetSourcesForCartItems(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart) error {
	if se.SourcingService == nil {
		return nil
	}

	itemUpdateCommands := make([]cartDomain.ItemUpdateCommand, 0)
	for _, decoratedDelivery := range decoratedCart.DecoratedDeliveries {
		for _, decoratedCartItem := range decoratedDelivery.DecoratedItems {
			source, err := se.sourceLocationForCartItem(ctx, session, decoratedCart, decoratedDelivery.Delivery.DeliveryInfo.Code, &decoratedCartItem)
			if err != nil {
				se.Logger.WithContext(ctx).WithField("subcategory", "SourcingEngine").Error(err)
				return fmt.Errorf("checkout.application.sourcingengine error: %v", err)
			}
			se.Logger.WithContext(ctx).WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Debug("SourcingEngine detected source %v for item %v", source.LocationCode, decoratedCartItem.Item.ID)

			itemUpdate := cartDomain.ItemUpdateCommand{
				SourceID: &source.LocationCode,
				ItemID:   decoratedCartItem.Item.ID,
				// ExternalSourceID contains the picking location used by an external system
				AdditionalData: map[string]string{ExternalSourceIDKey: source.ExternalLocationCode},
			}

			itemUpdateCommands = append(itemUpdateCommands, itemUpdate)
		}
	}

	err := se.Cartservice.UpdateItems(ctx, session, itemUpdateCommands)
	if err != nil {
		return errors.Wrap(err, "Could not update cart items")
	}

	return nil
}

func (se *SourcingEngine) sourceLocationForCartItem(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, item *decorator.DecoratedCartItem) (*Source, error) {
	detailService, ok := se.SourcingService.(SourcingServiceDetail)
	if ok {
		sources, err := detailService.GetSourcesForItem(ctx, session, decoratedCart, deliveryCode, item)
		if err != nil {
			return nil, err
		}

		if len(sources) == 0 {
			return nil, errors.New("no source locations found")
		}

		return &sources[0], nil
	}

	sourceID, err := se.SourcingService.GetSourceID(ctx, session, decoratedCart, deliveryCode, item)
	if err != nil {
		return nil, err
	}

	return &Source{LocationCode: sourceID}, nil
}
