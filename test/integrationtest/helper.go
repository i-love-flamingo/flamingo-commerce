package integrationtest

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"flamingo.me/dingo"
	"github.com/gavv/httpexpect/v2"

	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/framework/config"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	testModule struct {
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

// Side effect vars to get status and exchange stuff with the testModule
var rw sync.Mutex

var additionalConfig config.Map

// Configure for your testModule in the app
func (t *testModule) Inject(eventRouter flamingoFramework.EventRouter, router *web.Router) {
	t.eventRouter = eventRouter
	t.router = router
}

// Configure for your testModule in the app
func (t *testModule) Configure(i *dingo.Injector) {
	flamingoFramework.BindEventSubscriber(i).To(t)
}

// Notify gets notified by event router
func (t *testModule) Notify(_ context.Context, event flamingoFramework.Event) {
	switch event.(type) {
	case *flamingoFramework.ShutdownEvent:
		log.Printf("ShutdownEvent event received...")
	}
}

// DefaultConfig of test module
func (t *testModule) DefaultConfig() config.Map {
	return additionalConfig
}

// shutdown until
func (t *testModule) shutdownServer() {
	log.Printf("Trigger ServerShutdownEvent...")
	t.eventRouter.Dispatch(context.Background(), &flamingoFramework.ServerShutdownEvent{})
	_ = t.server.Shutdown(context.Background())
}

// returns the port or error
func (t *testModule) startServer(listenAndServeQuited chan struct{}) (string, error) {
	port := os.Getenv("INTEGRATION_TEST_PORT")
	if port == "" {
		port = "0"
	}

	t.eventRouter.Dispatch(context.Background(), &flamingoFramework.ServerStartEvent{})
	t.server = &http.Server{
		Addr: ":" + port,
	}

	t.server.Handler = t.router.Handler()
	listener, err := net.Listen("tcp", t.server.Addr)
	if err != nil {
		return "", err
	}

	listenerPort := listener.Addr().(*net.TCPAddr).Port

	log.Printf("startServer on port %v", listenerPort)
	go func() {
		_ = t.server.Serve(listener)
		listenAndServeQuited <- struct{}{}
	}()
	return strconv.Itoa(listenerPort), nil
}

// Bootup flamingo app with the given modules (and the config in folder given )
func Bootup(modules []dingo.Module, configDir string, config config.Map) BootupInfo {
	if configDir != "" {
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			panic("configdir: " + configDir + " does not exist")
		}
	}
	rw.Lock()
	defer rw.Unlock()
	//add testModule that listens
	modules = append(modules, new(testModule))
	//rootArea := rootArea("config")
	if len(os.Args) > 1 {
		os.Args[1] = "serve"
	}
	additionalConfig = config

	application, err := flamingo.NewApplication(modules, flamingo.ConfigDir(configDir))
	if err != nil {
		panic(fmt.Sprintf("unable to get flamingo application: %v", err))
	}

	testModuleInterface, err := application.ConfigArea().Injector.GetInstance(new(testModule))
	testModule := testModuleInterface.(*testModule)
	if err != nil {
		panic("unable to get testModule in flamingo execution area")
	}
	listenAndServeQuited := make(chan struct{})
	port, err := testModule.startServer(listenAndServeQuited)
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	return BootupInfo{
		func() {
			testModule.shutdownServer()
		},
		application,
		"localhost:" + port,
		listenAndServeQuited,
	}
}

// NewHTTPExpectWithCookies returns a new Expect object without printer with preset cookies to the base URL
func NewHTTPExpectWithCookies(testingTB httpexpect.TestingTB, baseURL string, cookies map[string]string) *httpexpect.Expect {
	cookieJar := httpexpect.NewCookieJar()
	parse, err := url.Parse(baseURL)
	if err != nil {
		testingTB.Errorf(err.Error())
	}
	theCookies := make([]*http.Cookie, 0, len(cookies))
	for key, value := range cookies {
		theCookies = append(theCookies, &http.Cookie{Name: key, Value: value})
	}

	cookieJar.SetCookies(parse, theCookies)

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: baseURL,
		Client: &http.Client{
			Jar: cookieJar,
		},
		Reporter: httpexpect.NewAssertReporter(testingTB),
		Printers: nil,
	})
}

// NewHTTPExpect returns a new Expect object without printer
func NewHTTPExpect(tb httpexpect.TestingTB, baseURL string) *httpexpect.Expect {
	return NewHTTPExpectWithCookies(tb, baseURL, nil)
}
