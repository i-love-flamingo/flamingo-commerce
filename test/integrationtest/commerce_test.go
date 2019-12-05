// +build integration

package integrationtest

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/assert.v1"
	"gopkg.in/headzoo/surf.v1"
	"testing"
)

func Test_ProductPage(t *testing.T) {
	t.Log("Booting Up Flamingo Commerce Test Project")
	bootup()
	bow := surf.NewBrowser()
	t.Log("Calling product PDP")
	err := bow.Open("http://localhost:3210/en/product/fake_configurable/typeconfigurable-product.html")
	require.NoError(t, err)
	assert.Equal(t, 200, bow.StatusCode())

}
