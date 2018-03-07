package application

import (
	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	UrlService struct {
		Router *router.Router `inject:""`
	}
)

func (s *UrlService) Get(product domain.BasicProduct, variantCode string) (string, error) {
	if product == nil {
		return "-", errors.New("no product given")
	}
	params := s.GetUrlParams(product, variantCode)
	url := s.Router.URL("product.view", params)
	return url.String(), nil
}

func (s *UrlService) GetUrlParams(product domain.BasicProduct, variantCode string) map[string]string {
	params := make(map[string]string)
	if product == nil {
		return params
	}
	if configurableProduct, ok := product.(domain.ConfigurableProduct); ok {
		params["marketplacecode"] = configurableProduct.ConfigurableBaseData().MarketPlaceCode
		if variantCode != "" && configurableProduct.HasVariant(variantCode) {
			params["variantcode"] = variantCode
			params["name"] = web.URLTitle(configurableProduct.ConfigurableBaseData().Title)
		}
		if configurableProduct.HasActiveVariant() {
			params["variantcode"] = configurableProduct.ActiveVariant.MarketPlaceCode
			params["name"] = web.URLTitle(configurableProduct.ActiveVariant.BaseData().Title)
		}
		if configurableProduct.TeaserData().PreSelectedVariantSku != "" {
			params["variantcode"] = configurableProduct.TeaserData().PreSelectedVariantSku
			params["name"] = web.URLTitle(configurableProduct.TeaserData().ShortTitle)
		}
		params["name"] = web.URLTitle(configurableProduct.ConfigurableBaseData().Title)
	} else {
		params["marketplacecode"] = product.BaseData().MarketPlaceCode
		params["name"] = web.URLTitle(product.BaseData().Title)
	}
	return params
}
