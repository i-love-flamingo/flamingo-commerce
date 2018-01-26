package application

import (
	"crypto/sha512"
	"encoding/base64"
	"strings"

	authApplication "go.aoe.com/flamingo/core/auth/application"
	canonicalUrlApplication "go.aoe.com/flamingo/core/canonicalUrl/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/core/w3cDatalayer/domain"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	/**
	Factory is used to build new datalayers
	*/
	Factory struct {
		Router              *router.Router                  `inject:""`
		DatalayerProvider   domain.DatalayerProvider        `inject:""`
		CanonicalUrlService canonicalUrlApplication.Service `inject:""`
		UserService         authApplication.UserService     `inject:""`

		PageNamePrefix  string `inject:"config:w3cDatalayer.pageNamePrefix,optional"`
		SiteName        string `inject:"config:w3cDatalayer.siteName,optional"`
		Locale          string `inject:"config:locale.locale,optional"`
		DefaultCurrency string `inject:"config:w3cDatalayer.defaultCurrency,optional"`
		HashUserValues  bool   `inject:"config:w3cDatalayer.hashUserValues,optional"`
	}
)

//Update
func (s Factory) BuildForCurrentRequest(ctx web.Context) domain.Datalayer {

	layer := s.DatalayerProvider()

	//get langiage from locale code configuration
	language := ""
	localeParts := strings.Split(s.Locale, "-")
	if len(localeParts) > 0 {
		language = localeParts[0]
	}

	layer.Page = &domain.Page{
		PageInfo: domain.PageInfo{
			PageID:         ctx.Request().URL.Path,
			PageName:       s.PageNamePrefix + ctx.Request().URL.Path,
			DestinationURL: s.CanonicalUrlService.GetCanonicalUrlForCurrentRequest(ctx),
			Language:       language,
		},
		Attributes: make(map[string]interface{}),
	}

	layer.Page.Attributes["currency"] = s.DefaultCurrency

	//Use the handler name as PageId if available
	if controllerHandler, ok := ctx.Value("HandlerName").(string); ok {
		layer.Page.PageInfo.PageID = controllerHandler
	}

	layer.SiteInfo = &domain.SiteInfo{
		SiteName: s.SiteName,
	}

	//Handle User
	layer.Page.Attributes["loggedIn"] = false
	if s.UserService.IsLoggedIn(ctx) {
		layer.Page.Attributes["loggedIn"] = true
		layer.Page.Attributes["logintype"] = "external"
		userData := s.getUser(ctx)
		if userData != nil {
			layer.User = append(layer.User, *userData)
		}
	}
	return *layer
}

func (s Factory) getUser(ctx web.Context) *domain.User {
	user := s.UserService.GetUser(ctx)
	if user == nil {
		return nil
	}

	dataLayerProfile := domain.UserProfile{
		ProfileInfo: domain.UserProfileInfo{
			EmailID:   user.Email,
			ProfileID: user.Sub,
		},
	}

	if s.HashUserValues {
		dataLayerProfile.ProfileInfo.EmailID = hashWithSHA512(dataLayerProfile.ProfileInfo.EmailID)
		dataLayerProfile.ProfileInfo.ProfileID = hashWithSHA512(dataLayerProfile.ProfileInfo.ProfileID)
	}

	dataLayerUser := domain.User{}
	dataLayerUser.Profile = append(dataLayerUser.Profile, dataLayerProfile)
	return &dataLayerUser
}

func (s Factory) BuildCartData(cart cart.DecoratedCart) *domain.Cart {
	cartData := domain.Cart{
		CartID: cart.Cart.ID,
		Price: &domain.CartPrice{
			Currency:       cart.Cart.CurrencyCode,
			BasePrice:      cart.Cart.SubTotal,
			CartTotal:      cart.Cart.GrandTotal,
			Shipping:       cart.Cart.ShippingItem.Price,
			ShippingMethod: cart.Cart.ShippingItem.Title,
		},
		Attributes: make(map[string]interface{}),
	}
	for _, item := range cart.DecoratedItems {
		level0 := ""
		level1 := ""
		if len(item.Product.BaseData().CategoryPath) > 0 {
			level0 = item.Product.BaseData().CategoryPath[0]
		}
		if len(item.Product.BaseData().CategoryPath) > 1 {
			level0 = item.Product.BaseData().CategoryPath[1]
		}
		productFamily := ""
		if item.Product.BaseData().HasAttribute("gs1Family") {
			productFamily = item.Product.BaseData().Attributes["gs1Family"].Value()
		}

		//Handle Variants if it is a Configurable
		var parentIdRef *string = nil
		var variantSelectedAttributeRef *string = nil
		if item.Product.Type() == productDomain.TYPECONFIGURABLE {
			if configurable, ok := item.Product.(productDomain.ConfigurableProduct); ok {
				parentId := configurable.BaseData().MarketPlaceCode
				parentIdRef = &parentId
				if configurable.HasActiveVariant() && len(configurable.VariantVariationAttributes) > 0 && configurable.ActiveVariant.HasAttribute(configurable.VariantVariationAttributes[0]) {
					variantSelectedAttribute := configurable.ActiveVariant.BaseData().Attributes[configurable.VariantVariationAttributes[0]].Value()
					variantSelectedAttributeRef = &variantSelectedAttribute
				}
			}
		}

		itemData := domain.CartItem{
			Category: &domain.CartItemCategory{
				PrimaryCategory: level0,
				SubCategory1:    level1,
				ProductType:     productFamily,
			},
			Quantity: item.Item.Qty,
			ProductInfo: domain.ProductInfo{
				ProductID:                item.GetDisplayMarketplaceCode(),
				ProductName:              item.GetDisplayTitle(),
				ProductThumbnail:         s.getItemImageUrl(item.Product),
				ProductType:              item.Product.Type(),
				ParentId:                 parentIdRef,
				VariantSelectedAttribute: variantSelectedAttributeRef,
			},
			Price: domain.CartItemPrice{
				BasePrice:    item.Item.Price,
				PriceWithTax: item.Item.PriceInclTax,
				TaxRate:      item.Item.TaxAmount,
				Currency:     cart.Cart.CurrencyCode,
			},
		}
		cartData.Item = append(cartData.Item, itemData)
	}
	return &cartData
}

func (s Factory) getItemImageUrl(product productDomain.BasicProduct) string {
	return "catalog/" + product.BaseData().GetListMedia().Reference
}

func hashWithSHA512(value string) string {
	newHash := sha512.New()
	newHash.Write([]byte(value))
	//the hash is a byte array
	result := newHash.Sum(nil)
	//since we want to uuse it in a variable we base64 encode it (other alternative would be Hexadecimal representation "% x", h.Sum(nil)
	return base64.URLEncoding.EncodeToString(result)
}
