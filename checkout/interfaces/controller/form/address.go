package form

type (
	AddressFormData struct {
		RegionCode   string `form:"regionCode" validate:"required" conform:"name"`
		CountryCode  string `form:"countryCode" validate:"required" conform:"name"`
		Company      string `form:"company" validate:"required" conform:"name"`
		Street       string `form:"street" validate:"required" conform:"name"`
		StreetNr     string `form:"streetNr" validate:"required" conform:"name"`
		AddressLine1 string `form:"addressLine1" validate:"required" conform:"name"`
		AddressLine2 string `form:"addressLine2" validate:"required" conform:"name"`
		Telephone    string `form:"telephone" validate:"required" conform:"name"`
		PostCode     string `form:"postCode" validate:"required" conform:"name"`
		City         string `form:"city" validate:"required" conform:"name"`
		Firstname    string `form:"firstname" validate:"required" conform:"name"`
		Lastname     string `form:"lastname" validate:"required" conform:"name"`
		Email        string `form:"email" validate:"required" conform:"name"`
	}
)
