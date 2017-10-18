package domain

import (
	"go.aoe.com/flamingo/core/product/domain"
)

type (
	// SearchResult defines our search model
	SearchResult struct {
		Results struct {
			Retailer struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []interface{} `json:"hits"`
			} `json:"retailer"`
			Location struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []struct {
					Document   Location `json:"document"`
					Highlights struct {
					} `json:"highlights"`
				} `json:"hits"`
			} `json:"location"`
			Brand struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []struct {
					Document   Brand `json:"document"`
					Highlights struct {
					} `json:"highlights"`
				} `json:"hits"`
			} `json:"brand"`
			Product struct {
				MetaData MetaData      `json:"metaData"`
				Facets   []interface{} `json:"facets"`
				Filters  []interface{} `json:"filters"`
				PageInfo PageInfo      `json:"pageInfo"`
				Hits     []struct {
					Document   domain.BasicProduct `json:"document"`
					Highlights struct {
					} `json:"highlights"`
				} `json:"hits"`
			} `json:"product"`
		} `json:"results"`
	}

	// Media is a generic media type
	Media struct {
		MimeType  string `json:"mimeType"`
		Reference string `json:"reference"`
		Title     string `json:"title"`
		Type      string `json:"type"`
		Usage     string `json:"usage"`
	}

	// Brand is a product brand
	Brand struct {
		Channel          string   `json:"channel"`
		ForeignID        string   `json:"foreignId"`
		FormatVersion    int      `json:"formatVersion"`
		Keywords         []string `json:"keywords"`
		Locale           string   `json:"locale"`
		Media            []Media  `json:"media"`
		ShortDescription string   `json:"shortDescription"`
		ShortTitle       string   `json:"shortTitle"`
		Teaser           string   `json:"teaser"`
		Title            string   `json:"title"`
	}

	// Location is a place
	Location struct {
		ForeignID        string `json:"foreignId"`
		Locale           string `json:"locale"`
		Channel          string `json:"channel"`
		Code             string `json:"code"`
		FormatVersion    int    `json:"formatVersion"`
		Title            string `json:"title"`
		ShortTitle       string `json:"shortTitle"`
		ShortDescription string `json:"shortDescription"`
		Description      string `json:"description"`
		Type             string `json:"type"`
		AirportZone      struct {
			Title    string `json:"title"`
			Area     string `json:"area"`
			Level    string `json:"level"`
			Terminal string `json:"terminal"`
			Landside bool   `json:"landside"`
			Schengen bool   `json:"schengen"`
		} `json:"airportZone"`
		Brands   []string `json:"brands"`
		Media    []Media  `json:"media"`
		Counter  string   `json:"counter"`
		Email    string   `json:"email"`
		Phone    string   `json:"phone"`
		Pickup   bool     `json:"pickup"`
		Keywords []string `json:"keywords"`
	}

	// MetaData is search related meta information
	MetaData struct {
		TotalHits    int    `json:"totalHits"`
		Took         int    `json:"took"`
		CurrentQuery string `json:"currentQuery"`
		FacetMapping []struct {
			DocumentType string        `json:"documentType"`
			FacetNames   []interface{} `json:"facetNames"`
		} `json:"facetMapping"`
		SortMapping []struct {
			DocumentType string `json:"documentType"`
			Sorts        struct {
			} `json:"sorts"`
		} `json:"sortMapping"`
	}

	// PageInfo gives information about multiple page results
	PageInfo struct {
		CurrentPage      int `json:"currentPage"`
		TotalPages       int `json:"totalPages"`
		VisiblePageLinks int `json:"visiblePageLinks"`
		Padding          int `json:"padding"`
	}
)
