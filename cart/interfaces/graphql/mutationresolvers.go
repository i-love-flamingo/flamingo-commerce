package graphql

import (
	"context"
	"errors"
	cartForms "flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	formApplication "flamingo.me/form/application"
	"flamingo.me/form/domain"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCartMutationResolver resolves cart mutations
type CommerceCartMutationResolver struct {
	q                            *CommerceCartQueryResolver
	applicationCartService       *application.CartService
	billingAddressFormController *cartForms.BillingAddressFormController
	simplePaymentFormController  *cartForms.SimplePaymentFormController
	formDataEncoderFactory       formApplication.FormDataEncoderFactory
	cartService                  *application.CartService
}

// Inject dependencies
func (r *CommerceCartMutationResolver) Inject(q *CommerceCartQueryResolver,
	applicationCartService *application.CartService,
	billingAddressFormController *cartForms.BillingAddressFormController,
	formDataEncoderFactory formApplication.FormDataEncoderFactory,
	simplePaymentFormController *cartForms.SimplePaymentFormController,
	cartService *application.CartService) *CommerceCartMutationResolver {
	r.q = q
	r.applicationCartService = applicationCartService
	r.billingAddressFormController = billingAddressFormController
	r.formDataEncoderFactory = formDataEncoderFactory
	r.simplePaymentFormController = simplePaymentFormController
	r.cartService = cartService
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

//CommerceCartUpdateBillingAddress resolver method
func (r *CommerceCartMutationResolver) CommerceCartUpdateBillingAddress(ctx context.Context, address *cartForms.BillingAddressForm) (*dto.BillingAddressForm, error) {
	newRequest := web.CreateRequest(web.RequestFromContext(ctx).Request(), web.SessionFromContext(ctx))
	v, err := r.formDataEncoderFactory.CreateByNamedEncoder("commerce.cart.billingFormService").Encode(ctx, address)
	if err != nil {
		return nil, err
	}
	newRequest.Request().Form = v

	form, success, err := r.billingAddressFormController.HandleFormAction(ctx, newRequest)
	if err != nil {
		return nil, err
	}
	return mapCommerceBillingAddressForm(form, success)

}

//CommerceCartUpdateSelectedPayment resolver method
func (r *CommerceCartMutationResolver) CommerceCartUpdateSelectedPayment(ctx context.Context, gateway string, method string) (*dto.SelectedPaymentResult, error) {
	newRequest := web.CreateRequest(web.RequestFromContext(ctx).Request(), web.SessionFromContext(ctx))
	urlValues := make(url.Values)
	urlValues["gateway"] = []string{gateway}
	urlValues["method"] = []string{method}
	newRequest.Request().Form = urlValues

	form, success, err := r.simplePaymentFormController.HandleFormAction(ctx, newRequest)
	if err != nil {
		return nil, err
	}

	return &dto.SelectedPaymentResult{
		Processed: success,
		ValidationInfo: dto.ValidationInfo{
			GeneralErrors: form.ValidationInfo.GetGeneralErrors(),
			FieldErrors:   mapFieldErrors(form.ValidationInfo),
		},
	}, nil

}

func (r *CommerceCartMutationResolver) CommerceCartAddCouponCode(ctx context.Context, couponCode string) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	_, err := r.cartService.ApplyAny(ctx, req.Session(), couponCode)

	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

//mapCommerce_Cart_BillingAddressForm - helper to map the graphql type Commerce_Cart_BillingAddressForm from common form
func mapCommerceBillingAddressForm(form *domain.Form, success bool) (*dto.BillingAddressForm, error) {
	billingFormData, ok := form.Data.(cartForms.BillingAddressForm)
	if !ok {
		return nil, errors.New("unexpected form data")
	}

	return &dto.BillingAddressForm{
		FormData:  billingFormData,
		Processed: success,
		ValidationInfo: dto.ValidationInfo{
			GeneralErrors: form.ValidationInfo.GetGeneralErrors(),
			FieldErrors:   mapFieldErrors(form.ValidationInfo),
		},
	}, nil
}

func mapFieldErrors(validationInfo domain.ValidationInfo) []dto.FieldError {
	var fieldErrors []dto.FieldError
	for fieldName, currentFieldErrors := range validationInfo.GetErrorsForAllFields() {
		for _, currentFieldError := range currentFieldErrors {
			fieldErrors = append(fieldErrors, dto.FieldError{
				MessageKey:   currentFieldError.MessageKey,
				DefaultLabel: currentFieldError.DefaultLabel,
				FieldName:    fieldName,
			})
		}
	}
	return fieldErrors
}
