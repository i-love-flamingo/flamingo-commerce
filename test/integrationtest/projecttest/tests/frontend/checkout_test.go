//go:build integration
// +build integration

package frontend_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/payment"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/placeorder"

	"github.com/stretchr/testify/assert"
)

func Test_Checkout_SubmitCheckoutAction(t *testing.T) {
	t.Run("empty cart should lead to an error", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		response := e.GET(routeCheckoutSubmit).Expect()
		response.Status(http.StatusOK).Body().IsEqual("null\n")
	})

	t.Run("cart and valid form should lead to redirect to review page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")

		response := SubmitCheckoutForm(t, e, map[string]interface{}{
			"billingAddress": map[string]interface{}{
				"firstname": "firstname",
				"lastname":  "lastname",
				"email":     "test@test.com",
			},
			"deliveries": map[string]interface{}{
				"inflight": map[string]interface{}{
					"deliveryAddress": map[string]interface{}{
						"firstname": "firstname",
						"lastname":  "lastname",
						"email":     "test@test.com",
					},
				},
			},
			"payment": map[string]interface{}{
				"gateway": payment.FakePaymentGateway,
				"method":  domain.PaymentFlowStatusCompleted,
			},
		})

		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())
	})

	t.Run("checkout with invalid form should lead to page with form errors", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")

		response := SubmitCheckoutForm(t, e, nil)

		assert.Equal(t, routeCheckoutSubmit, response.Raw().Request.URL.RequestURI())

		form := response.JSON().Object().Value("Form").Object()
		form.Value("BillingAddressForm").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().IsFalse()
		form.Value("DeliveryForms").Object().Value("inflight").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().IsFalse()
		form.Value("SimplePaymentForm").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().IsFalse()
	})

	t.Run("checkout with cart requires payment", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")

		response := SubmitCheckoutForm(t, e, map[string]interface{}{
			"billingAddress": map[string]interface{}{
				"firstname": "firstname",
				"lastname":  "lastname",
				"email":     "test@test.com",
			},
			"deliveries": map[string]interface{}{
				"inflight": map[string]interface{}{
					"deliveryAddress": map[string]interface{}{
						"firstname": "firstname",
						"lastname":  "lastname",
						"email":     "test@test.com",
					},
				},
			},
		})

		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutSubmit, response.Raw().Request.URL.RequestURI())

		form := response.JSON().Object().Value("Form").Object()
		form.Value("SimplePaymentForm").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().IsFalse()
	})

	t.Run("checkout with zero cart possible without payment", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")
		CartApplyVoucher(t, e, "100-percent-off")

		response := SubmitCheckoutForm(t, e, map[string]interface{}{
			"billingAddress": map[string]interface{}{
				"firstname": "firstname",
				"lastname":  "lastname",
				"email":     "test@test.com",
			},
			"deliveries": map[string]interface{}{
				"inflight": map[string]interface{}{
					"deliveryAddress": map[string]interface{}{
						"firstname": "firstname",
						"lastname":  "lastname",
						"email":     "test@test.com",
					},
				},
			},
		})

		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())
	})
}

func Test_Checkout_ReviewActionAndPlaceOrderAction(t *testing.T) {
	t.Run("valid payment should lead to success page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")

		// submit checkout form
		response := SubmitCheckoutForm(t, e, map[string]interface{}{
			"billingAddress": map[string]interface{}{
				"firstname": "firstname",
				"lastname":  "lastname",
				"email":     "test@test.com",
			},
			"deliveries": map[string]interface{}{
				"inflight": map[string]interface{}{
					"deliveryAddress": map[string]interface{}{
						"firstname": "firstname",
						"lastname":  "lastname",
						"email":     "test@test.com",
					},
				},
			},
			"payment": map[string]interface{}{
				"gateway": payment.FakePaymentGateway,
				"method":  domain.PaymentFlowStatusCompleted,
			},
		})

		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())

		// submit review form
		response = SubmitReviewForm(t, e, map[string]interface{}{
			"proceed":            "1",
			"termsAndConditions": "1",
			"privacyPolicy":      "1",
		})

		assert.Equal(t, routeCheckoutSuccess, response.Raw().Request.URL.RequestURI())
		response.JSON().Object().Value("PaymentInfos").NotNull()
		response.JSON().Object().Value("PlacedOrderInfos").Array().Value(0).Object().Value("OrderNumber").String().NotEmpty()
	})

	t.Run("zero cart without payment should lead to success page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")
		CartApplyVoucher(t, e, "100-percent-off")

		// submit checkout form
		response := SubmitCheckoutForm(t, e, map[string]interface{}{
			"billingAddress": map[string]interface{}{
				"firstname": "firstname",
				"lastname":  "lastname",
				"email":     "test@test.com",
			},
			"deliveries": map[string]interface{}{
				"inflight": map[string]interface{}{
					"deliveryAddress": map[string]interface{}{
						"firstname": "firstname",
						"lastname":  "lastname",
						"email":     "test@test.com",
					},
				},
			},
		})

		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())

		// submit review form
		response = SubmitReviewForm(t, e, map[string]interface{}{
			"proceed":            "1",
			"termsAndConditions": "1",
			"privacyPolicy":      "1",
		})

		assert.Equal(t, routeCheckoutSuccess, response.Raw().Request.URL.RequestURI())
		response.JSON().Object().Value("PaymentInfos").IsNull()
	})

	t.Run("error during payment should lead to checkout page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")

		// submit checkout form
		response := SubmitCheckoutForm(t, e, map[string]interface{}{
			"billingAddress": map[string]interface{}{
				"firstname": "firstname",
				"lastname":  "lastname",
				"email":     "test@test.com",
			},
			"deliveries": map[string]interface{}{
				"inflight": map[string]interface{}{
					"deliveryAddress": map[string]interface{}{
						"firstname": "firstname",
						"lastname":  "lastname",
						"email":     "test@test.com",
					},
				},
			},
			"payment": map[string]interface{}{
				"gateway": payment.FakePaymentGateway,
				"method":  domain.PaymentFlowStatusFailed,
			},
		})

		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())

		// submit review form
		response = SubmitReviewForm(t, e, map[string]interface{}{
			"proceed":            "1",
			"termsAndConditions": "1",
			"privacyPolicy":      "1",
		})

		assert.Equal(t, routeCheckoutSubmit, response.Raw().Request.URL.RequestURI())
		response.JSON().Object().Value("ErrorInfos").Object().Value("HasPaymentError").Boolean().IsTrue()
		response.JSON().Object().Value("ErrorInfos").Object().Value("HasError").Boolean().IsTrue()
		response.JSON().Object().Value("ErrorInfos").Object().Value("ErrorMessage").String().IsEqual(domain.PaymentErrorCodeFailed)
	})

	t.Run("error during place order should lead to checkout page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		CartAddProduct(t, e, "fake_simple", 5, "", "inflight")

		// submit checkout form
		response := SubmitCheckoutForm(t, e, map[string]interface{}{
			"billingAddress": map[string]interface{}{
				"firstname": "firstname",
				"lastname":  "lastname",
				"email":     "test@test.com",
			},
			"personalData": map[string]interface{}{
				placeorder.CustomAttributesKeyPlaceOrderError: "generic error during place order",
			},
			"deliveries": map[string]interface{}{
				"inflight": map[string]interface{}{
					"deliveryAddress": map[string]interface{}{
						"firstname": "firstname",
						"lastname":  "lastname",
						"email":     "test@test.com",
					},
				},
			},
			"payment": map[string]interface{}{
				"gateway": payment.FakePaymentGateway,
				"method":  domain.PaymentFlowStatusCompleted,
			},
		})

		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())

		// submit review form
		response = SubmitReviewForm(t, e, map[string]interface{}{
			"proceed":            "1",
			"termsAndConditions": "1",
			"privacyPolicy":      "1",
		})

		assert.Equal(t, routeCheckoutSubmit, response.Raw().Request.URL.RequestURI())
		response.JSON().Object().Value("ErrorInfos").Object().Value("HasError").Boolean().IsTrue()
		response.JSON().Object().Value("ErrorInfos").Object().Value("HasPaymentError").Boolean().IsFalse()
		response.JSON().Object().Value("ErrorInfos").Object().Value("ErrorMessage").String().IsEqual("generic error during place order")
	})
}
