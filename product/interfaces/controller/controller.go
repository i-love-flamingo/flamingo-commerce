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
	titles := make(map[string]string)

	for _, attribute := range configurable.VariantVariationAttributes {
		// attribute -> value -> combinableAttribute -> combinaleValue -> true
		combinations[attribute] = make(map[string]map[string]map[string]bool)

		for _, variant := range configurable.Variants {
			for _, subattribute := range configurable.VariantVariationAttributes {
				if _, ok := variant.Attributes[attribute]; !ok {
					continue
				}
				if combinations[attribute][variant.Attributes[attribute].Value()] == nil {
					combinations[attribute][variant.Attributes[attribute].Value()] = make(map[string]map[string]bool)
				}

				titles[variant.Attributes[attribute].Value()] = variant.Attributes[attribute].Label
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

			label, ok := titles[optionCode]
			if !ok {
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

		variantCode, err := c.Param1("variantcode")

		if err != nil {
			// 1.A. No variant selected
			// normalize URL
			urlName := web.URLTitle(configurableProduct.ConfigurableBaseData().Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.RedirectPermanent(URL(c.MustParam1("marketplacecode"), urlName))
			}
			viewData = productViewData{
				ConfigurableProduct: configurableProduct,
				VariantSelected:     false,
				RenderContext:       "configurable",
			}
		} else {
			activeVariant, _ = configurableProduct.Variant(variantCode)
			// 1.B. Variant selected
			// normalize URL
			urlName := web.URLTitle(activeVariant.BasicProductData.Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.RedirectPermanent(URLWithVariant(c.MustParam1("marketplacecode"), urlName, variantCode))
			}
			configurableProduct.ActiveVariant = activeVariant
			viewData = productViewData{
				ConfigurableProduct: configurableProduct,
				VariantSelected:     true,
				RenderContext:       "configurable_with_activevariant",
			}
		}

		viewData.VariantSelection = vc.variantSelection(configurableProduct, activeVariant)

	} else {
		// 2. Handle Simples
		// normalize URL
		urlName := web.URLTitle(product.BaseData().Title)
		if urlName != c.MustParam1("name") && skipnamecheck == "" {
			return vc.RedirectPermanent(URL(c.MustParam1("marketplacecode"), urlName))
		}

		simpleProduct := product.(domain.SimpleProduct)
		viewData = productViewData{SimpleProduct: simpleProduct, RenderContext: "simple"}
	}

	vc.addBreadCrum(product, c)

	return vc.Render(c, vc.Template, viewData)
}

// addBreadCrum
func (vc *View) addBreadCrum(product domain.BasicProduct, c web.Context) {
	paths := product.BaseData().CategoryToCodeMapping
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

// URL for a product
func URL(marketplacecode, name string) (string, map[string]string) {
	name = web.URLTitle(name)
	return "product.view", map[string]string{"marketplacecode": marketplacecode, "name": name}
}

// URLWithVariant for a product with a selected variant
func URLWithVariant(marketplacecode, name, variantcode string) (string, map[string]string) {
	return "product.view", map[string]string{"marketplacecode": marketplacecode, "name": name, "variantcode": variantcode}
}
