package placeorder_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces/mocks"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func provideProcessFactory(t *testing.T) *process.Factory {
	t.Helper()
	factory := &process.Factory{}
	factory.Inject(
		func() *process.Process {
			return &process.Process{}
		},
		&struct {
			StartState  process.State `inject:"startState"`
			FailedState process.State `inject:"failedState"`
		}{
			StartState:  &states.New{},
			FailedState: &states.Failed{},
		},
	)
	return factory
}

func provideCartWithPaymentSelection(t *testing.T) cartDomain.Cart {
	t.Helper()
	cart := cartDomain.Cart{}
	paymentSelection, err := cartDomain.NewDefaultPaymentSelection("test", map[string]string{price.ChargeTypeMain: "main"}, cart)
	require.NoError(t, err)
	cart.PaymentSelection = paymentSelection
	return cart
}

func paymentServiceHelper(t *testing.T, gateway interfaces.WebCartPaymentGateway) *application.PaymentService {
	t.Helper()
	paymentService := &application.PaymentService{}

	paymentService.Inject(func() map[string]interfaces.WebCartPaymentGateway {
		return map[string]interfaces.WebCartPaymentGateway{
			"test": gateway,
		}
	})
	return paymentService
}

func TestPaymentValidator(t *testing.T) {
	type flowStatusResult struct {
		flowStatus *domain.FlowStatus
		err        error
	}

	type want struct {
		runResult process.RunResult
		state     string
		stateData process.StateData
	}

	tests := []struct {
		name       string
		flowStatus flowStatusResult
		want       want
	}{
		{
			name: "generic payment error during FlowStatus request",
			flowStatus: flowStatusResult{
				err: errors.New("generic_error"),
			},
			want: want{
				runResult: process.RunResult{Failed: process.ErrorOccurredReason{Error: "generic_error"}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: unapproved, action: show iframe",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionShowIframe,
					ActionData: domain.FlowActionData{
						URL: &url.URL{Scheme: "https", Host: "iframe-url.com"},
					},
				},
			},
			want: want{
				runResult: process.RunResult{Failed: nil},
				state:     states.ShowIframe{}.Name(),
				stateData: process.StateData(&url.URL{Scheme: "https", Host: "iframe-url.com"}),
			},
		},
		{
			name: "status: unapproved, action: show iframe - URL missing",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionShowIframe,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{Error: placeorder.ValidatePaymentErrorNoActionURL}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: unapproved, action: show html",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionShowHTML,
					ActionData: domain.FlowActionData{
						DisplayData: "<h2>Payment Form<h2><form><input type=\"hidden\" /></form>",
					},
				},
			},
			want: want{
				runResult: process.RunResult{Failed: nil},
				state:     states.ShowHTML{}.Name(),
				stateData: process.StateData("<h2>Payment Form<h2><form><input type=\"hidden\" /></form>"),
			},
		},
		{
			name: "status: unapproved, action: show html - html missing",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionShowHTML,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{Error: placeorder.ValidatePaymentErrorNoActionDisplayData}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: unapproved, action: redirect",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionRedirect,
					ActionData: domain.FlowActionData{
						URL: &url.URL{Scheme: "https", Host: "redirect-url.com"},
					},
				},
			},
			want: want{
				runResult: process.RunResult{Failed: nil},
				state:     states.Redirect{}.Name(),
				stateData: process.StateData(&url.URL{Scheme: "https", Host: "redirect-url.com"}),
			},
		},
		{
			name: "status: unapproved, action: redirect - URL missing",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionRedirect,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{Error: placeorder.ValidatePaymentErrorNoActionURL}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: unapproved, action: trigger_client_sdk",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionTriggerClientSDK,
					ActionData: domain.FlowActionData{
						URL:         &url.URL{Scheme: "https", Host: "redirect-url.com"},
						DisplayData: `{"foo": "bar"}`,
					},
				},
			},
			want: want{
				runResult: process.RunResult{Failed: nil},
				state:     states.TriggerClientSDK{}.Name(),
				stateData: process.StateData(states.TriggerClientSDKData{
					URL:  &url.URL{Scheme: "https", Host: "redirect-url.com"},
					Data: `{"foo": "bar"}`,
				}),
			},
		},
		{
			name: "status: unapproved, action: trigger_client_sdk - URL missing",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionTriggerClientSDK,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{Error: placeorder.ValidatePaymentErrorNoActionURL}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: unapproved, action: post redirect",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionPostRedirect,
					ActionData: domain.FlowActionData{
						URL: &url.URL{Scheme: "https", Host: "post-redirect-url.com"},
						FormParameter: map[string]domain.FormField{
							"form-field-0": {
								Value: []string{"value0", "value1"},
							},
							"form-field-1": {
								Value: []string{"value0"},
							},
						},
					},
				},
			},
			want: want{
				runResult: process.RunResult{Failed: nil},
				state:     states.PostRedirect{}.Name(),
				stateData: process.StateData(states.PostRedirectData{
					FormFields: map[string]states.FormField{
						"form-field-0": {
							Value: []string{"value0", "value1"},
						},
						"form-field-1": {
							Value: []string{"value0"},
						},
					},
					URL: &url.URL{Scheme: "https", Host: "post-redirect-url.com"},
				}),
			},
		},
		{
			name: "status: unapproved, action: show wallet payment",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionShowWalletPayment,
					ActionData: domain.FlowActionData{
						WalletDetails: &domain.WalletDetails{
							UsedPaymentMethod: "ApplePay",
							PaymentRequestAPI: domain.PaymentRequestAPI{
								Methods: `{
									"supportedMethods": "https://apple.com/apple-pay",
									"data": {
										"version": 3,
										"merchantIdentifier": "merchant.com.aoe.om3.kso",
										"merchantCapabilities": ["supports3DS"],
										"supportedNetworks": ["asterCard", "visa", "discover", "amex"],
										"countryCode": "DE"
									}
								}
								`,
								MerchantValidationURL: &url.URL{Scheme: "https", Host: "validate.example.com"},
							},
						},
					},
				},
			},
			want: want{
				runResult: process.RunResult{Failed: nil},
				state:     states.ShowWalletPayment{}.Name(),
				stateData: process.StateData(states.ShowWalletPaymentData{
					UsedPaymentMethod: "ApplePay",
					PaymentRequestAPI: domain.PaymentRequestAPI{
						Methods: `{
									"supportedMethods": "https://apple.com/apple-pay",
									"data": {
										"version": 3,
										"merchantIdentifier": "merchant.com.aoe.om3.kso",
										"merchantCapabilities": ["supports3DS"],
										"supportedNetworks": ["asterCard", "visa", "discover", "amex"],
										"countryCode": "DE"
									}
								}
								`,
						MerchantValidationURL: &url.URL{Scheme: "https", Host: "validate.example.com"},
					},
				}),
			},
		},
		{
			name: "status: unapproved, action: post redirect - URL missing",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionPostRedirect,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{Error: placeorder.ValidatePaymentErrorNoActionURL}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: unapproved, action: not supported",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: "unknown",
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{Error: "Payment action not supported: \"unknown\""}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: approved",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusApproved,
				},
			},
			want: want{
				runResult: process.RunResult{},
				state:     states.CompletePayment{}.Name(),
			},
		},
		{
			name: "status: completed",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusCompleted,
				},
			},
			want: want{
				runResult: process.RunResult{},
				state:     states.Success{}.Name(),
			},
		},
		{
			name: "status: completed",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusCompleted,
				},
			},
			want: want{
				runResult: process.RunResult{},
				state:     states.Success{}.Name(),
			},
		},
		{
			name: "status: aborted",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusAborted,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentCanceledByCustomerReason{}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: cancelled",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusCancelled,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: failed",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusFailed,
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{}},
				state:     states.New{}.Name(),
			},
		},
		{
			name: "status: wait for customer",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowWaitingForCustomer,
				},
			},
			want: want{
				runResult: process.RunResult{},
				state:     states.WaitForCustomer{}.Name(),
			},
		},
		{
			name: "status: unknown",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: "unknown",
				},
			},
			want: want{
				runResult: process.RunResult{Failed: process.PaymentErrorOccurredReason{Error: "Payment status not supported: \"unknown\""}},
				state:     states.New{}.Name(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := provideProcessFactory(t)
			p, _ := factory.New(&url.URL{}, provideCartWithPaymentSelection(t))
			gateway := mocks.NewWebCartPaymentGateway(t)
			gateway.EXPECT().FlowStatus(mock.Anything, mock.Anything, p.Context().UUID).Return(tt.flowStatus.flowStatus, tt.flowStatus.err).Once()

			paymentService := paymentServiceHelper(t, gateway)
			got := placeorder.PaymentValidator(context.Background(), p, paymentService)
			if diff := cmp.Diff(got, tt.want.runResult); diff != "" {
				t.Error("PaymentValidator() = -got +want", diff)
			}

			assert.Equal(t, p.Context().CurrentStateName, tt.want.state)

			if diff := cmp.Diff(p.Context().CurrentStateData, tt.want.stateData); diff != "" {
				t.Error("CurrentStateData = -got +want", diff)
			}
		})
	}
}
