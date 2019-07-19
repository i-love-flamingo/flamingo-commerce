package cart_test

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"gopkg.in/go-playground/assert.v1"
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
		want    []byte
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
					Method:     "m2",
					ChargeType: charge.Type,
				}
				result[firstQualifier] = charge
				result[secondQualifier] = charge
				return result
			}(),
			want:    []byte("{\"m1-t1\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\",\"Reference\":\"\"},\"m2-t1\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\",\"Reference\":\"\"}}"),
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
			want:    nil,
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaymentSplit.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
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
				data: []byte("{\"m1-t1\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\"},\"m2-t1\":{\"Price\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Value\":{\"Amount\":\"0\",\"Currency\":\"\"},\"Type\":\"t1\"}}"),
			},
			want: func() cart.PaymentSplit {
				result := cart.PaymentSplit{}
				charge := domain.Charge{
					Type: "t1",
				}
				firstQualifier := cart.SplitQualifier{
					Method:     "m1",
					ChargeType: charge.Type,
				}
				secondQualifier := cart.SplitQualifier{
					Method:     "m2",
					ChargeType: charge.Type,
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
