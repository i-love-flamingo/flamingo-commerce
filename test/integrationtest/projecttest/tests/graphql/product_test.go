// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CommerceProductGet(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	resp := helper.GraphQlRequest(t, e, loadGraphQL(t, "product_get", nil)).Expect()
	resp.Status(http.StatusOK)
	getValue(resp, "Commerce_Product", "baseData").Object().Value("title").String().Equal("TypeConfigurable product")
}
