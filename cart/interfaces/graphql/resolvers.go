package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCartQueryResolver resolver for carts
type CommerceCartQueryResolver struct {
	applicationCartReceiverService *application.CartReceiverService
	billingAddressFormController   *forms.BillingAddressFormController
}

// Inject dependencies
func (r *CommerceCartQueryResolver) Inject(applicationCartReceiverService *application.CartReceiverService, billingAddressFormController *forms.BillingAddressFormController) {
	r.applicationCartReceiverService = applicationCartReceiverService
	r.billingAddressFormController = billingAddressFormController
}

// CommerceCart getter for queries
func (r *CommerceCartQueryResolver) CommerceCart(ctx context.Context) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	return r.applicationCartReceiverService.ViewDecoratedCart(ctx, req.Session())
}
