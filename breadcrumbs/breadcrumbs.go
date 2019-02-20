package breadcrumbs

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Crumb defines a breadcrumb
	Crumb struct {
		Title string
		URL   string
		Code  string
	}

	// Controller defines the data controller
	Controller struct{}

	contextKeyTyp string
)

const requestKey contextKeyTyp = "breadcrumbs"

// Add a breadcrumb to the current context
func Add(ctx context.Context, b Crumb) {
	req := web.RequestFromContext(ctx)

	breadcrumbs, _ := req.Values.Load(requestKey)
	if breadcrumbs, ok := breadcrumbs.([]Crumb); ok {
		breadcrumbs = append(breadcrumbs, b)
		req.Values.Store(requestKey, breadcrumbs)
	} else {
		req.Values.Store(requestKey, []Crumb{b})
	}
}

// Data controller
func (bc *Controller) Data(ctx context.Context, _ *web.Request, _ web.RequestParams) interface{} {
	req := web.RequestFromContext(ctx)

	breadcrumbs, _ := req.Values.Load(requestKey)
	if breadcrumbs, ok := breadcrumbs.([]Crumb); ok {
		return breadcrumbs
	}
	return []Crumb{}
}
