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

package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
	authApp "github.com/tektoncd/hub/api/pkg/auth/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/testutils"
	"github.com/tektoncd/hub/api/pkg/token"
)

func TestLogin(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// Mocks the time
	token.Now = testutils.Now

	authSvc := New(tc)

	provider = "github"
	req, err := http.NewRequest("POST", "/auth/login?code=test-code", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(authSvc.HubAuthenticate)
	assert.NoError(t, err)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(res, req)

	var u *authApp.AuthenticateResult
	err = json.Unmarshal(res.Body.Bytes(), &u)
	assert.NoError(t, err)

	// expected access jwt for user
	user, accessToken, err := tc.UserWithScopes("foo", "foo@bar.com", "rating:read", "rating:write", "agent:create")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	// expected refresh jwt for user
	user, refreshToken, err := tc.RefreshTokenForUser("foo", "foo@bar.com")
	assert.Equal(t, user.Email, "foo@bar.com")
	assert.NoError(t, err)

	accessExpiryTime := testutils.Now().Add(tc.JWTConfig().AccessExpiresIn).Unix()
	refreshExpiryTime := testutils.Now().Add(tc.JWTConfig().RefreshExpiresIn).Unix()

	assert.Equal(t, accessToken, u.Data.Access.Token)
	assert.Equal(t, tc.JWTConfig().AccessExpiresIn.String(), u.Data.Access.RefreshInterval)
	assert.Equal(t, accessExpiryTime, u.Data.Access.ExpiresAt)

	assert.Equal(t, refreshToken, u.Data.Refresh.Token)
	assert.Equal(t, tc.JWTConfig().RefreshExpiresIn.String(), u.Data.Refresh.RefreshInterval)
	assert.Equal(t, refreshExpiryTime, u.Data.Refresh.ExpiresAt)
}

func TestInvalidLogin(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// Mocks the time
	token.Now = testutils.Now

	authSvc := New(tc)

	provider = "github"
	req, err := http.NewRequest("POST", "/auth/login?code=fake-code", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()
	http.HandlerFunc(authSvc.HubAuthenticate).ServeHTTP(res, req)

	assert.Equal(t, res.Body.String(), "record not found\n")
	assert.Equal(t, res.Code, 400)
}

func TestProviderList(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth/providers", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(List)
	assert.NoError(t, err)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(res, req)

	var provider *authApp.ProviderList
	err = json.Unmarshal(res.Body.Bytes(), &provider)
	assert.NoError(t, err)

	assert.Equal(t, "github", provider.Data[0].Name)
}

func TestInsertData_AccountExistsButNotEmail(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// Mocks the time
	token.Now = testutils.Now

	req := request{
		db:            tc.DB(),
		log:           tc.Logger("svc"),
		defaultScopes: []string{"rating:write", "rating:read"},
		jwtConfig:     tc.JWTConfig(),
		provider:      "bitbucket",
	}

	gitUser := goth.User{
		Email:     "bbuser@bar.com",
		NickName:  "bbuser",
		Name:      "bitbucketuser",
		AvatarURL: "http://bitbucketavatar",
	}

	err := req.insertData(gitUser, "code", req.provider)
	assert.NoError(t, err)

	userQuery := tc.DB().Model(&model.User{}).
		Where("email = ?", gitUser.Email)
	err = userQuery.First(&model.User{}).Error
	assert.NoError(t, err)
}

func TestInsertData_AddNewEmailAndAccount(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// Mocks the time
	token.Now = testutils.Now

	req := request{
		db:            tc.DB(),
		log:           tc.Logger("svc"),
		defaultScopes: []string{"rating:write", "rating:read"},
		jwtConfig:     tc.JWTConfig(),
		provider:      "bitbucket",
	}

	gitUser := goth.User{
		Email:     "bbuser@bar.com",
		NickName:  "bbnewuser",
		Name:      "newbitbucketuser",
		AvatarURL: "http://bitbucketavatar",
	}

	// check whether email already exists or not
	// it should return an error
	userQuery := tc.DB().Model(&model.User{}).
		Where("email = ?", gitUser.Email)
	err := userQuery.First(&model.User{}).Error
	assert.Error(t, err)

	// check whether username with that provider already exists or not
	// it should return an error
	accountQuery := tc.DB().Model(&model.Account{}).
		Where(model.Account{Name: gitUser.NickName, Provider: req.provider})
	err = accountQuery.First(&model.Account{}).Error
	assert.Error(t, err)

	err = req.insertData(gitUser, "code", req.provider)
	assert.NoError(t, err)
}

func TestInsertData_EmailExistsAddNewAccount(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	// Mocks the time
	token.Now = testutils.Now

	req := request{
		db:            tc.DB(),
		log:           tc.Logger("svc"),
		defaultScopes: []string{"rating:write", "rating:read"},
		jwtConfig:     tc.JWTConfig(),
		provider:      "gitlab",
	}

	gitUser := goth.User{
		Email:     "foo@bar.com",
		NickName:  "gitlabuser",
		Name:      "gitlabuser",
		AvatarURL: "http://gitlabavatar",
	}

	user := model.User{}

	// Email should exist
	// If it doesn't exists then test should fail
	userQuery := tc.DB().Model(&model.User{}).
		Where("email = ?", gitUser.Email)
	err := userQuery.Last(&user).Error
	assert.NoError(t, err)

	// Check for no of accounts associated with particular email
	// The count should be 2. If not 2 then test should fail
	accounts := []model.Account{}
	accountQuery := tc.DB().Model(&model.Account{}).Where("user_id = ?", user.ID)
	err = accountQuery.Find(&accounts).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(accounts))

	// check whether username with that provider already exists or not
	// Test should fail if account doesn't exists
	accountQuery = tc.DB().Model(&model.Account{}).
		Where(model.Account{Name: gitUser.NickName, Provider: req.provider})
	err = accountQuery.First(&model.Account{}).Error
	assert.Error(t, err)

	// Insert new account data
	err = req.insertData(gitUser, "code", req.provider)
	assert.NoError(t, err)

	// Check for no of accounts associated with particular email
	// The count should be 3 as a new account is inserted. If not 3 then test should fail
	accountQuery = tc.DB().Model(&model.Account{}).Where("user_id = ?", user.ID)
	err = accountQuery.Find(&accounts).Error
	assert.NoError(t, err)
	assert.Equal(t, 3, len(accounts))
}
