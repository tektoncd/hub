// Code generated by goa v3.17.2, DO NOT EDIT.
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
	QueryEndpoint                    goa.Endpoint
	ListEndpoint                     goa.Endpoint
	VersionsByIDEndpoint             goa.Endpoint
	ByCatalogKindNameVersionEndpoint goa.Endpoint
	ByVersionIDEndpoint              goa.Endpoint
	ByCatalogKindNameEndpoint        goa.Endpoint
	ByIDEndpoint                     goa.Endpoint
}

// NewClient initializes a "resource" service client given the endpoints.
func NewClient(query, list, versionsByID, byCatalogKindNameVersion, byVersionID, byCatalogKindName, byID goa.Endpoint) *Client {
	return &Client{
		QueryEndpoint:                    query,
		ListEndpoint:                     list,
		VersionsByIDEndpoint:             versionsByID,
		ByCatalogKindNameVersionEndpoint: byCatalogKindNameVersion,
		ByVersionIDEndpoint:              byVersionID,
		ByCatalogKindNameEndpoint:        byCatalogKindName,
		ByIDEndpoint:                     byID,
	}
}

// Query calls the "Query" endpoint of the "resource" service.
// Query may return the following errors:
//   - "internal-error" (type *goa.ServiceError): Internal Server Error
//   - "not-found" (type *goa.ServiceError): Resource Not Found Error
//   - error: internal error
func (c *Client) Query(ctx context.Context, p *QueryPayload) (res *QueryResult, err error) {
	var ires any
	ires, err = c.QueryEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*QueryResult), nil
}

// List calls the "List" endpoint of the "resource" service.
// List may return the following errors:
//   - "internal-error" (type *goa.ServiceError): Internal Server Error
//   - "not-found" (type *goa.ServiceError): Resource Not Found Error
//   - error: internal error
func (c *Client) List(ctx context.Context) (res *Resources, err error) {
	var ires any
	ires, err = c.ListEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*Resources), nil
}

// VersionsByID calls the "VersionsByID" endpoint of the "resource" service.
// VersionsByID may return the following errors:
//   - "internal-error" (type *goa.ServiceError): Internal Server Error
//   - "not-found" (type *goa.ServiceError): Resource Not Found Error
//   - error: internal error
func (c *Client) VersionsByID(ctx context.Context, p *VersionsByIDPayload) (res *VersionsByIDResult, err error) {
	var ires any
	ires, err = c.VersionsByIDEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*VersionsByIDResult), nil
}

// ByCatalogKindNameVersion calls the "ByCatalogKindNameVersion" endpoint of
// the "resource" service.
// ByCatalogKindNameVersion may return the following errors:
//   - "internal-error" (type *goa.ServiceError): Internal Server Error
//   - "not-found" (type *goa.ServiceError): Resource Not Found Error
//   - error: internal error
func (c *Client) ByCatalogKindNameVersion(ctx context.Context, p *ByCatalogKindNameVersionPayload) (res *ByCatalogKindNameVersionResult, err error) {
	var ires any
	ires, err = c.ByCatalogKindNameVersionEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ByCatalogKindNameVersionResult), nil
}

// ByVersionID calls the "ByVersionId" endpoint of the "resource" service.
// ByVersionID may return the following errors:
//   - "internal-error" (type *goa.ServiceError): Internal Server Error
//   - "not-found" (type *goa.ServiceError): Resource Not Found Error
//   - error: internal error
func (c *Client) ByVersionID(ctx context.Context, p *ByVersionIDPayload) (res *ByVersionIDResult, err error) {
	var ires any
	ires, err = c.ByVersionIDEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ByVersionIDResult), nil
}

// ByCatalogKindName calls the "ByCatalogKindName" endpoint of the "resource"
// service.
// ByCatalogKindName may return the following errors:
//   - "internal-error" (type *goa.ServiceError): Internal Server Error
//   - "not-found" (type *goa.ServiceError): Resource Not Found Error
//   - error: internal error
func (c *Client) ByCatalogKindName(ctx context.Context, p *ByCatalogKindNamePayload) (res *ByCatalogKindNameResult, err error) {
	var ires any
	ires, err = c.ByCatalogKindNameEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ByCatalogKindNameResult), nil
}

// ByID calls the "ById" endpoint of the "resource" service.
// ByID may return the following errors:
//   - "internal-error" (type *goa.ServiceError): Internal Server Error
//   - "not-found" (type *goa.ServiceError): Resource Not Found Error
//   - error: internal error
func (c *Client) ByID(ctx context.Context, p *ByIDPayload) (res *ByIDResult, err error) {
	var ires any
	ires, err = c.ByIDEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*ByIDResult), nil
}
