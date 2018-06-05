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
	"net/url"
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
		RenderContext    string
		Product          domain.BasicProduct
		VariantSelected  bool
		VariantSelection variantSelection
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

	combinationsOrder := make(map[string][]string)

	for _, attribute := range configurable.VariantVariationAttributes {
		// attribute -> value -> combinableAttribute -> combinableValue -> true
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
					combinationsOrder[attribute] = append(combinationsOrder[attribute], variant.Attributes[attribute].Value())
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

	for _, code := range configurable.VariantVariationAttributes {
		attribute := combinations[code]

		viewVariantAttribute := viewVariantAttribute{
			Key:   code,
			Title: strings.Title(code),
		}

		options := append([]string{}, configurable.VariantVariationAttributesSorting[viewVariantAttribute.Key]...)
		options = append(options, combinationsOrder[viewVariantAttribute.Key]...)
		knownOption := make(map[string]struct{}, len(options))

		for _, optionCode := range options {
			option, ok := attribute[optionCode]
			if _, known := knownOption[optionCode]; !ok || known {
				continue
			}
			knownOption[optionCode] = struct{}{}

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
	if product.Type() == domain.TYPECONFIGURABLE {
		configurableProduct := product.(domain.ConfigurableProduct)
		var activeVariant *domain.Variant

		viewData = productViewData{}
		variantCode, err := c.Param1("variantcode")

		if err != nil {
			//Redirect if url is not canonical
			redirect := vc.getRedirectIfRequired(configurableProduct, c, skipnamecheck)
			if redirect != nil {
				return redirect
			}
			// 1.A. No variant selected
			viewData.VariantSelected = false
			viewData.RenderContext = "configurable"
			viewData.Product = configurableProduct
		} else {

			configurableProductWithActiveVariant, err := configurableProduct.GetConfigurableWithActiveVariant(variantCode)
			if err != nil {
				return vc.ErrorNotFound(c, err)
			}
			activeVariant = &configurableProductWithActiveVariant.ActiveVariant
			//Redirect if url is not canonical
			redirect := vc.getRedirectIfRequired(configurableProductWithActiveVariant, c, skipnamecheck)
			if redirect != nil {
				return redirect
			}
			viewData.VariantSelected = true
			viewData.RenderContext = "configurable_with_activevariant"
			viewData.Product = configurableProductWithActiveVariant
		}
		viewData.VariantSelection = vc.variantSelection(configurableProduct, activeVariant)

	} else {
		//Redirect if url is not canonical
		redirect := vc.getRedirectIfRequired(product, c, skipnamecheck)
		if redirect != nil {
			return redirect
		}

		// 2. Handle Simples
		simpleProduct := product.(domain.SimpleProduct)
		viewData = productViewData{Product: simpleProduct, RenderContext: "simple"}
	}

	vc.addBreadCrumb(product, c)

	return vc.Render(c, vc.Template, viewData)
}

// addBreadCrumb
func (vc *View) addBreadCrumb(product domain.BasicProduct, c web.Context) {
	var paths []string
	if product.Type() == domain.TYPESIMPLE || product.Type() == domain.TYPECONFIGURABLE {
		paths = product.BaseData().CategoryToCodeMapping
	} else if configurableProduct, ok := product.(domain.ConfigurableProductWithActiveVariant); ok {
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
			stringHead = stringHead + name + "/"
		}
	}
}

func (vc *View) getRedirectIfRequired(product domain.BasicProduct, context web.Context, skipnamecheck string) web.Redirect {

	currentNameParameter := context.MustParam1("name")
	var allParams url.Values
	if context.QueryAll() != nil {
		allParams = url.Values(context.QueryAll())
	}

	if skipnamecheck != "" {
		return nil
	}
	//Redirect if url is not canonical
	if vc.UrlService.GetNameParam(product, "") != currentNameParameter {
		if redirectUrl, err := vc.UrlService.Get(product, ""); err == nil {
			newUrl, _ := url.Parse(redirectUrl)
			if len(allParams) > 0 {
				newUrl.RawQuery = allParams.Encode()
			}
			return vc.RedirectPermanentURL(newUrl.String())
		}
	}
	return nil
}
