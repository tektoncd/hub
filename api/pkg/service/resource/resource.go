// Copyright © 2020 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/parser"
)

type service struct {
	app.Service
}

type request struct {
	db  *gorm.DB
	log *log.Logger
}

var replaceGHtoRaw = strings.NewReplacer("github.com", "raw.githubusercontent.com", "/tree/", "/")

// Errors
var (
	fetchError    = resource.MakeInternalError(fmt.Errorf("failed to fetch resources"))
	notFoundError = resource.MakeNotFound(fmt.Errorf("resource not found"))
)

// New returns the resource service implementation.
func New(api app.BaseConfig) resource.Service {
	return &service{api.Service("resource")}
}

// Find resources based on name, kind or both
func (s *service) Query(ctx context.Context, p *resource.QueryPayload) (resource.ResourceCollection, error) {

	// Validate the kinds passed are supported by Hub
	for _, k := range p.Kinds {
		if !parser.IsSupportedKind(k) {
			return nil, invalidKindError(k)
		}
	}

	db := s.DB(ctx)

	// DISTINCT(resources.id) and resources.id is required as the
	// INNER JOIN of tags and resources returns duplicate records of
	// resources as a resource may have multiple tags, thus we have to
	// find DISTINCT on resource.id

	q := db.Select("DISTINCT(resources.id), resources.*").Scopes(
		filterByTags(p.Tags),
		filterByKinds(p.Kinds),
		filterResourceName(p.Match, p.Name),
		withResourceDetails,
	).Limit(p.Limit)

	req := request{db: q, log: s.Logger(ctx)}
	return req.findAllResources()
}

// List all resources sorted by rating and name
func (s *service) List(ctx context.Context, p *resource.ListPayload) (resource.ResourceCollection, error) {

	db := s.DB(ctx)

	q := db.Scopes(withResourceDetails).
		Limit(p.Limit)

	req := request{db: q, log: s.Logger(ctx)}
	return req.findAllResources()
}

// VersionsByID returns all versions of a resource given its resource id
func (s *service) VersionsByID(ctx context.Context, p *resource.VersionsByIDPayload) (*resource.Versions, error) {

	log := s.Logger(ctx)
	db := s.DB(ctx)

	q := db.Scopes(orderByVersion, filterByResourceID(p.ID))

	var all []model.ResourceVersion
	if err := q.Find(&all).Error; err != nil {
		log.Error(err)
		return nil, fetchError
	}

	if len(all) == 0 {
		return nil, notFoundError
	}

	res := &resource.Versions{}
	for _, r := range all {
		res.Versions = append(res.Versions, minVersionInfo(r))
	}
	res.Latest = minVersionInfo(all[len(all)-1])

	return res, nil
}

// find resource using name of catalog & name, kind and version of resource
func (s *service) ByCatalogKindNameVersion(ctx context.Context, p *resource.ByCatalogKindNameVersionPayload) (*resource.Version, error) {

	log := s.Logger(ctx)
	db := s.DB(ctx)

	q := db.Scopes(
		withVersionInfo(p.Version),
		filterByCatalog(p.Catalog),
		filterByKind(p.Kind),
		filterResourceName("exact", p.Name))

	var r model.Resource
	if err := findOne(q, log, &r); err != nil {
		return nil, err
	}

	switch count := len(r.Versions); {
	case count == 1:
		return versionInfoFromResource(r), nil
	case count == 0:
		return nil, notFoundError
	default:
		log.Warnf("expected to find one version but found %d", count)
		r.Versions = []model.ResourceVersion{r.Versions[0]}
		return versionInfoFromResource(r), nil
	}
}

// find a resource using its version's id
func (s *service) ByVersionID(ctx context.Context, p *resource.ByVersionIDPayload) (*resource.Version, error) {

	db := s.DB(ctx)
	q := db.Scopes(withResourceVersionDetails, filterByVersionID(p.VersionID))

	var v model.ResourceVersion
	if err := findOne(q, s.Logger(ctx), &v); err != nil {
		return nil, err
	}

	return versionInfoFromVersion(v), nil
}

// find resources using name of catalog, resource name and kind of resource
func (s *service) ByCatalogKindName(ctx context.Context, p *resource.ByCatalogKindNamePayload) (*resource.Resource, error) {

	db := s.DB(ctx)

	q := db.Scopes(
		withResourceDetails,
		filterByCatalog(p.Catalog),
		filterByKind(p.Kind),
		filterResourceName("exact", p.Name))

	req := &request{db: q, log: s.Logger(ctx)}
	return req.findSingleResource()
}

// Find a resource using it's id
func (s *service) ByID(ctx context.Context, p *resource.ByIDPayload) (*resource.Resource, error) {

	db := s.DB(ctx)

	q := db.Scopes(withResourceDetails, filterByID(p.ID))

	req := &request{db: q, log: s.Logger(ctx)}
	return req.findSingleResource()
}

func (r *request) findSingleResource() (*resource.Resource, error) {

	var resource model.Resource
	if err := findOne(r.db, r.log, &resource); err != nil {
		return nil, err
	}

	res := initResource(resource)
	for _, v := range resource.Versions {
		res.Versions = append(res.Versions, tinyVersionInfo(v))
	}

	return res, nil
}

func (r *request) findAllResources() (resource.ResourceCollection, error) {

	var rs []model.Resource
	if err := r.db.Find(&rs).Error; err != nil {
		r.log.Error(err)
		return nil, fetchError
	}

	if len(rs) == 0 {
		return nil, notFoundError
	}

	res := resource.ResourceCollection{}
	for _, r := range rs {
		res = append(res, initResource(r))
	}

	return res, nil
}

func initResource(r model.Resource) *resource.Resource {

	res := &resource.Resource{}
	res.ID = r.ID
	res.Name = r.Name
	res.Catalog = &resource.Catalog{
		ID:   r.Catalog.ID,
		Name: r.Catalog.Name,
		Type: r.Catalog.Type,
	}
	res.Kind = r.Kind
	res.Rating = r.Rating

	lv := (r.Versions)[len(r.Versions)-1]
	res.LatestVersion = &resource.Version{
		ID:                  lv.ID,
		Version:             lv.Version,
		Description:         lv.Description,
		DisplayName:         lv.DisplayName,
		MinPipelinesVersion: lv.MinPipelinesVersion,
		WebURL:              lv.URL,
		RawURL:              replaceGHtoRaw.Replace(lv.URL),
		UpdatedAt:           lv.UpdatedAt.UTC().String(),
	}
	for _, tag := range r.Tags {
		res.Tags = append(res.Tags, &resource.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return res
}

func tinyVersionInfo(r model.ResourceVersion) *resource.Version {

	res := &resource.Version{
		ID:      r.ID,
		Version: r.Version,
	}

	return res
}

func minVersionInfo(r model.ResourceVersion) *resource.Version {

	res := tinyVersionInfo(r)
	res.WebURL = r.URL
	res.RawURL = replaceGHtoRaw.Replace(r.URL)

	return res
}

func versionInfoFromResource(r model.Resource) *resource.Version {

	tags := []*resource.Tag{}
	for _, tag := range r.Tags {
		tags = append(tags, &resource.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	res := &resource.Resource{
		ID:     r.ID,
		Name:   r.Name,
		Kind:   r.Kind,
		Rating: r.Rating,
		Tags:   tags,
		Catalog: &resource.Catalog{
			ID:   r.Catalog.ID,
			Name: r.Catalog.Name,
			Type: r.Catalog.Type,
		},
	}

	v := r.Versions[0]
	ver := &resource.Version{
		ID:                  v.ID,
		Version:             v.Version,
		Description:         v.Description,
		DisplayName:         v.DisplayName,
		MinPipelinesVersion: v.MinPipelinesVersion,
		WebURL:              v.URL,
		RawURL:              replaceGHtoRaw.Replace(v.URL),
		UpdatedAt:           v.UpdatedAt.UTC().String(),
		Resource:            res,
	}

	return ver
}

func versionInfoFromVersion(v model.ResourceVersion) *resource.Version {

	// NOTE: we are not preloading all versions (optimisation) and we only
	// need to return version detials of v, thus manually populating only
	// the required info
	v.Resource.Versions = []model.ResourceVersion{v}
	return versionInfoFromResource(v.Resource)
}

func withCatalogAndTags(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Catalog").
		Preload("Tags", orderByTags)
}

// withResourceDetails defines a gorm scope to include all details of resource.
func withResourceDetails(db *gorm.DB) *gorm.DB {
	return db.
		Order("rating DESC, resources.name").
		Scopes(withCatalogAndTags).
		Preload("Versions", orderByVersion)
}

// withVersionInfo defines a gorm scope to include all details of resource
// and filtering a particular version of the resource.
func withVersionInfo(version string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(withResourceDetails).
			Preload("Versions", "version = ?", version)
	}
}

// withResourceVersionDetails defines a gorm scope to include all details of
// resource in resource version.
func withResourceVersionDetails(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Resource").
		Preload("Resource.Catalog").
		Preload("Resource.Tags", orderByTags)
}

func orderByTags(db *gorm.DB) *gorm.DB {
	return db.Order("tags.name ASC")
}

func orderByVersion(db *gorm.DB) *gorm.DB {
	return db.Order("string_to_array(version, '.')::int[];")
}

func findOne(db *gorm.DB, log *log.Logger, result interface{}) error {

	err := db.Find(result).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return notFoundError
		}
		log.Error(err)
		return fetchError
	}
	return nil
}

func filterByID(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func filterByCatalog(catalog string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.Resource{}).
			Joins("JOIN catalogs as c on c.id = resources.catalog_id").
			Where("lower(c.name) = ?", strings.ToLower(catalog))
	}
}

func filterByKind(t string) func(db *gorm.DB) *gorm.DB {
	if t == "" {
		return noop
	}

	t = strings.ToLower(t)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(kind) = ?", t)
	}
}

func filterByKinds(t []string) func(db *gorm.DB) *gorm.DB {
	if len(t) == 0 {
		return noop
	}

	return func(db *gorm.DB) *gorm.DB {
		t = lower(t)
		return db.Where("LOWER(kind) IN (?)", t)
	}
}

func filterByTags(tags []string) func(db *gorm.DB) *gorm.DB {
	if tags == nil {
		return noop
	}

	tags = lower(tags)
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.Resource{}).
			Joins("JOIN resource_tags as rt on rt.resource_id = resources.id").
			Joins("JOIN tags on tags.id = rt.tag_id").
			Where("lower(tags.name) in (?)", tags)
	}
}

func filterByResourceID(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("resource_id = ?", id)
	}
}

func filterByVersionID(versionID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", versionID)
	}
}

func filterResourceName(match, name string) func(db *gorm.DB) *gorm.DB {
	if name == "" {
		return noop
	}
	name = strings.ToLower(name)
	switch match {
	case "exact":
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("LOWER(resources.name) = ?", name)
		}
	default:
		likeName := "%" + name + "%"
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("LOWER(resources.name) LIKE ?", likeName)
		}
	}
}

func noop(db *gorm.DB) *gorm.DB {
	return db
}

// This function lowercase all the elements of an array
func lower(t []string) []string {
	for i := range t {
		t[i] = strings.ToLower(t[i])
	}
	return t
}

func invalidKindError(kind string) error {
	return resource.MakeInvalidKind(fmt.Errorf("resource kind '%s' not supported. Supported kinds are %v",
		kind, parser.SupportedKinds()))
}
