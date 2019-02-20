package breadcrumbs

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	crumb := Crumb{
		Title: "Test",
		URL:   "http://testurl/",
	}

	r := web.CreateRequest(nil, nil)
	ctx := web.ContextWithRequest(context.Background(), r)

	Add(ctx, crumb)
	Add(ctx, crumb)

	b, _ := r.Values.Load(requestKey)
	assert.Len(t, b, 2)
}
