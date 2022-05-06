package cart_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
	builder.AddCartItem("id", "method", domain.Charge{
		Type:  "type",
		Price: domain.NewFromInt(100, 1, "€"),
		Value: domain.NewFromInt(100, 1, "€"),
	})

	forGob := SomeTypeWithPaymentSelection{Selection: cart.NewPaymentSelection("gateway", builder.Build())}
	assert.Equal(t, domain.NewFromInt(100, 1, "€"), forGob.Selection.ItemSplit().Sum().TotalValue())

	err := enc.Encode(&forGob)
	if err != nil {
		t.Fatal("encode error:", err)
	}
	var received SomeTypeWithPaymentSelection
	err = dec.Decode(&received)
	if err != nil {
		t.Fatal("decode error 1:", err)
	}

	assert.Equal(t, "gateway", received.Selection.Gateway())
	assert.Equal(t, domain.NewFromInt(100, 1, "€"), received.Selection.ItemSplit().Sum().TotalValue())
}

func TestPaymentSplit_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		split   cart.PaymentSplit
		want    string
		wantErr bool
	}{
		{
			name: "marshall payment split",
			split: func() cart.PaymentSplit {
				result := cart.PaymentSplit{}
				charge := domain.Charge{
					Type: "t1",
				}
				firstQualifier := cart.SplitQualifier{
					Method:     "m1",
					ChargeType: charge.Type,
				}
				secondQualifier := cart.SplitQualifier{
					Method:          "m2",
					ChargeType:      charge.Type,
					ChargeReference: "r2",
				}
				result[firstQualifier] = charge
				result[secondQualifier] = charge
				return result
			}(),
			want:    `{"m1-t1-":{"Price":{"Amount":"0.00","Currency":""},"Value":{"Amount":"0.00","Currency":""},"Type":"t1","Reference":""},"m2-t1-r2":{"Price":{"Amount":"0.00","Currency":""},"Value":{"Amount":"0.00","Currency":""},"Type":"t1","Reference":""}}`,
			wantErr: false,
		},
		{
			name: "marshall payment split with empty values - error",
			split: func() cart.PaymentSplit {
				result := cart.PaymentSplit{}
				charge := domain.Charge{
					Type: "t1",
				}
				firstQualifier := cart.SplitQualifier{
					Method:     "m1",
					ChargeType: "",
				}
				secondQualifier := cart.SplitQualifier{
					Method:     "",
					ChargeType: charge.Type,
				}
				result[firstQualifier] = charge
				result[secondQualifier] = charge
				return result
			}(),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.split.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentSplit.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestPaymentSplit_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    cart.PaymentSplit
		wantErr bool
	}{
		{
			name: "unmarshall payment split",
			args: args{
				data: []byte("{\"m1-t1\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\"},\"m2-t1-r2\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\"}}"),
			},
			want: func() cart.PaymentSplit {
				result := cart.PaymentSplit{}
				charge := domain.Charge{
					Price:     domain.NewZero(""),
					Value:     domain.NewZero(""),
					Type:      "t1",
					Reference: "",
				}
				firstQualifier := cart.SplitQualifier{
					Method:     "m1",
					ChargeType: charge.Type,
				}
				secondQualifier := cart.SplitQualifier{
					Method:          "m2",
					ChargeType:      charge.Type,
					ChargeReference: "r2",
				}
				result[firstQualifier] = charge
				result[secondQualifier] = charge
				return result
			}(),
			wantErr: false,
		},
		{
			name: "unmarshall payment split empty method or type - error",
			args: args{
				data: []byte("{\"m1?t1\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\"},\"m2-t1\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\"}}"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			split := cart.PaymentSplit{}
			err := split.UnmarshalJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentSplit.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(split, tt.want) {
				t.Errorf("PaymentSplit.UnmarshalJSON() = %v, want %v", split, tt.want)
			}
		})
	}
}

func TestRemoveZeroCharges(t *testing.T) {
	chargeTypeToPaymentMethod := map[string]string{
		domain.ChargeTypeMain:     "cc",
		domain.ChargeTypeGiftCard: "giftcard",
		"loyalty":                 "loyalty",
	}

	selection := cart.DefaultPaymentSelection{
		GatewayProp:      "xy",
		ChargedItemsProp: cart.PaymentSplitByItem{},
	}

	builder := cart.PaymentSplitByItemBuilder{}

	builder.AddCartItem("item-1", "cc", domain.Charge{
		Price: domain.NewFromInt(25, 1, "$"),
		Value: domain.NewFromInt(25, 1, "$"),
		Type:  domain.ChargeTypeMain,
	})

	builder.AddCartItem("item-1", "loyalty", domain.Charge{
		Price: domain.NewFromInt(500, 1, "Points"),
		Value: domain.NewFromInt(5, 1, "$"),
		Type:  "loyalty",
	})

	builder.AddCartItem("item-1", "giftcard", domain.Charge{
		Price: domain.NewFromInt(0, 1, "$"),
		Value: domain.NewFromInt(0, 1, "$"),
		Type:  domain.ChargeTypeGiftCard,
	})

	builder.AddShippingItem("delivery-1", "loyalty", domain.Charge{
		Price: domain.NewFromInt(20, 1, "Points"),
		Value: domain.NewFromInt(5, 1, "$"),
		Type:  "loyalty",
	})

	builder.AddShippingItem("delivery-1", "cc", domain.Charge{
		Price: domain.NewFromInt(0, 1, "$"),
		Value: domain.NewFromInt(0, 1, "$"),
		Type:  domain.ChargeTypeMain,
	})

	selection.ChargedItemsProp = builder.Build()
	filteredSelection := cart.RemoveZeroCharges(selection, chargeTypeToPaymentMethod)
	_, found := filteredSelection.ItemSplit().CartItems["item-1"].ChargesByType().GetByType(domain.ChargeTypeGiftCard)

	if found == true {
		t.Errorf("item-1 shouldn't have charge of type %q", domain.ChargeTypeGiftCard)
	}

	_, found = filteredSelection.ItemSplit().ShippingItems["delivery-1"].ChargesByType().GetByType(domain.ChargeTypeMain)

	if found == true {
		t.Errorf("delivery-1 shouldn't have charge of type %q", domain.ChargeTypeMain)
	}

	assert.Regexp(t, "(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", filteredSelection.IdempotencyKey(), "IdempotencyKey looks not like a valid UUID v4")
	assert.NotEqual(t, uuid.Nil.String(), filteredSelection.IdempotencyKey())
}

func Test_NewDefaultPaymentSelection_IdempotencyKey(t *testing.T) {
	// NewDefaultPaymentSelection should generate a new idempotency key
	selection, _ := cart.NewDefaultPaymentSelection("", map[string]string{domain.ChargeTypeMain: "main"}, cart.Cart{})
	assert.Regexp(t, "(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", selection.IdempotencyKey(), "IdempotencyKey looks not like a valid UUID v4")
	assert.NotEqual(t, uuid.Nil.String(), selection.IdempotencyKey())

	// GenerateNewIdempotencyKey should return a payment selection with a different key
	newPaymentSelection, err := selection.GenerateNewIdempotencyKey()
	assert.NoError(t, err)
	assert.NotEqual(t, newPaymentSelection.IdempotencyKey(), selection.IdempotencyKey(), "IdempotencyKey should be not matching")
	assert.Equal(t, newPaymentSelection.CartSplit(), selection.CartSplit())
	assert.Equal(t, newPaymentSelection.Gateway(), selection.Gateway())
	assert.Equal(t, newPaymentSelection.TotalValue(), selection.TotalValue())
}

func Test_NewPaymentSelection_IdempotencyKey(t *testing.T) {
	// NewPaymentSelection should generate a new idempotency key
	selection := cart.NewPaymentSelection("", cart.PaymentSplitByItem{})
	assert.Regexp(t, "(?i)^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$", selection.IdempotencyKey(), "IdempotencyKey looks not like a valid UUID v4")
	assert.NotEqual(t, uuid.Nil.String(), selection.IdempotencyKey())
}

func TestDefaultPaymentSelection_MarshalJSON(t *testing.T) {
	selection, _ := cart.NewDefaultPaymentSelection("", map[string]string{domain.ChargeTypeMain: "main"}, cart.Cart{})

	expectedJSON := fmt.Sprintf("{\"GatewayProp\":\"\",\"ChargedItemsProp\":{\"CartItems\":{},\"ShippingItems\":{},\"TotalItems\":{}},\"IdempotencyKey\":\"%s\"}", selection.IdempotencyKey())

	actual, _ := json.Marshal(selection)
	actualJSON := string(actual)
	assert.Equal(t, expectedJSON, actualJSON)
}
