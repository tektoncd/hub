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

package provider

import (
	"fmt"
	"os"

	"github.com/markbates/goth/providers/github"
)

type provider struct {
	Url          string
	ProfileUrl   string
	EmailUrl     string
	AuthUrl      string
	TokenUrl     string
	ClientId     string
	ClientSecret string
	CallbackUrl  string
}

func GithubProvider(AUTH_URL string) provider {
	githubAuth := provider{
		Url:          "https://github.com",
		ProfileUrl:   github.ProfileURL,
		EmailUrl:     github.EmailURL,
		ClientId:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		CallbackUrl:  fmt.Sprintf(AUTH_URL, "github"),
	}

	if os.Getenv("GHE_URL") != "" {
		githubAuth.Url = os.Getenv("GHE_URL")
		githubAuth.ProfileUrl = fmt.Sprintf("%s/api/v3/user", githubAuth.Url)
		githubAuth.EmailUrl = fmt.Sprintf("%s/api/v3/user/emails", githubAuth.Url)
	}

	githubAuth.AuthUrl = fmt.Sprintf("%s/login/oauth/authorize", githubAuth.Url)
	githubAuth.TokenUrl = fmt.Sprintf("%s/login/oauth/access_token", githubAuth.Url)

	return githubAuth
}
