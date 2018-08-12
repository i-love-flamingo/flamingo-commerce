package breadcrumbs

import (
	"context"
	"testing"

	"flamingo.me/flamingo/framework/web"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	crumb := Crumb{
		Title: "Test",
		Url:   "http://testurl/",
	}

	r := &web.Request{Values: make(map[interface{}]interface{})}
	ctx := web.Context_(context.Background(), r)

	Add(ctx, crumb)
	Add(ctx, crumb)

	assert.Len(t, r.Values[requestKey], 2)
}
