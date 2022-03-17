// Code generated by goa v3.4.0, DO NOT EDIT.
//
// category HTTP server
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package server

import (
	"context"
	"net/http"

	category "github.com/tektoncd/hub/api/gen/category"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the category service endpoint HTTP handlers.
type Server struct {
	Mounts []*MountPoint
	List   http.Handler
	CORS   http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the category service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *category.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"List", "GET", "/categories"},
			{"List", "GET", "/v1/categories"},
			{"CORS", "OPTIONS", "/categories"},
			{"CORS", "OPTIONS", "/v1/categories"},
		},
		List: NewListHandler(e.List, mux, decoder, encoder, errhandler, formatter),
		CORS: NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "category" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.List = m(s.List)
	s.CORS = m(s.CORS)
}

// Mount configures the mux to serve the category endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountListHandler(mux, h.List)
	MountCORSHandler(mux, h.CORS)
}

// MountListHandler configures the mux to serve the "category" service "list"
// endpoint.
func MountListHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := handleCategoryOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/categories", f)
	mux.Handle("GET", "/v1/categories", f)
}

// NewListHandler creates a HTTP handler which loads the HTTP request and calls
// the "category" service "list" endpoint.
func NewListHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		encodeResponse = EncodeListResponse(encoder)
		encodeError    = EncodeListError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list")
		ctx = context.WithValue(ctx, goa.ServiceKey, "category")
		var err error
		res, err := endpoint(ctx, nil)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service category.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleCategoryOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/categories", f)
	mux.Handle("OPTIONS", "/v1/categories", f)
}

// NewCORSHandler creates a HTTP handler which returns a simple 200 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}

// handleCategoryOrigin applies the CORS response headers corresponding to the
// origin for the service category.
func handleCategoryOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "*") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
