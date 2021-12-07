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

package token

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
)

const (
	issuer           = "Tekton Hub"
	accessTokenType  = "access-token"
	refreshTokenType = "refresh-token"
	agentTokenType   = "agent-token"
)

type Request struct {
	User      *model.User
	Scopes    []string
	JWTConfig *app.JWTConfig
	Provider  string
}

// current time
var Now = time.Now

func (r *Request) AccessJWT() (string, int64, error) {

	expiresAt := Now().Add(r.JWTConfig.AccessExpiresIn).Unix()
	claim := jwt.MapClaims{
		"iss":      issuer,
		"id":       r.User.ID,
		"provider": r.Provider,
		"scopes":   r.Scopes,
		"type":     accessTokenType,
		"iat":      Now().Unix(),
		"exp":      expiresAt,
	}

	token, err := Create(claim, r.JWTConfig.SigningKey)
	if err != nil {
		return "", 0, err
	}

	return token, expiresAt, nil
}

func (r *Request) RefreshJWT() (string, int64, error) {

	expiresAt := Now().Add(r.JWTConfig.RefreshExpiresIn).Unix()
	claim := jwt.MapClaims{
		"iss":      issuer,
		"id":       r.User.ID,
		"provider": r.Provider,
		"type":     refreshTokenType,
		"scopes":   []string{"refresh:token"},
		"iat":      Now().Unix(),
		"exp":      expiresAt,
	}

	token, err := Create(claim, r.JWTConfig.SigningKey)
	if err != nil {
		return "", 0, err
	}

	return token, expiresAt, nil
}

func (r *Request) AgentJWT() (string, error) {

	claim := jwt.MapClaims{
		"iss":    issuer,
		"id":     r.User.ID,
		"name":   r.User.AgentName,
		"type":   agentTokenType,
		"scopes": r.Scopes,
		"iat":    Now().Unix(),
	}

	token, err := Create(claim, r.JWTConfig.SigningKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
