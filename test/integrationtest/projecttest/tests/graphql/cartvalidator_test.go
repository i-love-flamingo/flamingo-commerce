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
