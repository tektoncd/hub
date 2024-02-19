// Code generated by goa v3.15.0, DO NOT EDIT.
//
// catalog service
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/v1/design

package catalog

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// The catalog service provides details about catalogs.
type Service interface {
	// List all Catalogs
	List(context.Context) (res *ListResult, err error)
}

// APIName is the name of the API as defined in the design.
const APIName = "v1"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "1.0"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "catalog"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"List"}

type Catalog struct {
	// ID is the unique id of the catalog
	ID uint
	// Name of catalog
	Name string
	// Type of catalog
	Type string
	// URL of catalog
	URL string
	// Provider of catalog
	Provider string
}

// ListResult is the result type of the catalog service List method.
type ListResult struct {
	Data []*Catalog
}

// MakeInternalError builds a goa.ServiceError from an error.
func MakeInternalError(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "internal-error", false, false, false)
}
