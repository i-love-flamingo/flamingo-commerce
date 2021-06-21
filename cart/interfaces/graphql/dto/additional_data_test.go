package dto_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"github.com/stretchr/testify/assert"
)

func TestAttributes_Get(t *testing.T) {
	t.Run("should return entry of map as key value pair value or nil if not found", func(t *testing.T) {
		var tests = []struct {
			inMap map[string]string
			inKey string
			out   *dto.KeyValue
		}{
			{map[string]string{"key": "true"}, "key", &dto.KeyValue{Key: "key", Value: "true"}},
			{map[string]string{"foo": "bar"}, "foo", &dto.KeyValue{Key: "foo", Value: "bar"}},
			{map[string]string{"foo": "bar"}, "bar", nil},
			{map[string]string{"key": "true"}, "notfound", nil},
		}

		for _, tt := range tests {
			t.Run("", func(t *testing.T) {
				attr := dto.CustomAttributes{Attributes: tt.inMap}
				value := attr.Get(tt.inKey)
				assert.Equal(t, tt.out, value)
			})
		}
	})
}
