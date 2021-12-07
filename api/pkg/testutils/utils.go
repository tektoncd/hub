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

package testutils

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
)

// Now mocks the current time
var Now = func() time.Time {
	return time.Date(2020, 01, 01, 12, 00, 00, 01234567, time.UTC)
}

// NowAfterDuration returns a time one second after adding the duration in Now
func NowAfterDuration(dur time.Duration) func() time.Time {
	return func() time.Time {
		return Now().Add(dur).Add(1 * time.Second)
	}
}

// FormatJSON formats json string to be added to golden file
func FormatJSON(b []byte) (string, error) {
	var formatted bytes.Buffer
	err := json.Indent(&formatted, b, "", "\t")
	if err != nil {
		return "", err
	}
	return formatted.String(), nil
}

// UserWithScopes returns JWT for user with required scopes
// User will have same github login and github name in db
func (tc *TestConfig) UserWithScopes(name, email string, scopes ...string) (*model.User, string, error) {

	user := &model.User{Type: model.NormalUserType, Email: email}
	if err := tc.DB().Where(&model.User{Email: email}).
		FirstOrCreate(user).Error; err != nil {
		return nil, "", err
	}

	account := &model.Account{Name: name, UserName: name, UserID: user.ID, Provider: "github"}
	if err := tc.DB().Where(&model.Account{UserID: user.ID, UserName: name}).
		FirstOrCreate(account).Error; err != nil {
		return nil, "", err
	}

	if err := tc.AddScopesForUser(user.ID, scopes); err != nil {
		return nil, "", err
	}

	token.Now = Now

	req := token.Request{User: user, Scopes: scopes, JWTConfig: tc.JWTConfig(), Provider: "github"}
	accessToken, _, err := req.AccessJWT()
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
}

// RefreshTokenForUser returns refresh JWT for user with refresh:token scope
// User will have same github login and github name in db
func (tc *TestConfig) RefreshTokenForUser(name, email string) (*model.User, string, error) {

	user := &model.User{Type: model.NormalUserType, Email: email}
	if err := tc.DB().Where(&model.User{Email: email}).
		FirstOrCreate(user).Error; err != nil {
		return nil, "", err
	}

	token.Now = Now

	req := token.Request{User: user, JWTConfig: tc.JWTConfig(), Provider: "github"}
	refreshToken, _, err := req.RefreshJWT()
	if err != nil {
		return nil, "", err
	}

	return user, refreshToken, nil
}

// AgentWithScopes returns JWT for user with required scopes
func (tc *TestConfig) AgentWithScopes(name string, scopes ...string) (*model.User, string, error) {

	agent := &model.User{AgentName: name, Type: model.AgentUserType}
	if err := tc.DB().Where(&model.User{AgentName: name}).
		FirstOrCreate(agent).Error; err != nil {
		return nil, "", err
	}

	if err := tc.AddScopesForUser(agent.ID, scopes); err != nil {
		return nil, "", err
	}

	token.Now = Now

	req := token.Request{User: agent, Scopes: scopes, JWTConfig: tc.JWTConfig()}
	agentToken, err := req.AgentJWT()
	if err != nil {
		return nil, "", err
	}

	return agent, agentToken, nil
}

// AddScopesForUser adds scopes for passed User ID
func (tc *TestConfig) AddScopesForUser(userID uint, scopes []string) error {

	for _, s := range scopes {

		// rating scopes are not saved in db, they are in Data object read from config
		if s == "rating:read" || s == "rating:write" {
			continue
		}

		scope := &model.Scope{}
		if err := tc.DB().Where(&model.Scope{Name: s}).First(&scope).Error; err != nil {
			return err
		}

		us := &model.UserScope{UserID: userID, ScopeID: scope.ID}
		if err := tc.DB().FirstOrCreate(us).Error; err != nil {
			return err
		}
	}
	return nil
}
