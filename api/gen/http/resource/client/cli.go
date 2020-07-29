// Code generated by goa v3.2.0, DO NOT EDIT.
//
// resource HTTP client CLI support package
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"fmt"
	"strconv"

	resource "github.com/tektoncd/hub/api/gen/resource"
	goa "goa.design/goa/v3/pkg"
)

// BuildQueryPayload builds the payload for the resource Query endpoint from
// CLI flags.
func BuildQueryPayload(resourceQueryName string, resourceQueryType string, resourceQueryLimit string) (*resource.QueryPayload, error) {
	var err error
	var name string
	{
		if resourceQueryName != "" {
			name = resourceQueryName
		}
	}
	var type_ string
	{
		if resourceQueryType != "" {
			type_ = resourceQueryType
			if !(type_ == "task" || type_ == "pipeline" || type_ == "") {
				err = goa.MergeErrors(err, goa.InvalidEnumValueError("type_", type_, []interface{}{"task", "pipeline", ""}))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	var limit uint
	{
		if resourceQueryLimit != "" {
			var v uint64
			v, err = strconv.ParseUint(resourceQueryLimit, 10, 64)
			limit = uint(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for limit, must be UINT")
			}
		}
	}
	v := &resource.QueryPayload{}
	v.Name = name
	v.Type = type_
	v.Limit = limit

	return v, nil
}

// BuildListPayload builds the payload for the resource List endpoint from CLI
// flags.
func BuildListPayload(resourceListLimit string) (*resource.ListPayload, error) {
	var err error
	var limit uint
	{
		if resourceListLimit != "" {
			var v uint64
			v, err = strconv.ParseUint(resourceListLimit, 10, 64)
			limit = uint(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for limit, must be UINT")
			}
		}
	}
	v := &resource.ListPayload{}
	v.Limit = limit

	return v, nil
}

// BuildVersionsByIDPayload builds the payload for the resource VersionsByID
// endpoint from CLI flags.
func BuildVersionsByIDPayload(resourceVersionsByIDID string) (*resource.VersionsByIDPayload, error) {
	var err error
	var id uint
	{
		var v uint64
		v, err = strconv.ParseUint(resourceVersionsByIDID, 10, 64)
		id = uint(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be UINT")
		}
	}
	v := &resource.VersionsByIDPayload{}
	v.ID = id

	return v, nil
}

// BuildByTypeNameVersionPayload builds the payload for the resource
// ByTypeNameVersion endpoint from CLI flags.
func BuildByTypeNameVersionPayload(resourceByTypeNameVersionType string, resourceByTypeNameVersionName string, resourceByTypeNameVersionVersion string) (*resource.ByTypeNameVersionPayload, error) {
	var err error
	var type_ string
	{
		type_ = resourceByTypeNameVersionType
		if !(type_ == "task" || type_ == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("type_", type_, []interface{}{"task", "pipeline"}))
		}
		if err != nil {
			return nil, err
		}
	}
	var name string
	{
		name = resourceByTypeNameVersionName
	}
	var version string
	{
		version = resourceByTypeNameVersionVersion
	}
	v := &resource.ByTypeNameVersionPayload{}
	v.Type = type_
	v.Name = name
	v.Version = version

	return v, nil
}

// BuildByVersionIDPayload builds the payload for the resource ByVersionId
// endpoint from CLI flags.
func BuildByVersionIDPayload(resourceByVersionIDVersionID string) (*resource.ByVersionIDPayload, error) {
	var err error
	var versionID uint
	{
		var v uint64
		v, err = strconv.ParseUint(resourceByVersionIDVersionID, 10, 64)
		versionID = uint(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for versionID, must be UINT")
		}
	}
	v := &resource.ByVersionIDPayload{}
	v.VersionID = versionID

	return v, nil
}

// BuildByTypeNamePayload builds the payload for the resource ByTypeName
// endpoint from CLI flags.
func BuildByTypeNamePayload(resourceByTypeNameType string, resourceByTypeNameName string) (*resource.ByTypeNamePayload, error) {
	var err error
	var type_ string
	{
		type_ = resourceByTypeNameType
		if !(type_ == "task" || type_ == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("type_", type_, []interface{}{"task", "pipeline"}))
		}
		if err != nil {
			return nil, err
		}
	}
	var name string
	{
		name = resourceByTypeNameName
	}
	v := &resource.ByTypeNamePayload{}
	v.Type = type_
	v.Name = name

	return v, nil
}

// BuildByIDPayload builds the payload for the resource ById endpoint from CLI
// flags.
func BuildByIDPayload(resourceByIDID string) (*resource.ByIDPayload, error) {
	var err error
	var id uint
	{
		var v uint64
		v, err = strconv.ParseUint(resourceByIDID, 10, 64)
		id = uint(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be UINT")
		}
	}
	v := &resource.ByIDPayload{}
	v.ID = id

	return v, nil
}
