package application

import (
	"github.com/pkg/errors"

	"go.aoe.com/flamingo/core/w3cDatalayer/domain"
	"go.aoe.com/flamingo/framework/flamingo"

	"go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
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

func (s Service) SetCartData(cart cart.DecoratedCart) error {
	layer := s.Get()
	layer.Cart = s.Factory.BuildCartData(cart)
	return s.store(layer)
}

// AddProduct - appends the productData to the datalayer
func (s Service) AddProduct(product productDomain.BasicProduct) error {
	layer := s.Get()
	layer.Product = append(layer.Product, s.Factory.BuildProductData(product))
	return s.store(layer)
}

//AddEvent - adds an event with the given eventName to the datalayer
func (s Service) AddEvent(eventName string) error {
	layer := s.Get()
	layer.Event = append(layer.Event, domain.Event{
		EventInfo: domain.EventInfo{
			EventName: eventName,
		},
	})
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
