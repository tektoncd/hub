// Code generated by goa v3.14.5, DO NOT EDIT.
//
// admin client
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package admin

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "admin" service client.
type Client struct {
	UpdateAgentEndpoint   goa.Endpoint
	RefreshConfigEndpoint goa.Endpoint
}

// NewClient initializes a "admin" service client given the endpoints.
func NewClient(updateAgent, refreshConfig goa.Endpoint) *Client {
	return &Client{
		UpdateAgentEndpoint:   updateAgent,
		RefreshConfigEndpoint: refreshConfig,
	}
}

// UpdateAgent calls the "UpdateAgent" endpoint of the "admin" service.
// UpdateAgent may return the following errors:
//   - "invalid-payload" (type *goa.ServiceError): Invalid request body
//   - "invalid-token" (type *goa.ServiceError): Invalid User token
//   - "invalid-scopes" (type *goa.ServiceError): Invalid Token scopes
//   - "internal-error" (type *goa.ServiceError): Internal server error
//   - error: internal error
func (c *Client) UpdateAgent(ctx context.Context, p *UpdateAgentPayload) (res *UpdateAgentResult, err error) {
	var ires any
	ires, err = c.UpdateAgentEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*UpdateAgentResult), nil
}

// RefreshConfig calls the "RefreshConfig" endpoint of the "admin" service.
// RefreshConfig may return the following errors:
//   - "invalid-payload" (type *goa.ServiceError): Invalid request body
//   - "invalid-token" (type *goa.ServiceError): Invalid User token
//   - "invalid-scopes" (type *goa.ServiceError): Invalid Token scopes
//   - "internal-error" (type *goa.ServiceError): Internal server error
//   - error: internal error
func (c *Client) RefreshConfig(ctx context.Context, p *RefreshConfigPayload) (res *RefreshConfigResult, err error) {
	var ires any
	ires, err = c.RefreshConfigEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*RefreshConfigResult), nil
}
