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
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/catalog"
	"github.com/tektoncd/hub/api/gen/http/catalog/server"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

// Token for agent agent-001 with catalog:refresh scope
const agentTokenWithCatalogRefreshScope = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MzEsIm5hbWUiOiJhZ2VudC0wMDEiLCJzY29wZXMiOlsiY2F0YWxvZzpyZWZyZXNoIl0sInR5cGUiOiJhZ2VudCJ9." +
	"-PSkMT2YJuM_WRtNmYJm5fpnUh47gLIL4Cdykz2lTFE"

func RefreshChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, tc.JWTSigningKey())
	checker := goahttpcheck.New()
	checker.Mount(server.NewRefreshHandler,
		server.MountRefreshHandler,
		catalog.NewRefreshEndpoint(NewServiceTest(tc), service.JWTAuth))
	return checker
}

func TestRefresh_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	data := []byte(`{"name": "catalog-official","org":"tektoncd"}`)

	RefreshChecker(tc).Test(t, http.MethodPost, "/catalog/refresh").
		WithHeader("Authorization", agentTokenWithCatalogRefreshScope).
		WithBody(data).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var res *catalog.Job
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		assert.Equal(t, 10001, int(res.ID))
	})
}
