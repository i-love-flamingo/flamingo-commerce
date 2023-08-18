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
		TelephoneCountryCode   string
		TelephoneAreaCode      string
		TelephoneNumber        string
		// Deprecated: parts of number should be distinguished, please use TelephoneCountryCode, TelephoneAreaCode and TelephoneNumber
		Telephone string
		Email     string
	}
)

// FullName return Firstname Lastname
func (a Address) FullName() string {
	return a.Firstname + " " + a.Lastname
}

// IsEmpty checks all fields of the address if they are empty
func (a *Address) IsEmpty() bool {
	if a == nil {
		return true
	}

	for _, additionalLine := range a.AdditionalAddressLines {
		if additionalLine != "" {
			return false
		}
	}

	if a.Vat != "" ||
		a.Firstname != "" ||
		a.Lastname != "" ||
		a.MiddleName != "" ||
		a.Title != "" ||
		a.Salutation != "" ||
		a.Street != "" ||
		a.StreetNr != "" ||
		a.Company != "" ||
		a.City != "" ||
		a.PostCode != "" ||
		a.State != "" ||
		a.RegionCode != "" ||
		a.Country != "" ||
		a.CountryCode != "" ||
		a.Telephone != "" ||
		a.Email != "" {

		return false
	}

	return true
}
