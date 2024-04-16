// Code generated by goa v3.16.1, DO NOT EDIT.
//
// catalog HTTP server
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package server

import (
	"context"
	"net/http"

	catalog "github.com/tektoncd/hub/api/gen/catalog"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"goa.design/plugins/v3/cors"
)

// Server lists the catalog service endpoint HTTP handlers.
type Server struct {
	Mounts       []*MountPoint
	Refresh      http.Handler
	RefreshAll   http.Handler
	CatalogError http.Handler
	CORS         http.Handler
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

// New instantiates HTTP handlers for all the catalog service endpoints using
// the provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *catalog.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Refresh", "POST", "/catalog/{catalogName}/refresh"},
			{"RefreshAll", "POST", "/catalog/refresh"},
			{"CatalogError", "GET", "/catalog/{catalogName}/error"},
			{"CORS", "OPTIONS", "/catalog/{catalogName}/refresh"},
			{"CORS", "OPTIONS", "/catalog/refresh"},
			{"CORS", "OPTIONS", "/catalog/{catalogName}/error"},
		},
		Refresh:      NewRefreshHandler(e.Refresh, mux, decoder, encoder, errhandler, formatter),
		RefreshAll:   NewRefreshAllHandler(e.RefreshAll, mux, decoder, encoder, errhandler, formatter),
		CatalogError: NewCatalogErrorHandler(e.CatalogError, mux, decoder, encoder, errhandler, formatter),
		CORS:         NewCORSHandler(),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "catalog" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Refresh = m(s.Refresh)
	s.RefreshAll = m(s.RefreshAll)
	s.CatalogError = m(s.CatalogError)
	s.CORS = m(s.CORS)
}

// MethodNames returns the methods served.
func (s *Server) MethodNames() []string { return catalog.MethodNames[:] }

// Mount configures the mux to serve the catalog endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountRefreshHandler(mux, h.Refresh)
	MountRefreshAllHandler(mux, h.RefreshAll)
	MountCatalogErrorHandler(mux, h.CatalogError)
	MountCORSHandler(mux, h.CORS)
}

// Mount configures the mux to serve the catalog endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountRefreshHandler configures the mux to serve the "catalog" service
// "Refresh" endpoint.
func MountRefreshHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleCatalogOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/catalog/{catalogName}/refresh", f)
}

// NewRefreshHandler creates a HTTP handler which loads the HTTP request and
// calls the "catalog" service "Refresh" endpoint.
func NewRefreshHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeRefreshRequest(mux, decoder)
		encodeResponse = EncodeRefreshResponse(encoder)
		encodeError    = EncodeRefreshError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "Refresh")
		ctx = context.WithValue(ctx, goa.ServiceKey, "catalog")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
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

// MountRefreshAllHandler configures the mux to serve the "catalog" service
// "RefreshAll" endpoint.
func MountRefreshAllHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleCatalogOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/catalog/refresh", f)
}

// NewRefreshAllHandler creates a HTTP handler which loads the HTTP request and
// calls the "catalog" service "RefreshAll" endpoint.
func NewRefreshAllHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeRefreshAllRequest(mux, decoder)
		encodeResponse = EncodeRefreshAllResponse(encoder)
		encodeError    = EncodeRefreshAllError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "RefreshAll")
		ctx = context.WithValue(ctx, goa.ServiceKey, "catalog")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
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

// MountCatalogErrorHandler configures the mux to serve the "catalog" service
// "CatalogError" endpoint.
func MountCatalogErrorHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := HandleCatalogOrigin(h).(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/catalog/{catalogName}/error", f)
}

// NewCatalogErrorHandler creates a HTTP handler which loads the HTTP request
// and calls the "catalog" service "CatalogError" endpoint.
func NewCatalogErrorHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeCatalogErrorRequest(mux, decoder)
		encodeResponse = EncodeCatalogErrorResponse(encoder)
		encodeError    = EncodeCatalogErrorError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "CatalogError")
		ctx = context.WithValue(ctx, goa.ServiceKey, "catalog")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
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
// service catalog.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleCatalogOrigin(h)
	mux.Handle("OPTIONS", "/catalog/{catalogName}/refresh", h.ServeHTTP)
	mux.Handle("OPTIONS", "/catalog/refresh", h.ServeHTTP)
	mux.Handle("OPTIONS", "/catalog/{catalogName}/error", h.ServeHTTP)
}

// NewCORSHandler creates a HTTP handler which returns a simple 204 response.
func NewCORSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
}

// HandleCatalogOrigin applies the CORS response headers corresponding to the
// origin for the service catalog.
func HandleCatalogOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "*") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.WriteHeader(204)
				return
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
