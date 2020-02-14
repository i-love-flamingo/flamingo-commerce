// +build integration

package rest_test

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/testhelper"
)

func Test_AddToCart(t *testing.T) {
	t.Run("adding simple product", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)

		testhelper.CartAddProduct(e, "fake_simple", 5, "", "")
		item := testhelper.CartGetItems(e).MustContain(t, "fake_simple")

		assert.Equal(t, 5, item.Qty)
	})

	t.Run("adding configurable product", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)

		testhelper.CartAddProduct(e, "fake_configurable", 3, "shirt-red-s", "")
		item := testhelper.CartGetItems(e).MustContain(t, "fake_configurable")

		assert.Equal(t, 3, item.Qty)
	})
}
