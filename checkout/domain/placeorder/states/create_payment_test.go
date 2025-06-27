package states_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces/mocks"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
)

type rawTransactionData struct{}

func TestCreatePayment_Run(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		factory := provideProcessFactory(t)

		cart := provideCartWithPaymentSelection(t)
		p, _ := factory.New(&url.URL{}, cart)

		state := states.CreatePayment{}

		expectedPayment := &placeorder.Payment{Gateway: "test", RawTransactionData: &rawTransactionData{}}
		gateway := mocks.NewWebCartPaymentGateway(t)
		gateway.EXPECT().StartFlow(mock.Anything, mock.Anything, p.Context().UUID, p.Context().ReturnURL).Return(&domain.FlowResult{EarlyPlaceOrder: true}, nil).Once()
		gateway.EXPECT().OrderPaymentFromFlow(mock.Anything, mock.Anything, p.Context().UUID).Return(expectedPayment, nil).Once()
		paymentService := paymentServiceHelper(t, gateway)

		state.Inject(paymentService)

		expectedResult := process.RunResult{
			RollbackData: states.CreatePaymentRollbackData{
				Gateway:            expectedPayment.Gateway,
				PaymentID:          expectedPayment.PaymentID,
				RawTransactionData: &rawTransactionData{},
			},
		}
		result := state.Run(context.Background(), p)
		assert.Equal(t, result, expectedResult)
		assert.Equal(t, p.Context().CurrentStateName, states.CompleteCart{}.Name())
		gateway.AssertExpectations(t)
	})

	t.Run("error missing payment selection", func(t *testing.T) {
		factory := provideProcessFactory(t)

		cart := cartDomain.Cart{}

		p, _ := factory.New(&url.URL{}, cart)

		state := states.CreatePayment{}

		paymentService := paymentServiceHelper(t, nil)

		state.Inject(paymentService)

		result := state.Run(context.Background(), p)
		assert.NotNil(t, result.Failed, "Missing PaymentSelection in cart should lead to an error")
	})

	t.Run("error during gateway.StartFlow", func(t *testing.T) {
		factory := provideProcessFactory(t)

		cart := provideCartWithPaymentSelection(t)

		p, _ := factory.New(&url.URL{}, cart)

		state := states.CreatePayment{}

		expectedError := errors.New("StartFlow payment error")

		gateway := mocks.NewWebCartPaymentGateway(t)
		gateway.EXPECT().StartFlow(mock.Anything, mock.Anything, p.Context().UUID, p.Context().ReturnURL).Return(nil, expectedError).Once()
		paymentService := paymentServiceHelper(t, gateway)
		state.Inject(paymentService)

		expectedResult := process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: expectedError.Error()},
		}
		assert.Equal(t, state.Run(context.Background(), p), expectedResult)
		gateway.AssertExpectations(t)
	})

	t.Run("error during gateway.OrderPaymentFromFlow", func(t *testing.T) {
		factory := provideProcessFactory(t)

		cart := provideCartWithPaymentSelection(t)

		p, _ := factory.New(&url.URL{}, cart)

		state := states.CreatePayment{}

		expectedError := errors.New("OrderPaymentFromFlow payment error")

		gateway := mocks.NewWebCartPaymentGateway(t)
		gateway.EXPECT().StartFlow(mock.Anything, mock.Anything, p.Context().UUID, p.Context().ReturnURL).Return(&domain.FlowResult{}, nil).Once()
		gateway.EXPECT().OrderPaymentFromFlow(mock.Anything, mock.Anything, p.Context().UUID).Return(nil, expectedError).Once()

		paymentService := paymentServiceHelper(t, gateway)
		state.Inject(paymentService)

		expectedResult := process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: expectedError.Error()},
		}
		assert.Equal(t, state.Run(context.Background(), p), expectedResult)
		gateway.AssertExpectations(t)
	})
}

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
			StartState: &states.New{},
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

func TestCreatePayment_IsFinal(t *testing.T) {
	state := states.CreatePayment{}
	assert.False(t, state.IsFinal())
}

func TestCreatePayment_Rollback(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		state := states.CreatePayment{}

		var data interface{}

		payment := &placeorder.Payment{
			Gateway:            "test",
			Transactions:       nil,
			RawTransactionData: &rawTransactionData{},
			PaymentID:          "1234",
		}

		data = states.CreatePaymentRollbackData{
			Gateway:            payment.Gateway,
			PaymentID:          payment.PaymentID,
			RawTransactionData: payment.RawTransactionData,
		}

		gateway := mocks.NewWebCartPaymentGateway(t)
		gateway.EXPECT().CancelOrderPayment(mock.Anything, payment).Return(nil).Once()
		paymentService := paymentServiceHelper(t, gateway)
		state.Inject(paymentService)

		result := state.Rollback(context.Background(), data)
		assert.Nil(t, result)
		gateway.AssertExpectations(t)
	})

	t.Run("RollbackData not of type", func(t *testing.T) {
		state := states.CreatePayment{}

		assert.Error(t, state.Rollback(context.Background(), "string"))
	})

	t.Run("Error during payment selection", func(t *testing.T) {
		state := states.CreatePayment{}

		var data interface{}

		payment := &placeorder.Payment{
			Gateway:            "non-existing",
			RawTransactionData: &rawTransactionData{},
		}

		data = states.CreatePaymentRollbackData{
			Gateway: payment.Gateway, PaymentID: payment.PaymentID, RawTransactionData: payment.RawTransactionData}

		paymentService := paymentServiceHelper(t, nil)
		state.Inject(paymentService)
		assert.Error(t, state.Rollback(context.Background(), data), "Missing payment selection / gateway should lead to an error")
	})

	t.Run("Error during CancelOrderPayment", func(t *testing.T) {
		state := states.CreatePayment{}

		var data interface{}

		payment := &placeorder.Payment{
			Gateway:            "test",
			Transactions:       nil,
			RawTransactionData: &rawTransactionData{},
			PaymentID:          "1234",
		}

		data = states.CreatePaymentRollbackData{
			Gateway:            payment.Gateway,
			PaymentID:          payment.PaymentID,
			RawTransactionData: payment.RawTransactionData,
		}

		expectedError := errors.New("generic payment error")
		gateway := mocks.NewWebCartPaymentGateway(t)
		gateway.EXPECT().CancelOrderPayment(mock.Anything, payment).Return(expectedError).Once()
		paymentService := paymentServiceHelper(t, gateway)
		state.Inject(paymentService)
		assert.EqualError(t, state.Rollback(context.Background(), data), expectedError.Error())
		gateway.AssertExpectations(t)
	})

}
