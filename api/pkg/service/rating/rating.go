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

package rating

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"goa.design/goa/v3/security"

	"github.com/tektoncd/hub/api/gen/rating"
	"github.com/tektoncd/hub/api/pkg/app"
)

type service struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
	jwtKey string
}

// New returns the rating service implementation.
func New(api app.Config) rating.Service {
	return &service{api.Logger(), api.DB(), api.JWTSigningKey()}
}

// JWTAuth implements the authorization logic for service "rating" for the
// "jwt" security scheme.
func (s *service) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	return ctx, fmt.Errorf("not implemented")
}

// Find user's rating for a resource
func (s *service) Get(ctx context.Context, p *rating.GetPayload) (res *rating.GetResult, err error) {
	res = &rating.GetResult{}
	s.logger.Info("rating.Get")
	return
}
