package application

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"strconv"
	"strings"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo-commerce/w3cDatalayer/domain"
	authApplication "flamingo.me/flamingo/core/auth/application"
	canonicalUrlApplication "flamingo.me/flamingo/core/canonicalUrl/application"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/session"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
	"go.opencensus.io/tag"
)

// Factory is used to build new datalayers
type Factory struct {
	router              *router.Router
	datalayerProvider   domain.DatalayerProvider
	canonicalURLService *canonicalUrlApplication.Service
	userService         *authApplication.UserService

	pageInstanceIDPrefix           string
	pageInstanceIDStage            string
	productMediaBaseURL            string
	productMediaURLPrefix          string
	productMediaThumbnailURLPrefix string
	pageNamePrefix                 string
	siteName                       string
	locale                         string
	defaultCurrency                string
	hashUserValues                 bool
}

// Inject factory dependencies
func (s *Factory) Inject(
	router2 *router.Router,
	provider domain.DatalayerProvider,
	service *canonicalUrlApplication.Service,
	userService *authApplication.UserService,
	config *struct {
		PageInstanceIDPrefix           string `inject:"config:w3cDatalayer.pageInstanceIDPrefix,optional"`
		PageInstanceIDStage            string `inject:"config:w3cDatalayer.pageInstanceIDStage,optional"`
		ProductMediaBaseURL            string `inject:"config:w3cDatalayer.productMediaBaseUrl,optional"`
		ProductMediaURLPrefix          string `inject:"config:w3cDatalayer.productMediaUrlPrefix,optional"`
		ProductMediaThumbnailURLPrefix string `inject:"config:w3cDatalayer.productMediaThumbnailUrlPrefix,optional"`
		PageNamePrefix                 string `inject:"config:w3cDatalayer.pageNamePrefix,optional"`
		SiteName                       string `inject:"config:w3cDatalayer.siteName,optional"`
		Locale                         string `inject:"config:locale.locale,optional"`
		DefaultCurrency                string `inject:"config:w3cDatalayer.defaultCurrency,optional"`
		HashUserValues                 bool   `inject:"config:w3cDatalayer.hashUserValues,optional"`
	},
) {
	s.router = router2
	s.datalayerProvider = provider
	s.canonicalURLService = service
	s.userService = userService

	s.pageInstanceIDPrefix = config.PageInstanceIDPrefix
	s.pageInstanceIDStage = config.PageInstanceIDStage
	s.productMediaBaseURL = config.ProductMediaBaseURL
	s.productMediaURLPrefix = config.ProductMediaURLPrefix
	s.productMediaThumbnailURLPrefix = config.ProductMediaThumbnailURLPrefix
	s.pageNamePrefix = config.PageNamePrefix
	s.siteName = config.SiteName
	s.locale = config.Locale
	s.defaultCurrency = config.DefaultCurrency
	s.hashUserValues = config.HashUserValues
}

// BuildForCurrentRequest builds the datalayer for the current request
func (s Factory) BuildForCurrentRequest(ctx context.Context, request *web.Request) domain.Datalayer {
	layer := s.datalayerProvider()

	//get langiage from locale code configuration
	language := ""
	localeParts := strings.Split(s.locale, "-")
	if len(localeParts) > 0 {
		language = localeParts[0]
	}

	layer.Page = &domain.Page{
		PageInfo: domain.PageInfo{
			PageID:         request.Request().URL.Path,
			PageName:       s.pageNamePrefix + request.Request().URL.Path,
			DestinationURL: s.canonicalURLService.GetCanonicalUrlForCurrentRequest(ctx),
			Language:       language,
		},
		Attributes: make(map[string]interface{}),
	}

	layer.Page.Attributes["currency"] = s.defaultCurrency

	//Use the handler name as PageId if available
	if controllerHandler, ok := tag.FromContext(ctx).Value(router.ControllerKey); ok {
		layer.Page.PageInfo.PageID = controllerHandler
	}

	layer.SiteInfo = &domain.SiteInfo{
		SiteName: s.siteName,
	}

	layer.PageInstanceID = s.pageInstanceIDPrefix + s.pageInstanceIDStage

	//Handle User
	layer.Page.Attributes["loggedIn"] = false
	if s.userService.IsLoggedIn(ctx, request.Session().G()) {
		layer.Page.Attributes["loggedIn"] = true
		layer.Page.Attributes["logintype"] = "external"
		userData := s.getUser(ctx, request.Session().G())
		if userData != nil {
			layer.User = append(layer.User, *userData)
		}
	} else {
		layer.Page.Attributes["logintype"] = "guest"
	}
	return *layer
}

func (s Factory) getUser(ctx context.Context, session *sessions.Session) *domain.User {

	dataLayerProfile := s.getUserProfileForCurrentUser(ctx, session)
	if dataLayerProfile == nil {
		return nil
	}

	dataLayerUser := domain.User{}
	dataLayerUser.Profile = append(dataLayerUser.Profile, *dataLayerProfile)
	return &dataLayerUser
}

func (s Factory) getUserProfileForCurrentUser(ctx context.Context, session *sessions.Session) *domain.UserProfile {
	user := s.userService.GetUser(ctx, session)
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

// HashValueIfConfigured returns the hashed `value` if hashing is configured
func (s Factory) HashValueIfConfigured(value string) string {
	if s.hashUserValues && value != "" {
		return hashWithSHA512(value)
	}
	return value
}

// BuildCartData builds the domain cart data
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
	for _, item := range cart.GetAllDecoratedItems() {
		itemData := s.buildCartItem(item, cart.Cart.CartTotals.CurrencyCode)
		cartData.Item = append(cartData.Item, itemData)
	}
	return &cartData
}

// BuildTransactionData builds the domain transaction data
func (s Factory) BuildTransactionData(ctx context.Context, cartTotals cart.CartTotals, decoratedItems []cart.DecoratedCartItem, orderID string, email string) *domain.Transaction {
	var profile *domain.UserProfile
	session, _ := session.FromContext(ctx)
	if s.userService.IsLoggedIn(ctx, session.G()) {
		profile = s.getUserProfileForCurrentUser(ctx, session.G())
	} else {
		profile = s.getUserProfile(email, "")
	}

	transactionData := domain.Transaction{
		TransactionID: orderID,
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
			BasePrice:    item.Item.SinglePrice,
			PriceWithTax: item.Item.SinglePriceInclTax,
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

// BuildProductData builds the domain product data
func (s Factory) BuildProductData(product productDomain.BasicProduct) domain.Product {
	productData := domain.Product{
		ProductInfo: s.getProductInfo(product),
		Category:    s.getProductCategory(product),
		Attributes:  make(map[string]interface{}),
	}

	// set prices
	productData.Attributes["productPrice"] = strconv.FormatFloat(product.SaleableData().ActivePrice.GetFinalPrice(), 'f', 2, 64)

	// check for highstreet price
	if product.BaseData().HasAttribute("rrp") {
		productData.Attributes["highstreetPrice"] = product.BaseData().Attributes["rrp"].Value()
	}

	// if FinalPrice is discounted, add it to specialPrice
	if product.SaleableData().ActivePrice.IsDiscounted {
		productData.Attributes["specialPrice"] = strconv.FormatFloat(product.SaleableData().ActivePrice.Discounted, 'f', 2, 64)
		productData.Attributes["productPrice"] = strconv.FormatFloat(product.SaleableData().ActivePrice.Default, 'f', 2, 64)
	}

	// set badge
	productData.Attributes["badge"] = s.EvaluateBadgeHierarchy(product)

	if product.BaseData().HasAttribute("ispuLimitedToAreas") {
		replacer := strings.NewReplacer("[", "", "]", "", " ", "|")
		productData.Attributes["ispuLimitedToAreas"] = replacer.Replace(product.BaseData().Attributes["ispuLimitedToAreas"].Value())
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
	var parentIDRef *string
	var variantSelectedAttributeRef *string
	if product.Type() == productDomain.TYPECONFIGURABLE_WITH_ACTIVE_VARIANT {
		if configurableWithActiveVariant, ok := product.(productDomain.ConfigurableProductWithActiveVariant); ok {
			parentID := configurableWithActiveVariant.ConfigurableBaseData().MarketPlaceCode
			parentIDRef = &parentID
			variantSelectedAttribute := configurableWithActiveVariant.ActiveVariant.BaseData().Attributes[configurableWithActiveVariant.VariantVariationAttributes[0]].Value()
			variantSelectedAttributeRef = &variantSelectedAttribute
		}
	}
	if product.Type() == productDomain.TYPECONFIGURABLE {
		if configurable, ok := product.(productDomain.ConfigurableProduct); ok {
			parentID := configurable.BaseData().MarketPlaceCode
			parentIDRef = &parentID
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
	if baseData.HasAttribute("brandName") {
		brand = baseData.Attributes["brandName"].Value()
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
		ProductThumbnail:         s.getProductThumbnailURL(baseData),
		ProductImage:             s.getProductImageURL(baseData),
		ProductType:              product.Type(),
		ParentID:                 parentIDRef,
		VariantSelectedAttribute: variantSelectedAttributeRef,
		Retailer:                 baseData.RetailerName,
		Brand:                    brand,
		SKU:                      gtin,
		Manufacturer:             brand,
		Color:                    color,
		Size:                     size,
		InStock:                  strconv.FormatBool(baseData.IsInStock()),
	}
}

func (s Factory) getProductThumbnailURL(baseData productDomain.BasicProductData) string {
	media := baseData.GetMedia("thumbnail")
	if media.Reference != "" {
		return s.productMediaBaseURL + s.productMediaThumbnailURLPrefix + media.Reference
	}
	media = baseData.GetMedia("list")
	if media.Reference != "" {
		return s.productMediaBaseURL + s.productMediaThumbnailURLPrefix + media.Reference
	}
	return ""
}

func (s Factory) getProductImageURL(baseData productDomain.BasicProductData) string {
	media := baseData.GetMedia("detail")
	if media.Reference != "" {
		return s.productMediaBaseURL + s.productMediaURLPrefix + media.Reference
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

// BuildChangeQtyEvent builds the change quantity domain event
func (s Factory) BuildChangeQtyEvent(productIdentifier string, productName string, qty int, qtyBefore int, cartID string) domain.Event {
	event := domain.Event{EventInfo: make(map[string]interface{})}
	event.EventInfo["productId"] = productIdentifier
	event.EventInfo["productName"] = productName
	event.EventInfo["cartId"] = cartID

	if qty == 0 {
		event.EventInfo["eventName"] = "Remove Product"
	} else {
		event.EventInfo["eventName"] = "Update Quantity"
		event.EventInfo["quantity"] = qty
	}
	return event
}

// BuildAddToBagEvent builds the add to bag domain event
func (s Factory) BuildAddToBagEvent(productIdentifier string, productName string, qty int) domain.Event {
	event := domain.Event{EventInfo: make(map[string]interface{})}
	event.EventInfo["eventName"] = "Add To Bag"
	event.EventInfo["productId"] = productIdentifier
	event.EventInfo["productName"] = productName
	event.EventInfo["quantity"] = qty

	return event
}

// EvaluateBadgeHierarchy get the active badge by product
func (s *Factory) EvaluateBadgeHierarchy(product productDomain.BasicProduct) string {
	badge := ""

	if product.BaseData().HasAttribute("airportBadge") {
		badge = "airportBadge"
	} else if product.BaseData().HasAttribute("retailerBadge") {
		badge = "retailerBadge"
	} else if product.BaseData().HasAttribute("exclusiveProduct") && product.BaseData().Attributes["exclusiveProduct"].Value() == "true" {
		badge = "travellerExclusive"
	} else if product.BaseData().IsNew {
		badge = "new"
	}

	return badge
}
