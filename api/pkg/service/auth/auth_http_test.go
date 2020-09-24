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

package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	goa "goa.design/goa/v3/pkg"
	"gopkg.in/h2non/gock.v1"

	"github.com/tektoncd/hub/api/gen/auth"
	"github.com/tektoncd/hub/api/gen/http/auth/server"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

// Token for the user with github name "test-user" and github login "test"
const validToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MTAwMDEsImxvZ2luIjoidGVzdCIsIm5hbWUiOiJ0ZXN0LXVzZXIiLCJzY29wZXMiOlsicmF0aW5nOnJlYWQiLCJyYXRpbmc6d3JpdGUiXX0." +
	"zFztueyvZLLCyx3RD7WpzzfVaTrybzxgS5a_pDsq5M8"

// Token for the user with github name "foo-bar" and github login "foo"
// It has a extra scope "agent:create" along with default scope
const validTokenWithExtraScope = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJpZCI6MTEsImxvZ2luIjoiZm9vIiwibmFtZSI6ImZvby1iYXIiLCJzY29wZXMiOlsicmF0aW5nOnJlYWQiLCJyYXRpbmc6d3JpdGUiLCJhZ2VudDpjcmVhdGUiXX0." +
	"bKPINZyhzX2Ls1QM--UV56cC-vm5uymT8y-DmEhY1dE"

func AuthenticateChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	checker := goahttpcheck.New()
	checker.Mount(
		server.NewAuthenticateHandler,
		server.MountAuthenticateHandler,
		auth.NewAuthenticateEndpoint(New(tc)))
	return checker
}

func TestLogin_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Disable()
	defer gock.DisableNetworking()

	gock.New("/auth/login").
		EnableNetworking()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		JSON(map[string]string{
			"access_token": "test-token",
		})

	gock.New("https://api.github.com").
		Get("/user").
		Reply(200).
		JSON(map[string]string{
			"login": "test",
			"name":  "test-user",
		})

	AuthenticateChecker(tc).Test(t, http.MethodPost, "/auth/login?code=test").Check().
		HasStatus(http.StatusOK).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := auth.AuthenticateResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		assert.Equal(t, validToken, res.Token)
		assert.Equal(t, gock.IsDone(), true)
	})
}

func TestLogin_Http_InvalidCode(t *testing.T) {
	tc := testutils.Setup(t)

	defer gock.Disable()
	defer gock.DisableNetworking()

	gock.New("/auth/login").
		EnableNetworking()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		SetError(errors.New("oauth2: server response missing access_token"))

	AuthenticateChecker(tc).Test(t, http.MethodPost, "/auth/login?code=foo").Check().
		HasStatus(http.StatusBadRequest).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := &goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)

		assert.EqualError(t, err, "invalid authorization code")
		assert.Equal(t, gock.IsDone(), true)
	})
}

func TestLogin_Http_UserWithExtraScope(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Disable()
	defer gock.DisableNetworking()

	gock.New("/auth/login").
		EnableNetworking()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		JSON(map[string]string{
			"access_token": "foo-token",
		})

	gock.New("https://api.github.com").
		Get("/user").
		Reply(200).
		JSON(map[string]string{
			"login": "foo",
			"name":  "foo-bar",
		})

	AuthenticateChecker(tc).Test(t, http.MethodPost, "/auth/login?code=foo-test").Check().
		HasStatus(http.StatusOK).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := auth.AuthenticateResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		assert.Equal(t, validTokenWithExtraScope, res.Token)
		assert.Equal(t, gock.IsDone(), true)
	})
}
