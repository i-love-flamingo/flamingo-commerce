package infrastructure

import (
	checkoutApplication "go.aoe.com/flamingo/core/checkout/application"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// FakeSourceLocator represents the fake source locator
	FakeDeliveryLocationsService struct {
	}
)

var (
	_ checkoutApplication.DeliveryLocationsService = new(FakeDeliveryLocationsService)
)

// FakeDeliveryLocationsService provides fake delivery locations
func (sl *FakeDeliveryLocationsService) GetDeliveryLocations(ctx web.Context) (*checkoutApplication.DeliveryLocations, error) {
	return &checkoutApplication.DeliveryLocations{
		RetailerLocations: []checkoutApplication.RetailerLocationCollection{
			{
				Retailer: "om3CommonTestretailer",
				Locations: []checkoutApplication.Location{
					{
						Id:       "mock_ispu_location1",
						Priority: "1",
					},
					{
						Id:       "mock_ispu_location2",
						Priority: "2",
					},
					{
						Id:       "mock_ispu_location3",
						Priority: "3",
					},
				},
			},
		},
	}, nil
}
