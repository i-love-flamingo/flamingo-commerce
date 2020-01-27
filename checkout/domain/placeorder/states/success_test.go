package states

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"github.com/stretchr/testify/assert"
)

func TestSuccess_IsFinal(t *testing.T) {
	s := Success{}
	assert.True(t, s.IsFinal())
}

func TestSuccess_Name(t *testing.T) {
	s := Success{}
	assert.Equal(t, "Success", s.Name())
}

func TestSuccess_Rollback(t *testing.T) {
	s := Success{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestSuccess_Run(t *testing.T) {
	s := Success{}
	assert.Equal(t, s.Run(context.Background(), &process.Process{}), process.RunResult{})
}
