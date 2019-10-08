package graphql

import (
	"context"
	"errors"
	cartForms "flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/form/domain"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCartMutationResolver resolves cart mutations
type CommerceCartMutationResolver struct {
	q                      *CommerceCartQueryResolver
	applicationCartService *application.CartService
	billingAddressFormController *cartForms.BillingAddressFormController
}

// Inject dependencies
func (r *CommerceCartMutationResolver) Inject(q *CommerceCartQueryResolver,
	applicationCartService *application.CartService,
	billingAddressFormController *cartForms.BillingAddressFormController) *CommerceCartMutationResolver {
	r.q = q
	r.applicationCartService = applicationCartService
	r.billingAddressFormController = billingAddressFormController
	return r
}

// CommerceAddToCart mutation for adding products to the current users cart
func (r *CommerceCartMutationResolver) CommerceAddToCart(ctx context.Context, marketplaceCode string, qty int, deliveryCode string) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	addRequest := r.applicationCartService.BuildAddRequest(ctx, marketplaceCode, "", qty)

	_, err := r.applicationCartService.AddProduct(ctx, req.Session(), deliveryCode, addRequest)
	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceDeleteItem resolver
func (r *CommerceCartMutationResolver) CommerceDeleteItem(ctx context.Context, itemID string, deliveryCode string) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	err := r.applicationCartService.DeleteItem(ctx, req.Session(), itemID, deliveryCode)

	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceDeleteCartDelivery mutation for removing deliveries from current users cart
func (r *CommerceCartMutationResolver) CommerceDeleteCartDelivery(ctx context.Context, deliveryCode string) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	_, err := r.applicationCartService.DeleteDelivery(ctx, req.Session(), deliveryCode)
	if err != nil {
		return nil, err
	}
	return r.q.CommerceCart(ctx)
}

// CommerceUpdateItemQty mutation for updating item quantity
func (r *CommerceCartMutationResolver) CommerceUpdateItemQty(ctx context.Context, itemID string, deliveryCode string, qty int) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	err := r.applicationCartService.UpdateItemQty(ctx, req.Session(), itemID, deliveryCode, qty)
	if err != nil {
		return nil, err
	}
	return r.q.CommerceCart(ctx)
}



// CommerceUpdateBillingAddress
// CommerceCartUpdateBillingAddress(ctx context.Context, addressForm *CommerceBillingAddressFormInput) (*graphql1.Commerce_Cart_BillingAddressForm, error)
//
func (r *CommerceCartMutationResolver) CommerceCartUpdateBillingAddress(ctx context.Context, address *cartForms.BillingAddressForm) (*Commerce_Cart_BillingAddressForm, error) {
	newRequest := web.CreateRequest(web.RequestFromContext(ctx).Request(),web.SessionFromContext(ctx))
	basicAddress := cartForms.AddressForm(*address)
	newRequest.Request().Form = basicAddress.UrlValues()

	form, success, err := r.billingAddressFormController.HandleFormAction(ctx,newRequest)
	if err != nil {
		return nil, err
	}
	return mapCommerce_Cart_BillingAddressForm(form,success)

}

//mapCommerce_Cart_BillingAddressForm - helper to map the graphql type Commerce_Cart_BillingAddressForm from common form
func mapCommerce_Cart_BillingAddressForm(form *domain.Form, success bool) (*Commerce_Cart_BillingAddressForm, error) {
	billingFormData, ok := form.Data.(cartForms.BillingAddressForm)
	if !ok {
		return nil, errors.New("unexpected form data")
	}

	var fieldErrors []Commerce_Cart_FieldError
	for fieldName,currentFieldErrors := range form.ValidationInfo.GetErrorsForAllFields() {
		for _, currentFieldError := range currentFieldErrors {
			fieldErrors = append(fieldErrors,Commerce_Cart_FieldError{
				MessageKey:   currentFieldError.MessageKey,
				DefaultLabel: currentFieldError.DefaultLabel,
				FieldName:    fieldName,
			})
		}
	}
	return &Commerce_Cart_BillingAddressForm{
		FormData:billingFormData,
		Processed:success,
		ValidationInfo: Commerce_Cart_ValidationInfo{
			GeneralErrors: form.ValidationInfo.GetGeneralErrors(),
			FieldErrors:   fieldErrors,
		},
	}, nil
}
