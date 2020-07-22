// Code generated by goa v3.2.0, DO NOT EDIT.
//
// resource HTTP server types
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package server

import (
	resource "github.com/tektoncd/hub/api/gen/resource"
	resourceviews "github.com/tektoncd/hub/api/gen/resource/views"
	goa "goa.design/goa/v3/pkg"
)

// ResourceResponseCollection is the type of the "resource" service "Query"
// endpoint HTTP response body.
type ResourceResponseCollection []*ResourceResponse

// VersionsByIDResponseBody is the type of the "resource" service
// "VersionsByID" endpoint HTTP response body.
type VersionsByIDResponseBody struct {
	// Latest Version of resource
	Latest *VersionResponseBodyUrls `form:"latest" json:"latest" xml:"latest"`
	// List of all versions of resource
	Versions []*VersionResponseBodyUrls `form:"versions" json:"versions" xml:"versions"`
}

// ByTypeNameVersionResponseBody is the type of the "resource" service
// "ByTypeNameVersion" endpoint HTTP response body.
type ByTypeNameVersionResponseBody struct {
	// ID is the unique id of resource's version
	ID uint `form:"id" json:"id" xml:"id"`
	// Version of resource
	Version string `form:"version" json:"version" xml:"version"`
	// Description of version
	Description string `form:"description" json:"description" xml:"description"`
	// Minimum pipelines version the resource's version is compatible with
	MinPipelinesVersion string `form:"minPipelinesVersion" json:"minPipelinesVersion" xml:"minPipelinesVersion"`
	// Display name of version
	DisplayName string `form:"displayName" json:"displayName" xml:"displayName"`
	// Raw URL of resource's yaml file of the version
	RawURL string `form:"rawURL" json:"rawURL" xml:"rawURL"`
	// Web URL of resource's yaml file of the version
	WebURL string `form:"webURL" json:"webURL" xml:"webURL"`
	// Timestamp when version was last updated
	UpdatedAt string `form:"updatedAt" json:"updatedAt" xml:"updatedAt"`
	// Resource to which the version belongs
	Resource *ResourceResponseBodyInfo `form:"resource" json:"resource" xml:"resource"`
}

// ByVersionIDResponseBody is the type of the "resource" service "ByVersionId"
// endpoint HTTP response body.
type ByVersionIDResponseBody struct {
	// ID is the unique id of resource's version
	ID uint `form:"id" json:"id" xml:"id"`
	// Version of resource
	Version string `form:"version" json:"version" xml:"version"`
	// Description of version
	Description string `form:"description" json:"description" xml:"description"`
	// Minimum pipelines version the resource's version is compatible with
	MinPipelinesVersion string `form:"minPipelinesVersion" json:"minPipelinesVersion" xml:"minPipelinesVersion"`
	// Display name of version
	DisplayName string `form:"displayName" json:"displayName" xml:"displayName"`
	// Raw URL of resource's yaml file of the version
	RawURL string `form:"rawURL" json:"rawURL" xml:"rawURL"`
	// Web URL of resource's yaml file of the version
	WebURL string `form:"webURL" json:"webURL" xml:"webURL"`
	// Timestamp when version was last updated
	UpdatedAt string `form:"updatedAt" json:"updatedAt" xml:"updatedAt"`
	// Resource to which the version belongs
	Resource *ResourceResponseBodyInfo `form:"resource" json:"resource" xml:"resource"`
}

// QueryInternalErrorResponseBody is the type of the "resource" service "Query"
// endpoint HTTP response body for the "internal-error" error.
type QueryInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// QueryNotFoundResponseBody is the type of the "resource" service "Query"
// endpoint HTTP response body for the "not-found" error.
type QueryNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ListInternalErrorResponseBody is the type of the "resource" service "List"
// endpoint HTTP response body for the "internal-error" error.
type ListInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// VersionsByIDInternalErrorResponseBody is the type of the "resource" service
// "VersionsByID" endpoint HTTP response body for the "internal-error" error.
type VersionsByIDInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// VersionsByIDNotFoundResponseBody is the type of the "resource" service
// "VersionsByID" endpoint HTTP response body for the "not-found" error.
type VersionsByIDNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ByTypeNameVersionInternalErrorResponseBody is the type of the "resource"
// service "ByTypeNameVersion" endpoint HTTP response body for the
// "internal-error" error.
type ByTypeNameVersionInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ByTypeNameVersionNotFoundResponseBody is the type of the "resource" service
// "ByTypeNameVersion" endpoint HTTP response body for the "not-found" error.
type ByTypeNameVersionNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ByVersionIDInternalErrorResponseBody is the type of the "resource" service
// "ByVersionId" endpoint HTTP response body for the "internal-error" error.
type ByVersionIDInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ByVersionIDNotFoundResponseBody is the type of the "resource" service
// "ByVersionId" endpoint HTTP response body for the "not-found" error.
type ByVersionIDNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ByTypeNameInternalErrorResponseBody is the type of the "resource" service
// "ByTypeName" endpoint HTTP response body for the "internal-error" error.
type ByTypeNameInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ByTypeNameNotFoundResponseBody is the type of the "resource" service
// "ByTypeName" endpoint HTTP response body for the "not-found" error.
type ByTypeNameNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// ResourceResponse is used to define fields on response body types.
type ResourceResponse struct {
	// ID is the unique id of the resource
	ID uint `form:"id" json:"id" xml:"id"`
	// Name of resource
	Name string `form:"name" json:"name" xml:"name"`
	// Type of catalog to which resource belongs
	Catalog *CatalogResponse `form:"catalog" json:"catalog" xml:"catalog"`
	// Type of resource
	Type string `form:"type" json:"type" xml:"type"`
	// Latest version of resource
	LatestVersion *LatestVersionResponse `form:"latestVersion" json:"latestVersion" xml:"latestVersion"`
	// Tags related to resource
	Tags []*TagResponse `form:"tags" json:"tags" xml:"tags"`
	// Rating of resource
	Rating float64 `form:"rating" json:"rating" xml:"rating"`
}

// CatalogResponse is used to define fields on response body types.
type CatalogResponse struct {
	// ID is the unique id of the catalog
	ID uint `form:"id" json:"id" xml:"id"`
	// Type of catalog
	Type string `form:"type" json:"type" xml:"type"`
}

// LatestVersionResponse is used to define fields on response body types.
type LatestVersionResponse struct {
	// ID is the unique id of resource's version
	ID uint `form:"id" json:"id" xml:"id"`
	// Version of resource
	Version string `form:"version" json:"version" xml:"version"`
	// Display name of version
	DisplayName string `form:"displayName" json:"displayName" xml:"displayName"`
	// Description of version
	Description string `form:"description" json:"description" xml:"description"`
	// Minimum pipelines version the resource's version is compatible with
	MinPipelinesVersion string `form:"minPipelinesVersion" json:"minPipelinesVersion" xml:"minPipelinesVersion"`
	// Raw URL of resource's yaml file of the version
	RawURL string `form:"rawURL" json:"rawURL" xml:"rawURL"`
	// Web URL of resource's yaml file of the version
	WebURL string `form:"webURL" json:"webURL" xml:"webURL"`
	// Timestamp when version was last updated
	UpdatedAt string `form:"updatedAt" json:"updatedAt" xml:"updatedAt"`
}

// TagResponse is used to define fields on response body types.
type TagResponse struct {
	// ID is the unique id of tag
	ID uint `form:"id" json:"id" xml:"id"`
	// Name of tag
	Name string `form:"name" json:"name" xml:"name"`
}

// VersionResponseBodyUrls is used to define fields on response body types.
type VersionResponseBodyUrls struct {
	// ID is the unique id of resource's version
	ID uint `form:"id" json:"id" xml:"id"`
	// Version of resource
	Version string `form:"version" json:"version" xml:"version"`
	// Raw URL of resource's yaml file of the version
	RawURL string `form:"rawURL" json:"rawURL" xml:"rawURL"`
	// Web URL of resource's yaml file of the version
	WebURL string `form:"webURL" json:"webURL" xml:"webURL"`
}

// ResourceResponseBodyInfo is used to define fields on response body types.
type ResourceResponseBodyInfo struct {
	// ID is the unique id of the resource
	ID uint `form:"id" json:"id" xml:"id"`
	// Name of resource
	Name string `form:"name" json:"name" xml:"name"`
	// Type of catalog to which resource belongs
	Catalog *CatalogResponseBody `form:"catalog" json:"catalog" xml:"catalog"`
	// Type of resource
	Type string `form:"type" json:"type" xml:"type"`
	// Tags related to resource
	Tags []*TagResponseBody `form:"tags" json:"tags" xml:"tags"`
	// Rating of resource
	Rating float64 `form:"rating" json:"rating" xml:"rating"`
}

// CatalogResponseBody is used to define fields on response body types.
type CatalogResponseBody struct {
	// ID is the unique id of the catalog
	ID uint `form:"id" json:"id" xml:"id"`
	// Type of catalog
	Type string `form:"type" json:"type" xml:"type"`
}

// TagResponseBody is used to define fields on response body types.
type TagResponseBody struct {
	// ID is the unique id of tag
	ID uint `form:"id" json:"id" xml:"id"`
	// Name of tag
	Name string `form:"name" json:"name" xml:"name"`
}

// NewResourceResponseCollection builds the HTTP response body from the result
// of the "Query" endpoint of the "resource" service.
func NewResourceResponseCollection(res resourceviews.ResourceCollectionView) ResourceResponseCollection {
	body := make([]*ResourceResponse, len(res))
	for i, val := range res {
		body[i] = marshalResourceviewsResourceViewToResourceResponse(val)
	}
	return body
}

// NewVersionsByIDResponseBody builds the HTTP response body from the result of
// the "VersionsByID" endpoint of the "resource" service.
func NewVersionsByIDResponseBody(res *resourceviews.VersionsView) *VersionsByIDResponseBody {
	body := &VersionsByIDResponseBody{}
	if res.Latest != nil {
		body.Latest = marshalResourceviewsVersionViewToVersionResponseBodyUrls(res.Latest)
	}
	if res.Versions != nil {
		body.Versions = make([]*VersionResponseBodyUrls, len(res.Versions))
		for i, val := range res.Versions {
			body.Versions[i] = marshalResourceviewsVersionViewToVersionResponseBodyUrls(val)
		}
	}
	return body
}

// NewByTypeNameVersionResponseBody builds the HTTP response body from the
// result of the "ByTypeNameVersion" endpoint of the "resource" service.
func NewByTypeNameVersionResponseBody(res *resourceviews.VersionView) *ByTypeNameVersionResponseBody {
	body := &ByTypeNameVersionResponseBody{
		ID:                  *res.ID,
		Version:             *res.Version,
		DisplayName:         *res.DisplayName,
		Description:         *res.Description,
		MinPipelinesVersion: *res.MinPipelinesVersion,
		RawURL:              *res.RawURL,
		WebURL:              *res.WebURL,
		UpdatedAt:           *res.UpdatedAt,
	}
	if res.Resource != nil {
		body.Resource = marshalResourceviewsResourceViewToResourceResponseBodyInfo(res.Resource)
	}
	return body
}

// NewByVersionIDResponseBody builds the HTTP response body from the result of
// the "ByVersionId" endpoint of the "resource" service.
func NewByVersionIDResponseBody(res *resourceviews.VersionView) *ByVersionIDResponseBody {
	body := &ByVersionIDResponseBody{
		ID:                  *res.ID,
		Version:             *res.Version,
		DisplayName:         *res.DisplayName,
		Description:         *res.Description,
		MinPipelinesVersion: *res.MinPipelinesVersion,
		RawURL:              *res.RawURL,
		WebURL:              *res.WebURL,
		UpdatedAt:           *res.UpdatedAt,
	}
	if res.Resource != nil {
		body.Resource = marshalResourceviewsResourceViewToResourceResponseBodyInfo(res.Resource)
	}
	return body
}

// NewQueryInternalErrorResponseBody builds the HTTP response body from the
// result of the "Query" endpoint of the "resource" service.
func NewQueryInternalErrorResponseBody(res *goa.ServiceError) *QueryInternalErrorResponseBody {
	body := &QueryInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewQueryNotFoundResponseBody builds the HTTP response body from the result
// of the "Query" endpoint of the "resource" service.
func NewQueryNotFoundResponseBody(res *goa.ServiceError) *QueryNotFoundResponseBody {
	body := &QueryNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewListInternalErrorResponseBody builds the HTTP response body from the
// result of the "List" endpoint of the "resource" service.
func NewListInternalErrorResponseBody(res *goa.ServiceError) *ListInternalErrorResponseBody {
	body := &ListInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewVersionsByIDInternalErrorResponseBody builds the HTTP response body from
// the result of the "VersionsByID" endpoint of the "resource" service.
func NewVersionsByIDInternalErrorResponseBody(res *goa.ServiceError) *VersionsByIDInternalErrorResponseBody {
	body := &VersionsByIDInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewVersionsByIDNotFoundResponseBody builds the HTTP response body from the
// result of the "VersionsByID" endpoint of the "resource" service.
func NewVersionsByIDNotFoundResponseBody(res *goa.ServiceError) *VersionsByIDNotFoundResponseBody {
	body := &VersionsByIDNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewByTypeNameVersionInternalErrorResponseBody builds the HTTP response body
// from the result of the "ByTypeNameVersion" endpoint of the "resource"
// service.
func NewByTypeNameVersionInternalErrorResponseBody(res *goa.ServiceError) *ByTypeNameVersionInternalErrorResponseBody {
	body := &ByTypeNameVersionInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewByTypeNameVersionNotFoundResponseBody builds the HTTP response body from
// the result of the "ByTypeNameVersion" endpoint of the "resource" service.
func NewByTypeNameVersionNotFoundResponseBody(res *goa.ServiceError) *ByTypeNameVersionNotFoundResponseBody {
	body := &ByTypeNameVersionNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewByVersionIDInternalErrorResponseBody builds the HTTP response body from
// the result of the "ByVersionId" endpoint of the "resource" service.
func NewByVersionIDInternalErrorResponseBody(res *goa.ServiceError) *ByVersionIDInternalErrorResponseBody {
	body := &ByVersionIDInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewByVersionIDNotFoundResponseBody builds the HTTP response body from the
// result of the "ByVersionId" endpoint of the "resource" service.
func NewByVersionIDNotFoundResponseBody(res *goa.ServiceError) *ByVersionIDNotFoundResponseBody {
	body := &ByVersionIDNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewByTypeNameInternalErrorResponseBody builds the HTTP response body from
// the result of the "ByTypeName" endpoint of the "resource" service.
func NewByTypeNameInternalErrorResponseBody(res *goa.ServiceError) *ByTypeNameInternalErrorResponseBody {
	body := &ByTypeNameInternalErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewByTypeNameNotFoundResponseBody builds the HTTP response body from the
// result of the "ByTypeName" endpoint of the "resource" service.
func NewByTypeNameNotFoundResponseBody(res *goa.ServiceError) *ByTypeNameNotFoundResponseBody {
	body := &ByTypeNameNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewQueryPayload builds a resource service Query endpoint payload.
func NewQueryPayload(name string, type_ string, limit uint) *resource.QueryPayload {
	v := &resource.QueryPayload{}
	v.Name = name
	v.Type = type_
	v.Limit = limit

	return v
}

// NewListPayload builds a resource service List endpoint payload.
func NewListPayload(limit uint) *resource.ListPayload {
	v := &resource.ListPayload{}
	v.Limit = limit

	return v
}

// NewVersionsByIDPayload builds a resource service VersionsByID endpoint
// payload.
func NewVersionsByIDPayload(id uint) *resource.VersionsByIDPayload {
	v := &resource.VersionsByIDPayload{}
	v.ID = id

	return v
}

// NewByTypeNameVersionPayload builds a resource service ByTypeNameVersion
// endpoint payload.
func NewByTypeNameVersionPayload(type_ string, name string, version string) *resource.ByTypeNameVersionPayload {
	v := &resource.ByTypeNameVersionPayload{}
	v.Type = type_
	v.Name = name
	v.Version = version

	return v
}

// NewByVersionIDPayload builds a resource service ByVersionId endpoint payload.
func NewByVersionIDPayload(versionID uint) *resource.ByVersionIDPayload {
	v := &resource.ByVersionIDPayload{}
	v.VersionID = versionID

	return v
}

// NewByTypeNamePayload builds a resource service ByTypeName endpoint payload.
func NewByTypeNamePayload(type_ string, name string) *resource.ByTypeNamePayload {
	v := &resource.ByTypeNamePayload{}
	v.Type = type_
	v.Name = name

	return v
}
