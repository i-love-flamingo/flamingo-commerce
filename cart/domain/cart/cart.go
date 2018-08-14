package cart

import (
	"fmt"

	"time"

	"math"

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

	CouponCode struct {
		Code string
	}

	Person struct {
		Address         *Address
		PersonalDetails PersonalDetails
		//ExistingCustomerData if the current purchaser is an existing customer - this contains infos about existing customer
		ExistingCustomerData *ExistingCustomerData
	}

	ExistingCustomerData struct {
		//ID of the customer
		ID string
	}

	PersonalDetails struct {
		DateOfBirth     string
		PassportCountry string
		PassportNumber  string
		Nationality     string
	}

	//Delivery - represents the DeliveryInfo and the assigned Items
	Delivery struct {
		DeliveryInfo DeliveryInfo
		//Cartitems - list of cartitems
		Cartitems []Item
	}

	//DeliveryInfo - represents the Delivery
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

	FlightData struct {
		ScheduledDateTime  time.Time
		Direction          string
		FlightNumber       string
		AirportName        string
		DestinationCountry string
	}

	DeliveryLocation struct {
		Type string
		//Address - only set for type adress
		Address *Address
		//Code - optional idendifier of this location/destination - is used in special destination Types
		Code string
	}

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

		//DEPRICATED
		//DiscountAmount float64

		//TaxAmount = Sum of Item TaxAmount
		TaxAmount float64
		//CurrencyCode of the Total positions
		CurrencyCode string
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

	// DiscountItem
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

// GetMainShippingEMail
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

// GetByItemId gets an item by its id
func (Cart Cart) GetDeliveryByCode(deliveryCode string) (*Delivery, error) {
	for _, delivery := range Cart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			return &delivery, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Delivery with code %v in cart not existend", deliveryCode))
}

func (Cart Cart) HasDeliveryForCode(deliveryCode string) bool {
	for _, d := range Cart.Deliveries {
		if d.DeliveryInfo.Code == deliveryCode {
			return true
		}
	}
	return false
}

// GetByItemId gets an item by its id
func (Cart Cart) GetByItemId(itemId string, deliveryCode string) (*Item, error) {
	delivery, err := Cart.GetDeliveryByCode(deliveryCode)
	if err != nil {
		return nil, err
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

// check if it is a mixed cart with different delivery intents
//@todo - only non empty Deliveries should count
func (Cart Cart) HasMixedCart() bool {
	return len(Cart.Deliveries) > 1
}

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

func (Cart Cart) HasAppliedCouponCode() bool {
	return len(Cart.AppliedCouponCodes) > 0
}

func (ct CartTotals) GetTotalItemsByType(typeCode string) []Totalitem {
	var totalitems []Totalitem
	for _, item := range ct.Totalitems {
		if item.Type == typeCode {
			totalitems = append(totalitems, item)
		}
	}
	return totalitems
}

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

func (d DeliveryInfo) HasRelatedFlight() bool {
	return d.RelatedFlight != nil
}

func (fd *FlightData) GetScheduledDate() string {
	return fd.ScheduledDateTime.Format("2006-01-02")
}

func (fd *FlightData) GetScheduledDateTime() string {
	return fd.ScheduledDateTime.Format(time.RFC3339)
}

//GetScheduledDateTime string from ScheduledDateTime - used for display
func (f *FlightData) ParseScheduledDateTime() time.Time {
	//"scheduledDateTime": "2017-11-25T06:30:00Z",
	timeResult, e := time.Parse(time.RFC3339, f.GetScheduledDateTime())
	if e != nil {
		return time.Now()
	}
	return timeResult
}
