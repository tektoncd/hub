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
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/tektoncd/hub/api/pkg/app"
	"github.com/tektoncd/hub/api/pkg/db/model"
	"gorm.io/gorm"
)

// Migrate create tables and populates master tables
func Migrate(api *app.APIBase) error {

	log := api.Logger("migration")

	// NOTE: when writing a migration for a new table, add the same in InitSchema
	migration := gormigrate.New(
		api.DB(),
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			renameNameColumnToAgentNameInUserTable(log),
			createConfigTable(log),
			addRefreshTokenChecksumColumnInUserTable(log),
			updateCatalogBranchToMain(log),
		},
	)

	migration.InitSchema(func(db *gorm.DB) error {
		if err := db.AutoMigrate(
			&model.Category{},
			&model.Tag{},
			&model.Catalog{},
			&model.CatalogError{},
			&model.Resource{},
			&model.ResourceVersion{},
			&model.User{},
			&model.UserResourceRating{},
			&model.SyncJob{},
			&model.Scope{},
			&model.Config{},
		); err != nil {
			log.Error(err)
			return err
		}
		return nil
	})

	if err := migration.Migrate(); err != nil {
		log.Error(err, " failed to migrate")
		return err
	}

	log.Info("Migration ran successfully !!")
	return nil
}
