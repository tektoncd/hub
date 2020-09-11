// Code generated by goa v3.2.2, DO NOT EDIT.
//
// resource service
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package resource

import (
	"context"

	resourceviews "github.com/tektoncd/hub/api/gen/resource/views"
	goa "goa.design/goa/v3/pkg"
)

// The resource service provides details about all kind of resources
type Service interface {
	// Find resources by a combination of name, kind and tags
	Query(context.Context, *QueryPayload) (res ResourceCollection, err error)
	// List all resources sorted by rating and name
	List(context.Context, *ListPayload) (res ResourceCollection, err error)
	// Find all versions of a resource by its id
	VersionsByID(context.Context, *VersionsByIDPayload) (res *Versions, err error)
	// Find resource using name, kind and version of resource
	ByKindNameVersion(context.Context, *ByKindNameVersionPayload) (res *Version, err error)
	// Find a resource using its version's id
	ByVersionID(context.Context, *ByVersionIDPayload) (res *Version, err error)
	// Find resources using name and kind
	ByKindName(context.Context, *ByKindNamePayload) (res ResourceCollection, err error)
	// Find a resource using it's id
	ByID(context.Context, *ByIDPayload) (res *Resource, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "resource"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [7]string{"Query", "List", "VersionsByID", "ByKindNameVersion", "ByVersionId", "ByKindName", "ById"}

// QueryPayload is the payload type of the resource service Query method.
type QueryPayload struct {
	// Name of resource
	Name string
	// Kinds of resource to filter by
	Kinds []string
	// Tags associated with a resource to filter by
	Tags []string
	// Maximum number of resources to be returned
	Limit uint
	// Strategy used to find matching resources
	Match string
}

// ResourceCollection is the result type of the resource service Query method.
type ResourceCollection []*Resource

// ListPayload is the payload type of the resource service List method.
type ListPayload struct {
	// Maximum number of resources to be returned
	Limit uint
}

// VersionsByIDPayload is the payload type of the resource service VersionsByID
// method.
type VersionsByIDPayload struct {
	// ID of a resource
	ID uint
}

// Versions is the result type of the resource service VersionsByID method.
type Versions struct {
	// Latest Version of resource
	Latest *Version
	// List of all versions of resource
	Versions []*Version
}

// ByKindNameVersionPayload is the payload type of the resource service
// ByKindNameVersion method.
type ByKindNameVersionPayload struct {
	// kind of resource
	Kind string
	// name of resource
	Name string
	// version of resource
	Version string
}

// Version is the result type of the resource service ByKindNameVersion method.
type Version struct {
	// ID is the unique id of resource's version
	ID uint
	// Version of resource
	Version string
	// Display name of version
	DisplayName string
	// Description of version
	Description string
	// Minimum pipelines version the resource's version is compatible with
	MinPipelinesVersion string
	// Raw URL of resource's yaml file of the version
	RawURL string
	// Web URL of resource's yaml file of the version
	WebURL string
	// Timestamp when version was last updated
	UpdatedAt string
	// Resource to which the version belongs
	Resource *Resource
}

// ByVersionIDPayload is the payload type of the resource service ByVersionId
// method.
type ByVersionIDPayload struct {
	// Version ID of a resource's version
	VersionID uint
}

// ByKindNamePayload is the payload type of the resource service ByKindName
// method.
type ByKindNamePayload struct {
	// kind of resource
	Kind string
	// Name of resource
	Name string
}

// ByIDPayload is the payload type of the resource service ById method.
type ByIDPayload struct {
	// ID of a resource
	ID uint
}

// Resource is the result type of the resource service ById method.
type Resource struct {
	// ID is the unique id of the resource
	ID uint
	// Name of resource
	Name string
	// Type of catalog to which resource belongs
	Catalog *Catalog
	// Kind of resource
	Kind string
	// Latest version of resource
	LatestVersion *Version
	// Tags related to resource
	Tags []*Tag
	// Rating of resource
	Rating float64
	// List of all versions of a resource
	Versions []*Version
}

type Catalog struct {
	// ID is the unique id of the catalog
	ID uint
	// Type of catalog
	Type string
}

type Tag struct {
	// ID is the unique id of tag
	ID uint
	// Name of tag
	Name string
}

// MakeInternalError builds a goa.ServiceError from an error.
func MakeInternalError(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "internal-error",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "not-found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// NewResourceCollection initializes result type ResourceCollection from viewed
// result type ResourceCollection.
func NewResourceCollection(vres resourceviews.ResourceCollection) ResourceCollection {
	var res ResourceCollection
	switch vres.View {
	case "info":
		res = newResourceCollectionInfo(vres.Projected)
	case "withoutVersion":
		res = newResourceCollectionWithoutVersion(vres.Projected)
	case "default", "":
		res = newResourceCollection(vres.Projected)
	}
	return res
}

// NewViewedResourceCollection initializes viewed result type
// ResourceCollection from result type ResourceCollection using the given view.
func NewViewedResourceCollection(res ResourceCollection, view string) resourceviews.ResourceCollection {
	var vres resourceviews.ResourceCollection
	switch view {
	case "info":
		p := newResourceCollectionViewInfo(res)
		vres = resourceviews.ResourceCollection{Projected: p, View: "info"}
	case "withoutVersion":
		p := newResourceCollectionViewWithoutVersion(res)
		vres = resourceviews.ResourceCollection{Projected: p, View: "withoutVersion"}
	case "default", "":
		p := newResourceCollectionView(res)
		vres = resourceviews.ResourceCollection{Projected: p, View: "default"}
	}
	return vres
}

// NewVersions initializes result type Versions from viewed result type
// Versions.
func NewVersions(vres *resourceviews.Versions) *Versions {
	return newVersions(vres.Projected)
}

// NewViewedVersions initializes viewed result type Versions from result type
// Versions using the given view.
func NewViewedVersions(res *Versions, view string) *resourceviews.Versions {
	p := newVersionsView(res)
	return &resourceviews.Versions{Projected: p, View: "default"}
}

// NewVersion initializes result type Version from viewed result type Version.
func NewVersion(vres *resourceviews.Version) *Version {
	var res *Version
	switch vres.View {
	case "tiny":
		res = newVersionTiny(vres.Projected)
	case "min":
		res = newVersionMin(vres.Projected)
	case "withoutResource":
		res = newVersionWithoutResource(vres.Projected)
	case "default", "":
		res = newVersion(vres.Projected)
	}
	return res
}

// NewViewedVersion initializes viewed result type Version from result type
// Version using the given view.
func NewViewedVersion(res *Version, view string) *resourceviews.Version {
	var vres *resourceviews.Version
	switch view {
	case "tiny":
		p := newVersionViewTiny(res)
		vres = &resourceviews.Version{Projected: p, View: "tiny"}
	case "min":
		p := newVersionViewMin(res)
		vres = &resourceviews.Version{Projected: p, View: "min"}
	case "withoutResource":
		p := newVersionViewWithoutResource(res)
		vres = &resourceviews.Version{Projected: p, View: "withoutResource"}
	case "default", "":
		p := newVersionView(res)
		vres = &resourceviews.Version{Projected: p, View: "default"}
	}
	return vres
}

// NewResource initializes result type Resource from viewed result type
// Resource.
func NewResource(vres *resourceviews.Resource) *Resource {
	var res *Resource
	switch vres.View {
	case "info":
		res = newResourceInfo(vres.Projected)
	case "withoutVersion":
		res = newResourceWithoutVersion(vres.Projected)
	case "default", "":
		res = newResource(vres.Projected)
	}
	return res
}

// NewViewedResource initializes viewed result type Resource from result type
// Resource using the given view.
func NewViewedResource(res *Resource, view string) *resourceviews.Resource {
	var vres *resourceviews.Resource
	switch view {
	case "info":
		p := newResourceViewInfo(res)
		vres = &resourceviews.Resource{Projected: p, View: "info"}
	case "withoutVersion":
		p := newResourceViewWithoutVersion(res)
		vres = &resourceviews.Resource{Projected: p, View: "withoutVersion"}
	case "default", "":
		p := newResourceView(res)
		vres = &resourceviews.Resource{Projected: p, View: "default"}
	}
	return vres
}

// newResourceCollectionInfo converts projected type ResourceCollection to
// service type ResourceCollection.
func newResourceCollectionInfo(vres resourceviews.ResourceCollectionView) ResourceCollection {
	res := make(ResourceCollection, len(vres))
	for i, n := range vres {
		res[i] = newResourceInfo(n)
	}
	return res
}

// newResourceCollectionWithoutVersion converts projected type
// ResourceCollection to service type ResourceCollection.
func newResourceCollectionWithoutVersion(vres resourceviews.ResourceCollectionView) ResourceCollection {
	res := make(ResourceCollection, len(vres))
	for i, n := range vres {
		res[i] = newResourceWithoutVersion(n)
	}
	return res
}

// newResourceCollection converts projected type ResourceCollection to service
// type ResourceCollection.
func newResourceCollection(vres resourceviews.ResourceCollectionView) ResourceCollection {
	res := make(ResourceCollection, len(vres))
	for i, n := range vres {
		res[i] = newResource(n)
	}
	return res
}

// newResourceCollectionViewInfo projects result type ResourceCollection to
// projected type ResourceCollectionView using the "info" view.
func newResourceCollectionViewInfo(res ResourceCollection) resourceviews.ResourceCollectionView {
	vres := make(resourceviews.ResourceCollectionView, len(res))
	for i, n := range res {
		vres[i] = newResourceViewInfo(n)
	}
	return vres
}

// newResourceCollectionViewWithoutVersion projects result type
// ResourceCollection to projected type ResourceCollectionView using the
// "withoutVersion" view.
func newResourceCollectionViewWithoutVersion(res ResourceCollection) resourceviews.ResourceCollectionView {
	vres := make(resourceviews.ResourceCollectionView, len(res))
	for i, n := range res {
		vres[i] = newResourceViewWithoutVersion(n)
	}
	return vres
}

// newResourceCollectionView projects result type ResourceCollection to
// projected type ResourceCollectionView using the "default" view.
func newResourceCollectionView(res ResourceCollection) resourceviews.ResourceCollectionView {
	vres := make(resourceviews.ResourceCollectionView, len(res))
	for i, n := range res {
		vres[i] = newResourceView(n)
	}
	return vres
}

// newResourceInfo converts projected type Resource to service type Resource.
func newResourceInfo(vres *resourceviews.ResourceView) *Resource {
	res := &Resource{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.Kind != nil {
		res.Kind = *vres.Kind
	}
	if vres.Rating != nil {
		res.Rating = *vres.Rating
	}
	if vres.Catalog != nil {
		res.Catalog = transformResourceviewsCatalogViewToCatalog(vres.Catalog)
	}
	if vres.Tags != nil {
		res.Tags = make([]*Tag, len(vres.Tags))
		for i, val := range vres.Tags {
			res.Tags[i] = transformResourceviewsTagViewToTag(val)
		}
	}
	if vres.LatestVersion != nil {
		res.LatestVersion = newVersion(vres.LatestVersion)
	}
	return res
}

// newResourceWithoutVersion converts projected type Resource to service type
// Resource.
func newResourceWithoutVersion(vres *resourceviews.ResourceView) *Resource {
	res := &Resource{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.Kind != nil {
		res.Kind = *vres.Kind
	}
	if vres.Rating != nil {
		res.Rating = *vres.Rating
	}
	if vres.Catalog != nil {
		res.Catalog = transformResourceviewsCatalogViewToCatalog(vres.Catalog)
	}
	if vres.Tags != nil {
		res.Tags = make([]*Tag, len(vres.Tags))
		for i, val := range vres.Tags {
			res.Tags[i] = transformResourceviewsTagViewToTag(val)
		}
	}
	if vres.LatestVersion != nil {
		res.LatestVersion = newVersionWithoutResource(vres.LatestVersion)
	}
	return res
}

// newResource converts projected type Resource to service type Resource.
func newResource(vres *resourceviews.ResourceView) *Resource {
	res := &Resource{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Name != nil {
		res.Name = *vres.Name
	}
	if vres.Kind != nil {
		res.Kind = *vres.Kind
	}
	if vres.Rating != nil {
		res.Rating = *vres.Rating
	}
	if vres.Catalog != nil {
		res.Catalog = transformResourceviewsCatalogViewToCatalog(vres.Catalog)
	}
	if vres.Tags != nil {
		res.Tags = make([]*Tag, len(vres.Tags))
		for i, val := range vres.Tags {
			res.Tags[i] = transformResourceviewsTagViewToTag(val)
		}
	}
	if vres.Versions != nil {
		res.Versions = make([]*Version, len(vres.Versions))
		for i, val := range vres.Versions {
			res.Versions[i] = transformResourceviewsVersionViewToVersion(val)
		}
	}
	if vres.LatestVersion != nil {
		res.LatestVersion = newVersionWithoutResource(vres.LatestVersion)
	}
	return res
}

// newResourceViewInfo projects result type Resource to projected type
// ResourceView using the "info" view.
func newResourceViewInfo(res *Resource) *resourceviews.ResourceView {
	vres := &resourceviews.ResourceView{
		ID:     &res.ID,
		Name:   &res.Name,
		Kind:   &res.Kind,
		Rating: &res.Rating,
	}
	if res.Catalog != nil {
		vres.Catalog = transformCatalogToResourceviewsCatalogView(res.Catalog)
	}
	if res.Tags != nil {
		vres.Tags = make([]*resourceviews.TagView, len(res.Tags))
		for i, val := range res.Tags {
			vres.Tags[i] = transformTagToResourceviewsTagView(val)
		}
	}
	return vres
}

// newResourceViewWithoutVersion projects result type Resource to projected
// type ResourceView using the "withoutVersion" view.
func newResourceViewWithoutVersion(res *Resource) *resourceviews.ResourceView {
	vres := &resourceviews.ResourceView{
		ID:     &res.ID,
		Name:   &res.Name,
		Kind:   &res.Kind,
		Rating: &res.Rating,
	}
	if res.Catalog != nil {
		vres.Catalog = transformCatalogToResourceviewsCatalogView(res.Catalog)
	}
	if res.Tags != nil {
		vres.Tags = make([]*resourceviews.TagView, len(res.Tags))
		for i, val := range res.Tags {
			vres.Tags[i] = transformTagToResourceviewsTagView(val)
		}
	}
	if res.LatestVersion != nil {
		vres.LatestVersion = newVersionViewWithoutResource(res.LatestVersion)
	}
	return vres
}

// newResourceView projects result type Resource to projected type ResourceView
// using the "default" view.
func newResourceView(res *Resource) *resourceviews.ResourceView {
	vres := &resourceviews.ResourceView{
		ID:     &res.ID,
		Name:   &res.Name,
		Kind:   &res.Kind,
		Rating: &res.Rating,
	}
	if res.Catalog != nil {
		vres.Catalog = transformCatalogToResourceviewsCatalogView(res.Catalog)
	}
	if res.Tags != nil {
		vres.Tags = make([]*resourceviews.TagView, len(res.Tags))
		for i, val := range res.Tags {
			vres.Tags[i] = transformTagToResourceviewsTagView(val)
		}
	}
	if res.Versions != nil {
		vres.Versions = make([]*resourceviews.VersionView, len(res.Versions))
		for i, val := range res.Versions {
			vres.Versions[i] = transformVersionToResourceviewsVersionView(val)
		}
	}
	if res.LatestVersion != nil {
		vres.LatestVersion = newVersionViewWithoutResource(res.LatestVersion)
	}
	return vres
}

// newVersionTiny converts projected type Version to service type Version.
func newVersionTiny(vres *resourceviews.VersionView) *Version {
	res := &Version{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Version != nil {
		res.Version = *vres.Version
	}
	if vres.Resource != nil {
		res.Resource = newResource(vres.Resource)
	}
	return res
}

// newVersionMin converts projected type Version to service type Version.
func newVersionMin(vres *resourceviews.VersionView) *Version {
	res := &Version{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Version != nil {
		res.Version = *vres.Version
	}
	if vres.RawURL != nil {
		res.RawURL = *vres.RawURL
	}
	if vres.WebURL != nil {
		res.WebURL = *vres.WebURL
	}
	if vres.Resource != nil {
		res.Resource = newResource(vres.Resource)
	}
	return res
}

// newVersionWithoutResource converts projected type Version to service type
// Version.
func newVersionWithoutResource(vres *resourceviews.VersionView) *Version {
	res := &Version{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Version != nil {
		res.Version = *vres.Version
	}
	if vres.DisplayName != nil {
		res.DisplayName = *vres.DisplayName
	}
	if vres.Description != nil {
		res.Description = *vres.Description
	}
	if vres.MinPipelinesVersion != nil {
		res.MinPipelinesVersion = *vres.MinPipelinesVersion
	}
	if vres.RawURL != nil {
		res.RawURL = *vres.RawURL
	}
	if vres.WebURL != nil {
		res.WebURL = *vres.WebURL
	}
	if vres.UpdatedAt != nil {
		res.UpdatedAt = *vres.UpdatedAt
	}
	if vres.Resource != nil {
		res.Resource = newResource(vres.Resource)
	}
	return res
}

// newVersion converts projected type Version to service type Version.
func newVersion(vres *resourceviews.VersionView) *Version {
	res := &Version{}
	if vres.ID != nil {
		res.ID = *vres.ID
	}
	if vres.Version != nil {
		res.Version = *vres.Version
	}
	if vres.DisplayName != nil {
		res.DisplayName = *vres.DisplayName
	}
	if vres.Description != nil {
		res.Description = *vres.Description
	}
	if vres.MinPipelinesVersion != nil {
		res.MinPipelinesVersion = *vres.MinPipelinesVersion
	}
	if vres.RawURL != nil {
		res.RawURL = *vres.RawURL
	}
	if vres.WebURL != nil {
		res.WebURL = *vres.WebURL
	}
	if vres.UpdatedAt != nil {
		res.UpdatedAt = *vres.UpdatedAt
	}
	if vres.Resource != nil {
		res.Resource = newResourceInfo(vres.Resource)
	}
	return res
}

// newVersionViewTiny projects result type Version to projected type
// VersionView using the "tiny" view.
func newVersionViewTiny(res *Version) *resourceviews.VersionView {
	vres := &resourceviews.VersionView{
		ID:      &res.ID,
		Version: &res.Version,
	}
	return vres
}

// newVersionViewMin projects result type Version to projected type VersionView
// using the "min" view.
func newVersionViewMin(res *Version) *resourceviews.VersionView {
	vres := &resourceviews.VersionView{
		ID:      &res.ID,
		Version: &res.Version,
		RawURL:  &res.RawURL,
		WebURL:  &res.WebURL,
	}
	return vres
}

// newVersionViewWithoutResource projects result type Version to projected type
// VersionView using the "withoutResource" view.
func newVersionViewWithoutResource(res *Version) *resourceviews.VersionView {
	vres := &resourceviews.VersionView{
		ID:                  &res.ID,
		Version:             &res.Version,
		DisplayName:         &res.DisplayName,
		Description:         &res.Description,
		MinPipelinesVersion: &res.MinPipelinesVersion,
		RawURL:              &res.RawURL,
		WebURL:              &res.WebURL,
		UpdatedAt:           &res.UpdatedAt,
	}
	return vres
}

// newVersionView projects result type Version to projected type VersionView
// using the "default" view.
func newVersionView(res *Version) *resourceviews.VersionView {
	vres := &resourceviews.VersionView{
		ID:                  &res.ID,
		Version:             &res.Version,
		DisplayName:         &res.DisplayName,
		Description:         &res.Description,
		MinPipelinesVersion: &res.MinPipelinesVersion,
		RawURL:              &res.RawURL,
		WebURL:              &res.WebURL,
		UpdatedAt:           &res.UpdatedAt,
	}
	if res.Resource != nil {
		vres.Resource = newResourceViewInfo(res.Resource)
	}
	return vres
}

// newVersions converts projected type Versions to service type Versions.
func newVersions(vres *resourceviews.VersionsView) *Versions {
	res := &Versions{}
	if vres.Versions != nil {
		res.Versions = make([]*Version, len(vres.Versions))
		for i, val := range vres.Versions {
			res.Versions[i] = transformResourceviewsVersionViewToVersion(val)
		}
	}
	if vres.Latest != nil {
		res.Latest = newVersionMin(vres.Latest)
	}
	return res
}

// newVersionsView projects result type Versions to projected type VersionsView
// using the "default" view.
func newVersionsView(res *Versions) *resourceviews.VersionsView {
	vres := &resourceviews.VersionsView{}
	if res.Versions != nil {
		vres.Versions = make([]*resourceviews.VersionView, len(res.Versions))
		for i, val := range res.Versions {
			vres.Versions[i] = transformVersionToResourceviewsVersionView(val)
		}
	}
	if res.Latest != nil {
		vres.Latest = newVersionViewMin(res.Latest)
	}
	return vres
}

// transformResourceviewsCatalogViewToCatalog builds a value of type *Catalog
// from a value of type *resourceviews.CatalogView.
func transformResourceviewsCatalogViewToCatalog(v *resourceviews.CatalogView) *Catalog {
	if v == nil {
		return nil
	}
	res := &Catalog{
		ID:   *v.ID,
		Type: *v.Type,
	}

	return res
}

// transformResourceviewsTagViewToTag builds a value of type *Tag from a value
// of type *resourceviews.TagView.
func transformResourceviewsTagViewToTag(v *resourceviews.TagView) *Tag {
	if v == nil {
		return nil
	}
	res := &Tag{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// transformResourceviewsVersionViewToVersion builds a value of type *Version
// from a value of type *resourceviews.VersionView.
func transformResourceviewsVersionViewToVersion(v *resourceviews.VersionView) *Version {
	if v == nil {
		return nil
	}
	res := &Version{
		ID:                  *v.ID,
		Version:             *v.Version,
		DisplayName:         *v.DisplayName,
		Description:         *v.Description,
		MinPipelinesVersion: *v.MinPipelinesVersion,
		RawURL:              *v.RawURL,
		WebURL:              *v.WebURL,
		UpdatedAt:           *v.UpdatedAt,
	}
	if v.Resource != nil {
		res.Resource = transformResourceviewsResourceViewToResource(v.Resource)
	}

	return res
}

// transformResourceviewsResourceViewToResource builds a value of type
// *Resource from a value of type *resourceviews.ResourceView.
func transformResourceviewsResourceViewToResource(v *resourceviews.ResourceView) *Resource {
	res := &Resource{}
	if v.ID != nil {
		res.ID = *v.ID
	}
	if v.Name != nil {
		res.Name = *v.Name
	}
	if v.Kind != nil {
		res.Kind = *v.Kind
	}
	if v.Rating != nil {
		res.Rating = *v.Rating
	}
	if v.Catalog != nil {
		res.Catalog = transformResourceviewsCatalogViewToCatalog(v.Catalog)
	}
	if v.Tags != nil {
		res.Tags = make([]*Tag, len(v.Tags))
		for i, val := range v.Tags {
			res.Tags[i] = transformResourceviewsTagViewToTag(val)
		}
	}
	if v.Versions != nil {
		res.Versions = make([]*Version, len(v.Versions))
		for i, val := range v.Versions {
			res.Versions[i] = transformResourceviewsVersionViewToVersion(val)
		}
	}

	return res
}

// transformCatalogToResourceviewsCatalogView builds a value of type
// *resourceviews.CatalogView from a value of type *Catalog.
func transformCatalogToResourceviewsCatalogView(v *Catalog) *resourceviews.CatalogView {
	res := &resourceviews.CatalogView{
		ID:   &v.ID,
		Type: &v.Type,
	}

	return res
}

// transformTagToResourceviewsTagView builds a value of type
// *resourceviews.TagView from a value of type *Tag.
func transformTagToResourceviewsTagView(v *Tag) *resourceviews.TagView {
	res := &resourceviews.TagView{
		ID:   &v.ID,
		Name: &v.Name,
	}

	return res
}

// transformVersionToResourceviewsVersionView builds a value of type
// *resourceviews.VersionView from a value of type *Version.
func transformVersionToResourceviewsVersionView(v *Version) *resourceviews.VersionView {
	res := &resourceviews.VersionView{
		ID:                  &v.ID,
		Version:             &v.Version,
		DisplayName:         &v.DisplayName,
		Description:         &v.Description,
		MinPipelinesVersion: &v.MinPipelinesVersion,
		RawURL:              &v.RawURL,
		WebURL:              &v.WebURL,
		UpdatedAt:           &v.UpdatedAt,
	}
	if v.Resource != nil {
		res.Resource = transformResourceToResourceviewsResourceView(v.Resource)
	}

	return res
}

// transformResourceToResourceviewsResourceView builds a value of type
// *resourceviews.ResourceView from a value of type *Resource.
func transformResourceToResourceviewsResourceView(v *Resource) *resourceviews.ResourceView {
	res := &resourceviews.ResourceView{
		ID:     &v.ID,
		Name:   &v.Name,
		Kind:   &v.Kind,
		Rating: &v.Rating,
	}
	if v.Catalog != nil {
		res.Catalog = transformCatalogToResourceviewsCatalogView(v.Catalog)
	}
	if v.Tags != nil {
		res.Tags = make([]*resourceviews.TagView, len(v.Tags))
		for i, val := range v.Tags {
			res.Tags[i] = transformTagToResourceviewsTagView(val)
		}
	}
	if v.Versions != nil {
		res.Versions = make([]*resourceviews.VersionView, len(v.Versions))
		for i, val := range v.Versions {
			res.Versions[i] = transformVersionToResourceviewsVersionView(val)
		}
	}

	return res
}
