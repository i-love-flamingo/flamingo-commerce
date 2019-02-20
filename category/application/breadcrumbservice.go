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
func (bs *BreadcrumbService) AddBreadcrumb(ctx context.Context, category domain.Category) {
	if !category.Active() {
		return
	}
	if category.Code() != "" {
		u, _ := bs.router.URL(URLWithName(category.Code(), web.URLTitle(category.Name())))
		breadcrumbs.Add(ctx, breadcrumbs.Crumb{
			Title: category.Name(),
			URL:   u.String(),
			Code:  category.Code(),
		})
	}

	for _, subcat := range category.Categories() {
		bs.AddBreadcrumb(ctx, subcat)
	}
}
