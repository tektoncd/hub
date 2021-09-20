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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	authApp "github.com/tektoncd/hub/api/pkg/auth/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"gorm.io/gorm"
)

type service struct {
	app.Service
	api app.Config
}

type request struct {
	db            *gorm.DB
	log           *log.Logger
	defaultScopes []string
	jwtConfig     *app.JWTConfig
}

type AuthService struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Services struct {
	Service AuthService `json:"services"`
}

var (
	UI_URL string
)

type Service interface {
	AuthCallBack(res http.ResponseWriter, req *http.Request)
	HubAuthenticate(res http.ResponseWriter, req *http.Request)
}

// New returns the auth service implementation.
func New(api app.Config) Service {
	return &service{
		Service: api.Service("auth"),
		api:     api,
	}
}

// Return name and status of the services
func Status(res http.ResponseWriter, req *http.Request) {

	authSvc := Services{
		AuthService{
			Name:   "auth",
			Status: "ok",
		},
	}

	var log log.Logger
	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(authSvc); err != nil {
		log.Error(err)
	}
}

// Authenticates user with the speicified git provider
// using goth and calls the AuthCallback function
func Authenticate(res http.ResponseWriter, req *http.Request) {
	UI_URL = req.FormValue("redirect_uri")
	gothic.BeginAuthHandler(res, req)
}

// Once user is authenticated, store the user details in db
// and redirect to UI with the status code and auth code
func (s *service) AuthCallBack(res http.ResponseWriter, req *http.Request) {

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
	}

	ghUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	params := req.URL.Query()

	if err = r.insertData(ghUser, params.Get("code")); err != nil {
		r.log.Error(err)
		res.Header().Set("Location", fmt.Sprintf("%s?status=%d", UI_URL, http.StatusBadRequest))
		res.WriteHeader(http.StatusTemporaryRedirect)
	}

	res.Header().Set("Location", fmt.Sprintf("%s?status=%d&code=%s", UI_URL, http.StatusOK, params.Get("code")))
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// Checks for auth code in the headers, validates the
// auth code and returns the jwt token for user
func (s *service) HubAuthenticate(res http.ResponseWriter, req *http.Request) {

	// Get the auth code from params
	code := req.FormValue("code")

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
	}

	var user model.User
	// Check if user exist
	q := r.db.Model(&model.User{}).
		Where("code = ?", code)

	err := q.First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Error(err)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		} else {
			r.log.Error(err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := r.db.Model(&model.User{}).Where("github_login = ?", user.GithubLogin).Update("code", "").Error; err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// gets user scopes to add in jwt
	scopes, err := r.userScopes(&user)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	userTokens, err := r.createTokens(&user, scopes)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(userTokens); err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Provides a list of git provider present in auth server
func List(res http.ResponseWriter, req *http.Request) {
	// TODO: The values of the provider can be configured dynamically
	providers := authApp.ProviderList{
		Data: []authApp.Provider{
			{
				Name: "github",
			},
		},
	}

	var log log.Logger
	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(providers); err != nil {
		log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
