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

package admin

import (
	"context"

	"github.com/tektoncd/hub/api/gen/admin"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/service/auth"
)

type service struct {
	*auth.Service
}

// New returns the admin service implementation.
func New(api app.Config) admin.Service {
	return &service{auth.NewService(api, "admin")}
}

// Create or Update an agent user with required scopes
func (s *service) UpdateAgent(ctx context.Context, p *admin.UpdateAgentPayload) (*admin.UpdateAgentResult, error) {
	res := &admin.UpdateAgentResult{}
	return res, nil
}
