package product

import (
	"flamingo/core/flamingo/service_container"
	"flamingo/core/product/controller"
)

func Register(r *service_container.ServiceContainer) {
	r.Handle("product.view", new(controller.ViewController))
	r.Route("/product/{Uid}", "product.view")
}
