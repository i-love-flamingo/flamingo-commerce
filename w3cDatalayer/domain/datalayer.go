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
	}
	Page struct {
		PageInfo   PageInfo          `json:"pageInfo,omitempty"`
		Category   PageCategory      `json:"category,omitempty"`
		Attributes map[string]string `json:"attributes,omitempty"`
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
)
