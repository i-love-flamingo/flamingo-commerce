package controller

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/order/domain"
	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// DataControllerCustomerOrders for `get("orders.customerorders", ...)` requests
	DataControllerCustomerOrders struct {
		customerOrderService domain.CustomerOrderService
		authManager          *authApplication.AuthManager
	}
)

// Inject dependencies
func (dc *DataControllerCustomerOrders) Inject(
	customerOrderService domain.CustomerOrderService,
	authManager *authApplication.AuthManager,
) {
	dc.customerOrderService = customerOrderService
	dc.authManager = authManager
}

// Data controller for blocks
func (dc *DataControllerCustomerOrders) Data(c context.Context, r *web.Request, _ web.RequestParams) interface{} {
	auth, err := dc.authManager.Auth(c, r.Session())
	if err != nil {
		return nil
	}

	customerOrders, _ := dc.customerOrderService.Get(c, auth)

	return customerOrders
}
