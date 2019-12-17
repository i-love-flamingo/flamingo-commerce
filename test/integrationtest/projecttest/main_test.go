// +build integration

package projecttest

import (
	"os"
	"testing"
	//"time"
)

//TestMain - golang TestMain - used for setup and teardown
func TestMain(m *testing.M) {
	deferfunc, _ := BootupDemoProject()
	result := m.Run()
	deferfunc()
	//time.Sleep(100*time.Second)
	// call flag.Parse() here if TestMain uses flags
	os.Exit(result)
}
