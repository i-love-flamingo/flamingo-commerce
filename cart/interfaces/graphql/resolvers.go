package graphql

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

// CommerceCartQueryResolver resolver for carts
type CommerceCartQueryResolver struct {
	applicationCartReceiverService *application.CartReceiverService
	applicationCartService         *application.CartService
	restrictionService             *validation.RestrictionService
	productService                 domain.ProductService
	logger                         flamingo.Logger
}

// Inject dependencies
func (r *CommerceCartQueryResolver) Inject(
	applicationCartReceiverService *application.CartReceiverService,
	cartService *application.CartService,
	restrictionService *validation.RestrictionService,
	productService domain.ProductService,
	logger flamingo.Logger,
) {
	r.applicationCartReceiverService = applicationCartReceiverService
	r.applicationCartService = cartService
	r.restrictionService = restrictionService
	r.productService = productService
	r.logger = logger.WithField(flamingo.LogKeyModule, "cart").WithField(flamingo.LogKeyCategory, "graphql")
}

// CommerceCart getter for queries
func (r *CommerceCartQueryResolver) CommerceCart(ctx context.Context) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	dc, err := r.applicationCartReceiverService.ViewDecoratedCart(ctx, req.Session())
	if err != nil {
		r.logger.Error("Failed to view decorated cart", err)
		return nil, interfaces.ErrCartGeneral
	}

	return dto.NewDecoratedCart(dc), nil
}

// CommerceCartValidator to trigger the cart validation service
func (r *CommerceCartQueryResolver) CommerceCartValidator(ctx context.Context) (*validation.Result, error) {
	session := web.SessionFromContext(ctx)

	decoratedCart, err := r.applicationCartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		r.logger.Error("Failed to view decorated cart for validation", err)
		return nil, interfaces.ErrCartGeneral
	}

	result := r.applicationCartService.ValidateCart(ctx, session, decoratedCart)

	return &result, nil
}

// CommerceCartQtyRestriction checks if given sku is restricted in terms of qty
func (r *CommerceCartQueryResolver) CommerceCartQtyRestriction(ctx context.Context, marketplaceCode string, variantCode *string, deliveryCode string) (*validation.RestrictionResult, error) {
	session := web.SessionFromContext(ctx)

	product, err := r.productService.Get(ctx, marketplaceCode)
	if err != nil {
		r.logger.Error("Failed to get product for qty restriction", err)
		return nil, interfaces.ErrCartGeneral
	}
	if variantCode != nil {
		if configurableProduct, ok := product.(domain.ConfigurableProduct); ok {
			product, err = configurableProduct.GetConfigurableWithActiveVariant(*variantCode)
			if err != nil {
				r.logger.Error("Failed to get active variant product for qty restriction", err)
				return nil, interfaces.ErrCartGeneral
			}
		}
	}

	cart, err := r.applicationCartReceiverService.ViewCart(ctx, session)
	if err != nil {
		r.logger.Error("Failed to view cart for qty restriction", err)
		return nil, interfaces.ErrCartGeneral
	}
	result := r.restrictionService.RestrictQty(ctx, session, product, cart, deliveryCode)
	return result, nil
}

// CartSplit returns graphql specific cart split
func (r *CommerceCartQueryResolver) CartSplit(_ context.Context, paymentSelection *cart.DefaultPaymentSelection) ([]*dto.PaymentSelectionSplit, error) {
	if paymentSelection == nil {
		return nil, nil
	}

	paymentSelectionSplit := make([]*dto.PaymentSelectionSplit, 0)
	for qualifier, charge := range paymentSelection.CartSplit() {
		paymentSelectionSplit = append(paymentSelectionSplit, &dto.PaymentSelectionSplit{
			Qualifier: qualifier,
			Charge:    charge,
		})
	}

	return paymentSelectionSplit, nil
}

// CommerceCartAdditionalDataResolver resolver for custom attributes of cart
type CommerceCartAdditionalDataResolver struct{}

// CustomAttributes of cart
func (r *CommerceCartAdditionalDataResolver) CustomAttributes(_ context.Context, obj *cart.AdditionalData) (*dto.CustomAttributes, error) {
	return &dto.CustomAttributes{Attributes: obj.CustomAttributes}, nil
}

// CommerceCartDeliveryInfoResolver resolver for additional data of delivery info
type CommerceCartDeliveryInfoResolver struct{}

// AdditionalData of delivery info
func (r *CommerceCartDeliveryInfoResolver) AdditionalData(_ context.Context, obj *cart.DeliveryInfo) (*dto.CustomAttributes, error) {
	return &dto.CustomAttributes{Attributes: obj.AdditionalData}, nil
}
