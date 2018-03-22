package cart

type (
	PaymentInfo struct {
		//Method like "paymark" , "paypal",
		Method string
		//Status - Method specific status e.g. Auth, Captured, Open, ...
		Status string
		//TransactionId - The main reference of the payment that was done
		TransactionId string
	}
)
