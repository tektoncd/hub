// Code generated by goa v3.2.0, DO NOT EDIT.
//
// resource HTTP client types
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	resourceviews "github.com/tektoncd/hub/api/gen/resource/views"
	goa "goa.design/goa/v3/pkg"
)

// QueryResponseBody is the type of the "resource" service "Query" endpoint
// HTTP response body.
type QueryResponseBody []*ResourceResponse

// QueryInternalErrorResponseBody is the type of the "resource" service "Query"
// endpoint HTTP response body for the "internal-error" error.
type QueryInternalErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// QueryNotFoundResponseBody is the type of the "resource" service "Query"
// endpoint HTTP response body for the "not-found" error.
type QueryNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// ResourceResponse is used to define fields on response body types.
type ResourceResponse struct {
	// ID is the unique id of the resource
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of resource
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// Type of catalog to which resource belongs
	Catalog *CatalogResponse `form:"catalog,omitempty" json:"catalog,omitempty" xml:"catalog,omitempty"`
	// Type of resource
	Type *string `form:"type,omitempty" json:"type,omitempty" xml:"type,omitempty"`
	// Latest version of resource
	LatestVersion *VersionResponse `form:"latestVersion,omitempty" json:"latestVersion,omitempty" xml:"latestVersion,omitempty"`
	// Tags related to resource
	Tags []*TagResponse `form:"tags,omitempty" json:"tags,omitempty" xml:"tags,omitempty"`
	// Rating of resource
	Rating *float64 `form:"rating,omitempty" json:"rating,omitempty" xml:"rating,omitempty"`
}

// CatalogResponse is used to define fields on response body types.
type CatalogResponse struct {
	// ID is the unique id of the catalog
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Type of catalog
	Type *string `form:"type,omitempty" json:"type,omitempty" xml:"type,omitempty"`
}

// VersionResponse is used to define fields on response body types.
type VersionResponse struct {
	// ID is the unique id of resource's version
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Version of resource
	Version *string `form:"version,omitempty" json:"version,omitempty" xml:"version,omitempty"`
	// Display name of version
	DisplayName *string `form:"displayName,omitempty" json:"displayName,omitempty" xml:"displayName,omitempty"`
	// Description of version
	Description *string `form:"description,omitempty" json:"description,omitempty" xml:"description,omitempty"`
	// Minimum pipelines version the resource's version is compatible with
	MinPipelinesVersion *string `form:"minPipelinesVersion,omitempty" json:"minPipelinesVersion,omitempty" xml:"minPipelinesVersion,omitempty"`
	// Raw URL of resource's yaml file of the version
	RawURL *string `form:"rawURL,omitempty" json:"rawURL,omitempty" xml:"rawURL,omitempty"`
	// Web URL of resource's yaml file of the version
	WebURL *string `form:"webURL,omitempty" json:"webURL,omitempty" xml:"webURL,omitempty"`
	// Timestamp when version was last updated
	UpdatedAt *string `form:"updatedAt,omitempty" json:"updatedAt,omitempty" xml:"updatedAt,omitempty"`
}

// TagResponse is used to define fields on response body types.
type TagResponse struct {
	// ID is the unique id of tag
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of tag
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// NewQueryResourceCollectionOK builds a "resource" service "Query" endpoint
// result from a HTTP "OK" response.
func NewQueryResourceCollectionOK(body QueryResponseBody) resourceviews.ResourceCollectionView {
	v := make([]*resourceviews.ResourceView, len(body))
	for i, val := range body {
		v[i] = unmarshalResourceResponseToResourceviewsResourceView(val)
	}
	return v
}

// NewQueryInternalError builds a resource service Query endpoint
// internal-error error.
func NewQueryInternalError(body *QueryInternalErrorResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewQueryNotFound builds a resource service Query endpoint not-found error.
func NewQueryNotFound(body *QueryNotFoundResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// ValidateQueryInternalErrorResponseBody runs the validations defined on
// Query_internal-error_Response_Body
func ValidateQueryInternalErrorResponseBody(body *QueryInternalErrorResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateQueryNotFoundResponseBody runs the validations defined on
// Query_not-found_Response_Body
func ValidateQueryNotFoundResponseBody(body *QueryNotFoundResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateResourceResponse runs the validations defined on ResourceResponse
func ValidateResourceResponse(body *ResourceResponse) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.Catalog == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("catalog", "body"))
	}
	if body.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "body"))
	}
	if body.LatestVersion == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("latestVersion", "body"))
	}
	if body.Tags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("tags", "body"))
	}
	if body.Rating == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rating", "body"))
	}
	if body.Catalog != nil {
		if err2 := ValidateCatalogResponse(body.Catalog); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if body.LatestVersion != nil {
		if err2 := ValidateVersionResponse(body.LatestVersion); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range body.Tags {
		if e != nil {
			if err2 := ValidateTagResponse(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateCatalogResponse runs the validations defined on CatalogResponse
func ValidateCatalogResponse(body *CatalogResponse) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "body"))
	}
	if body.Type != nil {
		if !(*body.Type == "official" || *body.Type == "community") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.type", *body.Type, []interface{}{"official", "community"}))
		}
	}
	return
}

// ValidateVersionResponse runs the validations defined on VersionResponse
func ValidateVersionResponse(body *VersionResponse) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Version == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("version", "body"))
	}
	if body.Description == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("description", "body"))
	}
	if body.DisplayName == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("displayName", "body"))
	}
	if body.MinPipelinesVersion == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("minPipelinesVersion", "body"))
	}
	if body.RawURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("rawURL", "body"))
	}
	if body.WebURL == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("webURL", "body"))
	}
	if body.UpdatedAt == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("updatedAt", "body"))
	}
	if body.RawURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("body.rawURL", *body.RawURL, goa.FormatURI))
	}
	if body.WebURL != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("body.webURL", *body.WebURL, goa.FormatURI))
	}
	if body.UpdatedAt != nil {
		err = goa.MergeErrors(err, goa.ValidateFormat("body.updatedAt", *body.UpdatedAt, goa.FormatDateTime))
	}
	return
}

// ValidateTagResponse runs the validations defined on TagResponse
func ValidateTagResponse(body *TagResponse) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	return
}
