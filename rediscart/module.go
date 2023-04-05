package rediscart

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart"
	cartInfrastructure "flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"flamingo.me/flamingo-commerce/v3/rediscart/infrastructure"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
)

type (
	// Module for a cart storage using redis
	Module struct{}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(infrastructure.CartSerializer)).To(new(infrastructure.GobSerializer))
	injector.Override(new(cartInfrastructure.CartStorage), "").To(new(infrastructure.RedisStorage)).AsEagerSingleton()
	injector.BindMap(new(healthcheck.Status), "cart.storage.redis").To(new(infrastructure.RedisStorage))
}

// CueConfig defines the cart module configuration
func (*Module) CueConfig() string {
	return `
commerce: {
	cart: {
		redis: {
			keyPrefix: string | *"cart:"
			ttl: {
				guest: string | *"48h"
				customer: string | *"168h"
			}
			address: string | *""
			network: "unix" | *"tcp"
			password: string | *""
			idleConnections: number | *10
			database: float | int | *0
			tls: bool | *false
		}
	}
}`
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(cart.Module),
	}
}
