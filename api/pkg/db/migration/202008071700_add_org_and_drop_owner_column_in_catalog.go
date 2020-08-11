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

func addOrgAndDropOwnerColumnInCatalog(log *zap.SugaredLogger) *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "202008071700",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.Model(&model.Catalog{}).
				DropColumn("owner").Error; err != nil {
				log.Error(err)
				return err
			}
			if err := tx.AutoMigrate(
				&model.Catalog{}).Error; err != nil {
				log.Error(err)
				return err
			}
			if err := tx.Model(&model.Catalog{}).
				AddUniqueIndex("uix_name_org", "name", "org").Error; err != nil {
				log.Error(err)
				return err
			}

			catalogQuery := `ALTER TABLE catalogs
					ALTER COLUMN type SET NOT NULL,
					ALTER COLUMN url  SET NOT NULL,
					ALTER COLUMN revision SET NOT NULL`
			if err := tx.Exec(catalogQuery).Error; err != nil {
				log.Error(err)
				return err
			}

			// update existing record
			if err := tx.Model(&model.Catalog{}).
				Updates(map[string]interface{}{"name": "catalog", "org": "tektoncd"}).Error; err != nil {
				log.Error(err)
				return err
			}

			return nil
		},
	}
}
