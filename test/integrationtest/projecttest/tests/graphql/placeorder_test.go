//go:build integration
// +build integration

package graphql_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/cart"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/placeorder"
)

func Test_PlaceOrderWithPaymentService(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	tests := []struct {
		name                 string
		gatewayMethod        string
		expectedState        map[string]interface{}
		expectedGraphQLState string
	}{
		{
			name:          "Payment Completed",
			gatewayMethod: domain.PaymentFlowStatusCompleted,
			expectedState: map[string]interface{}{
				"name":       states.Success{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_Success",
			},
		},
		{
			name:          "Payment Cancelled",
			gatewayMethod: domain.PaymentFlowStatusCancelled,
			expectedState: map[string]interface{}{
				"name":       states.Failed{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_Failed",
				"reason": map[string]interface{}{
					"__typename": "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError",
					"reason":     "",
				},
			},
		},
		{
			name:          "Payment Approved",
			gatewayMethod: domain.PaymentFlowStatusApproved,
			expectedState: map[string]interface{}{
				"name":       states.Success{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_Success",
			},
		},
		{
			name:          "Payment Failed",
			gatewayMethod: domain.PaymentFlowStatusFailed,
			expectedState: map[string]interface{}{
				"name":       states.Failed{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_Failed",
				"reason": map[string]interface{}{
					"__typename": "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError",
					"reason":     "",
				},
			},
		},
		{
			name:          "Payment Aborted",
			gatewayMethod: domain.PaymentFlowStatusAborted,
			expectedState: map[string]interface{}{
				"name":       states.Failed{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_Failed",
				"reason": map[string]interface{}{
					"__typename": "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentCanceledByCustomer",
					"reason":     "Payment canceled by customer",
				},
			},
		},
		{
			name:          "Payment Waiting For Customer",
			gatewayMethod: domain.PaymentFlowWaitingForCustomer,
			expectedState: map[string]interface{}{
				"name":       states.WaitForCustomer{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_WaitForCustomer",
			},
		},
		{
			name:          "Payment Unapproved, Iframe",
			gatewayMethod: domain.PaymentFlowActionShowIframe,
			expectedState: map[string]interface{}{
				"name":       states.ShowIframe{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_ShowIframe",
				"URL":        "https://url.com",
			},
		},
		{
			name:          "Payment Unapproved, HTML",
			gatewayMethod: domain.PaymentFlowActionShowHTML,
			expectedState: map[string]interface{}{
				"name":       states.ShowHTML{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_ShowHTML",
				"HTML":       "<h2>test</h2>",
			},
		},
		{
			name:          "Payment Unapproved, Redirect",
			gatewayMethod: domain.PaymentFlowActionRedirect,
			expectedState: map[string]interface{}{
				"name":       states.Redirect{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_Redirect",
				"URL":        "https://url.com",
			},
		},
		{
			name:          "Payment Unapproved, Post Redirect",
			gatewayMethod: domain.PaymentFlowActionPostRedirect,
			expectedState: map[string]interface{}{
				"name":       states.PostRedirect{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_PostRedirect",
				"URL":        "https://url.com",
				"Parameters": []interface{}{}, // todo: better data in fake gateway
			},
		},
		{
			name:          "Payment Unapproved, Show Wallet Payment",
			gatewayMethod: domain.PaymentFlowActionShowWalletPayment,
			expectedState: map[string]interface{}{
				"name":          states.ShowWalletPayment{}.Name(),
				"__typename":    "Commerce_Checkout_PlaceOrderState_State_ShowWalletPayment",
				"paymentMethod": "ApplePay",
			},
		},
		{
			name:          "Payment Unapproved, Unknown",
			gatewayMethod: "unknown",
			expectedState: map[string]interface{}{
				"name":       states.Failed{}.Name(),
				"__typename": "Commerce_Checkout_PlaceOrderState_State_Failed",
				"reason": map[string]interface{}{
					"__typename": "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError",
					"reason":     `Payment action not supported: "unknown"`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := integrationtest.NewHTTPExpect(t, baseURL)
			prepareCartWithPaymentSelection(t, e, tt.gatewayMethod, nil)
			_, uuid := assertStartPlaceOrderWithValidUUID(t, e)

			helper.AsyncCheckWithTimeout(t, time.Second, func() error {
				return checkRefreshForExpectedState(t, e, uuid, tt.expectedState)
			})
		})
	}
}

func Test_PlaceOrderWithOrderService(t *testing.T) {
	t.Run("PlaceOrder fails due to payment error, rollback of place order fails, restart afterwards can succeed ", func(t *testing.T) {
		baseURL := "http://" + FlamingoURL

		e := integrationtest.NewHTTPExpect(t, baseURL)
		prepareCartWithPaymentSelection(t, e, domain.PaymentFlowStatusFailed, nil)
		placeorder.NextCancelFails = true

		_, firstUUID := assertStartPlaceOrderWithValidUUID(t, e)

		expectedState := map[string]interface{}{
			"name":       states.Failed{}.Name(),
			"__typename": "Commerce_Checkout_PlaceOrderState_State_Failed",
			"reason": map[string]interface{}{
				"__typename": "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError",
				"reason":     "",
			},
		}
		helper.AsyncCheckWithTimeout(t, time.Second, func() error {
			return checkRefreshForExpectedState(t, e, firstUUID, expectedState)
		})

		updatePaymentSelection(t, e, domain.PaymentFlowStatusApproved)

		var secondUUID string
		helper.AsyncCheckWithTimeout(t, time.Second, func() error {
			_, secondUUID = assertStartPlaceOrderWithValidUUID(t, e)
			if secondUUID == firstUUID {
				return errors.New("UUID didn't change during new start place order")
			}
			return nil
		})
		assert.NotEmpty(t, secondUUID, "start order should return uuid")

		expectedState = map[string]interface{}{
			"name":       states.Success{}.Name(),
			"__typename": "Commerce_Checkout_PlaceOrderState_State_Success",
		}
		helper.AsyncCheckWithTimeout(t, time.Second, func() error {
			return checkRefreshForExpectedState(t, e, secondUUID, expectedState)
		})
	})

}

func Test_StartPlaceOrder(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	t.Run("no payment selection", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, baseURL)
		prepareCart(t, e)
		assertStartPlaceOrderWithValidUUID(t, e)

		response, _ := assertRefreshPlaceOrder(t, e, true)

		actualState := getValue(response, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
		reason := actualState.Object().Value("reason").Object()
		reason.Value("__typename").IsEqual("Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError")
		reason.Value("reason").IsEqual("paymentSelection not set")
	})

	t.Run("payment selection invalid", func(t *testing.T) {
		e := integrationtest.NewHTTPExpectWithCookies(t, baseURL, map[string]string{cart.FakePaymentSelectionValidatorCookie: ""})

		prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe, nil)
		assertStartPlaceOrderWithValidUUID(t, e)

		response, _ := assertRefreshPlaceOrder(t, e, true)

		actualState := getValue(response, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
		fmt.Println(actualState.Raw())
		reason := actualState.Object().Value("reason").Object()
		reason.Value("__typename").IsEqual("Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError")
		reason.Value("reason").IsEqual("fake payment selection validator error")
	})

	t.Run("replace already running process", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, baseURL)
		prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe, nil)

		_, firstUUID := assertStartPlaceOrderWithValidUUID(t, e)

		_, refreshUUID := assertRefreshPlaceOrder(t, e, true)
		assert.Equal(t, firstUUID, refreshUUID)

		_, secondUUID := assertStartPlaceOrderWithValidUUID(t, e)

		assert.NotEqual(t, firstUUID, secondUUID, "already running process should be replaced by a new one")
	})

}

func Test_CancelPlaceOrder(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	tests := []struct {
		name          string
		gatewayMethod string
		prepareAndRun bool
		validator     func(*testing.T, *httpexpect.Object)
	}{
		{
			name:          "already final",
			gatewayMethod: domain.PaymentFlowStatusCompleted,
			prepareAndRun: true,
			validator: func(t *testing.T, response *httpexpect.Object) {
				err := response.Value("errors").Array().Value(0).Object()
				err.Value("message").IsEqual("process already in final state, cancel not possible")
				err.Value("path").Array().Value(0).IsEqual("Commerce_Checkout_CancelPlaceOrder")
			},
		},
		{
			name:          "not final",
			gatewayMethod: domain.PaymentFlowActionShowIframe,
			prepareAndRun: true,
			validator: func(t *testing.T, response *httpexpect.Object) {
				response.Value("data").Object().Value("Commerce_Checkout_CancelPlaceOrder").Boolean().IsTrue()
			},
		},
		{
			name:          "no running process",
			gatewayMethod: domain.PaymentFlowStatusCompleted,
			prepareAndRun: false,
			validator: func(t *testing.T, response *httpexpect.Object) {
				err := response.Value("errors").Array().Value(0).Object()
				err.Value("message").IsEqual("ErrNoPlaceOrderProcess")
				err.Value("path").Array().Value(0).IsEqual("Commerce_Checkout_CancelPlaceOrder")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := integrationtest.NewHTTPExpect(t, baseURL)
			if tt.prepareAndRun {
				prepareCartWithPaymentSelection(t, e, tt.gatewayMethod, nil)
				assertStartPlaceOrderWithValidUUID(t, e)
				assertRefreshPlaceOrder(t, e, true)
			}

			request := helper.GraphQlRequest(t, e, loadGraphQL(t, "cancel", nil))
			response := request.Expect().JSON().Object()
			tt.validator(t, response)

		})
	}
}

func Test_ClearPlaceOrder(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	tests := []struct {
		name          string
		gatewayMethod string
		prepareAndRun bool
	}{
		{
			name:          "final",
			gatewayMethod: domain.PaymentFlowStatusCompleted,
			prepareAndRun: true,
		},
		{
			name:          "not final",
			gatewayMethod: domain.PaymentFlowActionShowIframe,
			prepareAndRun: true,
		},
		{
			name:          "no process",
			gatewayMethod: domain.PaymentFlowStatusCompleted,
			prepareAndRun: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := integrationtest.NewHTTPExpect(t, baseURL)
			if tt.prepareAndRun {
				prepareCartWithPaymentSelection(t, e, tt.gatewayMethod, nil)
				assertStartPlaceOrderWithValidUUID(t, e)
				assertRefreshPlaceOrder(t, e, true)
			}

			request := helper.GraphQlRequest(t, e, loadGraphQL(t, "clear_place_order", nil))
			response := request.Expect().JSON().Object()
			response.Value("data").Object().Value("Commerce_Checkout_ClearPlaceOrder").Boolean().IsTrue()

			request = helper.GraphQlRequest(t, e, loadGraphQL(t, "refresh_blocking", nil))
			response = request.Expect().JSON().Object()
			err := response.Value("errors").Array().Value(0).Object()
			err.Value("message").IsEqual("ErrNoPlaceOrderProcess")
			err.Value("path").Array().Value(0).IsEqual("Commerce_Checkout_RefreshPlaceOrderBlocking")
		})
	}
}

func Test_RestartStartPlaceOrder(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowStatusFailed, nil)
	_, uuid1 := assertStartPlaceOrderWithValidUUID(t, e)
	// wait for fail
	res, _ := assertRefreshPlaceOrder(t, e, true)
	getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state").Object().Value("name").IsEqual("Failed")
	orderInfo1 := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "orderInfos")
	orderNumber1 := orderInfo1.Object().Value("placedOrderInfos").Array().Value(0).Object().Value("orderNumber").String()

	// restart
	_, uuid2 := assertStartPlaceOrderWithValidUUID(t, e)
	assert.NotEqual(t, uuid1, uuid2, "new process should have been started")
	res, _ = assertRefreshPlaceOrder(t, e, true)
	state := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
	state.Object().Value("name").IsEqual("Failed")
	// rollback of place order should lead to a new order number
	orderInfo2 := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "orderInfos")
	orderNumber2 := orderInfo2.Object().Value("placedOrderInfos").Array().Value(0).Object().Value("orderNumber").String()
	orderNumber2.NotEqual(orderNumber1.Raw())

	// payment selection should still be set, so we get the payment error (not PaymentSelection not set)
	reason := state.Object().Value("reason").Object()
	reason.Value("__typename").IsEqual("Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError")
	reason.Value("reason").IsEqual("")

	// update payment selection
	helper.GraphQlRequest(t, e, loadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": domain.PaymentFlowStatusCompleted})).Expect().Status(http.StatusOK)
	_, uuid3 := assertStartPlaceOrderWithValidUUID(t, e)
	assert.NotEqual(t, uuid2, uuid3, "new process should have been started")
	res, _ = assertRefreshPlaceOrder(t, e, true)
	getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state").Object().Value("name").IsEqual("Success")

	// rollback of place order should lead to a new order number
	orderInfo3 := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "orderInfos")
	orderNumber3 := orderInfo3.Object().Value("placedOrderInfos").Array().Value(0).Object().Value("orderNumber").String()
	orderNumber3.NotEqual(orderNumber1.Raw())
	orderNumber3.NotEqual(orderNumber2.Raw())
}

func Test_ActivePlaceOrder(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	query := loadGraphQL(t, "active_place_order", nil)

	// no process should be running at the start
	request := helper.GraphQlRequest(t, e, query)
	response := request.Expect()
	response.JSON().Object().Value("data").Object().Value("Commerce_Checkout_ActivePlaceOrder").Boolean().IsFalse()

	// let the process wait in iframe status
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe, nil)
	assertStartPlaceOrderWithValidUUID(t, e)
	// wait for goroutine to be finished
	assertRefreshPlaceOrder(t, e, true)

	// now we have a running process
	request = helper.GraphQlRequest(t, e, query)
	response = request.Expect()
	response.JSON().Object().Value("data").Object().Value("Commerce_Checkout_ActivePlaceOrder").Boolean().IsTrue()
}

func Test_GetCurrentState(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	// no current context before start
	request := helper.GraphQlRequest(t, e, loadGraphQL(t, "current_context", nil))
	request.Expect().JSON().Object().Value("errors").Array().Value(0).Object().Value("message").String().IsEqual("ErrNoPlaceOrderProcess")

	// prepare start and wait
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe, nil)
	_, uuid := assertStartPlaceOrderWithValidUUID(t, e)
	result, uuid2 := assertRefreshPlaceOrder(t, e, true)
	assert.Equal(t, uuid, uuid2)
	state := getValue(result, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")

	// now we can get the current state
	request = helper.GraphQlRequest(t, e, loadGraphQL(t, "current_context", nil))
	response := request.Expect()
	uuid3 := getValue(response, "Commerce_Checkout_CurrentContext", "uuid").Raw()
	assert.Equal(t, uuid, uuid3)
	state2 := getValue(response, "Commerce_Checkout_CurrentContext", "state")

	assert.Equal(t, state.Raw(), state2.Raw(), "current state must be the same as from RefreshPlaceOrderBlocking")
}
