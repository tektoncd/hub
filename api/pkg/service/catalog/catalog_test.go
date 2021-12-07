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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/catalog"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/service/validator"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

// NewServiceTest returns the catalog service implementation for test.
func NewServiceTest(api app.Config) catalog.Service {
	svc := validator.NewService(api, "catalog")
	wq := newSyncer(api)

	s := &service{
		svc,
		wq,
	}
	return s
}

func TestRefresh(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with catalog:refresh scope
	user, _, err := tc.UserWithScopes("foo", "foo@bar.com", "catalog:refresh")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	catalogSvc := NewServiceTest(tc)
	ctx := validator.WithUserID(context.Background(), user.ID)

	payload := &catalog.RefreshPayload{CatalogName: "catalog-official"}
	job, err := catalogSvc.Refresh(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, uint(10001), job.ID)
	assert.Equal(t, "queued", job.Status)
}

func TestRefresh_CatalogNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with catalog:refresh scope
	user, _, err := tc.UserWithScopes("foo", "foo@bar.com", "catalog:refresh")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	catalogSvc := NewServiceTest(tc)
	ctx := validator.WithUserID(context.Background(), user.ID)

	payload := &catalog.RefreshPayload{CatalogName: "abc"}
	_, err = catalogSvc.Refresh(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "abc catalog not found")
}

func TestRefreshAgain(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with catalog:refresh scopes
	user, _, err := tc.UserWithScopes("foo", "foo@bar.com", "catalog:refresh")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	catalogSvc := NewServiceTest(tc)
	ctx := validator.WithUserID(context.Background(), user.ID)

	payload := &catalog.RefreshPayload{CatalogName: "catalog-official"}
	res, err := catalogSvc.Refresh(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, uint(10001), res.ID)
	assert.Equal(t, "queued", res.Status)

	res, err = catalogSvc.Refresh(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, uint(10001), res.ID)
	assert.Equal(t, "queued", res.Status)
}

func TestRefresh_All(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with catalog:refresh scope
	user, _, err := tc.UserWithScopes("foo", "foo@bar.com", "catalog:refresh")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	catalogSvc := NewServiceTest(tc)
	ctx := validator.WithUserID(context.Background(), user.ID)

	payload := &catalog.RefreshAllPayload{}
	jobs, err := catalogSvc.RefreshAll(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(jobs))

	assert.Contains(t, []string{"catalog-official", "catalog-community", "catalog-enterprise"}, jobs[0].CatalogName)
	assert.Equal(t, "queued", jobs[0].Status)
	assert.Contains(t, []string{"catalog-official", "catalog-community", "catalog-enterprise"}, jobs[1].CatalogName)
	assert.Equal(t, "queued", jobs[1].Status)
	assert.Contains(t, []string{"catalog-official", "catalog-community", "catalog-enterprise"}, jobs[2].CatalogName)
	assert.Equal(t, "queued", jobs[2].Status)
}

func TestCatalogError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// User with catalog:refresh scope
	user, _, err := tc.UserWithScopes("foo", "foo@bar.com", "catalog:refresh")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	catalogSvc := NewServiceTest(tc)
	ctx := validator.WithUserID(context.Background(), user.ID)

	payload := &catalog.CatalogErrorPayload{CatalogName: "catalog-official"}
	res, err := catalogSvc.CatalogError(ctx, payload)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(res.Data))
	assert.Equal(t, "info", res.Data[1].Type)
	assert.Equal(t, "display name is missing for buildah task", res.Data[1].Errors[0])
}

func TestCatalogErrorHavingNoError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// User with catalog:refresh scope
	user, _, err := tc.UserWithScopes("foo", "foo@bar.com", "catalog:refresh")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	catalogSvc := NewServiceTest(tc)
	ctx := validator.WithUserID(context.Background(), user.ID)

	payload := &catalog.CatalogErrorPayload{CatalogName: "catalog-community"}
	res, err := catalogSvc.CatalogError(ctx, payload)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(res.Data))
}
