// Code generated by goa v3.15.1, DO NOT EDIT.
//
// status HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package server

import (
	"context"
	"net/http"

	status "github.com/tektoncd/hub/api/gen/status"
	goahttp "goa.design/goa/v3/http"
)

// EncodeStatusResponse returns an encoder for responses returned by the status
// Status endpoint.
func EncodeStatusResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, any) error {
	return func(ctx context.Context, w http.ResponseWriter, v any) error {
		res, _ := v.(*status.StatusResult)
		enc := encoder(ctx, w)
		body := NewStatusResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// marshalStatusHubServiceToHubServiceResponseBody builds a value of type
// *HubServiceResponseBody from a value of type *status.HubService.
func marshalStatusHubServiceToHubServiceResponseBody(v *status.HubService) *HubServiceResponseBody {
	if v == nil {
		return nil
	}
	res := &HubServiceResponseBody{
		Name:   v.Name,
		Status: v.Status,
		Error:  v.Error,
	}

	return res
}
