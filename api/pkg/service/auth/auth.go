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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

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
	jwtConfig     *app.JWTConfig
	gitConfig     *app.GitConfig
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
		jwtConfig:     s.api.JWTConfig(),
		gitConfig:     s.api.GitConfig(),
	}

	return req.authenticate(p.Code)
}

// Perform oauth request and get the user details from GitHub
func (r *request) authGithub(code string) (model.User, error) {

	// gets access_token for user using authorization_code
	token, err := r.oauth.Exchange(context.Background(), code)
	if err != nil {
		return model.User{}, invalidCode
	}

	// gets user details from github using the access_token
	oauthClient := r.oauth.Client(context.Background(), token)

	var ghClient *github.Client

	// check if the url is enterprise url and then create the
	// client accordingly
	if r.gitConfig.IsEnterprise {
		ghClient, err = github.NewEnterpriseClient(r.gitConfig.GhConfig.ApiUrl, r.gitConfig.GhConfig.UploadUrl, oauthClient)
		if err != nil {
			return model.User{}, err
		}
	} else {
		ghClient = github.NewClient(oauthClient)
	}

	ghUser, _, err := ghClient.Users.Get(context.Background(), "")
	if err != nil {
		r.log.Error(err)
		return model.User{}, internalError
	}

	return model.User{
		Name:        ghUser.GetName(),
		GitUsername: strings.ToLower(ghUser.GetLogin()),
		Type:        model.NormalUserType,
		AvatarURL:   ghUser.GetAvatarURL(),
	}, nil
}

func (r *request) authenticate(code string) (*auth.AuthenticateResult, error) {

	var remoteUser model.User
	var err error
	switch r.gitConfig.Provider {
	case "github":
		remoteUser, err = r.authGithub(code)
		if err != nil {
			return nil, err
		}
	}

	// adds user in db if not exist
	user, err := r.addUser(remoteUser)
	if err != nil {
		return nil, err
	}

	// gets user scopes to add in jwt
	scopes, err := r.userScopes(user)
	if err != nil {
		return nil, err
	}

	// creates tokens using user details
	return r.createTokens(user, scopes)
}

func (r *request) addUser(user model.User) (*model.User, error) {
	var dbUser model.User

	// Check if user exist
	q := r.db.Model(&model.User{}).
		Where("LOWER(git_username) = ?", user.GitUsername)
	err := q.First(&dbUser).Error
	if err != nil {
		// If user doesn't exist, create a new record
		if err == gorm.ErrRecordNotFound {

			err = r.db.Create(&user).Error
			if err != nil {
				r.log.Error(err)
				return nil, internalError
			}
			return &user, nil
		}
		r.log.Error(err)
		return nil, internalError
	}

	// User already exist, check if GitHub Name is empty
	// If Name is empty, then user is inserted through config.yaml
	// Update user with remaining details
	if dbUser.Name == "" {
		dbUser.Name = user.Name
		dbUser.Type = model.NormalUserType
	}
	// For existing user, check if URL is not added
	if dbUser.AvatarURL == "" {
		dbUser.AvatarURL = user.AvatarURL
		if err = r.db.Save(&dbUser).Error; err != nil {
			r.log.Error(err)
			return nil, err
		}
	}

	return &dbUser, nil
}

func (r *request) userScopes(user *model.User) ([]string, error) {

	var userScopes []string = r.defaultScopes

	q := r.db.Preload("Scopes").Where(&model.User{GitUsername: user.GitUsername})

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

func (r *request) createTokens(user *model.User, scopes []string) (*auth.AuthenticateResult, error) {

	req := token.Request{
		User:      user,
		Scopes:    scopes,
		JWTConfig: r.jwtConfig,
	}

	accessToken, accessExpiresAt, err := req.AccessJWT()
	if err != nil {
		r.log.Error(err)
		return nil, internalError
	}

	refreshToken, refreshExpiresAt, err := req.RefreshJWT()
	if err != nil {
		r.log.Error(err)
		return nil, internalError
	}

	user.RefreshTokenChecksum = createChecksum(refreshToken)

	if err = r.db.Save(user).Error; err != nil {
		r.log.Error(err)
		return nil, internalError
	}

	data := &auth.AuthTokens{
		Access: &auth.Token{
			Token:           accessToken,
			RefreshInterval: r.jwtConfig.AccessExpiresIn.String(),
			ExpiresAt:       accessExpiresAt,
		},
		Refresh: &auth.Token{
			Token:           refreshToken,
			RefreshInterval: r.jwtConfig.RefreshExpiresIn.String(),
			ExpiresAt:       refreshExpiresAt,
		},
	}

	return &auth.AuthenticateResult{Data: data}, nil
}

func createChecksum(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
