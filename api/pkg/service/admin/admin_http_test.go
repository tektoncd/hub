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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	goa "goa.design/goa/v3/pkg"

	"github.com/tektoncd/hub/api/gen/admin"
	"github.com/tektoncd/hub/api/gen/http/admin/server"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

// Token for the user with github name "foo-bar" and github login "foo"
// It has a scope "agent:create" along with default scope
const validTokenWithAgentCreateScope = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MTEsImxvZ2luIjoiZm9vIiwibmFtZSI6ImZvby1iYXIiLCJzY29wZXMiOlsicmF0aW5nOnJlYWQiLCJyYXRpbmc6d3JpdGUiLCJhZ2VudDpjcmVhdGUiXX0." +
	"bKPINZyhzX2Ls1QM--UV56cC-vm5uymT8y-DmEhY1dE"

// Token for the agent with name "agent-007" with scopes ["test:read"]
const agentToken007 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MTAwMDEsIm5hbWUiOiJhZ2VudC0wMDciLCJzY29wZXMiOlsidGVzdDpyZWFkIl0sInR5cGUiOiJhZ2VudCJ9." +
	"x2qjMYZT55-V6fH0z0hVVdM8jQAsiMzyxKQK2iEUiQA"

// Token for the agent with name "agent-007" with scopes ["test:read","agent:create"]
const agentToken007Updated = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MzEsIm5hbWUiOiJhZ2VudC0wMDEiLCJzY29wZXMiOlsidGVzdDpyZWFkIiwiYWdlbnQ6Y3JlYXRlIl0sInR5cGUiOiJhZ2VudCJ9." +
	"0u_SkfkJjuq8jyax8jwLdAzHUl7J0g0a6jkN9gLoJl4"

func UpdateAgentChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, tc.JWTSigningKey())
	checker := goahttpcheck.New()
	checker.Mount(server.NewUpdateAgentHandler,
		server.MountUpdateAgentHandler,
		admin.NewUpdateAgentEndpoint(New(tc), service.JWTAuth))
	return checker
}

func TestUpdateAgent_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	data := []byte(`{"name": "agent-007","scopes": ["test:read"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", validTokenWithAgentCreateScope).WithBody(data).
		Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &admin.UpdateAgentResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		assert.Equal(t, agentToken007, res.Token)
	})
}

func TestUpdateAgent_Http_NormalUserExistWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	data := []byte(`{"name": "foo-bar","scopes": ["test:read"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", validTokenWithAgentCreateScope).WithBody(data).
		Check().
		HasStatus(400).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := &goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.EqualError(t, err, "user exists with name: foo-bar")
	})
}

func TestUpdateAgent_Http_InvalidScopeCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	data := []byte(`{"name": "agent-001","scopes": ["invalid:scope"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", validTokenWithAgentCreateScope).WithBody(data).
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

	data := []byte(`{"name": "agent-001","scopes": ["test:read","agent:create"]}`)

	UpdateAgentChecker(tc).Test(t, http.MethodPut, "/system/user/agent").
		WithHeader("Authorization", validTokenWithAgentCreateScope).WithBody(data).
		Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &admin.UpdateAgentResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		assert.Equal(t, agentToken007Updated, res.Token)
	})
}
