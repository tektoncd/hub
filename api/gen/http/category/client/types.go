// Code generated by goa v3.2.2, DO NOT EDIT.
//
// category HTTP client types
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	category "github.com/tektoncd/hub/api/gen/category"
	goa "goa.design/goa/v3/pkg"
)

// ListResponseBody is the type of the "category" service "list" endpoint HTTP
// response body.
type ListResponseBody []*CategoryResponse

// ListInternalErrorResponseBody is the type of the "category" service "list"
// endpoint HTTP response body for the "internal-error" error.
type ListInternalErrorResponseBody struct {
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

// CategoryResponse is used to define fields on response body types.
type CategoryResponse struct {
	// ID is the unique id of the category
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of category
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// List of tags associated with the category
	Tags []*TagResponse `form:"tags,omitempty" json:"tags,omitempty" xml:"tags,omitempty"`
}

// TagResponse is used to define fields on response body types.
type TagResponse struct {
	// ID is the unique id of tag
	ID *uint `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Name of tag
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
}

// NewListCategoryOK builds a "category" service "list" endpoint result from a
// HTTP "OK" response.
func NewListCategoryOK(body []*CategoryResponse) []*category.Category {
	v := make([]*category.Category, len(body))
	for i, val := range body {
		v[i] = unmarshalCategoryResponseToCategoryCategory(val)
	}
	return v
}

// NewListInternalError builds a category service list endpoint internal-error
// error.
func NewListInternalError(body *ListInternalErrorResponseBody) *goa.ServiceError {
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

// ValidateListInternalErrorResponseBody runs the validations defined on
// list_internal-error_response_body
func ValidateListInternalErrorResponseBody(body *ListInternalErrorResponseBody) (err error) {
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

// ValidateCategoryResponse runs the validations defined on CategoryResponse
func ValidateCategoryResponse(body *CategoryResponse) (err error) {
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.Tags == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("tags", "body"))
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
