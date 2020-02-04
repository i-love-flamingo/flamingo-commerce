// +build integration

package graphql_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
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

func Test_PlaceOrderGraphQL(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	tests := []struct {
		name                 string
		gatewayMethod        string
		expectedState        string
		expectedGraphQLState string
	}{
		{
			name:                 "Payment Completed",
			gatewayMethod:        domain.PaymentFlowStatusCompleted,
			expectedState:        states.Success{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_Success",
		},
		{
			name:                 "Payment Cancelled",
			gatewayMethod:        domain.PaymentFlowStatusCancelled,
			expectedState:        states.Failed{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_Failed",
		},
		{
			name:                 "Payment Approved",
			gatewayMethod:        domain.PaymentFlowStatusApproved,
			expectedState:        states.Success{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_Success",
		},
		{
			name:                 "Payment Failed",
			gatewayMethod:        domain.PaymentFlowStatusFailed,
			expectedState:        states.Failed{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_Failed",
		},
		{
			name:                 "Payment Aborted",
			gatewayMethod:        domain.PaymentFlowStatusAborted,
			expectedState:        states.Failed{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_Failed",
		},
		{
			name:                 "Payment Waiting For Customer",
			gatewayMethod:        domain.PaymentFlowWaitingForCustomer,
			expectedState:        states.WaitForCustomer{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_WaitForCustomer",
		},
		{
			name:                 "Payment Unapproved, Iframe",
			gatewayMethod:        domain.PaymentFlowActionShowIframe,
			expectedState:        states.ShowIframe{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_ShowIframe",
		},
		{
			name:                 "Payment Unapproved, HTML",
			gatewayMethod:        domain.PaymentFlowActionShowHTML,
			expectedState:        states.ShowHTML{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_ShowHTML",
		},
		{
			name:                 "Payment Unapproved, Redirect",
			gatewayMethod:        domain.PaymentFlowActionRedirect,
			expectedState:        states.Redirect{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_Redirect",
		},
		{
			name:                 "Payment Unapproved, Post Redirect",
			gatewayMethod:        domain.PaymentFlowActionPostRedirect,
			expectedState:        states.PostRedirect{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_PostRedirect",
		},
		{
			name:                 "Payment Unapproved, Unknown",
			gatewayMethod:        "unknown",
			expectedState:        states.Failed{}.Name(),
			expectedGraphQLState: "Commerce_Checkout_PlaceOrderState_State_Failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := httpexpect.New(t, baseURL)
			prepareCartWithPaymentSelection(t, e, tt.gatewayMethod)
			mutation := `mutation {Commerce_Checkout_StartPlaceOrder(returnUrl: "placeorder") { uuid }}`

			request := helper.GraphQlRequest(t, e, mutation)
			response := request.Expect()
			response.Status(http.StatusOK)
			uuid := getValue(response, "Commerce_Checkout_StartPlaceOrder", "uuid").Raw()
			require.IsType(t, "string", uuid)
			assert.Regexp(t, "(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", uuid)

			var actualState, actualGraphlQLState interface{}

			helper.AsyncCheckWithTimeout(t, time.Second, func() error {
				mutation = `mutation { Commerce_Checkout_RefreshPlaceOrder { uuid, state { name, __typename, ... on Commerce_Checkout_PlaceOrderState_State_Failed { reason{ __typename reason } } } } }`
				request = helper.GraphQlRequest(t, e, mutation)
				response = request.Expect()
				fmt.Println(response.Body())
				refreshUUID := getValue(response, "Commerce_Checkout_RefreshPlaceOrder", "uuid").Raw()
				require.IsType(t, "string", refreshUUID)
				assert.Equal(t, uuid, refreshUUID, "uuid has changed")
				actualState = getValue(response, "Commerce_Checkout_RefreshPlaceOrder", "state").Object().Value("name").Raw()
				actualGraphlQLState = getValue(response, "Commerce_Checkout_RefreshPlaceOrder", "state").Object().Value("__typename").Raw()
				if actualState != tt.expectedState {
					return fmt.Errorf("timeout reached, actual state %q != result state %q", actualState, tt.expectedState)
				}

				if actualGraphlQLState != tt.expectedGraphQLState {
					return fmt.Errorf("timeout reached, actual GraphQL state %q != result GraphQL state %q", actualGraphlQLState, tt.expectedGraphQLState)
				}

				return nil
			})
		})
	}
}

func getValue(response *httpexpect.Response, queryName, key string) *httpexpect.Value {
	return response.JSON().Object().Value("data").Object().Value(queryName).Object().Value(key)
}
