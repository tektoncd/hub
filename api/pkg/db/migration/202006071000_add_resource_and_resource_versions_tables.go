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

package migration

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gopkg.in/gormigrate.v1"

	"github.com/tektoncd/hub/api/pkg/db/model"
)

func addResourceAndResourceVersionTable(log *zap.SugaredLogger) *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "202006071000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.Tag{}, &model.Catalog{},
				&model.Resource{}, &model.ResourceVersion{}).Error; err != nil {
				log.Error(err)
				return err
			}

			if err := fkey(log, tx, model.Resource{}, "catalog_id", "catalogs"); err != nil {
				return err
			}
			if err := fkey(log, tx, model.ResourceVersion{}, "resource_id", "resources"); err != nil {
				return err
			}
			if err := fkey(log, tx, model.ResourceTag{},
				"resource_id", "resources",
				"tag_id", "tags"); err != nil {
				return err
			}
			return nil
		},
	}
}
