package application

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"

	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/w3cdatalayer/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/pugtemplate/pugjs"
	"github.com/pkg/errors"
)

type (
	// ServiceProvider func
	ServiceProvider func() *Service

	// Service can be used from outside is expected to be initialized with the current request context
	// It stores a dataLayer Value object for the current request context and allows interaction with it
	Service struct {
		// currentContext need to be set when using the service
		currentContext context.Context
		logger         flamingo.Logger
		factory        *Factory
		productDomain.ProductService
	}
)

// constants
const (
	SessionEventsKey = "w3cdatalayer_events"
	DatalayerReqKey  = "w3cDatalayer"
)

// Inject method
func (s *Service) Inject(logger flamingo.Logger, factory *Factory, service productDomain.ProductService) {
	s.logger = logger
	s.factory = factory
	s.ProductService = service
}

// Init method - sets the context
func (s *Service) Init(ctx context.Context) {
	s.currentContext = ctx
}

// Get gets the data layer value object stored in the current context - or a freshly new build one if its the first call
func (s *Service) Get() domain.Datalayer {
	if s.currentContext == nil {
		s.logger.WithField("category", "w3cDatalayer").Error("Get called without context!")

		return domain.Datalayer{}
	}
	req := web.RequestFromContext(s.currentContext)
	if _, ok := req.Values.Load(DatalayerReqKey); !ok {
		_ = s.store(s.factory.BuildForCurrentRequest(s.currentContext, req))
	}

	_ = s.AddSessionEvents()

	layer, _ := req.Values.Load(DatalayerReqKey)
	if savedDataLayer, ok := layer.(domain.Datalayer); ok {
		return savedDataLayer
	}

	// error
	s.logger.WithField("category", "w3cDatalayer").Warn("Receiving datalayer from context failed")
	return domain.Datalayer{}
}

// SetBreadCrumb to data layer
func (s *Service) SetBreadCrumb(breadcrumb string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	if layer.Page != nil {
		layer.Page.PageInfo.BreadCrumbs = breadcrumb
	}

	return s.store(layer)
}

// AddSessionEvents to data layer
func (s *Service) AddSessionEvents() error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	session := web.SessionFromContext(s.currentContext)
	sessionEvents := session.Flashes(SessionEventsKey)

	// early return if there are no events
	if len(sessionEvents) == 0 {
		return nil
	}

	layer := s.Get()

	for _, event := range sessionEvents {
		if event, ok := event.(domain.Event); ok {
			s.logger.WithField("category", "w3cDatalayer").Debug("SESSION_EVENTS_KEY Event", event.EventInfo)
			layer.Event = append(layer.Event, event)
		}
	}

	err := s.store(layer)
	if err != nil {
		return err
	}

	return nil
}

// SetPageCategories to data layer
func (s *Service) SetPageCategories(category string, subcategory1 string, subcategory2 string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	if layer.Page == nil {
		layer.Page = &domain.Page{}
	}
	layer.Page.Category.PrimaryCategory = category
	layer.Page.Category.Section = category

	layer.Page.Category.SubCategory1 = subcategory1
	layer.Page.Category.SubCategory2 = subcategory2

	return s.store(layer)
}

// SetPageInfos to data layer
func (s *Service) SetPageInfos(pageID string, pageName string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	if layer.Page == nil {
		layer.Page = &domain.Page{}
	}
	if pageID != "" {
		layer.Page.PageInfo.PageID = pageID
	}
	if pageName != "" {
		layer.Page.PageInfo.PageName = pageName
	}
	return s.store(layer)
}

// SetUserEmail to a User object
func (s *Service) SetUserEmail(mail string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("category", "w3cDatalayer").Debug(fmt.Sprintf("Set Usermail %v", mail))
	layer := s.Get()
	layer.User = append(layer.User, domain.User{
		Profile: []domain.UserProfile{{
			ProfileInfo: domain.UserProfileInfo{
				EmailID: s.factory.HashValueIfConfigured(mail),
			},
		}},
	})
	return s.store(layer)
}

// SetPageInfoLanguage to the layer
func (s *Service) SetPageInfoLanguage(language string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("language", "w3cDatalayer").Debug(fmt.Sprintf("Set page language %v", language))
	layer := s.Get()

	if layer.Page == nil {
		return nil
	}
	layer.Page.PageInfo.Language = language
	return s.store(layer)
}

// SetSearchData to data layer
func (s *Service) SetSearchData(keyword string, results interface{}) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("category", "w3cDatalayer").Debug(fmt.Sprintf("SetSearchData Keyword %v Result: %#v", keyword, results))
	layer := s.Get()
	if layer.Page != nil {
		layer.Page.Search = domain.SearchInfo{
			SearchKeyword: keyword,
			Result:        results,
		}
	}
	return s.store(layer)
}

// SetCartData to data layer
func (s *Service) SetCartData(cart decorator.DecoratedCart) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("category", "w3cDatalayer").Debug(fmt.Sprintf("Set Cart Data for cart %v", cart.Cart.ID))
	layer := s.Get()
	layer.Cart = s.factory.BuildCartData(cart)
	return s.store(layer)
}

// SetTransaction information to data layer
func (s *Service) SetTransaction(cart decorator.DecoratedCart, decoratedItems []decorator.DecoratedCartItem, orderID string, email string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("category", "w3cDatalayer").Debug(fmt.Sprintf("Set Transaction Data for order %v mail %v", orderID, email))
	layer := s.Get()
	layer.Transaction = s.factory.BuildTransactionData(s.currentContext, cart, decoratedItems, orderID, email)
	return s.store(layer)
}

// AddTransactionAttribute to data layer
func (s *Service) AddTransactionAttribute(key string, value string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	if layer.Transaction != nil && layer.Transaction.Attributes != nil {
		layer.Transaction.Attributes[key] = value
	}
	return s.store(layer)
}

// AddProduct - appends the productData to the data layer
func (s *Service) AddProduct(product productDomain.BasicProduct) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	layer.Product = append(layer.Product, s.factory.BuildProductData(product))
	return s.store(layer)
}

// AddEvent - adds an event with the given eventName to the data layer
func (s *Service) AddEvent(eventName string, params ...*pugjs.Map) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()

	event := domain.Event{EventInfo: make(map[string]interface{})}
	event.EventInfo["eventName"] = eventName

	if len(params) == 1 {
		for k, v := range params[0].AsStringMap() {
			event.EventInfo[k] = v
		}
	}

	layer.Event = append(layer.Event, event)
	return s.store(layer)
}

// store data layer in current context
func (s *Service) store(layer domain.Datalayer) error {
	s.logger.Debug(fmt.Sprintf("Update %#v", layer))
	if s.currentContext == nil {
		s.logger.WithField("category", "w3cDatalayer").Error("Update called without context!")
		return errors.New("Update called without context")
	}
	req := web.RequestFromContext(s.currentContext)
	req.Values.Store(DatalayerReqKey, layer)

	return nil
}
