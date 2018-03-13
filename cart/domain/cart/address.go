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
)
