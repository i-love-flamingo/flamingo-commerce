package helper

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"flamingo.me/graphql"
	"github.com/gavv/httpexpect/v2"

	fakeCustomer "flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/customer"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/auth"
	fakeAuth "flamingo.me/flamingo/v3/core/auth/fake"
	"flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/core/locale"
	"flamingo.me/flamingo/v3/core/security"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/web/filter"

	"flamingo.me/flamingo-commerce/v3/customer"
	integrationCart "flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/cart"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/payment"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/modules/placeorder"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/core/robotstxt"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/prefixrouter"

	"flamingo.me/flamingo-commerce/v3/breadcrumbs"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/category"
	"flamingo.me/flamingo-commerce/v3/checkout"
	"flamingo.me/flamingo-commerce/v3/order"
	"flamingo.me/flamingo-commerce/v3/price"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/search"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	projectTestGraphql "flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/graphql"
	"flamingo.me/flamingo-commerce/v3/w3cdatalayer"
)

// modulesDemoProject return slice of modules that we want to have in our example app for testing
func modulesDemoProject() []dingo.Module {
	return []dingo.Module{
		new(framework.InitModule),
		new(cmd.Module),
		new(zap.Module),
		new(flamingoFramework.SessionModule),
		new(prefixrouter.Module),
		new(product.Module),
		new(locale.Module),
		new(customer.Module),
		new(fakeCustomer.Module),
		new(cart.Module),
		new(checkout.Module),
		new(search.Module),
		new(category.Module),
		new(requestlogger.Module),
		new(filter.DefaultCacheStrategyModule),
		new(auth.WebModule),
		new(fakeAuth.Module),
		new(breadcrumbs.Module),
		new(order.Module),
		new(healthcheck.Module),
		new(w3cdatalayer.Module),
		new(robotstxt.Module),
		new(security.Module),
		new(opencensus.Module),
		new(price.Module),
		new(projectTestGraphql.Module),
		new(graphql.Module),
		new(payment.Module),
		new(placeorder.Module),
		new(integrationCart.Module),
	}
}

// BootupDemoProject boots up a complete demo project
func BootupDemoProject(configDir string) integrationtest.BootupInfo {
	return integrationtest.Bootup(modulesDemoProject(), configDir, nil)
}

// GenerateGraphQL generates the graphql interfaces for the demo project and saves to filesystem.
// use via makefile - each time you modify the schema
func GenerateGraphQL() {
	application, err := flamingo.NewApplication(modulesDemoProject(), flamingo.ConfigDir("config"))
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

// GraphQlRequest helper to get a initialised httpexpect request with the graphql query - used in tests
func GraphQlRequest(t *testing.T, e *httpexpect.Expect, query string) *httpexpect.Request {
	t.Helper()
	query = strings.Replace(query, "\n", "", -1)
	query = strings.Replace(query, "\t", "", -1)
	query = strings.Replace(query, `"`, `\"`, -1)
	graphQlQuery := fmt.Sprintf(`{"variables":{},"query":"%v"}`, query)

	return e.POST("/en/graphql").WithHeader("Content-Type", "application/json").WithBytes([]byte(graphQlQuery))
}

// AsyncCheckWithTimeout calls fn every 10ms until timeout reached or fn doesn't return err anymore
func AsyncCheckWithTimeout(t *testing.T, timeoutAfter time.Duration, fn func() error) {
	t.Helper()
	timeout := time.NewTimer(timeoutAfter)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		err := fn()

		if err == nil {
			return
		}

		select {
		case <-timeout.C:
			t.Fatal(err)
		case <-ticker.C:
		}
	}
}
