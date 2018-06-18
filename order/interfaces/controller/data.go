package controller

import (
	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo-commerce/order/domain"
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
func (dc *DataControllerCustomerOrders) Data(c web.Context) interface{} {
	auth, err := dc.AuthManager.Auth(c)
	if err != nil {
		return nil
	}

	customerOrders, _ := dc.CustomerOrderService.Get(c, auth)

	return customerOrders
}
