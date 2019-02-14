package breadcrumbs

import (
	"context"
	"sync"
	"testing"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	crumb := Crumb{
		Title: "Test",
		Url:   "http://testurl/",
	}

	r := &web.Request{Values: new(sync.Map)}
	ctx := web.Context_(context.Background(), r)

	Add(ctx, crumb)

	b, _ := r.Values.Load(requestKey)
	assert.Len(t, b, 2)
}
