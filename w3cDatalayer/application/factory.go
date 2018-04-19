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
		Router              *router.Router                   `inject:""`
		DatalayerProvider   domain.DatalayerProvider         `inject:""`
		CanonicalUrlService *canonicalUrlApplication.Service `inject:""`
		UserService         *authApplication.UserService     `inject:""`

		PageInstanceIDPrefix           string `inject:"config:w3cDatalayer.pageInstanceIDPrefix,optional"`
		PageInstanceIDStage            string `inject:"config:w3cDatalayer.pageInstanceIDStage,optional"`
		ProductMediaBaseUrl            string `inject:"config:w3cDatalayer.productMediaBaseUrl,optional"`
		ProductMediaUrlPrefix          string `inject:"config:w3cDatalayer.productMediaUrlPrefix,optional"`
		ProductMediaThumbnailUrlPrefix string `inject:"config:w3cDatalayer.productMediaThumbnailUrlPrefix,optional"`
		PageNamePrefix                 string `inject:"config:w3cDatalayer.pageNamePrefix,optional"`
		SiteName                       string `inject:"config:w3cDatalayer.siteName,optional"`
		Locale                         string `inject:"config:locale.locale,optional"`
		DefaultCurrency                string `inject:"config:w3cDatalayer.defaultCurrency,optional"`
		HashUserValues                 bool   `inject:"config:w3cDatalayer.hashUserValues,optional"`
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

	layer.PageInstanceID = s.PageInstanceIDPrefix + s.PageInstanceIDStage

	//Handle User
	layer.Page.Attributes["loggedIn"] = false
	if s.UserService.IsLoggedIn(ctx) {
		layer.Page.Attributes["loggedIn"] = true
		layer.Page.Attributes["logintype"] = "external"
		userData := s.getUser(ctx)
		if userData != nil {
			layer.User = append(layer.User, *userData)
		}
	} else {
		layer.Page.Attributes["logintype"] = "guest"
	}
	return *layer
}

func (s Factory) getUser(ctx web.Context) *domain.User {

	dataLayerProfile := s.getUserProfileForCurrentUser(ctx)
	if dataLayerProfile == nil {
		return nil
	}

	dataLayerUser := domain.User{}
	dataLayerUser.Profile = append(dataLayerUser.Profile, *dataLayerProfile)
	return &dataLayerUser
}

func (s Factory) getUserProfileForCurrentUser(ctx web.Context) *domain.UserProfile {
	user := s.UserService.GetUser(ctx)
	if user == nil {
		return nil
	}
	return s.getUserProfile(user.Email, user.Sub)
}

func (s Factory) getUserProfile(email string, sub string) *domain.UserProfile {
	dataLayerProfile := domain.UserProfile{
		ProfileInfo: domain.UserProfileInfo{
			EmailID:   s.HashValueIfConfigured(email),
			ProfileID: s.HashValueIfConfigured(sub),
		},
	}
	return &dataLayerProfile
}

func (s Factory) HashValueIfConfigured(value string) string {
	if s.HashUserValues && value != "" {
		return hashWithSHA512(value)
	}
	return value
}

func (s Factory) BuildCartData(cart cart.DecoratedCart) *domain.Cart {
	cartData := domain.Cart{
		CartID: cart.Cart.ID,
		Price: &domain.CartPrice{
			Currency:       cart.Cart.CartTotals.CurrencyCode,
			BasePrice:      cart.Cart.CartTotals.SubTotal,
			CartTotal:      cart.Cart.CartTotals.GrandTotal,
			Shipping:       cart.Cart.CartTotals.TotalShippingItem.Price,
			ShippingMethod: cart.Cart.CartTotals.TotalShippingItem.Title,
			PriceWithTax:   cart.Cart.CartTotals.GrandTotal,
		},
		Attributes: make(map[string]interface{}),
	}
	for _, item := range cart.DecoratedItems {
		itemData := s.buildCartItem(item, cart.Cart.CartTotals.CurrencyCode)
		cartData.Item = append(cartData.Item, itemData)
	}
	return &cartData
}

func (s Factory) BuildTransactionData(ctx web.Context, cartTotals cart.CartTotals, decoratedItems []cart.DecoratedCartItem, orderId string, email string) *domain.Transaction {
	var profile *domain.UserProfile
	if s.UserService.IsLoggedIn(ctx) {
		profile = s.getUserProfileForCurrentUser(ctx)
	} else {
		profile = s.getUserProfile(email, "")
	}

	transactionData := domain.Transaction{
		TransactionID: orderId,
		Price: &domain.TransactionPrice{
			Currency:         cartTotals.CurrencyCode,
			BasePrice:        cartTotals.SubTotal,
			TransactionTotal: cartTotals.GrandTotal,
			Shipping:         cartTotals.TotalShippingItem.Price,
			ShippingMethod:   cartTotals.TotalShippingItem.Title,
		},
		Profile:    profile,
		Attributes: make(map[string]interface{}),
	}
	for _, item := range decoratedItems {
		itemData := s.buildCartItem(item, cartTotals.CurrencyCode)
		transactionData.Item = append(transactionData.Item, itemData)
	}
	return &transactionData
}

func (s Factory) buildCartItem(item cart.DecoratedCartItem, currencyCode string) domain.CartItem {
	cartItem := domain.CartItem{
		Category:    s.getProductCategory(item.Product),
		Quantity:    item.Item.Qty,
		ProductInfo: s.getProductInfo(item.Product),
		Price: domain.CartItemPrice{
			BasePrice:    item.Item.Price,
			PriceWithTax: item.Item.PriceInclTax,
			TaxRate:      item.Item.TaxAmount,
			Currency:     currencyCode,
		},
		Attributes: make(map[string]interface{}),
	}
	cartItem.Attributes["sourceId"] = item.Item.SourceId
	cartItem.Attributes["terminal"] = ""
	cartItem.Attributes["leadtime"] = ""
	return cartItem
}

func (s Factory) BuildProductData(product productDomain.BasicProduct) domain.Product {
	productData := domain.Product{
		ProductInfo: s.getProductInfo(product),
		Category:    s.getProductCategory(product),
		Attributes:  make(map[string]interface{}),
	}
	productData.Attributes["productPrice"] = product.SaleableData().ActivePrice.GetFinalPrice()
	if product.BaseData().HasAttribute("ispuLimitedToAreas") {
		productData.Attributes["ispuLimitedToAreas"] = product.BaseData().Attributes["ispuLimitedToAreas"].Value()
	}
	return productData
}

func (s Factory) getProductCategory(product productDomain.BasicProduct) *domain.ProductCategory {
	level0 := ""
	level1 := ""
	level2 := ""
	previousPathParts := 0
	longestFirstPath := ""

	categoryPaths := product.BaseData().CategoryPath
	baseData := product.BaseData()

	for _, path := range categoryPaths {
		pathParts := strings.Count(path, "/") + 1
		if pathParts > previousPathParts {
			previousPathParts = pathParts
		} else {
			//last path was longest first match
			longestFirstPath = path
		}
	}
	firstPathLevels := strings.Split(longestFirstPath, "/")
	if len(firstPathLevels) > 0 {
		level0 = firstPathLevels[0]
	}
	if len(firstPathLevels) > 1 {
		level1 = firstPathLevels[1]
	}
	if len(firstPathLevels) > 2 {
		level2 = firstPathLevels[2]
	}

	productFamily := ""
	if baseData.HasAttribute("gs1Family") {
		productFamily = baseData.Attributes["gs1Family"].Value()
	}
	return &domain.ProductCategory{
		PrimaryCategory: level0,
		SubCategory1:    level1,
		SubCategory2:    level2,
		SubCategory:     level1,
		ProductType:     productFamily,
	}
}

func (s Factory) getProductInfo(product productDomain.BasicProduct) domain.ProductInfo {
	baseData := product.BaseData()
	//Handle Variants if it is a Configurable
	var parentIdRef *string = nil
	var variantSelectedAttributeRef *string = nil
	if product.Type() == productDomain.TYPECONFIGURABLE_WITH_ACTIVE_VARIANT {
		if configurableWithActiveVariant, ok := product.(productDomain.ConfigurableProductWithActiveVariant); ok {
			parentId := configurableWithActiveVariant.ConfigurableBaseData().MarketPlaceCode
			parentIdRef = &parentId
			variantSelectedAttribute := configurableWithActiveVariant.ActiveVariant.BaseData().Attributes[configurableWithActiveVariant.VariantVariationAttributes[0]].Value()
			variantSelectedAttributeRef = &variantSelectedAttribute
		}
	}
	if product.Type() == productDomain.TYPECONFIGURABLE {
		if configurable, ok := product.(productDomain.ConfigurableProduct); ok {
			parentId := configurable.BaseData().MarketPlaceCode
			parentIdRef = &parentId
		}
	}
	// Search for some common product attributes to fill the productInfos (This maybe better to be configurable later)
	color := ""
	if baseData.HasAttribute("manufacturerColor") {
		color = baseData.Attributes["manufacturerColor"].Value()
	}
	if baseData.HasAttribute("baseColor") {
		color = baseData.Attributes["baseColor"].Value()
	}
	size := ""
	if baseData.HasAttribute("shoeSize") {
		size = baseData.Attributes["shoeSize"].Value()
	}
	if baseData.HasAttribute("clothingSize") {
		size = baseData.Attributes["clothingSize"].Value()
	}
	brand := ""
	if baseData.HasAttribute("brandCode") {
		brand = baseData.Attributes["brandCode"].Value()
	}
	gtin := ""
	if baseData.HasAttribute("gtin") {
		if baseData.Attributes["gtin"].HasMultipleValues() {
			gtins := baseData.Attributes["gtin"].Values()
			gtin = strings.Join(gtins, ",")
		} else {
			gtin = baseData.Attributes["gtin"].Value()
		}
	}
	return domain.ProductInfo{
		ProductID:                baseData.MarketPlaceCode,
		ProductName:              baseData.Title,
		ProductThumbnail:         s.getProductThumbnailUrl(baseData),
		ProductImage:             s.getProductImageUrl(baseData),
		ProductType:              product.Type(),
		ParentId:                 parentIdRef,
		VariantSelectedAttribute: variantSelectedAttributeRef,
		Retailer:                 baseData.RetailerCode,
		SKU:                      gtin,
		Manufacturer:             brand,
		Color:                    color,
		Size:                     size,
	}
}

func (s Factory) getProductThumbnailUrl(baseData productDomain.BasicProductData) string {
	media := baseData.GetMedia("thumbnail")
	if media.Reference != "" {
		return s.ProductMediaBaseUrl + s.ProductMediaThumbnailUrlPrefix + media.Reference
	}
	media = baseData.GetMedia("list")
	if media.Reference != "" {
		return s.ProductMediaBaseUrl + s.ProductMediaThumbnailUrlPrefix + media.Reference
	}
	return ""
}

func (s Factory) getProductImageUrl(baseData productDomain.BasicProductData) string {
	media := baseData.GetMedia("detail")
	if media.Reference != "" {
		return s.ProductMediaBaseUrl + s.ProductMediaUrlPrefix + media.Reference
	}
	return ""
}

func hashWithSHA512(value string) string {
	newHash := sha512.New()
	newHash.Write([]byte(value))
	//the hash is a byte array
	result := newHash.Sum(nil)
	//since we want to uuse it in a variable we base64 encode it (other alternative would be Hexadecimal representation "% x", h.Sum(nil)
	return base64.URLEncoding.EncodeToString(result)
}

func (s Factory) BuildChangeQtyEvent(productIdentifier string, productName string, qty int, qtyBefore int, cartId string) domain.Event {
	event := domain.Event{EventInfo: make(map[string]interface{})}
	event.EventInfo["productId"] = productIdentifier
	event.EventInfo["productName"] = productName
	event.EventInfo["cartId"] = cartId

	if qty == 0 {
		event.EventInfo["eventName"] = "Remove Product"
	} else {
		event.EventInfo["eventName"] = "Update Quantity"
		event.EventInfo["quantity"] = qty
	}
	return event
}

func (s Factory) BuildAddToBagEvent(productIdentifier string, productName string, qty int) domain.Event {

	event := domain.Event{EventInfo: make(map[string]interface{})}
	event.EventInfo["eventName"] = "Add To Bag"
	event.EventInfo["productId"] = productIdentifier
	event.EventInfo["productName"] = productName
	event.EventInfo["quantity"] = qty

	return event
}
