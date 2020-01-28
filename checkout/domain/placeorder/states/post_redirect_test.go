package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"github.com/stretchr/testify/assert"
)

func TestPostRedirect_IsFinal(t *testing.T) {
	s := states.PostRedirect{}
	assert.False(t, s.IsFinal())
}

func TestPostRedirect_Name(t *testing.T) {
	s := states.PostRedirect{}
	assert.Equal(t, "PostRedirect", s.Name())
}

func TestPostRedirect_Rollback(t *testing.T) {
	s := states.PostRedirect{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestPostRedirect_Run(t *testing.T) {

}
