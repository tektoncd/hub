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

package catalog

import (
	"context"
	"fmt"

	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"github.com/tektoncd/hub/api/v1/gen/catalog"
)

type service struct {
	app.Service
}

var (
	internalError = catalog.MakeInternalError(fmt.Errorf("failed to list all catalogs"))
)

// New returns the resource service implementation.
func New(api app.BaseConfig) catalog.Service {
	return &service{api.Service("catalog")}
}

// List of catalogs
func (s *service) List(ctx context.Context) (*catalog.ListResult, error) {

	log := s.Logger(ctx)
	db := s.DB(ctx)

	var all []model.Catalog
	if err := db.Order("id").Find(&all).Error; err != nil {
		log.Error(err)
		return nil, internalError
	}

	res := &catalog.ListResult{
		Data: []*catalog.Catalog{},
	}

	for _, c := range all {
		res.Data = append(res.Data,
			&catalog.Catalog{
				ID:       c.ID,
				Name:     c.Name,
				Type:     c.Type,
				URL:      c.URL,
				Provider: c.Provider,
			})
	}

	return res, nil
}
