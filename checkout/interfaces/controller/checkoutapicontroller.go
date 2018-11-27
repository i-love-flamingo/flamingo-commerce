package controller

import (
	"context"

	formApplicationService "flamingo.me/flamingo/core/form/application"
	formDomain "flamingo.me/flamingo/core/form/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (
	// CheckoutAPIController for onepage checkout
	CheckoutAPIController struct {
		responder   *web.Responder
		logger      flamingo.Logger
		formService formDomain.FormService
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
func (cac *CheckoutAPIController) Inject(r *web.Responder, l flamingo.Logger, fs formDomain.FormService) {
	cac.responder = r
	cac.logger = l
	cac.formService = fs
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

	return &web.JSONResponse{
		Data: JSONResult{
			Message:     "success",
			MessageCode: "200",
			Success:     true,
		},
	}
}
