package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	formDomain "flamingo.me/form/domain"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	Commerce_Cart_BillingAddressForm struct {
		FormData forms.BillingAddressForm
		Processed bool
		ValidationInfo Commerce_Cart_ValidationInfo
	}

	Commerce_Cart_ValidationInfo struct {
		GeneralErrors []formDomain.Error
		FieldErrors []Commerce_Cart_FieldError
	}

	Commerce_Cart_FieldError struct {
		// MessageKey - a key of the error message. Often used to pass to translation func in the template
		MessageKey string
		// DefaultLabel - a speaking error label. OFten used to show to end user - in case no translation exists
		DefaultLabel string
		FieldName string
	}

)
// CommerceCartQueryResolver resolver for carts
type CommerceCartQueryResolver struct {
	applicationCartReceiverService *application.CartReceiverService
	billingAddressFormController *forms.BillingAddressFormController
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

//Commerce_Cart_BillingAddressForm getter
func (r *CommerceCartQueryResolver) CommerceCartGetBillingAddressForm(ctx context.Context) (*Commerce_Cart_BillingAddressForm, error) {
	req := web.RequestFromContext(ctx)
	billingForm, err :=  r.billingAddressFormController.GetUnsubmittedForm(ctx, req)
	if err != nil {
		return nil, err
	}

	return mapCommerce_Cart_BillingAddressForm(billingForm,false)

}


