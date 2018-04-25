package cart

type (
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

func (cp *CartPayment) GetProviders() []string {
	var providers []string
	for _, info := range cp.PaymentInfos {
		providers = append(providers, info.Provider)
	}
	return providers
}
