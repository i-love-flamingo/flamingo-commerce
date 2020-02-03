// +build integration

package graphql_test

import (
	"flag"
	"os"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

var FlamingoURL string

// TestMain used for setup and teardown
func TestMain(m *testing.M) {
	flag.Parse()
	testing.Verbose()
	info := helper.BootupDemoProject("../../config/")
	FlamingoURL = info.BaseURL
	result := m.Run()
	info.ShutdownFunc()
	os.Exit(result)
}
