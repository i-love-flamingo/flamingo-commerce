package cart

import (
	"fmt"

	"github.com/gorilla/sessions"

	"time"

	"math"

	"flamingo.me/flamingo/framework/flamingo"
	"github.com/pkg/errors"
)

type (
	//CartProvider should be used to create the cart Value objects
	CartProvider func() *Cart

	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		//ID is the main identifier of the cart
		ID string
		//EntityID is a second identifier that may be used by some backends
		EntityID string

		//CartTotals - the cart totals (contain summary costs and discounts etc)
		CartTotals CartTotals
		//BillingAdress - the main billing address (relevant for all payments/invoices)
		BillingAdress Address

		//Purchaser - additional infos for the legal contact person in this order
		Purchaser Person

		//Deliveries - list of desired Deliverys (or Shippments) involved in this cart
		Deliveries []Delivery

		//AdditionalData   can be used for Custom attributes
		AdditionalData map[string]string

		//BelongsToAuthenticatedUser - false = Guest Cart true = cart from the authenticated user
		BelongsToAuthenticatedUser bool
		AuthenticatedUserId        string

		AppliedCouponCodes []CouponCode
	}

	// CouponCode value object
	CouponCode struct {
		Code string
	}

	// Person value object
	Person struct {
		Address         *Address
		PersonalDetails PersonalDetails
		//ExistingCustomerData if the current purchaser is an existing customer - this contains infos about existing customer
		ExistingCustomerData *ExistingCustomerData
	}

	// ExistingCustomerData value object
	ExistingCustomerData struct {
		//ID of the customer
		ID string
	}

	// PersonalDetails value object
	PersonalDetails struct {
		DateOfBirth     string
		PassportCountry string
		PassportNumber  string
		Nationality     string
	}

	// Delivery - represents the DeliveryInfo and the assigned Items
	Delivery struct {
		DeliveryInfo DeliveryInfo
		//Cartitems - list of cartitems
		Cartitems      []Item
		DeliveryTotals DeliveryTotals
	}

	// DeliveryInfo - represents the Delivery
	DeliveryInfo struct {
		Code             string
		Method           string
		Carrier          string
		DeliveryLocation DeliveryLocation
		ShippingItem     ShippingItem
		DesiredTime      time.Time
		AdditionalData   map[string]string
		RelatedFlight    *FlightData
	}

	// FlightData value object
	FlightData struct {
		ScheduledDateTime  time.Time
		Direction          string
		FlightNumber       string
		AirportName        string
		DestinationCountry string
		Terminal           string
		AirlineName        string
		Logger             flamingo.Logger
	}

	// DeliveryLocation value object
	DeliveryLocation struct {
		Type string
		//Address - only set for type adress
		Address *Address
		//Code - optional idendifier of this location/destination - is used in special destination Types
		Code string
	}

	// CartTotals value object
	CartTotals struct {
		Totalitems        []Totalitem
		TotalShippingItem ShippingItem
		//Final sum that need to be payed: GrandTotal = SubTotal + TaxAmount - DiscountAmount + SOME of Totalitems = (Sum of Items RowTotalWithDiscountInclTax) + SOME of Totalitems
		GrandTotal float64
		//SubTotal = SUM of Item RowTotal
		SubTotal float64
		//SubTotalInclTax = SUM of Item RowTotalInclTax
		SubTotalInclTax float64
		//SubTotalWithDiscounts = SubTotal - Sum of Item ItemRelatedDiscountAmount
		SubTotalWithDiscounts float64
		//SubTotalWithDiscountsAndTax= Sum of RowTotalWithItemRelatedDiscountInclTax
		SubTotalWithDiscountsAndTax float64

		//TotalDiscountAmount = SUM of Item TotalDiscountAmount
		TotalDiscountAmount float64
		//TotalNonItemRelatedDiscountAmount= SUM of Item NonItemRelatedDiscountAmount
		TotalNonItemRelatedDiscountAmount float64
		//TaxAmount = Sum of Item TaxAmount
		TaxAmount float64
		//CurrencyCode of the Total positions
		CurrencyCode string
	}

	// DeliveryTotals value object
	DeliveryTotals struct {
		//SubTotal = SUM of Item RowTotal
		SubTotal float64
		//SubTotalInclTax = SUM of Item RowTotalInclTax
		SubTotalInclTax float64
		//SubTotalWithDiscounts = SubTotal - Sum of Item ItemRelatedDiscountAmount
		SubTotalWithDiscounts float64
		//SubTotalWithDiscountsAndTax= Sum of RowTotalWithItemRelatedDiscountInclTax
		SubTotalWithDiscountsAndTax float64

		//TotalDiscountAmount = SUM of Item TotalDiscountAmount
		TotalDiscountAmount float64
		//TotalNonItemRelatedDiscountAmount= SUM of Item NonItemRelatedDiscountAmount
		TotalNonItemRelatedDiscountAmount float64
		//CurrencyCode of the Total positions
		CurrencyCode string
	}

	// Item for Cart
	Item struct {
		//ID of the item - need to be unique under a delivery
		ID string
		//
		UniqueId        string
		MarketplaceCode string
		//VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceId string

		Qty int

		CurrencyCode string

		AdditionalData map[string]string
		//brutto for single item
		SinglePrice float64
		//netto for single item
		SinglePriceInclTax float64
		//RowTotal = SinglePrice * Qty
		RowTotal float64
		//TaxAmount=Qty * (SinglePriceInclTax-SinglePrice)
		TaxAmount float64
		//RowTotalInclTax= RowTotal + TaxAmount
		RowTotalInclTax float64
		//AppliedDiscounts contains the details about the discounts applied to this item - they can be "itemrelated" or not
		AppliedDiscounts []ItemDiscount
		// TotalDiscountAmount = Sum of AppliedDiscounts = ItemRelatedDiscountAmount +NonItemRelatedDiscountAmount
		TotalDiscountAmount float64
		// ItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = True
		ItemRelatedDiscountAmount float64
		//NonItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = false
		NonItemRelatedDiscountAmount float64
		//RowTotalWithItemRelatedDiscountInclTax=RowTotal-ItemRelatedDiscountAmount
		RowTotalWithItemRelatedDiscount float64
		//RowTotalWithItemRelatedDiscountInclTax=RowTotalInclTax-ItemRelatedDiscountAmount
		RowTotalWithItemRelatedDiscountInclTax float64
		//This is the price the customer finaly need to pay for this item:  RowTotalWithDiscountInclTax = RowTotalInclTax-TotalDiscountAmount
		RowTotalWithDiscountInclTax float64
	}

	//ItemCartReference - value object that can be used to reference a Item in a Cart
	//@todo - Use in ServiePort methods...
	ItemCartReference struct {
		ItemId       string
		DeliveryCode string
	}

	// ItemDiscount value object
	ItemDiscount struct {
		Code  string
		Title string
		Price float64
		//IsItemRelated is a flag indicating if the discount should be displayed in the item or if it the result of a cart discount
		IsItemRelated bool
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price float64
		Type  string
	}

	// ShippingItem value object
	ShippingItem struct {
		Title string
		Price float64

		TaxAmount      float64
		DiscountAmount float64

		CurrencyCode string
	}

	// InvalidateCartEvent value object
	InvalidateCartEvent struct {
		Session *sessions.Session
	}
)

// Key constants
const (
	DELIVERY_METHOD_PICKUP      = "pickup"
	DELIVERY_METHOD_DELIVERY    = "delivery"
	DELIVERY_METHOD_UNSPECIFIED = "unspecified"

	DELIVERYLOCATION_TYPE_COLLECTIONPOINT = "collection-point"
	DELIVERYLOCATION_TYPE_STORE           = "store"
	DELIVERYLOCATION_TYPE_ADDRESS         = "address"
	DELIVERYLOCATION_TYPE_FREIGHTSTATION  = "freight-station"

	TOTALS_TYPE_DISCOUNT      = "totals_type_discount"
	TOTALS_TYPE_VOUCHER       = "totals_type_voucher"
	TOTALS_TYPE_TAX           = "totals_type_tax"
	TOTALS_TYPE_LOYALTYPOINTS = "totals_loyaltypoints"
	TOTALS_TYPE_SHIPPING      = "totals_type_shipping"
)

// GetMainShippingEMail returns the main shipping address email, empty string if not available
func (Cart Cart) GetMainShippingEMail() string {
	for _, deliveries := range Cart.Deliveries {
		if deliveries.DeliveryInfo.DeliveryLocation.Address != nil {
			if deliveries.DeliveryInfo.DeliveryLocation.Address.Email != "" {
				return deliveries.DeliveryInfo.DeliveryLocation.Address.Email
			}
		}
	}
	return ""
}

// GetDeliveryByCode gets a delivery by code
func (Cart Cart) GetDeliveryByCode(deliveryCode string) (*Delivery, bool) {
	for _, delivery := range Cart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			return &delivery, true
		}
	}
	return nil, false
}

// HasDeliveryForCode checks if a delivery with the given code exists in the cart
func (Cart Cart) HasDeliveryForCode(deliveryCode string) bool {
	_, found := Cart.GetDeliveryByCode(deliveryCode)
	return found == true
}

// GetDeliveryCodes returns a slice of all delivery codes in cart
func (Cart Cart) GetDeliveryCodes() []string {
	var deliveryCodes []string

	for deliveryIndex := range Cart.Deliveries {
		delivery := &Cart.Deliveries[deliveryIndex]
		if len(delivery.Cartitems) > 0 {
			deliveryCodes = append(deliveryCodes, delivery.DeliveryInfo.Code)
		}
	}

	return deliveryCodes
}

// GetByItemId gets an item by its id
func (Cart Cart) GetByItemId(itemId string, deliveryCode string) (*Item, error) {
	delivery, found := Cart.GetDeliveryByCode(deliveryCode)
	if found != true {
		return nil, errors.New(fmt.Sprintf("Delivery for code %v not found", deliveryCode))
	}
	for _, currentItem := range delivery.Cartitems {
		if currentItem.ID == itemId {
			return &currentItem, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("itemId %v in cart not existend", itemId))
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
	count := 0
	for _, delivery := range Cart.Deliveries {
		for _, item := range delivery.Cartitems {
			count += item.Qty
		}
	}

	return count
}

// GetItemCartReferences returns a slice of all ItemCartReferences
func (Cart Cart) GetItemCartReferences() []ItemCartReference {
	var ids []ItemCartReference
	for _, delivery := range Cart.Deliveries {
		for _, item := range delivery.Cartitems {
			ids = append(ids, ItemCartReference{
				ItemId:       item.ID,
				DeliveryCode: delivery.DeliveryInfo.Code,
			})
		}
	}
	return ids
}

// GetVoucherSavings returns the savings of all vouchers
func (Cart Cart) GetVoucherSavings() float64 {
	totalSavings := 0.0
	for _, item := range Cart.CartTotals.Totalitems {
		if item.Type == TOTALS_TYPE_VOUCHER {
			totalSavings = totalSavings + math.Abs(item.Price)
		}
	}

	if totalSavings < 0 {
		return 0.0
	}

	return totalSavings
}

// GetSavings retuns the total of all savings
func (Cart Cart) GetSavings() float64 {
	totalSavings := 0.0
	for _, item := range Cart.CartTotals.Totalitems {
		if item.Type == TOTALS_TYPE_DISCOUNT {
			totalSavings = totalSavings + math.Abs(item.Price)
		}
	}

	if totalSavings < 0 {
		return 0.0
	}

	return totalSavings
}

// HasAppliedCouponCode checks if a coupon code is applied to the cart
func (Cart Cart) HasAppliedCouponCode() bool {
	return len(Cart.AppliedCouponCodes) > 0
}

// GetTotalItemsByType gets a slice of all Totalitems by typeCode
func (ct CartTotals) GetTotalItemsByType(typeCode string) []Totalitem {
	var totalitems []Totalitem
	for _, item := range ct.Totalitems {
		if item.Type == typeCode {
			totalitems = append(totalitems, item)
		}
	}
	return totalitems
}

// GetSavingsByItem gets the savings by item
func (item Item) GetSavingsByItem() float64 {
	totalSavings := 0.0
	for _, discount := range item.AppliedDiscounts {
		totalSavings = totalSavings + math.Abs(discount.Price)
	}

	if totalSavings < 0 {
		return 0.0
	}

	return totalSavings
}

// HasRelatedFlight checks if a flight is related to the delivery
func (d DeliveryInfo) HasRelatedFlight() bool {
	return d.RelatedFlight != nil
}

// Inject dependencies
func (fd *FlightData) Inject(
	logger flamingo.Logger,
) {
	fd.Logger = logger
}

// GetScheduledDate returns the flights scheduled data
func (fd *FlightData) GetScheduledDate() string {
	return fd.ScheduledDateTime.Format("2006-01-02")
}

// GetScheduledDateTime returns the flights scheduled datetime
func (fd *FlightData) GetScheduledDateTime() string {
	return fd.ScheduledDateTime.Format(time.RFC3339)
}
