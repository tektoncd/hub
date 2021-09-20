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

package catalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/catalog"
	"github.com/tektoncd/hub/api/gen/http/catalog/server"
	"github.com/tektoncd/hub/api/pkg/service/validator"
	"github.com/tektoncd/hub/api/pkg/testutils"
	"gotest.tools/v3/golden"
)

func RefreshChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := validator.NewService(tc.APIConfig, "catalog")
	checker := goahttpcheck.New()
	checker.Mount(server.NewRefreshHandler,
		server.MountRefreshHandler,
		catalog.NewRefreshEndpoint(NewServiceTest(tc), service.JWTAuth))
	return checker
}

func TestRefresh_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with catalog:refresh scope
	agent, token, err := tc.AgentWithScopes("agent-001", "catalog:refresh")
	assert.Equal(t, agent.AgentName, "agent-001")
	assert.NoError(t, err)

	RefreshChecker(tc).Test(t, http.MethodPost, "/catalog/catalog-official/refresh").
		WithHeader("Authorization", token).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := catalog.Job{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		assert.Equal(t, uint(10001), res.ID)
	})
}

func RefreshAllChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := validator.NewService(tc.APIConfig, "catalog")
	checker := goahttpcheck.New()
	checker.Mount(server.NewRefreshAllHandler,
		server.MountRefreshAllHandler,
		catalog.NewRefreshAllEndpoint(NewServiceTest(tc), service.JWTAuth))
	return checker
}

func CatalogErrorChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := validator.NewService(tc.APIConfig, "catalog")
	checker := goahttpcheck.New()
	checker.Mount(server.NewCatalogErrorHandler,
		server.MountCatalogErrorHandler,
		catalog.NewCatalogErrorEndpoint(NewServiceTest(tc), service.JWTAuth))
	return checker
}

func TestRefreshAll_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with catalog:refresh scope
	agent, token, err := tc.AgentWithScopes("agent-001", "catalog:refresh")
	assert.Equal(t, agent.AgentName, "agent-001")
	assert.NoError(t, err)

	RefreshAllChecker(tc).Test(t, http.MethodPost, "/catalog/refresh").
		WithHeader("Authorization", token).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := []*catalog.Job{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		assert.Equal(t, 3, len(res))
	})
}

func TestCatalogError_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// User with catalog:refresh scope
	agent, token, err := tc.AgentWithScopes("agent-001", "catalog:refresh")
	assert.Equal(t, agent.AgentName, "agent-001")
	assert.NoError(t, err)

	CatalogErrorChecker(tc).Test(t, http.MethodGet, "/catalog/catalog-official/error").
		WithHeader("Authorization", token).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})

}

func TestCatalogError_HttpHavingNoError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// User with catalog:refresh scope
	agent, token, err := tc.AgentWithScopes("agent-001", "catalog:refresh")
	assert.Equal(t, agent.AgentName, "agent-001")
	assert.NoError(t, err)

	CatalogErrorChecker(tc).Test(t, http.MethodGet, "/catalog/catalog-community/error").
		WithHeader("Authorization", token).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})

}
