package payment

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
	// FakePaymentGateway gateway code
	FakePaymentGateway = "fake_payment_gateway"
)

type (
	// FakeGateway used for testing all payment states
	FakeGateway struct {
		CartIsCompleted map[string]bool
	}

	// Method ...
	Method struct {
		Title  string
		Status *domain.FlowStatus
	}
)

var (
	_       interfaces.WebCartPaymentGateway = &FakeGateway{}
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
		domain.PaymentFlowStatusAborted: {
			Title: "Payment aborted",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusAborted,
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
		domain.PaymentFlowActionShowWalletPayment: {
			Title: "Payment unapproved, wallet",
			Status: &domain.FlowStatus{
				Status: domain.PaymentFlowStatusUnapproved,
				Action: domain.PaymentFlowActionShowWalletPayment,
				ActionData: domain.FlowActionData{
					WalletDetails: &domain.WalletDetails{
						UsedPaymentMethod: "ApplePay",
						PaymentRequestAPI: domain.PaymentRequestAPI{
							Methods: `{"a": "b"}`,
						},
					},
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
func (g *FakeGateway) Inject() *FakeGateway {
	g.CartIsCompleted = make(map[string]bool)

	return g
}

// Methods returns all payment gateway methods
func (g *FakeGateway) Methods() []domain.Method {
	result := make([]domain.Method, 0, len(methods))

	for key, val := range methods {
		result = append(result, domain.Method{Code: key, Title: val.Title})
	}

	return result
}

func (g *FakeGateway) isSupportedPaymentMethod(method string) bool {
	for _, supportedMethod := range g.Methods() {
		if supportedMethod.Code == method {
			return true
		}
	}
	return false
}

// StartFlow starts a new Payment flow
func (g *FakeGateway) StartFlow(ctx context.Context, cart *cart.Cart, correlationID string, returnURL *url.URL) (*domain.FlowResult, error) {
	method := ""
	// just grab the first method we find and use it to decide between the different use cases
	for qualifier := range cart.PaymentSelection.CartSplit() {
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

// FlowStatus returns a payment with a state depending on the supplied payment method
func (g *FakeGateway) FlowStatus(ctx context.Context, cart *cart.Cart, correlationID string) (*domain.FlowStatus, error) {
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

// ConfirmResult mark payment as completed
func (g *FakeGateway) ConfirmResult(ctx context.Context, cart *cart.Cart, cartPayment *placeorder.Payment) error {
	g.CartIsCompleted[cart.ID] = true
	return nil
}

// OrderPaymentFromFlow return fake payment
func (g *FakeGateway) OrderPaymentFromFlow(ctx context.Context, cart *cart.Cart, correlationID string) (*placeorder.Payment, error) {
	return &placeorder.Payment{
		Gateway: FakePaymentGateway,
		Transactions: []placeorder.Transaction{
			{
				TransactionID:     correlationID,
				AdditionalData:    nil,
				AmountPayed:       cart.GrandTotal,
				ValuedAmountPayed: cart.GrandTotal,
			},
		},
		RawTransactionData: nil,
		PaymentID:          "",
	}, nil
}

// CancelOrderPayment does nothing
func (g *FakeGateway) CancelOrderPayment(ctx context.Context, cartPayment *placeorder.Payment) error {
	return nil
}
