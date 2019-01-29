package controller

import (
	"context"

	"flamingo.me/flamingo-commerce/order/domain"
	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/framework/web"
)

type (
	// DataControllerCustomerOrders for `get("orders.customerorders", ...)` requests
	DataControllerCustomerOrders struct {
		customerOrderService domain.CustomerOrderService
		authManager          *application.AuthManager
	}
)

// Inject dependencies
func (dc *DataControllerCustomerOrders) Inject(
	customerOrderService domain.CustomerOrderService,
	authManager *application.AuthManager,
) {
	dc.customerOrderService = customerOrderService
	dc.authManager = authManager
}

// Data controller for blocks
func (dc *DataControllerCustomerOrders) Data(c context.Context, r *web.Request) interface{} {
	auth, err := dc.authManager.Auth(c, r.Session().G())
	if err != nil {
		return nil
	}

	customerOrders, _ := dc.customerOrderService.Get(c, auth)

	return customerOrders
}
