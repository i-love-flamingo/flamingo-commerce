package graphql_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func loadGraphQL(t *testing.T, name string, replacements map[string]string) string {
	t.Helper()
	content, err := ioutil.ReadFile(path.Join("testdata", name+".graphql"))
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

func assertRefreshPlaceOrder(t *testing.T, e *httpexpect.Expect, blocking bool) (*httpexpect.Response, string) {
	t.Helper()
	mutationName := "Commerce_Checkout_RefreshPlaceOrder"
	fileName := "refresh"
	if blocking {
		mutationName = "Commerce_Checkout_RefreshPlaceOrderBlocking"
		fileName = "refresh_blocking"
	}
	mutation := loadGraphQL(t, fileName, nil)
	request := helper.GraphQlRequest(t, e, mutation)
	response := request.Expect()
	t.Log(response.Body())
	refreshUUID := getValue(response, mutationName, "uuid").Raw()
	require.IsType(t, "string", refreshUUID)
	return response, refreshUUID.(string)
}

func assertStartPlaceOrderWithValidUUID(t *testing.T, e *httpexpect.Expect) (*httpexpect.Response, string) {
	t.Helper()
	mutation := loadGraphQL(t, "start", nil)
	request := helper.GraphQlRequest(t, e, mutation)
	response := request.Expect()
	t.Log(response.Body())
	response.Status(http.StatusOK)
	uuid := getValue(response, "Commerce_Checkout_StartPlaceOrder", "uuid").Raw()
	require.IsType(t, "string", uuid)
	assert.Regexp(t, "(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", uuid)
	return response, uuid.(string)
}

func getValue(response *httpexpect.Response, queryName, key string) *httpexpect.Value {
	return response.JSON().Object().Value("data").Object().Value(queryName).Object().Value(key)
}
