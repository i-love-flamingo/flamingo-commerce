// +build integration

package integrationtest

import (
	"context"
	"flamingo.me/flamingo/v3/framework/config"
	"sync"

	"flamingo.me/flamingo/v3"
	"log"
	"os"

	// "flamingo.me/redirects"
	"flamingo.me/dingo"

	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	eventReceiver struct{}
	testmodule    struct {
		eventRouter flamingoFramework.EventRouter
	}
)

//Side effect vars to get status and exchange stuff with the testmodule
var rw sync.Mutex
var bootupReady chan struct{}
var testmoduleInstanceInApp *testmodule
var additionalConfig config.Map

func init() {
	bootupReady = make(chan struct{})
}

//Configure for your testmodule in the app
func (t *testmodule) Inject(eventRouter flamingoFramework.EventRouter) {
	t.eventRouter = eventRouter
	testmoduleInstanceInApp = t
}

//Configure for your testmodule in the app
func (t *testmodule) Configure(i *dingo.Injector) {
	flamingoFramework.BindEventSubscriber(i).To(t)
}

//Notify gets notified by event router
func (t *testmodule) Notify(ctx context.Context, event flamingoFramework.Event) {
	switch event.(type) {
	case *flamingoFramework.ServerStartEvent:
		log.Printf("ServerStartEvent event received...")
		bootupReady <- struct{}{}
	case *flamingoFramework.ShutdownEvent:
		log.Printf("ServerShutdownEvent event received...")
	}
}

// DefaultConfig enables inMemory cart service adapter
func (t *testmodule) DefaultConfig() config.Map {
	return additionalConfig
}

//WaitForStart until
func (t *testmodule) SendShutdown() {
	log.Printf("Trigger ServerShutdownEvent...")
	t.eventRouter.Dispatch(context.Background(), flamingoFramework.ShutdownEvent{})
}

//Bootup flamingo app with the given modules (and the config in folder given )
func Bootup(modules []dingo.Module, configDir string, config config.Map) (func(), string) {
	rw.Lock()
	defer rw.Unlock()
	//add testmodul that listens
	modules = append(modules, new(testmodule))
	//rootArea := rootArea("config")
	os.Args[1] = "serve"
	additionalConfig = config
	go flamingo.App(modules, flamingo.ConfigDir(configDir))
	<-bootupReady
	tm := testmoduleInstanceInApp
	//TODO - use new port on every bootup
	return func() { tm.SendShutdown() }, "localhost:3210"
}
