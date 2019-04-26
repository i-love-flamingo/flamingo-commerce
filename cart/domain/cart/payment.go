package cart

import (
	"errors"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
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
		AmountPayed domain.Price
		//ValuedPayed the value of the AmountPayed in the cart default currency
		ValuedAmountPayed domain.Price
		//CreditCardInfo Optional
		CreditCardInfo *CreditCardInfo
		//Title - speaking title - optional may describe the payment and may be shown to the customer
		Title string
		//RawTransactionData - place to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
		// ChargeAssignments - optional the assignment of this transaction to charges - this might be required for payments that are really only done
		ChargeAssignments *ChargeAssignments
	}

	// ChargeAssignments collects all assignments for item, shipment and total charges
	ChargeAssignments struct {
		ItemChargeAssignments     []ItemChargeAssignment
		TotalChargeAssignments    []TotalChargeAssignment
		ShipmentChargeAssignments []ShipmentChargeAssignment
	}
	//ChargeAssignments.GetSum() Charges

	// ChargeAssignmentsPerMethod collect charge assignments indexed by payment method
	ChargeAssignmentsPerMethod struct {
		perMethod map[string] ChargeAssignments
	}

	//CartChargeAssignment.GetForCartItem(itemId) map[string]Charge
	//CartChargeAssignment.GetForShippingItem(itemId) map[string]Charge
	//CartChargeAssignment.GetForTotalItem(itemId) map[string]Charge
	//CartChargeAssignment.GroupedSum() map[string]Charge


	// ItemChargeAssignment holds the information what amount was assigned to a specific chargetype of a specific item in the cart
	ItemChargeAssignment struct {
		ItemID string
		Charge domain.Charge
	}

	// TotalChargeAssignment holds the information what amount was assigned to a specific chargetype of a specific total in the cart
	TotalChargeAssignment struct {
		Type   string
		Charge domain.Charge
	}

	// ShipmentChargeAssignment holds the information what amount was assigned to a specific chargetype of a specific shipment in the cart
	ShipmentChargeAssignment struct {
		DeliveryCode string
		Charge       domain.Charge
	}

	// CreditCardInfo contains the necessary data
	CreditCardInfo struct {
		AnonymizedCardNumber string
		Type                 string
		CardHolder           string
		Expire               string
	}

	// PaymentSelection value object - that represents the payment selection on the cart
	PaymentSelection struct {
		//Gateway - the selected Gateway
		Gateway string
		//ChargeSplits - the selected split
		ChargeSplits []ChargeSplit
	}

	//ChargeSplit - amount by Method and chargetype
	ChargeSplit struct {
		//The Charge that need to be payed (with Amount, Value and Type in it)
		Charge domain.Charge

		//The selected payment method (code) that should be used
		Method string
    
		// ChargeAssignments - optional the assignment of this transaction to charges - this might be required for payments that are really only done
		ChargeAssignments *ChargeAssignments
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
func (cp *Payment) TotalValue() (domain.Price, error) {
	var prices []domain.Price

	for _, transaction := range cp.Transactions {
		prices = append(prices,transaction.ValuedAmountPayed)
	}

	return domain.SumAll(prices...)
}

//NewSimplePaymentSelection - returns a PaymentSelection that can be used to update the cart.
// 	multiple charges to pay the cart are not used here: The complete grandtotal is selected to be payed in one charge with the given paymentgateway and paymentmethod
func NewSimplePaymentSelection(gateway string, method string, grandTotal domain.Price) PaymentSelection {
	return PaymentSelection{
		Gateway: gateway,
		ChargeSplits: []ChargeSplit{
			{
				Charge: domain.Charge{
					Type:  domain.ChargeTypeMain,
					Price: grandTotal,
					Value: grandTotal,
				},
				Method: method,
			},
		},
	}
}

//IsSelected - returns true if a Gateway  is selected
func (s PaymentSelection) IsSelected() bool {
	return s.Gateway != ""
}

//GetCharges - sum per chargetype
func (s PaymentSelection) GetCharges() domain.Charges {
	result := domain.Charges{}

	for _, cs := range s.ChargeSplits {
		result = result.AddCharge(cs.Charge)
	}

	return result
}

//GetSelectedChargeAssignmentsPerMethod - returns the charge assignments per method
func (c Cart) GetSelectedChargeAssignmentsPerMethod() (*ChargeAssignmentsPerMethod, error) {
	if c.PaymentSelection == nil {
		return nil, errors.New("no payment selection")
	}

	if len(c.PaymentSelection.ChargeSplits) == 0 {
		return nil, errors.New("no chargesplit on selection")
	}

	chargeAssignmentsPerMethod := ChargeAssignmentsPerMethod{
		perMethod: make(map[string]ChargeAssignments),
	}

	if len(c.PaymentSelection.ChargeSplits) == 1 {
		chargeSplit := c.PaymentSelection.ChargeSplits[0]
		if chargeSplit.ChargeAssignments != nil {
			chargeAssignmentsPerMethod.perMethod[chargeSplit.Method] = *chargeSplit.ChargeAssignments
			return &chargeAssignmentsPerMethod, nil
		}

		return generateChargeAssignment(c)
	}

	for _, cs := range c.PaymentSelection.ChargeSplits {
		if cs.ChargeAssignments == nil {
			return nil, errors.New("no chargeassignments on paymentselection")
		}

		chargeAssignments := chargeAssignmentsPerMethod.perMethod[cs.Method]

		chargeAssignments.ItemChargeAssignments = append(chargeAssignments.ItemChargeAssignments, cs.ChargeAssignments.ItemChargeAssignments...)
		chargeAssignments.TotalChargeAssignments = append(chargeAssignments.TotalChargeAssignments, cs.ChargeAssignments.TotalChargeAssignments...)
		chargeAssignments.ShipmentChargeAssignments = append(chargeAssignments.ShipmentChargeAssignments, cs.ChargeAssignments.ShipmentChargeAssignments...)

		chargeAssignmentsPerMethod.perMethod[cs.Method] = chargeAssignments
	}

	return &chargeAssignmentsPerMethod, nil
}

// FilterByMethod returns Chargeassignments for a specific method
func (c *ChargeAssignmentsPerMethod) FilterByMethod(method string) ChargeAssignments {
	return c.perMethod[method]
}

func generateChargeAssignment(c Cart) (*ChargeAssignmentsPerMethod, error) {
	if c.PaymentSelection == nil {
		return nil, errors.New("no payment selection")
	}

	if len(c.PaymentSelection.ChargeSplits) != 1 {
		return nil, errors.New("too many chargesplits on paymentselections")
	}

	chargeSplit := c.PaymentSelection.ChargeSplits[0]

	if chargeSplit.Charge.Price.Currency() != chargeSplit.Charge.Value.Currency() {
		return nil, errors.New("currencies are different in charge - cannot generate")
	}

	chargeAssigment := ChargeAssignments{}

	for _, delivery := range c.Deliveries {
		if delivery.ShippingItem.TotalWithDiscountInclTax().Currency() != chargeSplit.Charge.Price.Currency() {
			return nil, errors.New("currencies are different in shipment - cannot generate")
		}

		shipmentCharge := domain.Charge{
			Price: delivery.ShippingItem.TotalWithDiscountInclTax(),
			Value: delivery.ShippingItem.TotalWithDiscountInclTax(),
			Type: domain.ChargeTypeMain,
		}

		chargeAssigment.ShipmentChargeAssignments = append(chargeAssigment.ShipmentChargeAssignments,ShipmentChargeAssignment{
			DeliveryCode: delivery.DeliveryInfo.Code,
			Charge:shipmentCharge,
		})

		for _, item := range delivery.Cartitems {
			if item.RowPriceGrossWithDiscount().Currency() != chargeSplit.Charge.Price.Currency() {
				return nil, errors.New("currencies are different in cart items - cannot generate")
			}

			itemCharge := domain.Charge{
				Price: item.RowPriceGrossWithDiscount(),
				Value: item.RowPriceGrossWithDiscount(),
				Type: domain.ChargeTypeMain,
			}

			chargeAssigment.ItemChargeAssignments = append(chargeAssigment.ItemChargeAssignments,ItemChargeAssignment{
				ItemID:item.ID,
				Charge:itemCharge,
			})
		}
	}

	chargeAssignmentsPerMethod := ChargeAssignmentsPerMethod{
		perMethod: make(map[string]ChargeAssignments),
	}

	chargeAssignmentsPerMethod.perMethod[chargeSplit.Method] = chargeAssigment

	return &chargeAssignmentsPerMethod, nil
}

// TotalValue - returns the Total Valued Price
func (s PaymentSelection) TotalValue() (domain.Price, error) {
	var prices []domain.Price
	for _, charge := range s.ChargeSplits {
		prices = append(prices,charge.Charge.Value)
	}

	return domain.SumAll(prices...)
}
