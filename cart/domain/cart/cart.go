package cart

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// Provider should be used to create the cart Value objects
	provider func() *Cart

	// Cart Value Object (immutable data - because the CartService is responsible to return a cart).
	Cart struct {
		// ID is the main identifier of the cart
		ID string
		// EntityID is a second identifier that may be used by some backends
		EntityID string

		// BillingAddress is the main billing address (relevant for all payments/invoices)
		BillingAddress *Address

		// Purchaser hold additional infos for the legal contact person in this order
		Purchaser *Person

		// Deliveries contains a list of desired Deliveries (or Shipments) involved in this cart
		Deliveries []Delivery

		// AdditionalData can be used for Custom attributes
		AdditionalData AdditionalData

		// PaymentSelection is used to store information on "how" the customer wants to pay
		PaymentSelection PaymentSelection

		// BelongsToAuthenticatedUser displays if the cart is guest cart (false) or from an authenticated user (true)
		BelongsToAuthenticatedUser bool

		// AuthenticatedUserID holds the potential customer ID
		AuthenticatedUserID string

		// AppliedCouponCodes hold the coupons or discount codes that are applied to the cart
		AppliedCouponCodes []CouponCode

		DefaultCurrency string

		// Additional non taxable totals
		Totalitems []Totalitem

		// AppliedGiftCards is a list of applied gift cards
		AppliedGiftCards []AppliedGiftCard
	}

	// Teaser represents some teaser infos for cart
	Teaser struct {
		ProductCount  int
		ItemCount     int
		DeliveryCodes []string
	}

	// CouponCode value object
	CouponCode struct {
		Code string
		// CustomAttributes can hold additional data for coupon code - keys and values are project specific
		CustomAttributes map[string]interface{}
	}

	// AppliedCouponCodes slice of applied coupon codes
	AppliedCouponCodes []CouponCode

	// Person value object
	Person struct {
		Address         *Address
		PersonalDetails PersonalDetails
		// ExistingCustomerData if the current purchaser is an existing customer - this contains infos about existing customer
		ExistingCustomerData *ExistingCustomerData
	}

	// ExistingCustomerData value object
	ExistingCustomerData struct {
		// ID of the customer
		ID string
	}

	// PersonalDetails value object
	PersonalDetails struct {
		DateOfBirth     string
		PassportCountry string
		PassportNumber  string
		Nationality     string
	}

	// Taxes is a list of Tax
	Taxes []Tax

	// Tax is the Tax represented by an Amount and optional the Rate.
	Tax struct {
		Amount domain.Price
		Type   string
		Rate   *big.Float `swaggertype:"string"`
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price domain.Price
		Type  string
	}

	// InvalidateCartEvent value object
	InvalidateCartEvent struct {
		Session *web.Session
	}

	// AdditionalData defines the supplementary cart data
	AdditionalData struct {
		// CustomAttributes list of key values
		CustomAttributes map[string]string
		// ReservedOrderID is an ID already known by the Cart of the future order ID
		ReservedOrderID string
	}

	// Builder is the main builder for a cart
	Builder struct {
		cartInBuilding *Cart
	}

	// BuilderProvider should be used to create the cart by using the Builder
	BuilderProvider func() *Builder

	// PricedItems - value object that contains items of the different possible types, that have a price
	PricedItems struct {
		cartItems     map[string]domain.Price
		shippingItems map[string]domain.Price
		totalItems    map[string]domain.Price
	}
)

var (
	// ErrAdditionalInfosNotFound is returned if the additional infos are not set
	ErrAdditionalInfosNotFound = errors.New("additional infos not found")
)

// Key constants
const (
	TotalsTypeDiscount      = "totals_type_discount"
	TotalsTypeVoucher       = "totals_type_voucher"
	TotalsTypeTax           = "totals_type_tax"
	TotalsTypeLoyaltypoints = "totals_loyaltypoints"
	TotalsTypeShipping      = "totals_type_shipping"
)

func init() {
	gob.Register(Cart{})
	gob.Register(DefaultPaymentSelection{})
}

// GetMainShippingEMail returns the main shipping address email, empty string if not available
func (c Cart) GetMainShippingEMail() string {
	for _, deliveries := range c.Deliveries {
		if deliveries.DeliveryInfo.DeliveryLocation.Address != nil {
			if deliveries.DeliveryInfo.DeliveryLocation.Address.Email != "" {
				return deliveries.DeliveryInfo.DeliveryLocation.Address.Email
			}
		}
	}

	return ""
}

// Clone the current cart
func (c Cart) Clone() (Cart, error) {
	cloned := Cart{}

	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(c)
	if err != nil {
		return Cart{}, err
	}
	err = gob.NewDecoder(b).Decode(&cloned)
	if err != nil {
		return Cart{}, err
	}

	return cloned, nil
}

// GetContactMail returns the contact mail from the shipping address with fall back to the billing address
func (c Cart) GetContactMail() string {
	// Get Email from either the cart
	shippingEmail := c.GetMainShippingEMail()
	if shippingEmail == "" && c.BillingAddress != nil {
		shippingEmail = c.BillingAddress.Email
	}

	return shippingEmail
}

// IsEmpty returns true if cart is empty
func (c Cart) IsEmpty() bool {
	return c.GetCartTeaser().ItemCount == 0
}

// GetDeliveryByCode gets a delivery by code
func (c Cart) GetDeliveryByCode(deliveryCode string) (*Delivery, bool) {
	for _, delivery := range c.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			return &delivery, true
		}
	}

	return nil, false
}

// GetDeliveryByCodeWithoutBool TODO: This func needs to be removed as soon as there's a solution for handling of boolean returns when gql expects an err
func (c Cart) GetDeliveryByCodeWithoutBool(deliveryCode string) *Delivery {
	delivery, _ := c.GetDeliveryByCode(deliveryCode)
	return delivery
}

// HasDeliveryForCode checks if a delivery with the given code exists in the cart
func (c Cart) HasDeliveryForCode(deliveryCode string) bool {
	_, found := c.GetDeliveryByCode(deliveryCode)

	return found == true
}

// GetDeliveryCodes returns a slice of all delivery codes in cart that have at least one cart item
func (c Cart) GetDeliveryCodes() []string {
	deliveryCodes := make([]string, 0, len(c.Deliveries))

	for _, delivery := range c.Deliveries {
		if len(delivery.Cartitems) > 0 {
			deliveryCodes = append(deliveryCodes, delivery.DeliveryInfo.Code)
		}
	}

	return deliveryCodes
}

// GetDeliveryByItemID returns a delivery by a given itemID
func (c Cart) GetDeliveryByItemID(itemID string) (*Delivery, error) {
	for _, delivery := range c.Deliveries {
		for _, cartItem := range delivery.Cartitems {
			if cartItem.ID == itemID {
				return &delivery, nil
			}
		}
	}

	return nil, errors.Errorf("delivery not found for %q", itemID)
}

// GetByItemID gets an item by its id
func (c Cart) GetByItemID(itemID string) (*Item, error) {
	for _, delivery := range c.Deliveries {
		for _, currentItem := range delivery.Cartitems {
			if currentItem.ID == itemID {
				return &currentItem, nil
			}
		}
	}

	return nil, errors.Errorf("itemId %q in cart does not exist", itemID)
}

// GetTotalQty for the product in the cart
func (c Cart) GetTotalQty(marketPlaceCode string, variantCode string) int {
	qty := int(0)
	for _, delivery := range c.Deliveries {
		for _, currentItem := range delivery.Cartitems {
			if currentItem.MarketplaceCode == marketPlaceCode && currentItem.VariantMarketPlaceCode == variantCode {
				qty = qty + currentItem.Qty
			}
		}
	}
	return qty
}

// GetByExternalReference gets an item by its external reference
func (c Cart) GetByExternalReference(ref string) (*Item, error) {
	for _, delivery := range c.Deliveries {
		for _, currentItem := range delivery.Cartitems {
			if currentItem.ExternalReference == ref {
				return &currentItem, nil
			}
		}
	}

	return nil, errors.Errorf("uitemID %v in cart not existing", ref)
}

// ItemCount returns amount of cart items in the current cart
func (c Cart) ItemCount() int {
	count := 0
	for _, delivery := range c.Deliveries {
		for _, item := range delivery.Cartitems {
			count += item.Qty
		}
	}

	return count
}

// ProductCount returns the amount of different products, but a product is counted twice if it is in different deliveries
func (c Cart) ProductCount() int {
	count := 0
	for _, delivery := range c.Deliveries {
		count += len(delivery.Cartitems)
	}

	return count
}

// ProductCountUnique returns the amount of unique products across all deliveries
func (c Cart) ProductCountUnique() int {
	marketplaceCodes := make(map[string]struct{})
	for _, delivery := range c.Deliveries {
		for _, item := range delivery.Cartitems {
			if _, ok := marketplaceCodes[item.MarketplaceCode]; !ok {
				marketplaceCodes[item.MarketplaceCode] = struct{}{}
			}
		}
	}
	return len(marketplaceCodes)
}

// IsPaymentSelected returns true if a valid payment is selected
func (c Cart) IsPaymentSelected() bool {
	return c.PaymentSelection != nil
}

// GetVoucherSavings returns the savings of all vouchers
func (c Cart) GetVoucherSavings() domain.Price {
	price := domain.Price{}
	var err error

	for _, item := range c.Totalitems {
		if item.Type == TotalsTypeVoucher {
			price, err = price.Add(item.Price)
			if err != nil {
				return price
			}
		}
	}

	if price.IsNegative() {
		return domain.Price{}
	}

	return price
}

// GrandTotal is the final sum (Valued price) that need to be paid: GrandTotal = SubTotal + TaxAmount - DiscountAmount + sum of Totalitems = (Sum of Items RowTotalWithDiscountInclTax) + sum of Totalitems
func (c Cart) GrandTotal() domain.Price {
	return c.GetAllPaymentRequiredItems().Sum()
}

// GetAllPaymentRequiredItems  returns all the Items (Cartitem, ShippingItem, TotalItems) that need to be paid with the final gross price
func (c Cart) GetAllPaymentRequiredItems() PricedItems {
	pricedItems := PricedItems{
		cartItems:     make(map[string]domain.Price),
		shippingItems: make(map[string]domain.Price, len(c.Deliveries)),
		totalItems:    make(map[string]domain.Price, len(c.Totalitems)),
	}
	for _, ti := range c.Totalitems {
		pricedItems.totalItems[ti.Code] = ti.Price
	}
	for _, del := range c.Deliveries {
		if !del.ShippingItem.TotalWithDiscountInclTax().IsZero() {
			pricedItems.shippingItems[del.DeliveryInfo.Code] = del.ShippingItem.TotalWithDiscountInclTax()
		}
		for _, ci := range del.Cartitems {
			pricedItems.cartItems[ci.ID] = ci.RowPriceGrossWithDiscount()
		}
	}
	return pricedItems
}

// SumShippingNet returns net sum price of deliveries ShippingItems
func (c Cart) SumShippingNet() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.ShippingItem.PriceNet)
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SumShippingNetWithDiscounts returns net sum price of deliveries ShippingItems with discounts
func (c Cart) SumShippingNetWithDiscounts() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		discount, _ := del.ShippingItem.AppliedDiscounts.Sum()
		prices = append(prices, del.ShippingItem.PriceNet.ForceAdd(discount))
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SumShippingGross returns gross sum price of deliveries ShippingItems
func (c Cart) SumShippingGross() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.ShippingItem.PriceGross)
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SumShippingGrossWithDiscounts returns gross sum price of deliveries ShippingItems with discounts
func (c Cart) SumShippingGrossWithDiscounts() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.ShippingItem.TotalWithDiscountInclTax())
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// HasShippingCosts returns true if cart HasShippingCosts
func (c Cart) HasShippingCosts() bool {
	return c.SumShippingNet().IsPositive()
}

// AllShippingTitles returns all ShippingItem titles
func (c Cart) AllShippingTitles() []string {
	label := make([]string, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		label = append(label, del.ShippingItem.Title)
	}

	return label
}

// SubTotalGross returns sum price of deliveries SubTotalGross
func (c Cart) SubTotalGross() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.SubTotalGross())
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SumTaxes returns sum of deliveries SumRowTaxes
func (c Cart) SumTaxes() Taxes {
	newTaxes := Taxes{}

	for _, del := range c.Deliveries {
		newTaxes = newTaxes.AddTaxesWithMerge(del.SumRowTaxes())
		if !del.ShippingItem.TaxAmount.IsZero() {
			newTaxes = newTaxes.AddTax(del.ShippingItem.Tax())
		}
	}

	return newTaxes
}

// SumTotalTaxAmount returns sum price of deliveries Taxes as total amount
func (c Cart) SumTotalTaxAmount() domain.Price {
	return c.SumTaxes().TotalAmount()
}

// SubTotalNet returns sum price of deliveries SubTotalNet
func (c Cart) SubTotalNet() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.SubTotalNet())
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SubTotalGrossWithDiscounts returns sum price of deliveries SubTotalGrossWithDiscounts
func (c Cart) SubTotalGrossWithDiscounts() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.SubTotalGrossWithDiscounts())
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SubTotalNetWithDiscounts returns sum price of deliveries SubTotalNetWithDiscounts
func (c Cart) SubTotalNetWithDiscounts() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.SubTotalNetWithDiscounts())
	}
	price, _ := domain.SumAll(prices...)

	return price
}

// SumTotalDiscountAmount returns sum price of deliveries SumTotalDiscountAmount
func (c Cart) SumTotalDiscountAmount() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.SumTotalDiscountAmount())
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SumNonItemRelatedDiscountAmount returns sum price of deliveries SumNonItemRelatedDiscountAmount
func (c Cart) SumNonItemRelatedDiscountAmount() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.SumNonItemRelatedDiscountAmount())
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// SumItemRelatedDiscountAmount returns sum price of deliveries SumItemRelatedDiscountAmount
func (c Cart) SumItemRelatedDiscountAmount() domain.Price {
	prices := make([]domain.Price, 0, len(c.Deliveries))

	for _, del := range c.Deliveries {
		prices = append(prices, del.SumItemRelatedDiscountAmount())
	}

	price, _ := domain.SumAll(prices...)

	return price
}

// HasAppliedCouponCode checks if a coupon code is applied to the cart
func (c Cart) HasAppliedCouponCode() bool {
	return len(c.AppliedCouponCodes) > 0
}

// GetCartTeaser returns the teaser
func (c Cart) GetCartTeaser() *Teaser {
	return &Teaser{
		DeliveryCodes: c.GetDeliveryCodes(),
		ItemCount:     c.ItemCount(),
		ProductCount:  c.ProductCount(),
	}
}

// GetPaymentReference returns a string that can be used as reference to pass to payment gateway. You may want to use it. It returns either the reserved Order id or the cart id/entityid
func (c Cart) GetPaymentReference() string {
	if c.AdditionalData.ReservedOrderID != "" {
		return c.AdditionalData.ReservedOrderID
	}
	return fmt.Sprintf("%v-%v", c.ID, c.EntityID)
}

// GetTotalItemsByType returns a slice of all Totalitems by typeCode
func (c Cart) GetTotalItemsByType(typeCode string) []Totalitem {
	totalitems := make([]Totalitem, 0, len(c.Totalitems))

	for _, item := range c.Totalitems {
		if item.Type == typeCode {
			totalitems = append(totalitems, item)
		}
	}

	return totalitems
}

// GrandTotalCharges is the final sum that need to be paid - split by the charges that need to be paid
func (c Cart) GrandTotalCharges() domain.Charges {
	// Check if a specific split was saved:
	if c.PaymentSelection != nil {
		charges := c.PaymentSelection.CartSplit().ChargesByType()
		// make sure we have valid main charge currency
		return charges.AddCharge(domain.Charge{
			Value: domain.NewFromInt(0, 1, c.DefaultCurrency),
			Price: domain.NewFromInt(0, 1, c.DefaultCurrency),
			Type:  domain.ChargeTypeMain,
		})
	}

	// else return the grand total as main charge
	charges := domain.Charges{}
	mainCharge := domain.Charge{
		Value: c.GrandTotal(),
		Price: c.GrandTotal(),
		Type:  domain.ChargeTypeMain,
	}

	charges = charges.AddCharge(mainCharge)

	return charges
}

// AddTax returns new Tax with this Tax added
func (t Taxes) AddTax(tax Tax) Taxes {
	newTaxes := Taxes(t)
	newTaxes = append(newTaxes, tax)

	return newTaxes
}

// AddTaxWithMerge returns new Taxes with this Tax added
func (t Taxes) AddTaxWithMerge(taxToAddOrMerge Tax) Taxes {
	newTaxes := Taxes(t)

	for k, tax := range newTaxes {
		if tax.Type == taxToAddOrMerge.Type {
			if tax.Rate == nil && taxToAddOrMerge.Rate == nil {
				newTaxes[k].Amount, _ = tax.Amount.Add(taxToAddOrMerge.Amount)
				return newTaxes
			} else if tax.Rate != nil && taxToAddOrMerge.Rate != nil && (tax.Rate.Cmp(taxToAddOrMerge.Rate) == 0) {
				newTaxes[k].Amount, _ = tax.Amount.Add(taxToAddOrMerge.Amount)
				return newTaxes
			}
		}
	}

	newTaxes = newTaxes.AddTax(taxToAddOrMerge)

	return newTaxes
}

// AddTaxesWithMerge returns new Taxes with the given Taxes all added or merged in
func (t Taxes) AddTaxesWithMerge(taxes Taxes) Taxes {
	newTaxes := Taxes(t)

	for _, tax := range taxes {
		newTaxes = newTaxes.AddTaxWithMerge(tax)
	}

	return newTaxes
}

// TotalAmount returns the sum of all taxes as price
func (t Taxes) TotalAmount() domain.Price {
	prices := make([]domain.Price, 0, len(t))

	for _, tax := range t {
		prices = append(prices, tax.Amount)
	}

	result, _ := domain.SumAll(prices...)

	return result
}

// Build is the main factory method
func (b *Builder) Build() (*Cart, error) {
	return b.reset(nil)
}

// SetIds sets the cart ids
func (b *Builder) SetIds(id string, entityID string) *Builder {
	b.init()
	b.cartInBuilding.ID = id
	b.cartInBuilding.EntityID = entityID

	return b
}

// SetReservedOrderID sets the optional order id
func (b *Builder) SetReservedOrderID(id string) *Builder {
	b.init()
	b.cartInBuilding.AdditionalData.ReservedOrderID = id

	return b
}

// SetBillingAddress sets the optional billing address
func (b *Builder) SetBillingAddress(a Address) *Builder {
	b.init()
	b.cartInBuilding.BillingAddress = &a

	return b
}

// SetPurchaser sets the optional purchaser
func (b *Builder) SetPurchaser(p Person) *Builder {
	b.init()
	b.cartInBuilding.Purchaser = &p

	return b
}

// AddDelivery adds a delivery - use the DeliveryBuilder to create one
// todo Make sure that item id is unique over the whole cart (that is an invariant so we need to ensure that invalid cart objects cannot be build)
func (b *Builder) AddDelivery(d Delivery) *Builder {
	b.init()
	b.cartInBuilding.Deliveries = append(b.cartInBuilding.Deliveries, d)

	return b
}

// SetAdditionalData can be used to add additional data
func (b *Builder) SetAdditionalData(d AdditionalData) *Builder {
	b.init()
	b.cartInBuilding.AdditionalData = d

	return b
}

// SetPaymentSelection sets the payment selection
func (b *Builder) SetPaymentSelection(d PaymentSelection) *Builder {
	b.init()
	b.cartInBuilding.PaymentSelection = d

	return b
}

// SetAuthenticatedUserID used to mark the cart as authenticated users cart
func (b *Builder) SetAuthenticatedUserID(id string) *Builder {
	b.init()
	b.cartInBuilding.AuthenticatedUserID = id
	b.cartInBuilding.BelongsToAuthenticatedUser = true

	return b
}

// SetBelongsToAuthenticatedUser mark the cart as authenticated users cart
func (b *Builder) SetBelongsToAuthenticatedUser(v bool) *Builder {
	b.init()
	b.cartInBuilding.BelongsToAuthenticatedUser = v

	return b
}

// AddAppliedCouponCode add the optional coupon that is applied for the cart
func (b *Builder) AddAppliedCouponCode(code CouponCode) *Builder {
	b.init()
	b.cartInBuilding.AppliedCouponCodes = append(b.cartInBuilding.AppliedCouponCodes, code)

	return b
}

// SetAppliedGiftCards sets the optional applied gift cards
func (b *Builder) SetAppliedGiftCards(gc []AppliedGiftCard) *Builder {
	b.init()
	b.cartInBuilding.AppliedGiftCards = gc

	return b
}

// SetDefaultCurrency sets the default currency
func (b *Builder) SetDefaultCurrency(d string) *Builder {
	b.init()
	b.cartInBuilding.DefaultCurrency = d

	return b
}

// AddTotalitem adds nontaxable extra totals on cart level
func (b *Builder) AddTotalitem(totali Totalitem) *Builder {
	b.init()
	b.cartInBuilding.Totalitems = append(b.cartInBuilding.Totalitems, totali)

	return b
}

func (b *Builder) init() {
	if b.cartInBuilding == nil {
		b.cartInBuilding = &Cart{}
	}
}

func (b *Builder) reset(err error) (*Cart, error) {
	cart := b.cartInBuilding
	b.cartInBuilding = nil

	return cart, err
}

// Sum returns Sum of all items in this struct
func (p PricedItems) Sum() domain.Price {
	prices := make([]domain.Price, 0, len(p.cartItems)+len(p.shippingItems)+len(p.totalItems))

	for _, p := range p.totalItems {
		prices = append(prices, p)
	}
	for _, p := range p.shippingItems {
		prices = append(prices, p)
	}
	for _, p := range p.cartItems {
		prices = append(prices, p)
	}
	sum, _ := domain.SumAll(prices...)

	return sum
}

// TotalItems returns the Price per Totalitem - map key is total type
func (p PricedItems) TotalItems() map[string]domain.Price {
	return p.totalItems
}

// ShippingItems returns the Price per ShippingItems - map key is delivery code
func (p PricedItems) ShippingItems() map[string]domain.Price {
	return p.shippingItems
}

// CartItems returns the Price per cartItems - map key is cart item ID
func (p PricedItems) CartItems() map[string]domain.Price {
	return p.cartItems
}

// ContainedIn returns if the coupon codes are contained in couponCodesToCompare
func (acc AppliedCouponCodes) ContainedIn(couponCodesToCompare AppliedCouponCodes) bool {
	for _, couponCode := range acc {
		contained := false
		for _, couponCodeToCompare := range couponCodesToCompare {
			if couponCode.Code == couponCodeToCompare.Code {
				contained = true
			}
		}
		if !contained {
			return false
		}
	}
	return true
}
