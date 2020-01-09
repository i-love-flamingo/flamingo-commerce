// +build integration

package tests

import (
	"flamingo.me/flamingo-commerce/v3/test/integrationtest/projecttest/helper"
	"os"
	"testing"
	//"time"
)

var FlamingoUrl string

//TestMain - golang TestMain - used for setup and teardown
func TestMain(m *testing.M) {
	info := helper.BootupDemoProject()
	FlamingoUrl = info.BaseUrl
	result := m.Run()
	info.ShutdownFunc()
	//time.Sleep(100*time.Second)
	// call flag.Parse() here if TestMain uses flags
	os.Exit(result)

}
