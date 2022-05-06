package forms

import (
	"context"
	"errors"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/application"
	"flamingo.me/form/domain"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// DeliveryForm the form for billing address
	DeliveryForm struct {
		DeliveryAddress AddressForm `form:"deliveryAddress"`
		// UseBillingAddress - the address should be taken from billing (only relevant for type address)
		UseBillingAddress bool      `form:"useBillingAddress"`
		ShippingMethod    string    `form:"shippingMethod"`
		ShippingCarrier   string    `form:"shippingCarrier"`
		LocationCode      string    `form:"locationCode"`
		DesiredTime       time.Time `form:"desiredTime"`
	}

	// DeliveryFormService implements Form(Data)Provider interface of form package
	DeliveryFormService struct {
		applicationCartReceiverService *cartApplication.CartReceiverService
	}

	// DeliveryFormController the (mini) MVC
	DeliveryFormController struct {
		responder                      *web.Responder
		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService
		logger                         flamingo.Logger
		formHandlerFactory             application.FormHandlerFactory
		billingAddressFormProvider     *BillingAddressFormService
	}
)

// MapToDeliveryInfo - updates some fields of the given DeliveryInfo with data from the form
func (d *DeliveryForm) MapToDeliveryInfo(currentInfo cartDomain.DeliveryInfo) cartDomain.DeliveryInfo {
	address := d.DeliveryAddress.MapToDomainAddress()
	currentInfo.DeliveryLocation.Address = &address
	currentInfo.DeliveryLocation.UseBillingAddress = d.UseBillingAddress
	currentInfo.DeliveryLocation.Code = d.LocationCode
	currentInfo.Method = d.ShippingMethod
	currentInfo.Carrier = d.ShippingCarrier
	currentInfo.DesiredTime = d.DesiredTime
	return currentInfo
}

// Inject - Inject
func (p *DeliveryFormService) Inject(applicationCartReceiverService *cartApplication.CartReceiverService) {
	p.applicationCartReceiverService = applicationCartReceiverService
}

// GetFormData from data provider
func (p *DeliveryFormService) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {

	cart, err := p.applicationCartReceiverService.ViewCart(ctx, req.Session())
	useBilling := false
	method := ""
	carrier := ""
	locationCode := ""
	deliveryAddress := AddressForm{}
	deliverycode := req.Params["deliveryCode"]
	if deliverycode != "" && err == nil {
		if delivery, found := cart.GetDeliveryByCode(deliverycode); found {
			if delivery.DeliveryInfo.DeliveryLocation.Address != nil {
				deliveryAddress.LoadFromCartAddress(*delivery.DeliveryInfo.DeliveryLocation.Address)
			}
			useBilling = delivery.DeliveryInfo.DeliveryLocation.UseBillingAddress
			method = delivery.DeliveryInfo.Method
			carrier = delivery.DeliveryInfo.Carrier
			locationCode = delivery.DeliveryInfo.DeliveryLocation.Code
		}
	}

	return DeliveryForm{
		DeliveryAddress:   deliveryAddress,
		UseBillingAddress: useBilling,
		ShippingMethod:    method,
		ShippingCarrier:   carrier,
		LocationCode:      locationCode,
	}, nil
}

// Validate form service
func (p *DeliveryFormService) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
	deliveryForm, ok := formData.(DeliveryForm)
	if !ok {
		return nil, errors.New("no BillingAddressForm given")
	}
	validationInfo := domain.ValidationInfo{}
	if !deliveryForm.UseBillingAddress {
		// Validate address only if no billing should be used
		validationInfo = validatorProvider.Validate(ctx, req, deliveryForm)
	}
	return &validationInfo, nil
}

// Inject dependencies
func (c *DeliveryFormController) Inject(responder *web.Responder,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	logger flamingo.Logger,
	formHandlerFactory application.FormHandlerFactory,
	billingAddressFormProvider *BillingAddressFormService) {
	c.responder = responder
	c.applicationCartReceiverService = applicationCartReceiverService
	c.applicationCartService = applicationCartService
	c.formHandlerFactory = formHandlerFactory
	c.logger = logger.WithField(flamingo.LogKeyModule, "cart").WithField(flamingo.LogKeyCategory, "deliveryform")
	c.billingAddressFormProvider = billingAddressFormProvider
}

// GetUnsubmittedForm returns the form with deliveryform data - without validation
func (c *DeliveryFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*domain.Form, error) {
	formHandler, err := c.getFormHandler()
	if err != nil {
		return nil, err
	}
	return formHandler.HandleUnsubmittedForm(ctx, r)
}

func (c *DeliveryFormController) getFormHandler() (domain.FormHandler, error) {
	// ##  Handle the submitted form (validation etc)
	formHandlerBuilder := c.formHandlerFactory.GetFormHandlerBuilder()
	err := formHandlerBuilder.SetNamedFormService("commerce.cart.deliveryFormService")
	if err != nil {
		return nil, err
	}
	return formHandlerBuilder.Build(), nil
}

// HandleFormAction handles submitted form and saves to cart
func (c *DeliveryFormController) HandleFormAction(ctx context.Context, r *web.Request) (form *domain.Form, actionSuccessful bool, err error) {
	session := web.SessionFromContext(ctx)

	deliverycode := r.Params["deliveryCode"]
	if deliverycode == "" {
		return nil, false, errors.New("no deliverycode parameter given")
	}
	formHandler, err := c.getFormHandler()
	if err != nil {
		return nil, false, err
	}
	form, err = formHandler.HandleSubmittedForm(ctx, r)
	if err != nil {
		return nil, false, err
	}
	deliveryForm, ok := form.Data.(DeliveryForm)
	if !ok {
		return form, false, errors.New("cannot convert to deliveryForm ")
	}
	if !form.IsValidAndSubmitted() {
		return form, false, nil
	}

	cart, err := c.applicationCartReceiverService.ViewCart(ctx, session)
	if err != nil {
		return form, false, err
	}
	var deliveryInfo cartDomain.DeliveryInfo
	delivery, found := cart.GetDeliveryByCode(deliverycode)
	if !found {
		initialDelIfno, err := c.applicationCartService.GetInitialDelivery(deliverycode)
		if err != nil {
			return form, false, err
		}
		deliveryInfo = *initialDelIfno
	} else {
		deliveryInfo = delivery.DeliveryInfo
	}

	deliveryInfo = deliveryForm.MapToDeliveryInfo(deliveryInfo)

	// update Cart
	err = c.applicationCartService.UpdateDeliveryInfo(ctx, session, deliverycode, cartDomain.CreateDeliveryInfoUpdateCommand(deliveryInfo))
	if err != nil {
		c.logger.WithContext(ctx).Error("UpdateDeliveryInfo  Error %v", err)
		return form, false, err
	}
	return form, true, nil
}
