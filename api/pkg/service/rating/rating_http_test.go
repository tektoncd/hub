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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/http/rating/server"
	"github.com/tektoncd/hub/api/gen/rating"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
	goa "goa.design/goa/v3/pkg"
)

func GetChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, tc.JWTSigningKey())
	checker := goahttpcheck.New()
	checker.Mount(server.NewGetHandler,
		server.MountGetHandler,
		rating.NewGetEndpoint(New(tc), service.JWTAuth))
	return checker
}

func TestGet_Http_InvalidToken(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	GetChecker(tc).Test(t, http.MethodGet, "/resource/1/rating").
		WithHeader("Authorization", "invalidToken").Check().
		HasStatus(401).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "invalid-token", err.Name)
	})
}

func TestGet_Http_InvalidScopes(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// invalid user token, does not have required scopes
	user, token, err := tc.UserWithScopes("abc", "catalog:refresh")
	assert.Equal(t, user.GithubLogin, "abc")
	assert.NoError(t, err)

	GetChecker(tc).Test(t, http.MethodGet, "/resource/1/rating").
		WithHeader("Authorization", token).Check().
		HasStatus(403).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "invalid-scopes", err.Name)
	})
}

func TestGet_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:read scope
	user, token, err := tc.UserWithScopes("foo", "rating:read")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	GetChecker(tc).Test(t, http.MethodGet, "/resource/1/rating").
		WithHeader("Authorization", token).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		rat := rating.GetResult{}
		marshallErr := json.Unmarshal([]byte(b), &rat)
		assert.NoError(t, marshallErr)

		assert.Equal(t, 5, rat.Rating)
	})
}

func TestGet_Http_RatingNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:read scope
	user, token, err := tc.UserWithScopes("foo", "rating:read")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	GetChecker(tc).Test(t, http.MethodGet, "/resource/3/rating").
		WithHeader("Authorization", token).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		rat := rating.GetResult{}
		marshallErr := json.Unmarshal([]byte(b), &rat)
		assert.NoError(t, marshallErr)

		assert.Equal(t, -1, rat.Rating)
	})
}

func TestGet_Http_ResourceNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:read scope
	user, token, err := tc.UserWithScopes("foo", "rating:read")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	GetChecker(tc).Test(t, http.MethodGet, "/resource/99/rating").
		WithHeader("Authorization", token).Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", err.Name)
	})
}

func UpdateChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, tc.JWTSigningKey())
	checker := goahttpcheck.New()
	checker.Mount(server.NewUpdateHandler,
		server.MountUpdateHandler,
		rating.NewUpdateEndpoint(New(tc), service.JWTAuth))
	return checker
}

func TestUpdate_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:write scope
	user, token, err := tc.UserWithScopes("foo", "rating:write")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	data := []byte(`{"rating": 5}`)

	UpdateChecker(tc).Test(t, http.MethodPut, "/resource/3/rating").
		WithHeader("Authorization", token).
		WithBody(data).Check().
		HasStatus(200)

	r := model.UserResourceRating{ResourceID: 3, UserID: user.ID}
	err = tc.DB().First(&r).Error
	assert.NoError(t, err)

	assert.Equal(t, uint(5), r.Rating)
}

func TestUpdate_Http_Existing(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:write scope
	user, token, err := tc.UserWithScopes("foo", "rating:write")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	data := []byte(`{"rating": 2}`)

	UpdateChecker(tc).Test(t, http.MethodPut, "/resource/1/rating").
		WithHeader("Authorization", token).
		WithBody(data).Check().
		HasStatus(200)

	r := model.UserResourceRating{ResourceID: 1, UserID: user.ID}
	err = tc.DB().First(&r).Error
	assert.NoError(t, err)

	assert.Equal(t, uint(2), r.Rating)
}

func TestUpdate_Http_ResourceNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user with rating:write scope
	user, token, err := tc.UserWithScopes("foo", "rating:write")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	data := []byte(`{"rating": 2}`)

	UpdateChecker(tc).Test(t, http.MethodPut, "/resource/99/rating").
		WithHeader("Authorization", token).
		WithBody(data).Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", err.Name)
	})
}
