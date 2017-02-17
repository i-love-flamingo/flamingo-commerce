package product

import (
	"flamingo/core/flamingo"
	"flamingo/core/product/controller"
)

func Register(r *flamingo.ServiceContainer) {
	r.Handle("product.view", new(controller.ViewController))
	r.Route("/product/{Uid}", "product.view")
}
