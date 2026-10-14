//go:build integration
// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CommerceProductGetBadges(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "product_get_badges", nil)).Expect().Status(http.StatusOK)

	expectedState := map[string]interface{}{
		"simple": map[string]interface{}{
			"badges": map[string]interface{}{
				"all": []interface{}{
					map[string]interface{}{"code": "new", "label": "New", "color": "#E30613", "priority": 100.0},
				},
			},
		},
		"active_variant": map[string]interface{}{
			"badges": map[string]interface{}{
				"all": []interface{}{
					map[string]interface{}{"code": "new", "label": "New", "color": nil, "priority": nil},
				},
			},
		},
	}

	assertResponseForExpectedState(t, response, expectedState)
}
