package product

import (
	di "flamingo/framework/dependencyinjection"
	"flamingo/framework/router"
	"flamingo/core/product/controller"
)

// Register Services for product package
func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		r.Handle("product.view", new(controller.ViewController))
		r.Route("/product/{uid}", "product.view")
	}, router.RouterRegister)
}
