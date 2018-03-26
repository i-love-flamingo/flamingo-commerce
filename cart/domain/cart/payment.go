package cart

type (
	CartPayment struct {
		PaymentInfos   []PaymentInfo
		ItemAssignment map[string]*PaymentInfo
	}

	PaymentInfo struct {
		//Method like "paymark" , "paypal",
		Method string
		//Status - Method specific status e.g. Auth, Captured, Open, ...
		Status string
		//TransactionId - The main reference of the payment that was done
		TransactionId string
	}
)

func (cp *CartPayment) AddPayment(paymentInfo PaymentInfo, items []string) {
	cp.PaymentInfos = append(cp.PaymentInfos, paymentInfo)
	if cp.ItemAssignment == nil {
		cp.ItemAssignment = make(map[string]*PaymentInfo)
	}
	for _, v := range items {
		cp.ItemAssignment[v] = &paymentInfo
	}
}
