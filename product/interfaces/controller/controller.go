package controller

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// View demonstrates a product view controller
	View struct {
		Responder             *web.Responder `inject:""`
		domain.ProductService `inject:""`
		URLService            *application.URLService `inject:""`

		Template string      `inject:"config:commerce.product.view.template"`
		Router   *web.Router `inject:""`
	}

	// productViewData is used for product rendering
	productViewData struct {
		// simple / configurable / configurable_with_variant
		RenderContext    string
		Product          domain.BasicProduct
		VariantSelected  bool
		VariantSelection variantSelection
		BackURL          string
	}

	// variantSelection for templating
	variantSelection struct {
		Attributes []viewVariantAttribute
		Variants   []viewVariant
	}

	viewVariantAttribute struct {
		Key       string
		Title     string
		CodeLabel string
		Options   []viewVariantOption
	}

	viewVariantOption struct {
		Key          string
		Title        string
		Combinations map[string][]string
		Selected     bool
		Unit         string
	}

	viewVariant struct {
		Attributes      map[string]string
		Marketplacecode string
		Title           string
		URL             string
		InStock         bool
	}
)

func (vc *View) variantSelection(configurable domain.ConfigurableProduct, activeVariant *domain.Variant) variantSelection {
	var variants variantSelection
	combinations := make(map[string]map[string]map[string]map[string]bool)
	titles := make(map[string]map[string]string)
	units := make(map[string]map[string]string)

	combinationsOrder := make(map[string][]string)

	for _, attribute := range configurable.VariantVariationAttributes {
		// attribute -> value -> combinableAttribute -> combinableValue -> true
		combinations[attribute] = make(map[string]map[string]map[string]bool)
		titles[attribute] = make(map[string]string)
		units[attribute] = make(map[string]string)

		for _, variant := range configurable.Variants {
			titles[attribute][variant.Attributes[attribute].Value()] = variant.Attributes[attribute].Label
			units[attribute][variant.Attributes[attribute].Value()] = variant.Attributes[attribute].UnitCode

			for _, subattribute := range configurable.VariantVariationAttributes {
				if _, ok := variant.Attributes[attribute]; !ok {
					continue
				}
				if combinations[attribute][variant.Attributes[attribute].Value()] == nil {
					combinations[attribute][variant.Attributes[attribute].Value()] = make(map[string]map[string]bool)
					combinationsOrder[attribute] = append(combinationsOrder[attribute], variant.Attributes[attribute].Value())
				}

				// titles[variant.Attributes[subattribute].Value()] = variant.Attributes[subattribute].Label
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

		codeLabel := ""
		if len(configurable.Variants) > 0 {
			codeLabel = configurable.Variants[0].BaseData().Attributes.Attribute(code).CodeLabel
		}

		viewVariantAttribute := viewVariantAttribute{
			Key:       code,
			Title:     cases.Title(language.Und).String(code),
			CodeLabel: codeLabel,
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
				label = cases.Title(language.Und).String(optionCode)
			}
			viewVariantAttribute.Options = append(viewVariantAttribute.Options, viewVariantOption{
				Key:          optionCode,
				Title:        label,
				Selected:     selected,
				Combinations: combinations,
				Unit:         units[code][optionCode],
			})
		}

		variants.Attributes = append(variants.Attributes, viewVariantAttribute)
	}

	for _, variant := range configurable.Variants {
		urlName := web.URLTitle(variant.BasicProductData.Title)
		variantURL, _ := vc.Router.URL("product.view", map[string]string{"marketplacecode": configurable.MarketPlaceCode, "variantcode": variant.MarketPlaceCode, "name": urlName})

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
			URL:             variantURL.String(),
			Attributes:      attributes,
			InStock:         variant.IsInStock(),
		})
	}

	return variants
}

// Get Response for Product matching sku param
func (vc *View) Get(c context.Context, r *web.Request) web.Result {
	product, err := vc.ProductService.Get(c, r.Params["marketplacecode"])
	skipnamecheck := r.Params["skipnamecheck"]

	// catch error
	if err != nil {
		switch errors.Cause(err).(type) {
		case domain.ProductNotFound:
			return vc.Responder.NotFound(err)

		default:
			return vc.Responder.ServerError(err)
		}
	}

	var viewData productViewData

	// 1. Handle Configurables
	if product.Type() == domain.TypeConfigurable {
		configurableProduct := product.(domain.ConfigurableProduct)
		var activeVariant *domain.Variant

		viewData = productViewData{}
		variantCode, ok := r.Params["variantcode"]

		if !ok {
			// Redirect if url is not canonical
			redirect := vc.getRedirectIfRequired(configurableProduct, r, skipnamecheck)
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
				return vc.Responder.NotFound(err)
			}
			activeVariant = &configurableProductWithActiveVariant.ActiveVariant
			// Redirect if url is not canonical
			redirect := vc.getRedirectIfRequired(configurableProductWithActiveVariant, r, skipnamecheck)
			if redirect != nil {
				return redirect
			}
			viewData.VariantSelected = true
			viewData.RenderContext = "configurable_with_activevariant"
			viewData.Product = configurableProductWithActiveVariant
		}
		viewData.VariantSelection = vc.variantSelection(configurableProduct, activeVariant)

	} else {
		// Redirect if url is not canonical
		redirect := vc.getRedirectIfRequired(product, r, skipnamecheck)
		if redirect != nil {
			return redirect
		}

		// 2. Handle Simples
		simpleProduct := product.(domain.SimpleProduct)
		viewData = productViewData{Product: simpleProduct, RenderContext: "simple"}
	}

	backURL, err := r.Query1("backurl")
	if err == nil {
		viewData.BackURL = backURL
	}

	return vc.Responder.Render(vc.Template, viewData)
}

func (vc *View) getRedirectIfRequired(product domain.BasicProduct, r *web.Request, skipnamecheck string) *web.URLRedirectResponse {
	currentNameParameter := r.Params["name"]
	var allParams url.Values
	if r.QueryAll() != nil {
		allParams = url.Values(r.QueryAll())
	}

	if skipnamecheck != "" {
		return nil
	}
	// Redirect if url is not canonical
	if vc.URLService.GetNameParam(product, "") != currentNameParameter {
		if redirectURL, err := vc.URLService.Get(product, ""); err == nil {
			newURL, _ := url.Parse(redirectURL)
			if len(allParams) > 0 {
				newURL.RawQuery = allParams.Encode()
			}
			return vc.Responder.URLRedirect(newURL).Permanent()
		}
	}
	return nil
}
