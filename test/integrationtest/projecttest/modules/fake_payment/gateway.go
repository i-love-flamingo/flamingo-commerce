package fake_payment

import (
	"context"
	"errors"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

const (
	FakePaymentGateway = "fake_payment_gateway"
)

type (
	// Gateway used for testing all payment states
	Gateway struct {
		CartIsCompleted map[string]bool
	}

	// Method ...
	Method struct {
		Title  string
		Status *domain.FlowStatus
	}
)

var (
	_       interfaces.WebCartPaymentGateway = &Gateway{}
	methods                                  = map[string]Method{
		domain.PaymentFlowStatusCompleted: {
			Title: "Payment completed",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusCompleted,
			},
		},
		domain.PaymentFlowStatusFailed: {
			Title: "Payment failed",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusFailed,
			},
		},
		domain.PaymentErrorAbortedByCustomer: {
			Title: "Payment aborted by customer",
			Status: &domain.FlowStatus{
				Status: domain.PaymentErrorAbortedByCustomer,
			},
		},
		domain.PaymentFlowStatusCancelled: {
			Title: "Payment canceled",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusCancelled,
			},
		},
		domain.PaymentFlowStatusApproved: {
			Title: "Payment approved",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusApproved,
			},
		},
		domain.PaymentFlowWaitingForCustomer: {
			Title: "Payment waiting for customer",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowWaitingForCustomer,
			},
		},
		domain.PaymentFlowActionShowIframe: {
			Title: "Payment unapproved, iframe",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusUnapproved,
				Action: domain.PaymentFlowActionShowIframe,
				ActionData: domain.FlowActionData{
					URL: &url.URL{Scheme: "https", Host: "url.com"},
				},
			},
		},
		domain.PaymentFlowActionRedirect: {
			Title: "Payment unapproved, redirect",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusUnapproved,
				Action: domain.PaymentFlowActionRedirect,
				ActionData: domain.FlowActionData{
					URL: &url.URL{Scheme: "https", Host: "url.com"},
				},
			},
		},
		domain.PaymentFlowActionPostRedirect: {
			Title: "Payment unapproved, post-redirect",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusUnapproved,
				Action: domain.PaymentFlowActionPostRedirect,
				ActionData: domain.FlowActionData{
					URL: &url.URL{Scheme: "https", Host: "url.com"},
				},
			},
		},
		domain.PaymentFlowActionShowHTML: {
			Title: "Payment unapproved, html",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusUnapproved,
				Action: domain.PaymentFlowActionShowHTML,
				ActionData: domain.FlowActionData{
					DisplayData: "<h2>test</h2>",
				},
			},
		},
		"unknown": {
			Title: "Payment unapproved, unknown",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusUnapproved,
				Action: "unknown",
			},
		},
	}
)

// Inject dependencies
func (g *Gateway) Inject() *Gateway {
	g.CartIsCompleted = make(map[string]bool)

	return g
}

// Methods returns all payment gateway methods
func (g *Gateway) Methods() []domain.Method {
	result := make([]domain.Method, 0, len(methods))

	for key, val := range methods {
		result = append(result, domain.Method{Code: key, Title: val.Title})
	}

	return result
}

func (g *Gateway) isSupportedPaymentMethod(method string) bool {
	for _, supportedMethod := range g.Methods() {
		if supportedMethod.Code == method {
			return true
		}
	}
	return false
}

func (g *Gateway) StartFlow(ctx context.Context, cart *cart.Cart, correlationID string, returnURL *url.URL) (*domain.FlowResult, error) {
	method := ""
	// just grab the first method we find and use it to decide between the different use cases
	for qualifier, _ := range cart.PaymentSelection.CartSplit() {
		method = qualifier.Method
		break
	}

	if !g.isSupportedPaymentMethod(method) {
		return nil, errors.New("specified method not supported by payment gateway: " + method)
	}

	return &domain.FlowResult{
		Status: domain.FlowStatus{
			Status: domain.PaymentFlowStatusUnapproved,
		},
	}, nil

}

func (g *Gateway) FlowStatus(ctx context.Context, cart *cart.Cart, correlationID string) (*domain.FlowStatus, error) {
	methodCode := ""
	// just grab the first method we find and use it to decide between the different use cases
	for qualifier := range cart.PaymentSelection.CartSplit() {
		methodCode = qualifier.Method
		break
	}

	if !g.isSupportedPaymentMethod(methodCode) {
		return nil, errors.New("specified method not supported by payment gateway: " + methodCode)
	}

	if g.CartIsCompleted[cart.ID] {
		return &domain.FlowStatus{
			Status: domain.PaymentFlowStatusCompleted,
		}, nil
	}

	return methods[methodCode].Status, nil
}

func (g *Gateway) ConfirmResult(ctx context.Context, cart *cart.Cart, cartPayment *placeorder.Payment) error {
	g.CartIsCompleted[cart.ID] = true
	return nil
}

func (g *Gateway) OrderPaymentFromFlow(ctx context.Context, cart *cart.Cart, correlationID string) (*placeorder.Payment, error) {
	return &placeorder.Payment{
		Gateway: FakePaymentGateway,
		Transactions: []placeorder.Transaction{
			{
				TransactionID:     correlationID,
				AdditionalData:    nil,
				AmountPayed:       cart.GrandTotal(),
				ValuedAmountPayed: cart.GrandTotal(),
			},
		},
		RawTransactionData: nil,
		PaymentID:          "",
	}, nil
}

func (g *Gateway) CancelOrderPayment(ctx context.Context, cartPayment *placeorder.Payment) error {
	return nil
}
