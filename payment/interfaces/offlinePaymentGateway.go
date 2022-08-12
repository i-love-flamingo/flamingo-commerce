package interfaces

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

// OfflineWebCartPaymentGateway provides an offline payment integration
type OfflineWebCartPaymentGateway struct {
	enabled   bool
	responder *web.Responder
}

const (
	// OfflineWebCartPaymentGatewayCode - the gateway code
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
	return []domain.Method{
		{
			Title: "cash on delivery",
			Code:  "offlinepayment_cashondelivery",
		},
		{
			Title: "cash in advance",
			Code:  "offlinepayment_cashinadvance",
		},
	}
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
	// Read the infos in the cart and check precondition
	if currentCart.PaymentSelection.Gateway() != OfflineWebCartPaymentGatewayCode {
		return errors.New("cart is not supposed to be paid by this gateway")
	}
	for qualifier := range currentCart.PaymentSelection.CartSplit() {
		if !o.isSupportedMethod(qualifier.Method) {
			return errors.New("cart payment method not supported by gateway")
		}
	}
	return nil
}

// StartFlow for offline payment and directly mark it as completed, since there is no online payment process
func (o *OfflineWebCartPaymentGateway) StartFlow(ctx context.Context, currentCart *cartDomain.Cart, correlationID string, returnURL *url.URL) (*domain.FlowResult, error) {
	err := o.checkCart(currentCart)
	if err != nil {
		return nil, err
	}
	return &domain.FlowResult{
		Status: domain.FlowStatus{
			Status: domain.PaymentFlowStatusCompleted,
		},
	}, nil
}

// FlowStatus for offline payment is always completed
func (o *OfflineWebCartPaymentGateway) FlowStatus(ctx context.Context, cart *cartDomain.Cart, correlationID string) (*domain.FlowStatus, error) {
	return &domain.FlowStatus{
		Status: domain.PaymentFlowStatusCompleted,
	}, nil
}

// ConfirmResult nothing to confirm  for offline payment
func (o *OfflineWebCartPaymentGateway) ConfirmResult(ctx context.Context, cart *cartDomain.Cart, cartPayment *placeorder.Payment) error {
	return nil
}

// OrderPaymentFromFlow create the order payment from the current cart/flow
func (o *OfflineWebCartPaymentGateway) OrderPaymentFromFlow(ctx context.Context, currentCart *cartDomain.Cart, correlationID string) (*placeorder.Payment, error) {
	err := o.checkCart(currentCart)
	if err != nil {
		return nil, err
	}

	cartPayment := placeorder.Payment{
		Gateway: OfflineWebCartPaymentGatewayCode,
	}

	var i int
	for qualifier, charge := range currentCart.PaymentSelection.CartSplit() {
		cartPayment.Transactions = append(cartPayment.Transactions, placeorder.Transaction{
			Method:            qualifier.Method,
			Status:            placeorder.PaymentStatusOpen,
			ValuedAmountPayed: charge.Value,
			AmountPayed:       charge.Price,
			TransactionID:     fmt.Sprintf("%v-%v", qualifier.Method, (i + 1)),
		})
	}

	return &cartPayment, nil
}

// CancelOrderPayment nothing to cancel for offline payment
func (o *OfflineWebCartPaymentGateway) CancelOrderPayment(ctx context.Context, cartPayment *placeorder.Payment) error {
	return nil
}
