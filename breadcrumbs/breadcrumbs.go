package breadcrumbs

import (
	"context"

	"flamingo.me/flamingo/framework/web"
)

type (
	// Crumb defines a breadcrumb
	Crumb struct {
		Title string
		Url   string
	}

	// Controller defines the data controller
	Controller struct{}

	contextKeyTyp string
)

const contextKey contextKeyTyp = "breadcrumbs"

// Add a breadcrumb to the current context
func Add(ctx_ context.Context, b Crumb) {
	ctx := web.ToContext(ctx_)
	if breadcrumbs, ok := ctx.Value(contextKey).([]Crumb); ok {
		ctx.WithValue(contextKey, append(breadcrumbs, b))
	} else {
		ctx.WithValue(contextKey, []Crumb{b})
	}
}

// Data controller
func (bc *Controller) Data(ctx context.Context, _ *web.Request) interface{} {
	if breadcrumbs, ok := ctx.Value(contextKey).([]Crumb); ok {
		return breadcrumbs
	}
	return []Crumb{}
}
