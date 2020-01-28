package placeorder_test

import (
	"context"
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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
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
			StartState: &states.Wait{},
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
			name: "satus: unapproved, action: show iframe",
			flowStatus: flowStatusResult{
				flowStatus: &domain.FlowStatus{
					Status: domain.PaymentFlowStatusUnapproved,
					Action: domain.PaymentFlowActionShowIFrame,
					ActionData: domain.FlowActionData{
						URL: &url.URL{
							Scheme: "http",
							Host:   "test.com",
						},
					},
				},
				err: nil,
			},
			want: want{
				runResult: process.RunResult{
					RollbackData: nil,
					Failed:       nil,
				},
				state: states.ShowIframe{}.Name(),
				stateData: process.StateData(url.URL{
					Scheme: "http",
					Host:   "test.com",
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := provideProcessFactory(t)
			p, _ := factory.New(&url.URL{}, provideCartWithPaymentSelection(t))

			gateway := &mocks.WebCartPaymentGateway{}
			gateway.On("FlowStatus", mock.Anything, mock.Anything, p.Context().UUID).Return(tt.flowStatus.flowStatus, tt.flowStatus.err).Once()

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
