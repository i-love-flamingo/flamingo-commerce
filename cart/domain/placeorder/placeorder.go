package placeorder


import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
	oauth "flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	// PlaceOrderService  interface - Secondary PORT
	PlaceOrderService interface {
		PlaceGuestCart(ctx context.Context, cart *cart.Cart, payment *Payment) (PlacedOrderInfos, error)
		PlaceCustomerCart(ctx context.Context, auth oauth.Auth, cart *cart.Cart, payment *Payment) (PlacedOrderInfos, error)
		ReserveOrderID(ctx context.Context, cart *cart.Cart) (string, error)
	}
	// Payment represents all payments done for the cart and which items have been purchased by what method
	Payment struct {
		//The name of the Gateway that has returned the Payment for the cart
		Gateway string
		//Transactions - the list of individual transactions -  most cases only one Transaction might be part of the payment
		Transactions []Transaction
		//RawTransactionData - place to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
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
		// ChargeAssignments - optional the assignment of this transaction to charges - this might be required for payments that are really only done
		ChargeAssignments *cart.ChargeAssignments
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
		prices = append(prices,transaction.ValuedAmountPayed)
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