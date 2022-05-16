// Code generated by goa v3.7.5, DO NOT EDIT.
//
// rating service
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package rating

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// The rating service exposes endpoints to read and write user's rating for
// resources
type Service interface {
	// Find user's rating for a resource
	Get(context.Context, *GetPayload) (res *GetResult, err error)
	// Update user's rating for a resource
	Update(context.Context, *UpdatePayload) (err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "rating"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"Get", "Update"}

// GetPayload is the payload type of the rating service Get method.
type GetPayload struct {
	// ID of a resource
	ID uint
	// JWT
	Token string
}

// GetResult is the result type of the rating service Get method.
type GetResult struct {
	// User rating for resource
	Rating int
}

// UpdatePayload is the payload type of the rating service Update method.
type UpdatePayload struct {
	// ID of a resource
	ID uint
	// User rating for resource
	Rating uint
	// JWT
	Token string
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not-found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInternalError builds a goa.ServiceError from an error.
func MakeInternalError(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "internal-error",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInvalidToken builds a goa.ServiceError from an error.
func MakeInvalidToken(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "invalid-token",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInvalidScopes builds a goa.ServiceError from an error.
func MakeInvalidScopes(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "invalid-scopes",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
