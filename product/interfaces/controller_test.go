package interfaces

import (
	"context"
	"errors"
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/testutil"
	"flamingo/framework/web"
	"testing"
)

type (
	MockProductService struct{}
)

func (mps *MockProductService) Get(ctx context.Context, foreignID string) (*domain.Product, error) {
	if foreignID == "fail" {
		return nil, errors.New("fail")
	}

	return &domain.Product{
		ForeignID:    foreignID,
		InternalName: foreignID,
	}, nil
}

func TestViewController_Get(t *testing.T) {
	var redirectedTo, redirectedName string
	var tplname string
	var errorHappened bool

	vc := &ViewController{
		ProductService: new(MockProductService),
		RedirectAware: &testutil.MockRedirectAware{
			CbRedirect: func(name string, args map[string]string) web.Response {
				redirectedTo = "product.view"
				redirectedName = args["name"]
				return nil
			},
		},
		RenderAware: &testutil.MockRenderAware{
			CbRender: func(context web.Context, tpl string, data interface{}) web.Response {
				tplname = tpl
				return nil
			},
		},
		ErrorAware: &testutil.MockErrorAware{
			CbError: func(context web.Context, err error) web.Response {
				errorHappened = true
				return nil
			},
		},
	}
	ctx := web.NewContext()

	ctx.LoadParams(router.P{"uid": "test", "name": "testname"})
	response := vc.Get(ctx)

	if redirectedTo != "product.view" {
		t.Errorf("Expected redirect to product.view, not %q", redirectedTo)
	}

	if redirectedName != "test" {
		t.Errorf("Expected redirect to name test, not %q", redirectedTo)
	}

	if response != nil {
		t.Errorf("Expected mocked response to be nil, not %T", response)
	}

	ctx.LoadParams(router.P{"uid": "test", "name": "test"})
	response = vc.Get(ctx)

	if tplname != "product/simple" {
		t.Errorf("expected to render %q", tplname)
	}

	if response != nil {
		t.Errorf("Expected mocked response to be nil, not %T", response)
	}

	ctx.LoadParams(router.P{"uid": "fail", "name": "fail"})
	response = vc.Get(ctx)

	if !errorHappened {
		t.Error("expected to error for 'fail' product")
	}

	if response != nil {
		t.Errorf("Expected mocked response to be nil, not %T", response)
	}
}
