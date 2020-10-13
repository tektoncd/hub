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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/admin"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestUpdateAgent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, _, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	adminSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &admin.UpdateAgentPayload{Name: "agent-007", Scopes: []string{"config:refresh"}}
	res, err := adminSvc.UpdateAgent(ctx, payload)
	assert.NoError(t, err)

	// expected jwt for agent-007
	agent, agentToken, err := tc.AgentWithScopes("agent-007", "config:refresh")
	assert.Equal(t, agent.AgentName, "agent-007")
	assert.NoError(t, err)

	assert.Equal(t, agentToken, res.Token)
}

func TestUpdateAgent_NormalUserExistsWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, _, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	adminSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &admin.UpdateAgentPayload{Name: "foo", Scopes: []string{"config:refresh", "agent:create"}}
	_, err = adminSvc.UpdateAgent(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "user exists with name: foo")
}

func TestUpdateAgent_InvalidScopeInPayload(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, _, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	adminSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &admin.UpdateAgentPayload{Name: "agent:007", Scopes: []string{"abc:read"}}
	_, err = adminSvc.UpdateAgent(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "scope does not exist: abc:read")
}

func TestUpdateAgent_UpdateScopesCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, _, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	adminSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &admin.UpdateAgentPayload{Name: "agent-001", Scopes: []string{"config:refresh", "agent:create"}}
	res, err := adminSvc.UpdateAgent(ctx, payload)
	assert.NoError(t, err)

	// expected jwt for agent-001 after updating scopes
	agent, agentToken, err := tc.AgentWithScopes("agent-001", "config:refresh", "agent:create")
	assert.Equal(t, agent.AgentName, "agent-001")
	assert.NoError(t, err)

	assert.Equal(t, agentToken, res.Token)
}
