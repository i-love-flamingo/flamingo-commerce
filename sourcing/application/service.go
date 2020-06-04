package application

import (
	"context"
	"errors"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/sourcing/domain"

	"flamingo.me/flamingo-commerce/v3/cart/application"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// SourcingApplication interface
	SourcingApplication interface {
		GetAvailableSourcesDeductedByCurrentCart(ctx context.Context, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSources, error)
		GetAvailableSources(ctx context.Context, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSources, error)
	}

	// Service to access the ourcing based on current cart
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
func (r *Service) Inject(
	l flamingo.Logger,
	cartReceiverService *application.CartReceiverService,
	sourcingService domain.SourcingService,
	deliveryInfoBuilder cart.DeliveryInfoBuilder,
) *Service {
	r.logger = l.WithField(flamingo.LogKeyModule, "sourcing").WithField(flamingo.LogKeyCategory, "Application.Service")

	r.cartReceiverService = cartReceiverService
	r.deliveryInfoBuilder = deliveryInfoBuilder
	r.sourcingService = sourcingService

	return r
}

// GetAvailableSourcesDeductedByCurrentCart fetches available sources minus those already allocated to the cart
func (r *Service) GetAvailableSourcesDeductedByCurrentCart(ctx context.Context, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSources, error) {
	if product == nil {
		r.logger.WithContext(ctx).Error("No product given for GetAvailableSourcesDeductedByCurrentCart")
		return nil, errors.New("no product given for GetAvailableSourcesDeductedByCurrentCart")
	}

	deliveryInfo, decoratedCart, err := r.getDeliveryInfo(ctx, deliveryCode)
	if err != nil {
		return nil, err
	}

	return r.sourcingService.GetAvailableSources(ctx, product, deliveryInfo, decoratedCart)
}

// GetAvailableSources without evaluating current cart items
func (r *Service) GetAvailableSources(ctx context.Context, product productDomain.BasicProduct, deliveryCode string) (domain.AvailableSources, error) {

	if product == nil {
		r.logger.WithContext(ctx).Error("No product given for GetAvailableSources")
		return nil, errors.New("no product given for GetAvailableSources")
	}

	deliveryInfo, _, err := r.getDeliveryInfo(ctx, deliveryCode)
	if err != nil {
		return nil, err
	}

	return r.sourcingService.GetAvailableSources(ctx, product, deliveryInfo, nil)
}

func (r *Service) getDeliveryInfo(ctx context.Context, deliveryCode string) (*cart.DeliveryInfo, *decorator.DecoratedCart, error) {
	session := web.SessionFromContext(ctx)
	decoratedCart, err := r.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		r.logger.WithContext(ctx).Error(err)
		return nil, nil, err
	}
	var deliveryInfo *cart.DeliveryInfo
	delivery, found := decoratedCart.Cart.GetDeliveryByCode(deliveryCode)
	if !found {
		deliveryInfo, err = r.deliveryInfoBuilder.BuildByDeliveryCode(deliveryCode)
		if err != nil {
			r.logger.WithContext(ctx).Error(err)
			return nil, decoratedCart, err
		}
	} else {
		deliveryInfo = &delivery.DeliveryInfo
	}
	return deliveryInfo, decoratedCart, nil
}
