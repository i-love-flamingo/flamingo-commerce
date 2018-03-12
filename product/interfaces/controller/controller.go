package controller

import (
	"strings"

	"go.aoe.com/flamingo/core/breadcrumbs"
	"go.aoe.com/flamingo/core/category"
	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/product/application"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		domain.ProductService   `inject:""`
		UrlService              *application.UrlService `inject:""`

		Template string         `inject:"config:core.product.view.template"`
		Router   *router.Router `inject:""`
	}

	// productViewData is used for product rendering
	productViewData struct {
		// simple / configurable / configurable_with_variant
		RenderContext       string
		SimpleProduct       domain.SimpleProduct
		ConfigurableProduct domain.ConfigurableProduct
		VariantSelected     bool
		VariantSelection    variantSelection
	}

	// variantSelection for templating
	variantSelection struct {
		Attributes []viewVariantAttribute
		Variants   []viewVariant
	}

	viewVariantAttribute struct {
		Key     string
		Title   string
		Options []viewVariantOption
	}

	viewVariantOption struct {
		Key          string
		Title        string
		Combinations map[string][]string
		Selected     bool
	}

	viewVariant struct {
		Attributes      map[string]string
		Marketplacecode string
		Title           string
		Url             string
	}
)

func (vc *View) variantSelection(configurable domain.ConfigurableProduct, activeVariant *domain.Variant) variantSelection {
	var variants variantSelection
	combinations := make(map[string]map[string]map[string]map[string]bool)
	titles := make(map[string]map[string]string)

	for _, attribute := range configurable.VariantVariationAttributes {
		// attribute -> value -> combinableAttribute -> combinaleValue -> true
		combinations[attribute] = make(map[string]map[string]map[string]bool)
		titles[attribute] = make(map[string]string)

		for _, variant := range configurable.Variants {
			titles[attribute][variant.Attributes[attribute].Value()] = variant.Attributes[attribute].Label

			for _, subattribute := range configurable.VariantVariationAttributes {
				if _, ok := variant.Attributes[attribute]; !ok {
					continue
				}
				if combinations[attribute][variant.Attributes[attribute].Value()] == nil {
					combinations[attribute][variant.Attributes[attribute].Value()] = make(map[string]map[string]bool)
				}

				//titles[variant.Attributes[subattribute].Value()] = variant.Attributes[subattribute].Label
				if subattribute != attribute {
					if _, ok := variant.Attributes[subattribute]; ok {
						if combinations[attribute][variant.Attributes[attribute].Value()][subattribute] == nil {
							combinations[attribute][variant.Attributes[attribute].Value()][subattribute] = make(map[string]bool)
						}
						combinations[attribute][variant.Attributes[attribute].Value()][subattribute][variant.Attributes[subattribute].Value()] = true
					}
				}
			}
		}
	}

	for code, attribute := range combinations {
		viewVariantAttribute := viewVariantAttribute{
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
			if activeVariant != nil && activeVariant.Attributes[code].Value() == optionCode {
				selected = true
			}

			label, ok := titles[code][optionCode]
			if !ok || label == "" {
				label = strings.Title(optionCode)
			}
			viewVariantAttribute.Options = append(viewVariantAttribute.Options, viewVariantOption{
				Key:          optionCode,
				Title:        label,
				Selected:     selected,
				Combinations: combinations,
			})
		}

		variants.Attributes = append(variants.Attributes, viewVariantAttribute)
	}

	for _, variant := range configurable.Variants {
		urlName := web.URLTitle(variant.BasicProductData.Title)
		variantUrl := vc.Router.URL("product.view", router.P{"marketplacecode": configurable.MarketPlaceCode, "variantcode": variant.MarketPlaceCode, "name": urlName}).String()

		attributes := make(map[string]string)

		for _, attr := range variants.Attributes {
			if !variant.HasAttribute(attr.Key) {
				continue
			}
			attributes[attr.Key] = variant.Attributes[attr.Key].Value()
		}

		variants.Variants = append(variants.Variants, viewVariant{
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

	var viewData productViewData

	// 1. Handle Configurables
	if product.Type() == "configurable" {
		configurableProduct := product.(domain.ConfigurableProduct)
		var activeVariant *domain.Variant

		viewData = productViewData{}
		variantCode, err := c.Param1("variantcode")

		if err != nil {
			//Redirect if url is not canonical
			redirect := vc.getRedirectIfRequired(configurableProduct, c.MustParam1("name"), skipnamecheck)
			if redirect != nil {
				return redirect
			}
			// 1.A. No variant selected
			viewData.VariantSelected = false
			viewData.RenderContext = "configurable"
			viewData.ConfigurableProduct = configurableProduct
		} else {
			// 1.B. Variant selected
			activeVariant, err = configurableProduct.Variant(variantCode)
			configurableProduct.ActiveVariant = activeVariant

			//Redirect if url is not canonical
			redirect := vc.getRedirectIfRequired(configurableProduct, c.MustParam1("name"), skipnamecheck)
			if redirect != nil {
				return redirect
			}
			viewData.VariantSelected = true
			viewData.RenderContext = "configurable_with_activevariant"
			viewData.ConfigurableProduct = configurableProduct
		}
		viewData.VariantSelection = vc.variantSelection(configurableProduct, activeVariant)

	} else {
		//Redirect if url is not canonical
		redirect := vc.getRedirectIfRequired(product, c.MustParam1("name"), skipnamecheck)
		if redirect != nil {
			return redirect
		}

		// 2. Handle Simples
		simpleProduct := product.(domain.SimpleProduct)
		viewData = productViewData{SimpleProduct: simpleProduct, RenderContext: "simple"}
	}

	vc.addBreadCrum(product, c)

	return vc.Render(c, vc.Template, viewData)
}

// addBreadCrum
func (vc *View) addBreadCrum(product domain.BasicProduct, c web.Context) {
	var paths []string
	if product.Type() == domain.TYPESIMPLE {
		paths = product.BaseData().CategoryToCodeMapping
	} else if configurableProduct, ok := product.(domain.ConfigurableProduct); ok {
		paths = configurableProduct.ConfigurableBaseData().CategoryToCodeMapping
	}

	//sort.Strings(paths)
	var stringHead string
	for _, p := range paths {
		parts := strings.Split(p, ":")
		name, code := parts[0], parts[1]
		if strings.HasPrefix(name, stringHead) {
			name = name[len(stringHead):]
			breadcrumbs.Add(c, breadcrumbs.Crumb{
				Title: name,
				Url:   vc.Router.URL(category.URLWithName(code, name)).String(),
			})
			stringHead = name + "/"
		}
	}
}

func (vc *View) getRedirectIfRequired(product domain.BasicProduct, currentNameParameter string, skipnamecheck string) web.Redirect {

	if skipnamecheck != "" {
		return nil
	}
	//Redirect if url is not canonical
	if vc.UrlService.GetNameParam(product, "") != currentNameParameter {
		if url, err := vc.UrlService.Get(product, ""); err == nil {
			return vc.RedirectPermanentURL(url)
		}
	}
	return nil
}
