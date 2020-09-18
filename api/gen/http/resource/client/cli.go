// Code generated by goa v3.2.2, DO NOT EDIT.
//
// resource HTTP client CLI support package
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	resource "github.com/tektoncd/hub/api/gen/resource"
	goa "goa.design/goa/v3/pkg"
)

// BuildQueryPayload builds the payload for the resource Query endpoint from
// CLI flags.
func BuildQueryPayload(resourceQueryName string, resourceQueryKinds string, resourceQueryTags string, resourceQueryLimit string, resourceQueryMatch string) (*resource.QueryPayload, error) {
	var err error
	var name string
	{
		if resourceQueryName != "" {
			name = resourceQueryName
		}
	}
	var kinds []string
	{
		if resourceQueryKinds != "" {
			err = json.Unmarshal([]byte(resourceQueryKinds), &kinds)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for kinds, example of valid JSON:\n%s", "'[\n      \"Tempora omnis et nihil aut quo quidem.\",\n      \"Qui nemo sint est.\",\n      \"Nesciunt sint cupiditate.\",\n      \"Ipsum tenetur unde et amet eum hic.\"\n   ]'")
			}
		}
	}
	var tags []string
	{
		if resourceQueryTags != "" {
			err = json.Unmarshal([]byte(resourceQueryTags), &tags)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for tags, example of valid JSON:\n%s", "'[\n      \"Explicabo enim adipisci.\",\n      \"Ipsa minus ut.\"\n   ]'")
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
	var match string
	{
		if resourceQueryMatch != "" {
			match = resourceQueryMatch
			if !(match == "exact" || match == "contains") {
				err = goa.MergeErrors(err, goa.InvalidEnumValueError("match", match, []interface{}{"exact", "contains"}))
			}
			if err != nil {
				return nil, err
			}
		}
	}
	v := &resource.QueryPayload{}
	v.Name = name
	v.Kinds = kinds
	v.Tags = tags
	v.Limit = limit
	v.Match = match

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

// BuildByCatalogKindNameVersionPayload builds the payload for the resource
// ByCatalogKindNameVersion endpoint from CLI flags.
func BuildByCatalogKindNameVersionPayload(resourceByCatalogKindNameVersionCatalog string, resourceByCatalogKindNameVersionKind string, resourceByCatalogKindNameVersionName string, resourceByCatalogKindNameVersionVersion string) (*resource.ByCatalogKindNameVersionPayload, error) {
	var err error
	var catalog string
	{
		catalog = resourceByCatalogKindNameVersionCatalog
	}
	var kind string
	{
		kind = resourceByCatalogKindNameVersionKind
		if !(kind == "task" || kind == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("kind", kind, []interface{}{"task", "pipeline"}))
		}
		if err != nil {
			return nil, err
		}
	}
	var name string
	{
		name = resourceByCatalogKindNameVersionName
	}
	var version string
	{
		version = resourceByCatalogKindNameVersionVersion
	}
	v := &resource.ByCatalogKindNameVersionPayload{}
	v.Catalog = catalog
	v.Kind = kind
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

// BuildByCatalogKindNamePayload builds the payload for the resource
// ByCatalogKindName endpoint from CLI flags.
func BuildByCatalogKindNamePayload(resourceByCatalogKindNameCatalog string, resourceByCatalogKindNameKind string, resourceByCatalogKindNameName string) (*resource.ByCatalogKindNamePayload, error) {
	var err error
	var catalog string
	{
		catalog = resourceByCatalogKindNameCatalog
	}
	var kind string
	{
		kind = resourceByCatalogKindNameKind
		if !(kind == "task" || kind == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("kind", kind, []interface{}{"task", "pipeline"}))
		}
		if err != nil {
			return nil, err
		}
	}
	var name string
	{
		name = resourceByCatalogKindNameName
	}
	v := &resource.ByCatalogKindNamePayload{}
	v.Catalog = catalog
	v.Kind = kind
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
