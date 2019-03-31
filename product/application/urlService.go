package application

import (
	"errors"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

// URLService to manage product urls
type URLService struct {
	Router *web.Router `inject:""`
}

// Get a product variant url
func (s *URLService) Get(product domain.BasicProduct, variantCode string) (string, error) {
	if product == nil {
		return "-", errors.New("no product given")
	}
	params := s.GetURLParams(product, variantCode)
	url, err := s.Router.Relative("product.view", params)
	return url.String(), err
}

// Get a product variant url
func (s *URLService) GetAbsolute(r *web.Request, product domain.BasicProduct, variantCode string) (string, error) {
	if product == nil {
		return "-", errors.New("no product given")
	}
	params := s.GetURLParams(product, variantCode)
	url, err := s.Router.Absolute(r, "product.view", params)
	return url.String(), err
}

// GetURLParams get product url params
func (s *URLService) GetURLParams(product domain.BasicProduct, variantCode string) map[string]string {
	params := make(map[string]string)
	if product == nil {
		return params
	}

	if product.Type() == domain.TypeSimple {
		params["marketplacecode"] = product.BaseData().MarketPlaceCode
		params["name"] = web.URLTitle(product.BaseData().Title)
	}
	if product.Type() == domain.TypeConfigurableWithActiveVariant {
		if configurableProduct, ok := product.(domain.ConfigurableProductWithActiveVariant); ok {
			params["marketplacecode"] = configurableProduct.ConfigurableBaseData().MarketPlaceCode
			params["name"] = web.URLTitle(configurableProduct.ConfigurableBaseData().Title)
			if variantCode != "" && configurableProduct.HasVariant(variantCode) {
				variantInstance, err := configurableProduct.Variant(variantCode)
				if err == nil {
					params["variantcode"] = variantCode
					params["name"] = web.URLTitle(variantInstance.BaseData().Title)
				}
			} else {
				params["variantcode"] = configurableProduct.ActiveVariant.MarketPlaceCode
				params["name"] = web.URLTitle(configurableProduct.ActiveVariant.BaseData().Title)
			}
		}
	}

	if product.Type() == domain.TypeConfigurable {
		if configurableProduct, ok := product.(domain.ConfigurableProduct); ok {
			params["marketplacecode"] = configurableProduct.BaseData().MarketPlaceCode
			params["name"] = web.URLTitle(configurableProduct.BaseData().Title)
			//if the teaser teasers a variant then link to this
			if configurableProduct.TeaserData().PreSelectedVariantSku != "" {
				params["variantcode"] = configurableProduct.TeaserData().PreSelectedVariantSku
				params["name"] = web.URLTitle(configurableProduct.TeaserData().ShortTitle)
			}
			//if a variantCode is given then link to that variant
			if variantCode != "" && configurableProduct.HasVariant(variantCode) {
				variantInstance, err := configurableProduct.Variant(variantCode)
				if err == nil {
					params["variantcode"] = variantCode
					params["name"] = web.URLTitle(variantInstance.BaseData().Title)
				}
			}
		}
	}

	return params
}

// GetNameParam retrieve the proper name parameter
func (s *URLService) GetNameParam(product domain.BasicProduct, variantCode string) string {
	params := s.GetURLParams(product, variantCode)
	if name, ok := params["name"]; ok {
		return name
	}
	return ""
}
