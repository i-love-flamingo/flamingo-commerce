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

	"sort"

	"flamingo/core/pug_template/pugjs"

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

func (vc *View) variantSelection(configurable domain.ConfigurableProduct, activeVariant *domain.Variant) VariantSelection {
	var variants VariantSelection
	combinations := make(map[string]map[string]map[string]map[string]bool)

	for _, attribute := range configurable.VariantVariationAttributes {
		// attribute -> value -> combinableAttribute -> combinaleValue -> true
		combinations[attribute] = make(map[string]map[string]map[string]bool)

		for _, variant := range configurable.Variants {
			for _, subattribute := range configurable.VariantVariationAttributes {
				if subattribute != attribute {
					if variant.Attributes[attribute] == nil {
						continue
					}

					if combinations[attribute][variant.Attributes[attribute].(string)] == nil {
						combinations[attribute][variant.Attributes[attribute].(string)] = make(map[string]map[string]bool)
					}

					if variant.Attributes[subattribute] != nil {
						if combinations[attribute][variant.Attributes[attribute].(string)][subattribute] == nil {
							combinations[attribute][variant.Attributes[attribute].(string)][subattribute] = make(map[string]bool)
						}
						combinations[attribute][variant.Attributes[attribute].(string)][subattribute][variant.Attributes[subattribute].(string)] = true
					}
				}
			}
		}
	}

	for code, attribute := range combinations {
		viewVariantAttribute := ViewVariantAttribute{
			Key:   code,
			Title: strings.Title(code),
		}

		for optionCode, option := range attribute {
			combinations := make(map[string][]string)
			for cattr, cvalues := range option {
				for cvalue := range cvalues {
					combinations[cattr] = append(combinations[cattr], cvalue)
				}
			}

			var selected bool
			if activeVariant != nil && activeVariant.Attributes[code] == optionCode {
				selected = true
			}
			viewVariantAttribute.Options = append(viewVariantAttribute.Options, ViewVariantOption{
				Key:          optionCode,
				Title:        strings.Title(optionCode),
				Selected:     selected,
				Combinations: combinations,
			})
		}

		variants.Attributes = append(variants.Attributes, viewVariantAttribute)
	}

	for _, variant := range configurable.Variants {
		urlName := makeUrlTitle(variant.BasicProductData.Title)
		variantUrl := vc.Router.URL("product.view", router.P{"marketplacecode": configurable.MarketPlaceCode, "variantcode": variant.MarketPlaceCode, "name": urlName}).String()

		attributes := make(map[string]string)

		for _, attr := range variants.Attributes {
			if variant.Attributes[attr.Key] == nil {
				continue
			}
			attributes[attr.Key] = variant.Attributes[attr.Key].(string)
		}

		variants.Variants = append(variants.Variants, ViewVariant{
			Title:           variant.Title,
			Marketplacecode: variant.MarketPlaceCode,
			Url:             variantUrl,
			Attributes:      attributes,
		})
	}

	return variants
}

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

		if err != nil {
			// 1.A. No variant selected
			// normalize URL
			urlName := makeUrlTitle(product.BaseData().Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "name": urlName})
			}
			viewData = ProductViewData{
				ConfigurableProduct: configurableProduct,
				VariantSelected:     false,
				RenderContext:       "configurable",
			}
		} else {
			log.Println("get variant by " + variantCode)
			activeVariant, _ = configurableProduct.Variant(variantCode)
			// 1.B. Variant selected
			// normalize URL
			urlName := makeUrlTitle(activeVariant.BasicProductData.Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "variantcode": variantCode, "name": urlName})
			}
			log.Printf("Variant Price %v / %v", activeVariant.ActivePrice.Default, activeVariant.ActivePrice)
			viewData = ProductViewData{
				ConfigurableProduct: configurableProduct,
				ActiveVariant:       *activeVariant,
				VariantSelected:     true,
				RenderContext:       "configurable_with_activevariant",
			}
		}

		viewData.VariantSelection = vc.variantSelection(configurableProduct, activeVariant)

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

	paths := product.BaseData().CategoryPath
	sort.Strings(paths)
	var stringHead string
	for _, p := range paths {
		if strings.HasPrefix(p, stringHead) {
			breadcrumbs.Add(c, breadcrumbs.Crumb{Title: p[len(stringHead):]})
			stringHead = p + "/"
		}
	}

	breadcrumbs.Add(c, breadcrumbs.Crumb{
		Title: product.BaseData().Title,
		URL:   vc.Router.URL("product.view", router.P{"marketplacecode": product.BaseData().MarketPlaceCode, "name": product.BaseData().Title}).String(),
	})

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
	vc.Get(c)

	return &web.JSONResponse{
		Data: pugjs.Convert(r.data),
	}
}
