package cart_test

import (
	"bytes"
	"encoding/gob"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestPrice_MarshalBinaryForGob(t *testing.T) {
	type (
		SomeTypeWithPaymentSelection struct {
			Selection cart.PaymentSelection
		}
	)
	gob.Register(SomeTypeWithPaymentSelection{})

	var network bytes.Buffer
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	builder := cart.PaymentSplitByItemBuilder{}
	builder.AddCartItem("id","method",domain.Charge{
		Type:"type",
		Price: domain.NewFromInt(100,1,"€"),
		Value: domain.NewFromInt(100,1,"€"),
	})

	forGob := SomeTypeWithPaymentSelection{Selection: cart.NewPaymentSelection("gateway",builder.Build())}
	assert.Equal(t, domain.NewFromInt(100,1,"€"),forGob.Selection.ItemSplit().Sum().TotalValue())


	err := enc.Encode(&forGob)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	var received SomeTypeWithPaymentSelection
	err = dec.Decode(&received)
	if err != nil {
		t.Fatal("decode error 1:", err)
	}

	assert.Equal(t, "gateway",received.Selection.Gateway())
	assert.Equal(t, domain.NewFromInt(100,1,"€"),received.Selection.ItemSplit().Sum().TotalValue())
}
