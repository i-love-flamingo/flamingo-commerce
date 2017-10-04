package controller

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder/mocks"

	"github.com/stretchr/testify/mock"
)

type (
	MockProductService struct{}
)

func (mps *MockProductService) Get(ctx context.Context, marketplacecode string) (domain.BasicProduct, error) {
	if marketplacecode == "fail" {
		return nil, errors.New("fail")
	}

	return domain.SimpleProduct{
		BasicProductData: domain.BasicProductData{Title: "My Product Title", MarketPlaceCode: marketplacecode},
	}, nil
}

func TestViewController_Get(t *testing.T) {
	expectedUrlTitle := "my-product-title"
	ctx := web.NewContext()

	redirectAware := new(mocks.RedirectAware)
	renderAware := new(mocks.RenderAware)
	errorAware := new(mocks.ErrorAware)

	vc := &View{
		ProductService: new(MockProductService),
		RedirectAware:  redirectAware,
		RenderAware:    renderAware,
		ErrorAware:     errorAware,
		Template:       "product/product",
		Router: &router.Router{
			RouterRegistry: router.NewRegistry(),
		},
	}
	u, _ := url.Parse(`http://test/`)
	vc.Router.SetBase(u)
	vc.Router.RouterRegistry.Route("/", `product.view(marketplacecode?="test", name?="test", variant?="test")`)

	redirectAware.On("Redirect", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return(nil)
	ctx.LoadParams(router.P{"marketplacecode": "test", "name": "testname"})
	vc.Get(ctx)
	redirectAware.AssertCalled(t, "Redirect", "product.view", map[string]string{"name": expectedUrlTitle, "marketplacecode": "test"})

	renderAware.On("Render", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("controller.productViewData")).Return(nil)
	ctx.LoadParams(router.P{"marketplacecode": "test", "name": expectedUrlTitle})
	vc.Get(ctx)
	renderAware.AssertCalled(t, "Render", ctx, vc.Template, mock.AnythingOfType("controller.productViewData"))

	errorAware.On("Error", ctx, mock.Anything).Return(nil)
	ctx.LoadParams(router.P{"marketplacecode": "fail", "name": "fail"})
	vc.Get(ctx)
	errorAware.AssertCalled(t, "Error", ctx, mock.Anything)
}
