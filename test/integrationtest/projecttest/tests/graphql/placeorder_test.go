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
	"github.com/stretchr/testify/require"
)

// prepareCartWithPaymentSelection adds a simple product via graphQl
func prepareCartWithPaymentSelection(t *testing.T, e *httpexpect.Expect, paymentMethod string) {
	t.Helper()
	helper.GraphQlRequest(t, e, loadGraphQL(t, "add_to_cart", nil)).Expect().Status(http.StatusOK)
	helper.GraphQlRequest(t, e, loadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": paymentMethod})).Expect().Status(http.StatusOK)
}

func updatePaymentSelection(t *testing.T, e *httpexpect.Expect, paymentMethod string) {
	t.Helper()
	query := loadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": paymentMethod})

	response := helper.GraphQlRequest(t, e, query).Expect()
	response.Status(http.StatusOK)
}

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

// TODO: Test place order with fake order service, success / fail / fail during rollback / success during rollback
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
			require.IsType(t, "string", refreshUUID)
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
			require.IsType(t, "string", refreshUUID)
			assert.Equal(t, uuid, refreshUUID, "uuid has changed")
			actualState = getValue(response, "Commerce_Checkout_RefreshPlaceOrder", "state").Raw()

			if diff := cmp.Diff(actualState, expectedState); diff != "" {
				return fmt.Errorf("timeout reached, -actual state +expected state =%v", diff)
			}
			return nil
		})
	})

}

// TODO:
// - Without cart
// - with invalid cart
// - when place order process already running
func Test_StartPlaceOrder(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	t.Run("no payment selection", func(t *testing.T) {
		e := httpexpect.New(t, baseURL)
		assertStartPlaceOrderWithValidUUID(t, e)

		response, _ := assertRefreshPlaceOrder(t, e, true)

		actualState := getValue(response, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
		reason := actualState.Object().Value("reason").Object()
		assert.Equal(t, "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError", reason.Value("__typename").Raw())
		assert.Equal(t, "PaymentSelection not set", reason.Value("reason").Raw())
	})

	t.Run("already running process", func(t *testing.T) {
		t.Skip("skip for now") // @todo fix race condition w/ session
		return
		e := httpexpect.New(t, baseURL)
		prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe)

		_, firstUUID := assertStartPlaceOrderWithValidUUID(t, e)

		_, refreshUUID := assertRefreshPlaceOrder(t, e, true)
		assert.Equal(t, firstUUID, refreshUUID)

		_, secondUUID := assertStartPlaceOrderWithValidUUID(t, e)

		assert.Equal(t, firstUUID, secondUUID)

	})

}

// TODO:
// - without running process
// - with running process
func Test_RefreshPlaceOrder(t *testing.T) {

}

// TODO:
// - without running process
// - with running process
func Test_RefreshBlockingPlaceOrder(t *testing.T) {

}

// TODO:
// - without running process
// - with running process in final state
// - with running process in non final state
func Test_CancelPlaceOrder(t *testing.T) {

}

func Test_RestartStartPlaceOrder(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := httpexpect.New(t, baseURL)
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowStatusFailed)
	_, uuid1 := assertStartPlaceOrderWithValidUUID(t, e)
	// wait for fail
	res, _ := assertRefreshPlaceOrder(t, e, true)
	state := getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
	assert.Equal(t, "Failed", state.Object().Value("name").Raw())

	// restart
	_, uuid2 := assertStartPlaceOrderWithValidUUID(t, e)
	assert.NotEqual(t, uuid1, uuid2, "new process should have been started")
	res, _ = assertRefreshPlaceOrder(t, e, true)
	state = getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
	assert.Equal(t, "Failed", state.Object().Value("name").Raw())
	reason := state.Object().Value("reason").Object()
	// payment selection should still be set, so we get the payment error (not PaymentSelection not set)
	assert.Equal(t, "Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError", reason.Value("__typename").String().Raw())
	assert.Equal(t, "", reason.Value("reason").String().Raw())

	// update payment selection
	helper.GraphQlRequest(t, e, loadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": domain.PaymentFlowStatusCompleted})).Expect().Status(http.StatusOK)
	_, uuid3 := assertStartPlaceOrderWithValidUUID(t, e)
	assert.NotEqual(t, uuid2, uuid3, "new process should have been started")
	res, _ = assertRefreshPlaceOrder(t, e, true)
	state = getValue(res, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state")
	assert.Equal(t, "Success", state.Object().Value("name").Raw())
}

func Test_ActivePlaceOrder(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := httpexpect.New(t, baseURL)
	query := loadGraphQL(t, "active_place_order", nil)

	// no process should be running at the start
	request := helper.GraphQlRequest(t, e, query)
	response := request.Expect()
	status := response.JSON().Object().Value("data").Object().Value("Commerce_Checkout_ActivePlaceOrder")
	assert.False(t, status.Boolean().Raw())

	// let the process wait in iframe status
	prepareCartWithPaymentSelection(t, e, domain.PaymentFlowActionShowIframe)
	assertStartPlaceOrderWithValidUUID(t, e)
	// wait for goroutine to be finished
	assertRefreshPlaceOrder(t, e, true)

	// now we have a running process
	request = helper.GraphQlRequest(t, e, query)
	response = request.Expect()
	status = response.JSON().Object().Value("data").Object().Value("Commerce_Checkout_ActivePlaceOrder")
	assert.True(t, status.Boolean().Raw())

}
