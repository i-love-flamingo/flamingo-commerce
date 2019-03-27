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
			{
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

// func getController(redirectAware *mocks.RedirectAware, renderAware *mocks.RenderAware, errorAware *mocks.ErrorAware) *View {
// 	u, _ := url.Parse(`http://test/`)
// 	r := &web.Router{
// 		RouterRegistry: web.NewRegistry(),
// 	}
// 	vc := &View{
// 		ProductService: new(MockProductService),
//
// 		RedirectAware: redirectAware,
// 		RenderAware:   renderAware,
// 		ErrorAware:    errorAware,
// 		URLService: &application.URLService{
// 			Router: r,
// 		},
// 		Template: "product/product",
// 		Router:   r,
// 	}
// 	vc.Router.SetBase(u)
// 	vc.Router.RouterRegistry.Route("/", `product.view(marketplacecode?="test", name?="test", variantcode?="test")`)
//
// 	return vc
// }

// This test is added to help better understand what variantSelection method is doing.
// Unfortunately the assert library is unable to compare []variantSelection slices because
// the order of values in some of the fields is not guaranteed (because how Go maps work).
// For this reason we try to compare manually using helper functions.
// func TestViewController_variantSelection(t *testing.T) {
// 	vc := getController(nil, nil, nil)
//
// 	testCases := []struct {
// 		variantVariationAttributesSorting map[string][]string
// 		variants                          []domain.Variant
// 		activeVariant                     *domain.Variant
//
// 		out variantSelection
// 	}{
// 		{
// 			variantVariationAttributesSorting: map[string][]string{"color": {"red", "blue"}, "size": {"s", "m", "l"}},
// 			variants: []domain.Variant{
// 				{
// 					BasicProductData: domain.BasicProductData{
// 						Attributes: domain.Attributes{"color": {Label: "Red", RawValue: "red"}, "size": {Label: "S", RawValue: "s"}},
// 						StockLevel: "",
// 					},
// 				},
// 				{
// 					BasicProductData: domain.BasicProductData{
// 						Attributes: domain.Attributes{"color": {Label: "Red", RawValue: "red"}, "size": {Label: "M", RawValue: "m"}},
// 						StockLevel: "out",
// 					},
// 				},
// 				{
// 					BasicProductData: domain.BasicProductData{
// 						Attributes: domain.Attributes{"color": {Label: "Red", RawValue: "red"}, "size": {Label: "L", RawValue: "l"}},
// 						StockLevel: "low",
// 					},
// 				},
// 				{
// 					BasicProductData: domain.BasicProductData{
// 						Attributes: domain.Attributes{"color": {Label: "Blue", RawValue: "blue"}, "size": {Label: "S", RawValue: "s"}},
// 						StockLevel: "high",
// 					},
// 				},
// 				{
// 					BasicProductData: domain.BasicProductData{
// 						Attributes: domain.Attributes{"color": {Label: "Blue", RawValue: "blue"}, "size": {Label: "M", RawValue: "m"}},
// 					},
// 				},
// 				//{
// 				//	BasicProductData: domain.BasicProductData{
// 				//		Attributes: domain.Attributes{"color": {Label: "Blue", RawValue: "blue"}, "size": {Label: "L", RawValue: "l"}},
// 				//	},
// 				//},
// 			},
// 			activeVariant: &domain.Variant{
// 				BasicProductData: domain.BasicProductData{
// 					Attributes: map[string]domain.Attribute{
// 						"color": {Label: "Blue", RawValue: "blue"}, "size": {Label: "M", RawValue: "m"},
// 					},
// 				},
// 			},
//
// 			out: variantSelection{
// 				Attributes: []viewVariantAttribute{
// 					{
// 						Key: "color", Title: "Color",
// 						Options: []viewVariantOption{
// 							{
// 								Key: "red", Title: "Red",
// 								Combinations: map[string][]string{"size": {"s", "m", "l"}},
// 							},
// 							{
// 								Key: "blue", Title: "Blue",
// 								Combinations: map[string][]string{"size": {"s", "m" /*, "l"*/}},
// 								Selected:     true,
// 							},
// 						},
// 					},
// 					{
// 						Key: "size", Title: "Size",
// 						Options: []viewVariantOption{
// 							{
// 								Key: "s", Title: "S",
// 								Combinations: map[string][]string{"color": {"red", "blue"}},
// 							},
// 							{
// 								Key: "m", Title: "M",
// 								Combinations: map[string][]string{"color": {"red", "blue"}},
// 								Selected:     true,
// 							},
// 							{
// 								Key: "l", Title: "L",
// 								Combinations: map[string][]string{"color": {"red" /*, "blue"*/}},
// 							},
// 						},
// 					},
// 				},
// 				Variants: []viewVariant{
// 					{
// 						Attributes: map[string]string{"color": "red", "size": "s"},
// 						URL:        "/?marketplacecode=&name=&variantcode=",
// 						InStock:    false,
// 					},
// 					{
// 						Attributes: map[string]string{"color": "red", "size": "m"},
// 						URL:        "/?marketplacecode=&name=&variantcode=",
// 						InStock:    false,
// 					},
// 					{
// 						Attributes: map[string]string{"color": "red", "size": "l"},
// 						URL:        "/?marketplacecode=&name=&variantcode=",
// 						InStock:    true,
// 					},
// 					{
// 						Attributes: map[string]string{"color": "blue", "size": "s"},
// 						URL:        "/?marketplacecode=&name=&variantcode=",
// 						InStock:    true,
// 					},
// 					{
// 						Attributes: map[string]string{"color": "blue", "size": "m"},
// 						URL:        "/?marketplacecode=&name=&variantcode=",
// 						InStock:    false,
// 					},
// 					//{
// 					//	Attributes: map[string]string{"color": "blue", "size": "l"},
// 					//	Url:        "/?marketplacecode=&name=&variantcode=",
// 					//},
// 				},
// 			},
// 		},
// 	}
//
// 	for _, tc := range testCases {
//
// 		var variantVariationAttributes []string
// 		for key := range tc.variantVariationAttributesSorting {
// 			variantVariationAttributes = append(variantVariationAttributes, key)
// 		}
//
// 		configurableProduct := domain.ConfigurableProduct{
// 			VariantVariationAttributes:        variantVariationAttributes,
// 			VariantVariationAttributesSorting: tc.variantVariationAttributesSorting,
// 			Variants:                          tc.variants,
// 		}
//
// 		vs := vc.variantSelection(configurableProduct, tc.activeVariant)
//
// 		assert.Len(t, vs.Attributes, len(tc.out.Attributes))
// 		assert.Len(t, vs.Variants, len(tc.out.Variants))
//
// 		aEqual := true
// 		for _, a1 := range vs.Attributes {
// 			var aFound bool
// 			for _, a2 := range tc.out.Attributes {
// 				aFound = aFound || viewVariantAttributesEqual(t, a1, a2)
// 			}
// 			assert.True(t, aFound, "attributes not the same: attribute '%v' not found in '%v'", a1, tc.out.Attributes)
// 			aEqual = aEqual && aFound
// 		}
// 		assert.True(t, aEqual, "attributes do not match")
//
// 		vEqual := true
// 		for _, v1 := range vs.Variants {
// 			var vFound bool
// 			for _, v2 := range tc.out.Variants {
// 				vFound = vFound || viewVariantsEqual(t, v1, v2)
// 			}
// 			assert.True(t, vFound, "variants not the same: variant '%v' not found in '%v'", v1, tc.out.Variants)
// 			vEqual = vEqual && vFound
// 		}
// 		assert.True(t, vEqual, "variants do not match")
// 	}
// }
//
// func viewVariantAttributesEqual(t *testing.T, a1, a2 viewVariantAttribute) bool {
// 	t.Helper()
// 	if a1.Title != a2.Title {
// 		return false
// 	}
// 	if a1.Key != a2.Key {
// 		return false
// 	}
// 	if len(a1.Options) != len(a2.Options) {
// 		return false
// 	}
// 	oEqual := true
// 	for _, o1 := range a1.Options {
// 		var oFound bool
// 		for _, o2 := range a2.Options {
// 			oFound = oFound || viewVariantOptionsEqual(t, o1, o2)
// 		}
// 		oEqual = oEqual && oFound
// 	}
// 	return oEqual
// }
//
// func viewVariantOptionsEqual(t *testing.T, o1, o2 viewVariantOption) bool {
// 	t.Helper()
// 	if o1.Title != o2.Title {
// 		return false
// 	}
// 	if o1.Key != o2.Key {
// 		return false
// 	}
// 	if o1.Selected != o2.Selected {
// 		return false
// 	}
// 	if len(o1.Combinations) != len(o2.Combinations) {
// 		return false
// 	}
// 	for k1, v1 := range o1.Combinations {
// 		v2, ok := o2.Combinations[k1]
// 		if !ok {
// 			return false
// 		}
// 		if !slicesEqual(t, v1, v2) {
// 			return false
// 		}
// 	}
// 	return true
// }
//
// func slicesEqual(t *testing.T, s1, s2 []string) bool {
// 	t.Helper()
// 	if len(s1) != len(s2) {
// 		return false
// 	}
// 	for _, v1 := range s1 {
// 		var sFound bool
// 		for _, v2 := range s2 {
// 			sFound = sFound || v1 == v2
// 		}
// 		if !sFound {
// 			return false
// 		}
// 	}
// 	return true
// }
//
// func viewVariantsEqual(t *testing.T, variant1, variant2 viewVariant) bool {
// 	t.Helper()
// 	if variant1.Title != variant2.Title {
// 		return false
// 	}
// 	if variant1.Url != variant2.Url {
// 		return false
// 	}
// 	if variant1.Marketplacecode != variant2.Marketplacecode {
// 		return false
// 	}
// 	if variant1.InStock != variant2.InStock {
// 		return false
// 	}
// 	if len(variant1.Attributes) != len(variant2.Attributes) {
// 		return false
// 	}
// 	for k1, v1 := range variant1.Attributes {
// 		v2, ok := variant2.Attributes[k1]
// 		if !ok {
// 			return false
// 		}
// 		if v1 != v2 {
// 			return false
// 		}
// 	}
// 	return true
// }
