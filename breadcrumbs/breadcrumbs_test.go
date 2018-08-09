package breadcrumbs

import (
	"testing"

	"flamingo.me/flamingo/framework/web/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	ctx := new(mocks.Context)
	crumb := Crumb{
		Title: "Test",
		Url:   "http://testurl/",
	}

	ctx.On("Value", contextKey).Once().Return(nil)
	ctx.On("WithValue", contextKey, []Crumb{crumb}).Once().Return(ctx)

	Add(ctx, crumb)

	ctx.On("Value", contextKey).Once().Return([]Crumb{crumb})
	ctx.On("WithValue", contextKey, []Crumb{crumb, crumb}).Once().Return(ctx)

	Add(ctx, crumb)

	ctx.AssertExpectations(t)
}

func TestController_Data(t *testing.T) {
	ctx := new(mocks.Context)
	crumb := Crumb{
		Title: "Test",
		Url:   "http://testurl/",
	}
	c := new(Controller)

	ctx.On("Value", contextKey).Once().Return([]Crumb{crumb, crumb})
	assert.Equal(t, c.Data(ctx, nil), []Crumb{crumb, crumb})

	ctx.On("Value", contextKey).Once().Return(nil)
	assert.Equal(t, c.Data(ctx, nil), []Crumb{})
}
