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

package rating

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/rating"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestGet(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:read scope
	user, _, err := tc.UserWithScopes("foo", "rating:read")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	ratingSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &rating.GetPayload{ID: 1}
	rat, err := ratingSvc.Get(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, 5, rat.Rating)
}

func TestGet_RatingNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:read scope
	user, _, err := tc.UserWithScopes("foo", "rating:read")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	ratingSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &rating.GetPayload{ID: 3}
	rat, err := ratingSvc.Get(ctx, payload)
	assert.NoError(t, err)
	assert.Equal(t, -1, rat.Rating)
}

func TestGet_ResourceNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:read scope
	user, _, err := tc.UserWithScopes("foo", "rating:read")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	ratingSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &rating.GetPayload{ID: 99}
	_, err = ratingSvc.Get(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}

func TestUpdate(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:write scope
	user, _, err := tc.UserWithScopes("foo", "rating:write")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	ratingSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &rating.UpdatePayload{ID: 1, Rating: 3}
	err = ratingSvc.Update(ctx, payload)
	assert.NoError(t, err)
}

func TestUpdate_ResourceNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:write scope
	user, _, err := tc.UserWithScopes("foo", "rating:write")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	ratingSvc := New(tc)
	ctx := auth.WithUserID(context.Background(), user.ID)
	payload := &rating.UpdatePayload{ID: 99, Rating: 3}
	err = ratingSvc.Update(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "resource not found")
}
