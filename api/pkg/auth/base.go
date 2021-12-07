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
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/auth/provider"
	auth "github.com/tektoncd/hub/api/pkg/auth/service"
)

// Auth Provider provides routes for authentication
// and also defines git providers using goth
func AuthProvider(r *mux.Router, api app.Config) {

	key := ""            // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30 // 30 days
	isProd := true       // Set to false when not serving over https
	if api.Environment() != app.EnvMode("production") {
		isProd = false
	}

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	var AUTH_BASE_URL = os.Getenv("AUTH_BASE_URL")
	var AUTH_URL = AUTH_BASE_URL + "/auth/%s/callback"

	githubAuth := provider.GithubProvider(AUTH_URL)
	gitlabAuth := provider.GitlabProvider(AUTH_URL)
	bitbucketAuth := provider.BitbucketProvider(AUTH_URL)

	goth.UseProviders(
		github.NewCustomisedURL(
			githubAuth.ClientId,
			githubAuth.ClientSecret,
			githubAuth.CallbackUrl,
			githubAuth.AuthUrl,
			githubAuth.TokenUrl,
			githubAuth.ProfileUrl,
			githubAuth.EmailUrl,
			"user"),

		gitlab.NewCustomisedURL(
			gitlabAuth.ClientId,
			gitlabAuth.ClientSecret,
			gitlabAuth.CallbackUrl,
			gitlabAuth.AuthUrl,
			gitlabAuth.TokenUrl,
			gitlabAuth.ProfileUrl,
		),

		bitbucket.New(
			bitbucketAuth.ClientId,
			bitbucketAuth.ClientSecret,
			bitbucketAuth.CallbackUrl,
		),
	)

	authSvc := auth.New(api)

	// Return name and status of the services
	r.HandleFunc("/", auth.Status)

	s := r.PathPrefix("/auth").Subrouter()

	// Provides a list of git provider present in auth server
	s.HandleFunc("/providers", auth.List)

	// Checks for auth code, validates it and returns users jwt token
	s.HandleFunc("/login", authSvc.HubAuthenticate).Methods(http.MethodPost)

	// Redirects to UI with the status code and auth code
	s.HandleFunc("/{provider}/callback", authSvc.AuthCallBack)

	// Authenticates user with the speicified git provider
	s.HandleFunc("/{provider}", auth.Authenticate)
}
