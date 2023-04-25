package application

import (
	"context"
	"errors"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/sourcing/domain"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// SourcingApplication interface
	SourcingApplication interface {
		GetAvailableSourcesDeductedByCurrentCart(ctx context.Context, session *web.Session, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSourcesPerProduct, error)
		GetAvailableSources(ctx context.Context, session *web.Session, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSourcesPerProduct, error)
	}

	// Service to access the sourcing based on current cart
	Service struct {
		logger              flamingo.Logger
		sourcingService     domain.SourcingService
		cartReceiverService *application.CartReceiverService
		deliveryInfoBuilder cart.DeliveryInfoBuilder
	}
)

var (
	_ SourcingApplication = new(Service)
)

// Inject dependencies
func (s *Service) Inject(
	l flamingo.Logger,
	cartReceiverService *application.CartReceiverService,
	sourcingService domain.SourcingService,
	deliveryInfoBuilder cart.DeliveryInfoBuilder,
) *Service {
	s.logger = l.WithField(flamingo.LogKeyModule, "sourcing").WithField(flamingo.LogKeyCategory, "Application.Service")

	s.cartReceiverService = cartReceiverService
	s.deliveryInfoBuilder = deliveryInfoBuilder
	s.sourcingService = sourcingService

	return s
}

// GetAvailableSourcesDeductedByCurrentCart fetches available sources minus those already allocated to the cart
func (s *Service) GetAvailableSourcesDeductedByCurrentCart(ctx context.Context, session *web.Session, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSourcesPerProduct, error) {
	if product == nil {
		s.logger.WithContext(ctx).Error("No product given for GetAvailableSourcesDeductedByCurrentCart")
		return nil, errors.New("no product given for GetAvailableSourcesDeductedByCurrentCart")
	}

	deliveryInfo, decoratedCart, err := s.getDeliveryInfo(ctx, session, deliveryCode)
	if err != nil {
		return nil, err
	}

	return s.sourcingService.GetAvailableSources(ctx, product, deliveryInfo, decoratedCart)
}

// GetAvailableSources without evaluating current cart items
func (s *Service) GetAvailableSources(ctx context.Context, session *web.Session, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSourcesPerProduct, error) {
	if product == nil {
		s.logger.WithContext(ctx).Error("No product given for GetAvailableSources")
		return nil, errors.New("no product given for GetAvailableSources")
	}

	deliveryInfo, _, err := s.getDeliveryInfo(ctx, session, deliveryCode)
	if err != nil {
		return nil, err
	}

	return s.sourcingService.GetAvailableSources(ctx, product, deliveryInfo, nil)
}

func (s *Service) getDeliveryInfo(ctx context.Context, session *web.Session, deliveryCode string) (*cart.DeliveryInfo, *decorator.DecoratedCart, error) {
	decoratedCart, err := s.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		s.logger.WithContext(ctx).Error(err)
		return nil, nil, err
	}
	var deliveryInfo *cart.DeliveryInfo
	delivery, found := decoratedCart.Cart.GetDeliveryByCode(deliveryCode)
	if !found {
		deliveryInfo, err = s.deliveryInfoBuilder.BuildByDeliveryCode(deliveryCode)
		if err != nil {
			s.logger.WithContext(ctx).Error(err)
			return nil, decoratedCart, err
		}
	} else {
		deliveryInfo = &delivery.DeliveryInfo
	}
	return deliveryInfo, decoratedCart, nil
}
