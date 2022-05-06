//go:build integration
// +build integration

package frontend_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"github.com/stretchr/testify/assert"
)

func Test_Cart_AddToCart(t *testing.T) {
	t.Run("adding simple product", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)

		CartAddProduct(t, e, "fake_simple", 5, "", "")
		item := CartGetItems(t, e).MustContain(t, "fake_simple")

		assert.Equal(t, 5, item.Qty)
	})

	t.Run("adding configurable product", func(t *testing.T) {
		e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)

		CartAddProduct(t, e, "fake_configurable", 3, "shirt-red-s", "")
		item := CartGetItems(t, e).MustContain(t, "fake_configurable")

		assert.Equal(t, 3, item.Qty)
	})
}
