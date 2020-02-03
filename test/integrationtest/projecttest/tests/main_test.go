// +build integration

package tests

import (
	"os"
	"testing"

	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
)

var FlamingoURL string

// TestMain used for setup and teardown
func TestMain(m *testing.M) {
	info := helper.BootupDemoProject()
	FlamingoURL = info.BaseURL
	result := m.Run()
	info.ShutdownFunc()
	os.Exit(result)
}
