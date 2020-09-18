// Code generated by goa v3.2.2, DO NOT EDIT.
//
// admin HTTP client types
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	admin "github.com/tektoncd/hub/api/gen/admin"
	goa "goa.design/goa/v3/pkg"
)

// UpdateAgentRequestBody is the type of the "admin" service "UpdateAgent"
// endpoint HTTP request body.
type UpdateAgentRequestBody struct {
	// Name of Agent
	Name string `form:"name" json:"name" xml:"name"`
	// Scopes required for Agent
	Scopes []string `form:"scopes" json:"scopes" xml:"scopes"`
}

// UpdateAgentResponseBody is the type of the "admin" service "UpdateAgent"
// endpoint HTTP response body.
type UpdateAgentResponseBody struct {
	// Agent JWT
	Token *string `form:"token,omitempty" json:"token,omitempty" xml:"token,omitempty"`
}

// UpdateAgentInvalidPayloadResponseBody is the type of the "admin" service
// "UpdateAgent" endpoint HTTP response body for the "invalid-payload" error.
type UpdateAgentInvalidPayloadResponseBody struct {
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

// UpdateAgentInvalidTokenResponseBody is the type of the "admin" service
// "UpdateAgent" endpoint HTTP response body for the "invalid-token" error.
type UpdateAgentInvalidTokenResponseBody struct {
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

// UpdateAgentInvalidScopesResponseBody is the type of the "admin" service
// "UpdateAgent" endpoint HTTP response body for the "invalid-scopes" error.
type UpdateAgentInvalidScopesResponseBody struct {
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

// UpdateAgentInternalErrorResponseBody is the type of the "admin" service
// "UpdateAgent" endpoint HTTP response body for the "internal-error" error.
type UpdateAgentInternalErrorResponseBody struct {
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

// NewUpdateAgentRequestBody builds the HTTP request body from the payload of
// the "UpdateAgent" endpoint of the "admin" service.
func NewUpdateAgentRequestBody(p *admin.UpdateAgentPayload) *UpdateAgentRequestBody {
	body := &UpdateAgentRequestBody{
		Name: p.Name,
	}
	if p.Scopes != nil {
		body.Scopes = make([]string, len(p.Scopes))
		for i, val := range p.Scopes {
			body.Scopes[i] = val
		}
	}
	return body
}

// NewUpdateAgentResultOK builds a "admin" service "UpdateAgent" endpoint
// result from a HTTP "OK" response.
func NewUpdateAgentResultOK(body *UpdateAgentResponseBody) *admin.UpdateAgentResult {
	v := &admin.UpdateAgentResult{
		Token: *body.Token,
	}

	return v
}

// NewUpdateAgentInvalidPayload builds a admin service UpdateAgent endpoint
// invalid-payload error.
func NewUpdateAgentInvalidPayload(body *UpdateAgentInvalidPayloadResponseBody) *goa.ServiceError {
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

// NewUpdateAgentInvalidToken builds a admin service UpdateAgent endpoint
// invalid-token error.
func NewUpdateAgentInvalidToken(body *UpdateAgentInvalidTokenResponseBody) *goa.ServiceError {
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

// NewUpdateAgentInvalidScopes builds a admin service UpdateAgent endpoint
// invalid-scopes error.
func NewUpdateAgentInvalidScopes(body *UpdateAgentInvalidScopesResponseBody) *goa.ServiceError {
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

// NewUpdateAgentInternalError builds a admin service UpdateAgent endpoint
// internal-error error.
func NewUpdateAgentInternalError(body *UpdateAgentInternalErrorResponseBody) *goa.ServiceError {
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

// ValidateUpdateAgentResponseBody runs the validations defined on
// UpdateAgentResponseBody
func ValidateUpdateAgentResponseBody(body *UpdateAgentResponseBody) (err error) {
	if body.Token == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("token", "body"))
	}
	return
}

// ValidateUpdateAgentInvalidPayloadResponseBody runs the validations defined
// on UpdateAgent_invalid-payload_Response_Body
func ValidateUpdateAgentInvalidPayloadResponseBody(body *UpdateAgentInvalidPayloadResponseBody) (err error) {
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

// ValidateUpdateAgentInvalidTokenResponseBody runs the validations defined on
// UpdateAgent_invalid-token_Response_Body
func ValidateUpdateAgentInvalidTokenResponseBody(body *UpdateAgentInvalidTokenResponseBody) (err error) {
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

// ValidateUpdateAgentInvalidScopesResponseBody runs the validations defined on
// UpdateAgent_invalid-scopes_Response_Body
func ValidateUpdateAgentInvalidScopesResponseBody(body *UpdateAgentInvalidScopesResponseBody) (err error) {
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

// ValidateUpdateAgentInternalErrorResponseBody runs the validations defined on
// UpdateAgent_internal-error_Response_Body
func ValidateUpdateAgentInternalErrorResponseBody(body *UpdateAgentInternalErrorResponseBody) (err error) {
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
