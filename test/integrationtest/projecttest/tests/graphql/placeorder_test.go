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
	"github.com/gavv/httpexpect/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// prepareCartWithPaymentSelection adds a simple product via graphQl
func prepareCartWithPaymentSelection(t *testing.T, e *httpexpect.Expect, paymentMethod string) {
	t.Helper()
	query := `mutation {
  Commerce_AddToCart(marketplaceCode: "fake_simple", qty: 1, deliveryCode: "delivery") {
    cart { id }
  }

  Commerce_Cart_UpdateSelectedPayment( gateway: "fake_payment_gateway", method: "` + paymentMethod + `" ) { processed }
}`

	helper.GraphQlRequest(t, e, query).Expect().Status(http.StatusOK)
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
				"Parameters": []interface{}{}, //todo: better data in fake gateway
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
			mutation := `mutation {Commerce_Checkout_StartPlaceOrder(returnUrl: "placeorder") { uuid }}`

			request := helper.GraphQlRequest(t, e, mutation)
			response := request.Expect()
			t.Log(response.Body())
			response.Status(http.StatusOK)
			uuid := getValue(response, "Commerce_Checkout_StartPlaceOrder", "uuid").Raw()
			require.IsType(t, "string", uuid)
			assert.Regexp(t, "(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", uuid)

			var actualState interface{}

			helper.AsyncCheckWithTimeout(t, time.Second, func() error {
				mutation = `mutation { Commerce_Checkout_RefreshPlaceOrderBlocking { uuid, state { name, __typename, ... on Commerce_Checkout_PlaceOrderState_State_Redirect { URL }, ... on Commerce_Checkout_PlaceOrderState_State_PostRedirect { URL Parameters { key value } }, ... on Commerce_Checkout_PlaceOrderState_State_ShowHTML { HTML }, ... on Commerce_Checkout_PlaceOrderState_State_ShowIframe { URL }, ... on Commerce_Checkout_PlaceOrderState_State_Failed { reason{ __typename reason } } } } }`
				request = helper.GraphQlRequest(t, e, mutation)
				response = request.Expect()
				t.Log(response.Body())
				refreshUUID := getValue(response, "Commerce_Checkout_RefreshPlaceOrderBlocking", "uuid").Raw()
				require.IsType(t, "string", refreshUUID)
				assert.Equal(t, uuid, refreshUUID, "uuid has changed")
				actualState = getValue(response, "Commerce_Checkout_RefreshPlaceOrderBlocking", "state").Raw()
				if diff := cmp.Diff(actualState, tt.expectedState); diff != "" {
					return fmt.Errorf("timeout reached, -actual state +expected state =%v", diff)
				}
				return nil
			})
		})
	}
}

func getValue(response *httpexpect.Response, queryName, key string) *httpexpect.Value {
	return response.JSON().Object().Value("data").Object().Value(queryName).Object().Value(key)
}

// TODO: Test place order with fake order service, success / fail / fail during rollback / success during rollback
func Test_PlaceOrderWithOrderService(t *testing.T) {

}

// TODO:
// - Without cart
// - with invalid cart
// - when place order process already running
func Test_StartPlaceOrder(t *testing.T) {

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

// TODO: should generally test that restore of cart works and user can directly start a new process
// start process which fails, then restart process which also fails
// start process which fails, then restart with e.g. new payment method and success
func Test_RestartStartPlaceOrder(t *testing.T) {

}

// TODO: check that running process can be detected with CommerceCheckoutActivePlaceOrder
// - with running
// - without running
func Test_ActivePlaceOrder(t *testing.T) {
}
