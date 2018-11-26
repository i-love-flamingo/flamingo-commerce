package controller

import (
	"context"

	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (
	// CheckoutAPIController for onepage checkout
	CheckoutAPIController struct {
		responder *web.Responder
		logger    flamingo.Logger
	}

	// JSONResult for ajax response
	JSONResult struct {
		Message     string
		MessageCode string
		Success     bool
	}
)

// Inject required dependencies
func (cac *CheckoutAPIController) Inject(r *web.Responder, l flamingo.Logger) {
	cac.responder = r
	cac.logger = l
}

// UpdateShippingAddressAction - save or update the shipping address
func (cac *CheckoutAPIController) UpdateShippingAddressAction(ctx context.Context, r *web.Request) web.Response {

	return &web.JSONResponse{
		Data: JSONResult{
			Message:     "success",
			MessageCode: "200",
			Success:     true,
		},
	}
}
