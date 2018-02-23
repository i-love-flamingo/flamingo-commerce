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
		//CurrentContext need to be set when using the service
		CurrentContext               web.Context
		Logger                       flamingo.Logger `inject:""`
		Factory                      *Factory        `inject:""`
		productDomain.ProductService `inject:""`
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

	s.AddSessionEvents()

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

func (s Service) AddSessionEvents() error {
	session := s.CurrentContext.Session()
	addToCartEvents := session.Flashes("addToCart")
	for _, event := range addToCartEvents {
		if addToCartEvent, ok := event.(cart.AddToCartEvent); ok {
			s.Logger.WithField("category", "w3cDatalayer").Println("addToCartEvent", addToCartEvent)
			product, err := s.ProductService.Get(s.CurrentContext, addToCartEvent.ProductIdentifier)

			if err != nil {
				return err
			}
			title := product.BaseData().Title

			s.AddToBagEvent(addToCartEvent.ProductIdentifier, title, addToCartEvent.Qty)
		}
	}

	changedQtyInCartEvents := session.Flashes("changedQtyInCart")
	for _, event := range changedQtyInCartEvents {
		if changedQtyInCartEvent, ok := event.(cart.ChangedQtyInCartEvent); ok {
			s.Logger.WithField("category", "w3cDatalayer").Println("changedQtyInCartEvent", changedQtyInCartEvent)
			product, err := s.ProductService.Get(s.CurrentContext, changedQtyInCartEvent.ProductIdentifier)

			if err != nil {
				return err
			}

			title := product.BaseData().Title
			s.AddChangeQtyEvent(changedQtyInCartEvent.ProductIdentifier, title, changedQtyInCartEvent.QtyAfter, changedQtyInCartEvent.QtyBefore, changedQtyInCartEvent.CartId)
		}
	}

	return nil
}

func (s Service) SetPageCategories(category string, subcategory1 string, subcategory2 string) error {
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

func (s Service) SetPageInfos(pageId string, pageName string) error {
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

func (s Service) SetCartData(cart cart.DecoratedCart) error {
	s.Logger.WithField("category", "w3cDatalayer").Debugf("Set Cart Data for cart %v", cart.Cart.ID)
	layer := s.Get()
	layer.Cart = s.Factory.BuildCartData(cart)
	return s.store(layer)
}

func (s Service) SetTransaction(cartTotals cart.CartTotals, decoratedItems []cart.DecoratedCartItem, orderId string) error {
	s.Logger.WithField("category", "w3cDatalayer").Debugf("Set Transaction Data for order %v", orderId)
	layer := s.Get()
	layer.Transaction = s.Factory.BuildTransactionData(cartTotals, decoratedItems, orderId)
	return s.store(layer)
}

// AddProduct - appends the productData to the datalayer
func (s Service) AddProduct(product productDomain.BasicProduct) error {
	layer := s.Get()
	layer.Product = append(layer.Product, s.Factory.BuildProductData(product))
	return s.store(layer)
}

//AddEvent - adds an event with the given eventName to the datalayer
func (s Service) AddEvent(eventName string, params ...*pugjs.Map) error {
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

func (s Service) AddToBagEvent(productIdentifier string, productName string, qty int) error {
	layer := s.Get()

	event := domain.Event{EventInfo: make(map[string]interface{})}
	event.EventInfo["eventName"] = "Add To Bag"
	event.EventInfo["productId"] = productIdentifier
	event.EventInfo["productName"] = productName
	event.EventInfo["quantity"] = qty

	layer.Event = append(layer.Event, event)
	return s.store(layer)
}

func (s Service) AddChangeQtyEvent(productIdentifier string, productName string, qty int, qtyBefore int, cartId string) error {
	layer := s.Get()

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

	layer.Event = append(layer.Event, event)
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
