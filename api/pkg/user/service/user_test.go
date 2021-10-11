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
	testUser, accessToken, err := tc.UserWithScopes("abc", "rating:read", "rating:write", "agent:create")
	assert.Equal(t, testUser.GithubLogin, "abc")
	assert.NoError(t, err)

	// Mocks the time
	jwt.TimeFunc = testutils.Now

	req, err := http.NewRequest("GET", "/user/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	userSvc := New(tc)
	jwt := UserService{
		JwtConfig: tc.JWTConfig(),
	}

	req.Header.Set("Authorization", accessToken)
	handler := http.HandlerFunc(jwt.JWTAuth(userSvc.Info))
	assert.NoError(t, err)

	handler.ServeHTTP(res, req)

	var u *userApp.InfoResult
	err = json.Unmarshal(res.Body.Bytes(), &u)
	assert.NoError(t, err)

	assert.Equal(t, "abc", u.Data.GithubID)
	assert.Equal(t, "https://bar", u.Data.AvatarURL)
}
