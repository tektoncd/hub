// Copyright © 2021 The Tekton Authors.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestQuery_DefaultLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "",
		Kinds: []string{},
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 8, len(res))
}

func TestQuery_ByLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "",
		Limit: 2,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "tekton", res[0].Name)
}

func TestQuery_ByName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "tekton",
		Kinds: []string{},
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, 3, len(res[0].Versions))
}

func TestQuery_ByPartialName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "build",
		Kinds: []string{},
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestQuery_ByKind(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "",
		Kinds: []string{"pipeline"},
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
}

func TestQuery_ByMultipleKinds(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "",
		Kinds: []string{"task", "pipeline"},
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 8, len(res))
}

func TestQuery_ByTags(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "",
		Kinds: []string{},
		Tags:  []string{"atag"},
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
}

func TestQuery_ByPlatforms(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:        tc.DB(),
		Log:       tc.Logger("resource"),
		Name:      "",
		Kinds:     []string{},
		Platforms: []string{"linux/amd64"},
		Limit:     100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "tkn-enterprise", res[0].Name)
	assert.Equal(t, "build-pipeline", res[1].Name)
	assert.Equal(t, "buildah", res[2].Name)
}

func TestQuery_ByCatalogs(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:       tc.DB(),
		Log:      tc.Logger("resource"),
		Name:     "",
		Kinds:    []string{},
		Catalogs: []string{"catalog-community"},
		Limit:    100,
	}

	res, err := req.Query()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(res))
	assert.Equal(t, "catalog-community", res[0].Catalog.Name)
}

func TestQuery_ByWrongCatalogs(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:       tc.DB(),
		Log:      tc.Logger("resource"),
		Name:     "",
		Kinds:    []string{},
		Catalogs: []string{"catalog"},
		Limit:    100,
	}

	res, err := req.Query()
	assert.EqualError(t, err, "resource not found")

	assert.Equal(t, 0, len(res))
}

func TestQuery_ByNameAndKind(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "build",
		Kinds: []string{"pipeline"},
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "build-pipeline", res[0].Name)
}

func TestQuery_ByNameTagsAndMultipleType(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "build",
		Kinds: []string{"task", "pipeline"},
		Tags:  []string{"atag", "ztag"},
		Match: "contains",
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestQuery_ByExactNameAndMultipleType(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "buildah",
		Kinds: []string{"task", "pipeline"},
		Match: "exact",
		Limit: 100,
	}

	res, err := req.Query()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))
}

func TestQuery_ExactNameNotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "build",
		Kinds: []string{},
		Match: "exact",
		Limit: 100,
	}

	_, err := req.Query()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestQuery_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Name:  "foo",
		Kinds: []string{},
		Limit: 100,
	}

	_, err := req.Query()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestList_ByLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:    tc.DB(),
		Log:   tc.Logger("resource"),
		Limit: 3,
	}

	res, err := req.AllResources()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "tekton", res[0].Name)
}

func TestVersionsByID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:  tc.DB(),
		Log: tc.Logger("resource"),
		ID:  1,
	}

	res, err := req.AllVersions()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "0.2", res[2].Version)
}

func TestVersionsByID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:  tc.DB(),
		Log: tc.Logger("resource"),
		ID:  111,
	}

	_, err := req.AllVersions()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindNameVersion(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "catalog-official",
		Kind:    "task",
		Name:    "tkn",
		Version: "0.1",
	}

	res, err := req.ByCatalogKindNameVersion()
	assert.NoError(t, err)
	assert.Equal(t, "tkn", res.Name)
	assert.Equal(t, "0.1", res.Versions[0].Version)
}

func TestByCatalogKindNameVersion_NoResourceWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "catalog-official",
		Kind:    "task",
		Name:    "foo",
		Version: "0.1",
	}

	_, err := req.ByCatalogKindNameVersion()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindNameVersion_NoCatalogWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "Abc",
		Kind:    "task",
		Name:    "foo",
		Version: "0.1",
	}

	_, err := req.ByCatalogKindNameVersion()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindNameVersion_ResourceVersionNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "catalog-official",
		Kind:    "task",
		Name:    "tekton",
		Version: "0.9",
	}

	res, err := req.ByCatalogKindNameVersion()
	assert.NoError(t, err)
	assert.Equal(t, "tekton", res.Name)
	assert.Equal(t, 0, len(res.Versions))
	//assert.Error(t, err)
	//assert.EqualError(t, err, "resource not found")
}

func TestByVersionID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:        tc.DB(),
		Log:       tc.Logger("resource"),
		VersionID: 6,
	}

	res, err := req.ByVersionID()
	assert.NoError(t, err)
	assert.Equal(t, "0.1.1", res.Version)
}

func TestByVersionID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:        tc.DB(),
		Log:       tc.Logger("resource"),
		VersionID: 111,
	}

	_, err := req.ByVersionID()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "catalog-community",
		Kind:    "task",
		Name:    "img",
	}

	res, err := req.ByCatalogKindName()
	assert.NoError(t, err)
	assert.Equal(t, "img", res.Name)
}

func TestByCatalogKindName_NoCatalogWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "abc",
		Kind:    "task",
		Name:    "foo",
	}

	_, err := req.ByCatalogKindName()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindName_ResourceNotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "catalog-community",
		Kind:    "task",
		Name:    "foo",
	}

	_, err := req.ByCatalogKindName()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:  tc.DB(),
		Log: tc.Logger("resource"),
		ID:  1,
	}

	res, err := req.ByID()
	assert.NoError(t, err)
	assert.Equal(t, "tekton", res.Name)
}

func TestByID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:  tc.DB(),
		Log: tc.Logger("resource"),
		ID:  77,
	}

	_, err := req.ByID()
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestGetLatestVersion(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "catalog-official",
		Kind:    "task",
		Name:    "img",
	}

	res, err := req.GetLatestVersion()
	assert.NoError(t, err)
	assert.Equal(t, "0.2", res)
}

func TestGetLatestVersion_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	req := Request{
		Db:      tc.DB(),
		Log:     tc.Logger("resource"),
		Catalog: "foo",
		Kind:    "task",
		Name:    "bar",
	}

	_, err := req.GetLatestVersion()
	assert.Equal(t, err.Error(), "record not found")
}
