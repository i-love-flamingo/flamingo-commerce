package controller

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/order/domain"
)

type (
	// DataControllerCustomerOrders for `get("orders.customerorders", ...)` requests
	DataControllerCustomerOrders struct {
		customerIdentityOrderService domain.CustomerIdentityOrderService
		webIdentityService           *auth.WebIdentityService
	}
)

// Inject dependencies
func (dc *DataControllerCustomerOrders) Inject(
	customerIdentityOrderService domain.CustomerIdentityOrderService,
	webIdentityService *auth.WebIdentityService,
) {
	dc.customerIdentityOrderService = customerIdentityOrderService
	dc.webIdentityService = webIdentityService
}

// Data controller for blocks
func (dc *DataControllerCustomerOrders) Data(ctx context.Context, r *web.Request, _ web.RequestParams) interface{} {
	identity := dc.webIdentityService.Identify(ctx, r)
	if identity == nil {
		return nil
	}

	customerOrders, _ := dc.customerIdentityOrderService.Get(ctx, identity)

	return customerOrders
}
