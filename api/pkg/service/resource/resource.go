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

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/app"
)

type service struct {
	app.Service
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
func (s *service) VersionsByID(ctx context.Context, p *resource.VersionsByIDPayload) (*resource.VersionsByIDResult, error) {
	return &resource.VersionsByIDResult{Location: fmt.Sprintf("/v1/resource/%d/versions", p.ID)}, nil
}

// Find resource using name of catalog & name, kind and version of resource
func (s *service) ByCatalogKindNameVersion(ctx context.Context, p *resource.ByCatalogKindNameVersionPayload) (res *resource.ByCatalogKindNameVersionResult, err error) {
	return &resource.ByCatalogKindNameVersionResult{Location: fmt.Sprintf("/v1/resource/%s/%s/%s/%s", p.Catalog, p.Kind, p.Name, p.Version)}, nil
}

// Find a resource using its version's id
func (s *service) ByVersionID(ctx context.Context, p *resource.ByVersionIDPayload) (*resource.ByVersionIDResult, error) {
	return &resource.ByVersionIDResult{Location: fmt.Sprintf("/v1/resource/version/%d", p.VersionID)}, nil
}

// Find resources using name of catalog, resource name and kind of resource
func (s *service) ByCatalogKindName(ctx context.Context, p *resource.ByCatalogKindNamePayload) (*resource.ByCatalogKindNameResult, error) {
	if p.Pipelinesversion != nil {
		return &resource.ByCatalogKindNameResult{Location: fmt.Sprintf("/v1/resource/%s/%s/%s?pipelinesversion=%s", p.Catalog, p.Kind, p.Name, *p.Pipelinesversion)}, nil
	}
	return &resource.ByCatalogKindNameResult{Location: fmt.Sprintf("/v1/resource/%s/%s/%s", p.Catalog, p.Kind, p.Name)}, nil
}
