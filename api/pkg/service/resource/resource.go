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
	"time"

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	res "github.com/tektoncd/hub/api/pkg/shared/resource"
)

type service struct {
	app.Service
}

var replaceGHtoRaw = strings.NewReplacer("github.com", "raw.githubusercontent.com", "/tree/", "/")

// New returns the resource service implementation.
func New(api app.BaseConfig) resource.Service {
	return &service{api.Service("resource")}
}

// Find resources based on name, kind or both
func (s *service) Query(ctx context.Context, p *resource.QueryPayload) (*resource.Resources, error) {

	req := res.Request{
		Db:    s.DB(ctx),
		Log:   s.Logger(ctx),
		Name:  p.Name,
		Kinds: p.Kinds,
		Tags:  p.Tags,
		Limit: p.Limit,
		Match: p.Match,
	}

	rArr, err := req.Query()
	if err != nil {
		if strings.Contains(err.Error(), "not supported") {
			return nil, resource.MakeInvalidKind(err)
		}
		if err == res.NotFoundError {
			return nil, resource.MakeNotFound(err)
		}
		if err == res.FetchError {
			return nil, resource.MakeInternalError(err)
		}
	}

	var rd []*resource.ResourceData
	for _, r := range rArr {
		rd = append(rd, initResource(r))
	}

	return &resource.Resources{Data: rd}, nil
}

// List all resources sorted by rating and name
func (s *service) List(ctx context.Context, p *resource.ListPayload) (*resource.Resources, error) {

	req := res.Request{
		Db:    s.DB(ctx),
		Log:   s.Logger(ctx),
		Limit: p.Limit,
	}

	rArr, err := req.AllResources()
	if err != nil {
		if err == res.NotFoundError {
			return nil, resource.MakeNotFound(err)
		}
		if err == res.FetchError {
			return nil, resource.MakeInternalError(err)
		}
	}

	var rd []*resource.ResourceData
	for _, r := range rArr {
		rd = append(rd, initResource(r))
	}

	return &resource.Resources{Data: rd}, nil
}

// VersionsByID returns all versions of a resource given its resource id
func (s *service) VersionsByID(ctx context.Context, p *resource.VersionsByIDPayload) (*resource.ResourceVersions, error) {

	req := res.Request{
		Db:  s.DB(ctx),
		Log: s.Logger(ctx),
		ID:  p.ID,
	}

	versions, err := req.AllVersions()
	if err != nil {
		if err == res.FetchError {
			return nil, resource.MakeInternalError(err)
		}
		if err == res.NotFoundError {
			return nil, resource.MakeNotFound(err)
		}
	}

	var rv resource.Versions
	rv.Versions = []*resource.ResourceVersionData{}
	for _, r := range versions {
		rv.Versions = append(rv.Versions, minVersionInfo(r))
	}
	rv.Latest = minVersionInfo(versions[len(versions)-1])

	return &resource.ResourceVersions{Data: &rv}, nil
}

// Find resource using name of catalog & name, kind and version of resource
func (s *service) ByCatalogKindNameVersion(ctx context.Context, p *resource.ByCatalogKindNameVersionPayload) (*resource.ResourceVersion, error) {

	req := res.Request{
		Db:      s.DB(ctx),
		Log:     s.Logger(ctx),
		Catalog: p.Catalog,
		Kind:    p.Kind,
		Name:    p.Name,
		Version: p.Version,
	}

	r, err := req.ByCatalogKindNameVersion()
	if err != nil {
		if err == res.FetchError {
			return nil, resource.MakeInternalError(err)
		}
		if err == res.NotFoundError {
			return nil, resource.MakeNotFound(err)
		}
	}

	switch count := len(r.Versions); {
	case count == 1:
		return versionInfoFromResource(r), nil
	case count == 0:
		return nil, resource.MakeNotFound(fmt.Errorf("resource not found"))
	default:
		s.Logger(ctx).Warnf("expected to find one version but found %d", count)
		r.Versions = []model.ResourceVersion{r.Versions[0]}
		return versionInfoFromResource(r), nil
	}
}

// Find a resource using its version's id
func (s *service) ByVersionID(ctx context.Context, p *resource.ByVersionIDPayload) (*resource.ResourceVersion, error) {

	req := res.Request{
		Db:        s.DB(ctx),
		Log:       s.Logger(ctx),
		VersionID: p.VersionID,
	}

	v, err := req.ByVersionID()
	if err != nil {
		if err == res.FetchError {
			return nil, resource.MakeInternalError(err)
		}
		if err == res.NotFoundError {
			return nil, resource.MakeNotFound(err)
		}
	}

	return versionInfoFromVersion(v), nil
}

// Find resources using name of catalog, resource name and kind of resource
func (s *service) ByCatalogKindName(ctx context.Context, p *resource.ByCatalogKindNamePayload) (*resource.Resource, error) {

	req := res.Request{
		Db:      s.DB(ctx),
		Log:     s.Logger(ctx),
		Catalog: p.Catalog,
		Kind:    p.Kind,
		Name:    p.Name,
	}

	r, err := req.ByCatalogKindName()
	if err != nil {
		if err == res.FetchError {
			return nil, resource.MakeInternalError(err)
		}
		if err == res.NotFoundError {
			return nil, resource.MakeNotFound(err)
		}
	}

	res := initResource(r)
	for _, v := range r.Versions {
		res.Versions = append(res.Versions, tinyVersionInfo(v))
	}

	return &resource.Resource{Data: res}, nil
}

// Find a resource using it's id
func (s *service) ByID(ctx context.Context, p *resource.ByIDPayload) (*resource.Resource, error) {

	req := res.Request{
		Db:  s.DB(ctx),
		Log: s.Logger(ctx),
		ID:  p.ID,
	}

	r, err := req.ByID()
	if err != nil {
		if err == res.FetchError {
			return nil, resource.MakeInternalError(err)
		}
		if err == res.NotFoundError {
			return nil, resource.MakeNotFound(err)
		}
	}

	res := initResource(r)
	for _, v := range r.Versions {
		res.Versions = append(res.Versions, tinyVersionInfo(v))
	}

	return &resource.Resource{Data: res}, nil
}

func initResource(r model.Resource) *resource.ResourceData {

	res := &resource.ResourceData{}
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
	res.LatestVersion = &resource.ResourceVersionData{
		ID:                  lv.ID,
		Version:             lv.Version,
		Description:         lv.Description,
		DisplayName:         lv.DisplayName,
		MinPipelinesVersion: lv.MinPipelinesVersion,
		WebURL:              lv.URL,
		RawURL:              replaceGHtoRaw.Replace(lv.URL),
		UpdatedAt:           lv.ModifiedAt.UTC().Format(time.RFC3339),
	}
	res.Tags = []*resource.Tag{}
	for _, tag := range r.Tags {
		res.Tags = append(res.Tags, &resource.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return res
}

func tinyVersionInfo(r model.ResourceVersion) *resource.ResourceVersionData {

	res := &resource.ResourceVersionData{
		ID:      r.ID,
		Version: r.Version,
	}

	return res
}

func minVersionInfo(r model.ResourceVersion) *resource.ResourceVersionData {

	res := tinyVersionInfo(r)
	res.WebURL = r.URL
	res.RawURL = replaceGHtoRaw.Replace(r.URL)

	return res
}

func versionInfoFromResource(r model.Resource) *resource.ResourceVersion {

	var tags []*resource.Tag
	for _, tag := range r.Tags {
		tags = append(tags, &resource.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	res := &resource.ResourceData{
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
	ver := &resource.ResourceVersionData{
		ID:                  v.ID,
		Version:             v.Version,
		Description:         v.Description,
		DisplayName:         v.DisplayName,
		MinPipelinesVersion: v.MinPipelinesVersion,
		WebURL:              v.URL,
		RawURL:              replaceGHtoRaw.Replace(v.URL),
		UpdatedAt:           v.ModifiedAt.UTC().Format(time.RFC3339),
		Resource:            res,
	}

	return &resource.ResourceVersion{Data: ver}
}

func versionInfoFromVersion(v model.ResourceVersion) *resource.ResourceVersion {

	// NOTE: we are not preloading all versions (optimisation) and we only
	// need to return version details of v, thus manually populating only
	// the required info
	v.Resource.Versions = []model.ResourceVersion{v}
	return versionInfoFromResource(v.Resource)
}
