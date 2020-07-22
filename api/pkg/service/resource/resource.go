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
func New(api app.Config) resource.Service {
	return &service{api.Logger(), api.DB()}
}

// Find resources based on name, type or both
func (s *service) Query(ctx context.Context, p *resource.QueryPayload) (res resource.ResourceCollection, err error) {

	q := s.db.Scopes(withResourceDetails).Limit(p.Limit)

	if p.Type != "" {
		q = q.Where("LOWER(type) = ?", p.Type)
	}

	if p.Name != "" {
		name := "%" + strings.ToLower(p.Name) + "%"
		q = q.Where("LOWER(name) LIKE ?", name)
	}

	return s.resourcesForQuery(q)
}

// List all resources sorted by rating and name
func (s *service) List(ctx context.Context, p *resource.ListPayload) (res resource.ResourceCollection, err error) {
	q := s.db.Scopes(withResourceDetails).Limit(p.Limit)
	return s.resourcesForQuery(q)
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

// VersionsByID returns all versions of a resource given its resource id
func (s *service) VersionsByID(ctx context.Context, p *resource.VersionsByIDPayload) (res *resource.Versions, err error) {

	q := s.db.Scopes(orderByVersion).Where("resource_id = ?", p.ID)

	var all []model.ResourceVersion
	if err := q.Find(&all).Error; err != nil {
		s.logger.Error(err)
		return nil, fetchError
	}

	if len(all) == 0 {
		return nil, notFoundError
	}

	var allVersions []*resource.Version
	for _, r := range all {
		allVersions = append(allVersions, initResourceVersion(r))
	}
	latestVersion := initResourceVersion(all[len(all)-1])

	res = &resource.Versions{
		Latest:   latestVersion,
		Versions: allVersions,
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
	res.Type = r.Type
	res.Rating = r.Rating

	lv := (r.Versions)[len(r.Versions)-1]
	res.LatestVersion = &resource.LatestVersion{
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

func initResourceVersion(r model.ResourceVersion) *resource.Version {
	res := &resource.Version{
		ID:      r.ID,
		Version: r.Version,
		WebURL:  r.URL,
		RawURL:  replaceStrings.Replace(r.URL),
	}

	return res
}

// withResourceDetails defines a gorm scope to include all details of resource.
func withResourceDetails(db *gorm.DB) *gorm.DB {
	return db.Order("rating DESC, name").
		Preload("Catalog").
		Preload("Versions", orderByVersion).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name ASC")
		})
}

func orderByVersion(db *gorm.DB) *gorm.DB {
	return db.Order("string_to_array(version, '.')::int[];")
}
