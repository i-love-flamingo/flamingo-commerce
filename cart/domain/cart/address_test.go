package cart_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"github.com/stretchr/testify/assert"
)

func TestAddress_FullName(t *testing.T) {
	address := cart.Address{Firstname: "first", Lastname: "last"}
	assert.Equal(t, "first last", address.FullName())
}

func TestAddress_IsEmpty(t *testing.T) {
	address := cart.Address{}
	assert.True(t, address.IsEmpty())

	address = cart.Address{Firstname: "hey"}
	assert.False(t, address.IsEmpty())

	address = cart.Address{AdditionalAddressLines: []string{"foo"}}
	assert.False(t, address.IsEmpty())

	var nilAddress *cart.Address
	assert.True(t, nilAddress.IsEmpty())
}
