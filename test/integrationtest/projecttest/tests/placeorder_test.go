// +build integration

package tests

import (
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

func Test_PlaceOrderGraphQL(t *testing.T) {
	e := httpexpect.New(t, "http://"+FlamingoUrl)
	t.Run("adding a simple product via graphQl", func(t *testing.T) {
		query := `
			mutation {
  Commerce_AddToCart(marketplaceCode: "fake_simple", qty: 1, deliveryCode: "delivery") {
    cart {
      id,
      itemCount
    }
  }
}`
		helper.GraphQlQueryRequest(t, e, query).
			Expect().
			Status(http.StatusOK).JSON().Object().Value("data").Object().Value("Commerce_AddToCart").Object().
			Value("cart").Object().
			Value("itemCount").Number().Equal(1)

		query = `
			query {
			  Commerce_Cart {
				cart {
				  id
				  itemCount
				  billingAddress {
					firstname
				  }
				  deliveries {
					deliveryInfo {
					  code
					}
					cartitems {
					  qty
					  productName
					}
				  }
				}
			  }
			}`
		helper.GraphQlQueryRequest(t, e, query).
			Expect().
			Status(http.StatusOK).JSON().Object().Value("data").Object().Value("Commerce_Cart").NotNull()
	})

}
