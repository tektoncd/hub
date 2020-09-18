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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/golden"

	"github.com/tektoncd/hub/api/gen/http/resource/server"
	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func QueryChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(server.NewQueryHandler,
		server.MountQueryHandler,
		resource.NewQueryEndpoint(New(tc)))
	return checker
}

func TestQuery_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=build&kinds=pipeline&limit=1").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQueryWithKinds_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?kinds=pipeline").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQueryWithTags_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?tags=ztag&tags=Atag").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQueryWithExactName_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=buildah&exact=true").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQueryWithNameAndKinds_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=build&kinds=task&kinds=pipeline").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQueryWithNameAndTags_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=build&tags=atag&tags=ztag").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQueryWithKindsAndTags_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=build&kinds=task&kinds=pipeline").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQueryWithAllParams_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=build&kinds=task&kinds=Pipeline&tags=ztag&tags=Atag&exact=false").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestQuery_Http_ErrorCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=foo").Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", jsonMap["name"])
	})
}

func ListChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(server.NewListHandler,
		server.MountListHandler,
		resource.NewListEndpoint(New(tc)))
	return checker
}

func TestList_Http_WithLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ListChecker(tc).Test(t, http.MethodGet, "/resources?limit=2").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestList_Http_NoLimit(t *testing.T) {
	// Test no limit returns some records
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ListChecker(tc).Test(t, http.MethodGet, "/resources").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func VersionsByIDChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(
		server.NewVersionsByIDHandler,
		server.MountVersionsByIDHandler,
		resource.NewVersionsByIDEndpoint(New(tc)))
	return checker
}

func TestVersionsByID_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	VersionsByIDChecker(tc).Test(t, http.MethodGet, "/resource/1/versions").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestVersionsByID_Http_ErrorCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	VersionsByIDChecker(tc).Test(t, http.MethodGet, "/resource/111/versions").Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", jsonMap["name"])
	})
}

func ByKindNameVersionChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(
		server.NewByKindNameVersionHandler,
		server.MountByKindNameVersionHandler,
		resource.NewByKindNameVersionEndpoint(New(tc)))
	return checker
}

func TestByKindNameVersion_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByKindNameVersionChecker(tc).Test(t, http.MethodGet, "/resource/task/tkn/0.1").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestByKindNameVersion_Http_ErrorCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByKindNameVersionChecker(tc).Test(t, http.MethodGet, "/resource/task/foo/0.1.1").Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", jsonMap["name"])
	})
}

func ByVersionIDChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(
		server.NewByVersionIDHandler,
		server.MountByVersionIDHandler,
		resource.NewByVersionIDEndpoint(New(tc)))
	return checker
}

func TestByVersionID_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByVersionIDChecker(tc).Test(t, http.MethodGet, "/resource/version/4").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestByVersionID_Http_ErrorCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByVersionIDChecker(tc).Test(t, http.MethodGet, "/resource/version/43").Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", jsonMap["name"])
	})
}

func ByKindNameChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(
		server.NewByKindNameHandler,
		server.MountByKindNameHandler,
		resource.NewByKindNameEndpoint(New(tc)))
	return checker
}

func TestByKindName_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByKindNameChecker(tc).Test(t, http.MethodGet, "/resource/task/img").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestByKindName_Http_ErrorCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByKindNameChecker(tc).Test(t, http.MethodGet, "/resource/task/foo").Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", jsonMap["name"])
	})
}

func ByIDChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(
		server.NewByIDHandler,
		server.MountByIDHandler,
		resource.NewByIDEndpoint(New(tc)))
	return checker
}

func TestByID_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByIDChecker(tc).Test(t, http.MethodGet, "/resource/1").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestByID_Http_ErrorCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByIDChecker(tc).Test(t, http.MethodGet, "/resource/77").Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", jsonMap["name"])
	})
}
