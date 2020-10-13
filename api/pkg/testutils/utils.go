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

	"github.com/dgrijalva/jwt-go"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
)

// FormatJSON formats json string to be added to golden file
func FormatJSON(b []byte) (string, error) {
	var formatted bytes.Buffer
	err := json.Indent(&formatted, b, "", "\t")
	if err != nil {
		return "", err
	}
	return string(formatted.Bytes()), nil
}

// UserWithScopes returns JWT for user with required scopes
// User will have same github login and github name in db
func (tc *TestConfig) UserWithScopes(name string, scopes ...string) (*model.User, string, error) {

	user := &model.User{GithubLogin: name, GithubName: name, Type: model.NormalUserType}
	if err := tc.DB().Where(&model.User{GithubLogin: name}).
		FirstOrCreate(user).Error; err != nil {
		return nil, "", err
	}

	if err := tc.AddScopesForUser(user.ID, scopes); err != nil {
		return nil, "", err
	}

	claim := jwt.MapClaims{
		"id":     user.ID,
		"login":  user.GithubLogin,
		"name":   user.GithubName,
		"scopes": scopes,
	}

	token, err := token.Create(claim, tc.JWTSigningKey())
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// AgentWithScopes returns JWT for user with required scopes
func (tc *TestConfig) AgentWithScopes(name string, scopes ...string) (*model.User, string, error) {

	user := &model.User{AgentName: name, Type: model.AgentUserType}
	if err := tc.DB().Where(&model.User{AgentName: name}).
		FirstOrCreate(user).Error; err != nil {
		return nil, "", err
	}

	if err := tc.AddScopesForUser(user.ID, scopes); err != nil {
		return nil, "", err
	}

	claim := jwt.MapClaims{
		"id":     user.ID,
		"name":   user.AgentName,
		"type":   user.Type,
		"scopes": scopes,
	}

	token, err := token.Create(claim, tc.JWTSigningKey())
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
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
