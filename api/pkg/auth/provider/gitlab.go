/*
Copyright 2022 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.

You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"fmt"
	"os"

	"github.com/markbates/goth/providers/gitlab"
)

func GitlabProvider(AUTH_URL string) provider {
	gitlabAuth := provider{
		Url:          "https://gitlab.com",
		AuthUrl:      gitlab.AuthURL,
		TokenUrl:     gitlab.TokenURL,
		ProfileUrl:   gitlab.ProfileURL,
		ClientId:     os.Getenv("GL_CLIENT_ID"),
		ClientSecret: os.Getenv("GL_CLIENT_SECRET"),
		CallbackUrl:  fmt.Sprintf(AUTH_URL, "gitlab"),
	}

	if os.Getenv("GLE_URL") != "" {
		gitlabAuth.Url = os.Getenv("GLE_URL")
		gitlabAuth.AuthUrl = fmt.Sprintf("%s/oauth/authorize", gitlabAuth.Url)
		gitlabAuth.TokenUrl = fmt.Sprintf("%s/oauth/token", gitlabAuth.Url)
		gitlabAuth.ProfileUrl = fmt.Sprintf("%s/api/v3/user", gitlabAuth.Url)
	}

	return gitlabAuth
}
