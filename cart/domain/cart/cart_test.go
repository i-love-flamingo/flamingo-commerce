package cart

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCartItem(t *testing.T) {
	var cart = new(Cart)

	cartItem := Item{MarketplaceCode: "code1", Qty: 5}
	cart.Cartitems = append(cart.Cartitems, cartItem)

	found, nr := cart.HasItem("code1", "")
	assert.True(t, found)
	assert.Equal(t, nr, 1)

	item, err := cart.GetByLineNr(1)
	assert.NoError(t, err)
	assert.Equal(t, &cartItem, item)
}
