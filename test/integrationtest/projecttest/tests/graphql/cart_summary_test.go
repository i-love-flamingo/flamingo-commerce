//go:build integration
// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CartSummary(t *testing.T) {
	t.Parallel()
	baseURL := "http://" + FlamingoURL
	tests := []struct {
		name                 string
		gatewayMethod        string
		marketPlaceCode      string
		expectedState        map[string]interface{}
		expectedGraphQLState string
	}{
		{
			name:            "sumPaymentSelectionCartSplitValueAmountByMethods",
			gatewayMethod:   "creditcard",
			marketPlaceCode: "fake_simple_with_fixed_price",
			expectedState: map[string]interface{}{
				"Commerce_Cart_DecoratedCart": map[string]interface{}{
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
			prepareCartWithPaymentSelection(t, e, tt.gatewayMethod, &tt.marketPlaceCode)

			response := helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_summary", map[string]string{"METHOD": tt.gatewayMethod})).Expect().Status(http.StatusOK)
			response.Status(http.StatusOK)

			assertResponseForExpectedState(t, response, tt.expectedState)
		})
	}
}
