package cart

type (

	//CartPayment represents the all payments done for the cart and which items have been purchased by what methhod
	CartPayment struct {
		PaymentInfos       []*PaymentInfo
		Assignments        []CartPaymentAssignment
		RawTransactionData interface{}
	}

	//CartPaymentAssignment - represents the infos required to reference a CartItem
	CartPaymentAssignment struct {
		ItemCartReference ItemCartReference
		PaymentInfo       *PaymentInfo
	}

	PaymentInfo struct {
		//Provider code like "paymark"
		Provider string
		//Method like "paymark_cc" , "paypal",
		Method string
		//Status - Method specific status e.g. Auth, Captured, Open, ...
		Status string
		//TransactionId - The main reference of the payment that was done
		TransactionId string
		//AdditionalData - room for AdditionalData - specific to the payment
		AdditionalData map[string]string
		//Amount optional the amount payed
		Amount float64
		//CreditCardInfo Optional
		CreditCardInfo *CreditCardInfo
		//Title - speaking title - optional may describe the payment and may be shown to the customer
		Title string
	}

	CreditCardInfo struct {
		AnonymizedCardNumber string
		Type                 string
		CardHolder           string
		Expire               string
	}
)

const (
	PAYMENT_STATUS_CAPTURED = "CAPTURED"
	PAYMENT_STATUS_OPEN     = "OPEN"
)

func (cp *CartPayment) AddPayment(paymentInfo PaymentInfo, cartItemReferences []ItemCartReference) {
	cp.PaymentInfos = append(cp.PaymentInfos, &paymentInfo)

	for _, cartItemReference := range cartItemReferences {
		cp.Assignments = append(cp.Assignments, CartPaymentAssignment{
			ItemCartReference: cartItemReference,
			PaymentInfo:       &paymentInfo,
		})
	}
}

//GetAssignmentsForPaymentInfo - returns the CartItemReferences that are payed with the given paymentInfo
func (cp *CartPayment) GetAssignmentsForPaymentInfo(paymentInfo *PaymentInfo) []ItemCartReference {
	var ids []ItemCartReference
	for _, v := range cp.Assignments {
		if v.PaymentInfo == paymentInfo {
			ids = append(ids, v.ItemCartReference)
		}
	}
	return ids
}

//GetProviders - gets (deduplicated) list of PaymentProvider names used in the CartPayment
func (cp *CartPayment) GetProviders() []string {

	providerMap := make(map[string]bool)

	for _, info := range cp.PaymentInfos {
		providerMap[info.Provider] = true
	}

	var providers []string

	for provider, _ := range providerMap {
		providers = append(providers, provider)
	}

	return providers
}
