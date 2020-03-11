// +build integration

package frontend_test

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/payment"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/testhelper"
	"gotest.tools/assert"
)

const routeCheckoutSubmit = "/en/checkout"
const routeCheckoutReview = "/en/checkout/review"
const routeCheckoutPlaceOrder = "/en/checkout/placeorder"
const routeCheckoutSuccess = "/en/checkout/success"

func Test_SubmitCheckoutAction(t *testing.T) {
	t.Run("empty cart should lead to an error", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		response := e.GET(routeCheckoutSubmit).Expect()
		response.Status(http.StatusOK).Body().Equal("null\n")
	})

	t.Run("cart and valid form should lead to redirect to review page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		testhelper.CartAddProduct(e, "fake_simple", 5, "", "inflight")

		response := e.POST(routeCheckoutSubmit).WithForm(map[string]interface{}{
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
		}).Expect()

		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())
	})

	t.Run("checkout with invalid form should lead to page with form errors", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		testhelper.CartAddProduct(e, "fake_simple", 5, "", "inflight")

		response := e.POST(routeCheckoutSubmit).WithForm(map[string]interface{}{
			"foo": "bar",
		}).Expect()

		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutSubmit, response.Raw().Request.URL.RequestURI())

		form := response.JSON().Object().Value("Form").Object()
		form.Value("BillingAddressForm").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().False()
		form.Value("DeliveryForms").Object().Value("inflight").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().False()
		form.Value("SimplePaymentForm").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().False()
	})

	t.Run("checkout with cart requires payment", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		testhelper.CartAddProduct(e, "fake_simple", 5, "", "inflight")

		response := e.POST(routeCheckoutSubmit).WithForm(map[string]interface{}{
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
		}).Expect()

		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutSubmit, response.Raw().Request.URL.RequestURI())

		form := response.JSON().Object().Value("Form").Object()
		form.Value("SimplePaymentForm").Object().Value("ValidationInfo").Object().Value("IsValid").Boolean().False()
	})

	t.Run("checkout with zero cart possible without payment", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		testhelper.CartAddProduct(e, "fake_simple", 5, "", "inflight")
		testhelper.CartApplyVoucher(e, "100-percent-off")

		response := e.POST(routeCheckoutSubmit).WithForm(map[string]interface{}{
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
		}).Expect()

		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())
	})
}

func Test_PlaceOrderAction(t *testing.T) {
	t.Run("valid payment should lead to success page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		testhelper.CartAddProduct(e, "fake_simple", 5, "", "inflight")

		// submit checkout form
		response := e.POST(routeCheckoutSubmit).WithForm(map[string]interface{}{
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
		}).Expect()

		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())

		// place order
		response = e.GET(routeCheckoutPlaceOrder).Expect()
		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutSuccess, response.Raw().Request.URL.RequestURI())
		response.JSON().Object().Value("PaymentInfos").NotNull()
	})

	t.Run("zero cart without payment should lead to success page", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
		// prepare cart
		testhelper.CartAddProduct(e, "fake_simple", 5, "", "inflight")
		testhelper.CartApplyVoucher(e, "100-percent-off")

		// submit checkout form
		response := e.POST(routeCheckoutSubmit).WithForm(map[string]interface{}{
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
		}).Expect()

		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutReview, response.Raw().Request.URL.RequestURI())

		// place order
		response = e.GET(routeCheckoutPlaceOrder).Expect()
		response.Status(http.StatusOK)
		assert.Equal(t, routeCheckoutSuccess, response.Raw().Request.URL.RequestURI())
		response.JSON().Object().Value("PaymentInfos").Null()
	})
}
