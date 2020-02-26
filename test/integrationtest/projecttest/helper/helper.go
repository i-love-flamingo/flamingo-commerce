package helper

import (
	"flamingo.me/flamingo-commerce/v3/breadcrumbs"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/category"
	"flamingo.me/flamingo-commerce/v3/checkout"
	"flamingo.me/flamingo-commerce/v3/order"
	"flamingo.me/flamingo-commerce/v3/payment"
	"flamingo.me/flamingo-commerce/v3/price"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/search"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	projectTestGraphql "flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/graphql"
	"flamingo.me/flamingo-commerce/v3/w3cdatalayer"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/core/locale"
	auth "flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/core/security"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/web/filter"
	"flamingo.me/graphql"
	"fmt"
	"github.com/gavv/httpexpect"
	"strings"
	"testing"

	// "flamingo.me/redirects"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/core/robotstxt"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

//modulesDemoProject return slice of modules that we want to have in our example app for testing
func modulesDemoProject() []dingo.Module {
	return []dingo.Module{
		new(framework.InitModule),
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
		new(projectTestGraphql.Module),
		new(graphql.Module),
	}
}

//BootupDemoProject - boots up a complete demo project
func BootupDemoProject() integrationtest.BootupInfo {
	return integrationtest.Bootup(modulesDemoProject(), "../config", nil)
}

//GenerateGraphQL - generates the graphql interfaces for the demo project and saves to filesystem.
// use via makefile - each time you modify the schema
func GenerateGraphQL() {
	application, err := flamingo.NewApplication(modulesDemoProject(), flamingo.ConfigDir("../config"))
	if err != nil {
		panic(err)
	}

	servicesI, err := application.ConfigArea().Injector.GetInstance(new([]graphql.Service))
	if err != nil {
		panic(err)
	}
	services := servicesI.([]graphql.Service)
	err = graphql.Generate(services, "graphql", "graphql/schema")
	if err != nil {
		panic(err)
	}
}

//GraphQlQueryRequest - helper to get a initialised httpexpect request with the graphql query - used in tests
func GraphQlQueryRequest(t *testing.T, e *httpexpect.Expect, query string) *httpexpect.Request {
	query = strings.Replace(query, "\n", "", -1)
	query = strings.Replace(query, "\t", "", -1)
	query = strings.Replace(query, `"`, `\"`, -1)
	graphQlQuery := fmt.Sprintf(`{"variables":{},"query":"%v"}`, query)
	t.Log("GraphQlQueryRequest", graphQlQuery)
	return e.POST("/en/graphql").WithHeader("Content-Type", "application/json").WithBytes([]byte(graphQlQuery))
}
