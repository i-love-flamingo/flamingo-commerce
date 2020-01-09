package integrationtest

import (
	"context"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"fmt"
	"net/http"
	"sync"
	"time"

	"log"
	"os"

	"flamingo.me/flamingo/v3"

	// "flamingo.me/redirects"
	"flamingo.me/dingo"

	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	eventReceiver struct{}
	testmodule    struct {
		eventRouter flamingoFramework.EventRouter
		router      *web.Router
		server      *http.Server
	}

	//BootupInfo about the booted app
	BootupInfo struct {
		ShutdownFunc func()
		Application  *flamingo.Application
		BaseURL      string
		Running      chan struct{}
	}
)

//Side effect vars to get status and exchange stuff with the testmodule
var rw sync.Mutex

var additionalConfig config.Map
var lastPort = 9999

//Configure for your testmodule in the app
func (t *testmodule) Inject(eventRouter flamingoFramework.EventRouter,
	router *web.Router) {
	t.eventRouter = eventRouter
	t.router = router
}

//Configure for your testmodule in the app
func (t *testmodule) Configure(i *dingo.Injector) {
	flamingoFramework.BindEventSubscriber(i).To(t)
}

//Notify gets notified by event router
func (t *testmodule) Notify(ctx context.Context, event flamingoFramework.Event) {
	switch event.(type) {
	case *flamingoFramework.ShutdownEvent:
		log.Printf("ShutdownEvent event received...")
	}
}

// DefaultConfig enables inMemory cart service adapter
func (t *testmodule) DefaultConfig() config.Map {
	return additionalConfig
}

//shutdown until
func (t *testmodule) shutdownServer() {
	log.Printf("Trigger ServerShutdownEvent...")
	t.eventRouter.Dispatch(context.Background(), &flamingoFramework.ServerShutdownEvent{})
	_ = t.server.Shutdown(context.Background())
}

func (t *testmodule) nextServerPort() string {
	lastPort++
	return fmt.Sprintf("%v", lastPort)
}

//returns the port or error
func (t *testmodule) startServer(listenAndServeQuited chan struct{}) (string, error) {

	t.eventRouter.Dispatch(context.Background(), &flamingoFramework.ServerStartEvent{})
	t.server = &http.Server{
		Addr: ":" + t.nextServerPort(),
	}
	log.Printf("startServer on port %v", t.server.Addr)
	t.server.Handler = t.router.Handler()
	go func() {
		_ = t.server.ListenAndServe()
		listenAndServeQuited <- struct{}{}
	}()
	return t.server.Addr, nil
}

//Bootup flamingo app with the given modules (and the config in folder given )
func Bootup(modules []dingo.Module, configDir string, config config.Map) BootupInfo {
	if configDir != "" {
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			panic("configdir: " + configDir + " does not exist")
		}
	}
	rw.Lock()
	defer rw.Unlock()
	//add testmodul that listens
	modules = append(modules, new(testmodule))
	//rootArea := rootArea("config")
	if len(os.Args) > 1 {
		os.Args[1] = "serve"
	}
	additionalConfig = config

	application, err := flamingo.NewApplication(modules, flamingo.ConfigDir(configDir))
	if err != nil {
		panic(fmt.Sprintf("unable to get flamingo application: %v", err))
	}

	testmoduli, err := application.ConfigArea().Injector.GetInstance(new(testmodule))
	testmodul := testmoduli.(*testmodule)
	if err != nil {
		panic("unable to get testmodul in flamingo execution area")
	}
	listenAndServeQuited := make(chan struct{})
	port, err := testmodul.startServer(listenAndServeQuited)
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	return BootupInfo{
		func() {
			testmodul.shutdownServer()
		},
		application,
		"localhost" + port,
		listenAndServeQuited,
	}
}
