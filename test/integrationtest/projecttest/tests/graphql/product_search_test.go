//go:build integration
// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CommerceProductSearchFacets(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	resp := helper.GraphQlRequest(t, e, loadGraphQL(t, "product_search", nil)).Expect()
	resp.Status(http.StatusOK)

	facets := getValue(resp, "Commerce_Product_Search", "facets").Array()
	facets.Length().IsEqual(3)

	for _, facet := range facets.Iter() {
		facetName := facet.Object().Value("name").String()

		switch facetName.Raw() {
		case "brandCode":
			facet.Object().Value("items").Array().Value(0).Object().Value("value").String().IsEqual("apple")
		case "retailerCode":
			facet.Object().Value("items").Array().Value(0).Object().Value("value").String().IsEqual("retailer")
		case "categoryCodes":
			facet.Object().Value("items").Array().Value(0).Object().Value("value").String().IsEqual("electronics")
		default:
			// Do nothing here
		}
	}
}
