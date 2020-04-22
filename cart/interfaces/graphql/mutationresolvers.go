package graphql

import (
	"context"
	"errors"
	"net/url"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	cartForms "flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	formApplication "flamingo.me/form/application"
	"flamingo.me/form/domain"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCartMutationResolver resolves cart mutations
type CommerceCartMutationResolver struct {
	q                            *CommerceCartQueryResolver
	cartService                  *application.CartService
	cartReceiverService          *application.CartReceiverService
	billingAddressFormController *cartForms.BillingAddressFormController
	deliveryFormController       *cartForms.DeliveryFormController
	simplePaymentFormController  *cartForms.SimplePaymentFormController
	formDataEncoderFactory       formApplication.FormDataEncoderFactory
}

// Inject dependencies
func (r *CommerceCartMutationResolver) Inject(q *CommerceCartQueryResolver,
	billingAddressFormController *cartForms.BillingAddressFormController,
	deliveryFormController *cartForms.DeliveryFormController,
	formDataEncoderFactory formApplication.FormDataEncoderFactory,
	simplePaymentFormController *cartForms.SimplePaymentFormController,
	cartService *application.CartService,
	cartReceiverService *application.CartReceiverService) *CommerceCartMutationResolver {
	r.q = q
	r.billingAddressFormController = billingAddressFormController
	r.deliveryFormController = deliveryFormController
	r.formDataEncoderFactory = formDataEncoderFactory
	r.simplePaymentFormController = simplePaymentFormController
	r.cartService = cartService
	r.cartReceiverService = cartReceiverService
	return r
}

// CommerceAddToCart mutation for adding products to the current users cart
func (r *CommerceCartMutationResolver) CommerceAddToCart(ctx context.Context, marketplaceCode string, qty int, deliveryCode string) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	addRequest := r.cartService.BuildAddRequest(ctx, marketplaceCode, "", qty, nil)

	_, err := r.cartService.AddProduct(ctx, req.Session(), deliveryCode, addRequest)
	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceDeleteItem resolver
func (r *CommerceCartMutationResolver) CommerceDeleteItem(ctx context.Context, itemID string, deliveryCode string) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	err := r.cartService.DeleteItem(ctx, req.Session(), itemID, deliveryCode)

	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceDeleteCartDelivery mutation for removing deliveries from current users cart
func (r *CommerceCartMutationResolver) CommerceDeleteCartDelivery(ctx context.Context, deliveryCode string) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	_, err := r.cartService.DeleteDelivery(ctx, req.Session(), deliveryCode)
	if err != nil {
		return nil, err
	}
	return r.q.CommerceCart(ctx)
}

// CommerceUpdateItemQty mutation for updating item quantity
func (r *CommerceCartMutationResolver) CommerceUpdateItemQty(ctx context.Context, itemID string, deliveryCode string, qty int) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	err := r.cartService.UpdateItemQty(ctx, req.Session(), itemID, deliveryCode, qty)
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

// CommerceCartApplyCouponCodeOrGiftCard – apply coupon code or gift card
func (r *CommerceCartMutationResolver) CommerceCartApplyCouponCodeOrGiftCard(ctx context.Context, code string) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	_, err := r.cartService.ApplyAny(ctx, req.Session(), code)

	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceCartRemoveCouponCode - remove coupon code
func (r *CommerceCartMutationResolver) CommerceCartRemoveCouponCode(ctx context.Context, couponCode string) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	_, err := r.cartService.RemoveVoucher(ctx, req.Session(), couponCode)

	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceCartRemoveGiftCard - remove gift card
func (r *CommerceCartMutationResolver) CommerceCartRemoveGiftCard(ctx context.Context, giftCardCode string) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	_, err := r.cartService.RemoveGiftCard(ctx, req.Session(), giftCardCode)

	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceCartUpdateDeliveryAddresses can be used to update or set one or multiple delivery addresses, uses the delivery form controller
func (r *CommerceCartMutationResolver) CommerceCartUpdateDeliveryAddresses(ctx context.Context, deliveryForms []*forms.DeliveryForm) (*dto.DeliveryAddressForms, error) {
	result := &dto.DeliveryAddressForms{}
	request := web.CreateRequest(web.RequestFromContext(ctx).Request(), web.SessionFromContext(ctx))
	for _, deliveryForm := range deliveryForms {
		encodedForm, err := r.formDataEncoderFactory.CreateByNamedEncoder("commerce.cart.billingFormService").Encode(ctx, deliveryForm)
		if err != nil {
			return nil, err
		}
		request.Request().Form = encodedForm
		request.Params["deliveryCode"] = deliveryForm.LocationCode
		form, success, err := r.deliveryFormController.HandleFormAction(ctx, request)
		if err != nil {
			return nil, err
		}

		deliveryAddressForm, err := mapCommerceDeliveryAddressForm(form, success)
		if err != nil {
			return nil, err
		}

		result.Forms = append(result.Forms, deliveryAddressForm)
	}

	return result, nil
}

// CommerceCartUpdateDeliveryShippingOptions updates the method/carrier of one or multiple existing deliveries
func (r *CommerceCartMutationResolver) CommerceCartUpdateDeliveryShippingOptions(ctx context.Context, shippingOptions []*dto.DeliveryShippingOption) (*dto.DeliveryAddressForms, error) {
	result := &dto.DeliveryAddressForms{}

	session := web.SessionFromContext(ctx)
	cart, err := r.cartReceiverService.ViewCart(ctx, session)
	if err != nil {
		return nil, err
	}

	for _, shippingOption := range shippingOptions {
		formResult := dto.DeliveryAddressForm{
			Processed:    false,
			DeliveryCode: shippingOption.DeliveryCode,
			Method:       shippingOption.Method,
			Carrier:      shippingOption.Carrier,
		}

		delivery, found := cart.GetDeliveryByCode(shippingOption.DeliveryCode)
		if !found {
			formResult.ValidationInfo.GeneralErrors = []domain.Error{
				{
					MessageKey: "delivery_not_found",
				},
			}
			result.Forms = append(result.Forms, formResult)
			continue
		}

		deliveryInfo := delivery.DeliveryInfo
		deliveryInfo.Carrier = shippingOption.Carrier
		deliveryInfo.Method = shippingOption.Method

		err = r.cartService.UpdateDeliveryInfo(ctx, session, shippingOption.DeliveryCode, cartDomain.CreateDeliveryInfoUpdateCommand(deliveryInfo))
		if err != nil {
			return nil, err
		}
		formResult.Processed = true
		result.Forms = append(result.Forms, formResult)
	}

	return result, nil
}

func mapCommerceDeliveryAddressForm(form *domain.Form, success bool) (dto.DeliveryAddressForm, error) {
	formData, ok := form.Data.(cartForms.DeliveryForm)
	if !ok {
		return dto.DeliveryAddressForm{}, errors.New("unexpected form data")
	}

	return dto.DeliveryAddressForm{
		FormData:  formData.DeliveryAddress,
		Processed: success,
		ValidationInfo: dto.ValidationInfo{
			GeneralErrors: form.ValidationInfo.GetGeneralErrors(),
			FieldErrors:   mapFieldErrors(form.ValidationInfo),
		},
		UseBillingAddress: formData.UseBillingAddress,
		DeliveryCode:      formData.LocationCode,
		Method:            formData.ShippingMethod,
		Carrier:           formData.ShippingCarrier,
	}, nil
}

// mapCommerceBillingAddressForm helper to map the graphql type Commerce_Cart_BillingAddressForm from common form
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
