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
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/gen/user"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/token"
	"gorm.io/gorm"
)

type service struct {
	*auth.Service
	api app.Config
}

type request struct {
	db            *gorm.DB
	log           *log.Logger
	user          *model.User
	defaultScopes []string
	jwtConfig     *app.JWTConfig
}

var (
	invalidRefreshToken = user.MakeInvalidToken(fmt.Errorf("invalid refresh token"))
	refreshError        = user.MakeInternalError(fmt.Errorf("failed to refresh token"))
)

// New returns the user service implementation.
func New(api app.Config) user.Service {
	return &service{auth.NewService(api, "user"), api}
}

func (s *service) newRequest(ctx context.Context, user *model.User) *request {
	return &request{
		db:            s.DB(ctx),
		log:           s.Logger(ctx),
		user:          user,
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
	}
}

// Refreshes the access token of User
func (s *service) RefreshAccessToken(ctx context.Context, p *user.RefreshAccessTokenPayload) (*user.RefreshAccessTokenResult, error) {

	user, err := s.validateRefreshToken(ctx, p.RefreshToken)
	if err != nil {
		return nil, err
	}

	return s.newRequest(ctx, user).refreshAccessToken()
}

func (s *service) validateRefreshToken(ctx context.Context, token string) (*model.User, error) {

	user, err := s.User(ctx)
	if err != nil {
		return nil, err
	}

	if user.RefreshTokenChecksum != createChecksum(token) {
		return nil, invalidRefreshToken
	}

	return user, nil
}

func (r *request) refreshAccessToken() (*user.RefreshAccessTokenResult, error) {

	scopes, err := r.userScopes()
	if err != nil {
		return nil, err
	}

	req := token.Request{
		User:      r.user,
		Scopes:    scopes,
		JWTConfig: r.jwtConfig,
	}

	accessToken, accessExpiresAt, err := req.AccessJWT()
	if err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	data := &user.AccessToken{
		Access: &user.Token{
			Token:           accessToken,
			RefreshInterval: r.jwtConfig.AccessExpiresIn.String(),
			ExpiresAt:       accessExpiresAt,
		},
	}

	return &user.RefreshAccessTokenResult{Data: data}, nil
}

func (r *request) userScopes() ([]string, error) {

	var userScopes []string = r.defaultScopes

	q := r.db.Preload("Scopes").Where(&model.User{GitUsername: r.user.GitUsername})

	dbUser := &model.User{}
	if err := q.Find(dbUser).Error; err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	for _, s := range dbUser.Scopes {
		userScopes = append(userScopes, s.Name)
	}

	return userScopes, nil
}

func createChecksum(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Refreshes the refresh token of User
func (s *service) NewRefreshToken(ctx context.Context, p *user.NewRefreshTokenPayload) (*user.NewRefreshTokenResult, error) {

	user, err := s.validateRefreshToken(ctx, p.RefreshToken)
	if err != nil {
		return nil, err
	}

	return s.newRequest(ctx, user).newRefreshToken()
}

func (r *request) newRefreshToken() (*user.NewRefreshTokenResult, error) {

	req := token.Request{
		User:      r.user,
		JWTConfig: r.jwtConfig,
	}

	refreshToken, refreshExpiresAt, err := req.RefreshJWT()
	if err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	err = r.db.Model(r.user).UpdateColumn("refresh_token_checksum", createChecksum(refreshToken)).Error
	if err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	data := &user.RefreshToken{
		Refresh: &user.Token{
			Token:           refreshToken,
			RefreshInterval: r.jwtConfig.RefreshExpiresIn.String(),
			ExpiresAt:       refreshExpiresAt,
		},
	}

	return &user.NewRefreshTokenResult{Data: data}, nil
}

// Get the user Info
func (s *service) Info(ctx context.Context, p *user.InfoPayload) (*user.InfoResult, error) {

	data, err := s.User(ctx)
	if err != nil {
		return nil, err
	}
	res := &user.InfoResult{Data: &user.UserData{
		GithubID:  data.GitUsername,
		Name:      data.Name,
		AvatarURL: data.AvatarURL,
	},
	}
	return res, nil
}
