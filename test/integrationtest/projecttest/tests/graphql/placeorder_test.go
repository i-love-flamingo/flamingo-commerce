// +build integration

package graphql_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/placeorder"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func Test_PlaceOrderWithPaymentService(t *testing.T) {
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
			e := httpexpect.New(t, baseURL)
			prepareCartWithPaymentSelection(t, e, tt.gatewayMethod)
			_, uuid := assertStartPlaceOrderWithValidUUID(t, e)

			var actualState interface{}
			helper.AsyncCheckWithTimeout(t, time.Second, func() error {
				response, refreshUUID := assertRefreshPlaceOrder(t, e, false)
				assert.Equal(t, uuid, refreshUUID, "uuid has changed")
				actualState = getValue(response, "Commerce_Checkout_RefreshPlaceOrder", "state").Raw()
				if diff := cmp.Diff(actualState, tt.expectedState); diff != "" {
					return fmt.Errorf("timeout reached, -actual state +expected state =%v", diff)
				}
				return nil
			})
		})
	}
}

func Test_PlaceOrderWithOrderService(t *testing.T) {
	t.Run("PlaceOrder fails due to payment error, rollback of place order fails, restart afterwards can succeed ", func(t *testing.T) {
		baseURL := "http://" + FlamingoURL

		e := httpexpect.New(t, baseURL)
		prepareCartWithPaymentSelection(t, e, domain.PaymentFlowStatusFailed)
		placeorder.NextCancelFails = true

		_, uuid := assertStartPlaceOrderWithValidUUID(t, e)

		var actualState interface{}
		expectedState := map[string]interface{}{
			"name":       states.Failed{}.Name(),
			"__typename": "Commerce_Checkout_PlaceOrderState_State_Failed",
			"reason": map[string]interface{}{
				"__typename": "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError",
				"reason":     "",
			},
		}
		helper.AsyncCheckWithTimeout(t, time.Second, func() error {
			response, refreshUUID := assertRefreshPlaceOrder(t, e, false)
			assert.Equal(t, uuid, refreshUUID, "uuid has changed")
			actualState = getValue(response, "Commerce_Checkout_RefreshPlaceOrder", "state").Raw()

			if diff := cmp.Diff(actualState, expectedState); diff != "" {
				return fmt.Errorf("timeout reached, -actual state +expected state =%v", diff)
			}
			return nil
		})

		updatePaymentSelection(t, e, domain.PaymentFlowStatusApproved)

		_, uuid = assertStartPlaceOrderWithValidUUID(t, e)

		expectedState = map[string]interface{}{
			"name":       states.Success{}.Name(),
			"__typename": "Commerce_Checkout_PlaceOrderState_State_Success",
		}
		helper.AsyncCheckWithTimeout(t, time.Second, func() error {
			response, refreshUUID := assertRefreshPlaceOrder(t, e, false)
			assert.Equal(t, uuid, refreshUUID, "uuid has changed")
			actualState = getValue(response, "Commerce_Checkout_RefreshPlaceOrder", "state").Raw()

			if diff := cmp.Diff(actualState, expectedState); diff != "" {
				return fmt.Errorf("timeout reached, -actual state +expected state =%v", diff)
			}
			return nil
		})
	})

}

func Test_StartPlaceOrder(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	t.Run("no payment selection", func(t *testing.T) {
		e := httpexpect.New(t, baseURL)
		assertStartPlaceOrderWithValidUUID(t, e)

		response, _ := assertRefreshPlaceOrder(t, e, true)

		actualState := getValue(response, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
		reason := actualState.Object().Value("reason").Object()
		reason.Value("__typename").Equal("Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError")
		reason.Value("reason").Equal("PaymentSelection not set")
	})

	t.Run("already running process", func(t *testing.T) {
		e := httpexpect.New(t, baseURL)
		prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe)

		_, firstUUID := assertStartPlaceOrderWithValidUUID(t, e)

		_, refreshUUID := assertRefreshPlaceOrder(t, e, true)
		assert.Equal(t, firstUUID, refreshUUID)

		_, secondUUID := assertStartPlaceOrderWithValidUUID(t, e)

		assert.Equal(t, firstUUID, secondUUID)
	})

}

func Test_CancelPlaceOrder(t *testing.T) {
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
				err := response.Value("errors").Array().First().Object()
				err.Value("message").Equal("process already in final state, cancel not possible")
				err.Value("path").Array().First().Equal("Commerce_Checkout_CancelPlaceOrder")
			},
		},
		{
			name:          "not final",
			gatewayMethod: domain.PaymentFlowActionShowIframe,
			prepareAndRun: true,
			validator: func(t *testing.T, response *httpexpect.Object) {
				response.Value("data").Object().Value("Commerce_Checkout_CancelPlaceOrder").Boolean().True()
			},
		},
		{
			name:          "no running process",
			gatewayMethod: domain.PaymentFlowStatusCompleted,
			prepareAndRun: false,
			validator: func(t *testing.T, response *httpexpect.Object) {
				err := response.Value("errors").Array().First().Object()
				err.Value("message").Equal("ErrNoPlaceOrderProcess")
				err.Value("path").Array().First().Equal("Commerce_Checkout_CancelPlaceOrder")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := httpexpect.New(t, baseURL)
			if tt.prepareAndRun {
				prepareCartWithPaymentSelection(t, e, tt.gatewayMethod)
				assertStartPlaceOrderWithValidUUID(t, e)
				assertRefreshPlaceOrder(t, e, true)
			}

			request := helper.GraphQlRequest(t, e, loadGraphQL(t, "cancel", nil))
			response := request.Expect().JSON().Object()
			tt.validator(t, response)

		})
	}
}

func Test_RestartStartPlaceOrder(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := httpexpect.New(t, baseURL)
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowStatusFailed)
	_, uuid1 := assertStartPlaceOrderWithValidUUID(t, e)
	// wait for fail
	res, _ := assertRefreshPlaceOrder(t, e, true)
	getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state").Object().Value("name").Equal("Failed")
	orderInfo1 := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "orderInfos")
	orderNumber1 := orderInfo1.Object().Value("placedOrderInfos").Array().First().Object().Value("orderNumber").String()

	// restart
	_, uuid2 := assertStartPlaceOrderWithValidUUID(t, e)
	assert.NotEqual(t, uuid1, uuid2, "new process should have been started")
	res, _ = assertRefreshPlaceOrder(t, e, true)
	state := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
	state.Object().Value("name").Equal("Failed")
	// rollback of place order should lead to a new order number
	orderInfo2 := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "orderInfos")
	orderNumber2 := orderInfo2.Object().Value("placedOrderInfos").Array().First().Object().Value("orderNumber").String()
	orderNumber2.NotEqual(orderNumber1.Raw())

	// payment selection should still be set, so we get the payment error (not PaymentSelection not set)
	reason := state.Object().Value("reason").Object()
	reason.Value("__typename").Equal("Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError")
	reason.Value("reason").Equal("")

	// update payment selection
	helper.GraphQlRequest(t, e, loadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": domain.PaymentFlowStatusCompleted})).Expect().Status(http.StatusOK)
	_, uuid3 := assertStartPlaceOrderWithValidUUID(t, e)
	assert.NotEqual(t, uuid2, uuid3, "new process should have been started")
	res, _ = assertRefreshPlaceOrder(t, e, true)
	getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state").Object().Value("name").Equal("Success")

	// rollback of place order should lead to a new order number
	orderInfo3 := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "orderInfos")
	orderNumber3 := orderInfo3.Object().Value("placedOrderInfos").Array().First().Object().Value("orderNumber").String()
	orderNumber3.NotEqual(orderNumber1.Raw())
	orderNumber3.NotEqual(orderNumber2.Raw())
}

func Test_ActivePlaceOrder(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := httpexpect.New(t, baseURL)
	query := loadGraphQL(t, "active_place_order", nil)

	// no process should be running at the start
	request := helper.GraphQlRequest(t, e, query)
	response := request.Expect()
	response.JSON().Object().Value("data").Object().Value("Commerce_Checkout_ActivePlaceOrder").Boolean().False()

	// let the process wait in iframe status
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe)
	assertStartPlaceOrderWithValidUUID(t, e)
	// wait for goroutine to be finished
	assertRefreshPlaceOrder(t, e, true)

	// now we have a running process
	request = helper.GraphQlRequest(t, e, query)
	response = request.Expect()
	response.JSON().Object().Value("data").Object().Value("Commerce_Checkout_ActivePlaceOrder").Boolean().True()
}

func Test_GetCurrentState(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := httpexpect.New(t, baseURL)
	// no current context before start
	request := helper.GraphQlRequest(t, e, loadGraphQL(t, "current_context", nil))
	request.Expect().JSON().Object().Value("errors").Array().First().Object().Value("message").String().Equal("ErrNoPlaceOrderProcess")

	// prepare start and wait
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe)
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

	assert.Equal(t, state, state2, "current state must be the same as from RefreshPlaceOrderBlocking")
}
