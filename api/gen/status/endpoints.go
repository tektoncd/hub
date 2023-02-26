// Code generated by goa v3.11.1, DO NOT EDIT.
//
// status endpoints
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package status

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "status" service endpoints.
type Endpoints struct {
	Status goa.Endpoint
}

// NewEndpoints wraps the methods of the "status" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Status: NewStatusEndpoint(s),
	}
}

// Use applies the given middleware to all the "status" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Status = m(e.Status)
}

// NewStatusEndpoint returns an endpoint function that calls the method
// "Status" of service "status".
func NewStatusEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Status(ctx)
	}
}
