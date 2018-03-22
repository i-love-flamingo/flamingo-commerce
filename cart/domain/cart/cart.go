package cart

import (
	"fmt"
	"log"

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
		//BillingAdress - the main billing address (relevant for all payments)
		BillingAdress Address
		//PaymentInfos - list of Payments involved in this cart - referenced from the items
		PaymentInfos []PaymentInfo
		//DeliveryInfos - list of desired Deliverys (or Shippments) involved in this cart - referenced from the items
		DeliveryInfos []DeliveryInfo
	}

	//DeliveryInfo - informations for a wished delivery
	DeliveryInfo struct {
		Method       string
		Carrier      string
		Destination  Destination
		ShippingItem ShippingItem
	}

	CartTotals struct {
		Totalitems []Totalitem

		GrandTotal     float64
		SubTotal       float64
		DiscountAmount float64
		TaxAmount      float64
		CurrencyCode   string
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

		ShippingMethodReference *DeliveryInfo
		PaymentReference        *PaymentInfo
		CurrencyCode            string
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
