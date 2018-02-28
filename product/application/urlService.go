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
	if configurableProduct, ok := product.(domain.ConfigurableProduct); ok {
		if variantCode != "" && configurableProduct.HasVariant(variantCode) {
			return s.urlWithVariant(configurableProduct.ConfigurableBaseData().MarketPlaceCode, variantCode, configurableProduct.ConfigurableBaseData().Title)
		}
		if configurableProduct.HasActiveVariant() {
			return s.urlWithVariant(configurableProduct.ConfigurableBaseData().MarketPlaceCode, configurableProduct.ActiveVariant.MarketPlaceCode, configurableProduct.ActiveVariant.BaseData().Title)
		}
		if configurableProduct.TeaserData().PreSelectedVariantSku != "" {
			return s.urlWithVariant(configurableProduct.ConfigurableBaseData().MarketPlaceCode, configurableProduct.TeaserData().PreSelectedVariantSku, configurableProduct.TeaserData().ShortTitle)
		}
		return s.url(configurableProduct.ConfigurableBaseData().MarketPlaceCode, configurableProduct.ConfigurableBaseData().Title)
	}
	return s.url(product.BaseData().MarketPlaceCode, product.BaseData().Title)
}

// URL for a product
func (s *UrlService) url(marketplacecode, name string) (string, error) {
	name = web.URLTitle(name)
	url := s.Router.URL("product.view", map[string]string{"marketplacecode": marketplacecode, "name": name})
	return url.String(), nil
}

// URL for a product
func (s *UrlService) urlWithVariant(marketplacecode, variantcode, name string) (string, error) {
	name = web.URLTitle(name)
	url := s.Router.URL("product.view", map[string]string{"marketplacecode": marketplacecode, "name": name, "variantcode": variantcode})
	return url.String(), nil
}
