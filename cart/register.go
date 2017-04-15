package cart

import (
	di "flamingo/framework/dependencyinjection"
	"flamingo/framework/router"
	"flamingo/core/cart/interfaces/controller"
	"flamingo/core/cart/infrastructure"
	"fmt"
	"flamingo/core/cart/application"
)

// Register adds handles for cms page routes.
func Register(c *di.Container) {
	fmt.Println("is called now 1")
	c.Register(func(r *router.Router) {
		fmt.Println("is called now 2")
		// Cart View
		r.Handle("cart.view", new(controller.CartViewController))
		r.Route("/cart", "cart.view")

		// Updating and Adding Items
		r.Handle("cart.api", new(controller.CartApiController))
		r.Route("/api/cart", "cart.api")  				// Implements - Get JSON (Get JSON)  And Put+Post= (Update Cart)

		r.Handle("cart.item.add.api", new(controller.CartItemAddApiController))
		r.Route("/api/cart/item/add", "cart.item.add.api")			// Post (Add)


		//For test
		r.Handle("logintest", new(controller.TestLoginController))
		r.Route("/logintest", "logintest")

	}, "router.register")

	// Basti - that did not work:
	 //c.RegisterFactory(infrastructure.FakecartrepositoryFactory)
	c.Register(new(infrastructure.Fakecartrepository))

	//Basti - also application.Cartservice is injected even without registration?
	c.Register(new(application.Cartservice))


	// Better way to discuss with Basti
	// Would like to have something where "Subscriber" can register and then get executed after DI is ready? But it need to belong to "event2" package somehow
	eventOrchestrator := new(application.EventOrchestration)
	c.Resolve(eventOrchestrator)
	eventOrchestrator.AddSubscriptions()

}
