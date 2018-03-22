package cart

type (
	Address struct {
		RegionCode             string
		CountryCode            string
		Company                string
		Street                 string
		StreetNr               string
		AdditionalAddressLines []string
		Telephone              string
		PostCode               string
		City                   string
		Firstname              string
		Lastname               string
		Email                  string
		Salutation             string
	}

	Destination struct {
		Type    string
		Address Address
		//Code - optional idendifier of this location/destination - might be used by special destinations
		Code string
	}
)

const (
	DESTINATION_TYPE_COLLECTIONPOINT = "collectionpoint"
	DESTINATION_TYPE_ADDRESS         = "address"
	DESTINATION_TYPE_FREIGHTSTATION  = "freight-station"
)
