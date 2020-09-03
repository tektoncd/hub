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
	"go.uber.org/zap"

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
)

type service struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

var replaceStrings = strings.NewReplacer("github.com", "raw.githubusercontent.com", "/tree/", "/")

// Errors
var (
	fetchError    = resource.MakeInternalError(fmt.Errorf("Failed to fetch resources"))
	notFoundError = resource.MakeNotFound(fmt.Errorf("Resource not found"))
)

// New returns the resource service implementation.
func New(api app.BaseConfig) resource.Service {
	return &service{api.Logger(), api.DB()}
}

// Find resources based on name, kind or both
func (s *service) Query(ctx context.Context, p *resource.QueryPayload) (res resource.ResourceCollection, err error) {

	q := s.db.Scopes(
		withResourceDetails,
		filterByKind(p.Kind),
		matchesName(p.Name),
	).Limit(p.Limit)

	return s.resourcesForQuery(q)
}

// List all resources sorted by rating and name
func (s *service) List(ctx context.Context, p *resource.ListPayload) (res resource.ResourceCollection, err error) {

	q := s.db.Scopes(withResourceDetails).
		Limit(p.Limit)

	return s.resourcesForQuery(q)
}

// VersionsByID returns all versions of a resource given its resource id
func (s *service) VersionsByID(ctx context.Context, p *resource.VersionsByIDPayload) (res *resource.Versions, err error) {

	q := s.db.Scopes(orderByVersion, filterByResourceID(p.ID))

	var all []model.ResourceVersion
	if err := q.Find(&all).Error; err != nil {
		s.logger.Error(err)
		return nil, fetchError
	}

	if len(all) == 0 {
		return nil, notFoundError
	}

	res = &resource.Versions{}
	for _, r := range all {
		res.Versions = append(res.Versions, minVersionInfo(r))
	}
	res.Latest = minVersionInfo(all[len(all)-1])

	return res, nil
}

func (s *service) ByKindNameVersion(ctx context.Context, p *resource.ByKindNameVersionPayload) (res *resource.Version, err error) {

	q := s.db.Scopes(
		withVersionInfo(p.Version),
		filterByKind(p.Kind),
		filterByName(p.Name))

	var r model.Resource
	if err := findOne(q, &r); err != nil {
		return nil, err
	}

	switch count := len(r.Versions); {
	case count == 1:
		return versionInfoFromResource(r), nil
	case count == 0:
		return nil, notFoundError
	default:
		s.logger.Warnf("expected to find one version but found %d", count)
		r.Versions = []model.ResourceVersion{r.Versions[0]}
		return versionInfoFromResource(r), nil
	}
}

// find a resource using its version's id
func (s *service) ByVersionID(ctx context.Context, p *resource.ByVersionIDPayload) (res *resource.Version, err error) {

	q := s.db.Scopes(withResourceVersionDetails, filterByVersionID(p.VersionID))

	var v model.ResourceVersion
	if err := findOne(q, &v); err != nil {
		return nil, err
	}

	return versionInfoFromVersion(v), nil
}

// find resources using name and kind
func (s *service) ByKindName(ctx context.Context, p *resource.ByKindNamePayload) (res resource.ResourceCollection, err error) {

	q := s.db.Scopes(
		withResourceDetails,
		filterByKind(p.Kind),
		filterByName(p.Name))

	return s.resourcesForQuery(q)
}

// Find a resource using it's id
func (s *service) ByID(ctx context.Context, p *resource.ByIDPayload) (res *resource.Resource, err error) {

	q := s.db.Scopes(withResourceDetails,
		filterByID(p.ID))

	var r model.Resource
	if err := findOne(q, &r); err != nil {
		return nil, err
	}

	res = initResource(r)

	for _, v := range r.Versions {
		res.Versions = append(res.Versions, tinyVersionInfo(v))
	}

	return res, nil
}

func (s *service) resourcesForQuery(q *gorm.DB) (resource.ResourceCollection, error) {

	var rs []model.Resource
	if err := q.Find(&rs).Error; err != nil {
		s.logger.Error(err)
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
		RawURL:              replaceStrings.Replace(lv.URL),
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
	res.RawURL = replaceStrings.Replace(r.URL)

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
		RawURL:              replaceStrings.Replace(v.URL),
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
		Order("rating DESC, name").
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

func findOne(db *gorm.DB, result interface{}) error {

	err := db.Find(result).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return notFoundError
		}
		return fetchError
	}
	return nil
}

func filterByID(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
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

func filterByName(name string) func(db *gorm.DB) *gorm.DB {
	if name == "" {
		return noop
	}

	name = strings.ToLower(name)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(name) = ?", name)
	}
}

func matchesName(name string) func(db *gorm.DB) *gorm.DB {
	if name == "" {
		return noop
	}

	likeName := "%" + strings.ToLower(name) + "%"
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("LOWER(name) LIKE ?", likeName)
	}
}

func noop(db *gorm.DB) *gorm.DB {
	return db
}
