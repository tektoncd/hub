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
	"crypto/sha256"
	"encoding/hex"

	"github.com/markbates/goth"
	"github.com/tektoncd/hub/api/pkg/auth/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
	"gorm.io/gorm"
)

func (r *request) userScopes(user *model.User) ([]string, error) {

	var userScopes []string = r.defaultScopes

	q := r.db.Preload("Scopes").Where(&model.User{GithubLogin: user.GithubLogin})

	dbUser := model.User{}
	if err := q.Find(&dbUser).Error; err != nil {
		r.log.Error(err)
		return nil, err
	}

	for _, s := range dbUser.Scopes {
		userScopes = append(userScopes, s.Name)
	}

	return userScopes, nil
}

func (r *request) createTokens(user *model.User, scopes []string) (*app.AuthenticateResult, error) {

	req := token.Request{
		User:      user,
		Scopes:    scopes,
		JWTConfig: r.jwtConfig,
	}

	accessToken, accessExpiresAt, err := req.AccessJWT()
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := req.RefreshJWT()
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	user.RefreshTokenChecksum = createChecksum(refreshToken)

	if err = r.db.Save(user).Error; err != nil {
		r.log.Error(err)
		return nil, err
	}

	data := &app.AuthTokens{
		Access: &app.Token{
			Token:           accessToken,
			RefreshInterval: r.jwtConfig.AccessExpiresIn.String(),
			ExpiresAt:       accessExpiresAt,
		},
		Refresh: &app.Token{
			Token:           refreshToken,
			RefreshInterval: r.jwtConfig.RefreshExpiresIn.String(),
			ExpiresAt:       refreshExpiresAt,
		},
	}

	return &app.AuthenticateResult{Data: data}, nil
}

func createChecksum(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *request) insertData(ghUser goth.User, code string) error {

	// Check if user exist
	q := r.db.Model(&model.User{}).
		Where("github_login = ?", ghUser.NickName)

	var user model.User
	err := q.First(&user).Error
	if err != nil {
		// If user doesn't exist, create a new record
		if err == gorm.ErrRecordNotFound {

			user.GithubName = ghUser.Name
			user.GithubLogin = ghUser.NickName
			user.Type = model.NormalUserType
			user.AvatarURL = ghUser.AvatarURL
			user.Code = code

			err = r.db.Create(&user).Error
			if err != nil {
				r.log.Error(err)
				return err
			}
		}
	} else {
		if err := r.db.Model(&model.User{}).Where("github_login = ?", ghUser.NickName).Update("code", code).Error; err != nil {
			r.log.Error(err)
			return err
		}
	}

	// User already exist, check if GitHub Name is empty
	// If Name is empty, then user is inserted through config.yaml
	// Update user with remaining details

	if user.GithubName == "" {
		user.GithubName = ghUser.Name
		user.Type = model.NormalUserType
	}

	// For existing user, check if URL is not added
	if user.AvatarURL == "" {
		user.AvatarURL = ghUser.AvatarURL
		if err = r.db.Save(&user).Error; err != nil {
			r.log.Error(err)
			return err
		}
	}

	return nil
}
