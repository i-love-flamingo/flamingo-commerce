//go:build integration
// +build integration

package graphql_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

func Test_CartUpdateDeliveryAddresses(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	// check response of update delivery mutation
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "update_delivery_addresses", nil)).Expect()
	response.Status(http.StatusOK)
	forms := getArray(response, "Commerce_Cart_UpdateDeliveryAddresses")
	forms.Length().IsEqual(3)

	address := forms.Value(0).Object()
	address.Value("deliveryCode").String().IsEqual("foo")
	address.Value("carrier").String().IsEqual("carrier")
	address.Value("method").String().IsEqual("method")
	address.Value("processed").Boolean().IsEqual(true)
	address.Value("useBillingAddress").Boolean().IsEqual(false)
	formData := address.Value("formData").Object()
	formData.Value("firstname").String().IsEqual("Foo-firstname")
	formData.Value("lastname").String().IsEqual("Foo-lastname")
	formData.Value("email").String().IsEqual("foo@flamingo.me")
	validation := address.Value("validationInfo").Object()
	validation.Value("generalErrors").IsNull()
	validation.Value("fieldErrors").IsNull()

	address = forms.Value(1).Object()
	address.Value("deliveryCode").IsEqual("bar")
	address.Value("carrier").String().IsEqual("")
	address.Value("method").String().IsEqual("")
	address.Value("processed").Boolean().IsEqual(true)
	address.Value("useBillingAddress").Boolean().IsEqual(true)
	formData = address.Value("formData").Object()
	formData.Value("firstname").String().IsEqual("")
	formData.Value("lastname").String().IsEqual("")
	formData.Value("email").String().IsEqual("")
	validation = address.Value("validationInfo").Object()
	validation.Value("generalErrors").IsNull()
	validation.Value("fieldErrors").IsNull()

	address = forms.Value(2).Object()
	address.Value("deliveryCode").IsEqual("invalid-email-address")
	address.Value("carrier").String().IsEqual("")
	address.Value("method").String().IsEqual("")
	address.Value("processed").Boolean().IsEqual(false)
	address.Value("useBillingAddress").Boolean().IsEqual(false)
	validation = address.Value("validationInfo").Object()
	validation.Value("generalErrors").IsNull()
	validation.Value("fieldErrors").NotNull()

	// check that deliveries are saved to cart
	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect()
	response.Status(http.StatusOK)
	getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("deliveries").Array().Length().IsEqual(2)
}

func Test_CommerceCartUpdateDeliveryShippingOptions(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	// add some deliveries
	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "update_delivery_addresses", nil)).Expect()
	response.Status(http.StatusOK)
	forms := getArray(response, "Commerce_Cart_UpdateDeliveryAddresses")
	forms.Length().IsEqual(3)

	// update shipping options
	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "update_delivery_shipping_options", nil)).Expect()
	response.Status(http.StatusOK)
	getValue(response, "Commerce_Cart_UpdateDeliveryShippingOptions", "processed").Boolean().IsEqual(true)

	// check cart
	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect()
	response.Status(http.StatusOK)
	deliveries := getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("deliveries").Array()
	deliveries.Length().IsEqual(2)
	deliveries.Value(0).Object().Value("deliveryInfo").Object().Value("carrier").String().IsEqual("foo-carrier")
	deliveries.Value(0).Object().Value("deliveryInfo").Object().Value("method").String().IsEqual("foo-method")
	deliveries.Value(1).Object().Value("deliveryInfo").Object().Value("carrier").String().IsEqual("bar-carrier")
	deliveries.Value(1).Object().Value("deliveryInfo").Object().Value("method").String().IsEqual("bar-method")
}

func Test_CommerceCartClean(t *testing.T) {
	baseURL := "http://" + FlamingoURL
	e := integrationtest.NewHTTPExpect(t, baseURL)

	prepareCart(t, e)

	response := helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect()
	response.Status(http.StatusOK)
	getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("itemCount").IsEqual(1)

	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_clear", nil)).Expect().Status(http.StatusOK)

	response = helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_decorated_cart", nil)).Expect().Status(http.StatusOK)
	getValue(response, "Commerce_Cart_DecoratedCart", "cart").Object().Value("itemCount").IsEqual(0)
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
	response.IsEqual(expected)
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
	response.IsEqual(expected)
}

// prepareCartWithDeliveries adds a simple product with different delivery codes via graphQl
func prepareCartWithDeliveries(t *testing.T, e *httpexpect.Expect) {
	t.Helper()
	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": "fake_simple", "DELIVERY_CODE": "delivery1"})).Expect().Status(http.StatusOK)
	helper.GraphQlRequest(t, e, loadGraphQL(t, "cart_add_to_cart", map[string]string{"MARKETPLACE_CODE": "fake_simple", "DELIVERY_CODE": "delivery2"})).Expect().Status(http.StatusOK)
}

func TestAddBundleProductToCart(t *testing.T) {
	t.Parallel()

	t.Run("add to cart bundle product", func(t *testing.T) {
		t.Parallel()
		e := httpexpect.Default(t, "http://"+FlamingoURL)
		response := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart", map[string]string{
			"MARKETPLACE_CODE":          "fake_bundle",
			"DELIVERY_CODE":             "delivery",
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option2",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-s",
		}))

		body := response.Expect().Body()

		expected := `{
  "data": {
    "Commerce_Cart_AddToCart": {
      "decoratedDeliveries": [
        {
          "decoratedItems": [
            {
              "product": {
                "marketPlaceCode": "fake_bundle",
                "choices": [
                  {
                    "identifier": "identifier1",
                    "active": {
                      "marketPlaceCode": "simple_option1"
                    },
                    "activeOption": {
                      "product": {
                        "marketPlaceCode": "simple_option1"
                      },
                      "qty": 1
                    }
                  },
                  {
                    "identifier": "identifier2",
                    "active": {
                      "marketPlaceCode": "configurable_option2",
                      "variantMarketPlaceCode": "shirt-red-s"
                    },
                    "activeOption": {
                      "product": {
                        "marketPlaceCode": "configurable_option2",
                        "variantMarketPlaceCode": "shirt-red-s"
                      },
                      "qty": 1
                    }
                  }
                ]
              }
            }
          ]
        }
      ]
    }
  }
}`

		expected = spaceMap(expected)
		body.IsEqual(expected)
	})

	t.Run("add to cart bundle product, selected variant do not exists", func(t *testing.T) {
		t.Parallel()
		e := httpexpect.Default(t, "http://"+FlamingoURL)
		response := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart", map[string]string{
			"MARKETPLACE_CODE":          "fake_bundle",
			"DELIVERY_CODE":             "delivery",
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option2",
			"VARIANT_MARKETPLACE_CODE2": "there is no option like this",
		}))

		data := response.Expect().Status(http.StatusOK).JSON().Object()

		errorMessage := data.Value("errors").Array().Value(0).Object().Value("message").String().Raw()

		if !strings.Contains(errorMessage, "No Variant with code there is no option like this found") {
			t.Error("error do not contain: No Variant with code there is no option like this found ")
		}
	})

	t.Run("add to cart bundle product, required choice is not selected", func(t *testing.T) {
		t.Parallel()
		e := httpexpect.Default(t, "http://"+FlamingoURL)
		response := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart", map[string]string{
			"MARKETPLACE_CODE":  "fake_bundle",
			"DELIVERY_CODE":     "delivery",
			"IDENTIFIER1":       "identifier1",
			"MARKETPLACE_CODE1": "simple_option1",
		}))

		data := response.Expect().Status(http.StatusOK).JSON().Object()

		errorMessage := data.Value("errors").Array().Value(0).Object().Value("message").String().Raw()

		if !strings.Contains(errorMessage, "required choices are not selected") {
			t.Error("error do not contain: required choices are not selected")
		}
	})

	t.Run("add to cart bundle product, not existing marketplace code selected", func(t *testing.T) {
		t.Parallel()
		e := httpexpect.Default(t, "http://"+FlamingoURL)
		response := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart", map[string]string{
			"MARKETPLACE_CODE":          "fake_bundle",
			"DELIVERY_CODE":             "delivery",
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option1xxxxwrong",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option2",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-s",
		}))

		data := response.Expect().Status(http.StatusOK).JSON().Object()

		errorMessage := data.Value("errors").Array().Value(0).Object().Value("message").String().Raw()

		if !strings.Contains(errorMessage, "selected marketplace code does not exist") {
			t.Error("want: selected marketplace code does not exist, but have: ", errorMessage)
		}
	})
}

func TestUpdateBundleProductQty(t *testing.T) {
	t.Parallel()

	t.Run("update should update quantity of bundle product", func(t *testing.T) {
		t.Parallel()

		e := httpexpect.Default(t, "http://"+FlamingoURL)
		addResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart_Update_Qty_Helper", map[string]string{
			"MARKETPLACE_CODE":          "fake_bundle",
			"DELIVERY_CODE":             "delivery",
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option2",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-s",
		}))

		itemID := addResponse.Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().
			Value("Commerce_Cart_AddToCart").Object().Value("decoratedDeliveries").Array().
			Value(0).Object().Value("decoratedItems").Array().Value(0).Object().
			Value("item").Object().Value("id").String().Raw()

		updateResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "update_item_quantity", map[string]string{
			"ITEM_ID":       itemID,
			"DELIVERY_CODE": "inflight",
			"QTY":           "4",
		}))

		updateResponse.Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().
			Value("Commerce_Cart_UpdateItemQty").Object().Value("decoratedDeliveries").Array().
			Value(0).Object().Value("decoratedItems").Array().Value(0).Object().
			Value("item").Object().Value("qty").IsEqual(4)
	})
}

func TestUpdateBundleConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("update should update the bundle config", func(t *testing.T) {
		t.Parallel()

		e := httpexpect.Default(t, "http://"+FlamingoURL)
		addResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart_Update_Qty_Helper", map[string]string{
			"MARKETPLACE_CODE":          "fake_bundle",
			"DELIVERY_CODE":             "delivery",
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option2",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-s",
		}))

		itemID := addResponse.Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().
			Value("Commerce_Cart_AddToCart").Object().Value("decoratedDeliveries").Array().
			Value(0).Object().Value("decoratedItems").Array().Value(0).Object().
			Value("item").Object().Value("id").String().Raw()

		updateResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "Commerce_Cart_UpdateItemBundleConfig", map[string]string{
			"ITEM_ID":                   itemID,
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option2",
			"VARIANT_MARKETPLACE_CODE1": "",
			"QTY1":                      "1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option1",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-m",
			"QTY2":                      "1",
		}))

		body := updateResponse.Expect().Body()

		expected := `{
  "data": {
    "Commerce_Cart_UpdateItemBundleConfig": {
      "decoratedDeliveries": [
        {
          "decoratedItems": [
            {
              "product": {
                "marketPlaceCode": "fake_bundle",
                "choices": [
                  {
                    "identifier": "identifier1",
                    "active": {
                      "marketPlaceCode": "simple_option2"
                    },
                    "activeOption": {
                      "product": {
                        "marketPlaceCode": "simple_option2"
                      },
                      "qty": 1
                    }
                  },
                  {
                    "identifier": "identifier2",
                    "active": {
                      "marketPlaceCode": "configurable_option1",
                      "variantMarketPlaceCode": "shirt-red-m"
                    },
                    "activeOption": {
                      "product": {
                        "marketPlaceCode": "configurable_option1",
                        "variantMarketPlaceCode": "shirt-red-m"
                      },
                      "qty": 1
                    }
                  }
                ]
              }
            }
          ]
        }
      ]
    }
  }
}`

		expected = spaceMap(expected)
		body.IsEqual(expected)
	})

	t.Run("update should lead to an error if invalid bundle config (qty)", func(t *testing.T) {
		t.Parallel()

		e := httpexpect.Default(t, "http://"+FlamingoURL)
		addResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart_Update_Qty_Helper", map[string]string{
			"MARKETPLACE_CODE":          "fake_bundle",
			"DELIVERY_CODE":             "delivery",
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option2",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-s",
		}))

		itemID := addResponse.Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().
			Value("Commerce_Cart_AddToCart").Object().Value("decoratedDeliveries").Array().
			Value(0).Object().Value("decoratedItems").Array().Value(0).Object().
			Value("item").Object().Value("id").String().Raw()

		updateResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "Commerce_Cart_UpdateItemBundleConfig", map[string]string{
			"ITEM_ID":                   itemID,
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option2",
			"VARIANT_MARKETPLACE_CODE1": "",
			"QTY1":                      "1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option1",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-m",
			"QTY2":                      "0", // invalid qty
		}))

		data := updateResponse.Expect().Status(http.StatusOK).JSON().Object()

		errorMessage := data.Value("errors").Array().Value(0).Object().Value("message").String().Raw()

		if !strings.Contains(errorMessage, "selected quantity is out of range") {
			t.Error("want: selected quantity is out of range, but have: ", errorMessage)
		}
	})

	t.Run("update should lead to an error if invalid bundle config (not all required choices set)", func(t *testing.T) {
		t.Parallel()

		e := httpexpect.Default(t, "http://"+FlamingoURL)
		addResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "commerce_cart_AddBundleToCart_Update_Qty_Helper", map[string]string{
			"MARKETPLACE_CODE":          "fake_bundle",
			"DELIVERY_CODE":             "delivery",
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option1",
			"IDENTIFIER2":               "identifier2",
			"MARKETPLACE_CODE2":         "configurable_option2",
			"VARIANT_MARKETPLACE_CODE2": "shirt-red-s",
		}))

		itemID := addResponse.Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().
			Value("Commerce_Cart_AddToCart").Object().Value("decoratedDeliveries").Array().
			Value(0).Object().Value("decoratedItems").Array().Value(0).Object().
			Value("item").Object().Value("id").String().Raw()

		updateResponse := helper.GraphQlRequest(t, e, loadGraphQL(t, "Commerce_Cart_UpdateItemBundleConfig_one", map[string]string{
			"ITEM_ID":                   itemID,
			"IDENTIFIER1":               "identifier1",
			"MARKETPLACE_CODE1":         "simple_option2",
			"VARIANT_MARKETPLACE_CODE1": "",
			"QTY1":                      "1",
		}))

		data := updateResponse.Expect().Status(http.StatusOK).JSON().Object()

		errorMessage := data.Value("errors").Array().Value(0).Object().Value("message").String().Raw()

		if !strings.Contains(errorMessage, "required choices are not selected") {
			t.Error("want: required choices are not selected, but have: ", errorMessage)
		}
	})
}
