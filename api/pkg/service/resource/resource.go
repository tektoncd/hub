// Copyright © 2022 The Tekton Authors.
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
	"net/url"
	"strings"

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	res "github.com/tektoncd/hub/api/pkg/shared/resource"
)

type service struct {
	app.Service
}

var replacerStrings = []string{"github.com", "raw.githubusercontent.com", "/tree/", "/", "/blob/", "/raw/", "/src/", "/raw/"}

// Returns a replacer object which replaces a list of strings with replacements.
// This function basically helps create the raw URL
func getStringReplacer(resourceUrl, provider string) *strings.Replacer {
	if !strings.HasPrefix(resourceUrl, "https://github.com") && provider == "github" {
		parsedUrl, _ := url.Parse(resourceUrl)
		host := "raw." + parsedUrl.Host
		replacerStrings = append(replacerStrings, parsedUrl.Host, host)
	}
	return strings.NewReplacer(replacerStrings...)
}

// New returns the resource service implementation.
func New(api app.BaseConfig) resource.Service {
	return &service{api.Service("resource")}
}

// List all resources sorted by rating and name
func (s *service) List(ctx context.Context) (*resource.Resources, error) {
	return &resource.Resources{}, nil
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

func tinyVersionInfo(r model.ResourceVersion) *resource.ResourceVersionData {

	res := &resource.ResourceVersionData{
		ID:      r.ID,
		Version: r.Version,
	}

	return res
}

// This functions finds the minimum version information of
// the resource such as rawURL, webURL, platforms and HubURL
func minVersionInfo(r model.ResourceVersion) *resource.ResourceVersionData {

	res := tinyVersionInfo(r)
	res.WebURL = r.URL
	res.RawURL = getStringReplacer(r.URL, r.Resource.Catalog.Provider).Replace(r.URL)
	platforms := []*resource.Platform{}
	for _, platform := range r.Platforms {
		platforms = append(platforms, &resource.Platform{
			ID:   platform.ID,
			Name: platform.Name,
		})
	}
	res.Platforms = platforms

	res.HubURLPath = fmt.Sprintf("%s/%s/%s/%s", r.Resource.Catalog.Name, r.Resource.Kind, r.Resource.Name, r.Version)

	return res
}
