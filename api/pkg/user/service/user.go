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
	"net/http"
	"strconv"

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	userApp "github.com/tektoncd/hub/api/pkg/user/app"
	"gorm.io/gorm"
)

type request struct {
	db            *gorm.DB
	log           *log.Logger
	user          *model.User
	defaultScopes []string
	jwtConfig     *app.JWTConfig
}

type UserService struct {
	app.Service
	api       app.Config
	JwtConfig *app.JWTConfig
}

type Service interface {
	Info(res http.ResponseWriter, req *http.Request)
}

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

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := r.User(userId)
	if err != nil {
		r.log.Error(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	result := userApp.InfoResult{
		Data: &userApp.UserData{
			GithubID:  user.GithubLogin,
			Name:      user.GithubName,
			AvatarURL: user.AvatarURL,
		},
	}

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(result); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}
