package dto_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"github.com/stretchr/testify/assert"
)

func TestAttributes_Get(t *testing.T) {
	t.Run("should return boolean value", func(t *testing.T) {
		var tests = []struct {
			in  string
			out *dto.CustomAttributeBoolean
		}{
			{"true", &dto.CustomAttributeBoolean{Value: true}},
			{"True", &dto.CustomAttributeBoolean{Value: true}},
			{"false", &dto.CustomAttributeBoolean{Value: false}},
			{"False", &dto.CustomAttributeBoolean{Value: false}},
		}

		for _, tt := range tests {
			t.Run(tt.in, func(t *testing.T) {
				key := "key"
				attr := dto.CustomAttributes{Attributes: map[string]string{key: tt.in}}
				value := attr.Get(key)
				customAttribute, success := value.(*dto.CustomAttributeBoolean)
				assert.True(t, success)
				if customAttribute.Key() != key {
					t.Errorf("got %q, want %q", customAttribute.Key(), key)
				}
				if customAttribute.Value != tt.out.Value {
					t.Errorf("got %v, want %v", customAttribute.Value, tt.out.Value)
				}
			})
		}
	})

	t.Run("should return string value", func(t *testing.T) {
		var tests = []struct {
			in  string
			out *dto.CustomAttributeString
		}{
			{"0", &dto.CustomAttributeString{Value: "0"}},
			{"1", &dto.CustomAttributeString{Value: "1"}},
			{"bla", &dto.CustomAttributeString{Value: "bla"}},
		}

		for _, tt := range tests {
			t.Run(tt.in, func(t *testing.T) {
				key := "key"
				attr := dto.CustomAttributes{Attributes: map[string]string{key: tt.in}}
				value := attr.Get(key)
				customAttribute, success := value.(*dto.CustomAttributeString)
				assert.True(t, success)
				if customAttribute.Key() != key {
					t.Errorf("got %q, want %q", customAttribute.Key(), key)
				}
				if customAttribute.Value != tt.out.Value {
					t.Errorf("got %v, want %v", customAttribute.Value, tt.out.Value)
				}
			})
		}
	})
}
