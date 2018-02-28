package application

import (
	"fmt"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/core/w3cDatalayer/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	ServiceProvider func() *Service
	/*
		Service can be used from outside is expected to be initialized with the current request context
		It stores a dataLayer Value object for the current request context and allows interaction with it
	*/
	Service struct {
		//currentContext need to be set when using the service
		currentContext               web.Context
		Logger                       flamingo.Logger `inject:""`
		Factory                      *Factory        `inject:""`
		productDomain.ProductService `inject:""`
	}
)

const (
	SESSION_EVENTS_KEY = "w3cdatalayer_events"
	DATALAYER_CTX_KEY  = "w3cDatalayer"
)

func (s *Service) Init(ctx web.Context) {
	s.currentContext = ctx
}

//Get gets the datalayer value object stored in the current context - or a freshly new build one if its the first call
func (s *Service) Get() domain.Datalayer {
	if s.currentContext == nil {
		s.Logger.WithField("category", "w3cDatalayer").Error("Get called without context!")

		return domain.Datalayer{}
	}
	if _, ok := s.currentContext.Value(DATALAYER_CTX_KEY).(domain.Datalayer); !ok {
		s.store(s.Factory.BuildForCurrentRequest(s.currentContext))
	}

	s.AddSessionEvents()

	if savedDataLayer, ok := s.currentContext.Value(DATALAYER_CTX_KEY).(domain.Datalayer); ok {
		return savedDataLayer
	}
	//error
	s.Logger.WithField("category", "w3cDatalayer").Warn("Receiving datalayer from context failed %v")
	return domain.Datalayer{}
}

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

func (s *Service) AddSessionEvents() error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	sessionEvents := s.currentContext.Session().Flashes(SESSION_EVENTS_KEY)

	// early return if there are no events
	if len(sessionEvents) == 0 {
		return nil
	}

	layer := s.Get()

	for _, event := range sessionEvents {
		if event, ok := event.(domain.Event); ok {
			s.Logger.WithField("category", "w3cDatalayer").Debugf("SESSION_EVENTS_KEY Event", event.EventInfo)
			layer.Event = append(layer.Event, event)
		}
	}

	err := s.store(layer)
	if err != nil {
		return err
	}

	return nil
}

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

//SetUserEmail to a User object
func (s *Service) SetUserEmail(mail string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.Logger.WithField("category", "w3cDatalayer").Debugf("Set Usermail %v", mail)
	layer := s.Get()
	layer.User = append(layer.User, domain.User{
		Profile: []domain.UserProfile{domain.UserProfile{
			ProfileInfo: domain.UserProfileInfo{
				EmailID: mail,
			},
		}},
	})
	return s.store(layer)
}

func (s *Service) SetCartData(cart cart.DecoratedCart) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.Logger.WithField("category", "w3cDatalayer").Debugf("Set Cart Data for cart %v", cart.Cart.ID)
	layer := s.Get()
	layer.Cart = s.Factory.BuildCartData(cart)
	return s.store(layer)
}

func (s *Service) SetTransaction(cartTotals cart.CartTotals, decoratedItems []cart.DecoratedCartItem, orderId string) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	s.Logger.WithField("category", "w3cDatalayer").Debugf("Set Transaction Data for order %v", orderId)
	layer := s.Get()
	layer.Transaction = s.Factory.BuildTransactionData(cartTotals, decoratedItems, orderId)
	return s.store(layer)
}

// AddProduct - appends the productData to the datalayer
func (s *Service) AddProduct(product productDomain.BasicProduct) error {
	if s.currentContext == nil {
		return errors.New("Service can only be used with currentContext - call Init() first")
	}
	layer := s.Get()
	layer.Product = append(layer.Product, s.Factory.BuildProductData(product))
	return s.store(layer)
}

//AddEvent - adds an event with the given eventName to the datalayer
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

//store datalayer in current context
func (s *Service) store(layer domain.Datalayer) error {
	s.Logger.Debugf("Update %#v", layer)
	if s.currentContext == nil {
		s.Logger.WithField("category", "w3cDatalayer").Error("Update called without context!")
		return errors.New("Update called without context")
	}
	s.currentContext.WithValue(DATALAYER_CTX_KEY, layer)

	return nil
}
