package controller

import (
	"context"

	"flamingo.me/flamingo-commerce/cart/application"
	formApplicationService "flamingo.me/flamingo/core/form/application"
	formDomain "flamingo.me/flamingo/core/form/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	...CheckoutDomain ".../.../flamingo/src/checkout/domain"
)

type (
	// CheckoutAPIController for onepage checkout
	CheckoutAPIController struct {
		responder           *web.Responder
		logger              flamingo.Logger
		formService         formDomain.FormService
		cartReceiverService *application.CartReceiverService
	}

	// JSONResult for ajax response
	JSONResult struct {
		Message     string
		MessageCode string
		Success     bool
		FieldErrors map[string][]formDomain.Error
	}
)

// Inject required dependencies
func (cac *CheckoutAPIController) Inject(
	r *web.Responder,
	l flamingo.Logger,
	fs formDomain.FormService,
	crs *application.CartReceiverService,
) {
	cac.responder = r
	cac.logger = l
	cac.formService = fs
	cac.cartReceiverService = crs
}

// SubmitBillingAddressAction - save or update the shipping address
func (cac *CheckoutAPIController) SubmitBillingAddressAction(ctx context.Context, r *web.Request) web.Response {

	form, err := formApplicationService.ProcessFormRequest(ctx, r, cac.formService)
	if err != nil {
		cac.logger.Error("Error processing form ", form)
	}

	if !form.IsValidAndSubmitted() {

		return &web.JSONResponse{
			BasicResponse: web.BasicResponse{
				Status: 400,
			},
			Data: JSONResult{
				Message:     "error parsing one ore more fields",
				MessageCode: "form.invalid",
				Success:     false,
				FieldErrors: form.ValidationInfo.FieldErrors,
			},
		}

	}

	// submit valid data
	if billingFormData, ok := form.Data.(...CheckoutDomain.BillingAddressForm); ok {
		billingAddress := billingFormData.MapBillingAddress(billingFormData.AddressFormData)

		cart, cartOrderBehaviour, err := cac.cartReceiverService.GetCart(ctx, r.Session().G())
		if err != nil {
			cac.logger.Error("no cart to update")
			return &web.JSONResponse{
				BasicResponse: web.BasicResponse{
					Status: 400,
				},
				Data: JSONResult{
					Message:     "error saving information to cart",
					MessageCode: "submit.error",
					Success:     false,
					FieldErrors: nil,
				},
			}
		}

		_, cartUpdateErr := cartOrderBehaviour.UpdateBillingAddress(ctx, cart, billingAddress)
		if cartUpdateErr != nil {
			cac.logger.Error(cartUpdateErr)
			return &web.JSONResponse{
				BasicResponse: web.BasicResponse{
					Status: 400,
				},
				Data: JSONResult{
					Message:     "error saving information to cart",
					MessageCode: "submit.error",
					Success:     false,
					FieldErrors: nil,
				},
			}
		}

	}

	return &web.JSONResponse{
		Data: JSONResult{
			Message:     "success",
			MessageCode: "200",
			Success:     true,
		},
	}
}
