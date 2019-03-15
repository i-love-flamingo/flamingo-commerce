package form

import "time"

type (

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
