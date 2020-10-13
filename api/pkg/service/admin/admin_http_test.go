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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/admin"
	"github.com/tektoncd/hub/api/gen/http/admin/server"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
	goa "goa.design/goa/v3/pkg"
)

func UpdateAgentChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, tc.JWTSigningKey())
	checker := goahttpcheck.New()
	checker.Mount(server.NewUpdateAgentHandler,
		server.MountUpdateAgentHandler,
		admin.NewUpdateAgentEndpoint(New(tc), service.JWTAuth))
	return checker
}

func TestUpdateAgent_Http_NewAgent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, token, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	data := []byte(`{"name": "agent-007","scopes": ["catalog:refresh"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", token).WithBody(data).
		Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &admin.UpdateAgentResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		// expected jwt for agent-007
		agent, agentToken, err := tc.AgentWithScopes("agent-007", "catalog:refresh")
		assert.Equal(t, agent.AgentName, "agent-007")
		assert.NoError(t, err)

		assert.Equal(t, agentToken, res.Token)
	})
}

func TestUpdateAgent_Http_NormalUserExistWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, token, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	data := []byte(`{"name": "foo","scopes": ["catalog:refresh"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", token).WithBody(data).
		Check().
		HasStatus(400).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := &goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.EqualError(t, err, "user exists with name: foo")
	})
}

func TestUpdateAgent_Http_InvalidScopeCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, token, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	data := []byte(`{"name": "agent-001","scopes": ["invalid:scope"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", token).WithBody(data).
		Check().
		HasStatus(400).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := &goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.EqualError(t, err, "scope does not exist: invalid:scope")
	})
}

func TestUpdateAgent_Http_UpdateCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with agent:create scope
	user, token, err := tc.UserWithScopes("foo", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	data := []byte(`{"name": "agent-001","scopes": ["catalog:refresh","agent:create"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", token).WithBody(data).
		Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &admin.UpdateAgentResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		// expected jwt for agent-001 after updating scopes
		agent, agentToken, err := tc.AgentWithScopes("agent-001", "catalog:refresh", "agent:create")
		assert.Equal(t, agent.AgentName, "agent-001")
		assert.NoError(t, err)

		assert.Equal(t, agentToken, res.Token)
	})
}

func RefreshConfigChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(server.NewRefreshConfigHandler,
		server.MountRefreshConfigHandler,
		admin.NewRefreshConfigEndpoint(New(tc), New(tc).(admin.Auther).JWTAuth))
	return checker
}

func TestRefreshConfig_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with config:refresh scope
	user, token, err := tc.UserWithScopes("foo", "config:refresh")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	// DB is populated using text fixture, so it has by default value `testChecksum` in table
	config := &model.Config{}
	err = tc.DB().First(config).Error
	assert.NoError(t, err)
	assert.Equal(t, "testChecksum", config.Checksum)

	data := []byte(`{"force": false}`)

	RefreshConfigChecker(tc).Test(t, http.MethodPost, "/system/config/refresh").
		WithHeader("Authorization", token).
		WithBody(data).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &admin.RefreshConfigResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		// compute checksum of test config file which is reloaded
		checksum, err := computeChecksum()
		assert.NoError(t, err)

		assert.Equal(t, checksum, res.Checksum)
	})
}

func TestRefreshConfig_Http_ForceRefresh(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with config:refresh scope
	user, token, err := tc.UserWithScopes("foo", "config:refresh")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	// DB is populated using text fixture, so it has by default value `testChecksum` in table
	config := &model.Config{}
	err = tc.DB().First(config).Error
	assert.NoError(t, err)
	assert.Equal(t, "testChecksum", config.Checksum)

	data := []byte(`{"force": true}`)

	RefreshConfigChecker(tc).Test(t, http.MethodPost, "/system/config/refresh").
		WithHeader("Authorization", token).
		WithBody(data).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &admin.RefreshConfigResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		// compute checksum of test config file which is reloaded
		checksum, err := computeChecksum()
		assert.NoError(t, err)

		assert.Equal(t, checksum, res.Checksum)
	})
}

func computeChecksum() (string, error) {
	data, err := ioutil.ReadFile("../../../test/config/config.yaml")
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
