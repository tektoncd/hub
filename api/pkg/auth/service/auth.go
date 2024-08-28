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
	"os"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
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
	log           *app.Logger
	defaultScopes []string
	jwtConfig     *app.JWTConfig
	provider      string
}

type AuthService struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Services struct {
	Service AuthService `json:"services"`
}

var (
	UI_URL   string
	provider string
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

	var log app.Logger
	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(authSvc); err != nil {
		log.Error(err)
	}
}

// Authenticates user with the speicified git provider
// using goth and calls the AuthCallback function
func Authenticate(res http.ResponseWriter, req *http.Request) {
	UI_URL = req.FormValue("redirect_uri")
	provider = mux.Vars(req)["provider"]

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
		provider:      provider,
	}

	ghUser, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	params := req.URL.Query()

	if err = r.insertData(ghUser, params.Get("code"), provider); err != nil {
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

	var gitUser model.User
	// Check if user exist
	q := r.db.Model(&model.User{}).
		Where("code = ?", code)

	err := q.First(&gitUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		} else {
			r.log.Error(err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Once the user is authenticated clear the code from DB and user struct so that it can't be reused once the user logs in
	gitUser.Code = ""
	if err := r.db.Model(&model.User{}).Where("email = ?", gitUser.Email).Update("code", gitUser.Code).Error; err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var acc model.Account
	accountQuery := r.db.Model(&model.Account{}).Where(model.Account{UserID: gitUser.ID, Provider: provider})

	err = accountQuery.First(&acc).Error
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// gets user scopes to add in jwt
	scopes, err := r.userScopes(&acc)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	userTokens, err := r.createTokens(&gitUser, scopes, provider)
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

	providerList := make([]authApp.Provider, 0)

	if os.Getenv("GH_CLIENT_ID") != "" && os.Getenv("GH_CLIENT_SECRET") != "" {
		providerList = append(providerList, authApp.Provider{Name: "github"})
	}

	if os.Getenv("BB_CLIENT_ID") != "" && os.Getenv("BB_CLIENT_SECRET") != "" {
		providerList = append(providerList, authApp.Provider{Name: "bitbucket"})
	}

	if os.Getenv("GL_CLIENT_ID") != "" && os.Getenv("GL_CLIENT_SECRET") != "" {
		providerList = append(providerList, authApp.Provider{Name: "gitlab"})
	}

	providers := authApp.ProviderList{
		Data: providerList,
	}

	var log app.Logger
	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(providers); err != nil {
		log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
