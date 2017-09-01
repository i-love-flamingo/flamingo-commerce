package controller

import (
	"flamingo/core/breadcrumbs"
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"log"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		domain.ProductService   `inject:""`

		Template string         `inject:"config:core.product.view.template"`
		Router   *router.Router `inject:""`
	}

	DebugData View

	// ProductViewData is used for product rendering
	ProductViewData struct {
		// simple / configurable / configurable_with_variant
		RenderContext       string
		SimpleProduct       domain.SimpleProduct
		ConfigurableProduct domain.ConfigurableProduct
		ActiveVariant       domain.Variant
		VariantSelected     bool
		VariantSelection    VariantSelection
	}

	// VariantSelection for templating
	VariantSelection struct {
		Attributes []ViewVariantAttribute
		Variants   []ViewVariant
	}

	ViewVariantAttribute struct {
		Key     string
		Title   string
		Options []ViewVariantOption
	}

	ViewVariantOption struct {
		Key          string
		Title        string
		Combinations map[string][]string
		Selected     bool
	}

	ViewVariant struct {
		Attributes      map[string]string
		Marketplacecode string
		Title           string
		Url             string
	}
)

// Get Response for Product matching sku param
func (vc *View) Get(c web.Context) web.Response {
	product, err := vc.ProductService.Get(c, c.MustParam1("marketplacecode"))
	skipnamecheck, _ := c.Param1("skipnamecheck")

	// catch error
	if err != nil {
		switch errors.Cause(err).(type) {
		case domain.ProductNotFound:
			return vc.ErrorNotFound(c, err)

		default:
			return vc.Error(c, err)
		}
	}

	var viewData ProductViewData

	// 1. Handle Configurables
	if product.Type() == "configurable" {
		configurableProduct := product.(domain.ConfigurableProduct)
		var activeVariant *domain.Variant

		variantCode, err := c.Param1("variantcode")

		if err == nil {
			log.Println("get variant by " + variantCode)
			activeVariant, _ = configurableProduct.Variant(variantCode)
		}
		if activeVariant == nil {
			// 1.A. No variant selected
			// normalize URL
			urlName := makeUrlTitle(product.BaseData().Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "name": urlName})
			}
			viewData = ProductViewData{ConfigurableProduct: configurableProduct, VariantSelected: false, RenderContext: "configurable"}
		} else {
			// 1.B. Variant selected
			// normalize URL
			urlName := makeUrlTitle(activeVariant.BasicProductData.Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "variantcode": variantCode, "name": urlName})
			}
			log.Printf("Variant Price %v / %v", activeVariant.ActivePrice.Default, activeVariant.ActivePrice)
			viewData = ProductViewData{ConfigurableProduct: configurableProduct, ActiveVariant: *activeVariant, VariantSelected: true, RenderContext: "configurable_with_activevariant"}
		}

	} else {
		// 2. Handle Simples
		// normalize URL
		urlName := makeUrlTitle(product.BaseData().Title)
		if urlName != c.MustParam1("name") && skipnamecheck == "" {
			return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "name": urlName})
		}

		simpleProduct := product.(domain.SimpleProduct)
		viewData = ProductViewData{SimpleProduct: simpleProduct, RenderContext: "simple"}
	}

	for _, category := range product.BaseData().CategoryPath {
		breadcrumbs.Add(c, breadcrumbs.Crumb{Title: category})
	}

	breadcrumbs.Add(c, breadcrumbs.Crumb{
		Title: product.BaseData().Title,
		URL:   vc.Router.URL("product.view", router.P{"marketplacecode": product.BaseData().MarketPlaceCode, "name": product.BaseData().Title}).String(),
	})

	viewData.VariantSelection = VariantSelection{
		Attributes: []ViewVariantAttribute{
			{
				Key:   "baseColor",
				Title: "Color",
				Options: []ViewVariantOption{
					{
						Title: "Red",
						Key:   "red",
						Combinations: map[string][]string{
							"clothingSize": {"l", "xl"},
						},
                        Selected: true,
					},
					{
						Title: "Green",
						Key:   "green",
						Combinations: map[string][]string{
							"clothingSize": {"xl", "m"},
						},
					},
				},
			},
			{
				Key:   "clothingSize",
				Title: "Size",
				Options: []ViewVariantOption{
					{
						Title: "Size M",
						Key:   "m",
						Combinations: map[string][]string{
							"baseColor": {"green"},
						},
					},
					{
						Title: "Size L",
						Key:   "l",
						Combinations: map[string][]string{
							"baseColor": {"red"},
						},
                        Selected: true,
					},
					{
						Title: "Size XL",
						Key:   "xl",
						Combinations: map[string][]string{
							"baseColor": {"red", "green"},
						},
					},
				},
			},
		},
		Variants: []ViewVariant{
			{
				Title:           "Red Shirt L",
				Url:             "/",
				Marketplacecode: "red-shirt-l",
				Attributes: map[string]string{
					"baseColor":    "red",
					"clothingSize": "l",
				},
			},
			{
				Title:           "Red Shirt XL",
				Url:             "/",
				Marketplacecode: "red-shirt-xl",
				Attributes: map[string]string{
					"baseColor":    "red",
					"clothingSize": "xl",
				},
			},
			{
				Title:           "Green Shirt XL",
				Url:             "/",
				Marketplacecode: "green-shirt-xl",
				Attributes: map[string]string{
					"baseColor":    "green",
					"clothingSize": "xl",
				},
			},
			{
				Title:           "Green Shirt M",
				Url:             "/",
				Marketplacecode: "green-shirt-m",
				Attributes: map[string]string{
					"baseColor":    "green",
					"clothingSize": "m",
				},
			},
		},
	}

	return vc.Render(c, vc.Template, viewData)
}

func makeUrlTitle(title string) string {
	newTitle := strings.ToLower(strings.Replace(title, " ", "-", -1))
	newTitle = url.QueryEscape(newTitle)

	return newTitle
}

type dbgrender struct {
	data interface{}
}

func (d *dbgrender) Render(context web.Context, tpl string, data interface{}) web.Response {
	d.data = data
	return nil
}

func (d *DebugData) Get(c web.Context) web.Response {
	vc := (*View)(d)
	r := &dbgrender{}
	vc.RenderAware = r

	params := c.ParamAll()
	params["skipnamecheck"] = "1"
	params["name"] = ""
	c.LoadParams(params)
	rr := vc.Get(c)

	log.Println(rr)

	return &web.JSONResponse{
		Data: r.data,
	}
}
