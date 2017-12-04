package application

import (
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
	"fmt"
	"errors"
)

type (
	SourceLocator interface {
		DeliveryLocationsService(ctx web.Context) (DeliveryLocations, error)
	}

	SourcingEngine struct {
		DecoratedCart cart.DecoratedCart `inject:""`
		SourceLocator SourceLocator `inject:""`
		Logger flamingo.Logger `inject:""`
	}

	DeliveryLocations struct {
		RetailerLocations []RetailerLocationCollection
		CollectionPointLocations []Location
	}

	RetailerLocationCollection struct {
		Retailer string
		Locations []Location
	}

	Location struct {
		Id string
		Priority string
	}
)
// GetSources gets Sources and modifies the Cart Items
func (se *SourcingEngine) GetSources(ctx web.Context) error {
	locations, err := se.SourceLocator.DeliveryLocationsService(ctx)

	if err != nil {
		se.Logger.Error("checkout.application.sourcingengine: Get sources failed")
	}

	for _, decoratedCartItem := range se.DecoratedCart.DecoratedItems {
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
		decoratedCartItem.Item.SourceId = ispuLocations.Locations[0].Id
	}

	return nil
}

// getByRetailerCode returns just the RetailerLocationCollection for a given Retailer from a List of
func (dl *DeliveryLocations) getByRetailerCode (retailerCode string) (RetailerLocationCollection, error) {
	var result RetailerLocationCollection

	for _, retailerLocation := range dl.RetailerLocations {
		if retailerLocation.Retailer == retailerCode {
			result = retailerLocation
		}
	}

	if result.Retailer == "" {
		return result, errors.New(fmt.Sprintf("could not find ispu location for retailer %s", retailerCode))
	}

	return result, nil
}
