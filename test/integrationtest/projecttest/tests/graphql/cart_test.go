// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CartUpdateDeliveryAddresses(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	// check response of update delivery mutation
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "update_delivery_addresses", nil)).Expect()
	response.Status(http.StatusOK)
	forms := getValue(response, "Commerce_Cart_UpdateDeliveryAddresses", "forms").Array()
	forms.Length().Equal(3)

	address := forms.Element(0).Object()
	address.Value("deliveryCode").String().Equal("foo")
	address.Value("processed").Boolean().Equal(true)
	address.Value("useBillingAddress").Boolean().Equal(false)
	formData := address.Value("formData").Object()
	formData.Value("firstname").String().Equal("Foo-firstname")
	formData.Value("lastname").String().Equal("Foo-lastname")
	formData.Value("email").String().Equal("foo@flamingo.me")
	validation := address.Value("validationInfo").Object()
	validation.Value("generalErrors").Null()
	validation.Value("fieldErrors").Null()

	address = forms.Element(1).Object()
	address.Value("deliveryCode").Equal("bar")
	address.Value("processed").Boolean().Equal(true)
	address.Value("useBillingAddress").Boolean().Equal(true)
	formData = address.Value("formData").Object()
	formData.Value("firstname").String().Equal("")
	formData.Value("lastname").String().Equal("")
	formData.Value("email").String().Equal("")
	validation = address.Value("validationInfo").Object()
	validation.Value("generalErrors").Null()
	validation.Value("fieldErrors").Null()

	address = forms.Element(2).Object()
	address.Value("deliveryCode").Equal("invalid-email-address")
	address.Value("processed").Boolean().Equal(false)
	address.Value("useBillingAddress").Boolean().Equal(false)
	validation = address.Value("validationInfo").Object()
	validation.Value("generalErrors").Null()
	validation.Value("fieldErrors").NotNull()

	// check that deliveries are saved to cart
	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "cart", nil)).Expect()
	response.Status(http.StatusOK)
	getValue(response, "Commerce_Cart", "cart").Object().Value("deliveries").Array().Length().Equal(2)
}
