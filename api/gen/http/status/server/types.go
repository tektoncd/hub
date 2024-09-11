// Code generated by goa v3.19.0, DO NOT EDIT.
//
// status HTTP server types
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package server

import (
	status "github.com/tektoncd/hub/api/gen/status"
)

// StatusResponseBody is the type of the "status" service "Status" endpoint
// HTTP response body.
type StatusResponseBody struct {
	// List of services and their status
	Services []*HubServiceResponseBody `form:"services,omitempty" json:"services,omitempty" xml:"services,omitempty"`
}

// HubServiceResponseBody is used to define fields on response body types.
type HubServiceResponseBody struct {
	// Name of the service
	Name string `form:"name" json:"name" xml:"name"`
	// Status of the service
	Status string `form:"status" json:"status" xml:"status"`
	// Details of the error if any
	Error *string `form:"error,omitempty" json:"error,omitempty" xml:"error,omitempty"`
}

// NewStatusResponseBody builds the HTTP response body from the result of the
// "Status" endpoint of the "status" service.
func NewStatusResponseBody(res *status.StatusResult) *StatusResponseBody {
	body := &StatusResponseBody{}
	if res.Services != nil {
		body.Services = make([]*HubServiceResponseBody, len(res.Services))
		for i, val := range res.Services {
			body.Services[i] = marshalStatusHubServiceToHubServiceResponseBody(val)
		}
	}
	return body
}
