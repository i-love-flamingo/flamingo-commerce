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
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "product_get", nil)).Expect().Status(http.StatusOK)

	expectedState := map[string]interface{}{
		"simple": map[string]interface{}{
			"title":           "TypeSimple product",
			"marketPlaceCode": "fake_simple",
			"meta":            map[string]interface{}{"keywords": []interface{}{string("keywords")}},
			"price":           map[string]interface{}{"activeBase": map[string]interface{}{"amount": 1.0, "currency": "€"}},
		},
		"configurable": map[string]interface{}{
			"title":           "TypeConfigurable product",
			"marketPlaceCode": "fake_configurable",
			"meta":            map[string]interface{}{"keywords": []interface{}{string("keywords")}},
			"variationSelections": []interface{}{
				map[string]interface{}{"code": string("color"), "label": string("Color")},
				map[string]interface{}{"code": string("size"), "label": string("Size")},
			},
		},
		"active_variant": map[string]interface{}{
			"type":  "configurable_with_activevariant",
			"title": "Shirt Black L",
			"meta":  map[string]interface{}{"keywords": []interface{}{string("keywords")}},
			"variationSelections": []interface{}{
				map[string]interface{}{"code": string("color"), "label": string("Color")},
				map[string]interface{}{"code": string("size"), "label": string("Size")},
			},
			"activeVariationSelections": []interface{}{
				map[string]interface{}{"code": string("color"), "label": string("Color"), "value": string("Black")},
				map[string]interface{}{"code": string("size"), "label": string("Size"), "value": string("L")},
			},
		},
	}

	assertResponseForExpectedState(t, response, expectedState)
}
