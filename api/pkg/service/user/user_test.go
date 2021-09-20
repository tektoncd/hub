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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/user"
	"github.com/tektoncd/hub/api/pkg/service/validator"
	"github.com/tektoncd/hub/api/pkg/testutils"
	"github.com/tektoncd/hub/api/pkg/token"
)

func TestRefreshAccessToken(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("abc")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	// Mocks the time
	token.Now = testutils.Now

	userSvc := New(tc)
	ctx := validator.WithUserID(context.Background(), testUser.ID)
	payload := &user.RefreshAccessTokenPayload{RefreshToken: refreshToken}
	res, err := userSvc.RefreshAccessToken(ctx, payload)
	assert.NoError(t, err)

	// expected access jwt for user
	user, accessToken, err := tc.UserWithScopes("abc", "rating:read", "rating:write")
	assert.Equal(t, user.GithubLogin, "abc")
	assert.NoError(t, err)

	accessExpiryTime := testutils.Now().Add(tc.JWTConfig().AccessExpiresIn).Unix()

	assert.Equal(t, accessToken, res.Data.Access.Token)
	assert.Equal(t, tc.JWTConfig().AccessExpiresIn.String(), res.Data.Access.RefreshInterval)
	assert.Equal(t, accessExpiryTime, res.Data.Access.ExpiresAt)
}

func TestRefreshAccessToken_RefreshTokenChecksumIsDifferent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("foo")
	assert.Equal(t, testUser.GithubLogin, "foo")
	assert.NoError(t, err)

	userSvc := New(tc)
	ctx := validator.WithUserID(context.Background(), testUser.ID)
	payload := &user.RefreshAccessTokenPayload{RefreshToken: refreshToken}
	_, err = userSvc.RefreshAccessToken(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid refresh token")
}

func TestNewRefreshToken(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("abc")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	// Mocks the time
	token.Now = testutils.Now

	userSvc := New(tc)
	ctx := validator.WithUserID(context.Background(), testUser.ID)
	payload := &user.NewRefreshTokenPayload{RefreshToken: refreshToken}
	res, err := userSvc.NewRefreshToken(ctx, payload)
	assert.NoError(t, err)

	// user refresh token
	testUser, refreshToken, err = tc.RefreshTokenForUser("abc")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	refreshExpiryTime := testutils.Now().Add(tc.JWTConfig().RefreshExpiresIn).Unix()

	assert.Equal(t, refreshToken, res.Data.Refresh.Token)
	assert.Equal(t, tc.JWTConfig().RefreshExpiresIn.String(), res.Data.Refresh.RefreshInterval)
	assert.Equal(t, refreshExpiryTime, res.Data.Refresh.ExpiresAt)
}

func TestNewRefreshToken_RefreshTokenChecksumIsDifferent(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user refresh token
	testUser, refreshToken, err := tc.RefreshTokenForUser("foo")
	assert.Equal(t, testUser.GithubLogin, "foo")
	assert.NoError(t, err)

	userSvc := New(tc)
	ctx := validator.WithUserID(context.Background(), testUser.ID)
	payload := &user.NewRefreshTokenPayload{RefreshToken: refreshToken}
	_, err = userSvc.NewRefreshToken(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid refresh token")
}

func TestInfo(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// user Access Token
	testUser, accessToken, err := tc.UserWithScopes("abc", "rating:read", "rating:write")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	userSvc := New(tc)
	ctx := validator.WithUserID(context.Background(), testUser.ID)
	payload := &user.InfoPayload{AccessToken: accessToken}
	res, err := userSvc.Info(ctx, payload)

	assert.NoError(t, err)

	assert.Equal(t, "abc", res.Data.GithubID)
	assert.Equal(t, "https://bar", res.Data.AvatarURL)

}
