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
	"net/http"
	"strconv"

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	auth "github.com/tektoncd/hub/api/pkg/auth/service"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
	userApp "github.com/tektoncd/hub/api/pkg/user/app"
	"gorm.io/gorm"
)

type request struct {
	db            *gorm.DB
	log           *log.Logger
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
	GetAccessToken(res http.ResponseWriter, req *http.Request)
	Logout(res http.ResponseWriter, req *http.Request)
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

	userId, err := strconv.Atoi(id)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	gitUser, err := r.GitUser(userId)
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

	userId, err := strconv.Atoi(id)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, err := req.Cookie(auth.RefreshToken)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
	}
	refreshToken := cookie.Value

	user, err := s.validateRefreshToken(userId, refreshToken)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := s.newRequest(user, provider).refreshAccessToken(res, req)
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

func (r *request) refreshAccessToken(res http.ResponseWriter, req *http.Request) (*userApp.RefreshAccessTokenResult, error) {

	scopes, err := r.userScopes()
	if err != nil {
		return nil, err
	}

	request := token.Request{
		User:      r.user,
		Scopes:    scopes,
		JWTConfig: r.jwtConfig,
		Provider:  r.provider,
	}

	accessToken, accessExpiresAt, err := request.AccessJWT()
	if err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	http.SetCookie(res, &http.Cookie{
		Name:     auth.AccessToken,
		Value:    accessToken,
		MaxAge:   int(r.jwtConfig.AccessExpiresIn.Seconds()),
		Path:     "/",
		HttpOnly: true,
	})

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

	cookie, err := req.Cookie("refreshToken")
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
	}

	refreshToken := cookie.Value

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
		provider:      provider,
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := s.validateRefreshToken(userId, refreshToken)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// result, err := s.newRequest(user, provider).newRefreshToken()
	result, err := s.newRequest(user, provider).newRefreshToken(res, req)
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

func (r *request) newRefreshToken(res http.ResponseWriter, req *http.Request) (*userApp.NewRefreshTokenResult, error) {

	request := token.Request{
		User:      r.user,
		JWTConfig: r.jwtConfig,
		Provider:  r.provider,
	}

	refreshToken, refreshExpiresAt, err := request.RefreshJWT()
	if err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	err = r.db.Model(r.user).UpdateColumn("refresh_token_checksum", createChecksum(refreshToken)).Error
	if err != nil {
		r.log.Error(err)
		return nil, refreshError
	}

	http.SetCookie(res, &http.Cookie{
		Name:     auth.RefreshToken,
		Value:    refreshToken,
		MaxAge:   int(r.jwtConfig.RefreshExpiresIn.Seconds()),
		Path:     "/",
		HttpOnly: true,
	})

	data := &userApp.RefreshToken{
		Refresh: &userApp.Token{
			Token:           refreshToken,
			RefreshInterval: r.jwtConfig.RefreshExpiresIn.String(),
			ExpiresAt:       refreshExpiresAt,
		},
	}

	return &userApp.NewRefreshTokenResult{Data: data}, nil
}

func (s *UserService) GetAccessToken(res http.ResponseWriter, req *http.Request) {

	c, err := req.Cookie("accessToken")
	if err == http.ErrNoCookie {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken := c.Value

	result := userApp.ExitingAccessToken{Data: accessToken}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(result); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *UserService) Logout(res http.ResponseWriter, req *http.Request) {

	// Unset the cookie
	deleteCookie(res, auth.AccessToken)
	deleteCookie(res, auth.RefreshToken)

	result := userApp.ClearCookies{Data: true}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(result); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func deleteCookie(res http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(res, cookie)
}
