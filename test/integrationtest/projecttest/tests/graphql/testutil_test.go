package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertRefreshPlaceOrder(t *testing.T, e *httpexpect.Expect, blocking bool) (*httpexpect.Response, string) {
	t.Helper()
	mutationName := "Commerce_Checkout_RefreshPlaceOrder"

	if blocking {
		mutationName = "Commerce_Checkout_RefreshPlaceOrderBlocking"
	}
	mutation := `mutation { ` + mutationName + ` { uuid, state { name, __typename
... on Commerce_Checkout_PlaceOrderState_State_Redirect { URL }, ... on Commerce_Checkout_PlaceOrderState_State_PostRedirect { URL Parameters { key value } }, ... on Commerce_Checkout_PlaceOrderState_State_ShowHTML { HTML }, ... on Commerce_Checkout_PlaceOrderState_State_ShowIframe { URL }, ... on Commerce_Checkout_PlaceOrderState_State_Failed { reason{ __typename reason } } } } }`
	request := helper.GraphQlRequest(t, e, mutation)
	response := request.Expect()
	t.Log(response.Body())
	refreshUUID := getValue(response, mutationName, "uuid").Raw()
	require.IsType(t, "string", refreshUUID)
	return response, refreshUUID.(string)
}

func assertStartPlaceOrderWithValidUUID(t *testing.T, e *httpexpect.Expect) (*httpexpect.Response, interface{}) {
	t.Helper()
	mutation := `mutation {Commerce_Checkout_StartPlaceOrder(returnUrl: "placeorder") { uuid }}`
	request := helper.GraphQlRequest(t, e, mutation)
	response := request.Expect()
	t.Log(response.Body())
	response.Status(http.StatusOK)
	uuid := getValue(response, "Commerce_Checkout_StartPlaceOrder", "uuid").Raw()
	require.IsType(t, "string", uuid)
	assert.Regexp(t, "(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", uuid)
	return response, uuid
}

func getValue(response *httpexpect.Response, queryName, key string) *httpexpect.Value {
	return response.JSON().Object().Value("data").Object().Value(queryName).Object().Value(key)
}
