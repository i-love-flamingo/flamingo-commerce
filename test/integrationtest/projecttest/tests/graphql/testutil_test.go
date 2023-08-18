//go:build integration
// +build integration

package graphql_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"unicode"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func loadGraphQL(t *testing.T, name string, replacements map[string]string) string {
	t.Helper()
	content, err := os.ReadFile(path.Join("testdata", name+".graphql"))
	if err != nil {
		t.Fatal(err)
	}

	r := make([]string, 2*len(replacements))
	i := 0
	for key, val := range replacements {
		r[i] = fmt.Sprintf("###%s###", key)
		r[i+1] = val
		i = i + 2
	}
	replacer := strings.NewReplacer(r...)

	return replacer.Replace(string(content))
}

// loginTestCustomer performs a login of the configured fake user
func loginTestCustomer(t *testing.T, e *httpexpect.Expect) {
	t.Helper()
	// perform login
	e.POST("/en/core/auth/callback/fake").
		WithForm(map[string]interface{}{"username": "username", "password": "password"}).
		Expect()

	// check if login succeeded
	resp := e.GET("/en/core/auth/debug").Expect()
	resp.Status(http.StatusOK)
	body := resp.Body().Raw()
	if !strings.Contains(body, "fake: fake/username:") || strings.Contains(body, "fake: identity not saved in session") {
		t.Fatal("login failed")
	}
}

// prepareCart adds a simple product via graphQl
func prepareCart(t *testing.T, e *httpexpect.Expect) {
	t.Helper()
	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": "fake_simple", "DELIVERY_CODE": "delivery"})).Expect().Status(http.StatusOK)
}

// prepareCartWithPaymentSelection adds a simple product via graphQl
func prepareCartWithPaymentSelection(t *testing.T, e *httpexpect.Expect, paymentMethod string, marketPlaceCode *string) {
	t.Helper()

	code := "fake_simple"
	if marketPlaceCode != nil {
		code = *marketPlaceCode
	}

	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": code, "DELIVERY_CODE": "delivery"})).Expect().Status(http.StatusOK)
	helper.GraphQlRequest(t, e, loadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": paymentMethod})).Expect().Status(http.StatusOK)
}

func updatePaymentSelection(t *testing.T, e *httpexpect.Expect, paymentMethod string) {
	t.Helper()
	query := loadGraphQL(t, "update_payment_selection", map[string]string{"PAYMENT_METHOD": paymentMethod})

	response := helper.GraphQlRequest(t, e, query).Expect()
	response.Status(http.StatusOK)
}

func performRefreshPlaceOrder(t *testing.T, e *httpexpect.Expect, blocking bool) *httpexpect.Response {
	t.Helper()
	fileName := "refresh"
	if blocking {
		fileName = "refresh_blocking"
	}
	mutation := loadGraphQL(t, fileName, nil)
	request := helper.GraphQlRequest(t, e, mutation)

	return request.Expect()
}

func checkRefreshForExpectedState(t *testing.T, e *httpexpect.Expect, expectedUUID string, expectedState map[string]interface{}) error {
	t.Helper()
	response := performRefreshPlaceOrder(t, e, false)
	data := make(map[string]interface{})
	require.NoError(t, json.Unmarshal([]byte(response.Body().Raw()), &data))
	var theData interface{}
	var ok bool
	if theData, ok = data["data"]; !ok || theData == nil {
		return fmt.Errorf("no data in response: %s", response.Body().Raw())
	}
	data = theData.(map[string]interface{})

	if theData, ok = data["Commerce_Checkout_RefreshPlaceOrder"]; !ok {
		return fmt.Errorf("no data>Commerce_Checkout_RefreshPlaceOrder in response: %s", response.Body().Raw())
	}

	data = theData.(map[string]interface{})

	refreshUUID := data["uuid"]
	if expectedUUID != refreshUUID {
		return fmt.Errorf("uuid has changed: expected: %q, got from refresh: %q", expectedUUID, refreshUUID)
	}
	actualState := data["state"]
	if diff := cmp.Diff(actualState, expectedState); diff != "" {
		return fmt.Errorf("timeout reached, -actual state +expected state =%v", diff)
	}
	return nil
}

func assertRefreshPlaceOrder(t *testing.T, e *httpexpect.Expect, blocking bool) (*httpexpect.Response, string) {
	t.Helper()
	response := performRefreshPlaceOrder(t, e, blocking)
	mutationName := "Commerce_Checkout_RefreshPlaceOrder"
	if blocking {
		mutationName = "Commerce_Checkout_RefreshPlaceOrderBlocking"
	}
	refreshUUID := getValue(response, mutationName, "uuid").String().Raw()

	return response, refreshUUID
}

func assertStartPlaceOrderWithValidUUID(t *testing.T, e *httpexpect.Expect) (*httpexpect.Response, string) {
	t.Helper()
	mutation := loadGraphQL(t, "start", nil)
	request := helper.GraphQlRequest(t, e, mutation)
	response := request.Expect()
	response.Status(http.StatusOK)
	uuidMatches := getValue(response, "Commerce_Checkout_StartPlaceOrder", "uuid").String().
		Match("(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$")
	uuidMatches.Length().IsEqual(1)

	return response, uuidMatches.Index(0).Raw()
}

func getValue(response *httpexpect.Response, queryName, key string) *httpexpect.Value {
	return response.JSON().Object().Value("data").Object().Value(queryName).Object().Value(key)
}

func getArray(response *httpexpect.Response, queryName string) *httpexpect.Array {
	return response.JSON().Object().Value("data").Object().Value(queryName).Array()
}

func assertResponseForExpectedState(t *testing.T, response *httpexpect.Response, expectedState map[string]interface{}) {
	t.Helper()
	data := make(map[string]interface{})
	require.NoError(t, json.Unmarshal([]byte(response.Body().Raw()), &data))
	var theData interface{}
	var ok bool
	if theData, ok = data["data"]; !ok || theData == nil {
		t.Fatalf("no data in response: %s", response.Body().Raw())
		return
	}
	data = theData.(map[string]interface{})

	if diff := cmp.Diff(data, expectedState); diff != "" {
		t.Errorf("diff mismatch (-want +got):\\n%s", diff)
	}
}

// spaceMap strips all whitespace from given string
func spaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
