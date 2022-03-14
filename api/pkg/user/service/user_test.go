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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/testutils"
	userApp "github.com/tektoncd/hub/api/pkg/user/app"
)

func TestInfo(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user Access Token
	testUser, accessToken, err := tc.UserWithScopes("abc", "abc@bar.com", "rating:read", "rating:write", "agent:create")
	assert.Equal(t, testUser.Email, "abc@bar.com")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("GET", "/user/info", nil)
	if err != nil {
		t.Fatal(err)
	}
	// add cookie in the request with accessToken name
	req.AddCookie(tc.CreateCookie("accessToken", accessToken))

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.Info))

	assert.NoError(t, err)

	handler.ServeHTTP(res, req)

	var u *userApp.InfoResult

	err = json.Unmarshal(res.Body.Bytes(), &u)
	assert.NoError(t, err)

	assert.Equal(t, "abc", u.Data.UserName)
	assert.Equal(t, "https://abc", u.Data.AvatarURL)
}

func TestRefreshAccessToken(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("abc", "abc@bar.com")
	assert.Equal(t, testUser.Email, "abc@bar.com")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("POST", "/user/refresh/accesstoken", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add cookie in the request with refreshToken name
	req.AddCookie(tc.CreateCookie("refreshToken", refreshToken))

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.RefreshAccessToken))
	assert.NoError(t, err)

	handler.ServeHTTP(res, req)

	// expected access jwt for user
	user, accessToken, err := tc.UserWithScopes("abc", "abc@bar.com", "rating:read", "rating:write")
	assert.Equal(t, user.Email, "abc@bar.com")
	assert.NoError(t, err)

	var u *userApp.RefreshAccessTokenResult
	err = json.Unmarshal(res.Body.Bytes(), &u)
	assert.NoError(t, err)

	accessExpiryTime := testutils.Now().Add(tc.JWTConfig().AccessExpiresIn).Unix()

	assert.Equal(t, accessToken, u.Data.Access.Token)
	assert.Equal(t, tc.JWTConfig().AccessExpiresIn.String(), u.Data.Access.RefreshInterval)
	assert.Equal(t, accessExpiryTime, u.Data.Access.ExpiresAt)
}

func TestRefreshAccessToken_RefreshTokenChecksumIsDifferent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("foo", "foo@bar.com")
	assert.Equal(t, testUser.Email, "foo@bar.com")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("POST", "/user/refresh/accesstoken", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add cookie in the request with refreshToken name
	req.AddCookie(tc.CreateCookie("refreshToken", refreshToken))

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	// req.Header.Set("Authorization", refreshToken)
	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.RefreshAccessToken))
	handler.ServeHTTP(res, req)

	assert.Equal(t, res.Body.String(), "invalid refresh token\n")
}

func TestNewRefreshToken(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("abc", "abc@bar.com")
	assert.Equal(t, testUser.Email, "abc@bar.com")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("POST", "/user/refresh/refreshtoken", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add cookie in the request with accessToken name
	req.AddCookie(tc.CreateCookie("refreshToken", refreshToken))

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.NewRefreshToken))
	assert.NoError(t, err)

	handler.ServeHTTP(res, req)

	// user refresh token
	testUser, refreshToken, err = tc.RefreshTokenForUser("abc", "abc@bar.com")
	assert.Equal(t, testUser.Email, "abc@bar.com")
	assert.NoError(t, err)

	var u *userApp.NewRefreshTokenResult
	err = json.Unmarshal(res.Body.Bytes(), &u)
	assert.NoError(t, err)

	refreshExpiryTime := testutils.Now().Add(tc.JWTConfig().RefreshExpiresIn).Unix()

	assert.Equal(t, refreshToken, u.Data.Refresh.Token)
	assert.Equal(t, tc.JWTConfig().RefreshExpiresIn.String(), u.Data.Refresh.RefreshInterval)
	assert.Equal(t, refreshExpiryTime, u.Data.Refresh.ExpiresAt)
}

func TestNewRefreshToken_RefreshTokenChecksumIsDifferent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("foo", "foo@bar.com")

	assert.Equal(t, testUser.Email, "foo@bar.com")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("POST", "/user/refresh/refreshtoken", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add cookie in the request with accessToken name
	req.AddCookie(tc.CreateCookie("refreshToken", refreshToken))

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.NewRefreshToken))
	handler.ServeHTTP(res, req)

	assert.Equal(t, res.Body.String(), "invalid refresh token\n")
}

func TestLogout(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	testUser, accessToken, err := tc.UserWithScopes("abc", "abc@bar.com", "rating:read", "rating:write")
	assert.Equal(t, testUser.Email, "abc@bar.com")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("GET", "/user/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add cookie in the request with accessToken name
	req.AddCookie(tc.CreateCookie("accessToken", accessToken))

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.Logout))

	assert.NoError(t, err)

	handler.ServeHTTP(res, req)

	var u *userApp.ClearCookies

	err = json.Unmarshal(res.Body.Bytes(), &u)
	assert.NoError(t, err)

	assert.Equal(t, true, u.Data)
}

func TestGetAccessToken(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	testUser, accessToken, err := tc.UserWithScopes("abc", "abc@bar.com", "rating:read", "rating:write")
	assert.Equal(t, testUser.Email, "abc@bar.com")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("GET", "/user/accesstoken", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add cookie in the request with accessToken name
	req.AddCookie(tc.CreateCookie("accessToken", accessToken))

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.GetAccessToken))

	assert.NoError(t, err)

	handler.ServeHTTP(res, req)

	var u *userApp.ExitingAccessToken

	err = json.Unmarshal(res.Body.Bytes(), &u)
	assert.NoError(t, err)

	assert.Equal(t, accessToken, u.Data)
}
