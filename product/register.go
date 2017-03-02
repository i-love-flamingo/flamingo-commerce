package product

import (
	di "flamingo/core/flamingo/dependencyinjection"
	"flamingo/core/flamingo/router"
	"flamingo/core/product/controller"
)

// Register Services for product package
func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		r.Handle("product.view", new(controller.ViewController))
		r.Route("/product/{Uid}", "product.view")
	}, router.RouterRegister)
}
