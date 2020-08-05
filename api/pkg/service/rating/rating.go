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

	"github.com/tektoncd/hub/api/gen/rating"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/service/auth"
)

type service struct {
	*auth.Validator
	logger *zap.SugaredLogger
	db     *gorm.DB
}

var (
	fetchError    = rating.MakeInternalError(fmt.Errorf("failed to fetch rating"))
	notFoundError = rating.MakeNotFound(fmt.Errorf("resource not found"))
)

// New returns the rating service implementation.
func New(api app.Config) rating.Service {
	return &service{
		Validator: &auth.Validator{
			DB:     api.DB(),
			Logger: api.Logger().With("service", "validator"),
			JWTKey: api.JWTSigningKey(),
		},
		logger: api.Logger().With("service", "rating"),
		db:     api.DB(),
	}
}

// Find user's rating for a resource
func (s *service) Get(ctx context.Context, p *rating.GetPayload) (*rating.GetResult, error) {

	user, err := s.Validator.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var res model.Resource
	if err := s.db.First(&res, p.ID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, notFoundError
		}
		s.logger.Error(err)
		return nil, fetchError
	}

	q := s.db.Where("user_id = ? AND resource_id = ?", user.ID, p.ID)

	r := model.UserResourceRating{}
	if err := q.Find(&r).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &rating.GetResult{Rating: -1}, nil
		}

		s.logger.Error(err)
		return nil, fetchError
	}

	return &rating.GetResult{Rating: int(r.Rating)}, nil
}
