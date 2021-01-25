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
	"github.com/tektoncd/hub/api/pkg/service/auth"
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
	ctx := auth.WithUserID(context.Background(), testUser.ID)
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
	ctx := auth.WithUserID(context.Background(), testUser.ID)
	payload := &user.RefreshAccessTokenPayload{RefreshToken: refreshToken}
	_, err = userSvc.RefreshAccessToken(ctx, payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid refresh token")
}
