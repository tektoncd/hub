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
	"github.com/gorilla/mux"
	"github.com/tektoncd/hub/api/pkg/app"
	user "github.com/tektoncd/hub/api/pkg/user/service"
)

func User(r *mux.Router, api app.Config) {

	userSvc := user.New(api)
	s := r.PathPrefix("/user").Subrouter()

	jwt := user.UserService{
		JwtConfig: api.JWTConfig(),
	}

	s.HandleFunc("/info", jwt.JWTAuth(userSvc.Info))
	s.HandleFunc("/refresh/accesstoken", jwt.JWTAuth(userSvc.RefreshAccessToken))
}
