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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/auth"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/testutils"
	"github.com/tektoncd/hub/api/pkg/token"
	"gopkg.in/h2non/gock.v1"
)

func TestLogin(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Off()

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

	// Mocks the time
	token.Now = testutils.Now

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "test-code"}
	res, err := authSvc.Authenticate(context.Background(), payload)
	assert.NoError(t, err)

	// expected access jwt for user
	user, accessToken, err := tc.UserWithScopes("test", "rating:read", "rating:write")
	assert.Equal(t, user.GithubLogin, "test")
	assert.NoError(t, err)

	// expected refresh jwt for user
	user, refreshToken, err := tc.RefreshTokenForUser("test")
	assert.Equal(t, user.GithubLogin, "test")
	assert.NoError(t, err)

	accessExpiryTime := testutils.Now().Add(tc.JWTConfig().AccessExpiresIn).Unix()
	refreshExpiryTime := testutils.Now().Add(tc.JWTConfig().RefreshExpiresIn).Unix()

	assert.Equal(t, accessToken, res.Data.Access.Token)
	assert.Equal(t, tc.JWTConfig().AccessExpiresIn.String(), res.Data.Access.RefreshInterval)
	assert.Equal(t, accessExpiryTime, res.Data.Access.ExpiresAt)

	assert.Equal(t, refreshToken, res.Data.Refresh.Token)
	assert.Equal(t, tc.JWTConfig().RefreshExpiresIn.String(), res.Data.Refresh.RefreshInterval)
	assert.Equal(t, refreshExpiryTime, res.Data.Refresh.ExpiresAt)

	assert.Equal(t, gock.IsDone(), true)
}

func TestLogin_again(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Off()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		JSON(map[string]string{
			"access_token": "test-token",
		})

	gock.New("https://api.github.com").
		Get("/user").
		Times(2).
		Reply(200).
		JSON(map[string]string{
			"login": "test",
			"name":  "test-user",
		})

	// Mocks the time
	token.Now = testutils.Now

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "test-code"}
	res, err := authSvc.Authenticate(context.Background(), payload)
	assert.NoError(t, err)

	// expected access jwt for user
	user, accessToken, err := tc.UserWithScopes("test", "rating:read", "rating:write")
	assert.Equal(t, user.GithubLogin, "test")
	assert.NoError(t, err)

	// expected refresh jwt for user
	user, refreshToken, err := tc.RefreshTokenForUser("test")
	assert.Equal(t, user.GithubLogin, "test")
	assert.NoError(t, err)

	assert.Equal(t, accessToken, res.Data.Access.Token)
	assert.Equal(t, refreshToken, res.Data.Refresh.Token)

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		JSON(map[string]string{
			"access_token": "test-token-2",
		})

	payloadAgain := &auth.AuthenticatePayload{Code: "test-code-2"}
	resAgain, err := authSvc.Authenticate(context.Background(), payloadAgain)
	assert.NoError(t, err)

	assert.Equal(t, accessToken, resAgain.Data.Access.Token)
	assert.Equal(t, refreshToken, resAgain.Data.Refresh.Token)
	assert.Equal(t, gock.IsDone(), true)
}

func TestLogin_InvalidCode(t *testing.T) {
	tc := testutils.Setup(t)

	defer gock.Off()

	gock.New("https://github.com").
		Post("/login/oauth/access_token").
		Reply(200).
		SetError(errors.New("oauth2: server response missing access_token"))

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "test-code"}
	_, err := authSvc.Authenticate(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid authorization code")
	assert.Equal(t, gock.IsDone(), true)
}

func TestLogin_UserWithExtraScope(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Off()

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

	// Mocks the time
	token.Now = testutils.Now

	// foo user is fetched from db to check its existing token checksum
	user, _, err := tc.UserWithScopes("foo")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	// validate existing checksum
	ut := &model.User{}
	err = tc.DB().First(ut, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "test-checksum", ut.RefreshTokenChecksum)

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "foo-test"}
	res, err := authSvc.Authenticate(context.Background(), payload)
	assert.NoError(t, err)

	// expected access jwt for user
	user, accessToken, err := tc.UserWithScopes("foo", "rating:read", "rating:write", "agent:create")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	// expected refresh jwt for user
	user, refreshToken, err := tc.RefreshTokenForUser("foo")
	assert.Equal(t, user.GithubLogin, "foo")
	assert.NoError(t, err)

	assert.Equal(t, accessToken, res.Data.Access.Token)
	assert.Equal(t, refreshToken, res.Data.Refresh.Token)
	assert.Equal(t, gock.IsDone(), true)

	// validate the new checksum which overrides existing refresh token checksum
	ut = &model.User{}
	err = tc.DB().First(ut, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, ut.RefreshTokenChecksum, createChecksum(refreshToken))
}

func TestLogin_UserAddedByConfig(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	defer gock.Off()

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
			"login":      "Config-user",
			"name":       "config-user",
			"avatar_url": "http://config",
		})

	// Mocks the time
	token.Now = testutils.Now

	authSvc := New(tc)
	payload := &auth.AuthenticatePayload{Code: "test-code"}
	res, err := authSvc.Authenticate(context.Background(), payload)
	assert.NoError(t, err)


	// expected access jwt for user
	user, accessToken, err := tc.UserWithScopes("Config-user", "rating:read", "rating:write", "config:refresh")
	assert.Equal(t, user.GithubLogin, "Config-user")
	assert.NoError(t, err)

	// validate the avatar_url of user after login
	ut := &model.User{}
	err = tc.DB().First(ut, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Config-user", ut.GithubLogin)
	assert.Equal(t, "http://config", ut.AvatarURL)

	// expected refresh jwt for user
	user, refreshToken, err := tc.RefreshTokenForUser("Config-user")
	assert.Equal(t, user.GithubLogin, "Config-user")
	assert.NoError(t, err)

	accessExpiryTime := testutils.Now().Add(tc.JWTConfig().AccessExpiresIn).Unix()
	refreshExpiryTime := testutils.Now().Add(tc.JWTConfig().RefreshExpiresIn).Unix()

	assert.Equal(t, accessToken, res.Data.Access.Token)
	assert.Equal(t, tc.JWTConfig().AccessExpiresIn.String(), res.Data.Access.RefreshInterval)
	assert.Equal(t, accessExpiryTime, res.Data.Access.ExpiresAt)

	assert.Equal(t, refreshToken, res.Data.Refresh.Token)
	assert.Equal(t, tc.JWTConfig().RefreshExpiresIn.String(), res.Data.Refresh.RefreshInterval)
	assert.Equal(t, refreshExpiryTime, res.Data.Refresh.ExpiresAt)

	assert.Equal(t, gock.IsDone(), true)
}
