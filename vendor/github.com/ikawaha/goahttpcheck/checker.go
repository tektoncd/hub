package goahttpcheck

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

type (
	decoder      = func(*http.Request) goahttp.Decoder
	encoder      = func(context.Context, http.ResponseWriter) goahttp.Encoder
	errorHandler = func(context.Context, http.ResponseWriter, error)
	formatter    = func(error) goahttp.Statuser
	middleware   = func(http.Handler) http.Handler

	// HandlerBuilder represents the goa http handler builder.
	HandlerBuilder func(goa.Endpoint, goahttp.Muxer, decoder, encoder, errorHandler, formatter) http.Handler
	// HandlerMounter represents the goa http handler mounter.
	HandlerMounter func(goahttp.Muxer, http.Handler)
)

// APIChecker represents the API checker.
type APIChecker struct {
	Mux           goahttp.Muxer
	Middleware    []middleware
	Decoder       decoder
	Encoder       encoder
	ErrorHandler  errorHandler
	Formatter     formatter
	ClientOptions []httpcheck.Option
}

// Option represents options for API checker.
type Option func(c *APIChecker)

// ClientTimeout sets the client timeout.
func ClientTimeout(d time.Duration) Option {
	return func(c *APIChecker) {
		c.ClientOptions = append(c.ClientOptions, httpcheck.ClientTimeout(d))
	}
}

// CheckRedirect sets the policy of redirection to the HTTP client.
func CheckRedirect(policy func(req *http.Request, via []*http.Request) error) Option {
	return func(c *APIChecker) {
		c.ClientOptions = append(c.ClientOptions, httpcheck.CheckRedirect(policy))
	}
}

// NoRedirect is the alias of the following:
//
//  CheckRedirect(func(req *http.Request, via []*http.Request) error {
//      return http.ErrUseLastResponse
//  })
//
// Client returns ErrUseLastResponse, the next request is not sent and the most recent
// response is returned with its body unclosed.
func NoRedirect() Option {
	return CheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})
}

// Muxer sets the goa http multiplexer.
func Muxer(mux goahttp.Muxer) Option {
	return func(c *APIChecker) {
		c.Mux = mux
	}
}

// Decoder sets the decoder.
func Decoder(dec decoder) Option {
	return func(c *APIChecker) {
		c.Decoder = dec
	}
}

// Encoder sets the encoder.
func Encoder(enc encoder) Option {
	return func(c *APIChecker) {
		c.Encoder = enc
	}
}

// ErrorHandler sets the error handler.
func ErrorHandler(eh errorHandler) Option {
	return func(c *APIChecker) {
		c.ErrorHandler = eh
	}
}

// Formatter sets the error handler.
func Formatter(fm formatter) Option {
	return func(c *APIChecker) {
		c.Formatter = fm
	}
}

// New constructs API checker.
func New(options ...Option) *APIChecker {
	ret := &APIChecker{
		Mux:        goahttp.NewMuxer(),
		Middleware: []middleware{},
		Decoder:    goahttp.RequestDecoder,
		Encoder:    goahttp.ResponseEncoder,
		ErrorHandler: func(ctx context.Context, w http.ResponseWriter, err error) {
			log.Printf("ERROR: %v", err)
		},
	}
	for _, v := range options {
		v(ret)
	}
	return ret
}

// Mount mounts the endpoint handler.
func (c *APIChecker) Mount(builder HandlerBuilder, mounter HandlerMounter, endpoint goa.Endpoint, middlewares ...middleware) {
	handler := builder(endpoint, c.Mux, c.Decoder, c.Encoder, c.ErrorHandler, c.Formatter)
	for _, v := range middlewares {
		handler = v(handler)
	}
	mounter(c.Mux, handler)
}

// Use sets the middleware.
func (c *APIChecker) Use(middleware func(http.Handler) http.Handler) {
	c.Middleware = append(c.Middleware, middleware)
}

// Test returns a http checker that tests the endpoint.
// see. https://github.com/ikawaha/httpcheck/
func (c APIChecker) Test(t *testing.T, method, path string) *httpcheck.Tester {
	var handler http.Handler = c.Mux
	for _, v := range c.Middleware {
		handler = v(handler)
	}
	return httpcheck.New(handler, c.ClientOptions...).Test(t, method, path)
}
