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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestQuery_DefaultLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Kinds: []string{}, Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(all.Data))
}

func TestQuery_ByLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Limit: 2}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(all.Data))
	assert.Equal(t, "tekton", all.Data[0].Name)
}

func TestQuery_ByName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "tekton", Kinds: []string{}, Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all.Data))
	assert.Equal(t, "0.2", all.Data[0].LatestVersion.Version)
}

func TestQuery_ByPartialName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "build", Kinds: []string{}, Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(all.Data))
}

func TestQuery_ByKind(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Kinds: []string{"pipeline"}, Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all.Data))
}

func TestQuery_ByMultipleKinds(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Kinds: []string{"task", "pipeline"}, Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 7, len(all.Data))
}

func TestQuery_ByTags(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Kinds: []string{}, Tags: []string{"atag"}, Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all.Data))
}

func TestQuery_ByNameAndKind(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "build", Kinds: []string{"pipeline"}, Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all.Data))
	assert.Equal(t, "build-pipeline", all.Data[0].Name)
}

func TestQuery_ByNameTagsAndMultipleType(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "build", Kinds: []string{"task", "pipeline"}, Tags: []string{"atag", "ztag"}, Match: "contains", Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(all.Data))
}

func TestQuery_ByExactNameAndMultipleType(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "buildah", Kinds: []string{"task", "pipeline"}, Match: "exact", Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all.Data))
}

func TestQuery_ExactNameNotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "build", Kinds: []string{}, Match: "exact", Limit: 100}
	_, err := resourceSvc.Query(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestQuery_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "foo", Kinds: []string{}, Limit: 100}
	_, err := resourceSvc.Query(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestList_ByLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ListPayload{Limit: 3}
	all, err := resourceSvc.List(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(all.Data))
	assert.Equal(t, "tekton", all.Data[0].Name)
}

func TestVersionsByID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.VersionsByIDPayload{ID: 1}
	res, err := resourceSvc.VersionsByID(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(res.Data.Versions))
	assert.Equal(t, "0.2", res.Data.Latest.Version)
}

func TestVersionsByID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.VersionsByIDPayload{ID: 11}
	_, err := resourceSvc.VersionsByID(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindNameVersion(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByCatalogKindNameVersionPayload{Catalog: "catalog-official", Kind: "task", Name: "tkn", Version: "0.1"}
	res, err := resourceSvc.ByCatalogKindNameVersion(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, "0.1", res.Data.Version)
}

func TestByCatalogKindNameVersion_NoResourceWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByCatalogKindNameVersionPayload{Catalog: "catalog-official", Kind: "task", Name: "foo", Version: "0.1"}
	_, err := resourceSvc.ByCatalogKindNameVersion(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindNameVersion_NoCatalogWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByCatalogKindNameVersionPayload{Catalog: "Abc", Kind: "task", Name: "foo", Version: "0.1"}
	_, err := resourceSvc.ByCatalogKindNameVersion(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindNameVersion_ResourceVersionNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByCatalogKindNameVersionPayload{Catalog: "catalog-official", Kind: "task", Name: "tekton", Version: "0.9"}
	_, err := resourceSvc.ByCatalogKindNameVersion(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByVersionID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByVersionIDPayload{VersionID: 6}
	res, err := resourceSvc.ByVersionID(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, "0.1.1", res.Data.Version)
}

func TestByVersionID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByVersionIDPayload{VersionID: 111}
	_, err := resourceSvc.ByVersionID(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByCatalogKindNamePayload{Catalog: "catalog-community", Kind: "task", Name: "img"}
	res, err := resourceSvc.ByCatalogKindName(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, "img", res.Data.Name)
}

func TestByCatalogKindName_NoCatalogWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByCatalogKindNamePayload{Catalog: "abc", Kind: "task", Name: "foo"}
	_, err := resourceSvc.ByCatalogKindName(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByCatalogKindName_ResourceNotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByCatalogKindNamePayload{Catalog: "catalog-community", Kind: "task", Name: "foo"}
	_, err := resourceSvc.ByCatalogKindName(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestByID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByIDPayload{ID: 1}
	res, err := resourceSvc.ByID(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, "tekton", res.Data.Name)
}

func TestByID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByIDPayload{ID: 77}
	_, err := resourceSvc.ByID(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}
