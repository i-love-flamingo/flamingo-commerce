package application

import (
	"fmt"

	"github.com/pkg/errors"
	cartApplication "go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	DeliveryLocationsService interface {
		GetDeliveryLocations(ctx web.Context) (*DeliveryLocations, error)
	}

	SourcingEngine struct {
		Cartservice              cartApplication.CartService `inject:""`
		DecoratedCart            cart.DecoratedCart          `inject:""`
		DeliveryLocationsService DeliveryLocationsService    `inject:""`
		Logger                   flamingo.Logger             `inject:""`
	}

	DeliveryLocations struct {
		RetailerLocations        []RetailerLocationCollection
		CollectionPointLocations []Location
	}

	RetailerLocationCollection struct {
		Retailer  string
		Locations []Location
	}

	Location struct {
		Id       string
		Priority string
	}
)

// GetSources gets Sources and modifies the Cart Items
func (se *SourcingEngine) GetSources(ctx web.Context) error {
	locations, err := se.DeliveryLocationsService.GetDeliveryLocations(ctx)
	if err != nil {
		return errors.Wrap(err, "checkout.application.sourcingengine: Get sources failed")
	}

	decoratedCart, err := se.Cartservice.GetDecoratedCart(ctx)
	if err != nil {
		return errors.Wrap(err, "checkout.application.sourcingengine: Could not get the decorated Cart")
	}

	for _, decoratedCartItem := range decoratedCart.DecoratedItems {
		// get retailer code for item and then get the retailers ispu locations
		retailerCode := decoratedCartItem.Product.BaseData().RetailerCode
		ispuLocations, err := locations.getByRetailerCode(retailerCode)

		if err != nil {
			// todo: do we need additional error handling here?
			// cannot get location for product
			continue
		}

		// todo: get stock for product and check if a location with stock for the product is in ispulocations
		// currently just using the first locations id since there is no stock service to ask
		cartitem := decoratedCartItem.Item

		cartitem.SourceId = ispuLocations.Locations[0].Id
		err = decoratedCart.Cart.UpdateItem(ctx, cartitem)
		if err != nil {
			errors.Wrap(err, "masterdataportal.application.sourcelocator: Could not update cart item")
		}
	}

	return nil
}

// getByRetailerCode returns just the RetailerLocationCollection for a given Retailer from a List of
func (dl *DeliveryLocations) getByRetailerCode(retailerCode string) (RetailerLocationCollection, error) {
	var result RetailerLocationCollection

	for _, retailerLocation := range dl.RetailerLocations {
		if retailerLocation.Retailer == retailerCode {
			result = retailerLocation
			break
		}
	}

	if result.Retailer == "" {
		return result, fmt.Errorf("could not find ispu location for retailer %s", retailerCode)
	}

	return result, nil
}
