package interfaces

import (
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"net/url"

	"github.com/pkg/errors"

	"strings"
)

type (
	// ViewController demonstrates a product view controller
	ViewController struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		domain.ProductService   `inject:""`
	}

	// SimpleProductViewData is used for product rendering
	SimpleProductViewData struct {
		SimpleProduct domain.SimpleProduct
	}

	// ConfigurableProductViewData is used for product rendering
	ConfigurableProductViewData struct {
		ConfigurableProduct domain.ConfigurableProduct
		ActiveVariant       domain.Variant
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	product, err := vc.ProductService.Get(c, c.MustParam1("marketplacecode"))

	// catch error
	if err != nil {
		switch errors.Cause(err).(type) {
		case domain.ProductNotFound:

			return vc.ErrorNotFound(c, err)

		default:

			return vc.Error(c, err)
		}
	}

	// 1. Handle Configurables
	if product.GetType() == "configurable" {
		configurableProduct := product.(domain.ConfigurableProduct)
		variantCode, err := c.Param1("variantcode")
		var activeVariant domain.Variant
		if err == nil {
			activeVariant, _ = configurableProduct.GetVariant(variantCode)
		}
		if &activeVariant == nil {
			// 1.A. No variant selected
			// normalize URL
			urlName := makeUrlTitle(product.GetBaseData().Title)
			if urlName != c.MustParam1("name") {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "name": urlName})
			}
			return vc.Render(c, "product/configurable", ConfigurableProductViewData{ConfigurableProduct: configurableProduct})
		} else {
			// 1.B. Variant selected
			// normalize URL
			urlName := makeUrlTitle(activeVariant.BasicProductData.Title)
			if urlName != c.MustParam1("name") {
				return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "variantcode": c.MustParam1("variantcode"), "name": urlName})
			}
			return vc.Render(c, "product/configurable", ConfigurableProductViewData{ConfigurableProduct: configurableProduct, ActiveVariant: activeVariant})
		}

	} else {
		// 2. Handle Simples
		// normalize URL
		urlName := makeUrlTitle(product.GetBaseData().Title)
		if urlName != c.MustParam1("name") {
			return vc.Redirect("product.view", router.P{"marketplacecode": c.MustParam1("marketplacecode"), "name": urlName})
		}

		simpleProduct := product.(domain.SimpleProduct)
		return vc.Render(c, "product/simple", SimpleProductViewData{SimpleProduct: simpleProduct})
	}

}

func makeUrlTitle(title string) string {
	newTitle := strings.ToLower(strings.Replace(title, " ", "_", -1))
	newTitle = url.QueryEscape(newTitle)

	return newTitle
}
