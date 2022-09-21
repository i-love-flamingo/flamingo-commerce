package controller

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	MockProductService struct{}
	routes             struct{}
)

// Routes definition for the reverse routing functionality needed
func (r *routes) Routes(registry *web.RouterRegistry) {
	_, _ = registry.Route("/", `product.view(marketplacecode?="test", name?="test", variantcode?="test")`)
	registry.HandleGet("product.view", nil)
}

func getController() *View {
	r := new(web.Router)
	r.Inject(
		&struct {
			// base url configuration
			Scheme      string `inject:"config:flamingo.router.scheme,optional"`
			Host        string `inject:"config:flamingo.router.host,optional"`
			Path        string `inject:"config:flamingo.router.path,optional"`
			External    string `inject:"config:flamingo.router.external,optional"`
			SessionName string `inject:"config:flamingo.session.name,optional"`
		}{
			Scheme: "http://",
			Host:   "test",
		},
		nil,
		nil,
		func() []web.Filter {
			return nil
		},
		func() []web.RoutesModule {
			return []web.RoutesModule{&routes{}}
		},
		new(flamingo.NullLogger),
		nil,
		nil,
	)
	// create a new handler to initialize the router registry
	r.Handler()

	vc := &View{
		ProductService: new(MockProductService),
		Responder:      new(web.Responder),
		URLService:     new(application.URLService).Inject(r, nil),
		Template:       "product/product",
		Router:         r,
	}

	return vc
}

func (mps *MockProductService) Get(_ context.Context, marketplacecode string) (domain.BasicProduct, error) {
	if marketplacecode == "fail" {
		return nil, errors.New("fail")
	}
	if marketplacecode == "not_found" {
		return nil, domain.ProductNotFound{MarketplaceCode: "not_found"}
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
	vc := getController()

	// call with correct name parameter and expect Rendering
	ctx := context.Background()
	r := web.CreateRequest(&http.Request{}, nil)
	r.Request().URL = &url.URL{}
	r.Params = web.RequestParams{
		"marketplacecode": "simple",
		"name":            "my-product-title",
	}

	result := vc.Get(ctx, r)
	require.IsType(t, &web.RenderResponse{}, result)

	renderResponse := result.(*web.RenderResponse)
	assert.Equal(t, "product/product", renderResponse.Template)
	require.IsType(t, productViewData{}, renderResponse.DataResponse.Data)
	p, _ := new(MockProductService).Get(context.Background(), "simple")
	assert.Equal(t, p, renderResponse.DataResponse.Data.(productViewData).Product)
}

func TestViewController_GetNotFound(t *testing.T) {
	tests := []struct {
		name            string
		marketPlaceCode string
		expectedStatus  uint
	}{
		{
			name:            "error",
			marketPlaceCode: "fail",
			expectedStatus:  http.StatusInternalServerError,
		},
		{
			name:            "not found",
			marketPlaceCode: "not_found",
			expectedStatus:  http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		vc := getController()

		// call with correct name parameter and expect Rendering
		ctx := context.Background()
		r := web.CreateRequest(&http.Request{}, nil)
		r.Request().URL = &url.URL{}
		r.Params = web.RequestParams{
			"marketplacecode": tt.marketPlaceCode,
			"name":            tt.marketPlaceCode,
		}

		result := vc.Get(ctx, r)
		require.IsType(t, &web.ServerErrorResponse{}, result)
		assert.Equal(t, int(tt.expectedStatus), int(result.(*web.ServerErrorResponse).Response.Status))
	}
}

func TestViewController_ExpectRedirect(t *testing.T) {
	tests := []struct {
		name            string
		marketPlaceCode string
		productName     string
		variantCode     string
		expectedStatus  uint
		expectedURL     string
	}{
		{
			name:            "call simple with wrong name and expect redirect",
			marketPlaceCode: "simple",
			productName:     "testname",
			expectedStatus:  http.StatusMovedPermanently,
			expectedURL:     "/?marketplacecode=simple&name=my-product-title",
		},
		{
			name:            "call configurable with wrong name and expect redirect",
			marketPlaceCode: "configurable",
			productName:     "testname",
			expectedStatus:  http.StatusMovedPermanently,
			expectedURL:     "/?marketplacecode=configurable&name=my-configurable-product-title",
		},
		{
			name:            "call configurable_with_variant with wrong name and expect redirect",
			marketPlaceCode: "configurable",
			productName:     "testname",
			variantCode:     "configurable_1",
			expectedStatus:  http.StatusMovedPermanently,
			expectedURL:     "/?marketplacecode=configurable&name=my-variant-title&variantcode=configurable_1",
		},
	}
	for _, tt := range tests {
		vc := getController()

		r := web.CreateRequest(&http.Request{}, nil)
		r.Request().URL = &url.URL{}
		r.Params = web.RequestParams{
			"marketplacecode": tt.marketPlaceCode,
			"name":            tt.productName,
		}
		if tt.variantCode != "" {
			r.Params["variantcode"] = "configurable_1"
		}
		result := vc.Get(context.Background(), r)
		require.IsType(t, &web.URLRedirectResponse{}, result)
		redirectResponse := result.(*web.URLRedirectResponse)
		assert.Equal(t, int(tt.expectedStatus), int(redirectResponse.Response.Status))
		assert.Equal(t, tt.expectedURL, redirectResponse.URL.String())
	}
}

// This test is added to help better understand what variantSelection method is doing.
// Unfortunately the assert library is unable to compare []variantSelection slices because
// the order of values in some of the fields is not guaranteed (because how Go maps work).
// For this reason we try to compare manually using helper functions.
func TestViewController_variantSelection(t *testing.T) {
	vc := getController()

	testCases := []struct {
		variantVariationAttributesSorting map[string][]string
		variants                          []domain.Variant
		activeVariant                     *domain.Variant

		out variantSelection
	}{
		{
			variantVariationAttributesSorting: map[string][]string{"color": {"red", "blue"}, "size": {"s", "m", "l"}},
			variants: []domain.Variant{
				{
					BasicProductData: domain.BasicProductData{
						Attributes: domain.Attributes{"color": {Label: "Red", CodeLabel: "Colour", RawValue: "red"}, "size": {Label: "S", CodeLabel: "Clothing Size", RawValue: "s"}},
						Stock:      getStock(false, domain.StockLevelOutOfStock, 0),
					},
				},
				{
					BasicProductData: domain.BasicProductData{
						Attributes: domain.Attributes{"color": {Label: "Red", CodeLabel: "Colour", RawValue: "red"}, "size": {Label: "M", CodeLabel: "Clothing Size", RawValue: "m"}},
						Stock:      getStock(false, domain.StockLevelOutOfStock, 0),
					},
				},
				{
					BasicProductData: domain.BasicProductData{
						Attributes: domain.Attributes{"color": {Label: "Red", CodeLabel: "Colour", RawValue: "red"}, "size": {Label: "L", CodeLabel: "Clothing Size", RawValue: "l"}},
						Stock:      getStock(true, domain.StockLevelLowStock, 100),
					},
				},
				{
					BasicProductData: domain.BasicProductData{
						Attributes: domain.Attributes{"color": {Label: "Blue", CodeLabel: "Colour", RawValue: "blue"}, "size": {Label: "S", CodeLabel: "Clothing Size", RawValue: "s"}},
						Stock:      getStock(true, domain.StockLevelInStock, 999),
					},
				},
				{
					BasicProductData: domain.BasicProductData{
						Attributes: domain.Attributes{"color": {Label: "Blue", CodeLabel: "Colour", RawValue: "blue"}, "size": {Label: "M", CodeLabel: "Clothing Size", RawValue: "m"}},
						Stock:      getStock(false, domain.StockLevelOutOfStock, 0),
					},
				},
			},
			activeVariant: &domain.Variant{
				BasicProductData: domain.BasicProductData{
					Attributes: map[string]domain.Attribute{
						"color": {Label: "Blue", CodeLabel: "Colour", RawValue: "blue"}, "size": {Label: "M", CodeLabel: "Clothing Size", RawValue: "m"},
					},
				},
			},

			out: variantSelection{
				Attributes: []viewVariantAttribute{
					{
						Key:       "color",
						Title:     "Color",
						CodeLabel: "Colour",
						Options: []viewVariantOption{
							{
								Key: "red", Title: "Red",
								Combinations: map[string][]string{"size": {"s", "m", "l"}},
							},
							{
								Key: "blue", Title: "Blue",
								Combinations: map[string][]string{"size": {"s", "m"}},
								Selected:     true,
							},
						},
					},
					{
						Key:       "size",
						Title:     "Size",
						CodeLabel: "Clothing Size",
						Options: []viewVariantOption{
							{
								Key: "s", Title: "S",
								Combinations: map[string][]string{"color": {"red", "blue"}},
							},
							{
								Key: "m", Title: "M",
								Combinations: map[string][]string{"color": {"red", "blue"}},
								Selected:     true,
							},
							{
								Key: "l", Title: "L",
								Combinations: map[string][]string{"color": {"red"}},
							},
						},
					},
				},
				Variants: []viewVariant{
					{
						Attributes: map[string]string{"color": "red", "size": "s"},
						URL:        "/?marketplacecode=&name=&variantcode=",
						InStock:    false,
					},
					{
						Attributes: map[string]string{"color": "red", "size": "m"},
						URL:        "/?marketplacecode=&name=&variantcode=",
						InStock:    false,
					},
					{
						Attributes: map[string]string{"color": "red", "size": "l"},
						URL:        "/?marketplacecode=&name=&variantcode=",
						InStock:    true,
					},
					{
						Attributes: map[string]string{"color": "blue", "size": "s"},
						URL:        "/?marketplacecode=&name=&variantcode=",
						InStock:    true,
					},
					{
						Attributes: map[string]string{"color": "blue", "size": "m"},
						URL:        "/?marketplacecode=&name=&variantcode=",
						InStock:    false,
					},
				},
			},
		},
		{
			variantVariationAttributesSorting: map[string][]string{"volume": {"500", "1"}},
			variants: []domain.Variant{
				{
					BasicProductData: domain.BasicProductData{
						Attributes: domain.Attributes{"volume": {CodeLabel: "Volume", RawValue: "500", UnitCode: "MILLILITRE"}},
						Stock:      getStock(true, domain.StockLevelInStock, 999),
					},
				},
				{
					BasicProductData: domain.BasicProductData{
						Attributes: domain.Attributes{"volume": {CodeLabel: "Volume", RawValue: "1", UnitCode: "LITRE"}},
						Stock:      getStock(true, domain.StockLevelInStock, 999),
					},
				},
			},
			activeVariant: &domain.Variant{
				BasicProductData: domain.BasicProductData{
					Attributes: map[string]domain.Attribute{
						"volume": {CodeLabel: "Volume", RawValue: "500", UnitCode: "MILLILITRE"},
					},
				},
			},

			out: variantSelection{
				Attributes: []viewVariantAttribute{
					{
						Key:       "volume",
						Title:     "Volume",
						CodeLabel: "Volume",
						Options: []viewVariantOption{
							{
								Key:      "500",
								Title:    "500",
								Selected: true,
								Unit:     "MILLILITRE",
							},
							{
								Key:   "1",
								Title: "1",
								Unit:  "LITRE",
							},
						},
					},
				},
				Variants: []viewVariant{
					{
						Attributes: map[string]string{"volume": "500"},
						URL:        "/?marketplacecode=&name=&variantcode=",
						InStock:    true,
					},
					{
						Attributes: map[string]string{"volume": "1"},
						URL:        "/?marketplacecode=&name=&variantcode=",
						InStock:    true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {

		var variantVariationAttributes []string
		for key := range tc.variantVariationAttributesSorting {
			variantVariationAttributes = append(variantVariationAttributes, key)
		}

		configurableProduct := domain.ConfigurableProduct{
			VariantVariationAttributes:        variantVariationAttributes,
			VariantVariationAttributesSorting: tc.variantVariationAttributesSorting,
			Variants:                          tc.variants,
		}

		vs := vc.variantSelection(configurableProduct, tc.activeVariant)

		assert.Len(t, vs.Attributes, len(tc.out.Attributes))
		assert.Len(t, vs.Variants, len(tc.out.Variants))

		aEqual := true
		for _, a1 := range vs.Attributes {
			var aFound bool
			for _, a2 := range tc.out.Attributes {
				aFound = aFound || viewVariantAttributesEqual(t, a1, a2)
			}
			assert.True(t, aFound, "attributes not the same: attribute '%v' not found in '%v'", a1, tc.out.Attributes)
			aEqual = aEqual && aFound
		}
		assert.True(t, aEqual, "attributes do not match")

		vEqual := true
		for _, v1 := range vs.Variants {
			var vFound bool
			for _, v2 := range tc.out.Variants {
				vFound = vFound || viewVariantsEqual(t, v1, v2)
			}
			assert.True(t, vFound, "variants not the same: variant '%v' not found in '%v'", v1, tc.out.Variants)
			vEqual = vEqual && vFound
		}
		assert.True(t, vEqual, "variants do not match")
	}
}

func viewVariantAttributesEqual(t *testing.T, a1, a2 viewVariantAttribute) bool {
	t.Helper()
	if a1.Title != a2.Title {
		return false
	}
	if a1.Key != a2.Key {
		return false
	}
	if len(a1.Options) != len(a2.Options) {
		return false
	}
	oEqual := true
	for _, o1 := range a1.Options {
		var oFound bool
		for _, o2 := range a2.Options {
			oFound = oFound || viewVariantOptionsEqual(t, o1, o2)
		}
		oEqual = oEqual && oFound
	}
	return oEqual
}

func viewVariantOptionsEqual(t *testing.T, o1, o2 viewVariantOption) bool {
	t.Helper()
	if o1.Title != o2.Title {
		return false
	}
	if o1.Key != o2.Key {
		return false
	}
	if o1.Selected != o2.Selected {
		return false
	}
	if o1.Unit != o2.Unit {
		return false
	}
	if len(o1.Combinations) != len(o2.Combinations) {
		return false
	}
	for k1, v1 := range o1.Combinations {
		v2, ok := o2.Combinations[k1]
		if !ok {
			return false
		}
		if !slicesEqual(t, v1, v2) {
			return false
		}
	}
	return true
}

func slicesEqual(t *testing.T, s1, s2 []string) bool {
	t.Helper()
	if len(s1) != len(s2) {
		return false
	}
	for _, v1 := range s1 {
		var sFound bool
		for _, v2 := range s2 {
			sFound = sFound || v1 == v2
		}
		if !sFound {
			return false
		}
	}
	return true
}

func viewVariantsEqual(t *testing.T, variant1, variant2 viewVariant) bool {
	t.Helper()
	if variant1.Title != variant2.Title {
		return false
	}
	if variant1.URL != variant2.URL {
		return false
	}
	if variant1.Marketplacecode != variant2.Marketplacecode {
		return false
	}
	if variant1.InStock != variant2.InStock {
		return false
	}
	if len(variant1.Attributes) != len(variant2.Attributes) {
		return false
	}
	for k1, v1 := range variant1.Attributes {
		v2, ok := variant2.Attributes[k1]
		if !ok {
			return false
		}
		if v1 != v2 {
			return false
		}
	}
	return true
}

func getStock(inStock bool, level string, amount int) []domain.Stock {
	stock := make([]domain.Stock, 0)

	stock = append(stock, domain.Stock{
		Amount:       amount,
		InStock:      inStock,
		Level:        level,
		DeliveryCode: "",
	})

	return stock
}
