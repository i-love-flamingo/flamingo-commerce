package graphql

import (
	"context"
	"errors"
	"net/url"

	formApplication "flamingo.me/form/application"
	"flamingo.me/form/domain"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/application"
)

// CommerceCartMutationResolver resolves cart mutations
type CommerceCartMutationResolver struct {
	q                            *CommerceCartQueryResolver
	cartService                  *application.CartService
	cartReceiverService          *application.CartReceiverService
	billingAddressFormController *forms.BillingAddressFormController
	deliveryFormController       *forms.DeliveryFormController
	simplePaymentFormController  *forms.SimplePaymentFormController
	formDataEncoderFactory       formApplication.FormDataEncoderFactory
}

// Inject dependencies
func (r *CommerceCartMutationResolver) Inject(q *CommerceCartQueryResolver,
	billingAddressFormController *forms.BillingAddressFormController,
	deliveryFormController *forms.DeliveryFormController,
	formDataEncoderFactory formApplication.FormDataEncoderFactory,
	simplePaymentFormController *forms.SimplePaymentFormController,
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
func (r *CommerceCartMutationResolver) CommerceAddToCart(ctx context.Context, graphqlAddRequest dto.AddToCart) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	if graphqlAddRequest.Qty < 0 {
		graphqlAddRequest.Qty = 0
	}

	addRequest := cartDomain.AddRequest{
		MarketplaceCode:        graphqlAddRequest.MarketplaceCode,
		Qty:                    graphqlAddRequest.Qty,
		VariantMarketplaceCode: graphqlAddRequest.VariantMarketplaceCode,
		AdditionalData:         nil,
		BundleConfiguration:    dto.MapBundleConfigToDomain(graphqlAddRequest.BundleConfiguration),
	}

	_, err := r.cartService.AddProduct(ctx, req.Session(), graphqlAddRequest.DeliveryCode, addRequest)
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

// CommerceUpdateItemBundleConfig mutation for updating item quantity
func (r *CommerceCartMutationResolver) CommerceUpdateItemBundleConfig(ctx context.Context, itemID string, bundleConfig []*dto.ChoiceConfiguration) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	var bundleConfigDto []dto.ChoiceConfiguration

	for _, config := range bundleConfig {
		bundleConfigDto = append(bundleConfigDto, *config)
	}

	bundleConfigDomain := dto.MapBundleConfigToDomain(bundleConfigDto)

	updateCommand := cartDomain.ItemUpdateCommand{
		ItemID:              itemID,
		BundleConfiguration: bundleConfigDomain,
	}

	err := r.cartService.UpdateItemBundleConfig(ctx, req.Session(), updateCommand)
	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceCartUpdateBillingAddress resolver method
func (r *CommerceCartMutationResolver) CommerceCartUpdateBillingAddress(ctx context.Context, address *forms.AddressForm) (*dto.BillingAddressForm, error) {
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

// CommerceCartUpdateSelectedPayment resolver method
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
func (r *CommerceCartMutationResolver) CommerceCartUpdateDeliveryAddresses(ctx context.Context, deliveryForms []*forms.DeliveryForm) ([]*dto.DeliveryAddressForm, error) {
	result := make([]*dto.DeliveryAddressForm, 0, len(deliveryForms))
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

		result = append(result, &deliveryAddressForm)
	}

	return result, nil
}

// CommerceCartUpdateDeliveryShippingOptions updates the method/carrier of one or multiple existing deliveries
func (r *CommerceCartMutationResolver) CommerceCartUpdateDeliveryShippingOptions(ctx context.Context, shippingOptions []*dto.DeliveryShippingOption) (*dto.UpdateShippingOptionsResult, error) {
	session := web.SessionFromContext(ctx)
	cart, err := r.cartReceiverService.ViewCart(ctx, session)
	if err != nil {
		return nil, err
	}

	for _, shippingOption := range shippingOptions {
		delivery, found := cart.GetDeliveryByCode(shippingOption.DeliveryCode)
		if !found {
			return nil, cartDomain.ErrDeliveryCodeNotFound
		}

		deliveryInfo := delivery.DeliveryInfo
		deliveryInfo.Carrier = shippingOption.Carrier
		deliveryInfo.Method = shippingOption.Method

		err = r.cartService.UpdateDeliveryInfo(ctx, session, shippingOption.DeliveryCode, cartDomain.CreateDeliveryInfoUpdateCommand(deliveryInfo))
		if err != nil {
			return nil, err
		}
	}

	return &dto.UpdateShippingOptionsResult{Processed: true}, nil
}

// CartClean clears users cart
func (r *CommerceCartMutationResolver) CartClean(ctx context.Context) (bool, error) {
	err := r.cartService.Clean(ctx, web.SessionFromContext(ctx))
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateAdditionalData of cart
func (r *CommerceCartMutationResolver) UpdateAdditionalData(ctx context.Context, additionalDataList []*dto.KeyValue) (*dto.DecoratedCart, error) {
	session := web.SessionFromContext(ctx)
	additionalDataMap := map[string]string{}
	for _, additionalData := range additionalDataList {
		additionalDataMap[additionalData.Key] = additionalData.Value
	}

	_, err := r.cartService.UpdateAdditionalData(ctx, session, additionalDataMap)
	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// UpdateDeliveriesAdditionalData of cart
func (r *CommerceCartMutationResolver) UpdateDeliveriesAdditionalData(ctx context.Context, additionalDataList []*dto.DeliveryAdditionalData) (*dto.DecoratedCart, error) {
	session := web.SessionFromContext(ctx)
	for _, additionalData := range additionalDataList {
		additionalDataMap := map[string]string{}
		for _, deliveryAdditionalData := range additionalData.AdditionalData {
			additionalDataMap[deliveryAdditionalData.Key] = deliveryAdditionalData.Value
		}

		_, err := r.cartService.UpdateDeliveryAdditionalData(ctx, session, additionalData.DeliveryCode, additionalDataMap)
		if err != nil {
			return nil, err
		}
	}

	return r.q.CommerceCart(ctx)
}

func mapCommerceDeliveryAddressForm(form *domain.Form, success bool) (dto.DeliveryAddressForm, error) {
	formData, ok := form.Data.(forms.DeliveryForm)
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
		DesiredTime:       formData.DesiredTime,
	}, nil
}

// mapCommerceBillingAddressForm helper to map the graphql type Commerce_Cart_BillingAddressForm from common form
func mapCommerceBillingAddressForm(form *domain.Form, success bool) (*dto.BillingAddressForm, error) {
	billingFormData, ok := form.Data.(forms.BillingAddressForm)
	if !ok {
		return nil, errors.New("unexpected form data")
	}

	return &dto.BillingAddressForm{
		FormData:  forms.AddressForm(billingFormData),
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
