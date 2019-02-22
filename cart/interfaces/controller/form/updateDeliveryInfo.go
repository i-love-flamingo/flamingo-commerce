package form

import "time"

type (
	// TODO: this is temp moved code?
	// CheckoutFormData defines the default form data provided by the checkout form
	UpdateDeliveryInfoCommandDto struct {
		Code             string `form:"code"`
		Method           string	`form:"method"`
		Carrier          string `form:"carrier"`
		//DeliveryLocation DeliveryLocation
		DesiredTime      string `form:"desiredTime"`
		//AdditionalData   map[string]string
		RelatedFlight FlightData `form:"flightData"`
	}

	// FlightData value object
	FlightData struct {
		ScheduledDateTime  time.Time
		Direction          string
		FlightNumber       string
		AirportName        string
		DestinationCountry string
		Terminal           string
		AirlineName        string
		AirlineCode        string
		BookingCode        string
		LastName           string
		FirstName          string
	}
)
