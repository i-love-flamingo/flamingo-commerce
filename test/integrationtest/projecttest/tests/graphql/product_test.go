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

	for _, facet := range facets.Iter() {
		facetName := facet.Object().Value("name").String()

		switch facetName.Raw() {
		case "brandCode":
			facet.Object().Value("items").Array().First().Object().Value("value").String().Equal("apple")
		case "retailerCode":
			facet.Object().Value("items").Array().First().Object().Value("value").String().Equal("retailer")
		default:
			// Do nothing here
		}
	}
}
