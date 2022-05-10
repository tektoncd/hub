// Code generated by goa v3.7.3, DO NOT EDIT.
//
// status client HTTP transport
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the status service endpoint HTTP clients.
type Client struct {
	// Status Doer is the HTTP client used to make requests to the Status endpoint.
	StatusDoer goahttp.Doer

	// CORS Doer is the HTTP client used to make requests to the  endpoint.
	CORSDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the status service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		StatusDoer:          doer,
		CORSDoer:            doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// Status returns an endpoint that makes HTTP requests to the status service
// Status server.
func (c *Client) Status() goa.Endpoint {
	var (
		decodeResponse = DecodeStatusResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildStatusRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.StatusDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("status", "Status", err)
		}
		return decodeResponse(resp)
	}
}
