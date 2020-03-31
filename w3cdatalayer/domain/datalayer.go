package domain

import (
	"encoding/gob"
	"encoding/json"
)

type (
	// DatalayerProvider func
	DatalayerProvider func() *Datalayer

	// Datalayer Value object - represents the structure of the w3c Datalayer.
	// Therefore it has the json annotations and its intended to be directly converted to Json in the output
	Datalayer struct {
		PageInstanceID string
		Page           *Page
		SiteInfo       *SiteInfo
		Version        string `inject:"config:commerce.w3cDatalayer.version,optional"` // todo: version should not be injected here (domain layer)
		//User List of user(s) interacting with the page. (Although typically web data has a single user per recorded interaction, this object is an array and can capture multiple users.)
		User []User
		//The Cart object carries details about a shopping cart or basket and the products that have been added to it.
		Cart *Cart
		// The Event object collects information about an interaction event by the user. An event might be a button click, the addition of a portal widget, playing a video, adding a product to the shopping cart, etc. Any action on the page could be captured by an Event object.
		Event []Event
		//The Product object carries details about a particular product with frequently used properties listed below. This is intended for data about products displayed on pages or other content. For products added to a shopping cart or ordered in a transaction, see the Cart and Transaction objects below.
		Product []Product
		//The Transaction object is similar to the Cart object, but represents a completed order. The Transaction object contains analogous sub-objects to the Cart object as well as additional subobjects specific to completed orders.
		Transaction *Transaction
	}

	// Product struct
	Product struct {
		ProductInfo ProductInfo            `json:"productInfo"`
		Category    *ProductCategory       `json:"category,omitempty"`
		Attributes  map[string]interface{} `json:"attributes,omitempty"`
	}

	// Event struct
	Event struct {
		EventInfo map[string]interface{} `json:"eventInfo,omitempty"`
	}

	// Page struct
	Page struct {
		PageInfo   PageInfo               `json:"pageInfo,omitempty"`
		Category   PageCategory           `json:"category,omitempty"`
		Search     SearchInfo             `json:"search,omitempty"`
		Attributes map[string]interface{} `json:"attributes,omitempty"`
	}

	// SearchInfo struct
	SearchInfo struct {
		SearchKeyword string      `json:"searchKeyword,omitempty"`
		Result        interface{} `json:"result,omitempty"`
	}

	// PageInfo generall information about the page
	PageInfo struct {
		PageID         string `json:"pageID,omitempty"`
		DestinationURL string `json:"destinationURL,omitempty"`
		BreadCrumbs    string `json:"breadCrumbs,omitempty"`
		PageName       string `json:"pageName,omitempty"`
		ReferringURL   string `json:"referringUrl,omitempty"`
		Language       string `json:"language,omitempty"`
		ErrorName      string `json:"errorName,omitempty"`
	}

	// PageCategory struct for the category and subcategories
	PageCategory struct {
		PrimaryCategory string `json:"primaryCategory"`
		SubCategory1    string `json:"subCategory1,omitempty"`
		SubCategory2    string `json:"subCategory2,omitempty"`
		PageType        string `json:"pageType,omitempty"`
		Section         string `json:"section,omitempty"`
	}

	// SiteInfo - struct for SiteName and Domain
	SiteInfo struct {
		SiteName string `json:"siteName,omitempty"`
		Domain   string `json:"domain,omitempty"`
	}

	// User The User object captures the profile of a user who is interacting with the website.
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
		ProfileInfo     UserProfileInfo `json:"profileInfo,omitempty"`
		Address         *Address        `json:"address,omitempty"`
		ShippingAddress *Address        `json:"shippingAddress,omitempty"`
	}

	// Address basic address information
	Address struct {
		Line1         string `json:"line1,omitempty"`
		Line2         string `json:"line2,omitempty"`
		City          string `json:"city,omitempty"`
		StateProvince string `json:"stateProvince,omitempty"`
		PostalCode    string `json:"postalCode,omitempty"`
		Country       string `json:"country,omitempty"`
	}

	// UserProfileInfo An extensible object for providing information about the user.
	UserProfileInfo struct {
		EmailID   string `json:"emailID,omitempty"`
		UserName  string `json:"userName,omitempty"`
		ProfileID string `json:"profileID"`
		Rewards   string `json:"rewards,omitempty"`
	}

	// Cart cartInformation
	Cart struct {
		CartID     string                 `json:"cartID,omitempty"`
		Price      *CartPrice             `json:"price,omitempty"`
		Attributes map[string]interface{} `json:"attributes,omitempty"`
		Item       []CartItem             `json:"item,omitempty"`
	}

	// CartPrice used in Cart
	CartPrice struct {
		//The basePrice SHOULD be the price of the items before applicable discounts,shipping charges, and tax.
		BasePrice       float64 `json:"basePrice"`
		VoucherCode     string  `json:"voucherCode"`
		VoucherDiscount float64 `json:"voucherDiscount"`
		Currency        string  `json:"currency"`
		TaxRate         float64 `json:"taxRate"`
		Shipping        float64 `json:"shipping"`
		ShippingMethod  string  `json:"shippingMethod"`
		PriceWithTax    float64 `json:"priceWithTax"`
		//cartTotal SHOULD be the total price inclusive of all discounts, charges, and tax
		CartTotal float64 `json:"cartTotal"`
	}

	// CartItem used in Cart
	CartItem struct {
		ProductInfo ProductInfo            `json:"productInfo"`
		Quantity    int                    `json:"quantity"`
		Category    *ProductCategory       `json:"category,omitempty"`
		Price       CartItemPrice          `json:"price"`
		Attributes  map[string]interface{} `json:"attributes,omitempty"`
	}

	// CartItemPrice struct
	CartItemPrice struct {
		BasePrice    float64 `json:"basePrice"`
		Currency     string  `json:"currency"`
		TaxRate      float64 `json:"taxRate"`
		PriceWithTax float64 `json:"priceWithTax"`
	}

	// Transaction struct
	Transaction struct {
		TransactionID string                 `json:"transactionID,omitempty"`
		Profile       *UserProfile           `json:"profile,omitempty"`
		Price         *TransactionPrice      `json:"total,omitempty"`
		Item          []CartItem             `json:"item,omitempty"`
		Attributes    map[string]interface{} `json:"attributes,omitempty"`
	}

	// TransactionPrice struct
	TransactionPrice struct {
		//The basePrice SHOULD be the price of the items before applicable discounts,shipping charges, and tax.
		BasePrice        float64 `json:"basePrice"`
		VoucherCode      string  `json:"voucherCode"`
		VoucherDiscount  float64 `json:"voucherDiscount"`
		Currency         string  `json:"currency"`
		TaxRate          float64 `json:"taxRate"`
		Shipping         float64 `json:"shipping"`
		ShippingMethod   string  `json:"shippingMethod"`
		PriceWithTax     float64 `json:"priceWithTax"`
		TransactionTotal float64 `json:"transactionTotal"`
	}

	// ProductCategory struct
	ProductCategory struct {
		PrimaryCategory string `json:"primaryCategory,omitempty"`
		SubCategory1    string `json:"subCategory1,omitempty"`
		SubCategory     string `json:"subCategory,omitempty"`
		SubCategory2    string `json:"subCategory2,omitempty"`
		ProductType     string `json:"productType,omitempty"`
	}

	// ProductInfo dataLayer product information
	ProductInfo struct {
		ProductID   string `json:"productID"`
		SKU         string `json:"sku"`
		ProductName string `json:"productName"`
		//ProductURL               string  `json:"productURL"`
		ProductImage             string  `json:"productImage"`
		ProductThumbnail         string  `json:"productThumbnail"`
		Manufacturer             string  `json:"manufacturer"`
		Size                     string  `json:"size"`
		Color                    string  `json:"color"`
		ParentID                 *string `json:"parentId,omitempty"`
		VariantSelectedAttribute *string `json:"variantSelectedAttribute,omitempty"`
		ProductType              string  `json:"productType"`
		Retailer                 string  `json:"retailer"`
		Brand                    string  `json:"brand"`
		InStock                  string  `json:"inStock"`
	}
)

func init() {
	gob.Register(Event{})
}

// MarshalJSON - is here to make sure the renderingengine uses this json interface instead of own encoding
func (d Datalayer) MarshalJSON() ([]byte, error) {
	//myDataLayer should match the Datalayer struct and is just here to define the top level json marshal annotations
	// we need this since json.Marshal(&d) would result in endless loop
	type myDataLayer struct {
		PageInstanceID string    `json:"pageInstanceID"`
		Page           *Page     `json:"page,omitempty"`
		SiteInfo       *SiteInfo `json:"siteInfo,omitempty"`
		Version        string    `json:"version"`
		//User List of user(s) interacting with the page. (Although typically web data has a single user per recorded interaction, this object is an array and can capture multiple users.)
		User []User `json:"user"`
		//The Cart object carries details about a shopping cart or basket and the products that have been added to it.
		Cart *Cart `json:"cart,omitempty"`
		// The Event object collects information about an interaction event by the user. An event might be a button click, the addition of a portal widget, playing a video, adding a product to the shopping cart, etc. Any action on the page could be captured by an Event object.
		Event []Event `json:"event,omitempty"`
		//The Product object carries details about a particular product with frequently used properties listed below. This is intended for data about products displayed on pages or other content. For products added to a shopping cart or ordered in a transaction, see the Cart and Transaction objects below.
		Product []Product `json:"product,omitempty"`
		//The Transaction object is similar to the Cart object, but represents a completed order. The Transaction object contains analogous sub-objects to the Cart object as well as additional subobjects specific to completed orders.
		Transaction *Transaction `json:"transaction,omitempty"`
	}
	myDataLayerInstance := myDataLayer(d)
	return json.Marshal(&myDataLayerInstance)
}
