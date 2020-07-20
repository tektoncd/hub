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

	QueryChecker(tc).Test(t, http.MethodGet, "/query?name=build&type=pipeline&limit=1").Check().
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

func TestList_Http_WithLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	checker := goahttpcheck.New()
	checker.Mount(
		server.NewListHandler,
		server.MountListHandler,
		resource.NewListEndpoint(New(tc)))

	checker.Test(t, http.MethodGet, "/resources?limit=2").Check().
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

	checker := goahttpcheck.New()
	checker.Mount(
		server.NewListHandler,
		server.MountListHandler,
		resource.NewListEndpoint(New(tc)))

	checker.Test(t, http.MethodGet, "/resources").Check().
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

func ByTypeNameVersionChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(
		server.NewByTypeNameVersionHandler,
		server.MountByTypeNameVersionHandler,
		resource.NewByTypeNameVersionEndpoint(New(tc)))
	return checker
}

func TestByTypeNameVersion_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByTypeNameVersionChecker(tc).Test(t, http.MethodGet, "/resource/task/tekton/0.1.1").Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}

func TestByTypeNameVersion_Http_ErrorCase(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	ByTypeNameVersionChecker(tc).Test(t, http.MethodGet, "/resource/task/foo/0.1.1").Check().
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
