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

package category

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"go.uber.org/zap"

	"github.com/tektoncd/hub/api/gen/category"
	"github.com/tektoncd/hub/api/pkg/app"
)

type service struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

var (
	fetchError = category.MakeInternalError(fmt.Errorf("Failed to fetch categories"))
)

// New returns the category service implementation.
func New(api app.BaseConfig) category.Service {
	return &service{api.Logger(), api.DB()}
}

// List all categories along with their tags sorted by name
func (s *service) List(ctx context.Context) (res []*category.Category, err error) {
	var all []model.Category
	if err := s.db.Order("name").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.name ASC")
		}).
		Find(&all).Error; err != nil {
		s.logger.Error(err)
		return nil, fetchError
	}

	for _, c := range all {
		tags := []*category.Tag{}
		for _, t := range c.Tags {
			tags = append(tags, &category.Tag{ID: t.ID, Name: t.Name})
		}
		res = append(res, &category.Category{ID: c.ID, Name: c.Name, Tags: tags})
	}

	return res, nil
}
