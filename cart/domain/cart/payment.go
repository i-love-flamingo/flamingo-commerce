package cart

type (
	// Payment represents all payments done for the cart and which items have been purchased by what method
	Payment struct {
		PaymentInfos       []*PaymentInfo
		Assignments        []PaymentAssignment
		RawTransactionData interface{}
	}

	// PaymentAssignment - represents the infos required to reference a CartItem
	PaymentAssignment struct {
		ItemCartReference ItemCartReference
		PaymentInfo       *PaymentInfo
	}

	// PaymentInfo contains information about the used payment
	PaymentInfo struct {
		//Provider code like "paymark"
		Provider string
		//Method like "paymark_cc" , "paypal",
		Method string
		//Status - Method specific status e.g. Auth, Captured, Open, ...
		Status string
		//TransactionID - The main reference of the payment that was done
		TransactionID string
		//AdditionalData - room for AdditionalData - specific to the payment
		AdditionalData map[string]string
		//Amount optional the amount payed
		Amount float64
		//CreditCardInfo Optional
		CreditCardInfo *CreditCardInfo
		//Title - speaking title - optional may describe the payment and may be shown to the customer
		Title string
	}

	// CreditCardInfo contains the necessary data
	CreditCardInfo struct {
		AnonymizedCardNumber string
		Type                 string
		CardHolder           string
		Expire               string
	}
)

const (
	// PaymentStatusCaptured a payment which has been captured
	PaymentStatusCaptured = "CAPTURED"
	// PaymentStatusOpen payment is still open
	PaymentStatusOpen = "OPEN"
)

// AddPayment for a paymentInfo with items
func (cp *Payment) AddPayment(paymentInfo PaymentInfo, cartItemReferences []ItemCartReference) {
	cp.PaymentInfos = append(cp.PaymentInfos, &paymentInfo)

	for _, cartItemReference := range cartItemReferences {
		cp.Assignments = append(cp.Assignments, PaymentAssignment{
			ItemCartReference: cartItemReference,
			PaymentInfo:       &paymentInfo,
		})
	}
}

//GetAssignmentsForPaymentInfo - returns the CartItemReferences that are payed with the given paymentInfo
func (cp *Payment) GetAssignmentsForPaymentInfo(paymentInfo *PaymentInfo) []ItemCartReference {
	var ids []ItemCartReference
	for _, v := range cp.Assignments {
		if v.PaymentInfo == paymentInfo {
			ids = append(ids, v.ItemCartReference)
		}
	}
	return ids
}

//GetProviders - gets (deduplicated) list of PaymentProvider names used in the CartPayment
func (cp *Payment) GetProviders() []string {

	providerMap := make(map[string]bool)

	for _, info := range cp.PaymentInfos {
		providerMap[info.Provider] = true
	}

	var providers []string

	for provider := range providerMap {
		providers = append(providers, provider)
	}

	return providers
}
