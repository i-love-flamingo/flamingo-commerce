// +build integration

package graphql_test

import (
	"net/http"
	"testing"
	"time"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CartSummary(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	tests := []struct {
		name                 string
		gatewayMethod        string
		expectedState        map[string]interface{}
		expectedGraphQLState string
	}{
		{
			name:          "sumPaymentSelectionCartSplitValueAmountByMethods",
			gatewayMethod: "creditcard",
			expectedState: map[string]interface{}{
				"Commerce_Cart": map[string]interface{}{
					"cartSummary": map[string]interface{}{
						"total": map[string]interface{}{
							"amount":   10.49,
							"currency": "€",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := integrationtest.NewHTTPExpect(t, baseURL)
			prepareCartWithPaymentSelection(t, e, tt.gatewayMethod)

			response := helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_summary", nil)).Expect().Status(http.StatusOK)
			response.Status(http.StatusOK)

			helper.AsyncCheckWithTimeout(t, time.Second, func() error {
				return assertResponseForExpectedState(t, e, response, tt.expectedState)
			})
		})
	}
}
