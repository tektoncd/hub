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

	"github.com/tektoncd/hub/api/gen/log"
	"github.com/tektoncd/hub/api/gen/rating"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/pkg/service/auth"
)

var (
	fetchError    = rating.MakeInternalError(fmt.Errorf("failed to fetch rating"))
	updateError   = rating.MakeInternalError(fmt.Errorf("failed to update rating"))
	notFoundError = rating.MakeNotFound(fmt.Errorf("resource not found"))
)

type service struct {
	*auth.Service
}

type request struct {
	db   *gorm.DB
	log  *log.Logger
	user *model.User
}

// New returns the rating service implementation.
func New(api app.Config) rating.Service {
	return &service{auth.NewService(api, "rating")}
}

// Find user's rating for a resource
func (s *service) Get(ctx context.Context, p *rating.GetPayload) (*rating.GetResult, error) {

	user, err := s.User(ctx)
	if err != nil {
		return nil, err
	}

	req := request{
		db:   s.DB(ctx),
		log:  s.Logger(ctx),
		user: user,
	}

	return req.getRating(p.ID)
}

// Update user's rating for a resource
func (s *service) Update(ctx context.Context, p *rating.UpdatePayload) error {

	user, err := s.User(ctx)
	if err != nil {
		return err
	}

	req := request{
		db:   s.DB(ctx),
		log:  s.Logger(ctx),
		user: user,
	}

	return req.updateRating(p.ID, p.Rating)
}

func (r *request) getRating(resID uint) (*rating.GetResult, error) {

	_, err := r.validateResourceID(resID)
	if err != nil {
		return nil, err
	}

	q := r.db.Where(&model.UserResourceRating{UserID: r.user.ID, ResourceID: resID})

	userRating := model.UserResourceRating{}
	if err := q.Find(&userRating).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &rating.GetResult{Rating: -1}, nil
		}
		r.log.Error(err)
		return nil, fetchError
	}

	return &rating.GetResult{Rating: int(userRating.Rating)}, nil
}

func (r *request) updateRating(resID, userRating uint) error {

	res, err := r.validateResourceID(resID)
	if err != nil {
		return err
	}

	if err := r.updateUserRating(resID, userRating); err != nil {
		return err
	}

	return r.updateResourceRating(res)
}

func (r *request) updateUserRating(resID, rating uint) error {

	q := r.db.Where(&model.UserResourceRating{UserID: r.user.ID, ResourceID: resID})

	rat := &model.UserResourceRating{}
	if err := q.FirstOrInit(rat).Error; err != nil {
		r.log.Error(err)
		return updateError
	}

	rat.Rating = rating
	if err := r.db.Save(rat).Error; err != nil {
		r.log.Error(err)
		return updateError
	}

	return nil
}

func (r *request) updateResourceRating(res *model.Resource) error {

	q := r.db.Model(&model.UserResourceRating{}).
		Where(&model.UserResourceRating{ResourceID: res.ID}).
		Select("ROUND(AVG(rating),1)")

	var avg float64
	if err := q.Row().Scan(&avg); err != nil {
		r.log.Error(err)
		return updateError
	}

	res.Rating = avg
	if err := r.db.Save(res).Error; err != nil {
		r.log.Error(err)
		return updateError
	}

	return nil
}

func (r *request) validateResourceID(id uint) (*model.Resource, error) {

	res := &model.Resource{}
	if err := r.db.First(res, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, notFoundError
		}
		r.log.Error(err)
		return nil, fetchError
	}

	return res, nil
}
