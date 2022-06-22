// Copyright © 2020 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	admin "github.com/tektoncd/hub/api/gen/admin"
	catalog "github.com/tektoncd/hub/api/gen/catalog"
	category "github.com/tektoncd/hub/api/gen/category"
	"github.com/tektoncd/hub/api/gen/log"
	rating "github.com/tektoncd/hub/api/gen/rating"
	resource "github.com/tektoncd/hub/api/gen/resource"
	status "github.com/tektoncd/hub/api/gen/status"
	"github.com/tektoncd/hub/api/pkg/app"
	auth "github.com/tektoncd/hub/api/pkg/auth"
	"github.com/tektoncd/hub/api/pkg/db/initializer"
	adminsvc "github.com/tektoncd/hub/api/pkg/service/admin"
	catalogsvc "github.com/tektoncd/hub/api/pkg/service/catalog"
	categorysvc "github.com/tektoncd/hub/api/pkg/service/category"
	ratingsvc "github.com/tektoncd/hub/api/pkg/service/rating"
	resourcesvc "github.com/tektoncd/hub/api/pkg/service/resource"
	statussvc "github.com/tektoncd/hub/api/pkg/service/status"
	userSvc "github.com/tektoncd/hub/api/pkg/user"
	v1catalog "github.com/tektoncd/hub/api/v1/gen/catalog"
	v1resource "github.com/tektoncd/hub/api/v1/gen/resource"
	v1catalogsvc "github.com/tektoncd/hub/api/v1/service/catalog"
	v1resourcesvc "github.com/tektoncd/hub/api/v1/service/resource"

	// Go runtime is unaware of CPU quota which means it will set GOMAXPROCS
	// to underlying host vm node. This high value means that GO runtime
	// scheduler assumes that it has more threads and does context switching
	// when it might work with fewer threads.
	// This doesn't happen# with our other controllers and services because
	// sharedmain already import this package for them.
	_ "go.uber.org/automaxprocs"
)

func main() {
	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		httpPortF = flag.String("http-port", "", "HTTP port (overrides host HTTP port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	var (
		api    app.Config
		logger *log.Logger
		err    error
	)
	{
		api, err = app.FromEnv()
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: failed to initialise: %s", err)
			os.Exit(1)
		}

		logger = api.Logger("main")
		defer api.Cleanup()
	}

	// Populate Tables
	initializer := initializer.New(api)
	if _, err := initializer.Run(context.Background()); err != nil {
		logger.Fatalf("Failed to populate table: %v", err)
	}

	// Add apiserver-bot user account
	db := initializer.DB(context.Background())
	if err := initializer.CreateApiServerAccount(db, logger); err != nil {
		logger.Fatalf("Failed to add resources: %v", err)
	}

	// Initialize the services.
	var (
		adminSvc      admin.Service
		catalogSvc    catalog.Service
		v1catalogSvc  v1catalog.Service
		categorySvc   category.Service
		ratingSvc     rating.Service
		resourceSvc   resource.Service
		v1resourceSvc v1resource.Service
		statusSvc     status.Service
	)
	{
		adminSvc = adminsvc.New(api)
		catalogSvc = catalogsvc.New(api)
		v1catalogSvc = v1catalogsvc.New(api)
		categorySvc = categorysvc.New(api)
		ratingSvc = ratingsvc.New(api)
		resourceSvc = resourcesvc.New(api)
		v1resourceSvc = v1resourcesvc.New(api)
		statusSvc = statussvc.New(api)
	}

	// Wrap the services in endpoints that can be invoked from other services
	// potentially running in different processes.
	var (
		adminEndpoints      *admin.Endpoints
		catalogEndpoints    *catalog.Endpoints
		v1catalogEndpoints  *v1catalog.Endpoints
		categoryEndpoints   *category.Endpoints
		ratingEndpoints     *rating.Endpoints
		resourceEndpoints   *resource.Endpoints
		v1resourceEndpoints *v1resource.Endpoints
		statusEndpoints     *status.Endpoints
	)
	{
		adminEndpoints = admin.NewEndpoints(adminSvc)
		catalogEndpoints = catalog.NewEndpoints(catalogSvc)
		v1catalogEndpoints = v1catalog.NewEndpoints(v1catalogSvc)
		categoryEndpoints = category.NewEndpoints(categorySvc)
		ratingEndpoints = rating.NewEndpoints(ratingSvc)
		resourceEndpoints = resource.NewEndpoints(resourceSvc)
		v1resourceEndpoints = v1resource.NewEndpoints(v1resourceSvc)
		statusEndpoints = status.NewEndpoints(statusSvc)
	}

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// Start the servers and send errors (if any) to the error channel.
	switch *hostF {
	case "localhost":
		{
			addr := "http://:8000"
			u, err := url.Parse(addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid URL %#v: %s\n", addr, err)
				os.Exit(1)
			}
			if *secureF {
				u.Scheme = "https"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			if *httpPortF != "" {
				h := strings.Split(u.Host, ":")[0]
				u.Host = h + ":" + *httpPortF
			} else if u.Port() == "" {
				u.Host += ":80"
			}
			handleHTTPServer(
				ctx, u,
				adminEndpoints,
				catalogEndpoints,
				v1catalogEndpoints,
				categoryEndpoints,
				ratingEndpoints,
				resourceEndpoints,
				v1resourceEndpoints,
				statusEndpoints,
				&wg, errc, api.Logger("http"), *dbgF,
			)
		}

	default:
		fmt.Fprintf(os.Stderr, "invalid host argument: %q (valid hosts: localhost)\n", *hostF)
	}

	r := mux.NewRouter()

	authPort := "4200"

	auth.AuthProvider(r, api)
	userSvc.User(r, api)
	go func() {
		// start the web server on port and accept requests
		logger.Infof("AUTH server listening on port %q", authPort)
		logger.Fatal(http.ListenAndServe(":"+authPort,
			handlers.CORS(handlers.AllowedHeaders(
				[]string{"Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST"}),
				handlers.AllowedOrigins([]string{"*"}))(r)))
	}()

	// Wait for signal.
	logger.Infof("exiting (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Info("exited")
}
