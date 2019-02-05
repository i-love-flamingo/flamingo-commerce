package cart

type (
	// Address defines a cart address
	Address struct {
		Vat                    string
		Firstname              string
		Lastname               string
		MiddleName             string
		Title                  string
		Salutation             string
		Street                 string
		StreetNr               string
		AdditionalAddressLines []string
		Company                string
		City                   string
		PostCode               string
		State                  string
		RegionCode             string
		Country                string
		CountryCode            string
		Telephone              string
		Email                  string
		SameAsBilling          bool
	}
)
