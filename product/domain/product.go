package domain

import "time"

type (
	// Product is the basic product model, as read from searchperience
	Product struct {
		Locale                     string              `json:"locale"`
		Channel                    string              `json:"channel"`
		ForeignID                  string              `json:"foreignId"`
		MarketPlaceCode            string              `json:"marketPlaceCode"`
		ProductType                string              `json:"productType"`
		FormatVersion              int                 `json:"formatVersion"`
		StockLevel                 string              `json:"stockLevel"`
		InternalName               string              `json:"internalName"`
		ProductFamily              string              `json:"productFamily"`
		CreatedAt                  time.Time           `json:"createdAt"`
		UpdatedAt                  time.Time           `json:"updatedAt"`
		VisibleFrom                time.Time           `json:"visibleFrom"`
		VisibleTo                  time.Time           `json:"visibleTo"`
		CustomerRating             int                 `json:"customerRating"`
		TeaserData                 TeaserData          `json:"teaserData"`
		Sorting                    map[string]float64  `json:"sorting"`
		CategoryPath               []string            `json:"categoryPath"`
		CategoryCodes              []string            `json:"categoryCodes"`
		DocumentBoost              int                 `json:"documentBoost"`
		Highlights                 struct{}            `json:"highlights"`
		Keywords                   []string            `json:"keywords"`
		KeywordsImportant          []string            `json:"keywordsImportant"`
		ConfigurableProduct        ConfigurableProduct `json:"configurableProduct"`
		VariantVariationAttributes []string            `json:"variantVariationAttributes"`
		Variants                   []Variant           `json:"variants"`
	}

	// TeaserData is the teaser-information for product previews
	TeaserData struct {
		Title            string  `json:"title"`
		ShortTitle       string  `json:"shortTitle"`
		Teaser           string  `json:"teaser"`
		ShortDescription string  `json:"shortDescription"`
		Media            []Media `json:"media"`
	}

	// Media holds product media information
	Media struct {
		Type      string `json:"type"`
		MimeType  string `json:"mimeType"`
		Usage     string `json:"usage"`
		Title     string `json:"title"`
		Reference string `json:"reference"`
	}

	// ConfigurableProduct defines the variant setup
	ConfigurableProduct struct {
		InternalName string     `json:"internalName"`
		Title        string     `json:"title"`
		CreatedAt    time.Time  `json:"createdAt"`
		UpdatedAt    time.Time  `json:"updatedAt"`
		Attributes   Attributes `json:"attributes"`
		Media        []Media    `json:"media"`
	}

	// Variant is a concrete variant of a product
	Variant struct {
		InternalName          string     `json:"internalName"`
		Title                 string     `json:"title"`
		CreatedAt             time.Time  `json:"createdAt"`
		UpdatedAt             time.Time  `json:"updatedAt"`
		IsSaleable            bool       `json:"isSaleable"`
		SaleableFrom          time.Time  `json:"saleableFrom,omitempty"`
		SaleableTo            time.Time  `json:"saleableTo,omitempty"`
		Attributes            Attributes `json:"attributes"`
		ShortDescription      string     `json:"shortDescription"`
		Description           string     `json:"description"`
		CustomerRating        int        `json:"customerRating"`
		Media                 []Media    `json:"media"`
		MarketPlaceCode       string     `json:"marketPlaceCode"`
		RetailerSku           string     `json:"retailerSku"`
		RetailerCode          string     `json:"retailerCode"`
		RetailerName          string     `json:"retailerName"`
		RetailerImage         string     `json:"retailerImage"`
		OriginPrice           int        `json:"originPrice"`
		OriginBasePrice       float64    `json:"originBasePrice"`
		OriginBasePriceUnit   string     `json:"originBasePriceUnit"`
		OriginBasePriceAmount float64    `json:"originBasePriceAmount"`
		SpecialPrice          float64    `json:"specialPrice,omitempty"`
		SpecialPriceFrom      float64    `json:"specialPriceFrom,omitempty"`
		SpecialPriceTo        float64    `json:"specialPriceTo,omitempty"`
		Currency              string     `json:"currency,omitempty"`
		TaxClass              string     `json:"taxClass,omitempty"`
	}

	// Attributes is a generic map[string]interface{}
	Attributes map[string]interface{}
)
