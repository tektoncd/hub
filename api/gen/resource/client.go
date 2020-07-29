// Code generated by goa v3.2.0, DO NOT EDIT.
//
// resource client
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package resource

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "resource" service client.
type Client struct {
	QueryEndpoint             goa.Endpoint
	ListEndpoint              goa.Endpoint
	VersionsByIDEndpoint      goa.Endpoint
	ByTypeNameVersionEndpoint goa.Endpoint
	ByVersionIDEndpoint       goa.Endpoint
	ByTypeNameEndpoint        goa.Endpoint
	ByIDEndpoint              goa.Endpoint
}

// NewClient initializes a "resource" service client given the endpoints.
func NewClient(query, list, versionsByID, byTypeNameVersion, byVersionID, byTypeName, byID goa.Endpoint) *Client {
	return &Client{
		QueryEndpoint:             query,
		ListEndpoint:              list,
		VersionsByIDEndpoint:      versionsByID,
		ByTypeNameVersionEndpoint: byTypeNameVersion,
		ByVersionIDEndpoint:       byVersionID,
		ByTypeNameEndpoint:        byTypeName,
		ByIDEndpoint:              byID,
	}
}

// Query calls the "Query" endpoint of the "resource" service.
func (c *Client) Query(ctx context.Context, p *QueryPayload) (res ResourceCollection, err error) {
	var ires interface{}
	ires, err = c.QueryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(ResourceCollection), nil
}

// List calls the "List" endpoint of the "resource" service.
func (c *Client) List(ctx context.Context, p *ListPayload) (res ResourceCollection, err error) {
	var ires interface{}
	ires, err = c.ListEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(ResourceCollection), nil
}

// VersionsByID calls the "VersionsByID" endpoint of the "resource" service.
func (c *Client) VersionsByID(ctx context.Context, p *VersionsByIDPayload) (res *Versions, err error) {
	var ires interface{}
	ires, err = c.VersionsByIDEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Versions), nil
}

// ByTypeNameVersion calls the "ByTypeNameVersion" endpoint of the "resource"
// service.
func (c *Client) ByTypeNameVersion(ctx context.Context, p *ByTypeNameVersionPayload) (res *Version, err error) {
	var ires interface{}
	ires, err = c.ByTypeNameVersionEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Version), nil
}

// ByVersionID calls the "ByVersionId" endpoint of the "resource" service.
func (c *Client) ByVersionID(ctx context.Context, p *ByVersionIDPayload) (res *Version, err error) {
	var ires interface{}
	ires, err = c.ByVersionIDEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Version), nil
}

// ByTypeName calls the "ByTypeName" endpoint of the "resource" service.
func (c *Client) ByTypeName(ctx context.Context, p *ByTypeNamePayload) (res ResourceCollection, err error) {
	var ires interface{}
	ires, err = c.ByTypeNameEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(ResourceCollection), nil
}

// ByID calls the "ById" endpoint of the "resource" service.
func (c *Client) ByID(ctx context.Context, p *ByIDPayload) (res *Resource, err error) {
	var ires interface{}
	ires, err = c.ByIDEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Resource), nil
}
