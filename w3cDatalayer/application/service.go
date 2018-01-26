package application

import (
	"strings"

	"crypto/sha512"

	"encoding/base64"

	"github.com/pkg/errors"
	authApplication "go.aoe.com/flamingo/core/auth/application"
	canonicalUrlApplication "go.aoe.com/flamingo/core/canonicalUrl/application"
	"go.aoe.com/flamingo/core/w3cDatalayer/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	ServiceProvider func() *Service
	/*
		Service can be used from outside is expected to be initialized with the current request context
		It stores a dataLayer Value object for the current request context and allows interaction with it
	*/
	Service struct {
		//CurrentContext need to be set when using the service
		CurrentContext web.Context
		Logger         flamingo.Logger `inject:""`
		Factory        Factory         `inject:""`
	}

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

//Get gets the datalayer value object stored in the current context - or a freshly new build one if its the first call
func (s Service) Get() domain.Datalayer {
	if s.CurrentContext == nil {
		s.Logger.WithField("category", "w3cDatalayer").Error("Get called without context!")

		return domain.Datalayer{}
	}
	if _, ok := s.CurrentContext.Value("w3cDatalayer").(domain.Datalayer); !ok {
		s.store(s.Factory.BuildForCurrentRequest(s.CurrentContext))
	}

	if savedDataLayer, ok := s.CurrentContext.Value("w3cDatalayer").(domain.Datalayer); ok {
		return savedDataLayer
	}
	//error
	s.Logger.WithField("category", "w3cDatalayer").Warn("Receiving datalayer from context failed %v")
	return domain.Datalayer{}
}

func (s Service) SetBreadCrumb(breadcrumb string) error {
	layer := s.Get()
	layer.Page.PageInfo.BreadCrumbs = breadcrumb
	return s.store(layer)
}

func (s Service) SetPageCategories(category string, subcategory1 string, subcategory2 string) error {
	layer := s.Get()
	layer.Page.Category.PrimaryCategory = category
	layer.Page.Category.Section = category

	layer.Page.Category.SubCategory1 = subcategory1
	layer.Page.Category.SubCategory2 = subcategory2

	return s.store(layer)
}

//store datalayer in current context
func (s Service) store(layer domain.Datalayer) error {
	s.Logger.Debugf("Update %#v", layer)
	if s.CurrentContext == nil {
		s.Logger.WithField("category", "w3cDatalayer").Error("Update called without context!")
		return errors.New("Update called without context")
	}
	s.CurrentContext.WithValue("w3cDatalayer", layer)

	return nil
}

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

func hashWithSHA512(value string) string {
	newHash := sha512.New()
	newHash.Write([]byte(value))
	//the hash is a byte array
	result := newHash.Sum(nil)
	//since we want to uuse it in a variable we base64 encode it (other alternative would be Hexadecimal representation "% x", h.Sum(nil)
	return base64.URLEncoding.EncodeToString(result)
}
