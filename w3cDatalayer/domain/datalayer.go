package domain

type (
	DatalayerProvider func() *Datalayer
	/**
	Datalayer Value object - represents the structure of the w3c Datalayer.
	Therefore it has the json annotations and its intended to be directly converted to Json in the output
	*/
	Datalayer struct {
		PageInstanceID string    `json:"pageInstanceID" inject:"config:w3cDatalayer.pageInstanceID,optional"`
		Page           *Page     `json:"page,omitempty"`
		SiteInfo       *SiteInfo `json:"siteInfo,omitempty"`
		Version        string    `json:"version" inject:"config:w3cDatalayer.version,optional"`
		//User List of user(s) interacting with the page. (Although typically web data has a single user per recorded interaction, this object is an array and can capture multiple users.)
		User []User `json:"user,omitempty"`
	}
	Page struct {
		PageInfo   PageInfo               `json:"pageInfo,omitempty"`
		Category   PageCategory           `json:"category,omitempty"`
		Attributes map[string]interface{} `json:"attributes,omitempty"`
	}
	PageInfo struct {
		PageID         string `json:"pageID,omitempty"`
		DestinationURL string `json:"destinationURL,omitempty"`
		BreadCrumbs    string `json:"breadCrumbs,omitempty"`
		PageName       string `json:"pageName,omitempty"`
		ReferringUrl   string `json:"referringUrl,omitempty"`
		Language       string `json:"language,omitempty"`
	}
	PageCategory struct {
		PrimaryCategory string `json:"primaryCategory,omitempty"`
		SubCategory1    string `json:"subCategory1,omitempty"`
		SubCategory2    string `json:"subCategory2,omitempty"`
		PageType        string `json:"pageType,omitempty"`
		Section         string `json:"section,omitempty"`
	}

	SiteInfo struct {
		SiteName string `json:"siteName,omitempty"`
		Domain   string `json:"domain,omitempty"`
	}

	//User The User object captures the profile of a user who is interacting with the website.
	User struct {
		/**
		Profile
		A profile for information about the user, typically associated with a registered user. (Although
		typically a user might have only a single profile, this object is an array and can capture multiple
		profiles per user.)
		*/
		Profile []UserProfile `json:"profile,omitempty"`
		/**
		Segment This object provides population segmentation information for the user, such as premium versus
		basic membership, or any other forms of segmentation that are desirable. Any additional
		dimensions related to the user can be provided. All names are optional and should fit the
		individual implementation needs in both naming and values passed.
		*/
		Segment string `json:"segment,omitempty"`
	}
	// UserProfile A profile for information about the user
	UserProfile struct {
		ProfileInfo UserProfileInfo `json:"profileInfo,omitempty"`
	}
	//UserProfileInfo An extensible object for providing information about the user.
	UserProfileInfo struct {
		EmailID   string `json:"emailID,omitempty"`
		ProfileID string `json:"profileID,omitempty"`
		Rewards   string `json:"rewards,omitempty"`
	}
)
