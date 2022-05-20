package application

import (
	"errors"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// URLService to manage product urls
	URLService struct {
		router            *web.Router
		generateSlug      bool
		slugAttributecode string
	}
)

// Inject dependencies
func (s *URLService) Inject(
	r *web.Router,
	c *struct {
		GenerateSlug      bool   `inject:"config:commerce.product.generateSlug,optional"`
		SlugAttributecode string `inject:"config:commerce.product.slugAttributeCode,optional"`
	},
) *URLService {
	s.router = r

	if c != nil {
		s.generateSlug = c.GenerateSlug
		s.slugAttributecode = c.SlugAttributecode
	}

	return s
}

// Get a product variant url
func (s *URLService) Get(product domain.BasicProduct, variantCode string) (string, error) {
	if product == nil {
		return "-", errors.New("no product given")
	}
	params := s.GetURLParams(product, variantCode)
	url, err := s.router.Relative("product.view", params)
	return url.String(), err
}

// GetAbsolute url for a product variant url
func (s *URLService) GetAbsolute(r *web.Request, product domain.BasicProduct, variantCode string) (string, error) {
	if product == nil {
		return "-", errors.New("no product given")
	}
	params := s.GetURLParams(product, variantCode)
	url, err := s.router.Absolute(r, "product.view", params)
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
		params["name"] = s.getSlug(product.BaseData(), product.BaseData().Title)
	}
	if product.Type() == domain.TypeConfigurableWithActiveVariant {
		if configurableProduct, ok := product.(domain.ConfigurableProductWithActiveVariant); ok {
			params["marketplacecode"] = configurableProduct.ConfigurableBaseData().MarketPlaceCode
			params["name"] = s.getSlug(configurableProduct.ConfigurableBaseData(), configurableProduct.ConfigurableBaseData().Title)
			if variantCode != "" && configurableProduct.HasVariant(variantCode) {
				variantInstance, err := configurableProduct.Variant(variantCode)
				if err == nil {
					params["variantcode"] = variantCode
					params["name"] = s.getSlug(variantInstance.BaseData(), variantInstance.BaseData().Title)
				}
			} else {
				params["variantcode"] = configurableProduct.ActiveVariant.MarketPlaceCode
				params["name"] = s.getSlug(configurableProduct.ActiveVariant.BaseData(), configurableProduct.ActiveVariant.BaseData().Title)
			}
		}
	}

	if product.Type() == domain.TypeConfigurable {
		if configurableProduct, ok := product.(domain.ConfigurableProduct); ok {
			params["marketplacecode"] = configurableProduct.BaseData().MarketPlaceCode
			params["name"] = s.getSlug(configurableProduct.BaseData(), configurableProduct.BaseData().Title)
			// if the teaser teasers a variant then link to this
			if configurableProduct.TeaserData().PreSelectedVariantSku != "" {
				params["variantcode"] = configurableProduct.TeaserData().PreSelectedVariantSku
				params["name"] = func(d domain.TeaserData) string {
					if s.generateSlug {
						return web.URLTitle(d.ShortTitle)
					}

					if d.URLSlug == "" {
						return web.URLTitle(d.ShortTitle)
					}

					return d.URLSlug
				}(configurableProduct.TeaserData())
			}
			// if a variantCode is given then link to that variant
			if variantCode != "" && configurableProduct.HasVariant(variantCode) {
				variantInstance, err := configurableProduct.Variant(variantCode)
				if err == nil {
					params["variantcode"] = variantCode
					params["name"] = s.getSlug(variantInstance.BaseData(), variantInstance.BaseData().Title)
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

// getSlug fetches the slug from the BasicProductData if available, returns web.URLTitle encoded fallback if disabled, attribute does not exist or attribute is empty
func (s *URLService) getSlug(b domain.BasicProductData, fallback string) string {
	if s.generateSlug {
		return web.URLTitle(fallback)
	}

	if !b.HasAttribute(s.slugAttributecode) {
		return web.URLTitle(fallback)
	}

	if nil == b.Attributes[s.slugAttributecode].RawValue {
		return web.URLTitle(fallback)
	}

	return b.Attributes[s.slugAttributecode].Value()
}
