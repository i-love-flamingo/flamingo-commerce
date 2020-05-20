// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

// In Commerce, the secondary port of `Validator` from cart/domain/validation/cartValidator.go
// is not implemented. So the query will always return an empty result.
// This test checks just if the query works
func Test_CartValidator(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	helper.GraphQlRequest(t, e, loadGraphQL(t, "add_to_cart", nil)).Expect().Status(http.StatusOK)

	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "validate_cart", nil)).Expect()
	response.Status(http.StatusOK)
	getValue(response, "Commerce_Cart_Validator", "hasCommonError").Boolean().False()
}

// This test checks if qty Restrictions work
func Test_CartRestrictor(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	t.Run("restricted product", func(t *testing.T) {
		response := helper.GraphQlRequest(t, e, loadGraphQL(t, "validate_restrictor", map[string]string{"MARKETPLACECODE": "fake_simple"})).Expect()
		response.Status(http.StatusOK)
		getValue(response, "Commerce_Cart_QtyRestriction", "isRestricted").Boolean().True()
	})

	t.Run("unrestricted product", func(t *testing.T) {
		response := helper.GraphQlRequest(t, e, loadGraphQL(t, "validate_restrictor", map[string]string{"MARKETPLACECODE": "fake_configurable"})).Expect()
		response.Status(http.StatusOK)
		getValue(response, "Commerce_Cart_QtyRestriction", "isRestricted").Boolean().False()
	})

	t.Run("404 product", func(t *testing.T) {
		response := helper.GraphQlRequest(t, e, loadGraphQL(t, "validate_restrictor", map[string]string{"MARKETPLACECODE": "some"})).Expect()
		response.Status(http.StatusOK)
		response.JSON().Object().Value("errors").NotNull()
	})

}
