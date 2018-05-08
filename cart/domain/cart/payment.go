package cart

type (

	//CartPayment represents the all payments done for the cart and which items have been purchased by what methhod
	CartPayment struct {
		PaymentInfos       []*PaymentInfo
		ItemIDAssignment   map[string]*PaymentInfo
		RawTransactionData interface{}
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

func (cp *CartPayment) AddPayment(paymentInfo PaymentInfo, items []string) {
	cp.PaymentInfos = append(cp.PaymentInfos, &paymentInfo)
	if cp.ItemIDAssignment == nil {
		cp.ItemIDAssignment = make(map[string]*PaymentInfo)
	}
	for _, v := range items {
		cp.ItemIDAssignment[v] = &paymentInfo
	}
}

func (cp *CartPayment) GetItemIdsForPaymentInfo(paymentInfo *PaymentInfo) []string {
	var ids []string
	for k, v := range cp.ItemIDAssignment {
		if v == paymentInfo {
			ids = append(ids, k)
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
