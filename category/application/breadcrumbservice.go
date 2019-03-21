package application

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/breadcrumbs"
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// RouterRouter defines a interface for testing
	RouterRouter interface {
		URL(name string, params map[string]string) (*url.URL, error)
	}

	// BreadcrumbService struct
	BreadcrumbService struct {
		router RouterRouter
	}
)

// Inject required dependencies
func (bs *BreadcrumbService) Inject(router RouterRouter) {
	bs.router = router
}

// AddBreadcrumb - add a breadcrumb based on a root category
func (bs *BreadcrumbService) AddBreadcrumb(ctx context.Context, tree domain.Tree) {
	if !tree.Active() {
		return
	}
	if tree.Code() != "" {
		u, _ := bs.router.URL(URLWithName(tree.Code(), web.URLTitle(tree.Name())))
		breadcrumbs.Add(ctx, breadcrumbs.Crumb{
			Title: tree.Name(),
			URL:   u.String(),
			Code:  tree.Code(),
		})
	}

	for _, subcat := range tree.SubTrees() {
		bs.AddBreadcrumb(ctx, subcat)
	}
}
