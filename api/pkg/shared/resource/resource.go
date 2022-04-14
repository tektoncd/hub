// Copyright © 2021 The Tekton Authors.
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
	"fmt"
	"strings"

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/parser"
	"gorm.io/gorm"
)

type Request struct {
	Db         *gorm.DB
	Log        *log.Logger
	ID         uint
	Name       string
	Match      string
	Kinds      []string
	Catalogs   []string
	Categories []string
	Tags       []string
	Platforms  []string
	Limit      uint
	Version    string
	Kind       string
	Catalog    string
	VersionID  uint
}

var (
	FetchError    = fmt.Errorf("failed to fetch resources")
	NotFoundError = fmt.Errorf("resource not found")
)

// Query resources based on name, kind, tags.
// Match is the type of search: 'exact' or 'contains'
// Fields: name, []kinds, []Catalogs, []Categories, []Tags, Match, Limit
func (r *Request) Query() ([]model.Resource, error) {

	// Validate the kinds passed are supported by Hub
	for _, k := range r.Kinds {
		if !parser.IsSupportedKind(k) {
			return nil, invalidKindError(k)
		}
	}

	// DISTINCT(resources.id) and resources.id is required as the
	// INNER JOIN of tags and resources returns duplicate records of
	// resources as a resource may have multiple tags, thus we have to
	// find DISTINCT on resource.id

	r.Db = r.Db.Select("DISTINCT(resources.id), resources.*").Scopes(
		filterByTags(r.Tags),
		filterByPlatforms(r.Platforms),
		filterByKinds(r.Kinds),
		filterByCatalogs(r.Catalogs),
		filterByCategories(r.Categories),
		filterResourceName(r.Match, r.Name),
		withResourceDetails,
	).Limit(int(r.Limit))

	return r.findAllResources()
}

// AllResources returns all resources in db sorted by rating and name
// Limit can be passed for limiting number of resources
// Field: Limit
func (r *Request) AllResources() ([]model.Resource, error) {

	r.Db = r.Db.Scopes(withResourceDetails).
		Limit(int(r.Limit))

	return r.findAllResources()
}

// AllVersions returns all versions of a resource
// Fields: ID (Resource ID)
func (r *Request) AllVersions() ([]model.ResourceVersion, error) {

	q := r.Db.Scopes(orderByVersion, filterByResourceID(r.ID)).Preload("Platforms", orderByPlatforms).Preload("Resource").Preload("Resource.Catalog")

	var all []model.ResourceVersion
	if err := q.Find(&all).Error; err != nil {
		r.Log.Error(err)
		return nil, FetchError
	}

	if len(all) == 0 {
		return nil, NotFoundError
	}

	return all, nil
}

// ByCatalogKindNameVersion searches resource by catalog name, kind, resource name, and version
// Fields: Catalog, Kind, Name, Version
func (r *Request) ByCatalogKindNameVersion() (model.Resource, error) {

	q := r.Db.Scopes(
		withVersionInfo(r.Version),
		filterByCatalog(r.Catalog),
		filterByKind(r.Kind),
		filterResourceName("exact", r.Name))

	var res model.Resource
	if err := findOne(q, r.Log, &res); err != nil {
		return model.Resource{}, err
	}

	return res, nil
}

// ByVersionID searches resource version by its ID
// Field: VersionID
func (r *Request) ByVersionID() (model.ResourceVersion, error) {

	q := r.Db.Scopes(withResourceVersionDetails, filterByVersionID(r.VersionID)).Preload("Platforms", orderByPlatforms)

	var v model.ResourceVersion
	if err := findOne(q, r.Log, &v); err != nil {
		return model.ResourceVersion{}, err
	}

	return v, nil
}

// ByCatalogKindName searches resource by catalog name, kind and resource name
// Fields: Catalog, Kind, Name
func (r *Request) ByCatalogKindName() (model.Resource, error) {

	q := r.Db.Scopes(
		withResourceDetails,
		filterByCatalog(r.Catalog),
		filterByKind(r.Kind),
		filterResourceName("exact", r.Name))

	var res model.Resource
	if err := findOne(q, r.Log, &res); err != nil {
		return model.Resource{}, err
	}

	return res, nil
}

// ByID searches resource by its ID
// Field: ID
func (r *Request) ByID() (model.Resource, error) {

	q := r.Db.Scopes(withResourceDetails, filterByID(r.ID))

	var res model.Resource
	if err := findOne(q, r.Log, &res); err != nil {
		return model.Resource{}, err
	}

	return res, nil
}

func (r *Request) findAllResources() ([]model.Resource, error) {

	var rs []model.Resource
	if err := r.Db.Find(&rs).Error; err != nil {
		r.Log.Error(err)
		return nil, FetchError
	}

	if len(rs) == 0 {
		return nil, NotFoundError
	}

	return rs, nil
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
		Preload("Categories", orderByCategories).
		Scopes(withCatalogAndTags).
		Preload("Versions", orderByVersion).
		Preload("Versions.Platforms", orderByPlatforms).
		Preload("Platforms", orderByPlatforms)
}

// withResourceDetailsWithoutVersions defines a gorm scope to include all details of resource.
func withResourceDetailsWithoutVersions(db *gorm.DB) *gorm.DB {
	return db.
		Order("rating DESC, resources.name").
		Preload("Categories", orderByCategories).
		Scopes(withCatalogAndTags).
		Preload("Versions.Platforms", orderByPlatforms).
		Preload("Platforms", orderByPlatforms)
}

// withVersionInfo defines a gorm scope to include all details of resource
// and filtering a particular version of the resource.
func withVersionInfo(version string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(withResourceDetailsWithoutVersions).
			Preload("Versions", "version = ?", version).
			Preload("Versions.Platforms", orderByPlatforms)
	}
}

// withResourceVersionDetails defines a gorm scope to include all details of
// resource in resource version.
func withResourceVersionDetails(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Resource").
		Preload("Resource.Catalog").
		Preload("Resource.Categories", orderByCategories).
		Preload("Resource.Tags", orderByTags).
		Preload("Resource.Platforms", orderByPlatforms)
}

func orderByCategories(db *gorm.DB) *gorm.DB {
	return db.Order("categories.name ASC")
}

func orderByTags(db *gorm.DB) *gorm.DB {
	return db.Order("tags.name ASC")
}

func orderByPlatforms(db *gorm.DB) *gorm.DB {
	return db.Order("platforms.name ASC")
}

func orderByVersion(db *gorm.DB) *gorm.DB {
	return db.Order("string_to_array(version, '.')::int[]")
}

func findOne(db *gorm.DB, log *log.Logger, result interface{}) error {

	err := db.First(result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NotFoundError
		}
		log.Error(err)
		return FetchError
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

func filterByPlatforms(platforms []string) func(db *gorm.DB) *gorm.DB {
	if platforms == nil {
		return noop
	}

	platforms = lower(platforms)
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.Resource{}).
			Joins("JOIN resource_platforms as rp on rp.resource_id = resources.id").
			Joins("JOIN platforms on platforms.id = rp.platform_id").
			Where("lower(platforms.name) in (?)", platforms)
	}
}

func filterByCategories(categories []string) func(db *gorm.DB) *gorm.DB {
	if categories == nil {
		return noop
	}

	categories = lower(categories)
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.Resource{}).
			Joins("JOIN resource_categories as rc on rc.resource_id = resources.id").
			Joins("JOIN categories on categories.id = rc.category_id").
			Where("lower(categories.name) in (?)", categories)
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

func filterByCatalogs(catalogs []string) func(db *gorm.DB) *gorm.DB {
	if len(catalogs) == 0 {
		return noop
	}

	catalogs = lower(catalogs)
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.Resource{}).
			Joins("JOIN catalogs as ct on ct.id = resources.catalog_id").
			Where("lower(ct.name) in (?)", catalogs)
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
	return fmt.Errorf("resource kind '%s' not supported. Supported kinds are %v",
		kind, parser.SupportedKinds())
}
