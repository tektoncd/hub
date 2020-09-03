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
)

func addNotNullToResourceAndResourceVersion(log *zap.SugaredLogger) *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "202008121043_add_not_null_constraint_to_resource_and_resource_version",
		Migrate: func(db *gorm.DB) error {

			resourceQuery := `ALTER TABLE resources
						 ALTER COLUMN name set NOT NULL,
						 ALTER COLUMN kind set NOT NULL;`
			if err := db.Exec(resourceQuery).Error; err != nil {
				log.Error(err)
				return err
			}

			resourceVersionQuery := `ALTER TABLE resource_versions
						 ALTER COLUMN version set NOT NULL,
						 ALTER COLUMN description set NOT NULL,
						 ALTER COLUMN min_pipelines_version set NOT NULL;`
			if err := db.Exec(resourceVersionQuery).Error; err != nil {
				log.Error(err)
				return err
			}
			return nil
		},
	}
}
