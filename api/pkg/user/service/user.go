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
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
	userApp "github.com/tektoncd/hub/api/pkg/user/app"
	"gorm.io/gorm"
)

type request struct {
	db            *gorm.DB
	log           *app.Logger
	user          *model.User
	defaultScopes []string
	jwtConfig     *app.JWTConfig
	provider      string
}

type UserService struct {
	app.Service
	api       app.Config
	JwtConfig *app.JWTConfig
}

type Service interface {
	Info(res http.ResponseWriter, req *http.Request)
	RefreshAccessToken(res http.ResponseWriter, req *http.Request)
	NewRefreshToken(res http.ResponseWriter, req *http.Request)
}

var (
	invalidRefreshToken = fmt.Errorf("invalid refresh token")
	refreshError        = fmt.Errorf("failed to refresh token")
)

// New returns the auth service implementation.
func New(api app.Config) Service {
	return &UserService{
		Service:   api.Service("user"),
		api:       api,
		JwtConfig: api.JWTConfig(),
	}
}

// Get the user Info
func (s *UserService) Info(res http.ResponseWriter, req *http.Request) {

	id := req.Header.Get("UserID")
	provider := req.Header.Get("Provider")

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
		provider:      provider,
	}

	userId, err := ParseStringToFloat(id)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	gitUser, err := r.GitUser(int(userId))
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result := userApp.InfoResult{
		Data: &userApp.UserData{
			UserName:  gitUser.UserName,
			Name:      gitUser.Name,
			AvatarURL: gitUser.AvatarURL,
		},
	}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(result); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

// Refreshes the access token of User
func (s *UserService) RefreshAccessToken(res http.ResponseWriter, req *http.Request) {

	id := req.Header.Get("UserID")
	provider := req.Header.Get("Provider")

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
		provider:      provider,
	}

	userId, err := ParseStringToFloat(id)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken := req.Header.Get("Authorization")
	user, err := s.validateRefreshToken(int(userId), refreshToken)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := s.newRequest(user, provider).refreshAccessToken()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(result); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

}

func (r *request) refreshAccessToken() (*userApp.RefreshAccessTokenResult, error) {

	scopes, err := r.userScopes()
	if err != nil {
		return nil, err
	}

	req := token.Request{
		User:      r.user,
		Scopes:    scopes,
		JWTConfig: r.jwtConfig,
		Provider:  r.provider,
	}

	accessToken, accessExpiresAt, err := req.AccessJWT()
	if err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	data := &userApp.AccessToken{
		Access: &userApp.Token{
			Token:           accessToken,
			RefreshInterval: r.jwtConfig.AccessExpiresIn.String(),
			ExpiresAt:       accessExpiresAt,
		},
	}

	return &userApp.RefreshAccessTokenResult{Data: data}, nil
}

func (s *UserService) NewRefreshToken(res http.ResponseWriter, req *http.Request) {
	id := req.Header.Get("UserID")
	provider := req.Header.Get("Provider")

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
		provider:      provider,
	}

	userId, err := ParseStringToFloat(id)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken := req.Header.Get("Authorization")
	user, err := s.validateRefreshToken(int(userId), refreshToken)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := s.newRequest(user, provider).newRefreshToken()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(result); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

}

func (r *request) newRefreshToken() (*userApp.NewRefreshTokenResult, error) {

	req := token.Request{
		User:      r.user,
		JWTConfig: r.jwtConfig,
		Provider:  r.provider,
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

	data := &userApp.RefreshToken{
		Refresh: &userApp.Token{
			Token:           refreshToken,
			RefreshInterval: r.jwtConfig.RefreshExpiresIn.String(),
			ExpiresAt:       refreshExpiresAt,
		},
	}

	return &userApp.NewRefreshTokenResult{Data: data}, nil
}

func ParseStringToFloat(str string) (float64, error) {
	val, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return val, nil
	}

	//If user id is specifed in scientific notation
	pos := strings.IndexAny(str, "eE")
	if pos < 0 {
		return strconv.ParseFloat(str, 64)
	}

	var baseVal float64
	var expVal int64

	baseStr := str[0:pos]
	baseVal, err = strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return 0, err
	}

	expStr := str[(pos + 1):]
	expVal, err = strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return baseVal * math.Pow10(int(expVal)), nil
}
