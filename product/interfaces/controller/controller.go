package controller

import (
	"flamingo/core/breadcrumbs"
	"flamingo/core/category"
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"log"
	"sort"
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

	// productViewData is used for product rendering
	productViewData struct {
		// simple / configurable / configurable_with_variant
		RenderContext       string
		SimpleProduct       domain.SimpleProduct
		ConfigurableProduct domain.ConfigurableProduct
		ActiveVariant       domain.Variant
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
			if activeVariant != nil && activeVariant.Attributes[code] == optionCode {
				selected = true
			}
			viewVariantAttribute.Options = append(viewVariantAttribute.Options, viewVariantOption{
				Key:          optionCode,
				Title:        strings.Title(optionCode),
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
			if variant.Attributes[attr.Key] == nil {
				continue
			}
			attributes[attr.Key] = variant.Attributes[attr.Key].(string)
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
			urlName := web.URLTitle(product.BaseData().Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.Redirect(URL(c.MustParam1("marketplacecode"), urlName))
			}
			viewData = productViewData{
				ConfigurableProduct: configurableProduct,
				VariantSelected:     false,
				RenderContext:       "configurable",
			}
		} else {
			log.Println("get variant by " + variantCode)
			activeVariant, _ = configurableProduct.Variant(variantCode)
			// 1.B. Variant selected
			// normalize URL
			urlName := web.URLTitle(activeVariant.BasicProductData.Title)
			if urlName != c.MustParam1("name") && skipnamecheck == "" {
				return vc.Redirect(URLWithVariant(c.MustParam1("marketplacecode"), urlName, variantCode))
			}
			log.Printf("Variant Price %v / %v", activeVariant.ActivePrice.Default, activeVariant.ActivePrice)
			viewData = productViewData{
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
		urlName := web.URLTitle(product.BaseData().Title)
		if urlName != c.MustParam1("name") && skipnamecheck == "" {
			return vc.Redirect(URL(c.MustParam1("marketplacecode"), urlName))
		}

		simpleProduct := product.(domain.SimpleProduct)
		viewData = productViewData{SimpleProduct: simpleProduct, RenderContext: "simple"}
	}

	paths := product.BaseData().CategoryPath
	sort.Strings(paths)
	var stringHead string
	for _, p := range paths {
		if strings.HasPrefix(p, stringHead) {
			breadcrumbs.Add(c, breadcrumbs.Crumb{
				Title: p[len(stringHead):],
				URL:   vc.Router.URL(category.URLWithName(p[len(stringHead):], p[len(stringHead):])).String(),
			})
			stringHead = p + "/"
		}
	}

	return vc.Render(c, vc.Template, viewData)
}

// URL for a product
func URL(marketplacecode, name string) (string, map[string]string) {
	return "product.view", map[string]string{"marketplacecode": marketplacecode, "name": name}
}

// URLWithVariant for a product with a selected variant
func URLWithVariant(marketplacecode, name, variantcode string) (string, map[string]string) {
	return "product.view", map[string]string{"marketplacecode": marketplacecode, "name": name, "variantcode": variantcode}
}
