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
		CustomerOrderService domain.CustomerOrderService `inject:""`
		AuthManager          *application.AuthManager    `inject:""`
	}
)

// Data controller for blocks
func (dc *DataControllerCustomerOrders) Data(c context.Context, r *web.Request) interface{} {
	auth, err := dc.AuthManager.Auth(c, r.Session().G())
	if err != nil {
		return nil
	}

	customerOrders, _ := dc.CustomerOrderService.Get(c, auth)

	return customerOrders
}
