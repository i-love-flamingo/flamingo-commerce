package controller

import (
	"go.aoe.com/flamingo/core/auth/application"
	"go.aoe.com/flamingo/core/order/domain"
	"go.aoe.com/flamingo/framework/web"
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
