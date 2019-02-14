package application

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/w3cDatalayer/domain"
	"flamingo.me/pugtemplate/pugjs"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// ServiceProvider func
	ServiceProvider func() *Service

	// Service can be used from outside is expected to be initialized with the current request context
	// It stores a dataLayer Value object for the current request context and allows interaction with it
	Service struct {
		//currentContext need to be set when using the service
		currentContext context.Context
		logger         flamingo.Logger
		factory        *Factory
		productDomain.ProductService
	}
)

const (
	SESSION_EVENTS_KEY = "w3cdatalayer_events"
	DATALAYER_REQ_KEY  = "w3cDatalayer"
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

// Get gets the datalayer value object stored in the current context - or a freshly new build one if its the first call
func (s *Service) Get() domain.Datalayer {
	if s.currentContext == nil {
		s.logger.WithField("category", "w3cDatalayer").Error("Get called without context!")

		return domain.Datalayer{}
	}
	req, _ := web.FromContext(s.currentContext)
	if _, ok := req.Values.Load(DATALAYER_REQ_KEY); !ok {
		s.store(s.factory.BuildForCurrentRequest(s.currentContext, req))
	}

	s.AddSessionEvents()

	layer, _ := req.Values.Load(DATALAYER_REQ_KEY)
	if savedDataLayer, ok := layer.(domain.Datalayer); ok {
		return savedDataLayer
	}

	//error
	s.logger.WithField("category", "w3cDatalayer").Warn("Receiving datalayer from context failed %v")
	return domain.Datalayer{}
}

// SetBreadCrumb to datalayer
func (s *Service) SetBreadCrumb(breadcrumb string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	if layer.Page != nil {

	}
	layer.Page.PageInfo.BreadCrumbs = breadcrumb
	return s.store(layer)
}

// AddSessionEvents to datalayer
func (s *Service) AddSessionEvents() error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	session, _ := session.FromContext(s.currentContext)
	sessionEvents := session.Flashes(SESSION_EVENTS_KEY)

	// early return if there are no events
	if len(sessionEvents) == 0 {
		return nil
	}

	layer := s.Get()

	for _, event := range sessionEvents {
		if event, ok := event.(domain.Event); ok {
			s.logger.WithField("category", "w3cDatalayer").Debug("SESSION_EVENTS_KEY Event", flamingo.EventInfo)
			layer.Event = append(layer.Event, event)
		}
	}

	err := s.store(layer)
	if err != nil {
		return err
	}

	return nil
}

// SetPageCategories to datalayer
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

// SetPageInfos to datalayer
func (s *Service) SetPageInfos(pageId string, pageName string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	if layer.Page == nil {
		layer.Page = &domain.Page{}
	}
	if pageId != "" {
		layer.Page.PageInfo.PageID = pageId
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
	s.logger.WithField("category", "w3cDatalayer").Debug("Set Usermail %v", mail)
	layer := s.Get()
	layer.User = append(layer.User, domain.User{
		Profile: []domain.UserProfile{domain.UserProfile{
			ProfileInfo: domain.UserProfileInfo{
				EmailID: s.factory.HashValueIfConfigured(mail),
			},
		}},
	})
	return s.store(layer)
}

// SetSearchData to datalayer
func (s *Service) SetSearchData(keyword string, results interface{}) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("category", "w3cDatalayer").Debug("SetSearchData Keyword %v Result: %#v", keyword, results)
	layer := s.Get()
	if layer.Page != nil {
		layer.Page.Search = domain.SearchInfo{
			SearchKeyword: keyword,
			Result:        results,
		}
	}
	return s.store(layer)
}

// SetCartData to datalayer
func (s *Service) SetCartData(cart cart.DecoratedCart) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("category", "w3cDatalayer").Debug("Set Cart Data for cart %v", cart.Cart.ID)
	layer := s.Get()
	layer.Cart = s.factory.BuildCartData(cart)
	return s.store(layer)
}

// SetTransaction information to datalayer
func (s *Service) SetTransaction(cartTotals cart.CartTotals, decoratedItems []cart.DecoratedCartItem, orderId string, email string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.logger.WithField("category", "w3cDatalayer").Debug("Set Transaction Data for order %v mail %v", orderId, email)
	layer := s.Get()
	layer.Transaction = s.factory.BuildTransactionData(s.currentContext, cartTotals, decoratedItems, orderId, email)
	return s.store(layer)
}

// AddTransactionAttribute to datalayer
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

// AddProduct - appends the productData to the datalayer
func (s *Service) AddProduct(product productDomain.BasicProduct) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	layer.Product = append(layer.Product, s.factory.BuildProductData(product))
	return s.store(layer)
}

// AddEvent - adds an event with the given eventName to the datalayer
func (s *Service) AddEvent(eventName string, params ...*pugjs.Map) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()

	event := domain.Event{EventInfo: make(map[string]interface{})}
	event.EventInfo["eventName"] = eventName

	if len(params) == 1 {
		for k, v := range params[0].Items {
			event.EventInfo[k.String()] = fmt.Sprint(v)
		}
	}

	layer.Event = append(layer.Event, event)
	return s.store(layer)
}

// store datalayer in current context
func (s *Service) store(layer domain.Datalayer) error {
	s.logger.Debug("Update %#v", layer)
	if s.currentContext == nil {
		s.logger.WithField("category", "w3cDatalayer").Error("Update called without context!")
		return errors.New("Update called without context")
	}
	req, _ := web.FromContext(s.currentContext)
	req.Values.Store(DATALAYER_REQ_KEY, layer)

	return nil
}
