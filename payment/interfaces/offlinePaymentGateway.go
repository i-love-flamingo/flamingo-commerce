package interfaces

import (
	"context"
	"errors"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/payment/domain"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

	"flamingo.me/flamingo/v3/framework/web"
)

// OfflineWebCartPaymentGateway provides an offline payment integration
type OfflineWebCartPaymentGateway struct {
	enabled   bool
	responder *web.Responder
}

const (
	//OfflineWebCartPaymentGatewayCode - the gateway code
	OfflineWebCartPaymentGatewayCode = "offline"
)

var _ WebCartPaymentGateway = (*OfflineWebCartPaymentGateway)(nil)

// Inject for OfflineWebCartPaymentGateway
func (o *OfflineWebCartPaymentGateway) Inject(responder *web.Responder, config *struct {
	Enabled bool `inject:"config:checkout.enableOfflinePaymentProvider,optional"`
}) {
	o.responder = responder
	if config != nil {
		o.enabled = config.Enabled
	}
}

// Methods returns the Payment Providers available Payment Methods
func (o *OfflineWebCartPaymentGateway) Methods() []domain.Method {
	return []domain.Method{{
		Title: "cash on delivery",
		Code:  "offlinepayment_cashondelivery",
	},
		{
			Title: "cash in advance",
			Code:  "offlinepayment_cashinadvance",
		}}
}

func (o *OfflineWebCartPaymentGateway) isSupportedMethod(method string) bool {
	for _, m := range o.Methods() {
		if m.Code == method {
			return true
		}
	}
	return false
}

func (o *OfflineWebCartPaymentGateway) checkCart(currentCart *cartDomain.Cart) error {
	//Read the infos in the cart and check precondition
	if currentCart.PaymentSelection.Gateway != OfflineWebCartPaymentGatewayCode {
		return errors.New("Cart is not supposed to be payed by this gateway")
	}
	for _, split := range currentCart.PaymentSelection.ChargeSplits {
		if !o.isSupportedMethod(split.Method) {
			return errors.New("Cart payment method not supported by gateway")
		}
	}
	return nil
}

//GetStartFlowResult for offline payment
func (o *OfflineWebCartPaymentGateway) GetStartFlowResult(ctx context.Context, currentCart *cartDomain.Cart, correlationID string, returnURL *url.URL) (*domain.FlowResult, error) {
	err := o.checkCart(currentCart)
	if err != nil {
		return nil, err
	}
	return &domain.FlowResult{}, nil
}

//StartFlow for offline payment requires not much - just redirect to the returnUrl :-)
func (o *OfflineWebCartPaymentGateway) StartFlow(ctx context.Context, currentCart *cartDomain.Cart, correlationID string, returnURL *url.URL) (web.Result, error) {
	err := o.checkCart(currentCart)
	if err != nil {
		return nil, err
	}
	return o.responder.URLRedirect(returnURL), nil
}

// GetFlowResult for offline payment can always return a simple valid payment that matches the given cart
func (o *OfflineWebCartPaymentGateway) GetFlowResult(ctx context.Context, currentCart *cartDomain.Cart, correlationID string) (*cartDomain.Payment, error) {
	err := o.checkCart(currentCart)
	if err != nil {
		return nil, err
	}
	cartPayment := cartDomain.Payment{
		Gateway: OfflineWebCartPaymentGatewayCode,
	}

	for _, split := range currentCart.PaymentSelection.ChargeSplits {
		cartPayment.Transactions = append(cartPayment.Transactions, cartDomain.Transaction{
			Method:            split.Method,
			Status:            cartDomain.PaymentStatusOpen,
			ValuedAmountPayed: split.Amount,
			AmountPayed:       split.Amount,
		})
	}

	return &cartPayment, nil
}

//ConfirmResult - nothing to confirm  for offline payment
func (o *OfflineWebCartPaymentGateway) ConfirmResult(ctx context.Context, cart *cartDomain.Cart, cartPayment *cartDomain.Payment) error {
	return nil
}
