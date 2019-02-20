package controller

import (
	"context"
	"errors"
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	MockProductService struct{}
)

func (mps *MockProductService) Get(ctx context.Context, marketplacecode string) (domain.BasicProduct, error) {
	if marketplacecode == "fail" {
		return nil, errors.New("fail")
	}
	if marketplacecode == "simple" {
		return domain.SimpleProduct{
			BasicProductData: domain.BasicProductData{Title: "My Product Title", MarketPlaceCode: marketplacecode},
		}, nil
	}
	return domain.ConfigurableProduct{
		BasicProductData: domain.BasicProductData{Title: "My Configurable Product Title", MarketPlaceCode: marketplacecode},
		Variants: []domain.Variant{
			domain.Variant{
				BasicProductData: domain.BasicProductData{Title: "My Variant Title", MarketPlaceCode: marketplacecode + "_1"},
			},
		},
	}, nil
}

func TestViewController_Get(t *testing.T) {
	t.Skip("Broken")

	//redirectAware := new(mocks.RedirectAware)
	//renderAware := new(mocks.RenderAware)
	//errorAware := new(mocks.ErrorAware)
	//vc := getController(redirectAware, renderAware, errorAware)

	// Test 2: call with corrent name parameter and expect Rendering
	//expectedUrlTitleSimple := "my-product-title"
	//renderAware.On("Render", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("controller.productViewData")).Return(nil)
	//ctx.LoadParams(router.P{"marketplacecode": "simple", "name": expectedUrlTitleSimple})
	//vc.Get(web.ToRequest(ctx))
	//renderAware.AssertCalled(t, "Render", ctx, vc.Template, mock.AnythingOfType("controller.productViewData"))

	// Test 3: call with error by service
	//errorAware.On("Error", ctx, mock.Anything).Return(nil)
	//ctx.LoadParams(router.P{"marketplacecode": "fail", "name": "fail"})
	//vc.Get(web.ToRequest(ctx))
	//errorAware.AssertCalled(t, "Error", ctx, mock.Anything)

}

func TestViewController_ExpectRedirect(t *testing.T) {
	t.Skip("Broken")

	//ctx := web.NewContext()

	//redirectAware := new(mocks.RedirectAware)
	//renderAware := new(mocks.RenderAware)
	//errorAware := new(mocks.ErrorAware)
	//vc := getController(redirectAware, renderAware, errorAware)
	//
	//Test 1: call simple with wrong name and expect redirect
	//redirectAware.On("RedirectPermanentURL", mock.AnythingOfType("string")).Return(&web.RedirectResponse{})
	//ctx.LoadParams(router.P{"marketplacecode": "simple", "name": "testname"})
	//vc.Get(web.ToRequest(ctx))
	//redirectAware.AssertCalled(t, "RedirectPermanentURL", "/?marketplacecode=simple&name=my-product-title")
	//
	//Test 2: call configurable with wrong name and expect redirect
	//redirectAware.On("RedirectPermanentURL", mock.AnythingOfType("string")).Return(&web.RedirectResponse{})
	//ctx.LoadParams(router.P{"marketplacecode": "configurable", "name": "testname"})
	//vc.Get(web.ToRequest(ctx))
	//redirectAware.AssertCalled(t, "RedirectPermanentURL", "/?marketplacecode=configurable&name=my-configurable-product-title")
	//
	//Test 3: call configurable_with_variant with wrong name and expect redirect
	//redirectAware.On("RedirectPermanentURL", mock.AnythingOfType("string")).Return(&web.RedirectResponse{})
	//ctx.LoadParams(router.P{"marketplacecode": "configurable", "name": "testname", "variantcode": "configurable_1"})
	//vc.Get(web.ToRequest(ctx))
	//redirectAware.AssertCalled(t, "RedirectPermanentURL", "/?marketplacecode=configurable&name=my-variant-title&variantcode=configurable_1")
}

//func getController(redirectAware *mocks.RedirectAware, renderAware *mocks.RenderAware, errorAware *mocks.ErrorAware) *View {
//u, _ := url.Parse(`http://test/`)
//router := &router.Router{
//	RouterRegistry: router.NewRegistry(),
//}
//vc := &View{
//	ProductService: new(MockProductService),
//	RedirectAware:  redirectAware,
//	RenderAware:    renderAware,
//	ErrorAware:     errorAware,
//	URLService: &application.URLService{
//		Router: router,
//	},
//	Template: "product/product",
//	Router:   router,
//}
//vc.Router.SetBase(u)
//vc.Router.RouterRegistry.Route("/", `product.view(marketplacecode?="test", name?="test", variantcode?="test")`)
//
//return vc
//}
