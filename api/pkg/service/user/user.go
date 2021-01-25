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
	refreshError        = user.MakeInternalError(fmt.Errorf("failed to refresh access token"))
)

// New returns the user service implementation.
func New(api app.Config) user.Service {
	return &service{auth.NewService(api, "user"), api}
}

// Refreshes the access token of User
func (s *service) RefreshAccessToken(ctx context.Context, p *user.RefreshAccessTokenPayload) (*user.RefreshAccessTokenResult, error) {

	user, err := s.User(ctx)
	if err != nil {
		return nil, err
	}

	if user.RefreshTokenChecksum != createChecksum(p.RefreshToken) {
		return nil, invalidRefreshToken
	}

	req := request{
		db:            s.DB(ctx),
		log:           s.Logger(ctx),
		user:          user,
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
	}

	return req.refreshAccessToken()
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

	q := r.db.Preload("Scopes").Where(&model.User{GithubLogin: r.user.GithubLogin})

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
