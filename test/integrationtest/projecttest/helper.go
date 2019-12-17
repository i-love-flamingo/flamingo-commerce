// +build integration

package projecttest

import (
	"flamingo.me/flamingo-commerce/v3/breadcrumbs"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/category"
	"flamingo.me/flamingo-commerce/v3/checkout"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"

	"flamingo.me/flamingo-commerce/v3/order"
	"flamingo.me/flamingo-commerce/v3/payment"
	"flamingo.me/flamingo-commerce/v3/price"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/search"
	"flamingo.me/flamingo-commerce/v3/w3cdatalayer"
	"flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/core/locale"
	auth "flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/core/security"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/web/filter"

	// "flamingo.me/redirects"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/core/robotstxt"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/config"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

//modulesDemoProject return slice of modules that we want to have in our example app for testing
func modulesDemoProject() []dingo.Module {
	return []dingo.Module{
		new(framework.InitModule),
		new(config.Flags),
		new(cmd.Module),
		new(zap.Module),
		new(flamingoFramework.SessionModule),
		new(prefixrouter.Module),
		new(product.Module),
		new(locale.Module),
		new(cart.Module),
		new(checkout.Module),
		new(search.Module),
		new(category.Module),
		new(requestlogger.Module),
		new(filter.DefaultCacheStrategyModule),
		new(auth.Module),
		new(breadcrumbs.Module),
		new(order.Module),
		new(healthcheck.Module),
		new(w3cdatalayer.Module),
		new(robotstxt.Module),
		new(security.Module),
		new(opencensus.Module),
		new(price.Module),
		new(payment.Module),
	}
}

//BootupDemoProject - boots up a complete demo project
func BootupDemoProject() (func(), string) {
	return integrationtest.Bootup(modulesDemoProject(), "config", nil)
}
