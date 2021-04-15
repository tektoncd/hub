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

package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/http/user/server"
	"github.com/tektoncd/hub/api/gen/user"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/testutils"
	"github.com/tektoncd/hub/api/pkg/token"
	goa "goa.design/goa/v3/pkg"
)

func RefreshAccessTokenChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, "admin")
	checker := goahttpcheck.New()
	checker.Mount(server.NewRefreshAccessTokenHandler,
		server.MountRefreshAccessTokenHandler,
		user.NewRefreshAccessTokenEndpoint(New(tc), service.JWTAuth))
	return checker
}

func TestRefreshAccessToken_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("abc")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now
	token.Now = testutils.Now

	RefreshAccessTokenChecker(tc).Test(t, http.MethodPost, "/user/refresh/accesstoken").
		WithHeader("Authorization", refreshToken).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &user.RefreshAccessTokenResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		// expected access jwt
		user, accessToken, err := tc.UserWithScopes("abc", "rating:read", "rating:write")
		assert.Equal(t, user.GithubName, "abc")
		assert.NoError(t, err)

		assert.Equal(t, accessToken, res.Data.Access.Token)
	})
}

func TestRefreshAccessToken_Http_ExpiredRefreshToken(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("abc")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.NowAfterDuration(tc.JWTConfig().RefreshExpiresIn)

	RefreshAccessTokenChecker(tc).Test(t, http.MethodPost, "/user/refresh/accesstoken").
		WithHeader("Authorization", refreshToken).Check().
		HasStatus(401).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := &goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)
		assert.EqualError(t, err, "invalid or expired user token")
	})
}

func TestRefreshAccessToken_Http_RefreshTokenChecksumIsDifferent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("foo")
	assert.Equal(t, testUser.GithubLogin, "foo")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	RefreshAccessTokenChecker(tc).Test(t, http.MethodPost, "/user/refresh/accesstoken").
		WithHeader("Authorization", refreshToken).Check().
		HasStatus(401).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := &goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)
		assert.EqualError(t, err, "invalid refresh token")
	})
}

func NewRefreshTokenChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, "admin")
	checker := goahttpcheck.New()
	checker.Mount(server.NewNewRefreshTokenHandler,
		server.MountNewRefreshTokenHandler,
		user.NewNewRefreshTokenEndpoint(New(tc), service.JWTAuth))
	return checker
}

func TestNewRefreshToken_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("abc")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now
	token.Now = testutils.Now

	NewRefreshTokenChecker(tc).Test(t, http.MethodPost, "/user/refresh/refreshtoken").
		WithHeader("Authorization", refreshToken).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &user.NewRefreshTokenResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)

		// user refresh token
		testUser, refreshToken, err := tc.RefreshTokenForUser("abc")
		assert.Equal(t, testUser.GithubLogin, "abc")
		assert.NoError(t, err)

		refreshExpiryTime := testutils.Now().Add(tc.JWTConfig().RefreshExpiresIn).Unix()

		assert.Equal(t, refreshToken, res.Data.Refresh.Token)
		assert.Equal(t, tc.JWTConfig().RefreshExpiresIn.String(), res.Data.Refresh.RefreshInterval)
		assert.Equal(t, refreshExpiryTime, res.Data.Refresh.ExpiresAt)
	})
}

func TestNewRefreshToken_Http_RefreshTokenChecksumIsDifferent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("foo")
	assert.Equal(t, testUser.GithubLogin, "foo")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	NewRefreshTokenChecker(tc).Test(t, http.MethodPost, "/user/refresh/refreshtoken").
		WithHeader("Authorization", refreshToken).Check().
		HasStatus(401).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		err := &goa.ServiceError{}
		marshallErr := json.Unmarshal([]byte(b), &err)
		assert.NoError(t, marshallErr)
		assert.EqualError(t, err, "invalid refresh token")
	})
}

func InfoChecker(tc *testutils.TestConfig) *goahttpcheck.APIChecker {
	service := auth.NewService(tc.APIConfig, "admin")
	checker := goahttpcheck.New()
	checker.Mount(server.NewInfoHandler,
		server.MountInfoHandler,
		user.NewInfoEndpoint(New(tc), service.JWTAuth))
	return checker
}
func TestUserInfo_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, accessToken, err := tc.UserWithScopes("abc", "rating:read")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	InfoChecker(tc).Test(t, http.MethodGet, "/user/info").
		WithHeader("Authorization", accessToken).Check().
		HasStatus(200).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		res := &user.InfoResult{}
		marshallErr := json.Unmarshal([]byte(b), &res)
		assert.NoError(t, marshallErr)
		assert.Equal(t, "abc", res.Data.Name)
	})
}
