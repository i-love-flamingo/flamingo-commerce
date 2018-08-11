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

const requestKey contextKeyTyp = "breadcrumbs"

// Add a breadcrumb to the current context
func Add(ctx context.Context, b Crumb) {
	req, _ := web.FromContext(ctx)

	if breadcrumbs, ok := req.Values[requestKey].([]Crumb); ok {
		req.Values[requestKey] = append(breadcrumbs, b)
	} else {
		req.Values[requestKey] = []Crumb{b}
	}
}

// Data controller
func (bc *Controller) Data(ctx context.Context, _ *web.Request) interface{} {
	req, _ := web.FromContext(ctx)

	if breadcrumbs, ok := req.Values[requestKey].([]Crumb); ok {
		return breadcrumbs
	}
	return []Crumb{}
}
