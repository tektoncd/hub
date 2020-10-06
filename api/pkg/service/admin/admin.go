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

package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	"github.com/tektoncd/hub/api/gen/admin"
	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/service/auth"
	"github.com/tektoncd/hub/api/pkg/token"
)

type service struct {
	*auth.Service
	jwtSigningKey string
}

type agentRequest struct {
	db            *gorm.DB
	log           *log.Logger
	jwtSigningKey string
}

var (
	invalidTokenError  = admin.MakeInvalidToken(fmt.Errorf("invalid user token"))
	invalidScopesError = admin.MakeInvalidScopes(fmt.Errorf("user not authorized"))
	internalError      = admin.MakeInternalError(fmt.Errorf("failed to create agent"))
)

// New returns the admin service implementation.
func New(api app.Config) admin.Service {
	return &service{
		Service:       auth.NewService(api, "admin"),
		jwtSigningKey: api.JWTSigningKey(),
	}
}

// Create or Update an agent user with required scopes
func (s *service) UpdateAgent(ctx context.Context, p *admin.UpdateAgentPayload) (*admin.UpdateAgentResult, error) {

	user, err := s.User(ctx)
	if err != nil {
		return nil, err
	}

	req := agentRequest{
		db:            s.DB(ctx),
		log:           s.LoggerWith(ctx, "user-id", user.ID),
		jwtSigningKey: s.jwtSigningKey,
	}

	return req.run(p.Name, p.Scopes)
}

func (r *agentRequest) run(name string, scopes []string) (*admin.UpdateAgentResult, error) {

	r.db = r.db.Begin()

	token, err := r.updateAgent(name, scopes)
	if err != nil {
		if err := r.db.Rollback().Error; err != nil {
			r.log.Error(err)
			return nil, internalError
		}
		return nil, err
	}

	if err := r.db.Commit().Error; err != nil {
		r.log.Error(err)
		return nil, internalError
	}
	return &admin.UpdateAgentResult{Token: token}, nil
}

func (r *agentRequest) updateAgent(name string, scopes []string) (string, error) {

	// Check if a normal user exists with the agent name in payload
	if err := r.userExistWithAgentName(name); err != nil {
		return "", err
	}

	// Check if an agent already exist with the name
	q := r.db.Model(&model.User{}).
		Where(&model.User{AgentName: name, Type: model.AgentUserType})

	agent := &model.User{}
	if err := q.First(&agent).Error; err != nil {

		// If agent does not exist then create one
		if gorm.IsRecordNotFoundError(err) {
			return r.addNewAgent(name, scopes)
		}
		r.log.Error(err)
		return "", internalError
	}

	// If an agent with name already exist, then update the scopes of agent
	return r.updateExistingAgent(agent, scopes)
}

func (r *agentRequest) addNewAgent(name string, scopes []string) (string, error) {

	agent := &model.User{
		AgentName: name,
		Type:      model.AgentUserType,
	}
	if err := r.db.Create(agent).Error; err != nil {
		r.log.Error(err)
		return "", internalError
	}

	if err := r.addScopesForAgent(agent, scopes); err != nil {
		return "", err
	}

	return r.createJWT(agent, scopes)
}

func (r *agentRequest) updateExistingAgent(agent *model.User, scopes []string) (string, error) {

	// Delete all existing scopes of agent
	if err := r.db.Where(&model.UserScope{UserID: agent.ID}).
		Delete(&model.UserScope{}).Error; err != nil {
		r.log.Error(err)
		return "", internalError
	}

	// Add new scopes for agent
	if err := r.addScopesForAgent(agent, scopes); err != nil {
		return "", err
	}

	return r.createJWT(agent, scopes)
}

func (r *agentRequest) addScopesForAgent(agent *model.User, scopes []string) error {

	for _, sc := range scopes {

		scope := &model.Scope{}
		if err := r.db.Where(&model.Scope{Name: sc}).
			First(&scope).Error; err != nil {

			// If scope in payload does not exist then return
			if gorm.IsRecordNotFoundError(err) {
				return admin.MakeInvalidPayload(fmt.Errorf("scope does not exist: %s", sc))
			}
			r.log.Error(err)
			return internalError
		}

		// Add scopes for agent
		us := model.UserScope{UserID: agent.ID, ScopeID: scope.ID}
		if err := r.db.Create(us).Error; err != nil {
			r.log.Error(err)
			return internalError
		}
	}

	return nil
}

func (r *agentRequest) createJWT(user *model.User, scopes []string) (string, error) {

	claim := jwt.MapClaims{
		"id":     user.ID,
		"name":   user.AgentName,
		"type":   user.Type,
		"scopes": scopes,
	}

	token, err := token.Create(claim, r.jwtSigningKey)
	if err != nil {
		r.log.Error(err)
		return "", internalError
	}

	return token, nil
}

func (r *agentRequest) userExistWithAgentName(name string) error {

	user := &model.User{}
	q := r.db.Where("LOWER(github_name) = ?", strings.ToLower(name))

	if err := q.First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		r.log.Error(err)
		return internalError
	}

	return admin.MakeInvalidPayload(fmt.Errorf("user exists with name: %s", name))
}
