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

func Test_CommerceProductSearchFacets(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	resp := helper.GraphQlRequest(t, e, loadGraphQL(t, "product_search", nil)).Expect()
	resp.Status(http.StatusOK)

	facets := getValue(resp, "Commerce_Product_Search", "facets").Array()
	facets.Length().Equal(2)

	brandCodeFacet := facets.First().Object()
	brandCodeFacet.Value("name").String().Equal("brandCode")
	brandCodeFacet.Value("items").Array().First().Object().Value("value").String().Equal("apple")

	retailerCodeFacet := facets.Last().Object()
	retailerCodeFacet.Value("name").String().Equal("retailerCode")
	retailerCodeFacet.Value("items").Array().First().Object().Value("value").String().Equal("retailer")
}
