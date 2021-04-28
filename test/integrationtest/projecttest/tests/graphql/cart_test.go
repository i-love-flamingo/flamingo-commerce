// +build integration

package graphql_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"github.com/gavv/httpexpect/v2"
)

func Test_CartUpdateDeliveryAddresses(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	// check response of update delivery mutation
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "update_delivery_addresses", nil)).Expect()
	response.Status(http.StatusOK)
	forms := getArray(response, "Commerce_Cart_UpdateDeliveryAddresses")
	forms.Length().Equal(3)

	address := forms.Element(0).Object()
	address.Value("deliveryCode").String().Equal("foo")
	address.Value("carrier").String().Equal("carrier")
	address.Value("method").String().Equal("method")
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
	address.Value("carrier").String().Equal("")
	address.Value("method").String().Equal("")
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
	address.Value("carrier").String().Equal("")
	address.Value("method").String().Equal("")
	address.Value("processed").Boolean().Equal(false)
	address.Value("useBillingAddress").Boolean().Equal(false)
	validation = address.Value("validationInfo").Object()
	validation.Value("generalErrors").Null()
	validation.Value("fieldErrors").NotNull()

	// check that deliveries are saved to cart
	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect()
	response.Status(http.StatusOK)
	getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("deliveries").Array().Length().Equal(2)
}

func Test_CommerceCartUpdateDeliveryShippingOptions(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	// add some deliveries
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "update_delivery_addresses", nil)).Expect()
	response.Status(http.StatusOK)
	forms := getArray(response, "Commerce_Cart_UpdateDeliveryAddresses")
	forms.Length().Equal(3)

	// update shipping options
	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "update_delivery_shipping_options", nil)).Expect()
	response.Status(http.StatusOK)
	forms = getArray(response, "Commerce_Cart_UpdateDeliveryShippingOptions")
	forms.Length().Equal(3)

	address := forms.Element(0).Object()
	address.Value("deliveryCode").String().Equal("foo")
	address.Value("carrier").String().Equal("foo-carrier")
	address.Value("method").String().Equal("foo-method")
	address.Value("processed").Boolean().Equal(true)
	validation := address.Value("validationInfo").Object()
	validation.Value("generalErrors").Null()
	validation.Value("fieldErrors").Null()

	address = forms.Element(1).Object()
	address.Value("deliveryCode").Equal("bar")
	address.Value("carrier").String().Equal("bar-carrier")
	address.Value("method").String().Equal("bar-method")
	address.Value("processed").Boolean().Equal(true)
	validation = address.Value("validationInfo").Object()
	validation.Value("generalErrors").Null()
	validation.Value("fieldErrors").Null()

	address = forms.Element(2).Object()
	address.Value("deliveryCode").Equal("non-existing")
	address.Value("carrier").String().Equal("bar")
	address.Value("method").String().Equal("foo")
	address.Value("processed").Boolean().Equal(false)
	validation = address.Value("validationInfo").Object()
	validation.Value("generalErrors").NotNull()
	validation.Value("fieldErrors").Null()

	// check cart
	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect()
	response.Status(http.StatusOK)
	deliveries := getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("deliveries").Array()
	deliveries.Length().Equal(2)
	deliveries.Element(0).Object().Value("deliveryInfo").Object().Value("carrier").String().Equal("foo-carrier")
	deliveries.Element(0).Object().Value("deliveryInfo").Object().Value("method").String().Equal("foo-method")
	deliveries.Element(1).Object().Value("deliveryInfo").Object().Value("carrier").String().Equal("bar-carrier")
	deliveries.Element(1).Object().Value("deliveryInfo").Object().Value("method").String().Equal("bar-method")
}

func Test_CommerceCartClean(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	prepareCart(t, e)

	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect()
	response.Status(http.StatusOK)
	getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("itemCount").Equal(1)

	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_clear", nil)).Expect().Status(http.StatusOK)

	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect().Status(http.StatusOK)
	getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("itemCount").Equal(0)
}

func Test_CommerceCartUpdateAdditionalData(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	request := helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_update_additional_data", nil))
	response := request.Expect().Body()

	expected := `{
				   "data": {
					 "Commerce_Cart_UpdateAdditionalData": {
					   "cart": {
						 "additionalData": {
						   "customAttributes": {
							 "foo": {
							   "key": "foo",
							   "value": "bar"
							 },
							 "biz": {
							   "key": "biz",
							   "value": "baz"
							 }
						   }
						 }
					   }
					 }
				   }
				 }`

	expected = spaceMap(expected)
	response.Equal(expected)
}

func Test_CommerceCartUpdateDeliveriesAdditionalData(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	prepareCartWithDeliveries(t, e)
	request := helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_update_deliveries_additional_data", nil))
	response := request.Expect().Body()

	expected := `{
				   "data": {
					 "Commerce_Cart_UpdateDeliveriesAdditionalData": {
					   "cart": {
						 "deliveries": [
						   {
							 "deliveryInfo": {
							   "additionalData": {
								  "foo": {
								   "key": "foo",
								   "value": "bar"
								 },
								 "biz": {
								   "key": "biz",
								   "value": "baz"
								 },
								 "one": null,
								 "three": null
							   }
							 }
						   },
						   {
							 "deliveryInfo": {
							   "additionalData": {
								 "foo": null,
								 "biz": null,
								 "one": {
								   "key": "one",
								   "value": "two"
								 },
								 "three": {
								   "key": "three",
								   "value": "four"
								 }
							   }
							 }
						   }
						 ]
					   }
					 }
				   }
				 }`

	expected = spaceMap(expected)
	response.Equal(expected)
}

// prepareCartWithDeliveries adds a simple product with different delivery codes via graphQl
func prepareCartWithDeliveries(t *testing.T, e *httpexpect.Expect) {
	t.Helper()
	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": "fake_simple", "DELIVERY_CODE": "delivery1"})).Expect().Status(http.StatusOK)
	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": "fake_simple", "DELIVERY_CODE": "delivery2"})).Expect().Status(http.StatusOK)
}
