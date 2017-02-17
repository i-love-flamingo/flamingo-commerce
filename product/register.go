package product

import (
	"flamingo/core/app"
	"flamingo/core/product/controller"
)

func Register(r *app.ServiceContainer) {
	r.Handle("product.view", new(controller.ViewController))
	r.Route("/product/{Uid}", "product.view")
}
