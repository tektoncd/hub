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
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"

	admin "github.com/tektoncd/hub/api/gen/admin"
	auth "github.com/tektoncd/hub/api/gen/auth"
	catalog "github.com/tektoncd/hub/api/gen/catalog"
	category "github.com/tektoncd/hub/api/gen/category"
	adminsvr "github.com/tektoncd/hub/api/gen/http/admin/server"
	authsvr "github.com/tektoncd/hub/api/gen/http/auth/server"
	catalogsvr "github.com/tektoncd/hub/api/gen/http/catalog/server"
	categorysvr "github.com/tektoncd/hub/api/gen/http/category/server"
	ratingsvr "github.com/tektoncd/hub/api/gen/http/rating/server"
	resourcesvr "github.com/tektoncd/hub/api/gen/http/resource/server"
	statussvr "github.com/tektoncd/hub/api/gen/http/status/server"
	swaggersvr "github.com/tektoncd/hub/api/gen/http/swagger/server"
	usersvr "github.com/tektoncd/hub/api/gen/http/user/server"
	"github.com/tektoncd/hub/api/gen/log"
	rating "github.com/tektoncd/hub/api/gen/rating"
	resource "github.com/tektoncd/hub/api/gen/resource"
	status "github.com/tektoncd/hub/api/gen/status"
	user "github.com/tektoncd/hub/api/gen/user"
	v1resourcesvr "github.com/tektoncd/hub/api/v1/gen/http/resource/server"
	v1swaggersvr "github.com/tektoncd/hub/api/v1/gen/http/swagger/server"
	v1resource "github.com/tektoncd/hub/api/v1/gen/resource"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(
	ctx context.Context, u *url.URL,
	adminEndpoints *admin.Endpoints,
	authEndpoints *auth.Endpoints,
	catalogEndpoints *catalog.Endpoints,
	categoryEndpoints *category.Endpoints,
	ratingEndpoints *rating.Endpoints,
	resourceEndpoints *resource.Endpoints,
	v1resourceEndpoints *v1resource.Endpoints,
	statusEndpoints *status.Endpoints,
	userEndpoints *user.Endpoints,
	wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {

	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = logger
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		adminServer      *adminsvr.Server
		authServer       *authsvr.Server
		catalogServer    *catalogsvr.Server
		categoryServer   *categorysvr.Server
		ratingServer     *ratingsvr.Server
		resourceServer   *resourcesvr.Server
		v1resourceServer *v1resourcesvr.Server
		statusServer     *statussvr.Server
		swaggerServer    *swaggersvr.Server
		v1swaggerServer  *v1swaggersvr.Server
		userServer       *usersvr.Server
	)
	{
		eh := errorHandler(logger)
		adminServer = adminsvr.New(adminEndpoints, mux, dec, enc, eh, nil)
		authServer = authsvr.New(authEndpoints, mux, dec, enc, eh, nil)
		catalogServer = catalogsvr.New(catalogEndpoints, mux, dec, enc, eh, nil)
		categoryServer = categorysvr.New(categoryEndpoints, mux, dec, enc, eh, nil)
		ratingServer = ratingsvr.New(ratingEndpoints, mux, dec, enc, eh, nil)
		resourceServer = resourcesvr.New(resourceEndpoints, mux, dec, enc, eh, nil)
		v1resourceServer = v1resourcesvr.New(v1resourceEndpoints, mux, dec, enc, eh, nil)
		statusServer = statussvr.New(statusEndpoints, mux, dec, enc, eh, nil)
		swaggerServer = swaggersvr.New(nil, mux, dec, enc, eh, nil)
		v1swaggerServer = v1swaggersvr.New(nil, mux, dec, enc, eh, nil)
		userServer = usersvr.New(userEndpoints, mux, dec, enc, eh, nil)

		if debug {
			servers := goahttp.Servers{
				adminServer,
				authServer,
				catalogServer,
				categoryServer,
				ratingServer,
				resourceServer,
				v1resourceServer,
				statusServer,
				swaggerServer,
				v1swaggerServer,
				userServer,
			}
			servers.Use(httpmdlwr.Debug(mux, os.Stdout))
		}
	}
	// Configure the mux.
	adminsvr.Mount(mux, adminServer)
	authsvr.Mount(mux, authServer)
	catalogsvr.Mount(mux, catalogServer)
	categorysvr.Mount(mux, categoryServer)
	ratingsvr.Mount(mux, ratingServer)
	resourcesvr.Mount(mux, resourceServer)
	v1resourcesvr.Mount(mux, v1resourceServer)
	statussvr.Mount(mux, statusServer)
	swaggersvr.Mount(mux, swaggerServer)
	v1swaggersvr.Mount(mux, v1swaggerServer)
	usersvr.Mount(mux, userServer)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		handler = httpmdlwr.Log(adapter)(handler)
		handler = httpmdlwr.RequestID()(handler)
	}

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler}
	for _, m := range adminServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range authServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range catalogServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range categoryServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range ratingServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range resourceServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range v1resourceServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range statusServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range swaggerServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range v1swaggerServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range userServer.Mounts {
		logger.Infof("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Infof("HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Infof("shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.With("id", id).Error(err.Error())
	}
}
