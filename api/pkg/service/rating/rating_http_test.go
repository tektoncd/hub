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
	goa "goa.design/goa/v3/pkg"

	"github.com/tektoncd/hub/api/gen/http/rating/server"
	"github.com/tektoncd/hub/api/gen/rating"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

// Token for the user with github name "foo-bar" and github login "foo"
const validToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MTEsImxvZ2luIjoiZm9vIiwibmFtZSI6ImZvby1iYXIiLCJzY29wZXMiOlsicmF0aW5nOnJlYWQiLCJyYXRpbmc6d3JpdGUiXX0." +
	"AnQtXmPfyAE22XqYo95mzfyynsqr6pe5GADXutsmRaM"

// Token with Invalid Scopes
const tokenWithInvalidScope = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MTIzLCJsb2dpbiI6ImZvbyIsIm5hbWUiOiJmb28tYmFyIiwic2NvcGVzIjpbImludmFsaWQ6c2NvcGUiXX0." +
	"SEBQUE9aG8zHDuyAlV5R20h63-TBQjlDEyXxdPmCIX4"

func GetChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	validate := &auth.Validator{DB: tc.DB(), JWTKey: tc.JWTSigningKey()}
	checker := goahttpcheck.New()
	checker.Mount(server.NewGetHandler,
		server.MountGetHandler,
		rating.NewGetEndpoint(New(tc), validate.JWTAuth))
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

		var err *goa.ServiceError
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "invalid-token", err.Name)
	})
}

func TestGet_Http_InvalidScopes(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	GetChecker(tc).Test(t, http.MethodGet, "/resource/1/rating").
		WithHeader("Authorization", tokenWithInvalidScope).Check().
		HasStatus(403).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var err *goa.ServiceError
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "invalid-scopes", err.Name)
	})
}

func TestGet_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	GetChecker(tc).Test(t, http.MethodGet, "/resource/1/rating").
		WithHeader("Authorization", validToken).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var rat *rating.GetResult
		marshallErr := json.Unmarshal([]byte(b), &rat)
		assert.NoError(t, marshallErr)

		assert.Equal(t, 5, rat.Rating)
	})
}

func TestGet_Http_RatingNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	GetChecker(tc).Test(t, http.MethodGet, "/resource/3/rating").
		WithHeader("Authorization", validToken).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var rat *rating.GetResult
		marshallErr := json.Unmarshal([]byte(b), &rat)
		assert.NoError(t, marshallErr)

		assert.Equal(t, -1, rat.Rating)
	})
}

func TestGet_Http_ResourceNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	GetChecker(tc).Test(t, http.MethodGet, "/resource/99/rating").
		WithHeader("Authorization", validToken).Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var err *goa.ServiceError
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", err.Name)
	})
}

func UpdateChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	validate := &auth.Validator{DB: tc.DB(), JWTKey: tc.JWTSigningKey()}
	checker := goahttpcheck.New()
	checker.Mount(server.NewUpdateHandler,
		server.MountUpdateHandler,
		rating.NewUpdateEndpoint(New(tc), validate.JWTAuth))
	return checker
}

func TestUpdate_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	data := []byte(`{"rating": 2}`)

	UpdateChecker(tc).Test(t, http.MethodPut, "/resource/1/rating").
		WithHeader("Authorization", validToken).
		WithBody(data).Check().
		HasStatus(200)
}

func TestUpdate_Http_ResourceNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	data := []byte(`{"rating": 2}`)

	UpdateChecker(tc).Test(t, http.MethodPut, "/resource/99/rating").
		WithHeader("Authorization", validToken).
		WithBody(data).Check().
		HasStatus(404).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var err *goa.ServiceError
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "not-found", err.Name)
	})
}
