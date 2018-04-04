package cart

import (
	"fmt"
	"log"

	"time"

	"github.com/pkg/errors"
)

type (
	//CartProvider should be used to create the cart Value objects
	CartProvider func() *Cart

	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		//ID is the main idendifier of the cart
		ID string
		//EntityID is a second idendifier that may be used by some backends
		EntityID string
		//Cartitems - list of cartitems
		Cartitems []Item
		//CartTotals - the cart totals (contain summary costs and discounts etc)
		CartTotals CartTotals
		//BillingAdress - the main billing address (relevant for all payments/invoices)
		BillingAdress Address

		//Purchaser - additional infos for the legal contact person in this order
		Purchaser Person

		//DeliveryInfos - list of desired Deliverys (or Shippments) involved in this cart - referenced from the items
		DeliveryInfos []DeliveryInfo
		//AdditionalData   can be used for Custom attributes
		AdditionalData map[string]string

		//IsCustomerCart - false = Guest Cart true = cart from the authenticated user
		IsCustomerCart bool
	}

	Person struct {
		Address         Address
		PersonalDetails PersonalDetails
	}

	PersonalDetails struct {
		DateOfBirth     string
		PassportCountry string
		PassportNumber  string
	}

	//DeliveryInfo - represents the Delivery
	DeliveryInfo struct {
		ID               string
		Method           string
		Carrier          string
		DeliveryLocation DeliveryLocation
		ShippingItem     ShippingItem
		DesiredTime      time.Time
		AdditionalData   map[string]string
	}

	DeliveryLocation struct {
		Type    string
		Address Address
		//Code - optional idendifier of this location/destination - is used in special destination Types
		Code string
	}

	CartTotals struct {
		Totalitems        []Totalitem
		TotalShippingItem ShippingItem
		GrandTotal        float64
		SubTotal          float64
		DiscountAmount    float64
		TaxAmount         float64
		CurrencyCode      string
	}

	// Item for Cart
	Item struct {
		ID              string
		MarketplaceCode string
		//VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceId string

		Price float64
		Qty   int

		RowTotal       float64
		TaxAmount      float64
		DiscountAmount float64

		PriceInclTax    float64
		RowTotalInclTax float64

		DeliveryInfoReference *DeliveryInfo
		CurrencyCode          string

		//OriginalDeliveryIntent can be "delivery" for homedelivery or "pickup_locationcode" or "collectionpoint_locationcode"
		OriginalDeliveryIntent *DeliveryIntent

		AdditionalData map[string]string
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price float64
	}

	// ShippingItem
	ShippingItem struct {
		Title string
		Price float64

		TaxAmount      float64
		DiscountAmount float64

		CurrencyCode string
	}
)

const (
	DELIVERY_METHOD_PICKUP      = "pickup"
	DELIVERY_METHOD_DELIVERY    = "delivery"
	DELIVERY_METHOD_UNSPECIFIED = "unspecified"

	DELIVERYLOCATION_TYPE_COLLECTIONPOINT = "collection"
	DELIVERYLOCATION_TYPE_STORE           = "store"
	DELIVERYLOCATION_TYPE_ADDRESS         = "address"
	DELIVERYLOCATION_TYPE_FREIGHTSTATION  = "freight-station"
)

// GetByLineNr gets an item - starting with 1
func (Cart Cart) GetByLineNr(lineNr int) (*Item, error) {
	var item Item
	if len(Cart.Cartitems) >= lineNr && lineNr > 0 {
		return &Cart.Cartitems[lineNr-1], nil
	} else {
		return &item, errors.New("Line in cart not existend")
	}
}

// GetByItemId gets an item by its id
func (Cart Cart) GetByItemId(itemId string) (*Item, error) {
	for _, currentItem := range Cart.Cartitems {
		log.Println(currentItem.ID)
		if currentItem.ID == itemId {
			return &currentItem, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("itemId %v in cart not existend", itemId))
}

// HasItem checks if a cartitem for that sku exists and returns lineNr if found
func (cart Cart) HasItem(marketplaceCode string, variantMarketplaceCode string) (bool, int) {
	for lineNr, item := range cart.Cartitems {
		if item.MarketplaceCode == marketplaceCode && item.VariantMarketPlaceCode == variantMarketplaceCode {
			return true, lineNr + 1
		}
	}
	return false, 0
}

func inStruct(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// ItemCount - returns amount of Cartitems
func (Cart Cart) ItemCount() int {
	return len(Cart.Cartitems)
}

// GetItemIds - returns amount of Cartitems
func (Cart Cart) GetItemIds() []string {
	ids := []string{}
	for _, item := range Cart.Cartitems {
		ids = append(ids, item.ID)
	}
	return ids
}

// GetCartItemsByOriginalIntend - returns the cart Items grouped by the original DeliveryIntent
func (Cart Cart) GetCartItemsByOriginalDeliveryIntent() map[string][]Item {
	result := make(map[string][]Item)
	for _, item := range Cart.Cartitems {
		result[item.OriginalDeliveryIntent.String()] = append(result[item.OriginalDeliveryIntent.String()], item)
	}
	return result
}

// HasDeliveryMethodForIntent - returns if the cart has an item with the delivery intent
func (Cart Cart) HasDeliveryMethodForIntent(intent string) bool {
	for _, item := range Cart.Cartitems {
		if item.OriginalDeliveryIntent.String() == intent {
			return true
		}
	}
	return false
}
