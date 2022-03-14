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
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/token"
	"gorm.io/gorm"
)

type JWTScheme struct {
	// Name is the scheme name defined in the design.
	Name string
	// Scopes holds a list of scopes for the scheme.
	Scopes []string
	// RequiredScopes holds a list of scopes which are required
	// by the scheme. It is a subset of Scopes field.
	RequiredScopes []string
}

// JWTAuth acts as middleware and implements the authorization logic for services.
func (s *UserService) JWTAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		accessCookie, err1 := req.Cookie("accessToken")
		refreshCookie, err2 := req.Cookie("refreshToken")

		if err1 == http.ErrNoCookie && err2 == http.ErrNoCookie {
			http.Error(res, err1.Error(), http.StatusUnauthorized)
			return
		}

		var defalutJWT string
		if accessCookie == nil {
			defalutJWT = refreshCookie.Value
		} else {
			defalutJWT = accessCookie.Value
		}

		var jwt string
		switch req.RequestURI {
		case "/user/info":
			jwt = accessCookie.Value
		case "/user/refresh/accesstoken":
			jwt = refreshCookie.Value
		case "/user/refresh/refreshtoken":
			jwt = refreshCookie.Value
		case "/user/accesstoken", "/user/logout":
			jwt = accessCookie.Value
		default:
			jwt = defalutJWT
		}

		claims, err := token.Verify(jwt, s.JwtConfig.SigningKey)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		scheme := &JWTScheme{
			Name:   jwt,
			Scopes: []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
		}

		if req.RequestURI == "/user/info" {
			scheme.RequiredScopes = []string{"rating:read", "rating:write"}
		} else if req.RequestURI == "/refresh/accesstoken" || req.RequestURI == "/user/logout" || req.RequestURI == "/refresh/refreshtoken" {
			scheme.RequiredScopes = []string{"rating:read", "rating:write", "refresh:token"}
		}

		err = ValidateScopes(claims, scheme)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		provider, ok := claims["provider"].(string)
		if !ok {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		// Set the userId in the headers
		req.Header.Set("UserId", fmt.Sprintf("%v", userID))
		req.Header.Set("Provider", fmt.Sprintf("%v", provider))

		handler.ServeHTTP(res, req)
	}
}

// ValidateScopes takes user scopes and checks if it has the scope which
// is required for accessing the api
func ValidateScopes(claims jwt.MapClaims, scheme *JWTScheme) error {

	if claims["scopes"] == nil {
		return fmt.Errorf("invalid scopes")
	}

	scopes, ok := claims["scopes"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid scopes")
	}

	scopesInToken := make([]string, len(scopes))
	for _, scp := range scopes {
		scopesInToken = append(scopesInToken, scp.(string))
	}

	if err := scheme.Validate(scopesInToken); err != nil {
		return err
	}

	return nil
}

// Validate returns a non-nil error if scopes does not contain all of
// JWT scheme's required scopes.
func (s *JWTScheme) Validate(scopes []string) error {
	return validateScopes(s.RequiredScopes, scopes)
}

func validateScopes(expected, actual []string) error {
	var missing []string
	for _, r := range expected {
		found := false
		for _, s := range actual {
			if s == r {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, r)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing scopes: %s", strings.Join(missing, ", "))
}

func (r *request) User(id int) (*model.User, error) {

	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warnf("user not found for token: %s", err.Error())
			return nil, err
		}

		r.log.Errorf("error when looking up user. err: %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *request) GitUser(id int) (*model.Account, error) {

	var account model.Account
	if err := r.db.Where(model.Account{UserID: uint(id), Provider: r.provider}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Warnf("user not found for token: %s", err.Error())
			return nil, err
		}

		r.log.Errorf("error when looking up user. err: %s", err.Error())
		return nil, err
	}

	return &account, nil
}

func (s *UserService) newRequest(user *model.User, provider string) *request {
	return &request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		user:          user,
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
		provider:      provider,
	}
}

func (s *UserService) validateRefreshToken(id int, token string) (*model.User, error) {

	r := request{
		db:            s.DB(context.Background()),
		log:           s.Logger(context.Background()),
		defaultScopes: s.api.Data().Default.Scopes,
		jwtConfig:     s.api.JWTConfig(),
	}

	user, err := r.User(id)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	if user.RefreshTokenChecksum != createChecksum(token) {
		return nil, invalidRefreshToken
	}

	return user, nil
}

func createChecksum(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *request) userScopes() ([]string, error) {

	var userScopes []string = r.defaultScopes

	q := r.db.Preload("Scopes").Where(&model.User{Email: r.user.Email})

	dbUser := &model.User{}
	if err := q.Find(dbUser).Error; err != nil {
		return nil, refreshError
	}

	for _, s := range dbUser.Scopes {
		userScopes = append(userScopes, s.Name)
	}

	return userScopes, nil
}
