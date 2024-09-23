// Code generated by goa v3.19.1, DO NOT EDIT.
//
// category HTTP client encoders and decoders
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	category "github.com/tektoncd/hub/api/gen/category"
	goahttp "goa.design/goa/v3/http"
)

// BuildListRequest instantiates a HTTP request object with method and path set
// to call the "category" service "list" endpoint
func (c *Client) BuildListRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ListCategoryPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("category", "list", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeListResponse returns a decoder for responses returned by the category
// list endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeListResponse may return the following errors:
//   - "internal-error" (type *goa.ServiceError): http.StatusInternalServerError
//   - error: internal error
func DecodeListResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
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
				return nil, goahttp.ErrDecodingError("category", "list", err)
			}
			err = ValidateListResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("category", "list", err)
			}
			res := NewListResultOK(&body)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body ListInternalErrorResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("category", "list", err)
			}
			err = ValidateListInternalErrorResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("category", "list", err)
			}
			return nil, NewListInternalError(&body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("category", "list", resp.StatusCode, string(body))
		}
	}
}

// unmarshalCategoryResponseBodyToCategoryCategory builds a value of type
// *category.Category from a value of type *CategoryResponseBody.
func unmarshalCategoryResponseBodyToCategoryCategory(v *CategoryResponseBody) *category.Category {
	if v == nil {
		return nil
	}
	res := &category.Category{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}
