// Copyright © 2020 The Tekton Authors.
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
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"github.com/tektoncd/hub/api/gen/auth"
	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type service struct {
	app.Service
	api app.Config
}

type request struct {
	db            *gorm.DB
	log           *log.Logger
	oauth         *oauth2.Config
	defaultScopes []string
	jwtSigningKey string
}

var (
	invalidCode   = auth.MakeInvalidCode(fmt.Errorf("invalid authorization code"))
	internalError = auth.MakeInternalError(fmt.Errorf("failed to authenticate"))
)

// New returns the auth service implementation.
func New(api app.Config) auth.Service {
	return &service{
		Service: api.Service("auth"),
		api:     api,
	}
}

// Authenticates users against GitHub OAuth
func (s *service) Authenticate(ctx context.Context, p *auth.AuthenticatePayload) (*auth.AuthenticateResult, error) {

	req := request{
		db:            s.DB(ctx),
		log:           s.Logger(ctx),
		oauth:         s.api.OAuthConfig(),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtSigningKey: s.api.JWTSigningKey(),
	}

	return req.authenticate(p.Code)
}

func (r *request) authenticate(code string) (*auth.AuthenticateResult, error) {

	// gets access_token for user using authorization_code
	token, err := r.oauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, invalidCode
	}

	// gets user details from github using the access_token
	oauthClient := r.oauth.Client(oauth2.NoContext, token)
	ghClient := github.NewClient(oauthClient)
	ghUser, _, err := ghClient.Users.Get(oauth2.NoContext, "")
	if err != nil {
		r.log.Error(err)
		return nil, internalError
	}

	// adds user in db if not exist
	user, err := r.addUser(ghUser)
	if err != nil {
		return nil, err
	}

	// gets user scopes to add in jwt
	scopes, err := r.userScopes(user)
	if err != nil {
		return nil, err
	}

	// creates jwt using user details
	jwt, err := r.createJWT(user, scopes)
	if err != nil {
		return nil, err
	}

	return &auth.AuthenticateResult{Token: jwt}, nil
}

func (r *request) addUser(user *github.User) (*model.User, error) {

	q := r.db.Model(&model.User{}).Where(&model.User{GithubLogin: user.GetLogin()})

	newUser := model.User{
		GithubName:  user.GetName(),
		GithubLogin: user.GetLogin(),
		Type:        model.NormalUserType,
	}
	if err := q.FirstOrCreate(&newUser).Error; err != nil {
		r.log.Error(err)
		return nil, internalError
	}

	return &newUser, nil
}

func (r *request) userScopes(user *model.User) ([]string, error) {

	var userScopes []string = r.defaultScopes

	q := r.db.Preload("Scopes").Where(&model.User{GithubLogin: user.GithubLogin})

	dbUser := model.User{}
	if err := q.Find(&dbUser).Error; err != nil {
		r.log.Error(err)
		return nil, internalError
	}

	for _, s := range dbUser.Scopes {
		userScopes = append(userScopes, s.Name)
	}

	return userScopes, nil
}

func (r *request) createJWT(user *model.User, scopes []string) (string, error) {

	claim := jwt.MapClaims{
		"id":     user.ID,
		"login":  user.GithubLogin,
		"name":   user.GithubName,
		"scopes": scopes,
	}

	token, err := token.Create(claim, r.jwtSigningKey)
	if err != nil {
		r.log.Error(err)
		return "", internalError
	}

	return token, nil
}
