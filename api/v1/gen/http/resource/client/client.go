// Code generated by goa v3.17.2, DO NOT EDIT.
//
// resource client HTTP transport
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/v1/design

package client

import (
	"context"
	"net/http"

	resource "github.com/tektoncd/hub/api/v1/gen/resource"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the resource service endpoint HTTP clients.
type Client struct {
	// Query Doer is the HTTP client used to make requests to the Query endpoint.
	QueryDoer goahttp.Doer

	// List Doer is the HTTP client used to make requests to the List endpoint.
	ListDoer goahttp.Doer

	// VersionsByID Doer is the HTTP client used to make requests to the
	// VersionsByID endpoint.
	VersionsByIDDoer goahttp.Doer

	// ByCatalogKindNameVersion Doer is the HTTP client used to make requests to
	// the ByCatalogKindNameVersion endpoint.
	ByCatalogKindNameVersionDoer goahttp.Doer

	// ByCatalogKindNameVersionReadme Doer is the HTTP client used to make requests
	// to the ByCatalogKindNameVersionReadme endpoint.
	ByCatalogKindNameVersionReadmeDoer goahttp.Doer

	// ByCatalogKindNameVersionYaml Doer is the HTTP client used to make requests
	// to the ByCatalogKindNameVersionYaml endpoint.
	ByCatalogKindNameVersionYamlDoer goahttp.Doer

	// ByVersionID Doer is the HTTP client used to make requests to the ByVersionId
	// endpoint.
	ByVersionIDDoer goahttp.Doer

	// ByCatalogKindName Doer is the HTTP client used to make requests to the
	// ByCatalogKindName endpoint.
	ByCatalogKindNameDoer goahttp.Doer

	// ByID Doer is the HTTP client used to make requests to the ById endpoint.
	ByIDDoer goahttp.Doer

	// GetRawYamlByCatalogKindNameVersion Doer is the HTTP client used to make
	// requests to the GetRawYamlByCatalogKindNameVersion endpoint.
	GetRawYamlByCatalogKindNameVersionDoer goahttp.Doer

	// GetLatestRawYamlByCatalogKindName Doer is the HTTP client used to make
	// requests to the GetLatestRawYamlByCatalogKindName endpoint.
	GetLatestRawYamlByCatalogKindNameDoer goahttp.Doer

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

// NewClient instantiates HTTP clients for all the resource service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		QueryDoer:                              doer,
		ListDoer:                               doer,
		VersionsByIDDoer:                       doer,
		ByCatalogKindNameVersionDoer:           doer,
		ByCatalogKindNameVersionReadmeDoer:     doer,
		ByCatalogKindNameVersionYamlDoer:       doer,
		ByVersionIDDoer:                        doer,
		ByCatalogKindNameDoer:                  doer,
		ByIDDoer:                               doer,
		GetRawYamlByCatalogKindNameVersionDoer: doer,
		GetLatestRawYamlByCatalogKindNameDoer:  doer,
		CORSDoer:                               doer,
		RestoreResponseBody:                    restoreBody,
		scheme:                                 scheme,
		host:                                   host,
		decoder:                                dec,
		encoder:                                enc,
	}
}

// Query returns an endpoint that makes HTTP requests to the resource service
// Query server.
func (c *Client) Query() goa.Endpoint {
	var (
		encodeRequest  = EncodeQueryRequest(c.encoder)
		decodeResponse = DecodeQueryResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildQueryRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.QueryDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "Query", err)
		}
		return decodeResponse(resp)
	}
}

// List returns an endpoint that makes HTTP requests to the resource service
// List server.
func (c *Client) List() goa.Endpoint {
	var (
		encodeRequest  = EncodeListRequest(c.encoder)
		decodeResponse = DecodeListResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildListRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ListDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "List", err)
		}
		return decodeResponse(resp)
	}
}

// VersionsByID returns an endpoint that makes HTTP requests to the resource
// service VersionsByID server.
func (c *Client) VersionsByID() goa.Endpoint {
	var (
		decodeResponse = DecodeVersionsByIDResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildVersionsByIDRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.VersionsByIDDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "VersionsByID", err)
		}
		return decodeResponse(resp)
	}
}

// ByCatalogKindNameVersion returns an endpoint that makes HTTP requests to the
// resource service ByCatalogKindNameVersion server.
func (c *Client) ByCatalogKindNameVersion() goa.Endpoint {
	var (
		decodeResponse = DecodeByCatalogKindNameVersionResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildByCatalogKindNameVersionRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ByCatalogKindNameVersionDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "ByCatalogKindNameVersion", err)
		}
		return decodeResponse(resp)
	}
}

// ByCatalogKindNameVersionReadme returns an endpoint that makes HTTP requests
// to the resource service ByCatalogKindNameVersionReadme server.
func (c *Client) ByCatalogKindNameVersionReadme() goa.Endpoint {
	var (
		decodeResponse = DecodeByCatalogKindNameVersionReadmeResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildByCatalogKindNameVersionReadmeRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ByCatalogKindNameVersionReadmeDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "ByCatalogKindNameVersionReadme", err)
		}
		return decodeResponse(resp)
	}
}

// ByCatalogKindNameVersionYaml returns an endpoint that makes HTTP requests to
// the resource service ByCatalogKindNameVersionYaml server.
func (c *Client) ByCatalogKindNameVersionYaml() goa.Endpoint {
	var (
		decodeResponse = DecodeByCatalogKindNameVersionYamlResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildByCatalogKindNameVersionYamlRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ByCatalogKindNameVersionYamlDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "ByCatalogKindNameVersionYaml", err)
		}
		return decodeResponse(resp)
	}
}

// ByVersionID returns an endpoint that makes HTTP requests to the resource
// service ByVersionId server.
func (c *Client) ByVersionID() goa.Endpoint {
	var (
		decodeResponse = DecodeByVersionIDResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildByVersionIDRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ByVersionIDDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "ByVersionId", err)
		}
		return decodeResponse(resp)
	}
}

// ByCatalogKindName returns an endpoint that makes HTTP requests to the
// resource service ByCatalogKindName server.
func (c *Client) ByCatalogKindName() goa.Endpoint {
	var (
		encodeRequest  = EncodeByCatalogKindNameRequest(c.encoder)
		decodeResponse = DecodeByCatalogKindNameResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildByCatalogKindNameRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ByCatalogKindNameDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "ByCatalogKindName", err)
		}
		return decodeResponse(resp)
	}
}

// ByID returns an endpoint that makes HTTP requests to the resource service
// ById server.
func (c *Client) ByID() goa.Endpoint {
	var (
		decodeResponse = DecodeByIDResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildByIDRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.ByIDDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "ById", err)
		}
		return decodeResponse(resp)
	}
}

// GetRawYamlByCatalogKindNameVersion returns an endpoint that makes HTTP
// requests to the resource service GetRawYamlByCatalogKindNameVersion server.
func (c *Client) GetRawYamlByCatalogKindNameVersion() goa.Endpoint {
	var (
		decodeResponse = DecodeGetRawYamlByCatalogKindNameVersionResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildGetRawYamlByCatalogKindNameVersionRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetRawYamlByCatalogKindNameVersionDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "GetRawYamlByCatalogKindNameVersion", err)
		}
		_, err = decodeResponse(resp)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
		return &resource.GetRawYamlByCatalogKindNameVersionResponseData{Body: resp.Body}, nil
	}
}

// GetLatestRawYamlByCatalogKindName returns an endpoint that makes HTTP
// requests to the resource service GetLatestRawYamlByCatalogKindName server.
func (c *Client) GetLatestRawYamlByCatalogKindName() goa.Endpoint {
	var (
		decodeResponse = DecodeGetLatestRawYamlByCatalogKindNameResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v any) (any, error) {
		req, err := c.BuildGetLatestRawYamlByCatalogKindNameRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.GetLatestRawYamlByCatalogKindNameDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("resource", "GetLatestRawYamlByCatalogKindName", err)
		}
		_, err = decodeResponse(resp)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
		return &resource.GetLatestRawYamlByCatalogKindNameResponseData{Body: resp.Body}, nil
	}
}
