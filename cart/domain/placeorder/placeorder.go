package placeorder

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
	oauth "flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	// Service  interface - Secondary PORT
	Service interface {
		PlaceGuestCart(ctx context.Context, cart *cart.Cart, payment *Payment) (PlacedOrderInfos, error)
		PlaceCustomerCart(ctx context.Context, auth oauth.Auth, cart *cart.Cart, payment *Payment) (PlacedOrderInfos, error)
		ReserveOrderID(ctx context.Context, cart *cart.Cart) (string, error)
		// CancelOrder cancels an previously placed order and returns the used cart
		CancelOrder(ctx context.Context, orderInfos PlacedOrderInfos) error
	}
	// Payment represents all payments done for the cart and which items have been purchased by what method
	Payment struct {
		//The name of the Gateway that has returned the Payment for the cart
		Gateway string
		//Transactions - the list of individual transactions -  most cases only one Transaction might be part of the payment
		Transactions []Transaction
		//RawTransactionData - place to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
		//PaymentID - optional a reference of the Payment (that contains the Transactions)
		PaymentID string
	}

	// Transaction representing the transaction
	Transaction struct {
		//PaymentProvider - optional - the underling processor of this transaction (e.g. "paymark")
		PaymentProvider string
		//Method like "paymark_cc" , "paypal",
		Method string
		//Status - Method specific status e.g. Auth, Captured, Open, ...
		Status string
		//TransactionID - The main reference of the payment that was done
		TransactionID string
		//AdditionalData - room for AdditionalData - specific to the payment
		AdditionalData map[string]string
		//AmountPayed the amount that have been payed - eventually in a different currency
		AmountPayed price.Price
		//ValuedPayed the value of the AmountPayed in the cart default currency
		ValuedAmountPayed price.Price
		//CreditCardInfo Optional
		CreditCardInfo *CreditCardInfo
		//Title - speaking title - optional may describe the payment and may be shown to the customer
		Title string
		//RawTransactionData - place to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
		// ChargeAssignments - optional the assignment of this transaction to charges - this might be required for payments that are really only done for a certain item
		ChargeByItem *ChargeByItem
	}

	//ChargeByItem - the Charge that is payed for the individual items
	ChargeByItem struct {
		cartItems     map[string]price.Charge
		shippingItems map[string]price.Charge
		totalItems    map[string]price.Charge
	}

	// CreditCardInfo contains the necessary data
	CreditCardInfo struct {
		AnonymizedCardNumber string
		Type                 string
		CardHolder           string
		Expire               string
	}

	// PlacedOrderInfos represents a slice of PlacedOrderInfo
	PlacedOrderInfos []PlacedOrderInfo

	// PlacedOrderInfo defines the additional info struct for placed orders
	PlacedOrderInfo struct {
		OrderNumber  string
		DeliveryCode string
	}
)

const (
	// PaymentStatusCaptured a payment which has been captured
	PaymentStatusCaptured = "CAPTURED"
	// PaymentStatusAuthorized a payment which has been AUTHORIZED
	PaymentStatusAuthorized = "AUTHORIZED"
	// PaymentStatusOpen payment is still open
	PaymentStatusOpen = "OPEN"
)

// AddTransaction for a paymentInfo with items
func (cp *Payment) AddTransaction(transaction Transaction) {
	cp.Transactions = append(cp.Transactions, transaction)
}

// TotalValue - returns the Total Valued Price
func (cp *Payment) TotalValue() (price.Price, error) {
	var prices []price.Price

	for _, transaction := range cp.Transactions {
		prices = append(prices, transaction.ValuedAmountPayed)
	}

	return price.SumAll(prices...)
}

// GetOrderNumberForDeliveryCode returns the order number for a delivery code
func (poi PlacedOrderInfos) GetOrderNumberForDeliveryCode(deliveryCode string) string {
	for _, v := range poi {
		if v.DeliveryCode == deliveryCode {
			return v.OrderNumber
		}
	}

	return ""
}

//CartItems - return CartItems
func (c ChargeByItem) CartItems() map[string]price.Charge {
	return c.cartItems
}

//ChargeForCartItem returns Charge for a cartitem by id
func (c ChargeByItem) ChargeForCartItem(itemid string) (*price.Charge, bool) {
	if charge, ok := c.cartItems[itemid]; ok {
		return &charge, true
	}
	return nil, false
}

//ChargeForDeliveryCode returns Charge for a shippingitem by deliverycode
func (c ChargeByItem) ChargeForDeliveryCode(itemid string) (*price.Charge, bool) {
	if charge, ok := c.shippingItems[itemid]; ok {
		return &charge, true
	}
	return nil, false
}

//ChargeForTotal returns Charge for a totalitem by code
func (c ChargeByItem) ChargeForTotal(itemid string) (*price.Charge, bool) {
	if charge, ok := c.totalItems[itemid]; ok {
		return &charge, true
	}
	return nil, false
}

//TotalItems - returns totalItems
func (c ChargeByItem) TotalItems() map[string]price.Charge {
	return c.totalItems
}

//ShippingItems - returns ShippingItems
func (c ChargeByItem) ShippingItems() map[string]price.Charge {
	return c.shippingItems
}

//AddCartItem - modifies the current instance and adds a charge for an item
func (c ChargeByItem) AddCartItem(id string, charge price.Charge) ChargeByItem {
	if c.cartItems == nil {
		c.cartItems = make(map[string]price.Charge)
	}
	c.cartItems[id] = charge
	return c
}

//AddTotalItem - modifies the current instance and adds a charge for an item
func (c ChargeByItem) AddTotalItem(id string, charge price.Charge) ChargeByItem {
	if c.totalItems == nil {
		c.totalItems = make(map[string]price.Charge)
	}
	c.totalItems[id] = charge
	return c
}

//AddShippingItems - modifies the current instance and adds a charge for an item
func (c ChargeByItem) AddShippingItems(id string, charge price.Charge) ChargeByItem {
	if c.shippingItems == nil {
		c.shippingItems = make(map[string]price.Charge)
	}
	c.shippingItems[id] = charge
	return c
}
