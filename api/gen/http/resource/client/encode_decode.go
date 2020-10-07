// Code generated by goa v3.2.2, DO NOT EDIT.
//
// resource HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	resource "github.com/tektoncd/hub/api/gen/resource"
	resourceviews "github.com/tektoncd/hub/api/gen/resource/views"
	goahttp "goa.design/goa/v3/http"
)

// BuildQueryRequest instantiates a HTTP request object with method and path
// set to call the "resource" service "Query" endpoint
func (c *Client) BuildQueryRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: QueryResourcePath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("resource", "Query", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeQueryRequest returns an encoder for requests sent to the resource
// Query server.
func EncodeQueryRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*resource.QueryPayload)
		if !ok {
			return goahttp.ErrInvalidType("resource", "Query", "*resource.QueryPayload", v)
		}
		values := req.URL.Query()
		values.Add("name", p.Name)
		for _, value := range p.Kinds {
			values.Add("kinds", value)
		}
		for _, value := range p.Tags {
			values.Add("tags", value)
		}
		values.Add("limit", fmt.Sprintf("%v", p.Limit))
		values.Add("match", p.Match)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeQueryResponse returns a decoder for responses returned by the resource
// Query endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeQueryResponse may return the following errors:
//	- "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//	- "invalid-kind" (type *goa.ServiceError): http.StatusBadRequest
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- error: internal error
func DecodeQueryResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body QueryResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "Query", err)
			}
			p := NewQueryResourceCollectionOK(body)
			view := "withoutVersion"
			vres := resourceviews.ResourceCollection{Projected: p, View: view}
			if err = resourceviews.ValidateResourceCollection(vres); err != nil {
				return nil, goahttp.ErrValidationError("resource", "Query", err)
			}
			res := resource.NewResourceCollection(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body QueryInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "Query", err)
			}
			err = ValidateQueryInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "Query", err)
			}
			return nil, NewQueryInternalError(&body)
		case http.StatusBadRequest:
			var (
				body QueryInvalidKindResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "Query", err)
			}
			err = ValidateQueryInvalidKindResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "Query", err)
			}
			return nil, NewQueryInvalidKind(&body)
		case http.StatusNotFound:
			var (
				body QueryNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "Query", err)
			}
			err = ValidateQueryNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "Query", err)
			}
			return nil, NewQueryNotFound(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("resource", "Query", resp.StatusCode, string(body))
		}
	}
}

// BuildListRequest instantiates a HTTP request object with method and path set
// to call the "resource" service "List" endpoint
func (c *Client) BuildListRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ListResourcePath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("resource", "List", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeListRequest returns an encoder for requests sent to the resource List
// server.
func EncodeListRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, interface{}) error {
	return func(req *http.Request, v interface{}) error {
		p, ok := v.(*resource.ListPayload)
		if !ok {
			return goahttp.ErrInvalidType("resource", "List", "*resource.ListPayload", v)
		}
		values := req.URL.Query()
		values.Add("limit", fmt.Sprintf("%v", p.Limit))
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeListResponse returns a decoder for responses returned by the resource
// List endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeListResponse may return the following errors:
//	- "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//	- error: internal error
func DecodeListResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ListResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "List", err)
			}
			p := NewListResourceCollectionOK(body)
			view := "withoutVersion"
			vres := resourceviews.ResourceCollection{Projected: p, View: view}
			if err = resourceviews.ValidateResourceCollection(vres); err != nil {
				return nil, goahttp.ErrValidationError("resource", "List", err)
			}
			res := resource.NewResourceCollection(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ListInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "List", err)
			}
			err = ValidateListInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "List", err)
			}
			return nil, NewListInternalError(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("resource", "List", resp.StatusCode, string(body))
		}
	}
}

// BuildVersionsByIDRequest instantiates a HTTP request object with method and
// path set to call the "resource" service "VersionsByID" endpoint
func (c *Client) BuildVersionsByIDRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		id uint
	)
	{
		p, ok := v.(*resource.VersionsByIDPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("resource", "VersionsByID", "*resource.VersionsByIDPayload", v)
		}
		id = p.ID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: VersionsByIDResourcePath(id)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("resource", "VersionsByID", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeVersionsByIDResponse returns a decoder for responses returned by the
// resource VersionsByID endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeVersionsByIDResponse may return the following errors:
//	- "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- error: internal error
func DecodeVersionsByIDResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body VersionsByIDResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "VersionsByID", err)
			}
			p := NewVersionsByIDVersionsOK(&body)
			view := "default"
			vres := &resourceviews.Versions{Projected: p, View: view}
			if err = resourceviews.ValidateVersions(vres); err != nil {
				return nil, goahttp.ErrValidationError("resource", "VersionsByID", err)
			}
			res := resource.NewVersions(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body VersionsByIDInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "VersionsByID", err)
			}
			err = ValidateVersionsByIDInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "VersionsByID", err)
			}
			return nil, NewVersionsByIDInternalError(&body)
		case http.StatusNotFound:
			var (
				body VersionsByIDNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "VersionsByID", err)
			}
			err = ValidateVersionsByIDNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "VersionsByID", err)
			}
			return nil, NewVersionsByIDNotFound(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("resource", "VersionsByID", resp.StatusCode, string(body))
		}
	}
}

// BuildByCatalogKindNameVersionRequest instantiates a HTTP request object with
// method and path set to call the "resource" service
// "ByCatalogKindNameVersion" endpoint
func (c *Client) BuildByCatalogKindNameVersionRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		catalog string
		kind    string
		name    string
		version string
	)
	{
		p, ok := v.(*resource.ByCatalogKindNameVersionPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("resource", "ByCatalogKindNameVersion", "*resource.ByCatalogKindNameVersionPayload", v)
		}
		catalog = p.Catalog
		kind = p.Kind
		name = p.Name
		version = p.Version
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ByCatalogKindNameVersionResourcePath(catalog, kind, name, version)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("resource", "ByCatalogKindNameVersion", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeByCatalogKindNameVersionResponse returns a decoder for responses
// returned by the resource ByCatalogKindNameVersion endpoint. restoreBody
// controls whether the response body should be restored after having been read.
// DecodeByCatalogKindNameVersionResponse may return the following errors:
//	- "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- error: internal error
func DecodeByCatalogKindNameVersionResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ByCatalogKindNameVersionResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByCatalogKindNameVersion", err)
			}
			p := NewByCatalogKindNameVersionVersionOK(&body)
			view := "default"
			vres := &resourceviews.Version{Projected: p, View: view}
			if err = resourceviews.ValidateVersion(vres); err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByCatalogKindNameVersion", err)
			}
			res := resource.NewVersion(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ByCatalogKindNameVersionInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByCatalogKindNameVersion", err)
			}
			err = ValidateByCatalogKindNameVersionInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByCatalogKindNameVersion", err)
			}
			return nil, NewByCatalogKindNameVersionInternalError(&body)
		case http.StatusNotFound:
			var (
				body ByCatalogKindNameVersionNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByCatalogKindNameVersion", err)
			}
			err = ValidateByCatalogKindNameVersionNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByCatalogKindNameVersion", err)
			}
			return nil, NewByCatalogKindNameVersionNotFound(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("resource", "ByCatalogKindNameVersion", resp.StatusCode, string(body))
		}
	}
}

// BuildByVersionIDRequest instantiates a HTTP request object with method and
// path set to call the "resource" service "ByVersionId" endpoint
func (c *Client) BuildByVersionIDRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		versionID uint
	)
	{
		p, ok := v.(*resource.ByVersionIDPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("resource", "ByVersionId", "*resource.ByVersionIDPayload", v)
		}
		versionID = p.VersionID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ByVersionIDResourcePath(versionID)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("resource", "ByVersionId", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeByVersionIDResponse returns a decoder for responses returned by the
// resource ByVersionId endpoint. restoreBody controls whether the response
// body should be restored after having been read.
// DecodeByVersionIDResponse may return the following errors:
//	- "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- error: internal error
func DecodeByVersionIDResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ByVersionIDResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByVersionId", err)
			}
			p := NewByVersionIDVersionOK(&body)
			view := "default"
			vres := &resourceviews.Version{Projected: p, View: view}
			if err = resourceviews.ValidateVersion(vres); err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByVersionId", err)
			}
			res := resource.NewVersion(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ByVersionIDInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByVersionId", err)
			}
			err = ValidateByVersionIDInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByVersionId", err)
			}
			return nil, NewByVersionIDInternalError(&body)
		case http.StatusNotFound:
			var (
				body ByVersionIDNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByVersionId", err)
			}
			err = ValidateByVersionIDNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByVersionId", err)
			}
			return nil, NewByVersionIDNotFound(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("resource", "ByVersionId", resp.StatusCode, string(body))
		}
	}
}

// BuildByCatalogKindNameRequest instantiates a HTTP request object with method
// and path set to call the "resource" service "ByCatalogKindName" endpoint
func (c *Client) BuildByCatalogKindNameRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		catalog string
		kind    string
		name    string
	)
	{
		p, ok := v.(*resource.ByCatalogKindNamePayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("resource", "ByCatalogKindName", "*resource.ByCatalogKindNamePayload", v)
		}
		catalog = p.Catalog
		kind = p.Kind
		name = p.Name
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ByCatalogKindNameResourcePath(catalog, kind, name)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("resource", "ByCatalogKindName", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeByCatalogKindNameResponse returns a decoder for responses returned by
// the resource ByCatalogKindName endpoint. restoreBody controls whether the
// response body should be restored after having been read.
// DecodeByCatalogKindNameResponse may return the following errors:
//	- "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- error: internal error
func DecodeByCatalogKindNameResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ByCatalogKindNameResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByCatalogKindName", err)
			}
			p := NewByCatalogKindNameResourceCollectionOK(body)
			view := "withoutVersion"
			vres := resourceviews.ResourceCollection{Projected: p, View: view}
			if err = resourceviews.ValidateResourceCollection(vres); err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByCatalogKindName", err)
			}
			res := resource.NewResourceCollection(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ByCatalogKindNameInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByCatalogKindName", err)
			}
			err = ValidateByCatalogKindNameInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByCatalogKindName", err)
			}
			return nil, NewByCatalogKindNameInternalError(&body)
		case http.StatusNotFound:
			var (
				body ByCatalogKindNameNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ByCatalogKindName", err)
			}
			err = ValidateByCatalogKindNameNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ByCatalogKindName", err)
			}
			return nil, NewByCatalogKindNameNotFound(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("resource", "ByCatalogKindName", resp.StatusCode, string(body))
		}
	}
}

// BuildByIDRequest instantiates a HTTP request object with method and path set
// to call the "resource" service "ById" endpoint
func (c *Client) BuildByIDRequest(ctx context.Context, v interface{}) (*http.Request, error) {
	var (
		id uint
	)
	{
		p, ok := v.(*resource.ByIDPayload)
		if !ok {
			return nil, goahttp.ErrInvalidType("resource", "ById", "*resource.ByIDPayload", v)
		}
		id = p.ID
	}
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ByIDResourcePath(id)}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("resource", "ById", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeByIDResponse returns a decoder for responses returned by the resource
// ById endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeByIDResponse may return the following errors:
//	- "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//	- "not-found" (type *goa.ServiceError): http.StatusNotFound
//	- error: internal error
func DecodeByIDResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (interface{}, error) {
	return func(resp *http.Response) (interface{}, error) {
		if restoreBody {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body ByIDResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ById", err)
			}
			p := NewByIDResourceOK(&body)
			view := "default"
			vres := &resourceviews.Resource{Projected: p, View: view}
			if err = resourceviews.ValidateResource(vres); err != nil {
				return nil, goahttp.ErrValidationError("resource", "ById", err)
			}
			res := resource.NewResource(vres)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ByIDInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ById", err)
			}
			err = ValidateByIDInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ById", err)
			}
			return nil, NewByIDInternalError(&body)
		case http.StatusNotFound:
			var (
				body ByIDNotFoundResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("resource", "ById", err)
			}
			err = ValidateByIDNotFoundResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("resource", "ById", err)
			}
			return nil, NewByIDNotFound(&body)
		default:
			body, _ := ioutil.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("resource", "ById", resp.StatusCode, string(body))
		}
	}
}

// unmarshalResourceResponseToResourceviewsResourceView builds a value of type
// *resourceviews.ResourceView from a value of type *ResourceResponse.
func unmarshalResourceResponseToResourceviewsResourceView(v *ResourceResponse) *resourceviews.ResourceView {
	res := &resourceviews.ResourceView{
		ID:     v.ID,
		Name:   v.Name,
		Kind:   v.Kind,
		Rating: v.Rating,
	}
	res.Catalog = unmarshalCatalogResponseToResourceviewsCatalogView(v.Catalog)
	res.LatestVersion = unmarshalVersionResponseToResourceviewsVersionView(v.LatestVersion)
	res.Tags = make([]*resourceviews.TagView, len(v.Tags))
	for i, val := range v.Tags {
		res.Tags[i] = unmarshalTagResponseToResourceviewsTagView(val)
	}
	res.Versions = make([]*resourceviews.VersionView, len(v.Versions))
	for i, val := range v.Versions {
		res.Versions[i] = unmarshalVersionResponseToResourceviewsVersionView(val)
	}

	return res
}

// unmarshalCatalogResponseToResourceviewsCatalogView builds a value of type
// *resourceviews.CatalogView from a value of type *CatalogResponse.
func unmarshalCatalogResponseToResourceviewsCatalogView(v *CatalogResponse) *resourceviews.CatalogView {
	res := &resourceviews.CatalogView{
		ID:   v.ID,
		Name: v.Name,
		Type: v.Type,
	}

	return res
}

// unmarshalVersionResponseToResourceviewsVersionView builds a value of type
// *resourceviews.VersionView from a value of type *VersionResponse.
func unmarshalVersionResponseToResourceviewsVersionView(v *VersionResponse) *resourceviews.VersionView {
	res := &resourceviews.VersionView{
		ID:                  v.ID,
		Version:             v.Version,
		DisplayName:         v.DisplayName,
		Description:         v.Description,
		MinPipelinesVersion: v.MinPipelinesVersion,
		RawURL:              v.RawURL,
		WebURL:              v.WebURL,
		UpdatedAt:           v.UpdatedAt,
	}
	res.Resource = unmarshalResourceResponseToResourceviewsResourceView(v.Resource)

	return res
}

// unmarshalTagResponseToResourceviewsTagView builds a value of type
// *resourceviews.TagView from a value of type *TagResponse.
func unmarshalTagResponseToResourceviewsTagView(v *TagResponse) *resourceviews.TagView {
	res := &resourceviews.TagView{
		ID:   v.ID,
		Name: v.Name,
	}

	return res
}

// unmarshalVersionResponseBodyToResourceviewsVersionView builds a value of
// type *resourceviews.VersionView from a value of type *VersionResponseBody.
func unmarshalVersionResponseBodyToResourceviewsVersionView(v *VersionResponseBody) *resourceviews.VersionView {
	res := &resourceviews.VersionView{
		ID:                  v.ID,
		Version:             v.Version,
		DisplayName:         v.DisplayName,
		Description:         v.Description,
		MinPipelinesVersion: v.MinPipelinesVersion,
		RawURL:              v.RawURL,
		WebURL:              v.WebURL,
		UpdatedAt:           v.UpdatedAt,
	}
	res.Resource = unmarshalResourceResponseBodyToResourceviewsResourceView(v.Resource)

	return res
}

// unmarshalResourceResponseBodyToResourceviewsResourceView builds a value of
// type *resourceviews.ResourceView from a value of type *ResourceResponseBody.
func unmarshalResourceResponseBodyToResourceviewsResourceView(v *ResourceResponseBody) *resourceviews.ResourceView {
	res := &resourceviews.ResourceView{
		ID:     v.ID,
		Name:   v.Name,
		Kind:   v.Kind,
		Rating: v.Rating,
	}
	res.Catalog = unmarshalCatalogResponseBodyToResourceviewsCatalogView(v.Catalog)
	res.LatestVersion = unmarshalVersionResponseBodyToResourceviewsVersionView(v.LatestVersion)
	res.Tags = make([]*resourceviews.TagView, len(v.Tags))
	for i, val := range v.Tags {
		res.Tags[i] = unmarshalTagResponseBodyToResourceviewsTagView(val)
	}
	res.Versions = make([]*resourceviews.VersionView, len(v.Versions))
	for i, val := range v.Versions {
		res.Versions[i] = unmarshalVersionResponseBodyToResourceviewsVersionView(val)
	}

	return res
}

// unmarshalCatalogResponseBodyToResourceviewsCatalogView builds a value of
// type *resourceviews.CatalogView from a value of type *CatalogResponseBody.
func unmarshalCatalogResponseBodyToResourceviewsCatalogView(v *CatalogResponseBody) *resourceviews.CatalogView {
	res := &resourceviews.CatalogView{
		ID:   v.ID,
		Name: v.Name,
		Type: v.Type,
	}

	return res
}

// unmarshalTagResponseBodyToResourceviewsTagView builds a value of type
// *resourceviews.TagView from a value of type *TagResponseBody.
func unmarshalTagResponseBodyToResourceviewsTagView(v *TagResponseBody) *resourceviews.TagView {
	res := &resourceviews.TagView{
		ID:   v.ID,
		Name: v.Name,
	}

	return res
}

// unmarshalByCatalogKindNameVersionResponseBodyToResourceviewsVersionView
// builds a value of type *resourceviews.VersionView from a value of type
// *ByCatalogKindNameVersionResponseBody.
func unmarshalByCatalogKindNameVersionResponseBodyToResourceviewsVersionView(v *ByCatalogKindNameVersionResponseBody) *resourceviews.VersionView {
	res := &resourceviews.VersionView{
		ID:                  v.ID,
		Version:             v.Version,
		DisplayName:         v.DisplayName,
		Description:         v.Description,
		MinPipelinesVersion: v.MinPipelinesVersion,
		RawURL:              v.RawURL,
		WebURL:              v.WebURL,
		UpdatedAt:           v.UpdatedAt,
	}
	res.Resource = unmarshalResourceResponseBodyToResourceviewsResourceView(v.Resource)

	return res
}

// unmarshalByVersionIDResponseBodyToResourceviewsVersionView builds a value of
// type *resourceviews.VersionView from a value of type
// *ByVersionIDResponseBody.
func unmarshalByVersionIDResponseBodyToResourceviewsVersionView(v *ByVersionIDResponseBody) *resourceviews.VersionView {
	res := &resourceviews.VersionView{
		ID:                  v.ID,
		Version:             v.Version,
		DisplayName:         v.DisplayName,
		Description:         v.Description,
		MinPipelinesVersion: v.MinPipelinesVersion,
		RawURL:              v.RawURL,
		WebURL:              v.WebURL,
		UpdatedAt:           v.UpdatedAt,
	}
	res.Resource = unmarshalResourceResponseBodyToResourceviewsResourceView(v.Resource)

	return res
}

// unmarshalByIDResponseBodyToResourceviewsResourceView builds a value of type
// *resourceviews.ResourceView from a value of type *ByIDResponseBody.
func unmarshalByIDResponseBodyToResourceviewsResourceView(v *ByIDResponseBody) *resourceviews.ResourceView {
	res := &resourceviews.ResourceView{
		ID:     v.ID,
		Name:   v.Name,
		Kind:   v.Kind,
		Rating: v.Rating,
	}
	res.Catalog = unmarshalCatalogResponseBodyToResourceviewsCatalogView(v.Catalog)
	res.LatestVersion = unmarshalVersionResponseBodyToResourceviewsVersionView(v.LatestVersion)
	res.Tags = make([]*resourceviews.TagView, len(v.Tags))
	for i, val := range v.Tags {
		res.Tags[i] = unmarshalTagResponseBodyToResourceviewsTagView(val)
	}
	res.Versions = make([]*resourceviews.VersionView, len(v.Versions))
	for i, val := range v.Versions {
		res.Versions[i] = unmarshalVersionResponseBodyToResourceviewsVersionView(val)
	}

	return res
}
